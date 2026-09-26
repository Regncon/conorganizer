package puljefordeling

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestEmulateSeatings_ListsWhoGotAndLacksForstevalg(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt to billettholdere som vil ha den ene plassen på samme arrangement, og en tredje som bare er litt interessert.",
		When:  "Når puljefordelingen emuleres.",
		Then:  "Så skal puljen liste hvem som fikk førstevalg og hvem som fortsatt mangler det, med plass og om de ville ha førstevalg her.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_forstevalg_lists")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "ev1", "Drager", 1, fredag)
	seedEvent(t, db, "ev2", "Brett", 4, fredag)
	seedParticipant(t, db, 1, "Anne", "Aas")
	seedParticipant(t, db, 2, "Bjørn", "Berg")
	seedParticipant(t, db, 3, "Cato", "Carlsen")
	seedInterest(t, db, 1, "ev1", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 2, "ev1", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 3, "ev2", fredag, models.InterestLevelLow)

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	pulje := em.Puljer[0]
	if len(pulje.GotForstevalg) != 1 || pulje.GotForstevalg[0].EventTitle != "Drager" || pulje.GotForstevalg[0].Level != models.InterestLevelHigh {
		t.Fatalf("expected one participant with førstevalg in Drager, got %+v", pulje.GotForstevalg)
	}
	winner := pulje.GotForstevalg[0].Name
	loser := "Bjørn Berg"
	if winner == loser {
		loser = "Anne Aas"
	}
	if got := participantNames(pulje.WithoutForstevalg); !slices.Equal(got, []string{loser, "Cato Carlsen"}) {
		t.Fatalf("expected %s (wanted førstevalg here) first, then Cato, got %v", loser, got)
	}
	missed := pulje.WithoutForstevalg[0]
	if !missed.WantedForstevalg || missed.EventTitle != "" {
		t.Fatalf("expected %s to have wanted førstevalg and have no seat, got %+v", loser, missed)
	}
	cato := pulje.WithoutForstevalg[1]
	if cato.WantedForstevalg || cato.EventTitle != "Brett" || cato.Level != models.InterestLevelLow {
		t.Fatalf("expected Cato seated in Brett with low interest, got %+v", cato)
	}
}

func TestEmulateSeatings_WithoutForstevalgIsCumulativeAcrossPuljer(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder som får førstevalg i første pulje.",
		When:  "Når puljefordelingen emuleres over to puljer.",
		Then:  "Så skal hen ikke stå som uten førstevalg i den andre puljen.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_forstevalg_cumulative")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "ev1", "Drager", 4, models.PuljeFredagKveld)
	seedEvent(t, db, "ev2", "Brett", 4, models.PuljeLordagMorgen)
	seedParticipant(t, db, 1, "Anne", "Aas")
	seedInterest(t, db, 1, "ev1", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedInterest(t, db, 1, "ev2", models.PuljeLordagMorgen, models.InterestLevelLow)

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	if len(em.Puljer) != 2 || len(em.Puljer[1].WithoutForstevalg) != 0 || len(em.Puljer[1].GotForstevalg) != 0 {
		t.Fatalf("expected no one without førstevalg in the second pulje, got %+v", em.Puljer[1])
	}
}

func participantNames(deltakere []PuljeParticipant) []string {
	out := make([]string, 0, len(deltakere))
	for _, d := range deltakere {
		out = append(out, d.Name)
	}
	return out
}

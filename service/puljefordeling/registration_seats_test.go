package puljefordeling

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func seedRegistration(t *testing.T, db *sql.DB, bhID int, eventID string, pulje models.Pulje) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
	`, eventID, pulje, bhID, models.EventPlayerRolePlayer, models.EventPlayerSourceRegistration)
}

func TestEmulateSeatings_RegistrationIsShownAndExcludesOrdinarySelectionInPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en påmeldt billetthelder som også har en ordinær interesse i samme pulje.",
		When:  "Når puljefordelingen beregnes.",
		Then:  "Så skal påmeldingen vises, og billetthelderen skal ikke plasseres på den ordinære interessen.",
	})

	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_registration_exclusion")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evOpen", "Open", 100, fredag)
	seedEvent(t, db, "evOrdinary", "Ordinary", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedParticipant(t, db, 2, "Ola", "Nordmann")
	seedRegistration(t, db, 1, "evOpen", fredag)
	seedInterest(t, db, 1, "evOrdinary", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 2, "evOrdinary", fredag, models.InterestLevelHigh)

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	openEvent, ok := findEvent(em.Puljer[0], "evOpen")
	if !ok {
		t.Fatal("evOpen missing from result")
	}
	if !slices.Contains(playerNames(openEvent.AssignedPlayers), "Kari Nordmann") {
		t.Fatalf("registration should be shown on evOpen, got %v", playerNames(openEvent.AssignedPlayers))
	}
	if len(openEvent.AssignedPlayers) != 1 || !openEvent.AssignedPlayers[0].Registration {
		t.Fatalf("registration should be identified as confirmed registration, got %+v", openEvent.AssignedPlayers)
	}

	ordinaryEvent, ok := findEvent(em.Puljer[0], "evOrdinary")
	if !ok {
		t.Fatal("evOrdinary missing from result")
	}
	if slices.Contains(playerNames(ordinaryEvent.AssignedPlayers), "Kari Nordmann") {
		t.Fatalf("registered participant must be excluded from ordinary selection in the same pulje, got %v", playerNames(ordinaryEvent.AssignedPlayers))
	}
	if !slices.Contains(playerNames(ordinaryEvent.AssignedPlayers), "Ola Nordmann") {
		t.Fatalf("unregistered participant should remain eligible, got %v", playerNames(ordinaryEvent.AssignedPlayers))
	}
}

func TestEmulateSeatings_PreservesMultipleRegistrationsInOnePulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder som er påmeldt flere åpne arrangementer i samme pulje.",
		When:  "Når puljefordelingen beregnes.",
		Then:  "Så skal alle påmeldingene vises i den samme puljen.",
	})

	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_multiple_registrations")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evOpenA", "Open A", 100, fredag)
	seedEvent(t, db, "evOpenB", "Open B", 100, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedRegistration(t, db, 1, "evOpenA", fredag)
	seedRegistration(t, db, 1, "evOpenB", fredag)

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	for _, eventID := range []string{"evOpenA", "evOpenB"} {
		event, ok := findEvent(em.Puljer[0], eventID)
		if !ok {
			t.Fatalf("%s missing from result", eventID)
		}
		if len(event.AssignedPlayers) != 1 || event.AssignedPlayers[0].Name != "Kari Nordmann" || !event.AssignedPlayers[0].Registration {
			t.Fatalf("%s should show Kari's registration, got %+v", eventID, event.AssignedPlayers)
		}
	}
}

func TestEmulateSeatings_RegistrationOnlyExcludesItsOwnPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder som er påmeldt fredag og interessert i et arrangement lørdag.",
		When:  "Når begge puljene beregnes.",
		Then:  "Så skal fredagens påmelding ikke hindre ordinær plassering på lørdag.",
	})

	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_registration_pulje_scope")
	const fredag = models.PuljeFredagKveld
	const lordag = models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, lordag, "Lørdag Formiddag", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "evOpen", "Open", 100, fredag)
	seedEvent(t, db, "evSaturday", "Saturday", 4, lordag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedRegistration(t, db, 1, "evOpen", fredag)
	seedInterest(t, db, 1, "evSaturday", lordag, models.InterestLevelHigh)

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	saturdayEvent, ok := findEvent(em.Puljer[1], "evSaturday")
	if !ok {
		t.Fatal("evSaturday missing from result")
	}
	if !slices.Contains(playerNames(saturdayEvent.AssignedPlayers), "Kari Nordmann") {
		t.Fatalf("Friday registration must not exclude Saturday selection, got %v", playerNames(saturdayEvent.AssignedPlayers))
	}
}

package event

import (
	"slices"
	"strings"
	"testing"

	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetProgramInterestNotice_IncludesInterestsForAllTicketHoldersInPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at to tilknyttede billettholdere har meldt interesse for samme arrangement i en åpen pulje.",
		When:  "Når programadvarselen for puljen hentes.",
		Then:  "Så skal arrangementet nevnes én gang sammen med riktig pulje.",
	})

	// Given
	expectedNotice := ProgramInterestNotice{
		PuljeName:   "Fredag kveld",
		EventTitles: []string{"Interest Event"},
	}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)

	mustExecEventInterestTest(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (902, 'Second', 'Interest', 1, 'Ticket', 1, 7002, 8002)
	`)
	mustExecEventInterestTest(t, db, `
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (902, ?, ?, ?)
	`, fixture.eventID, fixture.puljeID, models.InterestLevelMedium)

	// When
	notice, err := getProgramInterestNotice(db, []ticketholder.BillettHolder{
		{Id: fixture.billettholderID},
		{Id: 902},
	}, string(fixture.puljeID))

	// Then
	if err != nil {
		t.Fatalf("get program interest notice: %v", err)
	}
	if notice.PuljeName != expectedNotice.PuljeName {
		t.Fatalf("pulje name = %q, want %q", notice.PuljeName, expectedNotice.PuljeName)
	}
	if !slices.Equal(notice.EventTitles, expectedNotice.EventTitles) {
		t.Fatalf("event titles = %v, want %v", notice.EventTitles, expectedNotice.EventTitles)
	}
}

func TestGetProgramInterestNotice_IgnoresProgramEvents(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at et programarrangement har en interesse i en åpen pulje.",
		When:  "Når programadvarselen for puljen hentes.",
		Then:  "Så skal programarrangementet ikke gi en raffle-advarsel.",
	})

	// Given
	expectedNotice := ProgramInterestNotice{}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	ticketHolders := []ticketholder.BillettHolder{{Id: fixture.billettholderID}}

	mustExecEventInterestTest(t, db, `UPDATE events SET is_in_puljefordeling = 0 WHERE id = ?`, fixture.eventID)

	// When
	notice, err := getProgramInterestNotice(db, ticketHolders, string(fixture.puljeID))

	// Then
	if err != nil {
		t.Fatalf("get program interest notice for program event: %v", err)
	}
	if !slices.Equal(notice.EventTitles, expectedNotice.EventTitles) {
		t.Fatalf("event titles = %v, want %v", notice.EventTitles, expectedNotice.EventTitles)
	}
}

func TestGetProgramInterestNotice_IgnoresCompletedPuljer(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har meldt interesse i en fullført pulje.",
		When:  "Når programadvarselen for puljen hentes.",
		Then:  "Så skal det ikke vises en advarsel om å bli valgt ut.",
	})

	// Given
	expectedNotice := ProgramInterestNotice{}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	ticketHolders := []ticketholder.BillettHolder{{Id: fixture.billettholderID}}

	mustExecEventInterestTest(t, db, `UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusCompleted, fixture.puljeID)

	// When
	notice, err := getProgramInterestNotice(db, ticketHolders, string(fixture.puljeID))

	// Then
	if err != nil {
		t.Fatalf("get program interest notice for completed pulje: %v", err)
	}
	if !slices.Equal(notice.EventTitles, expectedNotice.EventTitles) {
		t.Fatalf("event titles = %v, want %v", notice.EventTitles, expectedNotice.EventTitles)
	}
}

func TestProgramEventInterestPanel_RendersOpenMessageAndWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har interesse i ett arrangement i puljen.",
		When:  "Når informasjonspanelet for et programarrangement rendres.",
		Then:  "Så skal åpenhetsmeldingen og advarselen vises uten raffle-kontroller.",
	})

	// Given
	expectedMessageParts := []string{
		"åpent for alle",
		"reservere plass",
		"Dungeons & Dragons",
		"Lørdag kveld",
		"Se og endre interessene dine",
	}
	notice := ProgramInterestNotice{
		PuljeName:   "Lørdag kveld",
		EventTitles: []string{"Dungeons & Dragons"},
	}

	// When
	doc := templtest.Render(t, ProgramEventInterestPanel(notice))

	// Then
	message := doc.Find(".event-interest-program-message").Text()
	if !containsAll(message, expectedMessageParts...) {
		t.Fatalf("program interest message = %q", message)
	}
	if doc.Find(".event-interest-open-button, .interest-dialog").Length() != 0 {
		t.Fatal("program interest panel rendered raffle controls")
	}
}

func containsAll(value string, expected ...string) bool {
	for _, part := range expected {
		if !strings.Contains(value, part) {
			return false
		}
	}
	return true
}

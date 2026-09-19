package event

import (
	"strings"
	"testing"

	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetProgramInterestNotice_IncludesInterestsForAllTicketHoldersInPulje(t *testing.T) {
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

	notice, err := getProgramInterestNotice(db, []ticketholder.BillettHolder{
		{Id: fixture.billettholderID},
		{Id: 902},
	}, string(fixture.puljeID))
	if err != nil {
		t.Fatalf("get program interest notice: %v", err)
	}

	if notice.PuljeName != "Fredag kveld" {
		t.Fatalf("pulje name = %q, want %q", notice.PuljeName, "Fredag kveld")
	}
	if len(notice.EventTitles) != 1 || notice.EventTitles[0] != "Interest Event" {
		t.Fatalf("event titles = %v, want [Interest Event]", notice.EventTitles)
	}
}

func TestGetProgramInterestNotice_IgnoresProgramEventsAndCompletedPuljer(t *testing.T) {
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	ticketHolders := []ticketholder.BillettHolder{{Id: fixture.billettholderID}}

	mustExecEventInterestTest(t, db, `UPDATE events SET is_in_puljefordeling = 0 WHERE id = ?`, fixture.eventID)
	notice, err := getProgramInterestNotice(db, ticketHolders, string(fixture.puljeID))
	if err != nil {
		t.Fatalf("get program interest notice for program event: %v", err)
	}
	if notice.HasInterest() {
		t.Fatalf("program event interest should not produce a raffle warning: %+v", notice)
	}

	mustExecEventInterestTest(t, db, `UPDATE events SET is_in_puljefordeling = 1 WHERE id = ?`, fixture.eventID)
	mustExecEventInterestTest(t, db, `UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusCompleted, fixture.puljeID)
	notice, err = getProgramInterestNotice(db, ticketHolders, string(fixture.puljeID))
	if err != nil {
		t.Fatalf("get program interest notice for completed pulje: %v", err)
	}
	if notice.HasInterest() {
		t.Fatalf("completed pulje interest should not produce a raffle warning: %+v", notice)
	}
}

func TestProgramEventInterestPanel_RendersOpenMessageAndWarning(t *testing.T) {
	doc := templtest.Render(t, ProgramEventInterestPanel(ProgramInterestNotice{
		PuljeName:   "Lørdag kveld",
		EventTitles: []string{"Dungeons & Dragons"},
	}))

	message := doc.Find(".event-interest-program-message").Text()
	if !containsAll(message, "åpent for alle", "reservere plass", "Dungeons & Dragons", "Lørdag kveld", "Se og endre interessene dine") {
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

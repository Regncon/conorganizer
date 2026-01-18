package formsubmission

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

func TestGetInterestsForEvent_FirstChoiceRules(t *testing.T) {
	db, logger, createDBErr := testutil.CreateTemporaryDBAndLogger("test_first_choice", t)
	if createDBErr != nil {
		t.Fatalf("failed to create test database: %v", createDBErr)
	}
	defer db.Close()

	seedBaseTables(t, db)
	seedBillettholdere(t, db, []billettholderFixture{
		{id: 1, firstName: "Player", lastName: "One"},
		{id: 2, firstName: "GM", lastName: "Two"},
		{id: 3, firstName: "NotVery", lastName: "Three"},
		{id: 4, firstName: "NoAssign", lastName: "Four"},
		{id: 5, firstName: "SameEvent", lastName: "Five"},
		{id: 6, firstName: "GMPlayer", lastName: "Six"},
		{id: 7, firstName: "GMAndPlayer", lastName: "Seven"},
	})
	seedInterests(t, db, []interestFixture{
		{billettholderID: 1, eventID: "E2", puljeID: "P2", interestLevel: "Veldig interessert"},
		{billettholderID: 2, eventID: "E2", puljeID: "P2", interestLevel: "Veldig interessert"},
		{billettholderID: 3, eventID: "E2", puljeID: "P2", interestLevel: "Interessert"},
		{billettholderID: 4, eventID: "E2", puljeID: "P2", interestLevel: "Litt interessert"},
		{billettholderID: 5, eventID: "E2", puljeID: "P2", interestLevel: "Veldig interessert"},
		{billettholderID: 6, eventID: "E2", puljeID: "P2", interestLevel: "Veldig interessert"},
		{billettholderID: 7, eventID: "E2", puljeID: "P2", interestLevel: "Veldig interessert"},
		{billettholderID: 1, eventID: "E3", puljeID: "P3", interestLevel: "Veldig interessert"},
		{billettholderID: 2, eventID: "E3", puljeID: "P3", interestLevel: "Veldig interessert"},
		{billettholderID: 6, eventID: "E3", puljeID: "P3", interestLevel: "Veldig interessert"},
		{billettholderID: 1, eventID: "E4", puljeID: "P4", interestLevel: "Veldig interessert"},
		{billettholderID: 2, eventID: "E4", puljeID: "P4", interestLevel: "Veldig interessert"},
		{billettholderID: 4, eventID: "E4", puljeID: "P4", interestLevel: "Ikkje interessert"},
	})
	seedAssignments(t, db, []assignmentFixture{
		{eventID: "E1", puljeID: "P1", billettholderID: 1, isPlayer: 1, isGM: 0},
		{eventID: "E1", puljeID: "P1", billettholderID: 2, isPlayer: 0, isGM: 1},
		{eventID: "E1", puljeID: "P1", billettholderID: 3, isPlayer: 1, isGM: 0},
		{eventID: "E2", puljeID: "P2", billettholderID: 5, isPlayer: 1, isGM: 0},
		{eventID: "E1", puljeID: "P1", billettholderID: 6, isPlayer: 0, isGM: 1},
		{eventID: "E3", puljeID: "P3", billettholderID: 6, isPlayer: 1, isGM: 0},
		{eventID: "E1", puljeID: "P1", billettholderID: 7, isPlayer: 0, isGM: 1},
		{eventID: "E4", puljeID: "P4", billettholderID: 7, isPlayer: 1, isGM: 0},
	})

	interests, getInterestsErr := GetInterestsForEvent("E2", db, logger)
	if getInterestsErr != nil {
		t.Fatalf("GetInterestsForEvent error: %v", getInterestsErr)
	}

	got := indexInterests(interests)

	// E2 inclusion checks confirm interests are listed even if the player is already assigned
	// to the same event; assignment should only affect FirstChoice, not filtering.
	t.Run("E2 includes/excludes correct billettholders", func(t *testing.T) {
		sameEventAssigneeID := 5
		playerAssignedID := 1
		gmAssignedID := 2
		notVeryInterestedID := 3
		unassignedID := 4
		gmPlayerID := 6
		gmAndPlayerDifferentEventsID := 7

		if _, ok := got[sameEventAssigneeID]; !ok {
			t.Fatalf("expected assigned-to-same-event billettholder to be returned")
		}
		if _, ok := got[playerAssignedID]; !ok {
			t.Fatalf("expected player-assigned billettholder to be returned")
		}
		if _, ok := got[gmAssignedID]; !ok {
			t.Fatalf("expected gm-assigned billettholder to be returned")
		}
		if _, ok := got[notVeryInterestedID]; !ok {
			t.Fatalf("expected not-very-interested billettholder to be returned")
		}
		if _, ok := got[unassignedID]; !ok {
			t.Fatalf("expected unassigned billettholder to be returned")
		}
		if _, ok := got[gmPlayerID]; !ok {
			t.Fatalf("expected gm+player billettholder to be returned")
		}
		if _, ok := got[gmAndPlayerDifferentEventsID]; !ok {
			t.Fatalf("expected gm+player (different events) billettholder to be returned")
		}
	})

	// E2 first-choice checks focus on the CASE logic in queryFirstChoice:
	// - "Veldig interessert" + assigned as player in a different event => FirstChoice should be true.
	// - GM-only in a different event should NOT count as FirstChoice.
	// - Any interest below "Veldig interessert" should NOT be FirstChoice, even if assigned elsewhere.
	// - No assignment at all should NOT be FirstChoice.
	t.Run("E2 first-choice rules", func(t *testing.T) {
		playerAssignedID := 1
		gmAssignedID := 2
		notVeryInterestedID := 3
		unassignedID := 4
		gmPlayerID := 6
		gmAndPlayerDifferentEventsID := 7

		if got[playerAssignedID].FirstChoice != true {
			t.Errorf("player assigned to other event should be first choice")
		}
		if got[gmAssignedID].FirstChoice != false {
			t.Errorf("gm assigned to other event should not be first choice")
		}
		if got[notVeryInterestedID].FirstChoice != false {
			t.Errorf("not very interested should not be first choice")
		}
		if got[unassignedID].FirstChoice != false {
			t.Errorf("no assignment should not be first choice")
		}
		if got[gmPlayerID].FirstChoice != true {
			t.Errorf("gm+player with very interested should be first choice due to player assignment")
		}
		if got[gmAndPlayerDifferentEventsID].FirstChoice != true {
			t.Errorf("gm in one event and player in another should still be first choice")
		}
	})

	interestsE3, getInterestsE3Err := GetInterestsForEvent("E3", db, logger)
	if getInterestsE3Err != nil {
		t.Fatalf("GetInterestsForEvent E3 error: %v", getInterestsE3Err)
	}

	gotE3 := indexInterests(interestsE3)

	// E3 inclusion check confirms assignments to the same event do not filter interests out.
	t.Run("E3 includes/excludes correct billettholders", func(t *testing.T) {
		sameEventAssigneeID := 6
		if _, ok := gotE3[sameEventAssigneeID]; !ok {
			t.Fatalf("expected assigned-to-same-event billettholder to be returned for E3")
		}
	})

	// E3 first-choice checks re-run the same CASE rules against a different event to confirm
	// the logic is not accidentally tied to E2-only data setup.
	t.Run("E3 first-choice rules", func(t *testing.T) {
		playerAssignedID := 1
		gmAssignedID := 2

		if gotE3[playerAssignedID].FirstChoice != true {
			t.Errorf("player assigned to other event should be first choice for E3")
		}
		if gotE3[gmAssignedID].FirstChoice != false {
			t.Errorf("gm assigned to other event should not be first choice for E3")
		}
	})

	interestsE4, getInterestsE4Err := GetInterestsForEvent("E4", db, logger)
	if getInterestsE4Err != nil {
		t.Fatalf("GetInterestsForEvent E4 error: %v", getInterestsE4Err)
	}

	gotE4 := indexInterests(interestsE4)

	// E4 first-choice checks cover an interest mix with an explicit "no assignment" case to ensure
	// the FirstChoice flag remains false when the participant has no cross-event player assignment.
	t.Run("E4 first-choice rules", func(t *testing.T) {
		playerAssignedID := 1
		gmAssignedID := 2
		unassignedID := 4

		if gotE4[playerAssignedID].FirstChoice != true {
			t.Errorf("player assigned to other event should be first choice for E4")
		}
		if gotE4[gmAssignedID].FirstChoice != false {
			t.Errorf("gm assigned to other event should not be first choice for E4")
		}
		if gotE4[unassignedID].FirstChoice != false {
			t.Errorf("no assignment should not be first choice for E4")
		}
	})
}

type billettholderFixture struct {
	id        int
	firstName string
	lastName  string
}

type interestFixture struct {
	billettholderID int
	eventID         string
	puljeID         string
	interestLevel   string
}

type assignmentFixture struct {
	eventID         string
	puljeID         string
	billettholderID int
	isPlayer        int
	isGM            int
}

func indexInterests(interests []InterestWithHolder) map[int]InterestWithHolder {
	index := make(map[int]InterestWithHolder, len(interests))
	for _, interest := range interests {
		index[interest.BillettholderId] = interest
	}
	return index
}

func seedBaseTables(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExec(t, db, `INSERT INTO event_statuses(status) VALUES ('Godkjent')`)
	mustExec(t, db, `INSERT INTO events_types(event_type) VALUES ('Other')`)
	mustExec(t, db, `INSERT INTO age_groups(age_group) VALUES ('Default')`)
	mustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES ('Normal')`)
	mustExec(t, db, `INSERT INTO interest_levels(interest_level) VALUES ('Veldig interessert'), ('Interessert'), ('Litt interessert'), ('Ikkje interessert')`)
	mustExec(t, db, `
		INSERT INTO puljer (
			id, name, is_closed, is_published, start_time, end_time
		) VALUES
			('P1', 'Friday', 0, 1, '2025-10-03', '2025-10-03'),
			('P2', 'SaturdayMorning', 0, 1, '2025-10-04', '2025-10-04'),
			('P3', 'SaturdayEvening', 0, 1, '2025-10-04', '2025-10-04'),
			('P4', 'Sunday', 0, 1, '2025-10-05', '2025-10-05')
	`)
	mustExec(t, db, `
		INSERT INTO events (
			id, title, intro, description, image_url, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			pulje_name, max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES
			('E1','Event 1','intro','desc','', '', 'Other','Default','Normal','Host 1','h1@test.no','11111111','Friday',4,1,1,'Godkjent'),
			('E2','Event 2','intro','desc','', '', 'Other','Default','Normal','Host 2','h2@test.no','22222222','SaturdayMorning',4,1,1,'Godkjent'),
			('E3','Event 3','intro','desc','', '', 'Other','Default','Normal','Host 3','h3@test.no','33333333','SaturdayEvening',4,1,1,'Godkjent'),
			('E4','Event 4','intro','desc','', '', 'Other','Default','Normal','Host 4','h4@test.no','44444444','Sunday',4,1,1,'Godkjent')
	`)
}

func seedBillettholdere(t *testing.T, db *sql.DB, rows []billettholderFixture) {
	t.Helper()

	for _, row := range rows {
		orderID := 1000 + row.id
		ticketID := 2000 + row.id
		mustExec(t, db, `
			INSERT INTO billettholdere (
				id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
			) VALUES (?, ?, ?, 1, 'Test', 1, ?, ?)
		`, row.id, row.firstName, row.lastName, orderID, ticketID)
	}
}

func seedInterests(t *testing.T, db *sql.DB, rows []interestFixture) {
	t.Helper()

	for _, row := range rows {
		mustExec(t, db, `
			INSERT INTO interests (
				billettholder_id, event_id, pulje_id, interest_level
			) VALUES (?, ?, ?, ?)
		`, row.billettholderID, row.eventID, row.puljeID, row.interestLevel)
	}
}

func seedAssignments(t *testing.T, db *sql.DB, rows []assignmentFixture) {
	t.Helper()

	for _, row := range rows {
		mustExec(t, db, `
			INSERT INTO events_players (
				event_id, pulje_id, billettholder_id, is_player, is_gm
			) VALUES (?, ?, ?, ?, ?)
		`, row.eventID, row.puljeID, row.billettholderID, row.isPlayer, row.isGM)
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, execErr := db.Exec(query, args...); execErr != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", execErr, query)
	}
}

package formsubmission

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

const (
	eventE1 = "E1"
	eventE2 = "E2"
	eventE3 = "E3"
	eventE4 = "E4"

	puljeP1 = "P1"
	puljeP2 = "P2"
	puljeP3 = "P3"
	puljeP4 = "P4"

	idPlayerAssigned             = 1
	idGMAssigned                 = 2
	idNotVeryInterested          = 3
	idUnassigned                 = 4
	idSameEventAssignee          = 5
	idGMPlayer                   = 6
	idGMAndPlayerDifferentEvents = 7
)

// FirstChoice rules:
// - The event you were assigned to does not mark you as FirstChoice there because you are
//   already placed; the flag only becomes meaningful when you appear in other event interest lists.
// - GM status alone never sets FirstChoice; only player assignments do, and GM-only involvement
//   should never show you as FirstChoice.

func TestGetInterestsForEvent_FirstChoiceRules(t *testing.T) {
	db, logger, createDBErr := testutil.CreateTemporaryDBAndLogger("test_first_choice", t)
	if createDBErr != nil {
		t.Fatalf("failed to create test database: %v", createDBErr)
	}
	defer db.Close()

	seedBaseTables(t, db)
	seedBillettholdere(t, db, append(
		playerFixtures(),
		gmFixtures()...,
	))
	seedInterests(t, db, append(
		interestsForE1(),
		append(interestsForE2(), append(interestsForE3(), interestsForE4()...)...)...,
	))
	assignmentRows := append(assignmentsE1(), assignmentsE2()...)
	assignmentRows = append(assignmentRows, assignmentsE3()...)
	seedAssignments(t, db, assignmentRows)

	interestsE1, getInterestsE1Err := GetInterestsForEvent(eventE1, db, logger)
	if getInterestsE1Err != nil {
		t.Fatalf("GetInterestsForEvent E1 error: %v", getInterestsE1Err)
	}

	gotE1 := indexInterests(interestsE1)

	// E1 first-choice check confirms the assigned event does not mark FirstChoice there.
	t.Run("E1 first-choice rules", func(t *testing.T) {
		for _, tc := range []firstChoiceCase{
			{id: idPlayerAssigned, want: false, name: "player assigned in current event should not mark first choice here"},
		} {
			expectFirstChoice(t, gotE1, tc)
		}
	})

	interests, getInterestsErr := GetInterestsForEvent(eventE2, db, logger)
	if getInterestsErr != nil {
		t.Fatalf("GetInterestsForEvent error: %v", getInterestsErr)
	}

	got := indexInterests(interests)

	// E2 inclusion checks confirm interests are listed even if the player is already assigned
	// to the same event; assignment should only affect FirstChoice, not filtering.
	t.Run("E2 includes/excludes correct billettholders", func(t *testing.T) {
		expectPresent(t, got, idSameEventAssignee, "expected assigned-to-same-event billettholder to be returned")
		expectPresent(t, got, idPlayerAssigned, "expected player-assigned billettholder to be returned")
		expectPresent(t, got, idGMAssigned, "expected gm-assigned billettholder to be returned")
		expectPresent(t, got, idNotVeryInterested, "expected not-very-interested billettholder to be returned")
		expectPresent(t, got, idUnassigned, "expected unassigned billettholder to be returned")
		expectPresent(t, got, idGMPlayer, "expected gm+player billettholder to be returned")
		expectPresent(t, got, idGMAndPlayerDifferentEvents, "expected gm+player (different events) billettholder to be returned")
	})

	// E2 first-choice checks focus on the CASE logic in queryFirstChoice:
	// - "Veldig interessert" + assigned as player in a different event => FirstChoice should be true.
	// - GM-only in a different event should NOT count as FirstChoice.
	// - Any interest below "Veldig interessert" should NOT be FirstChoice, even if assigned elsewhere.
	// - No assignment at all should NOT be FirstChoice.
	t.Run("E2 first-choice rules", func(t *testing.T) {
		for _, tc := range []firstChoiceCase{
			{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
			{id: idGMAssigned, want: false, name: "gm assigned to other event"},
			{id: idNotVeryInterested, want: false, name: "not very interested"},
			{id: idUnassigned, want: false, name: "no assignment"},
			{id: idGMPlayer, want: true, name: "gm+player with very interested"},
			{id: idGMAndPlayerDifferentEvents, want: true, name: "gm in one event and player in another"},
		} {
			expectFirstChoice(t, got, tc)
		}
	})

	interestsE3, getInterestsE3Err := GetInterestsForEvent(eventE3, db, logger)
	if getInterestsE3Err != nil {
		t.Fatalf("GetInterestsForEvent E3 error: %v", getInterestsE3Err)
	}

	gotE3 := indexInterests(interestsE3)

	// E3 inclusion check confirms assignments to the same event do not filter interests out.
	t.Run("E3 includes/excludes correct billettholders", func(t *testing.T) {
		expectPresent(t, gotE3, idGMPlayer, "expected assigned-to-same-event billettholder to be returned for E3")
	})

	// E3 first-choice checks re-run the same CASE rules against a different event to confirm
	// the logic is not accidentally tied to E2-only data setup.
	t.Run("E3 first-choice rules", func(t *testing.T) {
		for _, tc := range []firstChoiceCase{
			{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
			{id: idGMAssigned, want: false, name: "gm assigned to other event"},
		} {
			expectFirstChoice(t, gotE3, tc)
		}
	})

	interestsE4, getInterestsE4Err := GetInterestsForEvent(eventE4, db, logger)
	if getInterestsE4Err != nil {
		t.Fatalf("GetInterestsForEvent E4 error: %v", getInterestsE4Err)
	}

	gotE4 := indexInterests(interestsE4)

	// E4 first-choice checks cover an interest mix with an explicit "no assignment" case to ensure
	// the FirstChoice flag remains false when the participant has no cross-event player assignment.
	t.Run("E4 first-choice rules", func(t *testing.T) {
		for _, tc := range []firstChoiceCase{
			{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
			{id: idGMAssigned, want: false, name: "gm assigned to other event"},
			{id: idUnassigned, want: false, name: "no assignment"},
		} {
			expectFirstChoice(t, gotE4, tc)
		}
	})
}

type firstChoiceCase struct {
	id   int
	want bool
	name string
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

func playerFixtures() []billettholderFixture {
	return []billettholderFixture{
		{id: idPlayerAssigned, firstName: "Player", lastName: "One"},
		{id: idNotVeryInterested, firstName: "NotVery", lastName: "Three"},
		{id: idUnassigned, firstName: "NoAssign", lastName: "Four"},
		{id: idSameEventAssignee, firstName: "SameEvent", lastName: "Five"},
	}
}

func gmFixtures() []billettholderFixture {
	return []billettholderFixture{
		{id: idGMAssigned, firstName: "GM", lastName: "Two"},
		{id: idGMPlayer, firstName: "GMPlayer", lastName: "Six"},
		{id: idGMAndPlayerDifferentEvents, firstName: "GMAndPlayer", lastName: "Seven"},
	}
}

func interestsForE2() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE2, puljeID: puljeP2, interestLevel: "Veldig interessert"},
		{billettholderID: idGMAssigned, eventID: eventE2, puljeID: puljeP2, interestLevel: "Veldig interessert"},
		{billettholderID: idNotVeryInterested, eventID: eventE2, puljeID: puljeP2, interestLevel: "Interessert"},
		{billettholderID: idUnassigned, eventID: eventE2, puljeID: puljeP2, interestLevel: "Litt interessert"},
		{billettholderID: idSameEventAssignee, eventID: eventE2, puljeID: puljeP2, interestLevel: "Veldig interessert"},
		{billettholderID: idGMPlayer, eventID: eventE2, puljeID: puljeP2, interestLevel: "Veldig interessert"},
		{billettholderID: idGMAndPlayerDifferentEvents, eventID: eventE2, puljeID: puljeP2, interestLevel: "Veldig interessert"},
	}
}

func interestsForE1() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE1, puljeID: puljeP1, interestLevel: "Veldig interessert"},
	}
}

func interestsForE3() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE3, puljeID: puljeP3, interestLevel: "Veldig interessert"},
		{billettholderID: idGMAssigned, eventID: eventE3, puljeID: puljeP3, interestLevel: "Veldig interessert"},
		{billettholderID: idGMPlayer, eventID: eventE3, puljeID: puljeP3, interestLevel: "Veldig interessert"},
	}
}

func interestsForE4() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE4, puljeID: puljeP4, interestLevel: "Veldig interessert"},
		{billettholderID: idGMAssigned, eventID: eventE4, puljeID: puljeP4, interestLevel: "Veldig interessert"},
		{billettholderID: idUnassigned, eventID: eventE4, puljeID: puljeP4, interestLevel: "Ikkje interessert"},
	}
}

func assignmentsE1() []assignmentFixture {
	return []assignmentFixture{
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idPlayerAssigned, isPlayer: 1, isGM: 0},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMAssigned, isPlayer: 0, isGM: 1},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idNotVeryInterested, isPlayer: 1, isGM: 0},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMPlayer, isPlayer: 0, isGM: 1},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMAndPlayerDifferentEvents, isPlayer: 0, isGM: 1},
	}
}

func assignmentsE2() []assignmentFixture {
	return []assignmentFixture{
		{eventID: eventE2, puljeID: puljeP2, billettholderID: idSameEventAssignee, isPlayer: 1, isGM: 0},
	}
}

func assignmentsE3() []assignmentFixture {
	return []assignmentFixture{
		{eventID: eventE3, puljeID: puljeP3, billettholderID: idGMPlayer, isPlayer: 1, isGM: 0},
		{eventID: eventE4, puljeID: puljeP4, billettholderID: idGMAndPlayerDifferentEvents, isPlayer: 1, isGM: 0},
	}
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

func expectPresent(t *testing.T, got map[int]InterestWithHolder, id int, message string) {
	t.Helper()
	if _, ok := got[id]; !ok {
		t.Fatal(message)
	}
}

func expectFirstChoice(t *testing.T, got map[int]InterestWithHolder, tc firstChoiceCase) {
	t.Helper()
	if got[tc.id].FirstChoice != tc.want {
		t.Errorf("%s should be first choice = %v", tc.name, tc.want)
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, execErr := db.Exec(query, args...); execErr != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", execErr, query)
	}
}

package formsubmission

import (
	"database/sql"
	"log/slog"
	"testing"

	"github.com/Regncon/conorganizer/models"
)

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
	interestLevel   models.InterestLevel
}

type assignmentFixture struct {
	eventID         string
	puljeID         string
	billettholderID int
	isPlayer        int
	isGM            int
}

func interestIndexForEvent(t *testing.T, eventID string, db *sql.DB, logger *slog.Logger) map[int]InterestWithHolder {
	t.Helper()

	interests, err := GetInterestsForEvent(eventID, db, logger)
	if err != nil {
		t.Fatalf("GetInterestsForEvent %s error: %v", eventID, err)
	}

	return indexInterests(t, interests)
}

func assigneeIndexForEvent(t *testing.T, eventID string, db *sql.DB, logger *slog.Logger) map[int]InterestWithHolder {
	t.Helper()

	assignees, err := GetAssigneesForEvent(eventID, db, logger)
	if err != nil {
		t.Fatalf("GetAssigneesForEvent %s error: %v", eventID, err)
	}

	return indexInterests(t, assignees)
}

func indexInterests(t *testing.T, interests []InterestWithHolder) map[int]InterestWithHolder {
	t.Helper()
	index := make(map[int]InterestWithHolder, len(interests))
	for _, interest := range interests {
		if prev, exists := index[interest.BillettholderId]; exists {
			t.Fatalf(
				"duplicate billettholder id %d found for event %s (pulje %s and %s)",
				interest.BillettholderId,
				interest.EventId,
				prev.PuljeId,
				interest.PuljeId,
			)
		}
		index[interest.BillettholderId] = interest
	}
	return index
}

func seedBaseTables(t *testing.T, db *sql.DB) {
	t.Helper()

	puljeStatus := models.PuljeStatusOpen

	mustExec(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	mustExec(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExec(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExec(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExec(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	mustExec(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, puljeStatus)
	mustExec(t, db, `
		INSERT INTO puljer (
			id, name, status, start_at, end_at
		) VALUES
			('P1', 'Friday', ?, '2025-10-03', '2025-10-03'),
			('P2', 'SaturdayMorning', ?, '2025-10-04', '2025-10-04'),
			('P3', 'SaturdayEvening', ?, '2025-10-04', '2025-10-04'),
			('P4', 'Sunday', ?, '2025-10-05', '2025-10-05')
	`, puljeStatus, puljeStatus, puljeStatus, puljeStatus)
	mustExec(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES
			('E1','Event 1','intro','desc','', ?,?,?,'Host 1','h1@test.no','11111111',4,1,1,?),
			('E2','Event 2','intro','desc','', ?,?,?,'Host 2','h2@test.no','22222222',4,1,1,?),
			('E3','Event 3','intro','desc','', ?,?,?,'Host 3','h3@test.no','33333333',4,1,1,?),
			('E4','Event 4','intro','desc','', ?,?,?,'Host 4','h4@test.no','44444444',4,1,1,?)
	`,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
	)
}

func playerFixtures() []billettholderFixture {
	return []billettholderFixture{
		{id: idPlayerAssigned, firstName: "Assigned", lastName: "One"},
		{id: idNotVeryInterested, firstName: "NotVery", lastName: "Three"},
		{id: idUnassigned, firstName: "NoAssign", lastName: "Four"},
		{id: idSameEventAssignee, firstName: "SameEvent", lastName: "Five"},
	}
}

func gmFixtures() []billettholderFixture {
	return []billettholderFixture{
		{id: idGMAssigned, firstName: "Gamemaster", lastName: "Two"},
		{id: idGMPlayer, firstName: "GMPlayer", lastName: "Six"},
		{id: idGMAndPlayerDifferentEvents, firstName: "GMAndPlayer", lastName: "Seven"},
		{id: idGMOnlyVeryInterestedOther, firstName: "GMOnlyVeryInterested", lastName: "Eight"},
	}
}

func interestsForE2() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMAssigned, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
		{billettholderID: idNotVeryInterested, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelMedium},
		{billettholderID: idUnassigned, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelLow},
		{billettholderID: idSameEventAssignee, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMPlayer, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMAndPlayerDifferentEvents, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMOnlyVeryInterestedOther, eventID: eventE2, puljeID: puljeP2, interestLevel: models.InterestLevelHigh},
	}
}

func interestsForE1() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE1, puljeID: puljeP1, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMOnlyVeryInterestedOther, eventID: eventE1, puljeID: puljeP1, interestLevel: models.InterestLevelHigh},
	}
}

func interestsForE3() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE3, puljeID: puljeP3, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMAssigned, eventID: eventE3, puljeID: puljeP3, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMPlayer, eventID: eventE3, puljeID: puljeP3, interestLevel: models.InterestLevelHigh},
	}
}

func interestsForE4() []interestFixture {
	return []interestFixture{
		{billettholderID: idPlayerAssigned, eventID: eventE4, puljeID: puljeP4, interestLevel: models.InterestLevelHigh},
		{billettholderID: idGMAssigned, eventID: eventE4, puljeID: puljeP4, interestLevel: models.InterestLevelHigh},
		{billettholderID: idUnassigned, eventID: eventE4, puljeID: puljeP4, interestLevel: models.InterestLevelLow},
	}
}

func assignmentsE1() []assignmentFixture {
	return []assignmentFixture{
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idPlayerAssigned, isPlayer: 1, isGM: 0},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMAssigned, isPlayer: 0, isGM: 1},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idNotVeryInterested, isPlayer: 1, isGM: 0},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMPlayer, isPlayer: 0, isGM: 1},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMAndPlayerDifferentEvents, isPlayer: 0, isGM: 1},
		{eventID: eventE1, puljeID: puljeP1, billettholderID: idGMOnlyVeryInterestedOther, isPlayer: 0, isGM: 1},
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
		role := models.EventPlayerRolePlayer
		if row.isGM == 1 {
			role = models.EventPlayerRoleGM
		}
		mustExec(t, db, `
			INSERT INTO relation_events_players (
				event_id, pulje_id, billettholder_id, role
			) VALUES (?, ?, ?, ?)
		`, row.eventID, row.puljeID, row.billettholderID, role)
	}
}

func hasInterest(got map[int]InterestWithHolder, id int) bool {
	_, ok := got[id]
	return ok
}

func expectPresent(t *testing.T, got map[int]InterestWithHolder, id int, message string) {
	t.Helper()
	if !hasInterest(got, id) {
		t.Fatal(message)
	}
}

func expectAbsent(t *testing.T, got map[int]InterestWithHolder, id int, message string) {
	t.Helper()
	if hasInterest(got, id) {
		t.Fatal(message)
	}
}

func expectFirstChoice(t *testing.T, got map[int]InterestWithHolder, tc firstChoiceCase) {
	t.Helper()
	interest, ok := got[tc.id]
	if !ok {
		t.Fatalf("%s: missing billettholder id %d", tc.name, tc.id)
	}
	if interest.FirstChoice != tc.want {
		t.Errorf("%s should be first choice = %v", tc.name, tc.want)
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, execErr := db.Exec(query, args...); execErr != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", execErr, query)
	}
}

package puljefordeling

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestGetFirstChoiceStatusesForEvent_DerivesCurrentOtherAndGMIgnored(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_statuses")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagMorgen, "Lordag morgen")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")

	seedFirstChoiceEvent(t, db, "friday-gm", "Friday GM", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-choice", "Saturday Choice", models.PuljeLordagMorgen)
	seedFirstChoiceEvent(t, db, "saturday-evening", "Saturday Evening", models.PuljeLordagKveld)

	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")
	seedFirstChoiceBillettholder(t, db, 2, "Bob", "Berg")
	seedFirstChoiceBillettholder(t, db, 3, "Cara", "Christensen")

	seedFirstChoiceInterest(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.EventPlayerRoleGM)
	seedFirstChoiceInterest(t, db, 1, "saturday-choice", models.PuljeLordagMorgen, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-choice", models.PuljeLordagMorgen, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	seedFirstChoiceInterest(t, db, 2, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 2, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	seedFirstChoiceInterest(t, db, 3, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelMedium)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-evening")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent: %v", err)
	}
	if len(statuses) != 3 {
		t.Fatalf("statuses length: want 3, got %d", len(statuses))
	}

	assertFirstChoiceStatus(t, statuses, 1, "saturday-evening", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: true,
		HasOtherPuljeFirstChoice:   true,
	})
	assertFirstChoiceStatus(t, statuses, 2, "saturday-evening", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: true,
		HasOtherPuljeFirstChoice:   false,
	})
	assertFirstChoiceStatus(t, statuses, 3, "saturday-evening", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: false,
		HasOtherPuljeFirstChoice:   false,
	})
}

func TestGetFirstChoiceStatusesForEvent_GMHighInterestDoesNotCountAsOtherPuljeFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_gm_ignored")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")

	seedFirstChoiceEvent(t, db, "friday-gm", "Friday GM", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-evening", "Saturday Evening", models.PuljeLordagKveld)

	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.EventPlayerRoleGM)
	seedFirstChoiceInterest(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelMedium)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-evening")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses length: want 1, got %d", len(statuses))
	}

	assertFirstChoiceStatus(t, statuses, 1, "saturday-evening", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: false,
		HasOtherPuljeFirstChoice:   false,
	})
}

func TestGetFirstChoiceStatusesForEvent_CurrentGMHighInterestDoesNotCountAsFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_current_gm_ignored")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "saturday-evening", "Saturday Evening", models.PuljeLordagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRoleGM)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-evening")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses length: want 1, got %d", len(statuses))
	}

	assertFirstChoiceStatus(t, statuses, 1, "saturday-evening", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: false,
		HasOtherPuljeFirstChoice:   false,
	})
}

func TestGetFirstChoiceStatusesForEvent_SameEventDifferentPuljeCountsAsOtherFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_same_event_other_pulje")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeLordagMorgen, "Lordag morgen")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "saturday-event", "Saturday Event", models.PuljeLordagKveld)
	seedFirstChoiceEventPulje(t, db, "saturday-event", models.PuljeLordagMorgen)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "saturday-event", models.PuljeLordagKveld, models.InterestLevelMedium)
	seedFirstChoiceInterest(t, db, 1, "saturday-event", models.PuljeLordagMorgen, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-event", models.PuljeLordagMorgen, models.EventPlayerRolePlayer)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-event")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("statuses length: want 2, got %d", len(statuses))
	}

	assertFirstChoiceStatus(t, statuses, 1, "saturday-event", models.PuljeLordagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: false,
		HasOtherPuljeFirstChoice:   true,
	})
	assertFirstChoiceStatus(t, statuses, 1, "saturday-event", models.PuljeLordagMorgen, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: true,
		HasOtherPuljeFirstChoice:   false,
	})
}

func TestGetFirstChoiceStatusesForEvent_DifferentEventSamePuljeCountsAsOtherFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_same_pulje_other_event")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "friday-event-a", "Friday Event A", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "friday-event-b", "Friday Event B", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "friday-event-a", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-event-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "friday-event-b", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "friday-event-b", models.PuljeFredagKveld, models.EventPlayerRolePlayer)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "friday-event-b")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("statuses length: want 1, got %d", len(statuses))
	}

	assertFirstChoiceStatus(t, statuses, 1, "friday-event-b", models.PuljeFredagKveld, FirstChoiceStatus{
		HasCurrentPuljeFirstChoice: false,
		HasOtherPuljeFirstChoice:   true,
	})
}

func TestSetAssignmentFirstChoice_SetsAndRemovesInterestWithoutChangingAssignment(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_toggle")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "event-1", "Event 1", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "event-1", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "event-1", models.PuljeFredagKveld, models.EventPlayerRolePlayer)

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, true); err != nil {
		t.Fatalf("SetAssignmentFirstChoice enable: %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelHigh {
		t.Fatalf("interest level after enable: want %q, got %q", models.InterestLevelHigh, got)
	}
	if got := queryFirstChoiceAssignmentRole(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.EventPlayerRolePlayer {
		t.Fatalf("assignment role after enable: want %q, got %q", models.EventPlayerRolePlayer, got)
	}

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, false); err != nil {
		t.Fatalf("SetAssignmentFirstChoice disable: %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelMedium {
		t.Fatalf("interest level after disable: want %q, got %q", models.InterestLevelMedium, got)
	}
	if got := queryFirstChoiceAssignmentRole(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.EventPlayerRolePlayer {
		t.Fatalf("assignment role after disable: want %q, got %q", models.EventPlayerRolePlayer, got)
	}
}

func TestSetAssignmentFirstChoice_RejectsGMAssignment(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_reject_gm")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "event-1", "Event 1", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "event-1", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "event-1", models.PuljeFredagKveld, models.EventPlayerRoleGM)

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, true); !errors.Is(err, ErrFirstChoiceGMAssignment) {
		t.Fatalf("SetAssignmentFirstChoice enable GM: want ErrFirstChoiceGMAssignment, got %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelHigh {
		t.Fatalf("interest level after rejected GM update: want %q, got %q", models.InterestLevelHigh, got)
	}
}

func TestSetAssignmentFirstChoice_RejectsInvalidInputWithSentinel(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_invalid_input")

	tests := []struct {
		name            string
		eventID         string
		puljeID         string
		billettholderID int
	}{
		{
			name:            "missing event id",
			eventID:         "",
			puljeID:         string(models.PuljeFredagKveld),
			billettholderID: 1,
		},
		{
			name:            "missing pulje id",
			eventID:         "event-1",
			puljeID:         "",
			billettholderID: 1,
		},
		{
			name:            "missing billettholder id",
			eventID:         "event-1",
			puljeID:         string(models.PuljeFredagKveld),
			billettholderID: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := SetAssignmentFirstChoice(db, test.eventID, test.puljeID, test.billettholderID, true); !errors.Is(err, ErrFirstChoiceInvalidInput) {
				t.Fatalf("SetAssignmentFirstChoice invalid input: want ErrFirstChoiceInvalidInput, got %v", err)
			}
		})
	}
}

func TestSetAssignmentFirstChoice_RejectsMissingAssignmentWithSentinel(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_missing_assignment")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "event-1", "Event 1", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")
	seedFirstChoiceInterest(t, db, 1, "event-1", models.PuljeFredagKveld, models.InterestLevelMedium)

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, true); !errors.Is(err, ErrFirstChoiceMissingAssignment) {
		t.Fatalf("SetAssignmentFirstChoice missing assignment: want ErrFirstChoiceMissingAssignment, got %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelMedium {
		t.Fatalf("interest level after rejected missing assignment update: want %q, got %q", models.InterestLevelMedium, got)
	}
}

func TestSetAssignmentFirstChoice_RejectsOtherPuljeFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_reject_other_pulje")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "friday-event", "Friday Event", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-event", "Saturday Event", models.PuljeLordagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "friday-event", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-event", models.PuljeFredagKveld, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "saturday-event", models.PuljeLordagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "saturday-event", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	if err := SetAssignmentFirstChoice(db, "saturday-event", string(models.PuljeLordagKveld), 1, true); !errors.Is(err, ErrFirstChoiceOtherPuljeFirstChoice) {
		t.Fatalf("SetAssignmentFirstChoice enable with other pulje first-choice: want ErrFirstChoiceOtherPuljeFirstChoice, got %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "saturday-event", models.PuljeLordagKveld); got != models.InterestLevelMedium {
		t.Fatalf("saturday interest level after rejected update: want %q, got %q", models.InterestLevelMedium, got)
	}
}

func TestSetAssignmentFirstChoice_RejectsSamePuljeOtherEventFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_reject_same_pulje_other_event")
	seedFirstChoiceLookups(t, db)

	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "friday-event-a", "Friday Event A", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "friday-event-b", "Friday Event B", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Andersen")

	seedFirstChoiceInterest(t, db, 1, "friday-event-a", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-event-a", models.PuljeFredagKveld, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "friday-event-b", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "friday-event-b", models.PuljeFredagKveld, models.EventPlayerRolePlayer)

	if err := SetAssignmentFirstChoice(db, "friday-event-b", string(models.PuljeFredagKveld), 1, true); !errors.Is(err, ErrFirstChoiceOtherPuljeFirstChoice) {
		t.Fatalf("SetAssignmentFirstChoice enable with same pulje other event first-choice: want ErrFirstChoiceOtherPuljeFirstChoice, got %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "friday-event-b", models.PuljeFredagKveld); got != models.InterestLevelMedium {
		t.Fatalf("event B interest level after rejected update: want %q, got %q", models.InterestLevelMedium, got)
	}
}

func assertFirstChoiceStatus(
	t *testing.T,
	statuses map[FirstChoiceKey]FirstChoiceStatus,
	billettholderID int,
	eventID string,
	pulje models.Pulje,
	want FirstChoiceStatus,
) {
	t.Helper()

	key := FirstChoiceKey{
		BillettholderID: billettholderID,
		EventID:         eventID,
		PuljeID:         string(pulje),
	}
	got, ok := statuses[key]
	if !ok {
		t.Fatalf("status for %+v missing", key)
	}
	if got != want {
		t.Fatalf("status for %+v: want %+v, got %+v", key, want, got)
	}
}

func queryFirstChoiceInterestLevel(
	t *testing.T,
	db *sql.DB,
	billettholderID int,
	eventID string,
	pulje models.Pulje,
) models.InterestLevel {
	t.Helper()

	var level models.InterestLevel
	if err := db.QueryRow(`
		SELECT interest_level
		FROM interests
		WHERE billettholder_id = ?
			AND event_id = ?
			AND pulje_id = ?
	`, billettholderID, eventID, pulje).Scan(&level); err != nil {
		t.Fatalf("query first-choice interest level: %v", err)
	}
	return level
}

func queryFirstChoiceAssignmentRole(
	t *testing.T,
	db *sql.DB,
	billettholderID int,
	eventID string,
	pulje models.Pulje,
) models.EventPlayerRole {
	t.Helper()

	var role models.EventPlayerRole
	if err := db.QueryRow(`
		SELECT role
		FROM relation_events_players
		WHERE billettholder_id = ?
			AND event_id = ?
			AND pulje_id = ?
	`, billettholderID, eventID, pulje).Scan(&role); err != nil {
		t.Fatalf("query first-choice assignment role: %v", err)
	}
	return role
}

func seedFirstChoiceLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.EventStatusApproved)
	testutil.MustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
	for _, level := range []models.InterestLevel{
		models.InterestLevelHigh,
		models.InterestLevelMedium,
		models.InterestLevelLow,
	} {
		testutil.MustExec(t, db, `INSERT INTO interest_levels(interest_level) VALUES (?) ON CONFLICT(interest_level) DO NOTHING`, level)
	}
	testutil.MustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.PuljeStatusOpen)
}

func seedFirstChoicePulje(t *testing.T, db *sql.DB, pulje models.Pulje, name string) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES (?, ?, ?, '2026-09-04T18:00:00Z', '2026-09-04T23:00:00Z')
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			status = excluded.status,
			start_at = excluded.start_at,
			end_at = excluded.end_at
	`, pulje, name, models.PuljeStatusOpen)
}

func seedFirstChoiceEvent(t *testing.T, db *sql.DB, eventID string, title string, pulje models.Pulje) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO events(
			id, title, intro, description, system, event_type, age_group, event_runtime,
			host_name, email, phone_number, max_players, beginner_friendly,
			can_be_run_in_english, status
		)
		VALUES (?, ?, 'Intro', 'Description', 'System', ?, ?, ?, 'Host', 'host@example.com',
			'12345678', 4, 1, 1, ?)
	`, eventID, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved)

	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
		ON CONFLICT(event_id, pulje_id) DO UPDATE SET
			is_in_pulje = excluded.is_in_pulje,
			is_published = excluded.is_published
	`, eventID, pulje)
}

func seedFirstChoiceEventPulje(t *testing.T, db *sql.DB, eventID string, pulje models.Pulje) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
		ON CONFLICT(event_id, pulje_id) DO UPDATE SET
			is_in_pulje = excluded.is_in_pulje,
			is_published = excluded.is_published
	`, eventID, pulje)
}

func seedFirstChoiceBillettholder(t *testing.T, db *sql.DB, id int, firstName string, lastName string) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO billettholdere(
			id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id
		)
		VALUES (?, ?, ?, ?, 'Festivalpass', ?, ?)
	`, id, firstName, lastName, 1000+id, 2000+id, 3000+id)
}

func seedFirstChoiceInterest(
	t *testing.T,
	db *sql.DB,
	billettholderID int,
	eventID string,
	pulje models.Pulje,
	level models.InterestLevel,
) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
	`, billettholderID, eventID, pulje, level)
}

func seedFirstChoiceAssignment(
	t *testing.T,
	db *sql.DB,
	billettholderID int,
	eventID string,
	pulje models.Pulje,
	role models.EventPlayerRole,
) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role)
		VALUES (?, ?, ?, ?)
	`, eventID, pulje, billettholderID, role)
}

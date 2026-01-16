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

	mustExec(t, db, `INSERT INTO event_statuses(status) VALUES ('Godkjent')`)
	mustExec(t, db, `INSERT INTO events_types(event_type) VALUES ('Other')`)
	mustExec(t, db, `INSERT INTO age_groups(age_group) VALUES ('Default')`)
	mustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES ('Normal')`)
	mustExec(t, db, `INSERT INTO interest_levels(interest_level) VALUES ('Veldig interessert'), ('Interessert')`)
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
	mustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES
			(1,'Player','One',1,'Test',1,1001,2001),
			(2,'GM','Two',1,'Test',1,1002,2002),
			(3,'NotVery','Three',1,'Test',1,1003,2003),
			(4,'NoAssign','Four',1,'Test',1,1004,2004),
			(5,'SameEvent','Five',1,'Test',1,1005,2005),
			(6,'GMPlayer','Six',1,'Test',1,1006,2006)
	`)
	mustExec(t, db, `
		INSERT INTO interests (
			billettholder_id, event_id, pulje_id, interest_level
		) VALUES
			(1,'E2','P2','Veldig interessert'),
			(2,'E2','P2','Veldig interessert'),
			(3,'E2','P2','Interessert'),
			(4,'E2','P2','Veldig interessert'),
			(5,'E2','P2','Veldig interessert'),
			(6,'E2','P2','Veldig interessert'),
			(1,'E3','P3','Veldig interessert'),
			(2,'E3','P3','Veldig interessert'),
			(6,'E3','P3','Veldig interessert'),
			(1,'E4','P4','Veldig interessert'),
			(2,'E4','P4','Veldig interessert'),
			(4,'E4','P4','Veldig interessert')
	`)
	mustExec(t, db, `
		INSERT INTO events_players (
			event_id, pulje_id, billettholder_id, is_player, is_gm
		) VALUES
			('E1','P1',1,1,0),
			('E1','P1',2,0,1),
			('E1','P1',3,1,0),
			('E2','P2',5,1,0),
			('E1','P1',6,0,1),
			('E3','P3',6,1,0)
	`)

	interests, getInterestsErr := GetInterestsForEvent("E2", db, logger)
	if getInterestsErr != nil {
		t.Fatalf("GetInterestsForEvent error: %v", getInterestsErr)
	}

	got := make(map[int]InterestWithHolder, len(interests))
	for _, interest := range interests {
		got[interest.BillettholderId] = interest
	}

	if _, ok := got[5]; ok {
		t.Fatalf("expected assigned-to-same-event billettholder to be excluded")
	}
	if _, ok := got[1]; !ok {
		t.Fatalf("expected player-assigned billettholder to be returned")
	}
	if _, ok := got[2]; !ok {
		t.Fatalf("expected gm-assigned billettholder to be returned")
	}
	if _, ok := got[3]; !ok {
		t.Fatalf("expected not-very-interested billettholder to be returned")
	}
	if _, ok := got[4]; !ok {
		t.Fatalf("expected unassigned billettholder to be returned")
	}
	if _, ok := got[6]; !ok {
		t.Fatalf("expected gm+player billettholder to be returned")
	}

	if got[1].FirstChoice != true {
		t.Errorf("player assigned to other event should be first choice")
	}
	if got[2].FirstChoice != false {
		t.Errorf("gm assigned to other event should not be first choice")
	}
	if got[3].FirstChoice != false {
		t.Errorf("not very interested should not be first choice")
	}
	if got[4].FirstChoice != false {
		t.Errorf("no assignment should not be first choice")
	}
	if got[6].FirstChoice != true {
		t.Errorf("gm+player with very interested should be first choice due to player assignment")
	}

	interestsE3, getInterestsE3Err := GetInterestsForEvent("E3", db, logger)
	if getInterestsE3Err != nil {
		t.Fatalf("GetInterestsForEvent E3 error: %v", getInterestsE3Err)
	}

	gotE3 := make(map[int]InterestWithHolder, len(interestsE3))
	for _, interest := range interestsE3 {
		gotE3[interest.BillettholderId] = interest
	}

	if _, ok := gotE3[6]; ok {
		t.Fatalf("expected assigned-to-same-event billettholder to be excluded for E3")
	}
	if gotE3[1].FirstChoice != true {
		t.Errorf("player assigned to other event should be first choice for E3")
	}
	if gotE3[2].FirstChoice != false {
		t.Errorf("gm assigned to other event should not be first choice for E3")
	}

	interestsE4, getInterestsE4Err := GetInterestsForEvent("E4", db, logger)
	if getInterestsE4Err != nil {
		t.Fatalf("GetInterestsForEvent E4 error: %v", getInterestsE4Err)
	}

	gotE4 := make(map[int]InterestWithHolder, len(interestsE4))
	for _, interest := range interestsE4 {
		gotE4[interest.BillettholderId] = interest
	}

	if gotE4[1].FirstChoice != true {
		t.Errorf("player assigned to other event should be first choice for E4")
	}
	if gotE4[2].FirstChoice != false {
		t.Errorf("gm assigned to other event should not be first choice for E4")
	}
	if gotE4[4].FirstChoice != false {
		t.Errorf("no assignment should not be first choice for E4")
	}
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()
	if _, execErr := db.Exec(query, args...); execErr != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", execErr, query)
	}
}

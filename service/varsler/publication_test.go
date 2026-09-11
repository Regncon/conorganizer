package varsler

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestUpdatePuljeStatus_FirstCompletionQueuesAssignedAndInterestedRecipients(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_first_publication")
	db.SetMaxOpenConns(1)
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	seedRecipient(t, db, 2, 22, "https://updates.push.services.mozilla.com/wpush/v2/two")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',? ,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (22,'event-a',?,'Litt interessert')`, models.PuljeFredagKveld)

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err != nil {
		t.Fatalf("expected completion to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs WHERE status = 'pending'`); got != 2 {
		t.Fatalf("expected one queued job for each affected recipient, got %d", got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM pulje_varsel_results`); got != 2 {
		t.Fatalf("expected durable results for both recipients, got %d", got)
	}
	var unassignedResult string
	if err := db.QueryRow(`SELECT result_json FROM pulje_varsel_results WHERE billettholder_id = 22`).Scan(&unassignedResult); err != nil {
		t.Fatalf("query interested recipient result: %v", err)
	}
	if !strings.Contains(unassignedResult, `"assignments":[]`) {
		t.Fatalf("expected an explicit empty assignment snapshot, got %s", unassignedResult)
	}
}

func TestUpdatePuljeStatus_UnpublishedProgramDoesNotCreatePublication(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_program_unpublished")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `UPDATE program_publishing_state SET is_published = 0 WHERE id = 1`)

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err != nil {
		t.Fatalf("expected status update to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM pulje_varsel_publications`); got != 0 {
		t.Fatalf("expected no publication while the program is hidden, got %d", got)
	}
}

func TestUpdatePuljeStatus_HiddenEventDoesNotCreateRecipientResult(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_event_unpublished")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET is_published = 0 WHERE event_id = 'event-a'`)

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err != nil {
		t.Fatalf("expected status update to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM pulje_varsel_results`); got != 0 {
		t.Fatalf("expected hidden event to produce no recipient result, got %d", got)
	}
}

func TestUpdatePuljeStatus_RepeatedCompletionDoesNotDuplicateJobs(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_repeated_publication")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("first completion failed: %v", err)
	}

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err != nil {
		t.Fatalf("expected repeated completion to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs`); got != 1 {
		t.Fatalf("expected the original job only, got %d", got)
	}
}

func TestUpdatePuljeStatus_UnpublishingCancelsPendingJobs(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_unpublish")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("completion failed: %v", err)
	}

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusLocked)

	// Then
	if err != nil {
		t.Fatalf("expected unpublish to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs WHERE status = 'canceled'`); got != 1 {
		t.Fatalf("expected the queued job to be canceled, got %d", got)
	}
}

func TestUpdatePuljeStatus_RepublishQueuesOnlyChangedResults(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_changed_republish")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	seedRecipient(t, db, 2, 22, "https://updates.push.services.mozilla.com/wpush/v2/two")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player'),('event-a',?,22,'GM')`, models.PuljeFredagKveld, models.PuljeFredagKveld)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("first completion failed: %v", err)
	}
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusLocked); err != nil {
		t.Fatalf("unpublish failed: %v", err)
	}
	testutil.MustExec(t, db, `UPDATE relation_events_players SET role = 'GM' WHERE billettholder_id = 11`)

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err != nil {
		t.Fatalf("expected republish to succeed: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs WHERE publication_id = (SELECT MAX(id) FROM pulje_varsel_publications)`); got != 1 {
		t.Fatalf("expected only the changed result to queue, got %d", got)
	}
}

func TestUpdatePuljeStatus_RollsBackStatusWhenQueueingFails(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_atomic")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `CREATE TRIGGER reject_varsel_jobs BEFORE INSERT ON web_push_jobs BEGIN SELECT RAISE(ABORT, 'queue unavailable'); END`)

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if err == nil {
		t.Fatal("expected queue failure")
	}
	var status models.PuljeStatus
	if scanErr := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, models.PuljeFredagKveld).Scan(&status); scanErr != nil {
		t.Fatalf("query pulje status: %v", scanErr)
	}
	if status != models.PuljeStatusLocked {
		t.Fatalf("expected status rollback to Locked, got %s", status)
	}
}

func TestUpdatePuljeStatus_MissingPuljeReturnsSentinel(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_missing_pulje")

	// When
	err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if !errors.Is(err, ErrPuljeNotFound) {
		t.Fatalf("expected ErrPuljeNotFound, got %v", err)
	}
}

func seedPublicationFixture(t *testing.T, db DBExecutor) {
	t.Helper()
	mustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES ('Open'),('Locked'),('Completed') ON CONFLICT(status) DO NOTHING`)
	mustExec(t, db, `INSERT INTO event_statuses(status) VALUES ('Annonsert') ON CONFLICT(status) DO NOTHING`)
	mustExec(t, db, `INSERT INTO interest_levels(interest_level) VALUES ('Litt interessert') ON CONFLICT(interest_level) DO NOTHING`)
	mustExec(t, db, `INSERT INTO program_publishing_state(id,is_published) VALUES (1,1) ON CONFLICT(id) DO UPDATE SET is_published=1`)
	mustExec(t, db, `INSERT INTO puljer(id,name,status,start_at,end_at) VALUES (?,'Fredag Kveld','Locked','2026-09-04T18:00:00Z','2026-09-04T23:00:00Z')`, models.PuljeFredagKveld)
	mustExec(t, db, `INSERT INTO rooms(id,room_number,name,floor,max_concurrent_games) VALUES (1,'101','Hovedrom',1,4)`)
	mustExec(t, db, `INSERT INTO events(id,title,intro,description,host_name,email,phone_number,max_players,status) VALUES ('event-a','Event A','','','','','',5,'Annonsert')`)
	mustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje,is_published,room_id) VALUES ('event-a',?,1,1,1)`, models.PuljeFredagKveld)
}

func seedRecipient(t *testing.T, db DBExecutor, userID, billettholderID int, endpoint string) {
	t.Helper()
	mustExec(t, db, `INSERT INTO users(id,external_id,email) VALUES (?,?,?)`, userID, "external-"+string(rune('0'+userID)), "user@example.com")
	mustExec(t, db, `INSERT INTO billettholdere(id,first_name,last_name,ticket_type_id,ticket_type,order_id,ticket_id) VALUES (?,'Ola','Nordmann',1,'Helg',1,?)`, billettholderID, billettholderID)
	mustExec(t, db, `INSERT INTO relation_billettholdere_users(billettholder_id,user_id) VALUES (?,?)`, billettholderID, userID)
	mustExec(t, db, `INSERT INTO web_push_subscriptions(user_id,endpoint,p256dh,auth) VALUES (?,?, 'p256dh','auth')`, userID, endpoint)
}

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func mustExec(t *testing.T, db DBExecutor, query string, args ...any) {
	t.Helper()
	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("seed fixture: %v\nquery: %s", err, query)
	}
}

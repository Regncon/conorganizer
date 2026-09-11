package varsler

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

type fakePushClient struct {
	status int
	err    error
	calls  int
	hook   func()
}

func (client *fakePushClient) Do(*http.Request) (*http.Response, error) {
	client.calls++
	if client.hook != nil {
		client.hook()
	}
	if client.err != nil {
		return nil, client.err
	}
	return &http.Response{
		StatusCode: client.status,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}, nil
}

func TestWorker_UnpublishDuringSendKeepsJobCanceled(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusServiceUnavailable)
	client.hook = func() {
		if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusLocked); err != nil {
			t.Errorf("unpublish during send: %v", err)
		}
	}

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected fenced completion: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected one send attempt, processed=%v calls=%d", processed, client.calls)
	}
	if got := queryJobStatus(t, db); got != "canceled" {
		t.Fatalf("expected unpublish cancellation to win, got %s", got)
	}
}

func TestWorker_SuccessMarksJobSent(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusCreated)

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected worker success: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected one processed HTTP request, processed=%v calls=%d", processed, client.calls)
	}
	if got := queryJobStatus(t, db); got != "sent" {
		t.Fatalf("expected sent job, got %s", got)
	}
}

func TestWorker_ClaimsJobCreatedEarlierInSameMillisecond(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusCreated)
	var nextAttempt string
	if err := db.QueryRow(`SELECT next_attempt_at FROM web_push_jobs`).Scan(&nextAttempt); err != nil {
		t.Fatalf("query job timestamp: %v", err)
	}
	createdAt, err := time.Parse(time.RFC3339Nano, nextAttempt)
	if err != nil {
		t.Fatalf("parse job timestamp: %v", err)
	}
	service.now = func() time.Time { return createdAt.Add(100 * time.Microsecond) }

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected worker success: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected due job to be claimed, processed=%v calls=%d", processed, client.calls)
	}
}

func TestWorker_TransientFailureSchedulesBoundedRetry(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusServiceUnavailable)
	var queuedAt string
	if err := db.QueryRow(`SELECT next_attempt_at FROM web_push_jobs`).Scan(&queuedAt); err != nil {
		t.Fatalf("query queued time: %v", err)
	}
	now, err := time.Parse(time.RFC3339Nano, queuedAt)
	if err != nil {
		t.Fatalf("parse queued time: %v", err)
	}
	now = now.Add(time.Second)
	service.now = func() time.Time { return now }

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected handled push failure: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected one processed HTTP request, processed=%v calls=%d", processed, client.calls)
	}
	var status, nextAttempt string
	var attempts int
	if err := db.QueryRow(`SELECT status, attempts, next_attempt_at FROM web_push_jobs`).Scan(&status, &attempts, &nextAttempt); err != nil {
		t.Fatalf("query retried job: %v", err)
	}
	if status != "retrying" || attempts != 1 || nextAttempt != now.Add(time.Minute).Format(time.RFC3339Nano) {
		t.Fatalf("expected first retry in one minute, got status=%s attempts=%d next=%s", status, attempts, nextAttempt)
	}
}

func TestWorker_FinalTransportFailureExhaustsJobWithoutPersistingEndpoint(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusCreated)
	client.err = errors.New("POST https://fcm.googleapis.com/fcm/send/secret-device: connection reset")
	testutil.MustExec(t, db, `UPDATE web_push_jobs SET attempts = 4`)

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected handled final failure: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected final send attempt, processed=%v calls=%d", processed, client.calls)
	}
	var status, lastError string
	if err := db.QueryRow(`SELECT status,last_error FROM web_push_jobs`).Scan(&status, &lastError); err != nil {
		t.Fatalf("query exhausted job: %v", err)
	}
	if status != "exhausted" {
		t.Fatalf("expected exhausted job, got %s", status)
	}
	if strings.Contains(lastError, "secret-device") || lastError != "Web Push transport failed" {
		t.Fatalf("expected sanitized transport error, got %q", lastError)
	}
}

func TestWorker_ExpiredFinalLeaseMarksAbandonedJobExhausted(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusCreated)
	now := time.Date(2026, 9, 10, 20, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }
	testutil.MustExec(t, db, `UPDATE web_push_jobs SET status = 'processing', attempts = 5, lease_until = ?, lease_token = 'dead-worker'`, now.Add(-time.Minute).Format(time.RFC3339Nano))

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected abandoned job cleanup: %v", err)
	}
	if processed || client.calls != 0 {
		t.Fatalf("expected cleanup without sending, processed=%v calls=%d", processed, client.calls)
	}
	if got := queryJobStatus(t, db); got != "exhausted" {
		t.Fatalf("expected abandoned job exhausted, got %s", got)
	}
}

func TestWorker_ExpiredEndpointRemovesSubscriptionAndCancelsJobs(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusGone)

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected expired endpoint handling: %v", err)
	}
	if !processed || client.calls != 1 {
		t.Fatalf("expected one processed HTTP request, processed=%v calls=%d", processed, client.calls)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_subscriptions`); got != 0 {
		t.Fatalf("expected expired subscription removal, got %d", got)
	}
	if got := queryJobStatus(t, db); got != "canceled" {
		t.Fatalf("expected expired endpoint job canceled, got %s", got)
	}
}

func TestWorker_RemovedOwnershipCancelsWithoutSending(t *testing.T) {
	// Given
	service, client, db := seededWorker(t, http.StatusCreated)
	testutil.MustExec(t, db, `DELETE FROM relation_billettholdere_users`)

	// When
	processed, err := service.processNext(context.Background())

	// Then
	if err != nil {
		t.Fatalf("expected stale ownership handling: %v", err)
	}
	if !processed || client.calls != 0 {
		t.Fatalf("expected stale job canceled without HTTP, processed=%v calls=%d", processed, client.calls)
	}
	if got := queryJobStatus(t, db); got != "canceled" {
		t.Fatalf("expected stale job canceled, got %s", got)
	}
}

func seededWorker(t *testing.T, status int) (*Service, *fakePushClient, *sql.DB) {
	t.Helper()
	db, logger := testutil.CreateTestDBAndLogger(t, "varsler_worker")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/worker-device")
	p256dh, auth := validSubscriptionKeys()
	testutil.MustExec(t, db, `UPDATE web_push_subscriptions SET p256dh = ?, auth = ?`, p256dh, auth)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("queue worker fixture: %v", err)
	}
	service, err := New(db, logger, testConfig(t))
	if err != nil {
		t.Fatalf("create varsler service: %v", err)
	}
	client := &fakePushClient{status: status}
	service.httpClient = client
	return service, client, db
}

func queryJobStatus(t *testing.T, db *sql.DB) string {
	t.Helper()
	var status string
	if err := db.QueryRow(`SELECT status FROM web_push_jobs ORDER BY id LIMIT 1`).Scan(&status); err != nil {
		t.Fatalf("query job status: %v", err)
	}
	return status
}

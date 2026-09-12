package varsler

import (
	"context"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestUpdateProgramPublished_PublishesCompletedPuljeThatWasCompletedWhileHidden(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_global_publish")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `UPDATE program_publishing_state SET is_published = 0 WHERE id = 1`)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("complete hidden pulje: %v", err)
	}

	// When
	err := UpdateProgramPublished(context.Background(), db, true)

	// Then
	if err != nil {
		t.Fatalf("publish program: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs WHERE status = 'pending'`); got != 1 {
		t.Fatalf("expected completed pulje result queued on global publication, got %d", got)
	}
}

func TestUpdateProgramPublished_RepeatedPublicationDoesNotDuplicateJobs(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_repeated_global_publish")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `UPDATE program_publishing_state SET is_published = 0 WHERE id = 1`)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("complete hidden pulje: %v", err)
	}
	if err := UpdateProgramPublished(context.Background(), db, true); err != nil {
		t.Fatalf("first global publication: %v", err)
	}

	// When
	err := UpdateProgramPublished(context.Background(), db, true)

	// Then
	if err != nil {
		t.Fatalf("repeat global publication: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs`); got != 1 {
		t.Fatalf("expected original job only, got %d", got)
	}
}

func TestUpdateProgramPublished_UnpublishingCancelsPendingJobs(t *testing.T) {
	// Given
	db := testutil.CreateTestDB(t, "varsler_global_unpublish")
	seedPublicationFixture(t, db)
	seedRecipient(t, db, 1, 11, "https://fcm.googleapis.com/fcm/send/one")
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('event-a',?,11,'Player')`, models.PuljeFredagKveld)
	if err := UpdatePuljeStatus(context.Background(), db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("complete pulje: %v", err)
	}

	// When
	err := UpdateProgramPublished(context.Background(), db, false)

	// Then
	if err != nil {
		t.Fatalf("unpublish program: %v", err)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs WHERE status = 'canceled'`); got != 1 {
		t.Fatalf("expected pending job canceled, got %d", got)
	}
}

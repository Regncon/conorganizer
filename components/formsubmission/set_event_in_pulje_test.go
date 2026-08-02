package formsubmission

import (
	"context"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
)

func TestSetEventInPulje_AddsThenRemovesMembership(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "set_event_in_pulje")
	testutil.MustExec(t, db,
		`INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4)`)

	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), true); err != nil {
		t.Fatalf("add: %v", err)
	}
	got := testutil.QueryInt(t, db,
		`SELECT is_in_pulje FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if got != 1 {
		t.Fatalf("after add is_in_pulje = %d, want 1", got)
	}

	// Idempotent + removal.
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), false); err != nil {
		t.Fatalf("remove: %v", err)
	}
	got = testutil.QueryInt(t, db,
		`SELECT is_in_pulje FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if got != 0 {
		t.Fatalf("after remove is_in_pulje = %d, want 0", got)
	}
}

// TestSetEventInPulje_NoOpReAddPreservesPublished verifies that re-dropping a
// card onto the pulje it already belongs to (isInPulje unchanged: true -> true)
// does not clear is_published. A dumb "reset is_published = 0 on every upsert"
// would silently unpublish an already-published game on a no-op re-drop.
func TestSetEventInPulje_NoOpReAddPreservesPublished(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "set_event_in_pulje_noop")
	testutil.MustExec(t, db,
		`INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4)`)
	testutil.MustExec(t, db,
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published) VALUES ('e1', ?, 1, 1)`,
		string(models.PuljeFredagKveld))

	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), true); err != nil {
		t.Fatalf("no-op re-add: %v", err)
	}

	gotInPulje := testutil.QueryInt(t, db,
		`SELECT is_in_pulje FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if gotInPulje != 1 {
		t.Fatalf("after no-op re-add is_in_pulje = %d, want 1", gotInPulje)
	}
	gotPublished := testutil.QueryInt(t, db,
		`SELECT is_published FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if gotPublished != 1 {
		t.Fatalf("after no-op re-add is_published = %d, want 1 (must survive a no-op re-drop)", gotPublished)
	}
}

// TestSetEventInPulje_RealChangeResetsPublished verifies that an actual
// membership flip (true -> false) DOES clear is_published, and that a
// subsequent genuine re-add (false -> true) does NOT auto-republish — the
// preservation rule only applies when the membership value is unchanged.
func TestSetEventInPulje_RealChangeResetsPublished(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "set_event_in_pulje_realchange")
	testutil.MustExec(t, db,
		`INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spel', '', '', 'Ola', 'ola@x.no', '', 4)`)
	testutil.MustExec(t, db,
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published) VALUES ('e1', ?, 1, 1)`,
		string(models.PuljeFredagKveld))

	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")

	// Real change: true -> false must reset is_published.
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), false); err != nil {
		t.Fatalf("remove: %v", err)
	}
	gotPublished := testutil.QueryInt(t, db,
		`SELECT is_published FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if gotPublished != 0 {
		t.Fatalf("after real removal is_published = %d, want 0", gotPublished)
	}

	// Genuine re-add: false -> true must NOT auto-republish.
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), true); err != nil {
		t.Fatalf("re-add: %v", err)
	}
	gotPublished = testutil.QueryInt(t, db,
		`SELECT is_published FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))
	if gotPublished != 0 {
		t.Fatalf("after genuine re-add is_published = %d, want 0 (not auto-republished)", gotPublished)
	}
}

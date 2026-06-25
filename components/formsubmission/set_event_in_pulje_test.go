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

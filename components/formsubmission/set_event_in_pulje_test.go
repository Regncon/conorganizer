package formsubmission

import (
	"context"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
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
		 VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4)`)

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

// TestSetEventInPulje_DoesNotChangeLegacyPublished verifies that membership
// updates leave the legacy is_published column untouched.
func TestSetEventInPulje_DoesNotChangeLegacyPublished(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en puljerad allerede er med i puljen og er markert som publisert.",
		When:  "Når samme medlemskap lagres på nytt.",
		Then:  "Så skal det gamle publiseringsflagget forbli uendret.",
	})

	// Given
	expectedPublished := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "set_event_in_pulje_noop")
	testutil.MustExec(t, db,
		`INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4)`)
	testutil.MustExec(t, db,
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published) VALUES ('e1', ?, 1, 1)`,
		string(models.PuljeFredagKveld))

	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
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

	// Then
	if gotPublished != expectedPublished {
		t.Fatalf("after no-op re-add is_published = %d, want %d", gotPublished, expectedPublished)
	}
}

// TestSetEventInPulje_RealChangeLeavesLegacyPublishedUntouched verifies that
// membership changes do not change the legacy is_published value.
func TestSetEventInPulje_RealChangeLeavesLegacyPublishedUntouched(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en puljerad er med i puljen og er markert som publisert.",
		When:  "Når medlemskapet fjernes og legges til igjen.",
		Then:  "Så skal det gamle publiseringsflagget forbli uendret gjennom begge endringene.",
	})

	// Given
	expectedPublished := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "set_event_in_pulje_realchange")
	testutil.MustExec(t, db,
		`INSERT INTO users (id, external_id, email, is_admin) VALUES (42, 'ext-42', 'admin@x.no', 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spill', '', '', 'Ola', 'ola@x.no', '', 4)`)
	testutil.MustExec(t, db,
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published) VALUES ('e1', ?, 1, 1)`,
		string(models.PuljeFredagKveld))

	ctx := authctx.WithUserToken(context.Background(), "ext-42", "admin@x.no")

	// When
	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), false); err != nil {
		t.Fatalf("remove: %v", err)
	}
	gotPublishedAfterRemoval := testutil.QueryInt(t, db,
		`SELECT is_published FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))

	if err := SetEventInPulje(ctx, db, logger, "e1", string(models.PuljeFredagKveld), true); err != nil {
		t.Fatalf("re-add: %v", err)
	}
	gotPublishedAfterReadd := testutil.QueryInt(t, db,
		`SELECT is_published FROM relation_event_puljer WHERE event_id='e1' AND pulje_id=?`,
		string(models.PuljeFredagKveld))

	// Then
	if gotPublishedAfterRemoval != expectedPublished {
		t.Fatalf("after real removal is_published = %d, want %d", gotPublishedAfterRemoval, expectedPublished)
	}
	if gotPublishedAfterReadd != expectedPublished {
		t.Fatalf("after genuine re-add is_published = %d, want %d", gotPublishedAfterReadd, expectedPublished)
	}
}

package eventservice

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Regncon/conorganizer/service"
	_ "modernc.org/sqlite"
)

func TestGetPreviousNext(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	imgDir := "" // GetEventImageUrl -> placeholder -> code blanks it

	// One DB for the whole test
	db := mustInitTestDB(t)
	defer db.Close()

	// Seed once (deterministic ordering). We clear events once here, not per subtest.
	mustExec(t, db, `DELETE FROM events;`)
	mustExec(t, db, `
		INSERT INTO events (
			id,
			title,
			intro,
			description,
			image_url,
			host_name,
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			status,
			inserted_time
		) VALUES
		('e1','Old','intro e1','desc e1','/img1','Host One','one@test.test','11111111',4,1,1,'Godkjent','2025-10-01 10:00:00'),
		('e2','Mid','intro e2','desc e2','',     'Host Two','two@test.test','22222222',5,0,1,'Innsendt','2025-10-02 10:00:00'),
		('e3','New','intro e3','desc e3','/img3','Host Tre','tre@test.test','33333333',6,1,0,'Godkjent','2025-10-03 10:00:00'),
		-- excluded row: not in ('Innsendt','Godkjent'); use 'Kladd' to satisfy FK
		('e4','KladdRow','intro e4','desc e4','', 'Host Four','four@test.test','44444444',3,0,0,'Kladd','2025-10-04 10:00:00')
	`)

	// Subtests run sequentially (no t.Parallel)
	tests := []struct {
		name        string
		currentID   string
		wantPrevID  string
		wantPrevTit string
		wantNextID  string
		wantNextTit string
	}{
		{"middle_has_both_neighbors", "e2", "e3", "New", "e1", "Old"},
		{"first_has_next_only", "e3", "", "", "e2", "Mid"},
		{"last_has_prev_only", "e1", "e2", "Mid", "", ""},
		{"excluded_status_returns_empty_neighbors", "e4", "", "", "", ""},
		{"missing_id_returns_empty_neighbors", "does-not-exist", "", "", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetPreviousNext(ctx, db, logger, tc.currentID, &imgDir)
			if err != nil {
				t.Fatalf("GetPreviousNext() error = %v", err)
			}

			if got.PreviousUrl != tc.wantPrevID {
				t.Errorf("PreviousUrl = %q, want %q", got.PreviousUrl, tc.wantPrevID)
			}
			if got.PreviousTitle != tc.wantPrevTit {
				t.Errorf("PreviousTitle = %q, want %q", got.PreviousTitle, tc.wantPrevTit)
			}
			if got.NextUrl != tc.wantNextID {
				t.Errorf("NextUrl = %q, want %q", got.NextUrl, tc.wantNextID)
			}
			if got.NextTitle != tc.wantNextTit {
				t.Errorf("NextTitle = %q, want %q", got.NextTitle, tc.wantNextTit)
			}

			// With no banner files present, URLs should be blanked by the function
			if got.PreviousImageURL != "" {
				t.Errorf("PreviousImageURL = %q, want empty", got.PreviousImageURL)
			}
			if got.NextImageURL != "" {
				t.Errorf("NextImageURL = %q, want empty", got.NextImageURL)
			}
		})
	}
}

func mustInitTestDB(t *testing.T) *sql.DB {
	t.Helper()

	uniqueDatabaseName := "test_prevnext_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
		_ = os.Remove(testDBPath)
	})
	return db
}

func mustExec(t *testing.T, db *sql.DB, q string, args ...any) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := db.ExecContext(ctx, q, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, q)
	}
}

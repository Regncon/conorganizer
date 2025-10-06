package eventservice

import (
	"context"
	"io"
	"log/slog"
	_ "modernc.org/sqlite"
	"testing"
)

// TestGetPreviousNextByPulje verifies pulje-aware navigation with and without admin visibility.
func TestGetPreviousNextByPulje(t *testing.T) {
	// ========== Arrange ==========
	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError}))
	imgDir := "" // no real images -> placeholder -> code should blank URLs

	db := mustInitTestDB(t)
	defer db.Close()

	// Make test deterministic: clear link table first, then events.
	mustExec(t, db, `DELETE FROM event_puljer;`)
	mustExec(t, db, `DELETE FROM events;`)
	// Puljer: ensure the two we need exist (INSERT OR REPLACE is safe for re-runs)
	mustExec(t, db, `
		INSERT OR REPLACE INTO puljer (id, name, start_time, end_time) VALUES
		('FredagKveld',  'Fredag kveld',  '2025-10-10 18:00:00', '2025-10-10 22:00:00'),
		('LordagMorgen', 'Lørdag morgen', '2025-10-11 09:00:00', '2025-10-11 13:00:00')
	`)

	// Events:
	// Ordering inside each pulje should follow inserted_time DESC, id DESC (typical pattern).
	// FredagKveld: e2 (19:00, published), e6 (18:30, UNPUBLISHED), e1 (18:00, published)
	// LørdagMorgen: e4 (11:00, UNPUBLISHED), e3 (10:00, published)
	// e5 is Kladd (ignored by status filter).
	mustExec(t, db, `
		INSERT INTO events (
			id, title, intro, description, image_url,
			host_name, email, phone_number, max_players,
			beginner_friendly, can_be_run_in_english,
			status, inserted_time
		) VALUES
		('e1','F-Old','intro e1','desc e1','/img1','Host One','one@test.test','11111111',4,1,1,'Godkjent','2025-10-10 18:00:00'),
		('e2','F-Newest','intro e2','desc e2','',   'Host Two','two@test.test','22222222',5,0,1,'Godkjent','2025-10-10 19:00:00'),
		('e3','L-Pub','intro e3','desc e3','/img3','Host Tre','tre@test.test','33333333',6,1,0,'Godkjent','2025-10-11 10:00:00'),
		('e4','L-Unpub','intro e4','desc e4','',    'Host Four','four@test.test','44444444',3,0,0,'Godkjent','2025-10-11 11:00:00'),
		('e5','IgnoredByStatus','intro e5','desc e5','', 'Host Five','five@test.test','55555555',2,0,0,'Kladd','2025-10-12 10:00:00'),
		('e6','F-Unpub','intro e6','desc e6','',    'Host Six','six@test.test','66666666',2,0,0,'Godkjent','2025-10-10 18:30:00')
	`)

	// Link into puljer with publication flags.
	mustExec(t, db, `
		INSERT INTO event_puljer (event_id, pulje_id, is_active, is_published, room) VALUES
		('e1','FredagKveld',  1, 1, ''),
		('e2','FredagKveld',  1, 1, ''),
		('e6','FredagKveld',  1, 0, ''),

		('e3','LordagMorgen', 1, 1, ''),
		('e4','LordagMorgen', 1, 0, ''),

		('e5','LordagMorgen', 1, 1, '') -- status Kladd in events -> ignored anyway
	`)

	type want struct {
		prevID  string
		prevTit string
		prevPul string
		nextID  string
		nextTit string
		nextPul string
	}

	// Cases for regular user (isAdmin=false): only published in each pulje, only status Godkjent
	userCases := []struct {
		name      string
		currentID string
		want      want
	}{
		{
			name:      "user:e2_first_in_FredagKveld_has_next_only_within_pulje",
			currentID: "e2", // first published in FredagKveld (e6 is unpublished so ignored)
			want: want{
				prevID: "", prevTit: "", prevPul: "",
				nextID: "e1", nextTit: "F-Old", nextPul: "FredagKveld",
			},
		},
		{
			name:      "user:e1_last_in_FredagKveld_moves_to_first_published_in_next_pulje",
			currentID: "e1", // last published in FredagKveld
			want: want{
				prevID: "e2", prevTit: "F-Newest", prevPul: "FredagKveld",
				nextID: "e3", nextTit: "L-Pub", nextPul: "LordagMorgen", // skips e4 (unpublished)
			},
		},
		{
			name:      "user:e3_first_in_LordagMorgen_has_no_previous_and_no_next_if_also_last",
			currentID: "e3", // in published list for LordagMorgen, it's both first and last (e4 unpublished)
			want: want{
				prevID: "", prevTit: "", prevPul: "",
				nextID: "", nextTit: "", nextPul: "",
			},
		},
	}

	// Cases for admin (isAdmin=true): can see unpublished, but still only status Godkjent
	adminCases := []struct {
		name      string
		currentID string
		want      want
	}{
		{
			name:      "admin:e2_first_in_FredagKveld_next_is_unpublished_e6",
			currentID: "e2",
			want: want{
				prevID: "", prevTit: "", prevPul: "",
				nextID: "e6", nextTit: "F-Unpub", nextPul: "FredagKveld",
			},
		},
		{
			name:      "admin:e6_middle_in_FredagKveld_prev_e2_next_e1",
			currentID: "e6",
			want: want{
				prevID: "e2", prevTit: "F-Newest", prevPul: "FredagKveld",
				nextID: "e1", nextTit: "F-Old", nextPul: "FredagKveld",
			},
		},
		{
			name:      "admin:e1_last_in_FredagKveld_moves_to_first_in_next_pulje_including_unpublished",
			currentID: "e1",
			want: want{
				prevID: "e6", prevTit: "F-Unpub", prevPul: "FredagKveld",
				nextID: "e4", nextTit: "L-Unpub", nextPul: "LordagMorgen", // includes unpublished for admin
			},
		},
		{
			name:      "admin:e4_first_in_LordagMorgen_has_no_previous_no_wrap_to_FredagKveld",
			currentID: "e4",
			want: want{
				prevID: "", prevTit: "", prevPul: "",
				nextID: "e3", nextTit: "L-Pub", nextPul: "LordagMorgen",
			},
		},
	}

	// ========== Act + Assert (user) ==========
	for _, tc := range userCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetPreviousNextByPulje(ctx, db, logger, tc.currentID, false /* isAdmin */, &imgDir)
			if err != nil {
				t.Fatalf("GetPreviousNextByPulje() error = %v", err)
			}

			// ========== Assert ==========
			if got.PreviousUrl != tc.want.prevID {
				t.Errorf("PreviousUrl = %q, want %q", got.PreviousUrl, tc.want.prevID)
			}
			if got.PreviousTitle != tc.want.prevTit {
				t.Errorf("PreviousTitle = %q, want %q", got.PreviousTitle, tc.want.prevTit)
			}
			if string(got.PreviousPulje) != tc.want.prevPul {
				t.Errorf("PreviousPulje = %q, want %q", got.PreviousPulje, tc.want.prevPul)
			}

			if got.NextUrl != tc.want.nextID {
				t.Errorf("NextUrl = %q, want %q", got.NextUrl, tc.want.nextID)
			}
			if got.NextTitle != tc.want.nextTit {
				t.Errorf("NextTitle = %q, want %q", got.NextTitle, tc.want.nextTit)
			}
			if string(got.NextPulje) != tc.want.nextPul {
				t.Errorf("NextPulje = %q, want %q", got.NextPulje, tc.want.nextPul)
			}

			// No images present -> should be blanked by the function
			if got.PreviousImageURL != "" {
				t.Errorf("PreviousImageURL = %q, want empty", got.PreviousImageURL)
			}
			if got.NextImageURL != "" {
				t.Errorf("NextImageURL = %q, want empty", got.NextImageURL)
			}
		})
	}

	// ========== Act + Assert (admin) ==========
	for _, tc := range adminCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetPreviousNextByPulje(ctx, db, logger, tc.currentID, true /* isAdmin */, &imgDir)
			if err != nil {
				t.Fatalf("GetPreviousNextByPulje() error = %v", err)
			}

			// ========== Assert ==========
			if got.PreviousUrl != tc.want.prevID {
				t.Errorf("PreviousUrl = %q, want %q", got.PreviousUrl, tc.want.prevID)
			}
			if got.PreviousTitle != tc.want.prevTit {
				t.Errorf("PreviousTitle = %q, want %q", got.PreviousTitle, tc.want.prevTit)
			}
			if string(got.PreviousPulje) != tc.want.prevPul {
				t.Errorf("PreviousPulje = %q, want %q", got.PreviousPulje, tc.want.prevPul)
			}

			if got.NextUrl != tc.want.nextID {
				t.Errorf("NextUrl = %q, want %q", got.NextUrl, tc.want.nextID)
			}
			if got.NextTitle != tc.want.nextTit {
				t.Errorf("NextTitle = %q, want %q", got.NextTitle, tc.want.nextTit)
			}
			if string(got.NextPulje) != tc.want.nextPul {
				t.Errorf("NextPulje = %q, want %q", got.NextPulje, tc.want.nextPul)
			}

			if got.PreviousImageURL != "" {
				t.Errorf("PreviousImageURL = %q, want empty", got.PreviousImageURL)
			}
			if got.NextImageURL != "" {
				t.Errorf("NextImageURL = %q, want empty", got.NextImageURL)
			}
		})
	}
}

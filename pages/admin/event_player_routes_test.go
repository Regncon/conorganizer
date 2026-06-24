package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

type fakeEventPlayerBroadcaster struct {
	err     error
	buckets []live.Bucket
}

func (b *fakeEventPlayerBroadcaster) Broadcast(ctx context.Context, buckets ...live.Bucket) error {
	b.buckets = append(b.buckets, buckets...)
	return b.err
}

func TestEventPlayerUpdateStatus_BroadcastsInterestsAfterSuccessfulAssignment(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/update_status", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
		"assignmentIsPlayer":        true,
		"assignmentIsGm":            false,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
}

func TestEventPlayerUpdateStatus_CreatesPlayerAssignmentFromSearchAndBroadcasts(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerUnassignedFixture(t, db, 1, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/update_status", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
		"assignmentIsPlayer":        true,
		"assignmentIsGm":            false,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
	assertAdminEventPlayerRole(t, db, 1, models.EventPlayerRolePlayer)
}

func TestEventPlayerUpdateStatus_WhenMissingBillettholder_ReturnsBadRequestAndDoesNotBroadcast(t *testing.T) {
	tests := []struct {
		name    string
		signals map[string]any
	}{
		{
			name: "zero billettholder",
			signals: map[string]any{
				"assignmentEventId":         "event-1",
				"assignmentPuljeId":         string(models.PuljeFredagKveld),
				"assignmentBillettholderId": 0,
				"assignmentIsPlayer":        true,
				"assignmentIsGm":            false,
			},
		},
		{
			name: "omitted billettholder",
			signals: map[string]any{
				"assignmentEventId":  "event-1",
				"assignmentPuljeId":  string(models.PuljeFredagKveld),
				"assignmentIsPlayer": true,
				"assignmentIsGm":     false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, router, broadcaster := setupEventPlayerRouteTest(t)

			recorder := putAdminEventPlayerSignals(t, router, "/update_status", test.signals)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
			}
			if len(broadcaster.buckets) != 0 {
				t.Fatalf("expected no broadcasts, got %v", broadcaster.buckets)
			}
		})
	}
}

func TestEventPlayerAddGM_BroadcastsInterestsAfterSuccessfulAssignment(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := postAdminEventPlayerSignals(t, router, "/post/add_gm", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
}

func TestEventPlayerFirstChoice_BroadcastsInterestsAfterSuccessfulMutation(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/first-choice", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
		"assignmentFirstChoice":     true,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
}

func TestEventPlayerFirstChoice_WhenMissingFirstChoiceSignal_ReturnsBadRequestAndDoesNotBroadcast(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/first-choice", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
	})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusBadRequest, recorder.Code, recorder.Body.String())
	}
	if len(broadcaster.buckets) != 0 {
		t.Fatalf("expected no broadcasts, got %v", broadcaster.buckets)
	}
}

func TestEventPlayerFirstChoice_WhenBroadcastFails_ReturnsServerError(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	broadcaster.err = errors.New("broadcast unavailable")
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/first-choice", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
		"assignmentFirstChoice":     true,
	})

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusInternalServerError, recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
}

func TestEventPlayerFirstChoice_WhenMutationFails_DoesNotBroadcast(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRoleGM, models.InterestLevelHigh)

	recorder := putAdminEventPlayerSignals(t, router, "/first-choice", map[string]any{
		"assignmentEventId":         "event-1",
		"assignmentPuljeId":         string(models.PuljeFredagKveld),
		"assignmentBillettholderId": 1,
		"assignmentFirstChoice":     true,
	})

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusConflict, recorder.Code, recorder.Body.String())
	}
	if len(broadcaster.buckets) != 0 {
		t.Fatalf("expected no broadcasts, got %v", broadcaster.buckets)
	}
}

func setupEventPlayerRouteTest(t testing.TB) (*sql.DB, chi.Router, *fakeEventPlayerBroadcaster) {
	t.Helper()

	db, logger := testutil.CreateTestDBAndLogger(t, "admin_event_player_routes")
	router := chi.NewRouter()
	broadcaster := &fakeEventPlayerBroadcaster{}
	setupEventPlayerRoutes(router, db, logger, logger, broadcaster)
	return db, router, broadcaster
}

func putAdminEventPlayerSignals(t testing.TB, router http.Handler, path string, signals map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	return adminEventPlayerSignals(t, router, http.MethodPut, path, signals)
}

func postAdminEventPlayerSignals(t testing.TB, router http.Handler, path string, signals map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	return adminEventPlayerSignals(t, router, http.MethodPost, path, signals)
}

func adminEventPlayerSignals(t testing.TB, router http.Handler, method string, path string, signals map[string]any) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(signals)
	if err != nil {
		t.Fatalf("failed to marshal Datastar signals: %v", err)
	}
	request := httptest.NewRequest(method, path, strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func seedAdminEventPlayerFixture(
	t testing.TB,
	db *sql.DB,
	billettholderID int,
	role models.EventPlayerRole,
	interest models.InterestLevel,
) {
	t.Helper()

	seedAdminEventPlayerUnassignedFixture(t, db, billettholderID, interest)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role)
		VALUES ('event-1', ?, ?, ?)
	`, models.PuljeFredagKveld, billettholderID, role)
}

func seedAdminEventPlayerUnassignedFixture(
	t testing.TB,
	db *sql.DB,
	billettholderID int,
	interest models.InterestLevel,
) {
	t.Helper()

	seedFirstChoiceLookupsForAdminRoute(t, db)
	testutil.MustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-09-04T18:00:00Z', '2026-09-04T23:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusOpen)
	testutil.MustExec(t, db, `
		INSERT INTO events(
			id, title, intro, description, system, event_type, age_group, event_runtime,
			host_name, email, phone_number, max_players, beginner_friendly,
			can_be_run_in_english, status
		)
		VALUES ('event-1', 'Event 1', 'Intro', 'Description', 'System', ?, ?, ?, 'Host',
			'host@example.com', '12345678', 4, 1, 1, ?)
	`, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved)
	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES ('event-1', ?, 1, 1)
	`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere(
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		)
		VALUES (?, 'Test', 'Participant', ?, 'Festivalpass', 1, ?, ?)
	`, billettholderID, 1000+billettholderID, 2000+billettholderID, 3000+billettholderID)
	testutil.MustExec(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, 'event-1', ?, ?)
	`, billettholderID, models.PuljeFredagKveld, interest)
}

func seedFirstChoiceLookupsForAdminRoute(t testing.TB, db *sql.DB) {
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

func assertBroadcastedInterests(t testing.TB, broadcaster *fakeEventPlayerBroadcaster) {
	t.Helper()

	if len(broadcaster.buckets) != 1 {
		t.Fatalf("expected one broadcast, got %v", broadcaster.buckets)
	}
	if broadcaster.buckets[0] != live.BucketInterests {
		t.Fatalf("expected broadcast bucket %q, got %q", live.BucketInterests, broadcaster.buckets[0])
	}
}

func assertAdminEventPlayerRole(t testing.TB, db *sql.DB, billettholderID int, want models.EventPlayerRole) {
	t.Helper()

	var got models.EventPlayerRole
	if err := db.QueryRow(`
		SELECT role
		FROM relation_events_players
		WHERE event_id = 'event-1'
			AND pulje_id = ?
			AND billettholder_id = ?
	`, models.PuljeFredagKveld, billettholderID).Scan(&got); err != nil {
		t.Fatalf("query event player role: %v", err)
	}
	if got != want {
		t.Fatalf("expected role %q, got %q", want, got)
	}
}

# Puljefordeling Event First-Choice Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Complete event puljefordeling first-choice handling by deriving first-choice status from existing assignment and interest data, splitting assignment from first-choice mutation, and verifying interest broadcasts.

**Final clarification from review:** first-choice is festival-wide. Any qualifying first-choice outside the exact current `(event_id, pulje_id)` row counts as already used, including another event in the same pulje. The code keeps the initial `HasOtherPuljeFirstChoice` name, but the UI copy is generic.

**Architecture:** Add a focused first-choice service in `service/puljefordeling` that owns the domain rule and batch status query. Keep `components/formsubmission/who_is_interested.templ` responsible for rendering and row assembly, using service-derived status instead of embedding first-choice SQL. Extract admin event-player route wiring enough to test broadcasts with a fake broadcaster.

**Tech Stack:** Go, SQLite, templ, Datastar signals, Chi routes, `log/slog`, existing `testutil` database helpers, existing `service/live` buckets.

---

## File Structure

- Create `service/puljefordeling/first_choice.go`: first-choice status query and first-choice interest mutation.
- Create `service/puljefordeling/first_choice_test.go`: service and mutation tests.
- Modify `components/formsubmission/who_is_interested.templ`: remove component-local first-choice SQL, apply service statuses, split first-choice controls from assignment controls, disable controls for GM/other-pulje first-choice.
- Modify `components/formsubmission/who_is_interested_test.go`: replace old SQL-first-choice expectations with row status application helper expectations.
- Modify `components/formsubmission/who_is_interested_test_helpers_test.go`: remove old `FirstChoice` helper assertions that become unused.
- Modify `pages/admin/admin.go`: extract event-player route setup, add first-choice route, route assignment and first-choice mutations through separate functions, broadcast `live.BucketInterests` after successful mutations.
- Create `pages/admin/event_player_routes_test.go`: route-level tests for broadcast success/failure and first-choice endpoint behavior.
- Generate templ output after changing `.templ` files with `go tool task build:templ`.

Git note: this sandbox currently cannot write `.git/index.lock`. Keep the commit steps in the plan for a normal writable checkout. If `.git` is still read-only during execution, run the verification commands and report that commits were skipped because staging failed with a read-only `.git`.

---

### Task 1: First-Choice Status Service

**Files:**
- Create: `service/puljefordeling/first_choice.go`
- Create: `service/puljefordeling/first_choice_test.go`

- [ ] **Step 1: Write failing service status tests**

Create `service/puljefordeling/first_choice_test.go` with tests that seed puljer, events, billettholdere, interests, and assignments, then assert status for the current event rows.

```go
package puljefordeling

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestGetFirstChoiceStatusesForEvent_DerivesCurrentOtherAndGMIgnored(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_status")
	seedFirstChoiceLookups(t, db)
	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagMorgen, "Lordag morgen")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "friday-gm", "Friday GM", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-choice", "Saturday Choice", models.PuljeLordagMorgen)
	seedFirstChoiceEvent(t, db, "saturday-evening", "Saturday Evening", models.PuljeLordagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "GMAndPlayer")
	seedFirstChoiceBillettholder(t, db, 2, "Bob", "CurrentChoice")
	seedFirstChoiceBillettholder(t, db, 3, "Cara", "Medium")

	seedFirstChoiceInterest(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.EventPlayerRoleGM)
	seedFirstChoiceInterest(t, db, 1, "saturday-choice", models.PuljeLordagMorgen, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-choice", models.PuljeLordagMorgen, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	seedFirstChoiceInterest(t, db, 2, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 2, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	seedFirstChoiceInterest(t, db, 3, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 3, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-evening")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent returned error: %v", err)
	}

	alice := statuses[FirstChoiceKey{BillettholderID: 1, EventID: "saturday-evening", PuljeID: string(models.PuljeLordagKveld)}]
	if !alice.HasCurrentPuljeFirstChoice {
		t.Fatalf("Alice should have current first-choice as player in Saturday evening")
	}
	if !alice.HasOtherPuljeFirstChoice {
		t.Fatalf("Alice should have other-pulje first-choice from Saturday morning")
	}

	bob := statuses[FirstChoiceKey{BillettholderID: 2, EventID: "saturday-evening", PuljeID: string(models.PuljeLordagKveld)}]
	if !bob.HasCurrentPuljeFirstChoice {
		t.Fatalf("Bob should have current first-choice")
	}
	if bob.HasOtherPuljeFirstChoice {
		t.Fatalf("Bob should not have other-pulje first-choice")
	}

	cara := statuses[FirstChoiceKey{BillettholderID: 3, EventID: "saturday-evening", PuljeID: string(models.PuljeLordagKveld)}]
	if cara.HasCurrentPuljeFirstChoice || cara.HasOtherPuljeFirstChoice {
		t.Fatalf("Cara has medium interest and should not have first-choice status: %+v", cara)
	}
}

func TestGetFirstChoiceStatusesForEvent_GMHighInterestDoesNotCountAsOtherPuljeFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "first_choice_gm_ignored")
	seedFirstChoiceLookups(t, db)
	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "friday-gm", "Friday GM", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-evening", "Saturday Evening", models.PuljeLordagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "GMOnly")
	seedFirstChoiceInterest(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-gm", models.PuljeFredagKveld, models.EventPlayerRoleGM)
	seedFirstChoiceInterest(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "saturday-evening", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	statuses, err := GetFirstChoiceStatusesForEvent(db, "saturday-evening")
	if err != nil {
		t.Fatalf("GetFirstChoiceStatusesForEvent returned error: %v", err)
	}

	got := statuses[FirstChoiceKey{BillettholderID: 1, EventID: "saturday-evening", PuljeID: string(models.PuljeLordagKveld)}]
	if got.HasOtherPuljeFirstChoice {
		t.Fatalf("GM high interest in another pulje must not count as other-pulje first-choice")
	}
}

func seedFirstChoiceLookups(t testing.TB, db *sql.DB) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusOpen)
}

func seedFirstChoicePulje(t testing.TB, db *sql.DB, pulje models.Pulje, name string) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES (?, ?, ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z')
	`, pulje, name, models.PuljeStatusOpen)
}

func seedFirstChoiceEvent(t testing.TB, db *sql.DB, eventID string, title string, pulje models.Pulje) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type, age_group, event_runtime,
			host_name, email, phone_number, max_players, beginner_friendly,
			can_be_run_in_english, status
		) VALUES (?, ?, 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`, eventID, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved)
	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
	`, eventID, pulje)
}

func seedFirstChoiceBillettholder(t testing.TB, db *sql.DB, id int, firstName string, lastName string) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, ?, ?, 1, 'Festivalpass', 1, ?, ?)
	`, id, firstName, lastName, 1000+id, 2000+id)
}

func seedFirstChoiceInterest(t testing.TB, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje, level models.InterestLevel) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
	`, billettholderID, eventID, pulje, level)
}

func seedFirstChoiceAssignment(t testing.TB, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje, role models.EventPlayerRole) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role)
		VALUES (?, ?, ?, ?)
	`, eventID, pulje, billettholderID, role)
}
```

- [ ] **Step 2: Run service status tests and verify they fail**

Run:

```bash
go test ./service/puljefordeling -run 'TestGetFirstChoiceStatusesForEvent' -count=1
```

Expected: FAIL with compiler errors for undefined `GetFirstChoiceStatusesForEvent`, `FirstChoiceKey`, and `FirstChoiceStatus`.

- [ ] **Step 3: Implement the status service**

Create `service/puljefordeling/first_choice.go` with this implementation:

```go
package puljefordeling

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
)

type FirstChoiceKey struct {
	BillettholderID int
	EventID         string
	PuljeID         string
}

type FirstChoiceStatus struct {
	HasCurrentPuljeFirstChoice bool
	HasOtherPuljeFirstChoice   bool
}

func GetFirstChoiceStatusesForEvent(db *sql.DB, eventID string) (map[FirstChoiceKey]FirstChoiceStatus, error) {
	const query = `
		WITH current_rows AS (
			SELECT billettholder_id, event_id, pulje_id
			FROM interests
			WHERE event_id = ?

			UNION

			SELECT billettholder_id, event_id, pulje_id
			FROM relation_events_players
			WHERE event_id = ?
		),
		qualifying_first_choices AS (
			SELECT
				ep.billettholder_id,
				ep.event_id,
				ep.pulje_id
			FROM relation_events_players AS ep
			INNER JOIN interests AS i
				ON i.billettholder_id = ep.billettholder_id
				AND i.event_id = ep.event_id
				AND i.pulje_id = ep.pulje_id
			WHERE ep.role = ?
				AND i.interest_level = ?
		)
		SELECT
			cr.billettholder_id,
			cr.event_id,
			cr.pulje_id,
			EXISTS (
				SELECT 1
				FROM qualifying_first_choices AS q
				WHERE q.billettholder_id = cr.billettholder_id
					AND q.event_id = cr.event_id
					AND q.pulje_id = cr.pulje_id
			) AS has_current_pulje_first_choice,
			EXISTS (
				SELECT 1
				FROM qualifying_first_choices AS q
				WHERE q.billettholder_id = cr.billettholder_id
					AND q.pulje_id <> cr.pulje_id
			) AS has_other_pulje_first_choice
		FROM current_rows AS cr
	`

	rows, err := db.Query(query, eventID, eventID, models.EventPlayerRolePlayer, models.InterestLevelHigh)
	if err != nil {
		return nil, fmt.Errorf("query first-choice statuses for event %s: %w", eventID, err)
	}
	defer rows.Close()

	statuses := make(map[FirstChoiceKey]FirstChoiceStatus)
	for rows.Next() {
		var key FirstChoiceKey
		var status FirstChoiceStatus
		if err := rows.Scan(
			&key.BillettholderID,
			&key.EventID,
			&key.PuljeID,
			&status.HasCurrentPuljeFirstChoice,
			&status.HasOtherPuljeFirstChoice,
		); err != nil {
			return nil, fmt.Errorf("scan first-choice status for event %s: %w", eventID, err)
		}
		statuses[key] = status
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate first-choice statuses for event %s: %w", eventID, err)
	}

	return statuses, nil
}
```

- [ ] **Step 4: Run service status tests and verify they pass**

Run:

```bash
go test ./service/puljefordeling -run 'TestGetFirstChoiceStatusesForEvent' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit the first-choice status service**

Run:

```bash
git add service/puljefordeling/first_choice.go service/puljefordeling/first_choice_test.go
git commit -m "feat: derive puljefordeling first-choice status"
```

Expected: commit succeeds in a writable checkout. In this sandbox, expect `git add` to fail if `.git` remains read-only; record that and continue without destructive workaround.

---

### Task 2: First-Choice Mutation Service

**Files:**
- Modify: `service/puljefordeling/first_choice.go`
- Modify: `service/puljefordeling/first_choice_test.go`

- [ ] **Step 1: Write failing mutation tests**

Append these tests to `service/puljefordeling/first_choice_test.go`:

```go
func TestSetAssignmentFirstChoice_SetsAndRemovesInterestWithoutChangingAssignment(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice")
	seedFirstChoiceLookups(t, db)
	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "event-1", "Event 1", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Player")
	seedFirstChoiceInterest(t, db, 1, "event-1", models.PuljeFredagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "event-1", models.PuljeFredagKveld, models.EventPlayerRolePlayer)

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, true); err != nil {
		t.Fatalf("set first-choice returned error: %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelHigh {
		t.Fatalf("set first-choice interest mismatch: want %s, got %s", models.InterestLevelHigh, got)
	}
	if got := queryFirstChoiceAssignmentRole(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.EventPlayerRolePlayer {
		t.Fatalf("assignment role changed after setting first-choice: got %s", got)
	}

	if err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, false); err != nil {
		t.Fatalf("remove first-choice returned error: %v", err)
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelMedium {
		t.Fatalf("remove first-choice interest mismatch: want %s, got %s", models.InterestLevelMedium, got)
	}
}

func TestSetAssignmentFirstChoice_RejectsGMAssignment(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_gm")
	seedFirstChoiceLookups(t, db)
	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoiceEvent(t, db, "event-1", "Event 1", models.PuljeFredagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "GM")
	seedFirstChoiceInterest(t, db, 1, "event-1", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "event-1", models.PuljeFredagKveld, models.EventPlayerRoleGM)

	err := SetAssignmentFirstChoice(db, "event-1", string(models.PuljeFredagKveld), 1, true)
	if err == nil {
		t.Fatalf("expected GM assignment to reject first-choice mutation")
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "event-1", models.PuljeFredagKveld); got != models.InterestLevelHigh {
		t.Fatalf("GM rejection should not change interest: got %s", got)
	}
}

func TestSetAssignmentFirstChoice_RejectsOtherPuljeFirstChoice(t *testing.T) {
	db := testutil.CreateTestDB(t, "set_first_choice_other")
	seedFirstChoiceLookups(t, db)
	seedFirstChoicePulje(t, db, models.PuljeFredagKveld, "Fredag kveld")
	seedFirstChoicePulje(t, db, models.PuljeLordagKveld, "Lordag kveld")
	seedFirstChoiceEvent(t, db, "friday-choice", "Friday Choice", models.PuljeFredagKveld)
	seedFirstChoiceEvent(t, db, "saturday-choice", "Saturday Choice", models.PuljeLordagKveld)
	seedFirstChoiceBillettholder(t, db, 1, "Alice", "Player")
	seedFirstChoiceInterest(t, db, 1, "friday-choice", models.PuljeFredagKveld, models.InterestLevelHigh)
	seedFirstChoiceAssignment(t, db, 1, "friday-choice", models.PuljeFredagKveld, models.EventPlayerRolePlayer)
	seedFirstChoiceInterest(t, db, 1, "saturday-choice", models.PuljeLordagKveld, models.InterestLevelMedium)
	seedFirstChoiceAssignment(t, db, 1, "saturday-choice", models.PuljeLordagKveld, models.EventPlayerRolePlayer)

	err := SetAssignmentFirstChoice(db, "saturday-choice", string(models.PuljeLordagKveld), 1, true)
	if err == nil {
		t.Fatalf("expected other-pulje first-choice to reject setting a second first-choice")
	}
	if got := queryFirstChoiceInterestLevel(t, db, 1, "saturday-choice", models.PuljeLordagKveld); got != models.InterestLevelMedium {
		t.Fatalf("rejected second first-choice should keep medium interest: got %s", got)
	}
}

func queryFirstChoiceInterestLevel(t testing.TB, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje) models.InterestLevel {
	t.Helper()
	var level models.InterestLevel
	if err := db.QueryRow(`
		SELECT interest_level
		FROM interests
		WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?
	`, billettholderID, eventID, pulje).Scan(&level); err != nil {
		t.Fatalf("query interest level: %v", err)
	}
	return level
}

func queryFirstChoiceAssignmentRole(t testing.TB, db *sql.DB, billettholderID int, eventID string, pulje models.Pulje) models.EventPlayerRole {
	t.Helper()
	var role models.EventPlayerRole
	if err := db.QueryRow(`
		SELECT role
		FROM relation_events_players
		WHERE billettholder_id = ? AND event_id = ? AND pulje_id = ?
	`, billettholderID, eventID, pulje).Scan(&role); err != nil {
		t.Fatalf("query assignment role: %v", err)
	}
	return role
}
```

- [ ] **Step 2: Run mutation tests and verify they fail**

Run:

```bash
go test ./service/puljefordeling -run 'TestSetAssignmentFirstChoice' -count=1
```

Expected: FAIL with compiler error `undefined: SetAssignmentFirstChoice`.

- [ ] **Step 3: Implement `SetAssignmentFirstChoice`**

Append this function to `service/puljefordeling/first_choice.go`:

```go
func SetAssignmentFirstChoice(db *sql.DB, eventID string, puljeID string, billettholderID int, enabled bool) error {
	if eventID == "" {
		return fmt.Errorf("event id is required")
	}
	if puljeID == "" {
		return fmt.Errorf("pulje id is required")
	}
	if billettholderID <= 0 {
		return fmt.Errorf("billettholder id must be greater than 0")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin first-choice mutation: %w", err)
	}
	defer tx.Rollback()

	var role models.EventPlayerRole
	if err := tx.QueryRow(`
		SELECT role
		FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
	`, eventID, puljeID, billettholderID).Scan(&role); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("cannot update first-choice without an assignment")
		}
		return fmt.Errorf("query assignment before first-choice mutation: %w", err)
	}
	if role == models.EventPlayerRoleGM {
		return fmt.Errorf("GM assignment cannot be first-choice")
	}

	if enabled {
		var otherCount int
		if err := tx.QueryRow(`
			SELECT COUNT(1)
			FROM relation_events_players AS ep
			INNER JOIN interests AS i
				ON i.billettholder_id = ep.billettholder_id
				AND i.event_id = ep.event_id
				AND i.pulje_id = ep.pulje_id
			WHERE ep.billettholder_id = ?
				AND ep.role = ?
				AND i.interest_level = ?
				AND ep.pulje_id <> ?
		`, billettholderID, models.EventPlayerRolePlayer, models.InterestLevelHigh, puljeID).Scan(&otherCount); err != nil {
			return fmt.Errorf("query other-pulje first-choice before mutation: %w", err)
		}
		if otherCount > 0 {
			return fmt.Errorf("billettholder already has first-choice in another pulje")
		}
	}

	level := models.InterestLevelMedium
	if enabled {
		level = models.InterestLevelHigh
	}

	if _, err := tx.Exec(`
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(billettholder_id, event_id, pulje_id) DO UPDATE SET
			interest_level = EXCLUDED.interest_level,
			updated_at = `+models.DBDateTimeNowSQL+`
	`, billettholderID, eventID, puljeID, level); err != nil {
		return fmt.Errorf("upsert first-choice interest: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit first-choice mutation: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Run service tests and verify they pass**

Run:

```bash
go test ./service/puljefordeling -run 'Test(GetFirstChoiceStatusesForEvent|SetAssignmentFirstChoice)' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit the mutation service**

Run:

```bash
git add service/puljefordeling/first_choice.go service/puljefordeling/first_choice_test.go
git commit -m "feat: update puljefordeling first-choice interest"
```

Expected: commit succeeds in a writable checkout. In this sandbox, expect `git add` to fail if `.git` remains read-only; record that and continue.

---

### Task 3: Admin Event-Player Routes And Broadcast Tests

**Files:**
- Modify: `pages/admin/admin.go`
- Create: `pages/admin/event_player_routes_test.go`

- [ ] **Step 1: Write failing route tests with a fake broadcaster**

Create `pages/admin/event_player_routes_test.go`:

```go
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

func (b *fakeEventPlayerBroadcaster) Broadcast(_ context.Context, buckets ...live.Bucket) error {
	b.buckets = append(b.buckets, buckets...)
	return b.err
}

func TestEventPlayerUpdateStatus_BroadcastsInterestsAfterSuccessfulAssignment(t *testing.T) {
	db, router, broadcaster := setupEventPlayerRouteTest(t)
	seedAdminEventPlayerFixture(t, db, 1, models.EventPlayerRolePlayer, models.InterestLevelMedium)

	recorder := putAdminEventPlayerSignals(t, router, "/update_status", map[string]any{
		"assignmentEventId":          "event-1",
		"assignmentPuljeId":          string(models.PuljeFredagKveld),
		"assignmentBillettholderId":  1,
		"assignmentIsPlayer":         true,
		"assignmentIsGm":             false,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
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
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	assertBroadcastedInterests(t, broadcaster)
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
		t.Fatalf("expected status 500, got %d", recorder.Code)
	}
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

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for GM first-choice mutation, got %d", recorder.Code)
	}
	if len(broadcaster.buckets) != 0 {
		t.Fatalf("expected no broadcast after failed mutation, got %v", broadcaster.buckets)
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
	body, err := json.Marshal(signals)
	if err != nil {
		t.Fatalf("marshal signals: %v", err)
	}
	request := httptest.NewRequest(http.MethodPut, path, strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func seedAdminEventPlayerFixture(t testing.TB, db *sql.DB, billettholderID int, role models.EventPlayerRole, interest models.InterestLevel) {
	t.Helper()
	seedFirstChoiceLookupsForAdminRoute(t, db)
	testutil.MustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusOpen)
	testutil.MustExec(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type, age_group, event_runtime,
			host_name, email, phone_number, max_players, beginner_friendly,
			can_be_run_in_english, status
		) VALUES ('event-1', 'Event 1', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved)
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere(id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id)
		VALUES (?, 'Route', 'Tester', 1, 'Festivalpass', 1, ?, ?)
	`, billettholderID, 1000+billettholderID, 2000+billettholderID)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role)
		VALUES ('event-1', ?, ?, ?)
	`, models.PuljeFredagKveld, billettholderID, role)
	testutil.MustExec(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, 'event-1', ?, ?)
	`, billettholderID, models.PuljeFredagKveld, interest)
}

func seedFirstChoiceLookupsForAdminRoute(t testing.TB, db *sql.DB) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusOpen)
}

func assertBroadcastedInterests(t testing.TB, broadcaster *fakeEventPlayerBroadcaster) {
	t.Helper()
	if len(broadcaster.buckets) != 1 || broadcaster.buckets[0] != live.BucketInterests {
		t.Fatalf("expected one interests broadcast, got %v", broadcaster.buckets)
	}
}
```

- [ ] **Step 2: Run route tests and verify they fail**

Run:

```bash
go test ./pages/admin -run 'TestEventPlayer' -count=1
```

Expected: FAIL with compiler error `undefined: setupEventPlayerRoutes`.

- [ ] **Step 3: Extract event-player route setup and add first-choice route**

Modify `pages/admin/admin.go`:

1. Add an interface near `SetupAdminRoute`:

```go
type eventPlayerBroadcaster interface {
	Broadcast(ctx context.Context, buckets ...live.Bucket) error
}
```

2. Replace the inline `apiRouter.Route("/event-players", func(eventPlayersRouter chi.Router) { ... })` block with:

```go
apiRouter.Route("/event-players", func(eventPlayersRouter chi.Router) {
	setupEventPlayerRoutes(eventPlayersRouter, db, logger, baseLogger, liveManager)
})
```

3. Add this helper below `SetupAdminRoute`:

```go
func setupEventPlayerRoutes(
	eventPlayersRouter chi.Router,
	db *sql.DB,
	logger *slog.Logger,
	baseLogger *slog.Logger,
	broadcaster eventPlayerBroadcaster,
) {
	eventPlayersRouter.Post("/post/add_gm", func(w http.ResponseWriter, r *http.Request) {
		type Store struct {
			BillettholderId int    `json:"assignmentBillettholderId"`
			EventId         string `json:"assignmentEventId"`
			PuljeId         string `json:"assignmentPuljeId"`
		}
		store := &Store{}

		if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
			http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
			return
		}
		if store.BillettholderId <= 0 {
			http.Error(w, fmt.Errorf("invalid assignmentBillettholderId %d: must be greater than 0", store.BillettholderId).Error(), http.StatusNotFound)
			return
		}

		if err := formsubmission.UpdatePlayerStatus(store.EventId, store.PuljeId, store.BillettholderId, false, true, db, baseLogger); err != nil {
			logger.Error(fmt.Errorf("failed to add player as GM: %w", err).Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := broadcaster.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast add GM update: %w", err).Error())
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
	})

	eventPlayersRouter.Put("/update_status", func(w http.ResponseWriter, r *http.Request) {
		type Store struct {
			BillettholderId int    `json:"assignmentBillettholderId"`
			EventId         string `json:"assignmentEventId"`
			PuljeId         string `json:"assignmentPuljeId"`
			IsPlayer        bool   `json:"assignmentIsPlayer"`
			IsGm            bool   `json:"assignmentIsGm"`
		}
		store := &Store{}

		if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
			http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
			return
		}
		if err := formsubmission.UpdatePlayerStatus(store.EventId, store.PuljeId, store.BillettholderId, store.IsPlayer, store.IsGm, db, baseLogger); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := broadcaster.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast player status update: %w", err).Error())
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
	})

	eventPlayersRouter.Put("/first-choice", func(w http.ResponseWriter, r *http.Request) {
		type Store struct {
			BillettholderId int    `json:"assignmentBillettholderId"`
			EventId         string `json:"assignmentEventId"`
			PuljeId         string `json:"assignmentPuljeId"`
			FirstChoice     bool   `json:"assignmentFirstChoice"`
		}
		store := &Store{}

		if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
			http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
			return
		}
		if err := puljefordeling.SetAssignmentFirstChoice(db, store.EventId, store.PuljeId, store.BillettholderId, store.FirstChoice); err != nil {
			logger.Error(fmt.Errorf("failed to update first-choice interest: %w", err).Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := broadcaster.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast first-choice update: %w", err).Error())
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
	})
}
```

4. Add this import alias to `pages/admin/admin.go`:

```go
puljefordeling "github.com/Regncon/conorganizer/service/puljefordeling"
```

5. Remove the old `/post/add_first_choice` route from `pages/admin/admin.go`.

- [ ] **Step 4: Run route tests and verify they pass**

Run:

```bash
go test ./pages/admin -run 'TestEventPlayer' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit route extraction and broadcast tests**

Run:

```bash
git add pages/admin/admin.go pages/admin/event_player_routes_test.go
git commit -m "test: verify event-player interest broadcasts"
```

Expected: commit succeeds in a writable checkout. In this sandbox, expect `git add` to fail if `.git` remains read-only; record that and continue.

---

### Task 4: Component Status Wiring And First-Choice Controls

**Files:**
- Modify: `components/formsubmission/who_is_interested.templ`
- Modify: `components/formsubmission/who_is_interested_test.go`
- Modify: `components/formsubmission/who_is_interested_test_helpers_test.go`

- [ ] **Step 1: Write failing component helper tests**

Replace the old `TestGetInterestsForEvent_FirstChoiceRules` in `components/formsubmission/who_is_interested_test.go` with focused helper tests:

```go
func TestFirstChoiceStatusText_PrefersOtherPuljeOverCurrentPulje(t *testing.T) {
	row := InterestWithHolder{
		HasCurrentPuljeFirstChoice: true,
		HasOtherPuljeFirstChoice:   true,
	}

	got := firstChoiceStatusText(row)

	if got != "Fått førsteval i anna pulje" {
		t.Fatalf("status text mismatch: %q", got)
	}
}

func TestFirstChoiceButtonAction_DisablesForGMAndOtherPuljeFirstChoice(t *testing.T) {
	setAction := "set"
	removeAction := "remove"

	gm := InterestWithHolder{IsGamemaster: true}
	gmButton := firstChoiceButtonAction(gm, setAction, removeAction)
	if !gmButton.Disabled {
		t.Fatalf("GM first-choice button should be disabled")
	}

	other := InterestWithHolder{HasOtherPuljeFirstChoice: true}
	otherButton := firstChoiceButtonAction(other, setAction, removeAction)
	if !otherButton.Disabled {
		t.Fatalf("other-pulje first-choice should disable setting first-choice")
	}

	current := InterestWithHolder{HasCurrentPuljeFirstChoice: true}
	currentButton := firstChoiceButtonAction(current, setAction, removeAction)
	if currentButton.Disabled || currentButton.Action != removeAction {
		t.Fatalf("current first-choice should allow removal, got %+v", currentButton)
	}
}

func TestApplyFirstChoiceStatusesToRows_AddsServiceStatusToRows(t *testing.T) {
	rows := []InterestWithHolder{{
		BillettholderId: 42,
		EventId:         "event-1",
		PuljeId:         string(puljeP1),
	}}
	statuses := map[puljefordeling.FirstChoiceKey]puljefordeling.FirstChoiceStatus{
		{BillettholderID: 42, EventID: "event-1", PuljeID: string(puljeP1)}: {
			HasCurrentPuljeFirstChoice: true,
			HasOtherPuljeFirstChoice:   false,
		},
	}

	applyFirstChoiceStatusesToRows(rows, statuses)

	if !rows[0].HasCurrentPuljeFirstChoice {
		t.Fatalf("expected current first-choice status to be applied")
	}
}
```

Add this import to `components/formsubmission/who_is_interested_test.go`:

```go
puljefordeling "github.com/Regncon/conorganizer/service/puljefordeling"
```

Remove `firstChoiceCase` and `expectFirstChoice` from `components/formsubmission/who_is_interested_test_helpers_test.go` after replacing the old SQL-first-choice tests.

- [ ] **Step 2: Run component tests and verify they fail**

Run:

```bash
go test ./components/formsubmission -run 'Test(FirstChoiceStatusText|FirstChoiceButtonAction|ApplyFirstChoiceStatusesToRows)' -count=1
```

Expected: FAIL with compiler errors for undefined `firstChoiceStatusText`, `firstChoiceButtonAction`, `applyFirstChoiceStatusesToRows`, and missing fields on `InterestWithHolder`.

- [ ] **Step 3: Update row types and remove component-local first-choice SQL**

In `components/formsubmission/who_is_interested.templ`:

1. Add the service import:

```go
puljefordeling "github.com/Regncon/conorganizer/service/puljefordeling"
```

2. Replace `FirstChoice bool` in `InterestWithHolder` with:

```go
HasCurrentPuljeFirstChoice bool
HasOtherPuljeFirstChoice   bool
```

3. Delete `queryFirstChoice`.

4. Remove `queryFirstChoice` from both `GetInterestsForEvent` and `GetAssigneesForEvent` SELECT lists.

5. Remove `&interest.FirstChoice` and `&assignment.FirstChoice` from both `rows.Scan` calls.

- [ ] **Step 4: Add status and button helpers**

Replace the existing `ButtonInfo` type and add these helpers near the existing button helper functions in `components/formsubmission/who_is_interested.templ`:

```go
type ButtonInfo struct {
	Label    string
	Action   string
	Disabled bool
	Title    string
}

func applyFirstChoiceStatusesToRows(rows []InterestWithHolder, statuses map[puljefordeling.FirstChoiceKey]puljefordeling.FirstChoiceStatus) {
	for index := range rows {
		key := puljefordeling.FirstChoiceKey{
			BillettholderID: rows[index].BillettholderId,
			EventID:         rows[index].EventId,
			PuljeID:         rows[index].PuljeId,
		}
		status := statuses[key]
		rows[index].HasCurrentPuljeFirstChoice = status.HasCurrentPuljeFirstChoice
		rows[index].HasOtherPuljeFirstChoice = status.HasOtherPuljeFirstChoice
	}
}

func firstChoiceStatusText(row InterestWithHolder) string {
	if row.HasOtherPuljeFirstChoice {
		return "Fått førsteval i anna pulje"
	}
	if row.HasCurrentPuljeFirstChoice {
		return "Fått førsteval"
	}
	return ""
}

func firstChoiceButtonAction(row InterestWithHolder, setAction string, removeAction string) ButtonInfo {
	if row.HasCurrentPuljeFirstChoice {
		return ButtonInfo{Label: "Fjern førsteval", Action: removeAction}
	}
	if row.IsGamemaster {
		return ButtonInfo{Label: "Set førsteval", Disabled: true, Title: "GM kan ikkje ha førsteval for dette arrangementet"}
	}
	if row.HasOtherPuljeFirstChoice {
		return ButtonInfo{Label: "Set førsteval", Disabled: true, Title: "Har allereie fått førsteval i anna pulje"}
	}
	return ButtonInfo{Label: "Set førsteval", Action: setAction}
}
```

Keep the existing `playerButtonAction` and `gmButtonAction`, updating their `ButtonInfo` literals to set only `Label` and `Action`.

- [ ] **Step 5: Apply service statuses in `WhoIsInterested`**

In `WhoIsInterested`, load statuses and apply them:

```go
{{ firstChoiceStatuses, firstChoiceStatusesErr := puljefordeling.GetFirstChoiceStatusesForEvent(db, eventId) }}
{{
	if firstChoiceStatusesErr == nil {
		applyFirstChoiceStatusesToRows(interests, firstChoiceStatuses)
		applyFirstChoiceStatusesToRows(assignees, firstChoiceStatuses)
	}
}}
```

Include `firstChoiceStatusesErr` in `loadErr`:

```go
if firstChoiceStatusesErr != nil {
	loadErr = firstChoiceStatusesErr
}
```

- [ ] **Step 6: Render first-choice status and controls**

Inside `interestInPulje`, add first-choice actions:

```go
setFirstChoiceAction := fmt.Sprintf("$assignmentEventId = '%s'; $assignmentPuljeId = '%s'; $assignmentBillettholderId = %d; $assignmentFirstChoice = true; @put('/admin/approval/api/event-players/first-choice')", eventId, puljeId, interest.BillettholderId)
removeFirstChoiceAction := fmt.Sprintf("$assignmentEventId = '%s'; $assignmentPuljeId = '%s'; $assignmentBillettholderId = %d; $assignmentFirstChoice = false; @put('/admin/approval/api/event-players/first-choice')", eventId, puljeId, interest.BillettholderId)
firstChoiceAction := firstChoiceButtonAction(interest, setFirstChoiceAction, removeFirstChoiceAction)
buttonsInfo := []ButtonInfo{playerAction, gmAction}
if isAssigned {
	buttonsInfo = append(buttonsInfo, firstChoiceAction)
}
```

Replace:

```templ
if interest.FirstChoice {
	<p class="first-choice">Fått førsteval</p>
}
```

with:

```templ
{{ firstChoiceText := firstChoiceStatusText(interest) }}
if firstChoiceText != "" {
	<p class="first-choice">{ firstChoiceText }</p>
}
```

Update the button render loop:

```templ
<button
	class={ actionButtonClass(buttonInfo.Label) }
	data-on:click={ buttonInfo.Action }
	disabled?={ buttonInfo.Disabled }
	title={ buttonInfo.Title }
>
	{ buttonInfo.Label }
</button>
```

Replace the search button currently labeled `Legg til som førsteval` with a player assignment button:

```templ
<button
	class="btn btn--secondary"
	data-on:click={ fmt.Sprintf("$assignmentPuljeId = '%s'; $assignmentIsPlayer = true; $assignmentIsGm = false; @put('/admin/approval/api/event-players/update_status'); evt.currentTarget.closest('.billettholder-search')?.querySelector('admin-billettholder-search')?.clearSearch(); $assignmentBillettholderId = 0; $clearInput = $clearInput + 1", puljeId) }
>
	Legg til som { models.EventPlayerRolePlayer.Label() }
</button>
```

- [ ] **Step 7: Run templ generation and component tests**

Run:

```bash
go tool task build:templ
go test ./components/formsubmission -count=1
```

Expected: PASS.

- [ ] **Step 8: Commit component UI wiring**

Run:

```bash
git add components/formsubmission/who_is_interested.templ components/formsubmission/who_is_interested_test.go components/formsubmission/who_is_interested_test_helpers_test.go
git add components/formsubmission/*_templ.go
git commit -m "feat: split event first-choice controls"
```

Expected: commit succeeds in a writable checkout. In this sandbox, expect `git add` to fail if `.git` remains read-only; record that and continue.

---

### Task 5: Full Verification And Cleanup

**Files:**
- Modify only files required by failing tests from prior tasks.

- [ ] **Step 1: Run focused package tests**

Run:

```bash
go test ./service/puljefordeling ./components/formsubmission ./pages/admin -count=1
```

Expected: PASS.

- [ ] **Step 2: Run full templ generation**

Run:

```bash
go tool task build:templ
```

Expected: command exits 0 and generated templ files are up to date.

- [ ] **Step 3: Run full test suite**

Run:

```bash
go test ./...
```

Expected: PASS.

- [ ] **Step 4: Inspect changed files**

Run:

```bash
git status --short
git diff --stat
```

Expected: changed files are limited to the service, component, admin route, tests, generated templ files, and the approved docs. The unrelated `.ai/threads/feature-complete-plujefordeling-event.md` remains untracked unless the user explicitly asks to include it.

- [ ] **Step 5: Final commit**

Run:

```bash
git add service/puljefordeling/first_choice.go service/puljefordeling/first_choice_test.go
git add components/formsubmission/who_is_interested.templ components/formsubmission/who_is_interested_test.go components/formsubmission/who_is_interested_test_helpers_test.go components/formsubmission/*_templ.go
git add pages/admin/admin.go pages/admin/event_player_routes_test.go pages/admin/*_templ.go
git commit -m "feat: complete event puljefordeling first-choice flow"
```

Expected: commit succeeds in a writable checkout. If previous task commits were made, this final commit should contain no changes and can be skipped. In this sandbox, expect `git add` to fail if `.git` remains read-only; record that verification passed and commit was skipped.

---

## Self-Review Against Spec

- Service ownership: Task 1 creates `service/puljefordeling/first_choice.go`.
- Derived status: Task 1 implements current and other-pulje first-choice from assignment plus interest data, ignoring GM rows.
- Split mutations: Task 2 implements first-choice interest mutation independent of assignment; Task 3 keeps assignment route separate.
- UI controls: Task 4 adds separate first-choice button, disables it for GM and other-pulje first-choice, and replaces combined search action.
- Broadcast checks: Task 3 tests `live.BucketInterests` broadcasts after assignment and first-choice mutations.
- Verification: Task 5 runs focused tests, templ generation, and full `go test ./...`.

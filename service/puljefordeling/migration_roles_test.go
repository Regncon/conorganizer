package puljefordeling

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	_ "modernc.org/sqlite"
)

const multipleEventRolesMigration = "../../migrations/20260915100000_allow_multiple_event_roles.sql"

func TestAllowMultipleEventRolesMigration_PreservesExistingAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given existing assignments including multiple Player rows in one pulje.",
		When:  "When the independent event roles migration is applied.",
		Then:  "Then every assignment and its metadata is preserved and a second role can be added.",
	})

	// Given
	const expectedRowsAfterDualRole = 4
	db := createLegacyEventRolesDB(t)
	mustExecMigrationSQL(t, db, `
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, inserted_at, source)
		VALUES
			(1, 'event-a', 'pulje-1', 'Player', '2026-09-15T08:00:00Z', 'solver'),
			(1, 'event-b', 'pulje-1', 'Player', '2026-09-15T08:01:00Z', 'manual'),
			(1, 'event-c', 'pulje-1', 'GM',     '2026-09-15T08:02:00Z', 'manual')`)
	upSQL, _ := readEventRolesMigration(t)

	// When
	mustExecMigrationSQL(t, db, upSQL)
	_, err := db.Exec(`
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, inserted_at, source)
		VALUES (1, 'event-a', 'pulje-1', 'GM', '2026-09-15T08:03:00Z', 'manual')`)

	// Then
	if err != nil {
		t.Fatalf("insert GM beside migrated Player: %v", err)
	}
	var rowCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&rowCount); err != nil {
		t.Fatalf("count migrated assignments: %v", err)
	}
	if rowCount != expectedRowsAfterDualRole {
		t.Fatalf("assignment count = %d, want %d", rowCount, expectedRowsAfterDualRole)
	}
	var insertedAt, source string
	if err := db.QueryRow(`
		SELECT inserted_at, source
		FROM relation_events_players
		WHERE billettholder_id = 1 AND event_id = 'event-a' AND pulje_id = 'pulje-1' AND role = 'Player'`,
	).Scan(&insertedAt, &source); err != nil {
		t.Fatalf("read migrated Player metadata: %v", err)
	}
	if insertedAt != "2026-09-15T08:00:00Z" || source != "solver" {
		t.Fatalf("migrated Player metadata = (%q, %q), want (%q, %q)", insertedAt, source, "2026-09-15T08:00:00Z", "solver")
	}
}

func TestAllowMultipleEventRolesMigration_DownFailsWithoutDiscardingDualRoles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a migrated assignment with both Player and GM roles.",
		When:  "When the migration is rolled back.",
		Then:  "Then rollback fails and both role rows remain intact.",
	})

	// Given
	const expectedRows = 2
	db := createLegacyEventRolesDB(t)
	mustExecMigrationSQL(t, db, `
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, source)
		VALUES (1, 'event-a', 'pulje-1', 'Player', 'manual')`)
	upSQL, downSQL := readEventRolesMigration(t)
	mustExecMigrationSQL(t, db, upSQL)
	mustExecMigrationSQL(t, db, `
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, source)
		VALUES (1, 'event-a', 'pulje-1', 'GM', 'manual')`)

	// When
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("begin down migration: %v", err)
	}
	_, downErr := tx.Exec(downSQL)
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		t.Fatalf("rollback failed down migration: %v", rollbackErr)
	}

	// Then
	if downErr == nil {
		t.Fatal("down migration succeeded with dual roles, want a primary-key violation")
	}
	if !strings.Contains(downErr.Error(), "UNIQUE constraint failed") {
		t.Fatalf("down migration error = %v, want a primary-key violation", downErr)
	}
	var rowCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&rowCount); err != nil {
		t.Fatalf("count assignments after failed down migration: %v", err)
	}
	if rowCount != expectedRows {
		t.Fatalf("assignment count after failed down migration = %d, want %d", rowCount, expectedRows)
	}
	if _, err := db.Exec(`
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, source)
		VALUES (1, 'event-a', 'pulje-1', 'GM', 'manual')`); err == nil {
		t.Fatal("duplicate GM insert succeeded after failed down migration, want primary-key violation")
	}
}

func TestAllowMultipleEventRolesMigration_DownRestoresOldKeyWithoutDualRoles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a migrated assignment without dual roles.",
		When:  "When the migration is rolled back.",
		Then:  "Then the assignment remains and the old three-column key is restored.",
	})

	// Given
	const expectedRows = 1
	db := createLegacyEventRolesDB(t)
	mustExecMigrationSQL(t, db, `
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, source)
		VALUES (1, 'event-a', 'pulje-1', 'Player', 'manual')`)
	upSQL, downSQL := readEventRolesMigration(t)
	mustExecMigrationSQL(t, db, upSQL)

	// When
	mustExecMigrationSQL(t, db, downSQL)

	// Then
	var rowCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&rowCount); err != nil {
		t.Fatalf("count assignments after down migration: %v", err)
	}
	if rowCount != expectedRows {
		t.Fatalf("assignment count after down migration = %d, want %d", rowCount, expectedRows)
	}
	if _, err := db.Exec(`
		INSERT INTO relation_events_players
			(billettholder_id, event_id, pulje_id, role, source)
		VALUES (1, 'event-a', 'pulje-1', 'GM', 'manual')`); err == nil {
		t.Fatal("GM insert beside Player succeeded after down migration, want old-key violation")
	}
}

func createLegacyEventRolesDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", t.TempDir()+"/migration.db")
	if err != nil {
		t.Fatalf("open migration database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("close migration database: %v", err)
		}
	})
	mustExecMigrationSQL(t, db, `
		CREATE TABLE billettholdere (id INTEGER PRIMARY KEY) STRICT;
		CREATE TABLE events (id TEXT PRIMARY KEY) STRICT;
		CREATE TABLE puljer (id TEXT PRIMARY KEY) STRICT;
		INSERT INTO billettholdere (id) VALUES (1);
		INSERT INTO events (id) VALUES ('event-a'), ('event-b'), ('event-c');
		INSERT INTO puljer (id) VALUES ('pulje-1');
		CREATE TABLE relation_events_players (
			event_id TEXT NOT NULL,
			pulje_id TEXT NOT NULL,
			billettholder_id INTEGER NOT NULL,
			role TEXT NOT NULL DEFAULT 'Player' CHECK (role IN ('Player', 'GM')),
			inserted_at TEXT DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now')),
			source TEXT NOT NULL DEFAULT 'manual' CHECK (source IN ('manual', 'solver')),
			PRIMARY KEY (billettholder_id, event_id, pulje_id),
			FOREIGN KEY (billettholder_id) REFERENCES billettholdere (id),
			FOREIGN KEY (event_id) REFERENCES events (id),
			FOREIGN KEY (pulje_id) REFERENCES puljer (id)
		) STRICT;`)
	return db
}

func readEventRolesMigration(t *testing.T) (string, string) {
	t.Helper()
	contents, err := os.ReadFile(multipleEventRolesMigration)
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	parts := strings.Split(string(contents), "-- +goose Down")
	if len(parts) != 2 {
		t.Fatalf("migration contains %d Down sections, want 1", len(parts)-1)
	}
	return strings.TrimPrefix(parts[0], "-- +goose Up"), parts[1]
}

func mustExecMigrationSQL(t *testing.T, db interface {
	Exec(query string, args ...any) (sql.Result, error)
}, query string) {
	t.Helper()
	if _, err := db.Exec(query); err != nil {
		t.Fatalf("execute migration SQL: %v", err)
	}
}

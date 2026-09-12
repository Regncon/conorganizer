package puljefordeling

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func TestCommitDistribution_ConcurrentManualSeatRejectsStaleDistribution(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder is being distributed to one event.",
		When:  "An admin pins them into another event before the distribution is saved.",
		Then:  "The stale commit fails, preserves the manual seat, and can be retried safely.",
	})

	// Given
	const expectedEvent = "evB"
	const pulje = models.PuljeFredagKveld
	dbPath := filepath.Join(t.TempDir(), "commit_consistency.db")
	db, err := service.InitTestDBFrom(dbPath)
	if err != nil {
		t.Fatalf("create database: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	testutil.MustExec(t, db, "PRAGMA journal_mode=WAL")
	seedPulje(t, db, pulje, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, pulje)
	seedEvent(t, db, expectedEvent, "Bravo", 4, pulje)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evA", pulje, models.InterestLevelHigh)

	connector, err := sqlite.NewConnector(dbPath)
	if err != nil {
		t.Fatalf("create connector: %v", err)
	}
	pinned := false
	commitDB := sql.OpenDB(&beforeSolverDeleteConnector{
		Connector: connector,
		beforeDelete: func() error {
			if pinned {
				return nil
			}
			pinned = true
			return AddManualSeat(db, pulje, expectedEvent, 1)
		},
	})
	t.Cleanup(func() { _ = commitDB.Close() })
	commitDB.SetMaxOpenConns(1)

	// When
	err = CommitDistribution(commitDB, pulje)

	// Then
	if !pinned {
		t.Fatal("concurrent manual seat was never inserted")
	}
	assertOnlyManualSeat(t, db, pulje, expectedEvent)
	var sqliteErr *sqlite.Error
	if !errors.As(err, &sqliteErr) || sqliteErr.Code() != sqlite3.SQLITE_BUSY_SNAPSHOT {
		t.Fatalf("expected stale snapshot error, got %v", err)
	}
	if err := CommitDistribution(commitDB, pulje); err != nil {
		t.Fatalf("retry distribution: %v", err)
	}
	assertOnlyManualSeat(t, db, pulje, expectedEvent)
}

func assertOnlyManualSeat(t *testing.T, db *sql.DB, pulje models.Pulje, eventID string) {
	t.Helper()
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND role = ?`, string(pulje), models.EventPlayerRolePlayer); count != 1 {
		t.Fatalf("expected only the manual seat after commit, got %d player seats", count)
	}
	var actualEvent, source string
	if err := db.QueryRow(`SELECT event_id, source FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = 1 AND role = ?`, string(pulje), models.EventPlayerRolePlayer).Scan(&actualEvent, &source); err != nil {
		t.Fatalf("read persisted seat: %v", err)
	}
	if actualEvent != eventID || source != SourceManual {
		t.Fatalf("expected manual seat in %s, got %s seat in %s", eventID, source, actualEvent)
	}
}

// The wrapper inserts the competing write after emulation, immediately before
// the real SQLite deletion. All queries and transactions use the real driver.
type beforeSolverDeleteConnector struct {
	driver.Connector
	beforeDelete func() error
}

func (c *beforeSolverDeleteConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &beforeSolverDeleteConn{Conn: conn, beforeDelete: c.beforeDelete}, nil
}

type beforeSolverDeleteConn struct {
	driver.Conn
	beforeDelete func() error
}

func (c *beforeSolverDeleteConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if strings.HasPrefix(strings.TrimSpace(query), "DELETE FROM relation_events_players WHERE pulje_id") {
		if err := c.beforeDelete(); err != nil {
			return nil, err
		}
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}

package admin

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func insertBoardEvent(t *testing.T, db *sql.DB, id, title, status, ageGroup string, beginner, userID int, email, host string) {
	t.Helper()
	var userCol any
	if userID > 0 {
		userCol = userID
	} else {
		userCol = nil
	}
	_, err := db.Exec(
		`INSERT INTO events (id, title, intro, description, status, age_group, beginner_friendly, user_id, host_name, email, phone_number, max_players)
		 VALUES (?, ?, '', '', ?, ?, ?, ?, ?, ?, '', 4)`,
		id, title, status, ageGroup, beginner, userCol, host, email)
	if err != nil {
		t.Fatalf("insert event %s: %v", id, err)
	}
}

func placeBoardInPulje(t *testing.T, db *sql.DB, eventID string, pulje models.Pulje) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, ?, 1)`,
		eventID, string(pulje)); err != nil {
		t.Fatalf("place %s in %s: %v", eventID, pulje, err)
	}
}

func TestBuildScheduleBoard_PoolStatsAndCollisions(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "build_schedule_board")

	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag Kveld", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeLordagKveld), "Lørdag Kveld", "Open", "2026-01-02 18:00", "2026-01-02 22:00")

	// Owner Ola (user 7) runs two adult games — both in Fredag => collision.
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (7, 'ext7', 'ola@x.no')`)
	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "AdultsOnly", 0, 7, "ola@x.no", "Ola")
	insertBoardEvent(t, db, "g2", "Demoner", "Annonsert", "AdultsOnly", 0, 7, "ola@x.no", "Ola")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g2", models.PuljeFredagKveld)

	// Beginner game by Kari, in Lørdag only.
	insertBoardEvent(t, db, "g3", "Nybegynner", "Godkjent", "Default", 1, 0, "kari@x.no", "Kari")
	placeBoardInPulje(t, db, "g3", models.PuljeLordagKveld)

	// Approved game in no pulje => pool.
	insertBoardEvent(t, db, "g4", "Ledig", "Godkjent", "Default", 0, 0, "per@x.no", "Per")

	// Draft must never appear.
	insertBoardEvent(t, db, "g5", "Kladd", "Kladd", "Default", 0, 0, "x@x.no", "X")

	board, err := buildScheduleBoard(db, logger)
	if err != nil {
		t.Fatalf("buildScheduleBoard: %v", err)
	}

	if len(board.Pool) != 1 || board.Pool[0].EventID != "g4" {
		t.Fatalf("pool = %+v, want only g4", board.Pool)
	}

	if len(board.Columns) != 2 ||
		board.Columns[0].Pulje.ID != models.PuljeFredagKveld ||
		board.Columns[1].Pulje.ID != models.PuljeLordagKveld {
		t.Fatalf("columns order wrong: %+v", board.Columns)
	}

	fredag := board.Columns[0]
	if fredag.Stats.Games != 2 || fredag.Stats.Adults != 2 || fredag.Stats.Beginner != 0 {
		t.Fatalf("fredag stats = %+v, want {2 2 0}", fredag.Stats)
	}
	if len(fredag.Collisions) != 1 || fredag.Collisions[0].HostName != "Ola" || fredag.Collisions[0].Count != 2 {
		t.Fatalf("fredag collisions = %+v, want one Ola x2", fredag.Collisions)
	}

	lordag := board.Columns[1]
	if lordag.Stats.Games != 1 || lordag.Stats.Adults != 0 || lordag.Stats.Beginner != 1 {
		t.Fatalf("lordag stats = %+v, want {1 0 1}", lordag.Stats)
	}
	if len(lordag.Collisions) != 0 {
		t.Fatalf("lordag collisions = %+v, want none", lordag.Collisions)
	}

	if board.CollisionCount != 1 {
		t.Fatalf("CollisionCount = %d, want 1", board.CollisionCount)
	}
}

func TestBuildScheduleBoard_CarriesEventTypeRuntimeEnglish(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "build_board_card_attrs")

	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, status, age_group, beginner_friendly,
		 event_type, event_runtime, can_be_run_in_english, user_id, host_name, email, phone_number, max_players)
		 VALUES ('g1', 'Drager', '', '', 'Godkjent', 'AdultsOnly', 0, 'Roleplay', 'LongRunning', 1, NULL, 'Ola', 'ola@x.no', '', 4)`)

	board, err := buildScheduleBoard(db, logger)
	if err != nil {
		t.Fatalf("buildScheduleBoard: %v", err)
	}
	if len(board.Pool) != 1 {
		t.Fatalf("pool = %+v, want one game", board.Pool)
	}
	g := board.Pool[0]
	if g.EventType != models.EventTypeRoleplay {
		t.Errorf("EventType = %q, want Roleplay", g.EventType)
	}
	if g.Runtime != models.RunTimeLongRunning {
		t.Errorf("Runtime = %q, want LongRunning", g.Runtime)
	}
	if !g.English {
		t.Errorf("English = false, want true")
	}
}

// TestBuildScheduleBoard_CollisionMatchesOwnerNotHostName guards against matching
// collisions on the free-text host_name field: two different owners who happen to
// type the same host_name must not be conflated, and each pulje's collisions must
// be keyed by owner identity, not by that shared label.
func TestBuildScheduleBoard_CollisionMatchesOwnerNotHostName(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "build_board_owner_not_hostname")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeLordagKveld), "Lørdag", "Open", "2026-01-02 18:00", "2026-01-02 22:00")

	// Pulje 1: owner A (user 11) is double-booked. Owner B has only one game here,
	// but both A and B typed the exact same host_name "Ola Nordmann" — that must
	// not make B count as part of the collision.
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (11, 'ext11', 'a@x.no')`)
	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "Default", 0, 11, "a@x.no", "Ola Nordmann")
	insertBoardEvent(t, db, "g2", "Demoner", "Godkjent", "Default", 0, 11, "a@x.no", "Ola Nordmann")
	insertBoardEvent(t, db, "g3", "Trollmenn", "Godkjent", "Default", 0, 0, "b@x.no", "Ola Nordmann")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g2", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g3", models.PuljeFredagKveld)

	// Pulje 2: owner C (email-only identity) is double-booked separately.
	insertBoardEvent(t, db, "g4", "Alver", "Godkjent", "Default", 0, 0, "c@x.no", "Cecilie")
	insertBoardEvent(t, db, "g5", "Dverger", "Godkjent", "Default", 0, 0, "c@x.no", "Cecilie")
	placeBoardInPulje(t, db, "g4", models.PuljeLordagKveld)
	placeBoardInPulje(t, db, "g5", models.PuljeLordagKveld)

	board, err := buildScheduleBoard(db, logger)
	if err != nil {
		t.Fatalf("buildScheduleBoard: %v", err)
	}

	fredag := board.Columns[0]
	if len(fredag.Collisions) != 1 {
		t.Fatalf("fredag collisions = %+v, want exactly one", fredag.Collisions)
	}
	if fredag.Collisions[0].OwnerKey != "user:11" || fredag.Collisions[0].Count != 2 {
		t.Fatalf("fredag collision = %+v, want owner user:11 count 2", fredag.Collisions[0])
	}

	lordag := board.Columns[1]
	if len(lordag.Collisions) != 1 || lordag.Collisions[0].OwnerKey != "email:c@x.no" || lordag.Collisions[0].Count != 2 {
		t.Fatalf("lordag collisions = %+v, want owner email:c@x.no count 2", lordag.Collisions)
	}

	if board.CollisionCount != 2 {
		t.Fatalf("CollisionCount = %d, want 2", board.CollisionCount)
	}
}

// TestBuildScheduleBoard_EmailIdentityNormalisesCaseAndWhitespace pins ownerKey's
// normalisation: the same email-identified owner spelled with different case and
// surrounding whitespace across two games must collapse to one owner and be
// detected as a collision.
func TestBuildScheduleBoard_EmailIdentityNormalisesCaseAndWhitespace(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "build_board_email_normalise")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag", "Open", "2026-01-01 18:00", "2026-01-01 22:00")

	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "Default", 0, 0, "  A@X.no ", "Anne")
	insertBoardEvent(t, db, "g2", "Demoner", "Godkjent", "Default", 0, 0, "a@x.no", "Anne")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g2", models.PuljeFredagKveld)

	board, err := buildScheduleBoard(db, logger)
	if err != nil {
		t.Fatalf("buildScheduleBoard: %v", err)
	}

	fredag := board.Columns[0]
	if len(fredag.Collisions) != 1 || fredag.Collisions[0].OwnerKey != "email:a@x.no" || fredag.Collisions[0].Count != 2 {
		t.Fatalf("fredag collisions = %+v, want one owner email:a@x.no count 2", fredag.Collisions)
	}
}

func TestBuildScheduleBoard_SameOwnerDifferentPuljerNoCollision(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "build_board_no_collision")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeLordagKveld), "Lørdag", "Open", "2026-01-02 18:00", "2026-01-02 22:00")

	// Same owner (null user_id, same email) but in different puljer => no collision.
	insertBoardEvent(t, db, "a", "A", "Godkjent", "Default", 0, 0, "same@x.no", "Sam")
	insertBoardEvent(t, db, "b", "B", "Godkjent", "Default", 0, 0, "same@x.no", "Sam")
	placeBoardInPulje(t, db, "a", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "b", models.PuljeLordagKveld)

	board, err := buildScheduleBoard(db, logger)
	if err != nil {
		t.Fatalf("buildScheduleBoard: %v", err)
	}
	if board.CollisionCount != 0 {
		t.Fatalf("CollisionCount = %d, want 0", board.CollisionCount)
	}
}

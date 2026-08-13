package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestScheduleBoardContent_RendersColumnsStatsAndPool(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "schedule_board_render")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag Kveld", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "AdultsOnly", 0, 0, "ola@x.no", "Ola")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	insertBoardEvent(t, db, "g4", "Ledig spel", "Godkjent", "Default", 0, 0, "per@x.no", "Per")

	doc := templtest.Render(t, ScheduleBoardContent(db, logger, nil))
	text := strings.Join(templtest.CollectTexts(doc, "#puljeoppsett-board"), " ")

	for _, want := range []string{"Fredag Kveld", "Drager", "Ledig spel", "ingen pulje", "1 spel", "18+"} {
		if !strings.Contains(text, want) {
			t.Fatalf("board text missing %q\ngot: %s", want, text)
		}
	}
}

func TestScheduleBoardContent_CardHasBannerAndBadges(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "schedule_board_card_visual")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag Kveld", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, status, age_group, beginner_friendly,
		 event_type, event_runtime, can_be_run_in_english, user_id, host_name, email, phone_number, max_players)
		 VALUES ('g1', 'Drager', '', '', 'Godkjent', 'AdultsOnly', 1, 'Roleplay', 'LongRunning', 1, NULL, 'Ola', 'ola@x.no', '', 4)`)
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)

	doc := templtest.Render(t, ScheduleBoardContent(db, logger, nil))

	if !templtest.HasSelector(doc, ".board-card img") {
		t.Errorf("card is missing its banner <img>")
	}
	if !templtest.HasSelector(doc, ".board-card .event-card-tagicon-container") {
		t.Errorf("card is missing its attribute badge icons")
	}
}

// TestScheduleBoardContent_IsCollidingMatchesOwnerNotHostName guards against the
// card-highlighting bug where cards were matched on the free-text host_name field
// instead of owner identity: a genuinely double-booked owner's two cards must both
// carry .is-colliding, while a different owner's single card that merely shares the
// same host_name string must not.
func TestScheduleBoardContent_IsCollidingMatchesOwnerNotHostName(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "schedule_board_collision_owner_key")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag Kveld", "Open", "2026-01-01 18:00", "2026-01-01 22:00")

	// Owner A (user 21) double-booked, host_name "Ola Nordmann".
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (21, 'ext21', 'a@x.no')`)
	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "Default", 0, 21, "a@x.no", "Ola Nordmann")
	insertBoardEvent(t, db, "g2", "Demoner", "Godkjent", "Default", 0, 21, "a@x.no", "Ola Nordmann")
	// Owner B, single game, same free-text host_name — must not be flagged.
	insertBoardEvent(t, db, "g3", "Trollmenn", "Godkjent", "Default", 0, 0, "b@x.no", "Ola Nordmann")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g2", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g3", models.PuljeFredagKveld)

	doc := templtest.Render(t, ScheduleBoardContent(db, logger, nil))

	if got := doc.Find(".board-card.is-colliding").Length(); got != 2 {
		t.Fatalf(".board-card.is-colliding count = %d, want 2 (owner A's two cards only)", got)
	}
}

func TestScheduleBoardContent_ShowsDMCollisionWarning(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "schedule_board_collision")
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?,?,?,?,?)`,
		string(models.PuljeFredagKveld), "Fredag Kveld", "Open", "2026-01-01 18:00", "2026-01-01 22:00")
	insertBoardEvent(t, db, "g1", "Drager", "Godkjent", "Default", 0, 0, "ola@x.no", "Ola")
	insertBoardEvent(t, db, "g2", "Demoner", "Godkjent", "Default", 0, 0, "ola@x.no", "Ola")
	placeBoardInPulje(t, db, "g1", models.PuljeFredagKveld)
	placeBoardInPulje(t, db, "g2", models.PuljeFredagKveld)

	doc := templtest.Render(t, ScheduleBoardContent(db, logger, nil))
	text := strings.Join(templtest.CollectTexts(doc, "#puljeoppsett-board"), " ")

	for _, want := range []string{"Spelleiar-kollisjon", "Ola"} {
		if !strings.Contains(text, want) {
			t.Fatalf("collision warning missing %q\ngot: %s", want, text)
		}
	}
}

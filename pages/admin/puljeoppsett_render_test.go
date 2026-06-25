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

	doc := templtest.Render(t, ScheduleBoardContent(db, logger))
	text := strings.Join(templtest.CollectTexts(doc, "#puljeoppsett-board"), " ")

	for _, want := range []string{"Fredag Kveld", "Drager", "Ledig spel", "ingen pulje", "1 spel", "18+"} {
		if !strings.Contains(text, want) {
			t.Fatalf("board text missing %q\ngot: %s", want, text)
		}
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

	doc := templtest.Render(t, ScheduleBoardContent(db, logger))
	text := strings.Join(templtest.CollectTexts(doc, "#puljeoppsett-board"), " ")

	for _, want := range []string{"Spelleiar-kollisjon", "Ola"} {
		if !strings.Contains(text, want) {
			t.Fatalf("collision warning missing %q\ngot: %s", want, text)
		}
	}
}

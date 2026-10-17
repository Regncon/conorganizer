package formsubmission

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeRow_DeclaresAssignedRoomSignalUnderPuljeID(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at et arrangement har fått rom i en pulje.",
		When:  "Når puljeraden i arrangementsskjemaet rendres.",
		Then:  "Så skal romvalget starte med rommets id under puljens id.",
	})

	// Given
	expectedSignals := `{"roomAssignment":{"FredagKveld":"42"}}`
	db := testutil.CreateTestDB(t, "pulje_row_room_signal")
	testutil.MustExec(t, db,
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'Fredag kveld', 'Open', '2026-01-01 18:00', '2026-01-01 22:00')`,
		string(models.PuljeFredagKveld))
	testutil.MustExec(t, db,
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES ('e1', 'Spill', '', '', 'Kari Nordmann', 'kari@example.com', '', 4)`)
	testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games) VALUES (42, 'Andre', '705', 7, 1)`)
	testutil.MustExec(t, db,
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, room_id) VALUES ('e1', ?, 1, 42)`,
		string(models.PuljeFredagKveld))
	pulje := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag kveld"}

	// When
	doc := templtest.Render(t, puljeRow("e1", pulje, db))
	actualSignals, _ := doc.Find(`label:has(select[data-bind="roomAssignment.FredagKveld"])`).Attr("data-signals")

	// Then
	if actualSignals != expectedSignals {
		t.Fatalf("room assignment signals mismatch\nexpected: %s\nactual:   %s", expectedSignals, actualSignals)
	}
}

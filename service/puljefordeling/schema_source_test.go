package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestRelationEventsPlayers_AcceptsRegistrationSource(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder and an event occurrence.",
		When:  "When confirmed attendance is stored with the registration source.",
		Then:  "Then the attendee relation preserves the registration source.",
	})

	// Given
	expectedSource := models.EventPlayerSourceRegistration
	db := testutil.CreateTestDB(t, "registration_source")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "open-event", "Blood in the Clock Tower", 8, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// When
	_, err := db.Exec(`
		INSERT INTO relation_events_players(
			event_id, pulje_id, billettholder_id, role, source
		) VALUES (?, ?, ?, ?, ?)
	`, "open-event", models.PuljeFredagKveld, 1, models.EventPlayerRolePlayer, expectedSource)
	var actualSource models.EventPlayerSource
	if err == nil {
		err = db.QueryRow(`
			SELECT source
			FROM relation_events_players
			WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
		`, "open-event", models.PuljeFredagKveld, 1).Scan(&actualSource)
	}

	// Then
	if err != nil {
		t.Fatalf("expected registration source to be stored: %v", err)
	}
	if actualSource != expectedSource {
		t.Fatalf("source mismatch\nexpected: %q\nactual:   %q", expectedSource, actualSource)
	}
}

func TestRelationEventsPlayers_RejectsUnknownSource(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder and an event occurrence.",
		When:  "When attendance is stored with an unknown source.",
		Then:  "Then the database rejects the attendee relation.",
	})

	// Given
	db := testutil.CreateTestDB(t, "unknown_attendee_source")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "event", "Call of Cthulhu", 5, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// When
	_, err := db.Exec(`
		INSERT INTO relation_events_players(
			event_id, pulje_id, billettholder_id, role, source
		) VALUES (?, ?, ?, ?, ?)
	`, "event", models.PuljeFredagKveld, 1, models.EventPlayerRolePlayer, "unknown")

	// Then
	if err == nil {
		t.Fatal("expected unknown attendee source to be rejected")
	}
}

func TestRelationEventsPlayersHasSourceColumn(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_schema_source")

	rows, err := db.Query(`PRAGMA table_info(relation_events_players)`)
	if err != nil {
		t.Fatalf("pragma table_info: %v", err)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var (
			cid        int
			name       string
			ctype      string
			notnull    int
			dflt       any
			primaryKey int
		)
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &primaryKey); err != nil {
			t.Fatalf("scan column: %v", err)
		}
		if name == "source" {
			found = true
			if ctype != "TEXT" {
				t.Errorf("source column has type %q, want \"TEXT\"", ctype)
			}
			if notnull != 1 {
				t.Errorf("source column notnull = %d, want 1 (NOT NULL)", notnull)
			}
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("iterate columns: %v", err)
	}
	if !found {
		t.Error("relation_events_players is missing the 'source' column")
	}
}

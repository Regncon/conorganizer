package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

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
		}
	}
	if !found {
		t.Error("relation_events_players is missing the 'source' column")
	}
}

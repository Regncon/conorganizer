package formsubmission

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeAssignmentSearch_RendersPickerAndButtons(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_pulje_assignment_search")

	if _, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (1, 'Anna', 'A', 0, '', 0, 1)`,
	); err != nil {
		t.Fatalf("seed billettholder: %v", err)
	}

	doc := templtest.Render(t, PuljeAssignmentSearch(db, logger, models.PuljeFredagKveld))

	if doc.Find("admin-billettholder-search").Length() == 0 {
		t.Errorf("expected the search web component in the picker")
	}
	if doc.Find("button:contains('Legg til som førsteval')").Length() == 0 {
		t.Errorf("expected the 'add as first choice' button")
	}

	// The puljefordeling picker only adds players. The DM is the game's owner,
	// not a participant you pick here, so it must not offer an "add as GM" action.
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}
	if strings.Contains(html, "add_gm") {
		t.Errorf("puljefordeling picker must not offer an 'add as GM' action")
	}
}

package event

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventDetails_LeavesOutPropertyBoxWithoutProperties(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event with the default age group and runtime that is neither beginner friendly nor run in English.",
		When:  "When the event details are rendered.",
		Then:  "Then no empty property box is rendered.",
	})

	// Given
	expectedPropertyBoxes := 0
	event := models.Event{
		ID:       "details-event",
		Title:    "Testarrangement",
		AgeGroup: models.AgeGroupDefault,
		Runtime:  models.RunTimeNormal,
	}

	// When
	doc := templtest.Render(t, Event_mobile(&event, components.PreviousNext{IsRemoved: true}, nil, nil))

	// Then
	if got := doc.Find(".event-details-bottom").Length(); got != expectedPropertyBoxes {
		t.Fatalf("expected %d property boxes, got %d", expectedPropertyBoxes, got)
	}
}

func TestEventDetails_ShowsPropertyBoxWithAProperty(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a beginner friendly event with the default age group and runtime.",
		When:  "When the event details are rendered.",
		Then:  "Then the property box shows the beginner friendly property.",
	})

	// Given
	expectedProperty := "Nybegynnervennlig"
	event := models.Event{
		ID:               "details-event",
		Title:            "Testarrangement",
		AgeGroup:         models.AgeGroupDefault,
		Runtime:          models.RunTimeNormal,
		BeginnerFriendly: true,
	}

	// When
	doc := templtest.Render(t, Event_mobile(&event, components.PreviousNext{IsRemoved: true}, nil, nil))

	// Then
	propertyBox := doc.Find(".event-details-bottom")
	if propertyBox.Length() != 1 {
		t.Fatalf("expected 1 property box, got %d", propertyBox.Length())
	}
	if !strings.Contains(propertyBox.Text(), expectedProperty) {
		t.Fatalf("expected the property box to contain %q, got %q", expectedProperty, strings.TrimSpace(propertyBox.Text()))
	}
}

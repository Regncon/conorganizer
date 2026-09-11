package profilecomponent

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestNotificationSettings_HiddenControlsHaveScopedDisplayRule(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given notification configuration has not enabled either action.",
		When:  "When notification settings are rendered alongside global button styles.",
		Then:  "Then both actions remain hidden through a scoped display rule.",
	})

	// Given
	expectedHiddenControls := 2

	// When
	doc := templtest.Render(t, NotificationSettings())
	settings := doc.Find("[data-varsler-settings]")
	styles := doc.Find("style").Text()

	// Then
	if got := settings.Find("button[hidden]").Length(); got != expectedHiddenControls {
		t.Fatalf("hidden notification controls = %d, want %d", got, expectedHiddenControls)
	}
	if !strings.Contains(styles, ".notification-settings") || !strings.Contains(styles, "[hidden]") || !strings.Contains(styles, "display: none") {
		t.Fatalf("expected notification settings to contain a scoped hidden display rule, got %q", styles)
	}
	if settings.ParentsFiltered(".notification-settings-container").Length() != 1 {
		t.Fatal("expected notification settings to have an ancestor container for responsive descendant styling")
	}
}

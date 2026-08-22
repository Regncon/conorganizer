package profilecomponent

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMyProgram_WhenPuljeIsNotCompleted_RendersInterestsAndHidesPlayerResult(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a player assignment and a wish in an open pulje.",
		When:  "When Mitt festivalprogram is rendered.",
		Then:  "Then the visible HTML shows the wish and hides the player allocation.",
	})

	// Given
	expectedVisibleText := "Visible Wish Event"
	hiddenVisibleText := "Hidden Player Result"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "hidden-player-result", hiddenVisibleText)
	insertProfileProgramPublishedEvent(t, db, "visible-wish-event", expectedVisibleText)
	insertProfileProgramPlayer(t, db, "hidden-player-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "visible-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if !strings.Contains(actualText, models.InterestLevelHigh.Label()) {
		t.Fatalf("expected rendered profile program to contain interest level %q\nactual text: %s", models.InterestLevelHigh.Label(), actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func TestMyProgram_WhenPuljeIsCompleted_RendersPlayerResult(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a player assignment in a completed pulje.",
		When:  "When Mitt festivalprogram is rendered.",
		Then:  "Then the visible HTML shows what the user is playing.",
	})

	// Given
	expectedVisibleText := "Completed Player Result"
	hiddenVisibleText := "Completed Wish Hidden By Result"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-player-result", expectedVisibleText)
	insertProfileProgramPublishedEvent(t, db, "completed-wish-event", hiddenVisibleText)
	insertProfileProgramPlayer(t, db, "completed-player-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "completed-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func TestMyProgram_WhenGMEventIsInNotCompletedPulje_RendersGMEventOverInterests(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a GM assignment and a wish in a locked pulje.",
		When:  "When Mitt festivalprogram is rendered.",
		Then:  "Then the visible HTML shows the GM event instead of the interests.",
	})

	// Given
	expectedVisibleText := "Locked GM Event"
	hiddenVisibleText := "Locked Wish Hidden By GM"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-gm-event", expectedVisibleText)
	insertProfileProgramPublishedEvent(t, db, "locked-wish-event", hiddenVisibleText)
	insertProfileProgramPlayer(t, db, "locked-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)
	insertProfileProgramInterest(t, db, "locked-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func TestMyProgram_WhenProgramIsNotReady_HidesPlayerResult(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given that the festival program is not published",
		When:  "When Mitt festivalprogram is rendered.",
		Then:  "Then the visible HTML hides what the user is playing and presents information text.",
	})

	// Given
	expectedVisibleText := "Completed Player Result"
	hiddenVisibleText := "Completed Wish Hidden By Result"
	expectedStatusText := "Programmet for Regncon er ikkje publisert enno"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, false)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-player-result", expectedVisibleText)
	insertProfileProgramPublishedEvent(t, db, "completed-wish-event", hiddenVisibleText)
	insertProfileProgramPlayer(t, db, "completed-player-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "completed-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedStatusText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
}

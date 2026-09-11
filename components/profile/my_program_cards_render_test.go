package profilecomponent

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMyProgram_WhenPlayerHasCompletedAssignment_RendersPersonalGroupCard(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a completed published pulje with an assigned player, room, two GMs, and group members.",
		When:  "When the player's festival program is rendered.",
		Then:  "Then the assigned event card shows the result, room, every GM, and the player's group.",
	})

	// Given
	expectedTexts := []string{
		"Du har fått plass",
		"Completed Group Event",
		"Rom 101 · Spillsalen",
		"Alice Arrangør",
		"Bob Arrangør",
		"Se gruppen din",
		"Profile",
		"Teammate",
	}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-group-event", "Completed Group Event")
	insertProfileProgramRoom(t, db, "completed-group-event", models.PuljeFredagKveld, "101", "Spillsalen")
	insertProfileProgramPlayer(t, db, "completed-group-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramNamedPlayer(t, db, "completed-group-event", models.PuljeFredagKveld, 1002, "Teammate", "Player", models.EventPlayerRolePlayer)
	insertProfileProgramNamedPlayer(t, db, "completed-group-event", models.PuljeFredagKveld, 1003, "Alice", "Arrangør", models.EventPlayerRoleGM)
	insertProfileProgramNamedPlayer(t, db, "completed-group-event", models.PuljeFredagKveld, 1004, "Bob", "Arrangør", models.EventPlayerRoleGM)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)
	card := doc.Find("#program-FredagKveld")

	// Then
	if doc.Find("#mitt-program").Length() != 1 {
		t.Fatal("expected Mitt festivalprogram to expose the stable #mitt-program root")
	}
	if card.AttrOr("data-program-pulje", "") != string(models.PuljeFredagKveld) {
		t.Fatalf("expected completed pulje card to expose data-program-pulje=%q", models.PuljeFredagKveld)
	}
	group := card.Find("[data-program-group]")
	if group.Length() != 1 {
		t.Fatal("expected completed assignment to render an identifiable group details element")
	}
	if group.AttrOr("data-program-group", "") != string(models.PuljeFredagKveld) {
		t.Fatalf("expected group details to identify pulje %q", models.PuljeFredagKveld)
	}
	if card.Find(`[data-program-group="FredagKveld"]`).Length() != 1 {
		t.Fatal("expected completed assignment group to identify its pulje for deep links")
	}
	if card.Find(`a[href="/event/completed-group-event"]`).Length() == 0 {
		t.Fatal("expected the card to link to the existing event details URL")
	}
	for _, expectedText := range expectedTexts {
		if !strings.Contains(actualText, expectedText) {
			t.Errorf("expected rendered profile program to contain %q\nactual text: %s", expectedText, actualText)
		}
	}
}

func TestMyProgram_WhenGMEventIsNotCompleted_HidesOtherPlayers(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a GM's published event in a pulje whose allocation is not completed.",
		When:  "When the GM's festival program is rendered.",
		Then:  "Then the GM card and all GM names are visible without another player's roster entry.",
	})

	// Given
	expectedTexts := []string{"Du er GM", "Host Event", "Rom 202 · Loftet", "Greta GM", "Hans GM"}
	hiddenText := "Hidden Player"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "host-event", "Host Event")
	insertProfileProgramRoom(t, db, "host-event", models.PuljeFredagKveld, "202", "Loftet")
	insertProfileProgramPlayer(t, db, "host-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)
	insertProfileProgramNamedPlayer(t, db, "host-event", models.PuljeFredagKveld, 1002, "Hidden", "Player", models.EventPlayerRolePlayer)
	insertProfileProgramNamedPlayer(t, db, "host-event", models.PuljeFredagKveld, 1003, "Greta", "GM", models.EventPlayerRoleGM)
	insertProfileProgramNamedPlayer(t, db, "host-event", models.PuljeFredagKveld, 1004, "Hans", "GM", models.EventPlayerRoleGM)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	for _, expectedText := range expectedTexts {
		if !strings.Contains(actualText, expectedText) {
			t.Errorf("expected rendered profile program to contain %q\nactual text: %s", expectedText, actualText)
		}
	}
	if strings.Contains(actualText, hiddenText) {
		t.Fatalf("expected player roster to stay private before completion\nactual text: %s", actualText)
	}
	if doc.Find("[data-program-group]").Length() != 0 {
		t.Fatal("expected no group details before the pulje is completed")
	}
}

func TestMyProgram_WhenCompletedInterestWasNotAssigned_RendersUnassignedMessage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a completed published pulje where the player has a wish and no assignment.",
		When:  "When the player's festival program is rendered.",
		Then:  "Then the card clearly says the player did not get a place and can contact the organizers.",
	})

	// Given
	expectedTexts := []string{
		"Du fikk dessverre ikke plass i denne puljen",
		"Ta kontakt med arrangørene",
		"Unassigned Wish",
	}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "unassigned-wish", "Unassigned Wish")
	insertProfileProgramInterest(t, db, "unassigned-wish", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	for _, expectedText := range expectedTexts {
		if !strings.Contains(actualText, expectedText) {
			t.Errorf("expected rendered profile program to contain %q\nactual text: %s", expectedText, actualText)
		}
	}
}

func TestMyProgram_WhenProgramIsUnpublished_RendersWishesWithoutAllocation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an unpublished program with a completed allocation and an existing wish.",
		When:  "When the player's festival program is rendered.",
		Then:  "Then the allocation stays hidden while the wish and unpublished status remain visible.",
	})

	// Given
	expectedTexts := []string{"Fordelingen er ikke publisert ennå", "Visible Unpublished Wish"}
	hiddenTexts := []string{"Hidden Unpublished Result", "Private Teammate", "Du har fått plass"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, false)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "hidden-unpublished-result", "Hidden Unpublished Result")
	insertProfileProgramPublishedEvent(t, db, "visible-unpublished-wish", "Visible Unpublished Wish")
	insertProfileProgramPlayer(t, db, "hidden-unpublished-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramNamedPlayer(t, db, "hidden-unpublished-result", models.PuljeFredagKveld, 1002, "Private", "Teammate", models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "visible-unpublished-wish", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	for _, expectedText := range expectedTexts {
		if !strings.Contains(actualText, expectedText) {
			t.Errorf("expected rendered profile program to contain %q\nactual text: %s", expectedText, actualText)
		}
	}
	for _, hiddenText := range hiddenTexts {
		if strings.Contains(actualText, hiddenText) {
			t.Errorf("expected rendered profile program to hide %q\nactual text: %s", hiddenText, actualText)
		}
	}
}

func TestMyProgram_WhenCompletedPuljeHasNoWishOrAssignment_RendersNeutralState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a completed published pulje where the player had no wish and no assignment.",
		When:  "When the player's festival program is rendered.",
		Then:  "Then the pulje remains neutral and does not claim the player failed to get a place.",
	})

	// Given
	hiddenText := "Du fikk dessverre ikke plass i denne puljen"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if strings.Contains(actualText, hiddenText) {
		t.Fatalf("expected a neutral empty state without %q\nactual text: %s", hiddenText, actualText)
	}
}

func TestMyProgram_WhenAnotherOwnedBillettholderHasAGroup_HidesThatGroup(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given two billettholdere owned by the user with different completed event groups.",
		When:  "When the program for one selected billettholder is rendered.",
		Then:  "Then only that billettholder's actual event and group are visible.",
	})

	// Given
	expectedText := "Selected Group Event"
	hiddenTexts := []string{"Other Group Event", "Secret Member"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, selectedBillettholderID := seedProfileProgramUser(t, db)
	otherBillettholderID := insertProfileProgramOwnedBillettholder(t, db, 1002, "Other", "Ticket")
	insertProfileProgram(t, db, true)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "selected-group-event", expectedText)
	insertProfileProgramPublishedEvent(t, db, "other-group-event", "Other Group Event")
	insertProfileProgramPlayer(t, db, "selected-group-event", models.PuljeFredagKveld, selectedBillettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramPlayer(t, db, "other-group-event", models.PuljeFredagKveld, otherBillettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramNamedPlayer(t, db, "other-group-event", models.PuljeFredagKveld, 1003, "Secret", "Member", models.EventPlayerRolePlayer)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, selectedBillettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedText, actualText)
	}
	for _, hiddenText := range hiddenTexts {
		if strings.Contains(actualText, hiddenText) {
			t.Errorf("expected selected billettholder program to hide %q\nactual text: %s", hiddenText, actualText)
		}
	}
}

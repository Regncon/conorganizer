package profilecomponent

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
)

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInOpenPulje_ReturnsInterestsInstead(t *testing.T) {
	// Given a player assignment in an open pulje,
	// when the profile program data is loaded,
	// then the player result is hidden and the user's interests are returned.

	// Given
	expectedEventTitles := []string{}
	expectedInterestNames := []string{"Open Wish Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "open-assigned-event", "Open Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "open-wish-event", "Open Wish Event")
	insertProfileProgramPlayer(t, db, "open-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "open-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInLockedPulje_ReturnsInterestsInstead(t *testing.T) {
	// Given a player assignment in a locked pulje,
	// when the profile program data is loaded,
	// then the player result is hidden and the user's interests are returned.

	// Given
	expectedEventTitles := []string{}
	expectedInterestNames := []string{"Locked Wish Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-assigned-event", "Locked Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "locked-wish-event", "Locked Wish Event")
	insertProfileProgramPlayer(t, db, "locked-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "locked-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInCompletedPulje_ReturnsPlayerResult(t *testing.T) {
	// Given a player assignment in a completed pulje,
	// when the profile program data is loaded,
	// then the assigned event is returned as the user's program.

	// Given
	expectedEventTitles := []string{"Completed Assigned Event"}
	expectedInterestNames := []string{}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-assigned-event", "Completed Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "completed-wish-event", "Completed Wish Event")
	insertProfileProgramPlayer(t, db, "completed-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "completed-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenGMEventIsInOpenPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in an open pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Open GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "open-gm-event", "Open GM Event")
	insertProfileProgramPlayer(t, db, "open-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

func TestGetAllEventsForUser_WhenGMEventIsInLockedPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in a locked pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Locked GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-gm-event", "Locked GM Event")
	insertProfileProgramPlayer(t, db, "locked-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

func TestGetAllEventsForUser_WhenGMEventIsInCompletedPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in a completed pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Completed GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-gm-event", "Completed GM Event")
	insertProfileProgramPlayer(t, db, "completed-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

package event

import (
	"database/sql"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestCanViewEvent_AllowsPublicAccessToAnnouncedEvent(t *testing.T) {
	// Gitt at arrangementet er annonsert,
	// når visningstilgang sjekkes,
	// så skal alle kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusAnnounced}
	userInfo := requestctx.UserRequestInfo{}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_AllowsAdminAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren er admin,
	// når visningstilgang sjekkes,
	// så skal admin kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusApproved}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, IsAdmin: true}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_AllowsCreatorAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren er eieren,
	// når visningstilgang sjekkes,
	// så skal eieren kunne se arrangementet.

	// Given
	expectedCanView := true

	db := createEventVisibilityTestDB(t)
	seedEventVisibilityUser(t, db, 101, "creator-user")
	event := &models.Event{
		Status: models.EventStatusApproved,
		UserID: sql.NullInt64{
			Int64: 101,
			Valid: true,
		},
	}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, Id: "creator-user"}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestCanViewEvent_DeniesAnonymousAccessToUnannouncedEvent(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren ikke er innlogget,
	// når visningstilgang sjekkes,
	// så skal brukeren ikke kunne se arrangementet.

	// Given
	expectedCanView := false

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusApproved}
	userInfo := requestctx.UserRequestInfo{}

	// When
	actualCanView, err := canViewEvent(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility check to succeed: %v", err)
	}
	if actualCanView != expectedCanView {
		t.Fatalf("visibility mismatch\nexpected: %v\nactual:   %v", expectedCanView, actualCanView)
	}
}

func TestEventPageContent_WhenUnannouncedEventIsHidden_RendersFriendlyMessage(t *testing.T) {
	// Gitt at arrangementet ikke er annonsert og brukeren ikke har tilgang,
	// når arrangementssiden vises,
	// så skal brukeren se en vennlig melding om at arrangementet ikke er annonsert ennå.

	// Given
	expectedMessages := []string{"Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter."}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "hidden-event", "Hidden Event", models.EventStatusApproved, sql.NullInt64{})
	request := httptest.NewRequest("GET", "/event/hidden-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("hidden-event", false, logger, db, nil, request))
	actualMessages := templtest.CollectTexts(doc, ".event-not-announced-message")

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("hidden event message mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
}

func TestEventPageContent_WhenProgramPublishingIsOff_DoesNotRenderInterestDialog(t *testing.T) {
	// Gitt at programmet ikke er publisert, men arrangementet er publisert i en pulje,
	// når arrangementssiden vises,
	// så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.

	// Given
	expectedDialogVisible := false
	expectedHrefs := []string{"/", "/profile/tickets"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "program-off-event", "Program Off Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "program-off-event", models.PuljeFredagKveld, true)
	setEventVisibilityProgramPublishing(t, db, false)
	request := httptest.NewRequest("GET", "/event/program-off-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("program-off-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestEventPageContent_WhenEventIsNotPublishedInPulje_DoesNotRenderInterestDialog(t *testing.T) {
	// Gitt at programmet er publisert, men arrangementet ikke er publisert i puljen,
	// når arrangementssiden vises,
	// så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.

	// Given
	expectedDialogVisible := false
	expectedHrefs := []string{"/", "/profile/tickets"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "unpublished-pulje-event", "Unpublished Pulje Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "unpublished-pulje-event", models.PuljeFredagKveld, false)
	setEventVisibilityProgramPublishing(t, db, true)
	request := httptest.NewRequest("GET", "/event/unpublished-pulje-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("unpublished-pulje-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestEventPageContent_WhenProgramAndPuljeArePublished_RendersInterestDialog(t *testing.T) {
	// Gitt at programmet er publisert og arrangementet er publisert i en pulje,
	// når arrangementssiden vises,
	// så skal dialogen for interessevalg rendres.

	// Given
	expectedDialogVisible := true

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "published-interest-event", "Published Interest Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "published-interest-event", models.PuljeFredagKveld, true)
	setEventVisibilityProgramPublishing(t, db, true)
	request := httptest.NewRequest("GET", "/event/published-interest-event?pulje=fredag_kveld", nil)

	// When
	doc := templtest.Render(t, event_page_content("published-interest-event", false, logger, db, nil, request))
	actualDialogVisible := templtest.HasSelector(doc, ".interest-dialog")

	// Then
	if actualDialogVisible != expectedDialogVisible {
		t.Fatalf("interest dialog visibility mismatch\nexpected: %v\nactual:   %v", expectedDialogVisible, actualDialogVisible)
	}
}

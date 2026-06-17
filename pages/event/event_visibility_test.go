package event

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestCanViewEvent_AllowsPublicAccessToAnnouncedEvent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet er annonsert.",
		When:  "Når visningstilgang sjekkes.",
		Then:  "Så skal alle kunne se arrangementet.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren er admin.",
		When:  "Når visningstilgang sjekkes.",
		Then:  "Så skal admin kunne se arrangementet.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren er eieren.",
		When:  "Når visningstilgang sjekkes.",
		Then:  "Så skal eieren kunne se arrangementet.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren ikke er innlogget.",
		When:  "Når visningstilgang sjekkes.",
		Then:  "Så skal brukeren ikke kunne se arrangementet.",
	})

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

func TestCanViewEvent_DeniesLoggedInNonOwnerAccessToUnannouncedEvent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren er innlogget, men ikke eier.",
		When:  "Når visningstilgang sjekkes.",
		Then:  "Så skal brukeren ikke kunne se arrangementet.",
	})

	// Given
	expectedCanView := false

	db := createEventVisibilityTestDB(t)
	seedEventVisibilityUser(t, db, 101, "creator-user")
	event := &models.Event{
		Status: models.EventStatusApproved,
		UserID: sql.NullInt64{
			Int64: 101,
			Valid: true,
		},
	}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, Id: "someone-else"}

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

func TestDecideEventView_WhenOwnerViewsUnannouncedEvent_ShowsWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren er eieren.",
		When:  "Når visningsstatus avgjøres.",
		Then:  "Så skal eieren kunne se arrangementet med varsel om at det ikke er annonsert.",
	})

	// Given
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
	decision, err := decideEventView(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility decision to succeed: %v", err)
	}
	if !decision.CanView {
		t.Fatalf("expected owner to view unannounced event")
	}
	if !decision.ShowUnannouncedWarning {
		t.Fatalf("expected owner view to show unannounced warning")
	}
}

func TestDecideEventView_WhenArchivedEventIsHidden_ReturnsGoneDecision(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet er forkastet og brukeren ikke har utvidet tilgang.",
		When:  "Når visningsstatus avgjøres.",
		Then:  "Så skal arrangementet skjules som Gone.",
	})

	// Given
	expectedStatusCode := http.StatusGone
	expectedHiddenReason := eventHiddenReasonArchivedGone

	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusArchived}
	userInfo := requestctx.UserRequestInfo{}

	// When
	decision, err := decideEventView(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility decision to succeed: %v", err)
	}
	if decision.CanView {
		t.Fatalf("expected archived event to be hidden")
	}
	if decision.HiddenResponseStatusCode != expectedStatusCode {
		t.Fatalf("status code mismatch\nexpected: %v\nactual:   %v", expectedStatusCode, decision.HiddenResponseStatusCode)
	}
	if decision.HiddenReason != expectedHiddenReason {
		t.Fatalf("hidden reason mismatch\nexpected: %v\nactual:   %v", expectedHiddenReason, decision.HiddenReason)
	}
}

func TestDecideEventView_WhenAdminViewsArchivedEvent_ShowsArchivedWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet er forkastet og brukeren er admin.",
		When:  "Når visningsstatus avgjøres.",
		Then:  "Så skal admin kunne se arrangementet med varsel.",
	})

	// Given
	db := createEventVisibilityTestDB(t)
	event := &models.Event{Status: models.EventStatusArchived}
	userInfo := requestctx.UserRequestInfo{IsLoggedIn: true, IsAdmin: true}

	// When
	decision, err := decideEventView(event, userInfo, db)

	// Then
	if err != nil {
		t.Fatalf("expected visibility decision to succeed: %v", err)
	}
	if !decision.CanView {
		t.Fatalf("expected admin to view archived event")
	}
	if !decision.ShowArchivedWarning {
		t.Fatalf("expected admin view to show archived warning")
	}
}

func TestEventPageContent_WhenUnannouncedEventIsHidden_RendersFriendlyMessage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren ikke har tilgang.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal brukeren se en vennlig melding om at arrangementet ikke er annonsert ennå.",
	})

	// Given
	expectedMessages := []string{"Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter."}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "hidden-event", "Hidden Event", models.EventStatusApproved, sql.NullInt64{})
	request := httptest.NewRequest("GET", "/event/hidden-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("hidden-event", false, logger, db, nil, request))
	actualMessages := templtest.CollectTexts(doc, ".event-not-announced-message")
	actualEventContentVisible := templtest.HasSelector(doc, ".event-page-wrapper")

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("hidden event message mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
	if actualEventContentVisible {
		t.Fatalf("expected hidden unannounced event content not to render")
	}
}

func TestEventPageContent_WhenAdminViewsUnannouncedEvent_RendersWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er annonsert og brukeren er admin.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal admin se arrangementet med et tydelig varsel.",
	})

	// Given
	expectedWarnings := []string{"Arrangementet er ikke annonsert"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "admin-warning-event", "Admin Warning Event", models.EventStatusApproved, sqlNullInt64(501))
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "admin-warning-event", models.PuljeFredagKveld, true)
	request := httptest.NewRequest("GET", "/event/admin-warning-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("admin-warning-event", true, logger, db, nil, request))
	actualWarnings := templtest.CollectTexts(doc, ".event-visibility-warning-title")
	actualEventTitles := templtest.CollectTexts(doc, "h1")

	// Then
	if !slices.Contains(actualWarnings, expectedWarnings[0]) {
		t.Fatalf("warning title mismatch\nexpected to contain: %v\nactual:              %v", expectedWarnings, actualWarnings)
	}
	if !slices.Contains(actualEventTitles, "Admin Warning Event") {
		t.Fatalf("expected rendered event title, got h1 texts: %v", actualEventTitles)
	}
}

func TestEventPageContent_WhenArchivedEventIsHidden_RendersArchivedMessage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet er forkastet og brukeren ikke har tilgang.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal brukeren se en melding om at arrangementet ikke er tilgjengelig.",
	})

	// Given
	expectedMessages := []string{"Dette arrangementet er ikke tilgjengelig lenger."}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "archived-hidden-event", "Archived Hidden Event", models.EventStatusArchived, sql.NullInt64{})
	request := httptest.NewRequest("GET", "/event/archived-hidden-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("archived-hidden-event", false, logger, db, nil, request))
	actualMessages := templtest.CollectTexts(doc, ".event-archived-message")
	actualEventContentVisible := templtest.HasSelector(doc, ".event-page-wrapper")

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("archived event message mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
	if actualEventContentVisible {
		t.Fatalf("expected hidden archived event content not to render")
	}
}

func TestEventLayoutRoute_WhenArchivedEventIsHidden_ReturnsGone(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet er forkastet og brukeren ikke har tilgang.",
		When:  "Når arrangementsruten åpnes.",
		Then:  "Så skal HTTP-status være 410 Gone.",
	})

	// Given
	expectedStatusCode := http.StatusGone

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEventVisibilityEvent(t, db, "archived-route-event", "Archived Route Event", models.EventStatusArchived, sql.NullInt64{})
	router := chi.NewRouter()
	eventLayoutRoute(router, db, logger, nil, nil)
	request := httptest.NewRequest(http.MethodGet, "/event/archived-route-event", nil)
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("status code mismatch\nexpected: %v\nactual:   %v", expectedStatusCode, recorder.Code)
	}
}

func TestEventPageContent_WhenProgramPublishingIsOff_DoesNotRenderInterestDialog(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet ikke er publisert, men arrangementet er publisert i en pulje.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er publisert, men arrangementet ikke er publisert i puljen.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal dialogen ikke rendres og brukeren fortsatt se lenken for å hente billett.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er publisert og arrangementet er publisert i en pulje.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så skal dialogen for interessevalg rendres.",
	})

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

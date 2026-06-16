package event

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventPageContent_WhenAdminOpensEvent_RendersAdminEditLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at admin åpner en arrangementsdetaljside.",
		When:  "Når siden rendres.",
		Then:  "Så skal inngangen til administrativ redigering vises uten å erstatte den vanlige visningen.",
	})

	// Given
	expectedHref := "/admin/approval/edit/admin-visible-event"
	expectedTextPart := "Administrator - rediger arrangementer"

	db := createEventVisibilityTestDB(t)
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEvent(t, db, "admin-visible-event", "Admin Visible Event", models.EventStatusApproved, sqlNullInt64(501))
	seedEventVisibilityEventPulje(t, db, "admin-visible-event", models.PuljeFredagKveld, true)
	request := httptest.NewRequest(http.MethodGet, "/event/admin-visible-event", nil)
	logger := testutil.NewTestLogger()

	// When
	doc := templtest.Render(t, event_page_content("admin-visible-event", true, logger, db, nil, request))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, "#event-page-content, a"), " ")

	// Then
	if !slices.Contains(actualHrefs, expectedHref) {
		t.Fatalf("expected admin edit href %q in %v", expectedHref, actualHrefs)
	}
	if !strings.Contains(actualText, expectedTextPart) {
		t.Fatalf("expected event page text to contain %q\nactual text: %s", expectedTextPart, actualText)
	}
}

func TestInterestErrorMessageFromError_ReturnsFriendlyMessages(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at interesseoppdatering feiler av kjente årsaker.",
		When:  "Når feilen oversettes til brukerbeskjed.",
		Then:  "Så skal meldingen være handlingsrettet og uten tekniske detaljer.",
	})

	// Given
	expectedMessages := []string{
		"Denne pulja er ikkje tilgjengeleg for dette arrangementet.",
		"Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar.",
		"Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk.",
		"Interessevalget er ikke åpnet ennå.",
	}
	errors := []error{
		errForInterestMessage("pulje fredag is not active and published for event abc"),
		errForInterestMessage("pulje fredag is locked for event abc"),
		errForInterestMessage("pulje fredag is completed for event abc"),
		errForInterestMessage("program is not published"),
	}

	// When
	actualMessages := make([]string, 0, len(errors))
	for _, err := range errors {
		actualMessages = append(actualMessages, interestErrorMessageFromError(err))
	}

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("interest error messages mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
}

func sqlNullInt64(value int64) sql.NullInt64 {
	return sql.NullInt64{Int64: value, Valid: true}
}

type errForInterestMessage string

func (err errForInterestMessage) Error() string {
	return string(err)
}

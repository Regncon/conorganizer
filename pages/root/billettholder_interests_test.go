package root

import (
	"bytes"
	"context"
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

const interestTestUserEmail = "forelder@example.com"

func TestRootPageContent_WhenBillettholdereOnTheAccountHaveInterest_ListsYouFirstThenTheMostInterested(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at tre billettholdere på kontoen har meldt interesse på et arrangement i en pulje, og Anders er kontoens egen billettholder.",
		When:  "Når forsiden vises.",
		Then:  "Så skal arrangementskortet vise ett hjerte per billettholder, og lista skal vise Anders som «Du» først og deretter de mest interesserte.",
	})

	// Given
	expectedDescriptions := []string{"Du er litt interessert", "Amalie er veldig interessert", "Bjørn er interessert"}
	expectedHeartCount := 3
	expectedSelectedHeartCount := 1

	db := createInterestRootPageTestDB(t)
	andersID := insertRootPageBillettholder(t, db, "Anders", "Andersen", 1, interestTestUserEmail)
	amalieID := insertRootPageBillettholder(t, db, "Amalie", "Berg", 2, "amalie@example.com")
	bjornID := insertRootPageBillettholder(t, db, "Bjørn", "Berg", 3, "bjorn@example.com")
	associateRootPageBillettholder(t, db, amalieID, interestTestUserEmail)
	associateRootPageBillettholder(t, db, bjornID, interestTestUserEmail)
	insertRootPageInterest(t, db, andersID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelLow)
	insertRootPageInterest(t, db, amalieID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelHigh)
	insertRootPageInterest(t, db, bjornID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelMedium)

	// When
	doc := renderRootPageAs(t, db, interestTestUserEmail)

	// Then
	assertTexts(t, "interest descriptions", expectedDescriptions, strings.Split(doc.Find(".interest-indicator").AttrOr("data-tippy-content", ""), "\n"))
	if actual := doc.Find(".interest-indicator-heart").Length(); actual != expectedHeartCount {
		t.Fatalf("interest heart count = %d, want %d", actual, expectedHeartCount)
	}
	if actual := doc.Find(".interest-indicator-heart.selected").Length(); actual != expectedSelectedHeartCount {
		t.Fatalf("selected interest heart count = %d, want %d", actual, expectedSelectedHeartCount)
	}
}

func TestRootPageContent_WhenOnlyALinkedBillettholderHasInterest_ShowsTheirNameInsteadOfYou(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder bare er knyttet til kontoen, men har billetten på en annen e-post.",
		When:  "Når forsiden vises.",
		Then:  "Så skal interessen vises med billettholderens navn, ikke som «Du».",
	})

	// Given
	expectedDescriptions := []string{"Amalie er veldig interessert"}

	db := createInterestRootPageTestDB(t)
	insertRootPageBillettholder(t, db, "Anders", "Andersen", 1, interestTestUserEmail)
	amalieID := insertRootPageBillettholder(t, db, "Amalie", "Berg", 2, "amalie@example.com")
	associateRootPageBillettholder(t, db, amalieID, interestTestUserEmail)
	insertRootPageInterest(t, db, amalieID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelHigh)

	// When
	doc := renderRootPageAs(t, db, interestTestUserEmail)

	// Then
	assertTexts(t, "interest descriptions", expectedDescriptions, strings.Split(doc.Find(".interest-indicator").AttrOr("data-tippy-content", ""), "\n"))
}

func TestRootPageContent_WhenTheTicketEmailDiffersFromTheLoginEmail_ShowsTheNameInsteadOfYou(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at kontoens eneste billettholder har billetten på en annen e-post enn den brukeren logger inn med.",
		When:  "Når forsiden vises.",
		Then:  "Så skal interessen vises med billettholderens navn, fordi vi ikke kan bekrefte at det er brukeren selv.",
	})

	// Given
	expectedDescriptions := []string{"Anders er veldig interessert"}

	db := createInterestRootPageTestDB(t)
	andersID := insertRootPageBillettholder(t, db, "Anders", "Andersen", 1, "annen-epost@example.com")
	associateRootPageBillettholder(t, db, andersID, interestTestUserEmail)
	insertRootPageInterest(t, db, andersID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelHigh)

	// When
	doc := renderRootPageAs(t, db, interestTestUserEmail)

	// Then
	assertTexts(t, "interest descriptions", expectedDescriptions, strings.Split(doc.Find(".interest-indicator").AttrOr("data-tippy-content", ""), "\n"))
}

func TestRootPageContent_WhenOnlyAnotherAccountHasInterest_ShowsNoInterestIndicator(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at bare en billettholder på en annen konto har meldt interesse på arrangementet.",
		When:  "Når forsiden vises.",
		Then:  "Så skal arrangementskortet ikke vise noen interesse.",
	})

	// Given
	expectedInterestIndicatorVisible := false

	db := createInterestRootPageTestDB(t)
	insertRootPageBillettholder(t, db, "Anders", "Andersen", 1, interestTestUserEmail)
	strangerID := insertRootPageBillettholder(t, db, "Frida", "Fremmed", 2, "fremmed@example.com")
	insertRootPageInterest(t, db, strangerID, "alpha-event", models.PuljeFredagKveld, models.InterestLevelHigh)

	// When
	doc := renderRootPageAs(t, db, interestTestUserEmail)

	// Then
	if actual := templtest.HasSelector(doc, ".interest-indicator"); actual != expectedInterestIndicatorVisible {
		t.Fatalf("interest indicator visible = %v, want %v", actual, expectedInterestIndicatorVisible)
	}
}

func TestRootPageContent_WhenInterestIsInAnotherPulje_ShowsNoInterestIndicatorOnThisPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder på kontoen har meldt interesse på arrangementet i en annen pulje.",
		When:  "Når forsiden viser arrangementet i fredag kveld.",
		Then:  "Så skal arrangementskortet ikke vise interessen fra den andre puljen.",
	})

	// Given
	expectedInterestIndicatorVisible := false

	db := createInterestRootPageTestDB(t)
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertRootPageEventPulje(t, db, "alpha-event", models.PuljeLordagMorgen, true)
	andersID := insertRootPageBillettholder(t, db, "Anders", "Andersen", 1, interestTestUserEmail)
	insertRootPageInterest(t, db, andersID, "alpha-event", models.PuljeLordagMorgen, models.InterestLevelHigh)

	// When
	doc := renderRootPageAs(t, db, interestTestUserEmail)

	// Then
	if actual := templtest.HasSelector(doc, ".interest-indicator"); actual != expectedInterestIndicatorVisible {
		t.Fatalf("interest indicator visible = %v, want %v", actual, expectedInterestIndicatorVisible)
	}
}

func createInterestRootPageTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	insertRootPageEvent(t, db, "alpha-event", "Alpha Event", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "alpha-event", models.PuljeFredagKveld, true)
	return db
}

func renderRootPageAs(t *testing.T, db *sql.DB, email string) *goquery.Document {
	t.Helper()

	ctx := authctx.WithUserToken(context.Background(), "interest-test-user", email)
	var html bytes.Buffer
	if err := rootPageContentForDate(db, nil, "2026-10-09").Render(ctx, &html); err != nil {
		t.Fatalf("render root page: %v", err)
	}
	doc, err := goquery.NewDocumentFromReader(&html)
	if err != nil {
		t.Fatalf("parse root page html: %v", err)
	}
	return doc
}

func insertRootPageBillettholder(t *testing.T, db *sql.DB, firstName string, lastName string, ticketID int, email string) int {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO billettholdere(first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES(?, ?, 1, 'Festivalpass', 1, ?)
	`, firstName, lastName, ticketID)
	if err != nil {
		t.Fatalf("failed to insert billettholder: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to read billettholder id: %v", err)
	}
	mustExec(t, db, `
		INSERT INTO relation_billettholder_emails(billettholder_id, email, kind)
		VALUES(?, ?, ?)
	`, id, email, models.BillettholderEmailKindTicket)
	return int(id)
}

func associateRootPageBillettholder(t *testing.T, db *sql.DB, billettholderID int, email string) {
	t.Helper()

	mustExec(t, db, `
		INSERT INTO relation_billettholder_emails(billettholder_id, email, kind)
		VALUES(?, ?, ?)
	`, billettholderID, email, models.BillettholderEmailKindAssociated)
}

func insertRootPageInterest(t *testing.T, db *sql.DB, billettholderID int, eventID string, puljeID models.Pulje, level models.InterestLevel) {
	t.Helper()

	mustExec(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES(?, ?, ?, ?)
	`, billettholderID, eventID, puljeID, level)
}

func assertTexts(t *testing.T, name string, expected []string, actual []string) {
	t.Helper()

	if !slices.Equal(expected, actual) {
		t.Fatalf("%s mismatch\nexpected: %v\nactual:   %v", name, expected, actual)
	}
}

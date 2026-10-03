package root

import (
	"slices"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/components/icons"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

const (
	closingAlertText   = "Interessevalget stenger snart! Hvis du vil endre valgene dine for kommende pulje, gjør det nå."
	closingAlertLink   = "Klikk her for å gå til puljen."
	lockedAlertText    = "Interessevalg for kommende pulje er nå låst og kan ikke endres. Vi jobber med å fordele spillere og publiserer resultatet snart!"
	completedAlertText = "Puljefordelingen er klar! Se hva du fikk på profilen din. Om du fikk en plass, møt opp ved rommet når puljen starter (kl 10 / kl 18)."
)

func TestRootPageContent_WhenPuljeHasAlert_ShowsItAtTopOfPage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fredagens pulje er låst.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal puljevarselet være det første på siden.",
	})

	// Given
	expectedFirstElementIsAlert := true
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked, false)
	insertRootPageEvent(t, db, "friday-raffle", "Friday Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", rootPageTestNow))
	actualFirstElementIsAlert := doc.Find("body").Children().Not("style").First().Is(".pulje-alerts")

	// Then
	if actualFirstElementIsAlert != expectedFirstElementIsAlert {
		t.Fatalf("pulje alerts first on page mismatch\nexpected: %v\nactual:   %v", expectedFirstElementIsAlert, actualFirstElementIsAlert)
	}
}

func TestRootPageContent_WhenDayHasSeveralPuljer_ShowsOneNamedAlertPerPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at lørdagens to puljer er låst eller i ferd med å bli låst, og har flere arrangementer.",
		When:  "Når forsiden vises for lørdag.",
		Then:  "Så skal hver pulje få ett varsel med puljenavn, ikke ett per arrangement.",
	})

	// Given
	expectedAlerts := []string{
		"Lørdag morgen: " + lockedAlertText,
		"Lørdag kveld: " + closingAlertText + " " + closingAlertLink,
	}
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")
	setRootPagePuljeStatus(t, db, models.PuljeLordagMorgen, models.PuljeStatusLocked, false)
	setRootPagePuljeStatus(t, db, models.PuljeLordagKveld, models.PuljeStatusOpen, true)
	for _, event := range []struct {
		id                 string
		puljeID            models.Pulje
		isInPuljefordeling bool
	}{
		{id: "morning-program-alpha", puljeID: models.PuljeLordagMorgen},
		{id: "morning-program-beta", puljeID: models.PuljeLordagMorgen},
		{id: "morning-raffle", puljeID: models.PuljeLordagMorgen, isInPuljefordeling: true},
		{id: "evening-program", puljeID: models.PuljeLordagKveld},
		{id: "evening-raffle-alpha", puljeID: models.PuljeLordagKveld, isInPuljefordeling: true},
		{id: "evening-raffle-beta", puljeID: models.PuljeLordagKveld, isInPuljefordeling: true},
	} {
		insertRootPageEvent(t, db, event.id, event.id, models.EventStatusAnnounced, event.isInPuljefordeling)
		insertRootPageEventPulje(t, db, event.id, event.puljeID, true)
	}

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-10", rootPageTestNow))
	actualAlerts := templtest.CollectTexts(doc, ".pulje-alert")

	// Then
	if !slices.Equal(expectedAlerts, actualAlerts) {
		t.Fatalf("pulje alerts mismatch\nexpected: %v\nactual:   %v", expectedAlerts, actualAlerts)
	}
}

func TestRootPageContent_WhenDayHasOnePulje_ShowsAlertWithoutPuljeName(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fredagens eneste pulje er låst og har flere arrangementer.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal varselet vises én gang uten puljenavn.",
	})

	// Given
	expectedAlerts := []string{lockedAlertText}
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked, false)
	insertRootPageEvent(t, db, "friday-program", "Friday Program", models.EventStatusAnnounced, false)
	insertRootPageEventPulje(t, db, "friday-program", models.PuljeFredagKveld, true)
	insertRootPageEvent(t, db, "friday-raffle-alpha", "Friday Raffle Alpha", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle-alpha", models.PuljeFredagKveld, true)
	insertRootPageEvent(t, db, "friday-raffle-beta", "Friday Raffle Beta", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle-beta", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", rootPageTestNow))
	actualAlerts := templtest.CollectTexts(doc, ".pulje-alert")

	// Then
	if !slices.Equal(expectedAlerts, actualAlerts) {
		t.Fatalf("pulje alerts mismatch\nexpected: %v\nactual:   %v", expectedAlerts, actualAlerts)
	}
}

func TestRootPageContent_WhenPuljeIsClosingSoon_LinksAlertToPuljeHeading(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at interessevalget for fredagens pulje stenger snart.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal varselet lenke til overskriften for puljen på samme side.",
	})

	// Given
	expectedHref := "#pulje-FredagKveld"
	expectedHeadingText := "Fredag kveld (18:00 - 23:00)"
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen, true)
	insertRootPageEvent(t, db, "friday-raffle", "Friday Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", rootPageTestNow))
	actualHref := doc.Find(".pulje-alert a").AttrOr("href", "")
	actualHeadingText := templtest.CollectTexts(doc, expectedHref)

	// Then
	if actualHref != expectedHref {
		t.Fatalf("pulje alert link mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
	if !slices.Equal([]string{expectedHeadingText}, actualHeadingText) {
		t.Fatalf("linked pulje heading mismatch\nexpected: %q\nactual:   %v", expectedHeadingText, actualHeadingText)
	}
}

func TestRootPageContent_WhenClosingPuljeHasNoRaffleEvents_ShowsAlertWithoutLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at interessevalget for fredagens pulje stenger snart, men puljen bare har programarrangementer.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal varselet vises uten lenke, fordi puljen ikke har noen overskrift å gå til.",
	})

	// Given
	expectedAlerts := []string{closingAlertText}
	expectedLinkCount := 0
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen, true)
	insertRootPageEvent(t, db, "friday-program", "Friday Program", models.EventStatusAnnounced, false)
	insertRootPageEventPulje(t, db, "friday-program", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", rootPageTestNow))
	actualAlerts := templtest.CollectTexts(doc, ".pulje-alert")
	actualLinkCount := doc.Find(".pulje-alert a").Length()

	// Then
	if !slices.Equal(expectedAlerts, actualAlerts) {
		t.Fatalf("pulje alerts mismatch\nexpected: %v\nactual:   %v", expectedAlerts, actualAlerts)
	}
	if actualLinkCount != expectedLinkCount {
		t.Fatalf("pulje alert link count mismatch\nexpected: %d\nactual:   %d", expectedLinkCount, actualLinkCount)
	}
}

func TestRootPageContent_WhenPuljerHaveDifferentStatuses_MarksEachAlertWithItsStatus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at lørdagens ene pulje er låst og den andre er i ferd med å bli låst.",
		When:  "Når forsiden vises for lørdag.",
		Then:  "Så skal hvert varsel merkes med sin status, slik at det får riktig farge.",
	})

	// Given
	expectedClasses := []string{"pulje-alert is-locked", "pulje-alert is-closing"}
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")
	setRootPagePuljeStatus(t, db, models.PuljeLordagMorgen, models.PuljeStatusLocked, false)
	setRootPagePuljeStatus(t, db, models.PuljeLordagKveld, models.PuljeStatusOpen, true)
	insertRootPageEvent(t, db, "morning-raffle", "Morning Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "morning-raffle", models.PuljeLordagMorgen, true)
	insertRootPageEvent(t, db, "evening-raffle", "Evening Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "evening-raffle", models.PuljeLordagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-10", rootPageTestNow))
	actualClasses := make([]string, 0, len(expectedClasses))
	doc.Find(".pulje-alert").Each(func(_ int, alert *goquery.Selection) {
		actualClasses = append(actualClasses, alert.AttrOr("class", ""))
	})

	// Then
	if !slices.Equal(expectedClasses, actualClasses) {
		t.Fatalf("pulje alert classes mismatch\nexpected: %v\nactual:   %v", expectedClasses, actualClasses)
	}
}

func TestRootPageContent_WhenPuljerAreOpenWithoutWarning_ShowsNoPuljeAlerts(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fredagens pulje er åpen uten varsel om låsing.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal det ikke vises noe puljevarsel.",
	})

	// Given
	expectedAlertsVisible := false
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	insertRootPageEvent(t, db, "friday-raffle", "Friday Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", rootPageTestNow))
	actualAlertsVisible := templtest.HasSelector(doc, ".pulje-alert")

	// Then
	if actualAlertsVisible != expectedAlertsVisible {
		t.Fatalf("pulje alerts visibility mismatch\nexpected: %v\nactual:   %v", expectedAlertsVisible, actualAlertsVisible)
	}
}

func TestRootPageContent_WhenPuljeStartedMoreThanAnHourAgo_HidesItsAlert(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fredagens pulje har et varsel (stenger snart, låst eller klar) og startet for mer enn én time siden.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal puljevarselet ikke lenger vises.",
	})

	// Given
	expectedAlertsVisible := false
	puljeStart := time.Date(2026, 10, 9, 20, 0, 0, 0, programLocation(t))
	now := puljeStart.Add(puljeAlertHideAfterStart + time.Minute)
	cases := []struct {
		name                 string
		status               models.PuljeStatus
		closingWarningActive bool
	}{
		{name: "closing soon", status: models.PuljeStatusOpen, closingWarningActive: true},
		{name: "locked", status: models.PuljeStatusLocked},
		{name: "completed", status: models.PuljeStatusCompleted},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := createRootPageTestDB(t)
			seedRootPageLookups(t, db)
			setProgramPublishing(t, db, true)
			insertFridayPuljeStartingAt(t, db, puljeStart)
			setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, tc.status, tc.closingWarningActive)
			insertRootPageEvent(t, db, "friday-raffle", "Friday Raffle", models.EventStatusAnnounced, true)
			insertRootPageEventPulje(t, db, "friday-raffle", models.PuljeFredagKveld, true)

			// When
			doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", now))
			actualAlertsVisible := templtest.HasSelector(doc, ".pulje-alert")

			// Then
			if actualAlertsVisible != expectedAlertsVisible {
				t.Fatalf("pulje alerts visibility mismatch\nexpected: %v\nactual:   %v", expectedAlertsVisible, actualAlertsVisible)
			}
		})
	}
}

func TestRootPageContent_WhenPuljeStartedLessThanAnHourAgo_StillShowsItsAlert(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at fredagens pulje er låst og startet for mindre enn én time siden.",
		When:  "Når forsiden vises for fredag.",
		Then:  "Så skal puljevarselet fortsatt vises.",
	})

	// Given
	expectedAlerts := []string{lockedAlertText}
	puljeStart := time.Date(2026, 10, 9, 20, 0, 0, 0, programLocation(t))
	now := puljeStart.Add(puljeAlertHideAfterStart - time.Minute)
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertFridayPuljeStartingAt(t, db, puljeStart)
	setRootPagePuljeStatus(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked, false)
	insertRootPageEvent(t, db, "friday-raffle", "Friday Raffle", models.EventStatusAnnounced, true)
	insertRootPageEventPulje(t, db, "friday-raffle", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, nil, "2026-10-09", now))
	actualAlerts := templtest.CollectTexts(doc, ".pulje-alert")

	// Then
	if !slices.Equal(expectedAlerts, actualAlerts) {
		t.Fatalf("pulje alerts mismatch\nexpected: %v\nactual:   %v", expectedAlerts, actualAlerts)
	}
}

func TestPuljeAlertFor_DescribesEachPuljeState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje med en status og et valgfritt varsel om låsing.",
		When:  "Når puljevarselet bestemmes.",
		Then:  "Så skal tekst, farge og ikon beskrive puljens tilstand, og åpne puljer uten varsel skal ikke få varsel.",
	})

	// Given
	closingAlert := puljeAlert{Message: closingAlertText, Class: "is-closing", Icon: icons.Warning, LinkHref: "#pulje-FredagKveld"}
	lockedAlert := puljeAlert{Message: lockedAlertText, Class: "is-locked", Icon: icons.ClockLock}
	completedAlert := puljeAlert{Message: completedAlertText, Class: "is-completed", Icon: icons.ProgressComplete}
	cases := []struct {
		name                 string
		status               models.PuljeStatus
		closingWarningActive bool
		expectedAlert        puljeAlert
		expectedOK           bool
	}{
		{name: "open", status: models.PuljeStatusOpen},
		{name: "open with closing warning", status: models.PuljeStatusOpen, closingWarningActive: true, expectedAlert: closingAlert, expectedOK: true},
		{name: "locked", status: models.PuljeStatusLocked, expectedAlert: lockedAlert, expectedOK: true},
		{name: "locked ignores closing warning", status: models.PuljeStatusLocked, closingWarningActive: true, expectedAlert: lockedAlert, expectedOK: true},
		{name: "completed", status: models.PuljeStatusCompleted, expectedAlert: completedAlert, expectedOK: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pulje := models.PuljeRow{ID: models.PuljeFredagKveld, Status: tc.status, ClosingWarningActive: tc.closingWarningActive}

			// When
			actualAlert, actualOK := puljeAlertFor(pulje)

			// Then
			if actualOK != tc.expectedOK || actualAlert != tc.expectedAlert {
				t.Fatalf("pulje alert mismatch\nexpected: %+v, %v\nactual:   %+v, %v", tc.expectedAlert, tc.expectedOK, actualAlert, actualOK)
			}
		})
	}
}

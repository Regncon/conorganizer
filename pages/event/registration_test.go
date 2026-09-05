package event

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

type registrationTestFixture struct {
	userExternalID  string
	billettholderID int
	eventID         string
	otherEventID    string
	puljeID         models.Pulje
}

func TestSetRegistration_WhenEligible_CreatesRegistrationAndRemovesOnlyMatchingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en kvalifisert billetthelder med interesser i to arrangementer i samme pulje.",
		When:  "Når billetthelderen melder seg på ett arrangement to ganger.",
		Then:  "Så skal én påmelding lagres og bare interessen for det arrangementet fjernes.",
	})

	// Given
	expectedSeatCount := 1
	expectedSource := models.EventPlayerSourceRegistration
	expectedMatchingInterests := 0
	expectedOtherInterests := 1
	db := createEventInterestTestDB(t)
	fixture := seedRegistrationTestFixture(t, db)
	change := registrationChange{
		UserExternalID:  fixture.userExternalID,
		BillettholderID: fixture.billettholderID,
		EventID:         fixture.eventID,
		PuljeID:         fixture.puljeID,
		IsRegistered:    true,
	}

	// When
	if err := setRegistration(context.Background(), db, change); err != nil {
		t.Fatalf("expected registration to succeed: %v", err)
	}
	if err := setRegistration(context.Background(), db, change); err != nil {
		t.Fatalf("expected repeated registration to be idempotent: %v", err)
	}
	actualSeatCount := registrationSeatCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualSource := registrationSeatSource(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualMatchingInterests := registrationInterestCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualOtherInterests := registrationInterestCount(t, db, fixture.otherEventID, fixture.puljeID, fixture.billettholderID)

	// Then
	if actualSeatCount != expectedSeatCount {
		t.Fatalf("seat count mismatch\nexpected: %d\nactual:   %d", expectedSeatCount, actualSeatCount)
	}
	if actualSource != expectedSource {
		t.Fatalf("source mismatch\nexpected: %s\nactual:   %s", expectedSource, actualSource)
	}
	if actualMatchingInterests != expectedMatchingInterests {
		t.Fatalf("matching interest count mismatch\nexpected: %d\nactual:   %d", expectedMatchingInterests, actualMatchingInterests)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}

func TestSetRegistration_WhenRegisteringForMultipleEventsInPulje_KeepsEveryRegistration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt to arrangementer med åpen påmelding i samme pulje.",
		When:  "Når billetthelderen melder seg på begge arrangementene.",
		Then:  "Så skal begge påmeldingene beholdes.",
	})

	// Given
	expectedRegistrationCount := 2
	db := createEventInterestTestDB(t)
	fixture := seedRegistrationTestFixture(t, db)

	// When
	for _, eventID := range []string{fixture.eventID, fixture.otherEventID} {
		if err := setRegistration(context.Background(), db, registrationChange{
			UserExternalID:  fixture.userExternalID,
			BillettholderID: fixture.billettholderID,
			EventID:         eventID,
			PuljeID:         fixture.puljeID,
			IsRegistered:    true,
		}); err != nil {
			t.Fatalf("expected registration for %s to succeed: %v", eventID, err)
		}
	}
	actualRegistrationCount := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_events_players
		WHERE billettholder_id = ?
			AND pulje_id = ?
			AND role = ?
			AND source = ?
	`, fixture.billettholderID, fixture.puljeID, models.EventPlayerRolePlayer, models.EventPlayerSourceRegistration)

	// Then
	if actualRegistrationCount != expectedRegistrationCount {
		t.Fatalf("registration count mismatch\nexpected: %d\nactual:   %d", expectedRegistrationCount, actualRegistrationCount)
	}
}

func TestSetRegistration_WhenDeregisteringManualAssignment_RemovesSeatAndMatchingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en manuell spillertildeling til et vanlig arrangement.",
		When:  "Når billetthelderen melder seg av.",
		Then:  "Så skal tildelingen og den samsvarende interessen fjernes uten å røre andre interesser.",
	})

	// Given
	expectedSeatCount := 0
	expectedMatchingInterests := 0
	expectedOtherInterests := 1
	db := createEventInterestTestDB(t)
	fixture := seedRegistrationTestFixture(t, db)
	testutil.MustExec(t, db, `UPDATE events SET is_open_registration = 0 WHERE id = ?`, fixture.eventID)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
	`, fixture.eventID, fixture.puljeID, fixture.billettholderID, models.EventPlayerRolePlayer, models.EventPlayerSourceManual)

	// When
	err := setRegistration(context.Background(), db, registrationChange{
		UserExternalID:  fixture.userExternalID,
		BillettholderID: fixture.billettholderID,
		EventID:         fixture.eventID,
		PuljeID:         fixture.puljeID,
		IsRegistered:    false,
	})
	actualSeatCount := registrationSeatCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualMatchingInterests := registrationInterestCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualOtherInterests := registrationInterestCount(t, db, fixture.otherEventID, fixture.puljeID, fixture.billettholderID)

	// Then
	if err != nil {
		t.Fatalf("expected deregistration from a manual assignment to succeed: %v", err)
	}
	if actualSeatCount != expectedSeatCount {
		t.Fatalf("seat count mismatch\nexpected: %d\nactual:   %d", expectedSeatCount, actualSeatCount)
	}
	if actualMatchingInterests != expectedMatchingInterests {
		t.Fatalf("matching interest count mismatch\nexpected: %d\nactual:   %d", expectedMatchingInterests, actualMatchingInterests)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}

func TestSetRegistration_WhenRegistrationRuleRejects_KeepsInterestAndSeatUnchanged(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder som ikke oppfyller ett av kravene for påmelding.",
		When:  "Når billetthelderen forsøker å melde seg på.",
		Then:  "Så skal forsøket avvises uten å endre interesser eller tildelinger.",
	})

	testCases := []struct {
		name        string
		expectedErr error
		mutate      func(t *testing.T, db *sql.DB, fixture registrationTestFixture)
	}{
		{
			name:        "program is not published",
			expectedErr: errRegistrationProgramNotPublished,
			mutate: func(t *testing.T, db *sql.DB, _ registrationTestFixture) {
				setEventInterestProgramPublishing(t, db, false)
			},
		},
		{
			name:        "event occurrence is not published",
			expectedErr: errRegistrationEventUnavailable,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `UPDATE relation_event_puljer SET is_published = 0 WHERE event_id = ? AND pulje_id = ?`, fixture.eventID, fixture.puljeID)
			},
		},
		{
			name:        "pulje is locked",
			expectedErr: errRegistrationPuljeNotOpen,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusLocked, fixture.puljeID)
			},
		},
		{
			name:        "event does not have open registration",
			expectedErr: errRegistrationEventNotOpen,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `UPDATE events SET is_open_registration = 0 WHERE id = ?`, fixture.eventID)
			},
		},
		{
			name:        "billettholder is underage for adults-only event",
			expectedErr: errRegistrationAdultsOnly,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupAdultsOnly)
				testutil.MustExec(t, db, `UPDATE events SET age_group = ? WHERE id = ?`, models.AgeGroupAdultsOnly, fixture.eventID)
				testutil.MustExec(t, db, `UPDATE billettholdere SET is_over_18 = 0 WHERE id = ?`, fixture.billettholderID)
			},
		},
		{
			name:        "user does not own billettholder",
			expectedErr: errRegistrationAccessDenied,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `DELETE FROM relation_billettholdere_users WHERE billettholder_id = ?`, fixture.billettholderID)
			},
		},
		{
			name:        "billettholder is gamemaster in pulje",
			expectedErr: errRegistrationGamemaster,
			mutate: func(t *testing.T, db *sql.DB, fixture registrationTestFixture) {
				testutil.MustExec(t, db, `
					INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
					VALUES (?, ?, ?, ?, ?)
				`, fixture.otherEventID, fixture.puljeID, fixture.billettholderID, models.EventPlayerRoleGM, models.EventPlayerSourceManual)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			expectedSeatCount := 0
			expectedInterestCount := 1
			db := createEventInterestTestDB(t)
			fixture := seedRegistrationTestFixture(t, db)
			testCase.mutate(t, db, fixture)

			// When
			err := setRegistration(context.Background(), db, registrationChange{
				UserExternalID:  fixture.userExternalID,
				BillettholderID: fixture.billettholderID,
				EventID:         fixture.eventID,
				PuljeID:         fixture.puljeID,
				IsRegistered:    true,
			})
			actualSeatCount := registrationSeatCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
			actualInterestCount := registrationInterestCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)

			// Then
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("error mismatch\nexpected: %v\nactual:   %v", testCase.expectedErr, err)
			}
			if actualSeatCount != expectedSeatCount {
				t.Fatalf("seat count mismatch\nexpected: %d\nactual:   %d", expectedSeatCount, actualSeatCount)
			}
			if actualInterestCount != expectedInterestCount {
				t.Fatalf("interest count mismatch\nexpected: %d\nactual:   %d", expectedInterestCount, actualInterestCount)
			}
		})
	}
}

func TestSetRegistration_WhenPuljeIsLocked_RejectsDeregistration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en eksisterende påmelding i en låst pulje.",
		When:  "Når billetthelderen forsøker å melde seg av.",
		Then:  "Så skal påmeldingen og interessen forbli uendret.",
	})

	// Given
	expectedSeatCount := 1
	expectedInterestCount := 1
	db := createEventInterestTestDB(t)
	fixture := seedRegistrationTestFixture(t, db)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES (?, ?, ?, ?, ?)
	`, fixture.eventID, fixture.puljeID, fixture.billettholderID, models.EventPlayerRolePlayer, models.EventPlayerSourceRegistration)
	testutil.MustExec(t, db, `UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusLocked, fixture.puljeID)

	// When
	err := setRegistration(context.Background(), db, registrationChange{
		UserExternalID:  fixture.userExternalID,
		BillettholderID: fixture.billettholderID,
		EventID:         fixture.eventID,
		PuljeID:         fixture.puljeID,
		IsRegistered:    false,
	})
	actualSeatCount := registrationSeatCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)
	actualInterestCount := registrationInterestCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)

	// Then
	if !errors.Is(err, errRegistrationPuljeNotOpen) {
		t.Fatalf("error mismatch\nexpected: %v\nactual:   %v", errRegistrationPuljeNotOpen, err)
	}
	if actualSeatCount != expectedSeatCount {
		t.Fatalf("seat count mismatch\nexpected: %d\nactual:   %d", expectedSeatCount, actualSeatCount)
	}
	if actualInterestCount != expectedInterestCount {
		t.Fatalf("interest count mismatch\nexpected: %d\nactual:   %d", expectedInterestCount, actualInterestCount)
	}
}

func TestRegistrationRoute_WhenRequestIsEligible_RegistersBillettholder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker med en kvalifisert billetthelder.",
		When:  "Når registreringsendepunktet mottar ønske om påmelding.",
		Then:  "Så skal billetthelderen registreres og endepunktet svare uten innhold.",
	})

	// Given
	expectedStatusCode := http.StatusNoContent
	expectedSeatCount := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "registration_route")
	seedEventInterestLookups(t, db)
	fixture := seedRegistrationTestFixture(t, db)
	router := chi.NewRouter()
	router.Route("/event/api/{idx}/registration", func(registrationRouter chi.Router) {
		setupRegistrationRoute(registrationRouter, db, &live.Manager{}, logger)
	})
	request := httptest.NewRequest(
		http.MethodPut,
		"/event/api/"+fixture.eventID+"/registration",
		strings.NewReader(`{"billettHolderId":901,"puljeId":"FredagKveld","isRegistered":true}`),
	)
	request = request.WithContext(authctx.WithUserToken(request.Context(), fixture.userExternalID, "registration-user@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)
	actualSeatCount := registrationSeatCount(t, db, fixture.eventID, fixture.puljeID, fixture.billettholderID)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d\nbody: %s", expectedStatusCode, recorder.Code, recorder.Body.String())
	}
	if actualSeatCount != expectedSeatCount {
		t.Fatalf("seat count mismatch\nexpected: %d\nactual:   %d", expectedSeatCount, actualSeatCount)
	}
}

func seedRegistrationTestFixture(t *testing.T, db *sql.DB) registrationTestFixture {
	t.Helper()

	fixture := registrationTestFixture{
		userExternalID:  "registration-user",
		billettholderID: 901,
		eventID:         "registration-event",
		otherEventID:    "other-registration-event",
		puljeID:         models.PuljeFredagKveld,
	}

	testutil.MustExec(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (601, ?, 'registration-user@example.com', 0)
	`, fixture.userExternalID)
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, 'Regina', 'Strant', 1, 'Ticket', 1, 7101, 8101)
	`, fixture.billettholderID)
	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		VALUES (?, 601)
	`, fixture.billettholderID)
	testutil.MustExec(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:30:00+02:00', '2026-10-09T23:00:00+02:00')
	`, fixture.puljeID, models.PuljeStatusOpen)

	for _, eventID := range []string{fixture.eventID, fixture.otherEventID} {
		testutil.MustExec(t, db, `
			INSERT INTO events (
				id, title, intro, description, system, event_type,
				age_group, event_runtime, host_name, email, phone_number,
				max_players, beginner_friendly, can_be_run_in_english,
				status, is_open_registration
			) VALUES (?, 'Registration Event', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?, 1)
		`, eventID, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusAnnounced)
		testutil.MustExec(t, db, `
			INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
			VALUES (?, ?, 1, 1)
		`, eventID, fixture.puljeID)
	}

	testutil.MustExec(t, db, `
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?), (?, ?, ?, ?)
	`,
		fixture.billettholderID, fixture.eventID, fixture.puljeID, models.InterestLevelHigh,
		fixture.billettholderID, fixture.otherEventID, fixture.puljeID, models.InterestLevelMedium,
	)

	return fixture
}

func registrationSeatCount(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, billettholderID int) int {
	t.Helper()
	return testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND role = ?
	`, eventID, puljeID, billettholderID, models.EventPlayerRolePlayer)
}

func registrationSeatSource(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, billettholderID int) models.EventPlayerSource {
	t.Helper()
	var source models.EventPlayerSource
	if err := db.QueryRow(`
		SELECT source
		FROM relation_events_players
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND role = ?
	`, eventID, puljeID, billettholderID, models.EventPlayerRolePlayer).Scan(&source); err != nil {
		t.Fatalf("read seat source: %v", err)
	}
	return source
}

func registrationInterestCount(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, billettholderID int) int {
	t.Helper()
	return testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM interests
		WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?
	`, eventID, puljeID, billettholderID)
}

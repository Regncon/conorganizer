package ticketholder

import (
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
)

func TestBuildPuljeInterestState_WhenPuljeIsLocked_ReturnsLockedStateAndDisablesEditing(t *testing.T) {
	// Gitt at en pulje er låst,
	// når interessetilstanden bygges,
	// så skal billettholderen se låst status og ikke kunne redigere interessen.

	// Given
	expectedAvailability := PuljeInterestLocked
	expectedCanEdit := false
	expectedMessagePart := "Vi jobber med å fordele spillere."

	pulje := buildPuljeInterestStateTestPulje(
		t,
		models.PuljeFredagKveld,
		"Fredag kveld",
		models.PuljeStatusLocked,
		"2026-10-09T18:30:00+02:00",
	)
	now := parsePuljeInterestStateTestTime(t, "2026-10-09T15:00:00+02:00")

	// When
	actualState := BuildPuljeInterestState(pulje, now)

	// Then
	if actualState.Availability != expectedAvailability {
		t.Fatalf("pulje availability mismatch\nexpected: %s\nactual:   %s", expectedAvailability, actualState.Availability)
	}
	if actualState.CanEdit != expectedCanEdit {
		t.Fatalf("can edit mismatch\nexpected: %v\nactual:   %v", expectedCanEdit, actualState.CanEdit)
	}
	if !strings.Contains(actualState.Message, expectedMessagePart) {
		t.Fatalf("locked message mismatch\nexpected to contain: %q\nactual:              %q", expectedMessagePart, actualState.Message)
	}
}

func TestBuildPuljeInterestState_WhenOpenPuljeIsInWarningWindow_ReturnsWarningWithLockTime(t *testing.T) {
	// Gitt at en åpen pulje nærmer seg låsing,
	// når interessetilstanden bygges,
	// så skal billettholderen se en advarsel med tidspunktet puljen låses.

	// Given
	expectedAvailability := PuljeInterestWarning
	expectedMessage := "Puljen låses snart, kl 18:00."

	pulje := buildPuljeInterestStateTestPulje(
		t,
		models.PuljeFredagKveld,
		"Fredag kveld",
		models.PuljeStatusOpen,
		"2026-10-09T18:30:00+02:00",
	)
	now := parsePuljeInterestStateTestTime(t, "2026-10-09T16:15:00+02:00")

	// When
	actualState := BuildPuljeInterestState(pulje, now)

	// Then
	if actualState.Availability != expectedAvailability {
		t.Fatalf("pulje availability mismatch\nexpected: %s\nactual:   %s", expectedAvailability, actualState.Availability)
	}
	if actualState.Message != expectedMessage {
		t.Fatalf("warning message mismatch\nexpected: %q\nactual:   %q", expectedMessage, actualState.Message)
	}
}

func TestBuildPuljeInterestState_WhenOpenPuljeIsInUrgentWarningWindow_ReturnsUrgentWarningWithLockTime(t *testing.T) {
	// Gitt at en åpen pulje er svært nær låsing,
	// når interessetilstanden bygges,
	// så skal billettholderen se en tydelig hasteadvarsel.

	// Given
	expectedAvailability := PuljeInterestUrgentWarning
	expectedMessage := "Puljen låses straks, kl 18:00. Gjør endringer nå hvis du vil endre interessen din."

	pulje := buildPuljeInterestStateTestPulje(
		t,
		models.PuljeFredagKveld,
		"Fredag kveld",
		models.PuljeStatusOpen,
		"2026-10-09T18:30:00+02:00",
	)
	now := parsePuljeInterestStateTestTime(t, "2026-10-09T17:45:00+02:00")

	// When
	actualState := BuildPuljeInterestState(pulje, now)

	// Then
	if actualState.Availability != expectedAvailability {
		t.Fatalf("pulje availability mismatch\nexpected: %s\nactual:   %s", expectedAvailability, actualState.Availability)
	}
	if actualState.Message != expectedMessage {
		t.Fatalf("urgent warning message mismatch\nexpected: %q\nactual:   %q", expectedMessage, actualState.Message)
	}
}

func TestBuildMostUrgentPuljeInterestState_WhenWarningAndLockedPuljerExist_ReturnsWarningState(t *testing.T) {
	// Gitt at noen puljer er låst og en åpen pulje snart låses,
	// når den viktigste meldingen velges,
	// så skal tidsadvarselen vises i stedet for låst status.

	// Given
	expectedHasState := true
	expectedPuljeID := models.PuljeLordagMorgen
	expectedAvailability := PuljeInterestUrgentWarning

	now := parsePuljeInterestStateTestTime(t, "2026-10-10T09:15:00+02:00")
	puljer := []models.PuljeRow{
		buildPuljeInterestStateTestPulje(
			t,
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusLocked,
			"2026-10-09T18:30:00+02:00",
		),
		buildPuljeInterestStateTestPulje(
			t,
			models.PuljeLordagMorgen,
			"Lørdag morgen",
			models.PuljeStatusOpen,
			"2026-10-10T10:00:00+02:00",
		),
	}

	// When
	actualState, actualHasState := BuildMostUrgentPuljeInterestState(puljer, now)

	// Then
	if actualHasState != expectedHasState {
		t.Fatalf("has urgent state mismatch\nexpected: %v\nactual:   %v", expectedHasState, actualHasState)
	}
	if actualState.PuljeID != expectedPuljeID {
		t.Fatalf("pulje id mismatch\nexpected: %s\nactual:   %s", expectedPuljeID, actualState.PuljeID)
	}
	if actualState.Availability != expectedAvailability {
		t.Fatalf("pulje availability mismatch\nexpected: %s\nactual:   %s", expectedAvailability, actualState.Availability)
	}
}

func buildPuljeInterestStateTestPulje(t *testing.T, id models.Pulje, name string, status models.PuljeStatus, startAt string) models.PuljeRow {
	t.Helper()

	startTime, err := time.Parse(time.RFC3339, startAt)
	if err != nil {
		t.Fatalf("failed to parse test start_at %q: %v", startAt, err)
	}

	return models.PuljeRow{
		ID:      id,
		Name:    name,
		Status:  status,
		StartAt: models.NewDBDateTime(startTime),
		EndAt:   models.NewDBDateTime(startTime.Add(4 * time.Hour)),
	}
}

func parsePuljeInterestStateTestTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse test time %q: %v", value, err)
	}
	return parsed
}

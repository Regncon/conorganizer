package ticketholder

import (
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestNewBillettholderOptions_WhenUserHasMatchingBillettholder_SelectsItAsDefault(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren har flere tilknyttede billettholdere og én matcher brukerens e-post.",
		When:  "Når billettholdervalgene bygges.",
		Then:  "Så skal den matchende billettholderen være standard og brukeren skal kunne bytte.",
	})

	// Given
	expectedDefaultID := 2
	expectedHasBillettholder := true
	expectedCanSwitchBillettholder := true
	userInfo := requestctx.UserRequestInfo{Email: "user@example.com"}
	associated := []BillettHolder{
		{Id: 1, Email: "other@example.com"},
		{Id: expectedDefaultID, Email: userInfo.Email},
	}

	// When
	actualOptions := NewBillettholderOptions(userInfo, associated)

	// Then
	if actualOptions.Default.Id != expectedDefaultID {
		t.Fatalf("default billettholder ID mismatch\nexpected: %d\nactual:   %d", expectedDefaultID, actualOptions.Default.Id)
	}
	if actualOptions.HasBillettholder() != expectedHasBillettholder {
		t.Fatalf("has billettholder mismatch\nexpected: %v\nactual:   %v", expectedHasBillettholder, actualOptions.HasBillettholder())
	}
	if actualOptions.CanSwitchBillettholder() != expectedCanSwitchBillettholder {
		t.Fatalf("can switch billettholder mismatch\nexpected: %v\nactual:   %v", expectedCanSwitchBillettholder, actualOptions.CanSwitchBillettholder())
	}
}

func TestNewBillettholderOptions_WhenNoEmailMatches_UsesFirstBillettholderAsDefault(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at ingen tilknyttet billettholder matcher brukerens e-post.",
		When:  "Når billettholdervalgene bygges.",
		Then:  "Så skal den første tilknyttede billettholderen brukes som standard.",
	})

	// Given
	expectedDefaultID := 1
	userInfo := requestctx.UserRequestInfo{Email: "user@example.com"}
	associated := []BillettHolder{
		{Id: expectedDefaultID, Email: "first@example.com"},
		{Id: 2, Email: "second@example.com"},
	}

	// When
	actualOptions := NewBillettholderOptions(userInfo, associated)

	// Then
	if actualOptions.Default.Id != expectedDefaultID {
		t.Fatalf("default billettholder ID mismatch\nexpected: %d\nactual:   %d", expectedDefaultID, actualOptions.Default.Id)
	}
}

func TestNewBillettholderOptions_WhenNoBillettholdereExist_ReturnsUnavailableOptions(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke har noen tilknyttede billettholdere.",
		When:  "Når billettholdervalgene bygges.",
		Then:  "Så skal ingen billettholder være tilgjengelig og bytting være deaktivert.",
	})

	// Given
	expectedDefaultID := 0
	expectedHasBillettholder := false
	expectedCanSwitchBillettholder := false
	userInfo := requestctx.UserRequestInfo{Email: "user@example.com"}

	// When
	actualOptions := NewBillettholderOptions(userInfo, nil)

	// Then
	if actualOptions.Default.Id != expectedDefaultID {
		t.Fatalf("default billettholder ID mismatch\nexpected: %d\nactual:   %d", expectedDefaultID, actualOptions.Default.Id)
	}
	if actualOptions.HasBillettholder() != expectedHasBillettholder {
		t.Fatalf("has billettholder mismatch\nexpected: %v\nactual:   %v", expectedHasBillettholder, actualOptions.HasBillettholder())
	}
	if actualOptions.CanSwitchBillettholder() != expectedCanSwitchBillettholder {
		t.Fatalf("can switch billettholder mismatch\nexpected: %v\nactual:   %v", expectedCanSwitchBillettholder, actualOptions.CanSwitchBillettholder())
	}
}

func TestBuildPuljeInterestState_WhenPuljeIsLocked_ReturnsLockedStateAndDisablesEditing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en pulje er låst.",
		When:  "Når interessetilstanden bygges.",
		Then:  "Så skal billettholderen se låst status og ikke kunne redigere interessen.",
	})

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

func TestBuildPuljeInterestState_WhenOpenPuljeIsBeforeWarningWindow_ReturnsOpenStateWithoutWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en åpen pulje ikke nærmer seg låsing.",
		When:  "Når interessetilstanden bygges.",
		Then:  "Så skal billettholderen ikke se noen låseadvarsel.",
	})

	// Given
	expectedAvailability := PuljeInterestOpen
	expectedCanEdit := true
	expectedMessage := ""

	pulje := buildPuljeInterestStateTestPulje(
		t,
		models.PuljeFredagKveld,
		"Fredag kveld",
		models.PuljeStatusOpen,
		"2026-10-09T18:30:00+02:00",
	)
	now := parsePuljeInterestStateTestTime(t, "2026-10-09T15:59:00+02:00")

	// When
	actualState := BuildPuljeInterestState(pulje, now)

	// Then
	if actualState.Availability != expectedAvailability {
		t.Fatalf("pulje availability mismatch\nexpected: %s\nactual:   %s", expectedAvailability, actualState.Availability)
	}
	if actualState.CanEdit != expectedCanEdit {
		t.Fatalf("can edit mismatch\nexpected: %v\nactual:   %v", expectedCanEdit, actualState.CanEdit)
	}
	if actualState.Message != expectedMessage {
		t.Fatalf("message mismatch\nexpected: %q\nactual:   %q", expectedMessage, actualState.Message)
	}
}

func TestBuildPuljeInterestState_WhenOpenPuljeIsInWarningWindow_ReturnsWarningWithLockTime(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en åpen pulje nærmer seg låsing.",
		When:  "Når interessetilstanden bygges.",
		Then:  "Så skal billettholderen se en advarsel med tidspunktet puljen låses.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en åpen pulje er svært nær låsing.",
		When:  "Når interessetilstanden bygges.",
		Then:  "Så skal billettholderen se en tydelig hasteadvarsel.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at noen puljer er låst og en åpen pulje snart låses.",
		When:  "Når den viktigste meldingen velges.",
		Then:  "Så skal tidsadvarselen vises i stedet for låst status.",
	})

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

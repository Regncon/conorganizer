package event

import (
	"reflect"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/nats-io/nats.go/jetstream"
)

func TestPuljeScheduledBroadcastJetStreamConfig_UsesMemoryStorage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the pulje scheduled broadcast JetStream configuration.",
		When:  "When the stream and consumer configs are built.",
		Then:  "Then both the scheduled stream and durable consumer keep state in memory.",
	})

	// Given
	expectedStreamStorage := jetstream.MemoryStorage
	expectedConsumerMemoryStorage := true

	// When
	streamConfig := puljeScheduledBroadcastStreamConfig()
	consumerConfig := puljeScheduledBroadcastConsumerConfig()

	// Then
	if streamConfig.Storage != expectedStreamStorage {
		t.Fatalf("stream storage mismatch\nexpected: %s\nactual:   %s", expectedStreamStorage, streamConfig.Storage)
	}
	if consumerConfig.MemoryStorage != expectedConsumerMemoryStorage {
		t.Fatalf("expected consumer memory storage to be enabled")
	}
}

func TestBuildPuljeScheduledBroadcasts_SchedulesWarningAndUrgentThresholdsButNotLockTime(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en pulje starter senere enn varslingsvinduene.",
		When:  "Når planlagte puljeoppdateringer bygges.",
		Then:  "Så skal bare varsel og hastevarsel planlegges.",
	})

	// Given
	expectedBroadcasts := []expectedPuljeScheduledBroadcast{
		{PuljeID: models.PuljeFredagKveld, Threshold: "warning", ScheduleAt: "2026-10-09T16:00:00+02:00"},
		{PuljeID: models.PuljeFredagKveld, Threshold: "urgent", ScheduleAt: "2026-10-09T17:30:00+02:00"},
	}

	now := parsePuljeScheduledBroadcastTestTime(t, "2026-10-09T15:00:00+02:00")
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			parsePuljeScheduledBroadcastTestTime(t, "2026-10-09T18:30:00+02:00"),
		),
	}

	// When
	actualBroadcasts := normalizePuljeScheduledBroadcasts(buildPuljeScheduledBroadcasts(puljer, now))

	// Then
	if !reflect.DeepEqual(expectedBroadcasts, actualBroadcasts) {
		t.Fatalf("scheduled broadcasts mismatch\nexpected: %#v\nactual:   %#v", expectedBroadcasts, actualBroadcasts)
	}
}

func TestBuildPuljeScheduledBroadcasts_WhenWarningIsPast_OnlySchedulesFutureUrgentThreshold(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at varselgrensen allerede har passert.",
		When:  "Når planlagte puljeoppdateringer bygges.",
		Then:  "Så skal bare fremtidige terskler planlegges.",
	})

	// Given
	expectedBroadcasts := []expectedPuljeScheduledBroadcast{
		{PuljeID: models.PuljeFredagKveld, Threshold: "urgent", ScheduleAt: "2026-10-09T17:30:00+02:00"},
	}

	now := parsePuljeScheduledBroadcastTestTime(t, "2026-10-09T16:30:00+02:00")
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			parsePuljeScheduledBroadcastTestTime(t, "2026-10-09T18:30:00+02:00"),
		),
	}

	// When
	actualBroadcasts := normalizePuljeScheduledBroadcasts(buildPuljeScheduledBroadcasts(puljer, now))

	// Then
	if !reflect.DeepEqual(expectedBroadcasts, actualBroadcasts) {
		t.Fatalf("scheduled broadcasts mismatch\nexpected: %#v\nactual:   %#v", expectedBroadcasts, actualBroadcasts)
	}
}

type expectedPuljeScheduledBroadcast struct {
	PuljeID    models.Pulje
	Threshold  string
	ScheduleAt string
}

func normalizePuljeScheduledBroadcasts(broadcasts []puljeScheduledBroadcast) []expectedPuljeScheduledBroadcast {
	normalized := make([]expectedPuljeScheduledBroadcast, 0, len(broadcasts))
	for _, broadcast := range broadcasts {
		normalized = append(normalized, expectedPuljeScheduledBroadcast{
			PuljeID:    broadcast.PuljeID,
			Threshold:  broadcast.Threshold,
			ScheduleAt: broadcast.ScheduleAt.Format(time.RFC3339),
		})
	}
	return normalized
}

func parsePuljeScheduledBroadcastTestTime(t *testing.T, value string) time.Time {
	t.Helper()

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatalf("failed to parse test time %q: %v", value, err)
	}
	return parsed
}

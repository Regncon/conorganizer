package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/keyvalue"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	puljeScheduledBroadcastStream          = "EVENT_PULJE_SCHEDULES"
	puljeScheduledBroadcastScheduleSubject = "events.pulje.schedule"
	puljeScheduledBroadcastTargetSubject   = "events.pulje.schedule.due"
	puljeScheduledBroadcastConsumer        = "events-pulje-schedule-broadcast"
)

var puljeScheduledBroadcastConsumeContext jetstream.ConsumeContext

type puljeScheduledBroadcast struct {
	PuljeID    models.Pulje `json:"pulje_id"`
	Threshold  string       `json:"threshold"`
	ScheduleAt time.Time    `json:"schedule_at"`
}

func setupPuljeScheduledBroadcasts(ctx context.Context, js jetstream.JetStream, kv jetstream.KeyValue, db *sql.DB, logger *slog.Logger) error {
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:              puljeScheduledBroadcastStream,
		Description:       "Scheduled event page updates for pulje warning thresholds",
		Subjects:          []string{puljeScheduledBroadcastScheduleSubject, puljeScheduledBroadcastTargetSubject},
		AllowMsgSchedules: true,
		Retention:         jetstream.LimitsPolicy,
		MaxAge:            400 * 24 * time.Hour,
		MaxBytes:          1024 * 1024,
		Duplicates:        400 * 24 * time.Hour,
	})
	if err != nil {
		return fmt.Errorf("create pulje schedule stream: %w", err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          puljeScheduledBroadcastConsumer,
		Durable:       puljeScheduledBroadcastConsumer,
		Description:   "Broadcasts event page updates when pulje warning thresholds are reached",
		FilterSubject: puljeScheduledBroadcastTargetSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("create pulje schedule consumer: %w", err)
	}

	consumeContext, err := consumer.Consume(func(msg jetstream.Msg) {
		if err := keyvalue.BroadcastUpdateContext(context.Background(), kv); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast scheduled pulje update: %w", err).Error())
			if nakErr := msg.NakWithDelay(10 * time.Second); nakErr != nil {
				logger.Error(fmt.Errorf("failed to nack scheduled pulje update: %w", nakErr).Error())
			}
			return
		}
		if err := msg.Ack(); err != nil {
			logger.Error(fmt.Errorf("failed to ack scheduled pulje update: %w", err).Error())
		}
	}, jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
		logger.Error(fmt.Errorf("pulje scheduled update consumer error: %w", err).Error())
	}))
	if err != nil {
		return fmt.Errorf("consume pulje schedule messages: %w", err)
	}
	puljeScheduledBroadcastConsumeContext = consumeContext

	puljer, err := puljerService.GetAllPuljer(db)
	if err != nil {
		return fmt.Errorf("load puljer for scheduling: %w", err)
	}

	for _, broadcast := range buildPuljeScheduledBroadcasts(puljer, time.Now()) {
		payload, err := json.Marshal(broadcast)
		if err != nil {
			return fmt.Errorf("marshal pulje scheduled broadcast: %w", err)
		}
		if _, err := js.Publish(
			ctx,
			puljeScheduledBroadcastScheduleSubject,
			payload,
			jetstream.WithScheduleAt(broadcast.ScheduleAt),
			jetstream.WithScheduleTarget(puljeScheduledBroadcastTargetSubject),
			jetstream.WithMsgID(puljeScheduledBroadcastMessageID(broadcast)),
		); err != nil {
			return fmt.Errorf("schedule pulje %s %s broadcast at %s: %w", broadcast.PuljeID, broadcast.Threshold, broadcast.ScheduleAt.Format(time.RFC3339), err)
		}
		logger.Debug(
			"Scheduled pulje event page update",
			"pulje_id", broadcast.PuljeID,
			"threshold", broadcast.Threshold,
			"schedule_at", broadcast.ScheduleAt.Format(time.RFC3339),
		)
	}

	return nil
}

func buildPuljeScheduledBroadcasts(puljer []models.PuljeRow, now time.Time) []puljeScheduledBroadcast {
	broadcasts := make([]puljeScheduledBroadcast, 0, len(puljer)*2)

	for _, pulje := range puljer {
		if pulje.StartAt.IsZero() {
			continue
		}

		lockAt := pulje.StartAt.TimeOrZero().Add(-30 * time.Minute)
		warningAt := lockAt.Add(-2 * time.Hour)
		urgentAt := lockAt.Add(-30 * time.Minute)

		if warningAt.After(now) {
			broadcasts = append(broadcasts, puljeScheduledBroadcast{
				PuljeID:    pulje.ID,
				Threshold:  "warning",
				ScheduleAt: warningAt,
			})
		}
		if urgentAt.After(now) {
			broadcasts = append(broadcasts, puljeScheduledBroadcast{
				PuljeID:    pulje.ID,
				Threshold:  "urgent",
				ScheduleAt: urgentAt,
			})
		}
	}

	return broadcasts
}

func puljeScheduledBroadcastMessageID(broadcast puljeScheduledBroadcast) string {
	return fmt.Sprintf(
		"pulje-schedule:%s:%s:%s",
		broadcast.PuljeID,
		broadcast.Threshold,
		broadcast.ScheduleAt.UTC().Format(time.RFC3339),
	)
}

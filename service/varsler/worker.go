package varsler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
)

const (
	jobLeaseDuration = 30 * time.Second
	maxJobAttempts   = 5
)

type claimedJob struct {
	ID               int64
	PublicationID    int64
	PuljeID          string
	Revision         int
	BillettholderID  int64
	UserID           int64
	SubscriptionID   int64
	Attempts         int
	Payload          string
	Endpoint         string
	P256dh           string
	Auth             string
	PuljeStatus      string
	ProgramPublished bool
	LatestRevision   int
	Owned            bool
	Subscribed       bool
	LeaseToken       string
}

type claimedJobRef struct {
	ID         int64
	LeaseToken string
}

func (service *Service) Run(ctx context.Context) {
	if !service.config.Enabled || service.db == nil {
		<-ctx.Done()
		return
	}
	service.logger.Info("Web Push worker started")
	defer service.logger.Info("Web Push worker stopped")

	ticker := time.NewTicker(service.wake)
	defer ticker.Stop()
	for {
		processed, err := service.processNext(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			service.logger.Error(err.Error())
		}
		if processed {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (service *Service) processNext(ctx context.Context) (bool, error) {
	jobRef, err := service.claimNext(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	job, err := service.loadClaimedJob(ctx, jobRef)
	if errors.Is(err, sql.ErrNoRows) {
		return true, service.cancelJob(ctx, jobRef, "job dependencies no longer exist")
	}
	if err != nil {
		return true, err
	}
	if !job.isCurrent() || validateEndpoint(job.Endpoint) != nil {
		return true, service.cancelJob(ctx, jobRef, "job is stale or subscription is unavailable")
	}

	response, sendErr := webpush.SendNotificationWithContext(ctx, []byte(job.Payload), &webpush.Subscription{
		Endpoint: job.Endpoint,
		Keys: webpush.Keys{
			P256dh: job.P256dh,
			Auth:   job.Auth,
		},
	}, &webpush.Options{
		HTTPClient:      service.httpClient,
		Subscriber:      webPushSubscriber(service.config.Subject),
		VAPIDPublicKey:  service.config.PublicKey,
		VAPIDPrivateKey: service.config.PrivateKey,
		TTL:             3600,
	})
	if response != nil && response.Body != nil {
		defer response.Body.Close()
	}
	if sendErr != nil {
		return true, service.scheduleRetry(ctx, job, "Web Push transport failed")
	}
	if response.StatusCode >= http.StatusOK && response.StatusCode < http.StatusMultipleChoices {
		return true, service.markSent(ctx, job)
	}
	if response.StatusCode == http.StatusNotFound || response.StatusCode == http.StatusGone {
		return true, service.expireSubscription(ctx, job)
	}
	return true, service.scheduleRetry(ctx, job, fmt.Sprintf("push service returned HTTP %d", response.StatusCode))
}

func (service *Service) claimNext(ctx context.Context) (claimedJobRef, error) {
	now := service.now().UTC()
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return claimedJobRef{}, fmt.Errorf("begin Web Push job claim: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'exhausted', lease_until = NULL, lease_token = NULL, last_error = 'worker lease expired after maximum attempts'
		WHERE status = 'processing' AND attempts >= ? AND julianday(lease_until) <= julianday(?)
	`, maxJobAttempts, now.Format(time.RFC3339Nano)); err != nil {
		return claimedJobRef{}, fmt.Errorf("exhaust abandoned Web Push jobs: %w", err)
	}

	var jobID int64
	if err := tx.QueryRowContext(ctx, `
		SELECT id
		FROM web_push_jobs
		WHERE attempts < ? AND ((status IN ('pending','retrying') AND julianday(next_attempt_at) <= julianday(?))
		   OR (status = 'processing' AND julianday(lease_until) <= julianday(?)))
		ORDER BY next_attempt_at, id
		LIMIT 1
	`, maxJobAttempts, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)).Scan(&jobID); errors.Is(err, sql.ErrNoRows) {
		if err := tx.Commit(); err != nil {
			return claimedJobRef{}, fmt.Errorf("commit abandoned Web Push job cleanup: %w", err)
		}
		return claimedJobRef{}, sql.ErrNoRows
	} else if err != nil {
		return claimedJobRef{}, err
	}
	leaseToken := uuid.NewString()
	result, err := tx.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'processing', attempts = attempts + 1, lease_until = ?, lease_token = ?
		WHERE id = ? AND attempts < ? AND ((status IN ('pending','retrying') AND julianday(next_attempt_at) <= julianday(?)) OR (status = 'processing' AND julianday(lease_until) <= julianday(?)))
	`, now.Add(jobLeaseDuration).Format(time.RFC3339Nano), leaseToken, jobID, maxJobAttempts, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	if err != nil {
		return claimedJobRef{}, fmt.Errorf("claim Web Push job %d: %w", jobID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return claimedJobRef{}, fmt.Errorf("read Web Push job claim result: %w", err)
	}
	if rowsAffected != 1 {
		return claimedJobRef{}, sql.ErrNoRows
	}
	if err := tx.Commit(); err != nil {
		return claimedJobRef{}, fmt.Errorf("commit Web Push job claim: %w", err)
	}
	return claimedJobRef{ID: jobID, LeaseToken: leaseToken}, nil
}

func (service *Service) loadClaimedJob(ctx context.Context, jobRef claimedJobRef) (claimedJob, error) {
	var job claimedJob
	job.LeaseToken = jobRef.LeaseToken
	err := service.db.QueryRowContext(ctx, `
		SELECT job.id, job.publication_id, publication.pulje_id, publication.revision,
		       job.billettholder_id, job.user_id, job.subscription_id, job.attempts,
		       result.payload_json,
		       COALESCE(subscription.endpoint, ''), COALESCE(subscription.p256dh, ''), COALESCE(subscription.auth, ''),
		       pulje.status,
		       COALESCE((SELECT is_published FROM program_publishing_state WHERE id = 1), 0),
		       (SELECT MAX(revision) FROM pulje_varsel_publications WHERE pulje_id = publication.pulje_id),
		       EXISTS(SELECT 1 FROM relation_billettholdere_users relation WHERE relation.billettholder_id = job.billettholder_id AND relation.user_id = job.user_id),
		       subscription.id IS NOT NULL AND subscription.user_id = job.user_id
		FROM web_push_jobs job
		JOIN pulje_varsel_publications publication ON publication.id = job.publication_id
		JOIN pulje_varsel_results result ON result.publication_id = job.publication_id AND result.billettholder_id = job.billettholder_id
		JOIN puljer pulje ON pulje.id = publication.pulje_id
		LEFT JOIN web_push_subscriptions subscription ON subscription.id = job.subscription_id
		WHERE job.id = ? AND job.status = 'processing' AND job.lease_token = ?
	`, jobRef.ID, jobRef.LeaseToken).Scan(
		&job.ID, &job.PublicationID, &job.PuljeID, &job.Revision,
		&job.BillettholderID, &job.UserID, &job.SubscriptionID, &job.Attempts,
		&job.Payload, &job.Endpoint, &job.P256dh, &job.Auth, &job.PuljeStatus,
		&job.ProgramPublished, &job.LatestRevision, &job.Owned, &job.Subscribed,
	)
	if err != nil {
		return claimedJob{}, fmt.Errorf("load claimed Web Push job %d: %w", jobRef.ID, err)
	}
	return job, nil
}

func (job claimedJob) isCurrent() bool {
	return job.Revision == job.LatestRevision &&
		job.PuljeStatus == "Completed" &&
		job.ProgramPublished &&
		job.Owned &&
		job.Subscribed
}

func (service *Service) markSent(ctx context.Context, job claimedJob) error {
	if _, err := service.db.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'sent', sent_at = ?, lease_until = NULL, lease_token = NULL, last_error = ''
		WHERE id = ? AND status = 'processing' AND lease_token = ?
	`, service.now().UTC().Format(time.RFC3339Nano), job.ID, job.LeaseToken); err != nil {
		return fmt.Errorf("mark Web Push job %d sent: %w", job.ID, err)
	}
	return nil
}

func (service *Service) cancelJob(ctx context.Context, jobRef claimedJobRef, reason string) error {
	if _, err := service.db.ExecContext(ctx, `
		UPDATE web_push_jobs SET status = 'canceled', lease_until = NULL, lease_token = NULL, last_error = ?
		WHERE id = ? AND status = 'processing' AND lease_token = ?
	`, reason, jobRef.ID, jobRef.LeaseToken); err != nil {
		return fmt.Errorf("cancel Web Push job %d: %w", jobRef.ID, err)
	}
	return nil
}

func (service *Service) scheduleRetry(ctx context.Context, job claimedJob, reason string) error {
	if job.Attempts >= maxJobAttempts {
		result, err := service.db.ExecContext(ctx, `
			UPDATE web_push_jobs SET status = 'exhausted', lease_until = NULL, lease_token = NULL, last_error = ?
			WHERE id = ? AND status = 'processing' AND lease_token = ?
		`, reason, job.ID, job.LeaseToken)
		if err != nil {
			return fmt.Errorf("exhaust Web Push job %d: %w", job.ID, err)
		}
		if changed, err := oneRowChanged(result); err != nil {
			return fmt.Errorf("read exhausted Web Push job %d result: %w", job.ID, err)
		} else if changed {
			service.logger.Warn("Web Push retries exhausted", "pulje_id", job.PuljeID, "billettholder_id", job.BillettholderID)
		}
		return nil
	}
	delay := time.Minute << (job.Attempts - 1)
	if delay > time.Hour {
		delay = time.Hour
	}
	result, err := service.db.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'retrying', next_attempt_at = ?, lease_until = NULL, lease_token = NULL, last_error = ?
		WHERE id = ? AND status = 'processing' AND lease_token = ?
	`, service.now().UTC().Add(delay).Format(time.RFC3339Nano), reason, job.ID, job.LeaseToken)
	if err != nil {
		return fmt.Errorf("schedule Web Push job %d retry: %w", job.ID, err)
	}
	if changed, err := oneRowChanged(result); err != nil {
		return fmt.Errorf("read Web Push job %d retry result: %w", job.ID, err)
	} else if changed {
		service.logger.Debug("Web Push retry scheduled", "pulje_id", job.PuljeID, "billettholder_id", job.BillettholderID, "attempt", job.Attempts)
	}
	return nil
}

func oneRowChanged(result sql.Result) (bool, error) {
	rowsAffected, err := result.RowsAffected()
	return rowsAffected == 1, err
}

func (service *Service) expireSubscription(ctx context.Context, job claimedJob) error {
	tx, err := service.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin expired Web Push subscription cleanup: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	result, err := tx.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'canceled', lease_until = NULL, lease_token = NULL, last_error = 'push endpoint expired'
		WHERE id = ? AND status = 'processing' AND lease_token = ?
	`, job.ID, job.LeaseToken)
	if err != nil {
		return fmt.Errorf("fence expired Web Push subscription cleanup: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read expired Web Push subscription cleanup result: %w", err)
	}
	if rowsAffected == 0 {
		return nil
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE web_push_jobs
		SET status = 'canceled', lease_until = NULL, lease_token = NULL, last_error = 'push endpoint expired'
		WHERE subscription_id = ? AND status IN ('pending','retrying','processing')
	`, job.SubscriptionID); err != nil {
		return fmt.Errorf("cancel jobs for expired Web Push subscription %d: %w", job.SubscriptionID, err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM web_push_subscriptions WHERE id = ? AND user_id = ?`, job.SubscriptionID, job.UserID); err != nil {
		return fmt.Errorf("remove expired Web Push subscription %d: %w", job.SubscriptionID, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit expired Web Push subscription cleanup: %w", err)
	}
	return nil
}

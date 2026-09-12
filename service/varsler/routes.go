package varsler

import (
	"crypto/ecdh"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Regncon/conorganizer/service/userctx"
)

type subscriptionRequest struct {
	Endpoint       string          `json:"endpoint"`
	ExpirationTime json.RawMessage `json:"expirationTime"`
	Keys           struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

func (service *Service) getConfig(w http.ResponseWriter, r *http.Request) {
	if !authenticatedExternalID(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Enabled   bool   `json:"enabled"`
		PublicKey string `json:"publicKey"`
	}{Enabled: service.config.Enabled, PublicKey: service.config.PublicKey})
}

func (service *Service) getSubscription(w http.ResponseWriter, r *http.Request) {
	externalID := userctx.GetUserRequestInfo(r.Context()).Id
	if strings.TrimSpace(externalID) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !service.config.Enabled || service.db == nil {
		writeJSON(w, http.StatusOK, struct {
			Subscribed bool `json:"subscribed"`
		}{})
		return
	}
	endpoint := strings.TrimSpace(r.URL.Query().Get("endpoint"))
	if err := validateEndpoint(endpoint); err != nil {
		http.Error(w, "Invalid push endpoint", http.StatusBadRequest)
		return
	}
	var subscribed bool
	err := service.db.QueryRowContext(r.Context(), `
		SELECT EXISTS(
			SELECT 1 FROM web_push_subscriptions subscription
			JOIN users user ON user.id = subscription.user_id
			WHERE subscription.endpoint = ? AND user.external_id = ?
		)
	`, endpoint, externalID).Scan(&subscribed)
	if err != nil {
		service.logger.Error(fmt.Errorf("query current user's push subscription: %w", err).Error())
		http.Error(w, "Unable to read push subscription", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Subscribed bool `json:"subscribed"`
	}{Subscribed: subscribed})
}

func (service *Service) putSubscription(w http.ResponseWriter, r *http.Request) {
	externalID := userctx.GetUserRequestInfo(r.Context()).Id
	if strings.TrimSpace(externalID) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !service.config.Enabled || service.db == nil {
		http.Error(w, "Push notifications are unavailable", http.StatusServiceUnavailable)
		return
	}
	var request subscriptionRequest
	if err := decodeJSON(w, r, &request); err != nil {
		http.Error(w, "Invalid push subscription", http.StatusBadRequest)
		return
	}
	if err := validateEndpoint(request.Endpoint); err != nil {
		http.Error(w, "Invalid push endpoint", http.StatusBadRequest)
		return
	}
	if err := validateSubscriptionKeys(request.Keys.P256dh, request.Keys.Auth); err != nil {
		http.Error(w, "Invalid push subscription keys", http.StatusBadRequest)
		return
	}

	var userID int64
	if err := service.db.QueryRowContext(r.Context(), `SELECT id FROM users WHERE external_id = ?`, externalID).Scan(&userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		service.logger.Error(fmt.Errorf("find user for push subscription: %w", err).Error())
		http.Error(w, "Unable to save push subscription", http.StatusInternalServerError)
		return
	}

	var ownerID int64
	err := service.db.QueryRowContext(r.Context(), `SELECT user_id FROM web_push_subscriptions WHERE endpoint = ?`, request.Endpoint).Scan(&ownerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		service.logger.Error(fmt.Errorf("check push endpoint ownership: %w", err).Error(), "user_id", userID)
		http.Error(w, "Unable to save push subscription", http.StatusInternalServerError)
		return
	}
	if err == nil && ownerID != userID {
		http.Error(w, "Push endpoint belongs to another account", http.StatusConflict)
		return
	}
	result, err := service.db.ExecContext(r.Context(), `
		INSERT INTO web_push_subscriptions(user_id,endpoint,p256dh,auth)
		VALUES (?,?,?,?)
		ON CONFLICT(endpoint) DO UPDATE SET
			p256dh = excluded.p256dh,
			auth = excluded.auth,
			updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
		WHERE web_push_subscriptions.user_id = excluded.user_id
	`, userID, request.Endpoint, request.Keys.P256dh, request.Keys.Auth)
	if err != nil {
		service.logger.Error(fmt.Errorf("save push subscription: %w", err).Error(), "user_id", userID)
		http.Error(w, "Unable to save push subscription", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		service.logger.Error(fmt.Errorf("read saved push subscription result: %w", err).Error(), "user_id", userID)
		http.Error(w, "Unable to save push subscription", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Push endpoint belongs to another account", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (service *Service) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	externalID := userctx.GetUserRequestInfo(r.Context()).Id
	if strings.TrimSpace(externalID) == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if !service.config.Enabled || service.db == nil {
		http.Error(w, "Push notifications are unavailable", http.StatusServiceUnavailable)
		return
	}
	var request struct {
		Endpoint string `json:"endpoint"`
	}
	if err := decodeJSON(w, r, &request); err != nil || validateEndpoint(request.Endpoint) != nil {
		http.Error(w, "Invalid push endpoint", http.StatusBadRequest)
		return
	}
	tx, err := service.db.BeginTx(r.Context(), nil)
	if err != nil {
		service.logger.Error(fmt.Errorf("begin push subscription deletion: %w", err).Error())
		http.Error(w, "Unable to delete push subscription", http.StatusInternalServerError)
		return
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(r.Context(), `
		UPDATE web_push_jobs
		SET status = 'canceled', lease_until = NULL, lease_token = NULL, last_error = 'subscription removed'
		WHERE subscription_id = (
			SELECT subscription.id FROM web_push_subscriptions subscription
			JOIN users user ON user.id = subscription.user_id
			WHERE subscription.endpoint = ? AND user.external_id = ?
		) AND status IN ('pending','retrying','processing')
	`, request.Endpoint, externalID); err != nil {
		service.logger.Error(fmt.Errorf("cancel jobs for removed push subscription: %w", err).Error())
		http.Error(w, "Unable to delete push subscription", http.StatusInternalServerError)
		return
	}
	if _, err := tx.ExecContext(r.Context(), `
		DELETE FROM web_push_subscriptions
		WHERE endpoint = ? AND user_id = (SELECT id FROM users WHERE external_id = ?)
	`, request.Endpoint, externalID); err != nil {
		service.logger.Error(fmt.Errorf("delete current user's push subscription: %w", err).Error())
		http.Error(w, "Unable to delete push subscription", http.StatusInternalServerError)
		return
	}
	if err := tx.Commit(); err != nil {
		service.logger.Error(fmt.Errorf("commit push subscription deletion: %w", err).Error())
		http.Error(w, "Unable to delete push subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func authenticatedExternalID(r *http.Request) bool {
	return strings.TrimSpace(userctx.GetUserRequestInfo(r.Context()).Id) != ""
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain one JSON value")
	}
	return nil
}

func validateEndpoint(endpoint string) error {
	parsed, err := url.Parse(strings.TrimSpace(endpoint))
	if err != nil || parsed.Scheme != "https" || parsed.Hostname() == "" || parsed.User != nil {
		return errors.New("push endpoint must be an HTTPS URL")
	}
	if port := parsed.Port(); port != "" && port != "443" {
		return errors.New("push endpoint must use the default HTTPS port")
	}
	host := strings.ToLower(parsed.Hostname())
	allowed := host == "fcm.googleapis.com" ||
		host == "updates.push.services.mozilla.com" ||
		host == "push.services.mozilla.com" ||
		host == "web.push.apple.com" ||
		strings.HasSuffix(host, ".notify.windows.com")
	if !allowed {
		return errors.New("push endpoint host is not allowed")
	}
	return nil
}

func validateSubscriptionKeys(p256dh, auth string) error {
	publicBytes, err := decodeBase64URL(p256dh)
	if err != nil || len(publicBytes) != 65 || publicBytes[0] != 4 {
		return errors.New("invalid p256dh key")
	}
	if _, err := ecdh.P256().NewPublicKey(publicBytes); err != nil {
		return errors.New("invalid p256dh point")
	}
	authBytes, err := decodeBase64URL(auth)
	if err != nil || len(authBytes) < 16 {
		return errors.New("invalid auth secret")
	}
	return nil
}

func decodeBase64URL(value string) ([]byte, error) {
	if decoded, err := base64.RawURLEncoding.DecodeString(value); err == nil {
		return decoded, nil
	}
	return base64.URLEncoding.DecodeString(value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

package authctx

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5/middleware"
)

func TestLoggingSessionValidator_ValidationPreservesResultAndContext(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A validator returns a session token.",
		When:  "Session-validation timing is recorded.",
		Then:  "The result and context are unchanged, and the log includes correlation without token contents.",
	})

	// Given
	expectedToken := &descope.Token{ID: "user-123", JWT: "secret-session-token"}
	delegate := &fakeSessionValidator{sessionOK: true, sessionToken: expectedToken}
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "auth")
	validator := &loggingSessionValidator{validator: delegate, logger: logger}
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "request-123")

	// When
	valid, token, err := validator.ValidateSessionWithToken(ctx, expectedToken.JWT)

	// Then
	if err != nil || !valid || token != expectedToken {
		t.Fatalf("expected unchanged successful validation, got valid=%v token=%v error=%v", valid, token, err)
	}
	if delegate.sessionContext != ctx {
		t.Fatal("expected the original context to reach the validator")
	}
	assertSessionTimingLog(t, logs.Bytes(), "DEBUG", "session_validation", true)
	if strings.Contains(logs.String(), expectedToken.JWT) {
		t.Fatal("timing log contains the session token")
	}
}

func TestLoggingSessionValidator_RefreshPreservesFailureAndContext(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A refresh operation returns an error.",
		When:  "Refresh timing is recorded.",
		Then:  "The error and context are preserved without logging the token or underlying error again.",
	})

	// Given
	expectedError := errors.New("upstream refresh failure")
	delegate := &fakeSessionValidator{refreshErr: expectedError}
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug})).With("component", "auth")
	validator := &loggingSessionValidator{validator: delegate, logger: logger}
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "request-123")

	// When
	valid, token, err := validator.RefreshSessionWithToken(ctx, "secret-refresh-token")

	// Then
	if valid || token != nil || err != expectedError {
		t.Fatalf("expected unchanged refresh failure, got valid=%v token=%v error=%v", valid, token, err)
	}
	if delegate.refreshContext != ctx {
		t.Fatal("expected the original context to reach the validator")
	}
	assertSessionTimingLog(t, logs.Bytes(), "DEBUG", "session_refresh", false)
	if strings.Contains(logs.String(), "secret-refresh-token") || strings.Contains(logs.String(), expectedError.Error()) {
		t.Fatal("timing log contains the token or duplicates the underlying error")
	}
}

func TestLoggingSessionValidator_SlowOperationsAreVisibleAtDefaultLogLevel(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Authentication takes at least one second.",
		When:  "The operation duration is logged.",
		Then:  "A warning identifies the operation and its duration even when it succeeds.",
	})

	// Given
	expectedDuration := slowAuthenticationThreshold
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil)).With("component", "auth")
	validator := &loggingSessionValidator{logger: logger}
	ctx := context.WithValue(context.Background(), middleware.RequestIDKey, "request-123")

	// When
	validator.logDuration(ctx, "session_refresh", expectedDuration, true)

	// Then
	assertSessionTimingLog(t, logs.Bytes(), "WARN", "session_refresh", true)
	if !strings.Contains(logs.String(), `"duration_ms":1000`) {
		t.Fatalf("expected duration in milliseconds, got %s", logs.String())
	}
}

func TestLoggingSessionValidator_NormalOperationsAreQuietAtDefaultLogLevel(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Authentication completes in less than one second.",
		When:  "The operation duration is logged at the default log level.",
		Then:  "No routine timing log is emitted.",
	})

	// Given
	expectedLog := ""
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	validator := &loggingSessionValidator{logger: logger}

	// When
	validator.logDuration(context.Background(), "session_validation", slowAuthenticationThreshold-time.Nanosecond, true)

	// Then
	if logs.String() != expectedLog {
		t.Fatalf("expected no log, got %s", logs.String())
	}
}

func assertSessionTimingLog(t *testing.T, data []byte, level, operation string, succeeded bool) {
	t.Helper()
	var entry struct {
		Level      string `json:"level"`
		Component  string `json:"component"`
		Operation  string `json:"operation"`
		DurationMS *int64 `json:"duration_ms"`
		Succeeded  bool   `json:"succeeded"`
		RequestID  string `json:"request_id"`
	}
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("decode timing log: %v", err)
	}
	if entry.Level != level || entry.Component != "auth" || entry.Operation != operation || entry.Succeeded != succeeded || entry.RequestID != "request-123" {
		t.Fatalf("unexpected timing log: %s", data)
	}
	if entry.DurationMS == nil || *entry.DurationMS < 0 {
		t.Fatalf("expected non-negative duration_ms: %s", data)
	}
}

package authctx

import (
	"context"
	"log/slog"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5/middleware"
)

const slowAuthenticationThreshold = time.Second

type loggingSessionValidator struct {
	validator SessionValidator
	logger    *slog.Logger
}

func (v *loggingSessionValidator) ValidateSessionWithToken(ctx context.Context, sessionToken string) (bool, *descope.Token, error) {
	start := time.Now()
	valid, token, err := v.validator.ValidateSessionWithToken(ctx, sessionToken)
	v.logDuration(ctx, "session_validation", time.Since(start), err == nil && valid && token != nil)
	return valid, token, err
}

func (v *loggingSessionValidator) RefreshSessionWithToken(ctx context.Context, refreshToken string) (bool, *descope.Token, error) {
	start := time.Now()
	valid, token, err := v.validator.RefreshSessionWithToken(ctx, refreshToken)
	v.logDuration(ctx, "session_refresh", time.Since(start), err == nil && valid && token != nil)
	return valid, token, err
}

func (v *loggingSessionValidator) logDuration(ctx context.Context, operation string, duration time.Duration, succeeded bool) {
	level := slog.LevelDebug
	message := "authentication operation completed"
	if duration >= slowAuthenticationThreshold {
		level = slog.LevelWarn
		message = "slow authentication operation"
	}

	v.logger.Log(ctx, level, message,
		"operation", operation,
		"duration_ms", duration.Milliseconds(),
		"succeeded", succeeded,
		"request_id", middleware.GetReqID(ctx),
	)
}

package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func RequestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			responseWriter := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(responseWriter, r)

			statusCode := responseWriter.Status()
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			logLevel := requestLogLevel(statusCode)

			attributes := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status_code", statusCode),
				slog.Int64("duration_ms", time.Since(startTime).Milliseconds()),
			}

			requestID := middleware.GetReqID(r.Context())
			if requestID != "" {
				attributes = append(attributes, slog.String("request_id", requestID))
			}

			logger.LogAttrs(r.Context(), logLevel, "http request completed", attributes...)
		})
	}
}

func requestLogLevel(statusCode int) slog.Level {
	if statusCode >= http.StatusInternalServerError {
		return slog.LevelError
	}

	// These statuses are normal web control flow: logged-out users,
	// forbidden pages, stale links, missing assets, and crawlers. Keep
	// them visible without treating them as operational warnings.
	if statusCode == http.StatusUnauthorized ||
		statusCode == http.StatusForbidden ||
		statusCode == http.StatusNotFound {
		return slog.LevelInfo
	}

	if statusCode >= http.StatusBadRequest {
		return slog.LevelWarn
	}

	return slog.LevelInfo
}

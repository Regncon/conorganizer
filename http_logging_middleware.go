package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type statusCodeResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusCodeResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func RequestLoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			responseWriter := &statusCodeResponseWriter{ResponseWriter: w}

			next.ServeHTTP(responseWriter, r)

			statusCode := responseWriter.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			logLevel := slog.LevelInfo
			if statusCode >= http.StatusInternalServerError {
				logLevel = slog.LevelError
			} else if statusCode >= http.StatusBadRequest {
				logLevel = slog.LevelWarn
			}

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

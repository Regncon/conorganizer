package testutil

import (
	"context"
	"log/slog"
)

// ————————————————————————————
// 1. Define the abstraction used by production code
// ————————————————————————————
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
}

// ————————————————————————————
// 2. A lightweight stub that records calls
// ————————————————————————————
type StubLogger struct {
	calls []struct {
		msg           string
		keysAndValues []interface{}
	}
}

func (s *StubLogger) Info(msg string, keysAndValues ...interface{}) {
	s.calls = append(s.calls, struct {
		msg           string
		keysAndValues []interface{}
	}{msg, keysAndValues})
}

// 3. Adapter to make stubLogger compatible with *slog.Logger
// ————————————————————————————
type stubLoggerHandler struct {
	stub *StubLogger
}

func (h *stubLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *stubLoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract message
	msg := r.Message

	// Extract key-value pairs
	var keyValues []interface{}
	r.Attrs(func(attr slog.Attr) bool {
		keyValues = append(keyValues, attr.Key, attr.Value.Any())
		return true
	})

	// Forward to stub logger
	h.stub.Info(msg, keyValues...)
	return nil
}

func (h *stubLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *stubLoggerHandler) WithGroup(name string) slog.Handler {
	return h
}

func NewSlogAdapter(stub *StubLogger) *slog.Logger {
	handler := &stubLoggerHandler{stub: stub}
	return slog.New(handler)
}

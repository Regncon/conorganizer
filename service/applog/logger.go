package applog

import (
	"log/slog"
	"os"
	"strings"
)

const logLevelEnvName = "LOG_LEVEL"

func NewJSONLogger() *slog.Logger {
	logLevel := parseLogLevel(os.Getenv(logLevelEnvName))
	handlerOptions := &slog.HandlerOptions{Level: logLevel}
	return slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
}

func parseLogLevel(logLevelValue string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(logLevelValue)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	case "INFO", "":
		return slog.LevelInfo
	default:
		return slog.LevelInfo
	}
}

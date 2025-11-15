package logger

import (
	"log/slog"
	"os"
	"pr/internal/logctx"
)

const (
	envDebug = "debug"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Level
	switch env {
	case envDebug:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	}

	return slog.New(logctx.NewHandlerMiddleware(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	))
}

package loggers

import (
	"log/slog"
	"microservices-conversion/internal/configs"
	"os"
)

func NewAppLogger(appConfig configs.AppConfig) *AppLogger {
	var handler slog.Handler

	if appConfig.App.Prod {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return &AppLogger{Logger: slog.New(handler).With("app", appConfig.App.Name)}
}

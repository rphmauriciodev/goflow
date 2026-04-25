package platform

import (
	"log/slog"
	"os"
)

func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

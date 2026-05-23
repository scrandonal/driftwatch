package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/watcher"
	"github.com/yourorg/driftwatch/internal/webhook"
)

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to driftwatch config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	hookClient := webhook.New(cfg.WebhookURL)

	w, err := watcher.New(cfg.WatchPaths, cfg.PollInterval)
	if err != nil {
		slog.Error("failed to initialise watcher", "error", err)
		os.Exit(1)
	}

	slog.Info("driftwatch started",
		"paths", cfg.WatchPaths,
		"poll_interval", cfg.PollInterval,
		"webhook_url", cfg.WebhookURL,
	)

	for event := range w.Events() {
		slog.Info("drift detected", "path", event.Path, "old_hash", event.OldHash, "new_hash", event.NewHash)
		if err := hookClient.Send(event); err != nil {
			slog.Error("webhook delivery failed", "error", err, "path", event.Path)
		}
	}
}

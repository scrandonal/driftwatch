package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestLoad_Valid(t *testing.T) {
	path := writeTempConfig(t, `
watch_paths:
  - /etc/app/config.yaml
webhook_url: https://example.com/hook
poll_interval: 10s
log_level: debug
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.WebhookURL != "https://example.com/hook" {
		t.Errorf("expected webhook_url to be set")
	}
	if cfg.PollInterval != 10*time.Second {
		t.Errorf("expected poll_interval 10s, got %v", cfg.PollInterval)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTempConfig(t, `
watch_paths:
  - /tmp/file
webhook_url: https://example.com/hook
`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.PollInterval != 30*time.Second {
		t.Errorf("expected default poll_interval 30s, got %v", cfg.PollInterval)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default log_level info, got %s", cfg.LogLevel)
	}
}

func TestLoad_MissingWebhookURL(t *testing.T) {
	path := writeTempConfig(t, `
watch_paths:
  - /tmp/file
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing webhook_url")
	}
}

func TestLoad_MissingWatchPaths(t *testing.T) {
	path := writeTempConfig(t, `
webhook_url: https://example.com/hook
`)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing watch_paths")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

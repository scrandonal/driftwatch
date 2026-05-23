// Package config handles loading and validating driftwatch configuration.
package config

import (
	"errors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level application configuration.
type Config struct {
	WatchPaths  []string      `yaml:"watch_paths"`
	WebhookURL  string        `yaml:"webhook_url"`
	PollInterval time.Duration `yaml:"poll_interval"`
	LogLevel    string        `yaml:"log_level"`
}

// Load reads and parses a YAML config file at the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that required fields are present and sensible.
func (c *Config) Validate() error {
	if len(c.WatchPaths) == 0 {
		return errors.New("config: at least one watch_path is required")
	}
	if c.WebhookURL == "" {
		return errors.New("config: webhook_url is required")
	}
	if c.PollInterval <= 0 {
		c.PollInterval = 30 * time.Second
	}
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	return nil
}

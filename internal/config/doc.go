// Package config provides loading, parsing, and validation of driftwatch
// configuration files.
//
// Configuration is expressed as YAML and supports the following fields:
//
//	watch_paths:   ([]string, required) list of file paths to monitor
//	webhook_url:   (string, required)   URL to POST change notifications to
//	poll_interval: (duration, optional) how often to check files; default 30s
//	log_level:     (string, optional)   logging verbosity; default "info"
//
// Valid log_level values are: "debug", "info", "warn", "error".
//
// The poll_interval field accepts Go duration strings such as "10s", "1m",
// or "2m30s". Values below 1s are rejected to prevent excessive polling.
//
// Example:
//
//	watch_paths:
//	  - /etc/myapp/settings.yaml
//	  - /etc/myapp/secrets.env
//	webhook_url: https://hooks.example.com/driftwatch
//	poll_interval: 15s
//	log_level: debug
package config

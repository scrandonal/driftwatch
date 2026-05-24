package supervisor

import "time"

// Config controls the restart behaviour of a [Supervisor].
type Config struct {
	// MaxRestarts is the maximum number of times the supervised function may
	// be restarted before the supervisor gives up and returns an error.
	// A value of zero means restart indefinitely.
	MaxRestarts int

	// RestartDelay is the duration the supervisor waits between a worker
	// exit and the next restart attempt. Defaults to one second when zero.
	RestartDelay time.Duration
}

// defaults fills in zero values with sensible production defaults.
func (c *Config) defaults() {
	if c.RestartDelay == 0 {
		c.RestartDelay = time.Second
	}
}

// DefaultConfig returns a Config with conservative defaults suitable for
// production use: unlimited restarts with a one-second back-off.
func DefaultConfig() Config {
	return Config{
		MaxRestarts:  0,
		RestartDelay: time.Second,
	}
}

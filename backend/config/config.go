// Package config loads runtime configuration from the environment.
// Values are never logged; validation errors report only missing or
// invalid variable names, per docs/configuration.md.
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultSessionIdleTTL      = 12 * time.Hour
	defaultSessionAbsoluteTTL  = 168 * time.Hour
	defaultSessionCookieSecure = true
)

type Config struct {
	BootstrapToken      string
	MasterKey           string
	DatabaseDSN         string
	SessionIdleTTL      time.Duration
	SessionAbsoluteTTL  time.Duration
	SessionCookieSecure bool
}

var (
	ErrMissingConfig = errors.New("missing required configuration")
	ErrInvalidConfig = errors.New("invalid configuration value")
)

func Load() (Config, error) {
	cfg := Config{
		BootstrapToken: os.Getenv("AKRITAS_BOOTSTRAP_TOKEN"),
		MasterKey:      os.Getenv("AKRITAS_MASTER_KEY"),
		DatabaseDSN:    os.Getenv("AKRITAS_DB_DSN"),
	}

	var missing []string
	if cfg.BootstrapToken == "" {
		missing = append(missing, "AKRITAS_BOOTSTRAP_TOKEN")
	}
	if cfg.MasterKey == "" {
		missing = append(missing, "AKRITAS_MASTER_KEY")
	}
	if cfg.DatabaseDSN == "" {
		missing = append(missing, "AKRITAS_DB_DSN")
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("%w: %v", ErrMissingConfig, missing)
	}

	idleTTL, err := parseDurationOrDefault("AKRITAS_SESSION_IDLE_TTL", defaultSessionIdleTTL)
	if err != nil {
		return Config{}, err
	}
	cfg.SessionIdleTTL = idleTTL

	absoluteTTL, err := parseDurationOrDefault("AKRITAS_SESSION_ABSOLUTE_TTL", defaultSessionAbsoluteTTL)
	if err != nil {
		return Config{}, err
	}
	cfg.SessionAbsoluteTTL = absoluteTTL

	cookieSecure, err := parseBoolOrDefault("AKRITAS_SESSION_COOKIE_SECURE", defaultSessionCookieSecure)
	if err != nil {
		return Config{}, err
	}
	cfg.SessionCookieSecure = cookieSecure

	return cfg, nil
}

func parseDurationOrDefault(name string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidConfig, name)
	}
	return parsed, nil
}

func parseBoolOrDefault(name string, fallback bool) (bool, error) {
	value := os.Getenv(name)
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%w: %s", ErrInvalidConfig, name)
	}
	return parsed, nil
}

// Package config loads runtime configuration from the environment.
// Values are never logged; validation errors report only missing
// variable names, per docs/configuration.md.
package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	BootstrapToken string
	MasterKey      string
	DatabaseDSN    string
}

var ErrMissingConfig = errors.New("missing required configuration")

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

	return cfg, nil
}

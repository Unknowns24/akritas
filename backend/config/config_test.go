package config

import (
	"errors"
	"strings"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	t.Setenv(key, value)
}

func TestLoadSucceedsWithAllVariablesSet(t *testing.T) {
	setEnv(t, "AKRITAS_BOOTSTRAP_TOKEN", "deployment-bootstrap-secret-not-a-totp-seed")
	setEnv(t, "AKRITAS_MASTER_KEY", "a-high-entropy-master-key")
	setEnv(t, "AKRITAS_DB_DSN", "postgres://localhost:5432/akritas?sslmode=disable")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.BootstrapToken == "" || cfg.MasterKey == "" || cfg.DatabaseDSN == "" {
		t.Fatalf("unexpected empty config: %+v", cfg)
	}
}

func TestLoadFailsFastWithoutLeakingValues(t *testing.T) {
	setEnv(t, "AKRITAS_BOOTSTRAP_TOKEN", "")
	setEnv(t, "AKRITAS_MASTER_KEY", "a-high-entropy-master-key-that-must-not-appear-in-errors")
	setEnv(t, "AKRITAS_DB_DSN", "")

	_, err := Load()
	if !errors.Is(err, ErrMissingConfig) {
		t.Fatalf("expected ErrMissingConfig, got %v", err)
	}
	if !strings.Contains(err.Error(), "AKRITAS_BOOTSTRAP_TOKEN") || !strings.Contains(err.Error(), "AKRITAS_DB_DSN") {
		t.Fatalf("error must name the missing variables, got %q", err.Error())
	}
	if strings.Contains(err.Error(), "a-high-entropy-master-key-that-must-not-appear-in-errors") {
		t.Fatal("error must never contain a configured value")
	}
}

package config

import (
	"errors"
	"strings"
	"testing"
	"time"
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

func TestLoadUsesSessionDefaultsWhenUnset(t *testing.T) {
	setEnv(t, "AKRITAS_BOOTSTRAP_TOKEN", "deployment-bootstrap-secret-not-a-totp-seed")
	setEnv(t, "AKRITAS_MASTER_KEY", "a-high-entropy-master-key")
	setEnv(t, "AKRITAS_DB_DSN", "postgres://localhost:5432/akritas?sslmode=disable")
	setEnv(t, "AKRITAS_SESSION_IDLE_TTL", "")
	setEnv(t, "AKRITAS_SESSION_ABSOLUTE_TTL", "")
	setEnv(t, "AKRITAS_SESSION_COOKIE_SECURE", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SessionIdleTTL != defaultSessionIdleTTL {
		t.Fatalf("SessionIdleTTL = %v, want default %v", cfg.SessionIdleTTL, defaultSessionIdleTTL)
	}
	if cfg.SessionAbsoluteTTL != defaultSessionAbsoluteTTL {
		t.Fatalf("SessionAbsoluteTTL = %v, want default %v", cfg.SessionAbsoluteTTL, defaultSessionAbsoluteTTL)
	}
	if cfg.SessionCookieSecure != defaultSessionCookieSecure {
		t.Fatalf("SessionCookieSecure = %v, want default %v", cfg.SessionCookieSecure, defaultSessionCookieSecure)
	}
}

func TestLoadParsesExplicitSessionValues(t *testing.T) {
	setEnv(t, "AKRITAS_BOOTSTRAP_TOKEN", "deployment-bootstrap-secret-not-a-totp-seed")
	setEnv(t, "AKRITAS_MASTER_KEY", "a-high-entropy-master-key")
	setEnv(t, "AKRITAS_DB_DSN", "postgres://localhost:5432/akritas?sslmode=disable")
	setEnv(t, "AKRITAS_SESSION_IDLE_TTL", "1h")
	setEnv(t, "AKRITAS_SESSION_ABSOLUTE_TTL", "24h")
	setEnv(t, "AKRITAS_SESSION_COOKIE_SECURE", "false")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.SessionIdleTTL != time.Hour {
		t.Fatalf("SessionIdleTTL = %v, want 1h", cfg.SessionIdleTTL)
	}
	if cfg.SessionAbsoluteTTL != 24*time.Hour {
		t.Fatalf("SessionAbsoluteTTL = %v, want 24h", cfg.SessionAbsoluteTTL)
	}
	if cfg.SessionCookieSecure {
		t.Fatal("SessionCookieSecure = true, want false")
	}
}

func TestLoadRejectsInvalidSessionValuesWithoutLeaking(t *testing.T) {
	setEnv(t, "AKRITAS_BOOTSTRAP_TOKEN", "deployment-bootstrap-secret-not-a-totp-seed")
	setEnv(t, "AKRITAS_MASTER_KEY", "a-high-entropy-master-key")
	setEnv(t, "AKRITAS_DB_DSN", "postgres://localhost:5432/akritas?sslmode=disable")
	setEnv(t, "AKRITAS_SESSION_IDLE_TTL", "not-a-duration")

	_, err := Load()
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
	if !strings.Contains(err.Error(), "AKRITAS_SESSION_IDLE_TTL") {
		t.Fatalf("error must name the invalid variable, got %q", err.Error())
	}
	if strings.Contains(err.Error(), "not-a-duration") {
		t.Fatal("error must never contain the configured value")
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

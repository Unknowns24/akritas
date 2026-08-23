package config

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestLoadFromViperAppliesDefaultsAndValidatesSecurityValues(t *testing.T) {
	t.Parallel()

	v := viper.New()
	v.Set("AKRITAS_DATABASE_URL", "postgres://user:password@db:5432/akritas")
	v.Set("AKRITAS_PUBLIC_URL", "https://akritas.example.com")
	v.Set("AKRITAS_MASTER_KEY", base64.StdEncoding.EncodeToString(make([]byte, 32)))
	v.Set("AKRITAS_PAGINATION_SECRET", "01234567890123456789012345678901")
	v.Set("AKRITAS_BOOTSTRAP_TOKEN", "01234567890123456789012345678901")

	configuration, err := loadFromViper(v)
	if err != nil {
		t.Fatalf("loadFromViper() error = %v", err)
	}
	if configuration.PaginationTTL != 15*time.Minute {
		t.Fatalf("PaginationTTL = %s, want 15m", configuration.PaginationTTL)
	}
	if configuration.DatabaseMaxOpenConnections != 10 || configuration.DatabaseMaxIdleConnections != 5 {
		t.Fatalf("unexpected pool defaults: %+v", configuration)
	}
	if configuration.SessionIdleTTL != 12*time.Hour || configuration.SessionAbsoluteTTL != 7*24*time.Hour || !configuration.SessionCookieSecure {
		t.Fatalf("unexpected session defaults: %+v", configuration)
	}
	if configuration.AuthRateLimitAttempts != 5 || configuration.AuthRateLimitWindow != 15*time.Minute || configuration.AuthRateLimitMaxKeys != 4096 {
		t.Fatalf("unexpected auth rate-limit defaults: %+v", configuration)
	}
	if configuration.MonitoringPollInterval != 10*time.Second || configuration.MonitoringConcurrency != 4 {
		t.Fatalf("unexpected monitoring defaults: %+v", configuration)
	}
	if len(configuration.AllowedOrigins) != 1 || configuration.AllowedOrigins[0] != configuration.PublicURL {
		t.Fatalf("public URL must be an allowed origin: %v", configuration.AllowedOrigins)
	}
	if configuration.MasterKeyEncoded != "" || configuration.PaginationSecretValue != "" || configuration.BootstrapTokenValue != "" || configuration.AllowedOriginsValue != "" {
		t.Fatal("raw secret and origin values must be cleared")
	}
}

func TestLoadFromViperRejectsUnsafeSessionAndOrigins(t *testing.T) {
	t.Parallel()
	base := func() *viper.Viper {
		v := viper.New()
		v.Set("AKRITAS_DATABASE_URL", "postgres://user:password@db:5432/akritas")
		v.Set("AKRITAS_PUBLIC_URL", "https://akritas.example.com")
		v.Set("AKRITAS_MASTER_KEY", base64.StdEncoding.EncodeToString(make([]byte, 32)))
		v.Set("AKRITAS_PAGINATION_SECRET", "01234567890123456789012345678901")
		v.Set("AKRITAS_BOOTSTRAP_TOKEN", "01234567890123456789012345678901")
		return v
	}
	cases := map[string]func(*viper.Viper){
		"insecure cookie":            func(v *viper.Viper) { v.Set("AKRITAS_SESSION_COOKIE_SECURE", false) },
		"idle exceeds max":           func(v *viper.Viper) { v.Set("AKRITAS_SESSION_IDLE_TTL", 8*24*time.Hour) },
		"wildcard origin":            func(v *viper.Viper) { v.Set("AKRITAS_ALLOWED_ORIGINS", "*") },
		"origin with path":           func(v *viper.Viper) { v.Set("AKRITAS_ALLOWED_ORIGINS", "https://app.example.com/path") },
		"short bootstrap key":        func(v *viper.Viper) { v.Set("AKRITAS_BOOTSTRAP_TOKEN", "short") },
		"zero monitor interval":      func(v *viper.Viper) { v.Set("AKRITAS_MONITORING_POLL_INTERVAL", 0) },
		"zero monitor concurrency":   func(v *viper.Viper) { v.Set("AKRITAS_MONITORING_CONCURRENCY", 0) },
		"excess monitor concurrency": func(v *viper.Viper) { v.Set("AKRITAS_MONITORING_CONCURRENCY", 5) },
		"zero auth attempts":         func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_ATTEMPTS", 0) },
		"excess auth attempts":       func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_ATTEMPTS", 101) },
		"zero auth window":           func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_WINDOW", 0) },
		"excess auth window":         func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_WINDOW", 25*time.Hour) },
		"zero auth keys":             func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_MAX_KEYS", 0) },
		"excess auth keys":           func(v *viper.Viper) { v.Set("AKRITAS_AUTH_RATE_LIMIT_MAX_KEYS", 100001) },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			v := base()
			mutate(v)
			if _, err := loadFromViper(v); err == nil {
				t.Fatal("unsafe configuration accepted")
			}
		})
	}
}

func TestConfigurationDoesNotAcceptLegacyDatabaseDSN(t *testing.T) {
	t.Parallel()
	for _, key := range configurationKeys() {
		if key == "AKRITAS_DB_DSN" {
			t.Fatal("legacy AKRITAS_DB_DSN must not be part of centralized configuration")
		}
	}
}

func TestViperEnvironmentOverridesOptionalConfigFile(t *testing.T) {
	t.Setenv("AKRITAS_PUBLIC_URL", "https://environment.example.com")
	t.Setenv("AKRITAS_AUTH_RATE_LIMIT_ATTEMPTS", "7")
	t.Setenv("AKRITAS_AUTH_RATE_LIMIT_WINDOW", "5m")
	t.Setenv("AKRITAS_AUTH_RATE_LIMIT_MAX_KEYS", "512")

	v := viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()
	bindEnvironment(v)
	file := strings.NewReader(strings.Join([]string{
		"AKRITAS_DATABASE_URL=postgres://user:password@db:5432/akritas",
		"AKRITAS_PUBLIC_URL=https://file.example.com",
		"AKRITAS_MASTER_KEY=" + base64.StdEncoding.EncodeToString(make([]byte, 32)),
		"AKRITAS_PAGINATION_SECRET=01234567890123456789012345678901",
		"AKRITAS_BOOTSTRAP_TOKEN=01234567890123456789012345678901",
	}, "\n"))
	if err := v.ReadConfig(file); err != nil {
		t.Fatal(err)
	}

	configuration, err := loadFromViper(v)
	if err != nil {
		t.Fatalf("loadFromViper() error = %v", err)
	}
	if configuration.PublicURL != "https://environment.example.com" {
		t.Fatalf("PublicURL = %q, environment must override app.env", configuration.PublicURL)
	}
	if configuration.AuthRateLimitAttempts != 7 || configuration.AuthRateLimitWindow != 5*time.Minute || configuration.AuthRateLimitMaxKeys != 512 {
		t.Fatalf("rate-limit environment overrides not applied: %+v", configuration)
	}
}

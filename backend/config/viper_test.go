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
}

func TestViperEnvironmentOverridesOptionalConfigFile(t *testing.T) {
	t.Setenv("AKRITAS_PUBLIC_URL", "https://environment.example.com")

	v := viper.New()
	v.SetConfigType("env")
	v.AutomaticEnv()
	bindEnvironment(v)
	file := strings.NewReader(strings.Join([]string{
		"AKRITAS_DATABASE_URL=postgres://user:password@db:5432/akritas",
		"AKRITAS_PUBLIC_URL=https://file.example.com",
		"AKRITAS_MASTER_KEY=" + base64.StdEncoding.EncodeToString(make([]byte, 32)),
		"AKRITAS_PAGINATION_SECRET=01234567890123456789012345678901",
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
}

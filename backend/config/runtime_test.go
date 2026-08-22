package config

import (
	"encoding/base64"
	"errors"
	"testing"
)

func TestLoadValuesFailsClosedForMissingOrInvalidSecurityConfiguration(t *testing.T) {
	values := map[string]string{
		"AKRITAS_DATABASE_URL":      "postgres://user:password@db:5432/akritas",
		"AKRITAS_PUBLIC_URL":        "https://akritas.example.com",
		"AKRITAS_MASTER_KEY":        base64.StdEncoding.EncodeToString(make([]byte, 32)),
		"AKRITAS_PAGINATION_SECRET": "01234567890123456789012345678901",
	}
	getenv := func(name string) string { return values[name] }
	configuration, err := LoadValues(getenv)
	if err != nil || len(configuration.MasterKey) != 32 {
		t.Fatalf("valid configuration rejected: %v", err)
	}
	values["AKRITAS_MASTER_KEY"] = "not-a-key"
	if _, err := LoadValues(getenv); !errors.Is(err, ErrInvalidMasterKey) {
		t.Fatalf("invalid master key must stop startup: %v", err)
	}
	values["AKRITAS_MASTER_KEY"] = base64.StdEncoding.EncodeToString(make([]byte, 32))
	values["AKRITAS_PUBLIC_URL"] = "http://public.example.com"
	if _, err := LoadValues(getenv); !errors.Is(err, ErrInvalidRuntimeConfiguration) {
		t.Fatalf("non-HTTPS public URL must stop startup: %v", err)
	}
}

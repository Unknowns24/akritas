package config

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestParseMasterKey(t *testing.T) {
	t.Parallel()

	raw := []byte("0123456789abcdef0123456789abcdef")
	encoded := base64.StdEncoding.EncodeToString(raw)
	got, err := ParseMasterKey(encoded)
	if err != nil {
		t.Fatalf("ParseMasterKey() error = %v", err)
	}
	if string(got) != string(raw) {
		t.Fatal("ParseMasterKey() returned a different key")
	}

	for _, value := range []string{"", "not-base64", base64.StdEncoding.EncodeToString([]byte("too-short"))} {
		value := value
		t.Run(value, func(t *testing.T) {
			t.Parallel()
			_, err := ParseMasterKey(value)
			if err == nil {
				t.Fatal("ParseMasterKey() expected error")
			}
			if strings.Contains(err.Error(), value) && value != "" {
				t.Fatal("ParseMasterKey() leaked configured value")
			}
		})
	}
}

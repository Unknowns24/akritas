package utils

import (
	"testing"
	"time"
)

func TestParseAndFormatISODuration(t *testing.T) {
	parsed, err := ParseISODuration("PT30M")
	if err != nil || parsed != 30*time.Minute {
		t.Fatalf("parse PT30M: %v %s", err, parsed)
	}
	if got := FormatISODuration(30 * time.Minute); got != "PT30M" {
		t.Fatalf("format: %s", got)
	}
	if _, err := ParseISODuration("PT"); err == nil {
		t.Fatal("expected invalid duration")
	}
	if _, err := ParseISODuration("nope"); err == nil {
		t.Fatal("expected invalid duration")
	}
}

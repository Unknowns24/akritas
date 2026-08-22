package security

import (
	"testing"
	"time"
)

func TestClockNow(t *testing.T) {
	t.Parallel()

	now := NewClock().Now()
	if now.IsZero() {
		t.Fatal("Now() must not be zero")
	}
	if now.Location() != time.UTC {
		t.Fatalf("Now() must be in UTC, got %v", now.Location())
	}
}

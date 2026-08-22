package security

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiterAllowsUpToLimitThenDenies(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(5, 15*time.Minute, func() time.Time { return now })

	for i := 0; i < 5; i++ {
		allowed, err := limiter.Allow(context.Background(), "203.0.113.10")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("attempt %d should be allowed", i+1)
		}
	}

	allowed, err := limiter.Allow(context.Background(), "203.0.113.10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatal("attempt beyond the limit must be denied")
	}
}

func TestRateLimiterDoesNotMixKeys(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(1, 15*time.Minute, func() time.Time { return now })

	if allowed, err := limiter.Allow(context.Background(), "203.0.113.10"); err != nil || !allowed {
		t.Fatalf("first key should be allowed, got allowed=%v err=%v", allowed, err)
	}
	if allowed, err := limiter.Allow(context.Background(), "203.0.113.10"); err != nil || allowed {
		t.Fatalf("first key should be denied on second attempt, got allowed=%v err=%v", allowed, err)
	}
	if allowed, err := limiter.Allow(context.Background(), "198.51.100.20"); err != nil || !allowed {
		t.Fatalf("a different key must not be affected by another key's limit, got allowed=%v err=%v", allowed, err)
	}
}

func TestRateLimiterResetsAfterWindowElapses(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(1, 15*time.Minute, func() time.Time { return now })

	if allowed, _ := limiter.Allow(context.Background(), "203.0.113.10"); !allowed {
		t.Fatal("first attempt should be allowed")
	}
	if allowed, _ := limiter.Allow(context.Background(), "203.0.113.10"); allowed {
		t.Fatal("second attempt within the window should be denied")
	}

	now = now.Add(15 * time.Minute)
	if allowed, _ := limiter.Allow(context.Background(), "203.0.113.10"); !allowed {
		t.Fatal("attempt exactly at window boundary should be allowed again")
	}
}

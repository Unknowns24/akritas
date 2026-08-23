package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiterAllowsUpToLimitThenDenies(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(5, 15*time.Minute, 128, func() time.Time { return now })

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
	limiter := newRateLimiterWithClock(1, 15*time.Minute, 128, func() time.Time { return now })

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
	limiter := newRateLimiterWithClock(1, 15*time.Minute, 128, func() time.Time { return now })

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

func TestRateLimiterCleansExpiredKeys(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(1, time.Minute, cleanupEvery, func() time.Time { return now }).(*windowRateLimiter)
	for index := 0; index < cleanupEvery; index++ {
		if allowed, err := limiter.Allow(context.Background(), string(rune(index+1))); err != nil || !allowed {
			t.Fatalf("seed key %d: allowed=%v err=%v", index, allowed, err)
		}
	}
	now = now.Add(2 * time.Minute)
	if allowed, err := limiter.Allow(context.Background(), "fresh"); err != nil || !allowed {
		t.Fatalf("fresh key: allowed=%v err=%v", allowed, err)
	}
	if len(limiter.entries) != 1 {
		t.Fatalf("expired limiter entries were not cleaned: %d remain", len(limiter.entries))
	}
}

func TestRateLimiterFailsClosedAtKeyCapacityWithoutGrowing(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	limiter := newRateLimiterWithClock(2, time.Minute, 2, func() time.Time { return now }).(*windowRateLimiter)

	for _, key := range []string{"first", "second"} {
		if allowed, err := limiter.Allow(context.Background(), key); err != nil || !allowed {
			t.Fatalf("seed %q: allowed=%v err=%v", key, allowed, err)
		}
	}
	if allowed, err := limiter.Allow(context.Background(), "third"); err != nil || allowed {
		t.Fatalf("new key at capacity: allowed=%v err=%v", allowed, err)
	}
	if len(limiter.entries) != 2 {
		t.Fatalf("limiter grew beyond cap: %d", len(limiter.entries))
	}
	if allowed, err := limiter.Allow(context.Background(), "first"); err != nil || !allowed {
		t.Fatalf("existing bucket should retain its budget: allowed=%v err=%v", allowed, err)
	}

	now = now.Add(time.Minute)
	if allowed, err := limiter.Allow(context.Background(), "third"); err != nil || !allowed {
		t.Fatalf("expired buckets should free capacity: allowed=%v err=%v", allowed, err)
	}
	if len(limiter.entries) != 1 {
		t.Fatalf("expired buckets not removed: %d", len(limiter.entries))
	}
}

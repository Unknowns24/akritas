package security

import (
	"context"
	"sync"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type rateLimiterEntry struct {
	windowStart time.Time
	attempts    int
}

// windowRateLimiter is a fixed-window, in-memory limiter. It does not survive
// restarts and does not coordinate across instances — sufficient for the MVP
// single-admin deployment; advanced rate limiting is out of scope (PB-065).
type windowRateLimiter struct {
	maxAttempts int
	window      time.Duration
	now         func() time.Time

	mu      sync.Mutex
	entries map[string]*rateLimiterEntry
}

func NewRateLimiter(maxAttempts int, window time.Duration) out.RateLimiter {
	return newRateLimiterWithClock(maxAttempts, window, time.Now)
}

func newRateLimiterWithClock(maxAttempts int, window time.Duration, now func() time.Time) out.RateLimiter {
	return &windowRateLimiter{
		maxAttempts: maxAttempts,
		window:      window,
		now:         now,
		entries:     make(map[string]*rateLimiterEntry),
	}
}

func (l *windowRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	entry, exists := l.entries[key]
	if !exists || now.Sub(entry.windowStart) >= l.window {
		l.entries[key] = &rateLimiterEntry{windowStart: now, attempts: 1}
		return true, nil
	}

	if entry.attempts >= l.maxAttempts {
		return false, nil
	}
	entry.attempts++
	return true, nil
}

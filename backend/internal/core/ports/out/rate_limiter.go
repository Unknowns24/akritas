package out

import "context"

// RateLimiter reports whether a request identified by key is allowed under
// the configured attempt/window budget for its caller.
type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

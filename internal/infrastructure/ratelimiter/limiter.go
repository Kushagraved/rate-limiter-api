package ratelimiter

import (
	"context"
	"time"
)

// RateLimit holds the policy for a single rate-limit check.
type RateLimit struct {
	MaxRequests int
	Window      time.Duration
}

// RateLimiter is the strategy interface. Config is passed per-call so the
// same instance can serve multiple routes with different policies.
type RateLimiter interface {
	Allow(ctx context.Context, key string, rate RateLimit) (allowed bool, remaining int, resetAt time.Time, err error)
}

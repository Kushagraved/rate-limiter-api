package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisLimiter implements RateLimiter using a Redis sorted-set sliding window.
type redisLimiter struct {
	client *redis.Client
}

func NewRedisLimiter(client *redis.Client) RateLimiter {
	return &redisLimiter{client: client}
}

// redisScript atomically:
//  1. Removes timestamps outside the window
//  2. Counts remaining
//  3. Adds current timestamp if under limit
//  4. Sets TTL on the key
//
// Returns: {allowed (0|1), count_after, oldest_ms}
var redisScript = redis.NewScript(`
local key          = KEYS[1]
local now          = tonumber(ARGV[1])
local window_ms    = tonumber(ARGV[2])
local max          = tonumber(ARGV[3])
local window_start = now - window_ms

redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)
local count = redis.call('ZCARD', key)

if count < max then
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window_ms)
    return {1, count + 1, 0}
end

local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')
return {0, count, tonumber(oldest[2])}
`)

func (r *redisLimiter) Allow(ctx context.Context, key string, rate RateLimit) (bool, int, time.Time, error) {
	redisKey := fmt.Sprintf("rl:%s", key)
	now := time.Now()
	windowMs := rate.Window.Milliseconds()

	res, err := redisScript.Run(ctx, r.client, []string{redisKey},
		now.UnixMilli(), windowMs, rate.MaxRequests,
	).Int64Slice()
	if err != nil {
		return false, 0, time.Time{}, fmt.Errorf("redis rate limiter: %w", err)
	}

	allowed := res[0] == 1
	count := int(res[1])
	oldestMs := res[2]

	remaining := rate.MaxRequests - count
	if remaining < 0 {
		remaining = 0
	}

	if allowed {
		return true, remaining, time.Time{}, nil
	}

	resetAt := time.UnixMilli(oldestMs).Add(rate.Window)
	return false, 0, resetAt, nil
}

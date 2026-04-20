# Rate-Limited API Service

A Go service implementing per-user rate limiting backed by Redis using a sliding window algorithm.

---

## API

### `POST /v1/request`

Rate-limited per user. Default: 5 requests per user per 60-second sliding window.

**Request**
```json
{
  "user_id": "user_123",
  "payload": { "any": "data" }
}
```

**Response `200 OK`**
```json
{
  "success": true
}
```

**Response `429 Too Many Requests`**
```json
{
  "error": "TOO_MANY_REQUESTS",
  "message": "rate limit exceeded"
}
```

**Rate-limit headers** (set on every response):
```
X-RateLimit-Limit:     5
X-RateLimit-Remaining: 3
X-RateLimit-Reset:     1713600720   (unix timestamp — when the current window expires)
Retry-After:           1713600720   (unix timestamp, only on 429)
```

### `GET /v1/stats?user_id=user_123`

Returns per-user request statistics. Requires `user_id` query param.

**Response `200 OK`**
```json
{
  "user_id": "user_123",
  "total_count": 12,
  "window_count": 5,
  "last_request": 1713600715,
  "retry_after": 1713600775
}
```

- `total_count` — all-time requests for this user
- `window_count` — requests within the current sliding window
- `last_request` — unix timestamp of the most recent request
- `retry_after` — unix timestamp of when the user can make a new request (only present when at or over the limit)

### `GET /health`

Health check endpoint.

---

## Running

### Docker Compose

```bash
docker compose up app redis
```

This starts the API on `http://localhost:8080` backed by Redis.

### Configuration

All config lives in `config/dev/config.yaml`:

```yaml
settings:
  rate_limiter:
    max_requests: 5
    window_secs: 60
```

Rate limit values are passed inline at the route level — each route constructs its own `ratelimiter.RateLimit` from config, so different routes can have independent limits.

---

## Architecture

```
cmd/server/main.go          → entrypoint, wiring
router/                     → route registration, rate limit config per route
internal/middleware/         → rate limit middleware (body peek, headers, 429)
internal/infrastructure/
  ratelimiter/              → RateLimiter interface + Redis implementation
  redis/                    → Redis client setup
internal/controller/        → HTTP handlers
internal/core/request/      → business logic
internal/repository/        → in-memory request storage
internal/config/            → YAML config loading
```

### Rate Limiter Interface

```go
type RateLimiter interface {
    Allow(ctx context.Context, key string, rate RateLimit) (allowed bool, remaining int, resetAt time.Time, err error)
}
```

The `RateLimit` policy (max requests, window) is passed per-call — the limiter instance is stateless with respect to config. A single limiter serves all routes, each with its own policy.

### Sliding Window

Uses a Redis Lua script with sorted sets (`ZADD`/`ZREMRANGEBYSCORE`) for atomic, multi-instance-safe sliding window enforcement.

---

## Things to Improve

- **Strategy pattern for algorithms** — the `RateLimiter` interface already supports this. Add token bucket or leaky bucket implementations behind the same interface for burst-friendly use cases.
- **PostgreSQL for persistence** — swap the in-memory request repository for Postgres so request data and stats survive restarts.
- **Auth middleware** — extract `user_id` from a JWT instead of trusting the request body.


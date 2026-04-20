package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"rate-limiter-api/internal/constants"
	"rate-limiter-api/internal/infrastructure/ratelimiter"

	"github.com/gin-gonic/gin"
)

// RateLimit returns a Gin middleware that enforces per-user rate limiting.
// rate defines the policy (max requests + window) for this specific route.
func RateLimit(limiter ratelimiter.RateLimiter, rate ratelimiter.RateLimit) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Body == nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   constants.ErrBadRequest,
				"message": "request body is required",
			})
			return
		}

		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   constants.ErrBadRequest,
				"message": "failed to read request body",
			})
			return
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var peek struct {
			UserID string `json:"user_id"`
		}
		if err := json.Unmarshal(bodyBytes, &peek); err != nil || peek.UserID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   constants.ErrBadRequest,
				"message": "user_id is required",
			})
			return
		}

		allowed, remaining, resetAt, err := limiter.Allow(ctx, peek.UserID, rate)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   constants.ErrInternalServer,
				"message": "rate limiter unavailable",
			})
			return
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(rate.MaxRequests))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !allowed {
			ctx.Header("Retry-After", strconv.FormatInt(resetAt.Unix(), 10))
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   constants.ErrTooManyRequests,
				"message": "rate limit exceeded",
			})
			return
		}

		ctx.Next()
	}
}

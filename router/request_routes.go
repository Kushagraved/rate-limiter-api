package router

import (
	"time"

	"rate-limiter-api/internal/config"
	"rate-limiter-api/internal/controller"
	"rate-limiter-api/internal/infrastructure/ratelimiter"
	"rate-limiter-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func registerRequestRoutes(rg *gin.RouterGroup, srv *controller.Server, limiter ratelimiter.RateLimiter) {
	rl := config.Get().Settings.RateLimiter

	// POST /request — rate-limited per user
	rg.POST("/request", middleware.RateLimit(limiter, ratelimiter.RateLimit{
		MaxRequests: rl.MaxRequests,
		Window:      time.Duration(rl.WindowSecs) * time.Second,
	}), srv.PostRequest)

	// GET /stats — no rate limit (read-only, internal)
	rg.GET("/stats", srv.GetStats)
}

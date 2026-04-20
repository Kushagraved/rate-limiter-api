package router

import (
	"rate-limiter-api/internal/controller"
	"rate-limiter-api/internal/infrastructure/ratelimiter"
	"rate-limiter-api/internal/logger"
	"rate-limiter-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Init(srv *controller.Server, log logger.Logger) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger(log))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "ok"})
	})

	return r
}

func RegisterV1(r *gin.Engine, srv *controller.Server, limiter ratelimiter.RateLimiter) {
	v1 := r.Group("/v1")
	registerRequestRoutes(v1, srv, limiter)
}

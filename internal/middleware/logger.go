package middleware

import (
	"time"

	"rate-limiter-api/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RequestLogger(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		log.Info(ctx.Request.Context(), "http", "RequestLogger", "request completed",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", time.Since(start)),
		)
	}
}

package controller

import (
	"rate-limiter-api/internal/core/request"
	"rate-limiter-api/internal/logger"

	"github.com/gin-gonic/gin"
)

type Services struct {
	Request request.Service
}

type Server struct {
	services Services
	logger   logger.Logger
}

func New(services Services, log logger.Logger) *Server {
	return &Server{services: services, logger: log}
}

func throwError(ctx *gin.Context, code int, errCode, msg string) {
	ctx.JSON(code, map[string]interface{}{
		"error":   errCode,
		"message": msg,
	})
}

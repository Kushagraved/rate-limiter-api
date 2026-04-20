package controller

import (
	"net/http"

	apirequest "rate-limiter-api/api/request"
	"rate-limiter-api/internal/constants"
	"rate-limiter-api/internal/logger"

	"github.com/gin-gonic/gin"
)

// PostRequest handles POST /request
func (s *Server) PostRequest(ctx *gin.Context) {
	var req apirequest.SubmitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		throwError(ctx, http.StatusBadRequest, constants.ErrBadRequest, err.Error())
		return
	}

	resp, err := s.services.Request.Submit(ctx, req)
	if err != nil {
		s.logger.Error(ctx, "request", "PostRequest", "failed to submit", logger.Err(err))
		throwError(ctx, http.StatusInternalServerError, constants.ErrInternalServer, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetStats handles GET /stats — optional ?user_id= to filter to a single user.
func (s *Server) GetStats(ctx *gin.Context) {
	userID := ctx.Query("user_id")
	if userID == "" {
		throwError(ctx, http.StatusBadRequest, constants.ErrBadRequest, "user_id is required")
		return
	}

	resp, err := s.services.Request.Stats(ctx, userID)
	if err != nil {
		throwError(ctx, http.StatusInternalServerError, constants.ErrInternalServer, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

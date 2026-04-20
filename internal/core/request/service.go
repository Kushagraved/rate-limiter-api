package request

import (
	"context"
	"strconv"
	"time"

	apirequest "rate-limiter-api/api/request"
	"rate-limiter-api/internal/config"
	"rate-limiter-api/internal/helpers/timehelpers"
	"rate-limiter-api/internal/models"
	"rate-limiter-api/internal/repository"
)

// Service accepts and returns API-layer types directly, keeping the controller thin.
type Service interface {
	Submit(ctx context.Context, req apirequest.SubmitRequest) (apirequest.SubmitResponse, error)
	Stats(ctx context.Context, userID string) (apirequest.StatsResponse, error)
}

type service struct {
	repo  repository.RequestRepository
	clock timehelpers.Clock
}

func NewService(repo repository.RequestRepository, clock timehelpers.Clock) Service {
	return &service{repo: repo, clock: clock}
}

func (s *service) Submit(ctx context.Context, req apirequest.SubmitRequest) (apirequest.SubmitResponse, error) {
	record := models.UserRequest{
		UserID:    req.UserID,
		Payload:   req.Payload,
		CreatedAt: s.clock.Now(),
	}
	if err := s.repo.Save(ctx, record); err != nil {
		return apirequest.SubmitResponse{}, err
	}
	return apirequest.SubmitResponse{
		Success: true,
	}, nil
}

func (s *service) Stats(ctx context.Context, userID string) (apirequest.StatsResponse, error) {
	reqs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return apirequest.StatsResponse{}, err
	}
	return s.buildStats(userID, reqs), nil
}

func (s *service) buildStats(userID string, reqs []models.UserRequest) apirequest.StatsResponse {
	resp := apirequest.StatsResponse{UserID: userID}
	if len(reqs) == 0 {
		return resp
	}

	rl := config.Get().Settings.RateLimiter
	window := time.Duration(rl.WindowSecs) * time.Second
	now := s.clock.Now()
	windowStart := now.Add(-window)

	var windowReqs []models.UserRequest
	for _, r := range reqs {
		if r.CreatedAt.After(windowStart) {
			windowReqs = append(windowReqs, r)
		}
	}

	resp.TotalCount = len(reqs)
	resp.WindowCount = len(windowReqs)

	last := reqs[len(reqs)-1].CreatedAt
	if !last.IsZero() {
		resp.LastRequest = strconv.FormatInt(last.Unix(), 10)
	}

	// If at or over the limit, tell the client when the oldest window entry expires.
	if len(windowReqs) >= rl.MaxRequests {
		retryAt := windowReqs[0].CreatedAt.Add(window)
		resp.RetryAfter = strconv.FormatInt(retryAt.Unix(), 10)
	}

	return resp
}

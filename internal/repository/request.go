package repository

import (
	"context"
	"rate-limiter-api/internal/models"
	"sync"
)

// RequestRepository is the data-access contract.
// It stores and retrieves raw UserRequest records only.
// Stats computation belongs in the core layer.
type RequestRepository interface {
	Save(ctx context.Context, req models.UserRequest) error
	GetByUserID(ctx context.Context, userID string) ([]models.UserRequest, error)
	GetAll(ctx context.Context) (map[string][]models.UserRequest, error)
}

type inMemoryRequestRepo struct {
	mu       sync.RWMutex
	requests map[string][]models.UserRequest
}

func NewInMemoryRequestRepo() RequestRepository {
	return &inMemoryRequestRepo{requests: make(map[string][]models.UserRequest)}
}

func (r *inMemoryRequestRepo) Save(_ context.Context, req models.UserRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.requests[req.UserID] = append(r.requests[req.UserID], req)
	return nil
}

func (r *inMemoryRequestRepo) GetByUserID(_ context.Context, userID string) ([]models.UserRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	src := r.requests[userID]
	out := make([]models.UserRequest, len(src))
	copy(out, src)
	return out, nil
}

func (r *inMemoryRequestRepo) GetAll(_ context.Context) (map[string][]models.UserRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string][]models.UserRequest, len(r.requests))
	for k, v := range r.requests {
		cp := make([]models.UserRequest, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out, nil
}

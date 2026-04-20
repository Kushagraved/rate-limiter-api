package request

// SubmitRequest is the JSON body for POST /request.
type SubmitRequest struct {
	UserID  string      `json:"user_id"  binding:"required"`
	Payload interface{} `json:"payload"  binding:"required"`
}

// SubmitResponse is returned on a successful POST /request.
type SubmitResponse struct {
	Success bool `json:"success"`
}

// StatsResponse is returned by GET /stats.
type StatsResponse struct {
	UserID      string `json:"user_id"`
	TotalCount  int    `json:"total_count"`
	WindowCount int    `json:"window_count"`
	LastRequest string `json:"last_request,omitempty"`
	RetryAfter  string `json:"retry_after,omitempty"`
}

// AllStatsResponse wraps a list of per-user stats.
type AllStatsResponse struct {
	Stats []StatsResponse `json:"stats"`
	Total int             `json:"total"`
}

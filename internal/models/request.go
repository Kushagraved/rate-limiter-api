package models

import "time"

// UserRequest is the core↔repository data model for a single accepted request.
type UserRequest struct {
	UserID    string      `json:"user_id"`
	Payload   interface{} `json:"payload"`
	CreatedAt time.Time   `json:"created_at"`
}

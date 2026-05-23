package models

import "time"

type RequestBody struct {
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
}

type UserStats struct {
	AcceptedCount int
	RejectedCount int
	WindowStart   time.Time
}

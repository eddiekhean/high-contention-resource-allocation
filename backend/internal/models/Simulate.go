package models

import "time"

type Simulation struct {
	ID            string    `json:"id"`
	Policy        string    `json:"policy"`
	Slots         int       `json:"slots"`
	TotalRequests int       `json:"total_requests"`
	Seed          int       `json:"seed"`
	CreatedAt     time.Time `json:"created_at"`
}

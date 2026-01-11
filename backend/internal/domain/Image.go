package domain

import "time"

type Image struct {
	ID        int64
	S3Key     string
	DHash     string
	CreatedAt time.Time
	SignedURL string `json:"-"` // runtime-only data
}

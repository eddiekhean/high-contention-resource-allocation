package domain

import "time"

type Image struct {
	ID        int64
	URL       string
	DHash     uint64
	CreatedAt time.Time
}

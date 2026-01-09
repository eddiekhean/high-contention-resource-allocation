package domain

import "time"

type Image struct {
	ID          int64
	Key         string // s3 object key
	Bucket      string // bucket name
	URL         string // public hoặc CDN URL
	Size        int64  // bytes
	ContentType string // image/jpeg, image/png
	ETag        string // checksum từ S3
	CreatedAt   time.Time
}

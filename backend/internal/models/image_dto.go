package models

type MatchImageRequest struct {
	DHash uint64 `json:"dhash" binding:"required"`
}

type ImageDTO struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	DHash uint64 `json:"dhash"`
}

type MatchImageResponse struct {
	Matched bool      `json:"matched"`
	Image   *ImageDTO `json:"image,omitempty"`
}

type UploadURLResponse struct {
	UploadURL string `json:"uploadUrl"`
	S3Key     string `json:"s3Key"`
}

type CommitImageRequest struct {
	DHash uint64 `json:"dhash" binding:"required"`
	S3Key string `json:"s3Key" binding:"required"`
}

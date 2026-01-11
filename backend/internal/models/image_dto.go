package models

type MatchImageRequest struct {
	DHash string `json:"dhash" binding:"required"`
}

type ImageDTO struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	DHash string `json:"dhash"`
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
	DHash string `json:"dhash" binding:"required"`
	S3Key string `json:"s3Key" binding:"required"`
}

type PublicImagesResponse struct {
	Images []*ImageDTO `json:"images"`
}

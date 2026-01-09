package handler

import (
	"net/http"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MazeHandler struct {
	imageService service.ImageService
	logger       *logrus.Logger
}

func NewMazeHandler(imageService service.ImageService, logger *logrus.Logger) *MazeHandler {
	return &MazeHandler{
		imageService: imageService,
		logger:       logger,
	}
}

func (h *MazeHandler) MatchImage(c *gin.Context) {
	var req models.MatchImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	match, matched, err := h.imageService.MatchImage(c.Request.Context(), req.DHash)
	if err != nil {
		h.logger.WithError(err).Error("match image failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if !matched {
		c.JSON(http.StatusOK, models.MatchImageResponse{Matched: false})
		return
	}

	c.JSON(http.StatusOK, models.MatchImageResponse{
		Matched: true,
		Image: &models.ImageDTO{
			ID:    match.ID,
			URL:   match.URL, // signed URL
			DHash: match.DHash,
		},
	})
}

// POST /upload-url
func (h *MazeHandler) GetUploadURL(c *gin.Context) {
	uploadURL, s3Key, err := h.imageService.GetUploadURL(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("get upload url failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, models.UploadURLResponse{
		UploadURL: uploadURL,
		S3Key:     s3Key,
	})
}

// POST / (commit)
func (h *MazeHandler) CommitImage(c *gin.Context) {
	var req models.CommitImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	img, err := h.imageService.CommitImage(c.Request.Context(), req.DHash, req.S3Key)
	if err != nil {
		h.logger.WithError(err).Error("commit image failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.ImageDTO{
		ID:    img.ID,
		URL:   img.URL, // s3 key
		DHash: img.DHash,
	})
}

package handler

import (
	"net/http"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/apperr"
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
		c.JSON(apperr.ToResponse(apperr.ErrInvalidRequest))
		return
	}

	match, matched, err := h.imageService.MatchImage(c.Request.Context(), req.DHash)
	if err != nil {
		h.logger.WithError(err).Error("match image failed")
		c.JSON(apperr.ToResponse(err))
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
			URL:   match.SignedURL, // signed URL from service
			DHash: match.DHash,
		},
	})
}

// POST /upload-url
func (h *MazeHandler) GetUploadURL(c *gin.Context) {
	uploadURL, s3Key, err := h.imageService.GetUploadURL(c.Request.Context())
	if err != nil {
		h.logger.WithError(err).Error("get upload url failed")
		c.JSON(apperr.ToResponse(err))
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
		h.logger.WithError(err).Warn("invalid commit request body")
		c.JSON(apperr.ToResponse(apperr.ErrInvalidRequest))
		return
	}

	img, err := h.imageService.CommitImage(c.Request.Context(), req.DHash, req.S3Key)
	if err != nil {
		h.logger.WithError(err).Error("commit image failed")
		c.JSON(apperr.ToResponse(err))
		return
	}

	c.JSON(http.StatusCreated, models.ImageDTO{
		ID:    img.ID,
		URL:   img.S3Key, // s3 key for commit response (or we could return signed URL if needed)
		DHash: img.DHash,
	})
}

// GET /public/images
func (h *MazeHandler) GetPublicImages(c *gin.Context) {
	images, err := h.imageService.GetRandomImages(c.Request.Context(), 4)
	if err != nil {
		h.logger.WithError(err).Error("get public images failed")
		c.JSON(apperr.ToResponse(err))
		return
	}

	var dtos []*models.ImageDTO
	for _, img := range images {
		dtos = append(dtos, &models.ImageDTO{
			ID:    img.ID,
			URL:   img.SignedURL, // signed URL from service
			DHash: img.DHash,
		})
	}

	c.JSON(http.StatusOK, models.PublicImagesResponse{
		Images: dtos,
	})
}

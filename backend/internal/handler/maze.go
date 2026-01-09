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
	var req models.MatchMazeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	img, dist, matched, err := h.imageService.Match(
		c.Request.Context(),
		req.DHash,
	)
	if err != nil {
		h.logger.WithError(err).Error("image match failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	if !matched {
		c.JSON(http.StatusOK, models.MazeMatchResponse{
			Matched: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.MazeMatchResponse{
		Matched:  true,
		Distance: dist,
		Maze: &models.MazeDTO{
			ID:  img.ID,
			URL: img.URL,
		},
	})
}

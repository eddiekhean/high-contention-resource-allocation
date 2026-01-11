package handler

import (
	"net/http"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MazeHandler struct {
	service *service.MazeService
	logger  *logrus.Logger
}

func NewMazeHandler(service *service.MazeService, logger *logrus.Logger) *MazeHandler {
	return &MazeHandler{
		service: service,
		logger:  logger,
	}
}

func (h *MazeHandler) Generate(c *gin.Context) {
	var input models.MazeConfig

	// Bind JSON body if available, otherwise use defaults/query params?
	// The requirement is POST, so JSON body is expected.
	if err := c.ShouldBindJSON(&input); err != nil {
		// Just log, maybe they sent empty body which is fine if we have defaults?
		// But ShouldBindJSON returns error on empty body usually?
		// Use BindJSON?
		// Actually, if we want to allow empty body to use defaults, we might handle EOF.
		// For now, let's assume they might send partial data.
		// If error is basic format error, return 400.
		h.logger.Warnf("Invalid input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	maze, err := h.service.GenerateMaze(input)
	if err != nil {
		h.logger.Errorf("Failed to generate maze: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate maze"})
		return
	}

	c.JSON(http.StatusOK, maze)
}

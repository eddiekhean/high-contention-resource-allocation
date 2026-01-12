package handler

import (
	"errors"
	"net/http"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service/maze"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MazeHandler struct {
	service *maze.MazeService
	logger  *logrus.Logger
}

func NewMazeHandler(service *maze.MazeService, logger *logrus.Logger) *MazeHandler {
	return &MazeHandler{
		service: service,
		logger:  logger,
	}
}
func (h *MazeHandler) Submit(c *gin.Context) {
	var req models.MazeSolveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := utils.FormatValidationError(err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	result, err := h.service.SolveMaze(&req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *MazeHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, utils.ErrInvalidRequest),
		errors.Is(err, utils.ErrInvalidStrategy):
		c.JSON(http.StatusBadRequest, utils.NewAPIError(http.StatusBadRequest, err.Error()))

	case errors.Is(err, utils.ErrMazeUnreachable):
		c.JSON(http.StatusUnprocessableEntity, utils.NewAPIError(http.StatusUnprocessableEntity, err.Error()))

	default:
		h.logger.Errorf("internal error: %+v", err)
		c.JSON(http.StatusInternalServerError, utils.NewAPIError(http.StatusInternalServerError, "internal server error"))
	}
}
func (h *MazeHandler) Generate(c *gin.Context) {
	var input models.MazeConfig

	// Bind JSON body if available, otherwise use defaults/query params?
	// The requirement is POST, so JSON body is expected.
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warnf("Invalid input: %v", err)
		apiErr := utils.FormatValidationError(err)
		c.JSON(apiErr.Code, apiErr)
		return
	}

	maze, err := h.service.GenerateMaze(input)
	if err != nil {
		h.logger.Errorf("Failed to generate maze: %v", err)
		c.JSON(http.StatusInternalServerError, utils.NewAPIError(http.StatusInternalServerError, "Failed to generate maze"))
		return
	}

	c.JSON(http.StatusOK, maze)
}

package handler

import (
	"net/http"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/apperr"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/service"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SimulateHandler struct {
	service *service.SimulateService
	logger  *logrus.Logger
}

func NewSimulateHandler(
	service *service.SimulateService,
	logger *logrus.Logger,
) *SimulateHandler {
	return &SimulateHandler{
		service: service,
		logger:  logger,
	}
}
func (h *SimulateHandler) Simulate(c *gin.Context) {

	var input models.SimulationRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(apperr.ToResponse(apperr.ErrInvalidRequest))
		return
	}

	events, err := h.service.RunSimulation(input)
	if err != nil {
		c.JSON(apperr.ToResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
	})
}

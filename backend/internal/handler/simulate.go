package handler

import (
	"net/http"

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
func Simulate(c *gin.Context) {
	var input models.SimulateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON không hợp lệ"})
		return
	}

	simulationID, events, err := service.RunSimulation(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulation_id": simulationID,
		"events":        events,
	})
}

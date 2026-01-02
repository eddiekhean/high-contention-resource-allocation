package service

import (
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SimulateService struct {
	logger    *logrus.Logger
	slotStore *storage.SlotStore
}

func NewSimulateService(logger *logrus.Logger, store *storage.SlotStore) *SimulateService {
	return &SimulateService{
		logger:    logger,
		slotStore: store,
	}
}

func (s *SimulateService) RunSimulation(ctx *gin.Context, input models.SimulateRequest) (string, []models.Event, error) {
	// 1. Tạo simulation_id
	simulationID := uuid.New().String()

	//2. Init slot trong Redis
	if err := s.slotStore.InitSlot(ctx, simulationID, input.Slots); err != nil {
		return "", nil, err
	}

	// 3. Generate request giả lập
	requests := scheduler.GenerateRequests(
		input.TotalRequests,
		input.Seed,
	)
	//
	//// 4. Chạy scheduler + Redis
	//events := scheduler.RunSimulationWithRedis(
	//	simulationID,
	//	requests,
	//	input.Policy,
	//)

	return simulationID, events, nil
}

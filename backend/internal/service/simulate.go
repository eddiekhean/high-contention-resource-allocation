package service

import (
	"context"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/scheduler"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/storage"
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
func (s *SimulateService) RunSimulation(
	input models.SimulationRequest,
) (*models.SimulateResponse, error) {

	ctx := context.Background()
	simID := uuid.New().String()

	// 1. Init voucher
	if err := s.slotStore.InitSlot(ctx, simID, input.TotalVouchers); err != nil {
		return nil, err
	}

	// 2. Generate workload
	clients := scheduler.GenerateClients(input.Seed, input.TotalClients)
	requests, arrivalOrder := scheduler.GenerateRequests(clients, input.Seed)

	// 3. Scheduler quyết định thứ tự
	decisions := scheduler.RunHybridScheduler(requests)

	// 4. Execute decisions
	var events []models.Event

	for _, d := range decisions {
		ok, err := s.slotStore.TryAcquire(ctx, simID, 1)
		if err != nil {
			return nil, err
		}

		action := "allocated"
		if !ok {
			action = "rejected"
		}

		events = append(events, models.Event{
			Tick:      d.Tick,
			RequestID: d.Request.ID,
			ClientID:  d.Request.ClientID,
			Priority:  d.Request.Priority,
			Score:     d.Score,
			Action:    action,
		})

		if !ok {
			break
		}
	}

	s.slotStore.Clear(ctx, simID)

	// 5. BUILD RESPONSE
	resp := &models.SimulateResponse{
		Simulation: models.Simulation{
			ID: simID,
		},
		ArrivalOrder: arrivalOrder,
		Events:       events,
	}

	return resp, nil
}

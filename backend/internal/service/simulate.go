package service

import (
	"context"
	"time"

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

	// 3. Scheduler selects strategy
	strategyFactory := scheduler.NewStrategyFactory()
	strategy := strategyFactory.Get(input.Policy)
	if strategy == nil {
		// Fallback or error. For now, defaulting to hybrid if not found, or maybe just error?
		// Given validation in DTO, input.Policy should be valid.
		// However, let's default to hybrid if something goes wrong or for "fairness" later.
		strategy = strategyFactory.Get("hybrid")
	}

	decisions := strategy.Schedule(requests)

	// 4. Execute decisions
	var events []models.Event

	for _, d := range decisions {
		action := "rejected"

		ok, err := s.slotStore.TryAcquire(ctx, simID, 1)
		if err != nil {
			return nil, err
		}

		if ok {
			action = "allocated"
		}

		events = append(events, models.Event{
			Tick:      d.Tick,
			RequestID: d.Request.ID,
			ClientID:  d.Request.ClientID,
			Priority:  d.Request.Priority,
			Score:     d.Score,
			Action:    action,
		})
	}

	s.slotStore.Clear(ctx, simID)

	// 5. BUILD RESPONSE
	resp := &models.SimulateResponse{
		Simulation: models.Simulation{
			ID:            simID,
			Slots:         input.TotalVouchers,
			TotalRequests: len(requests),
			Policy:        "hybrid",
			CreatedAt:     time.Now(),
		},
		ArrivalOrder: arrivalOrder,
		Events:       events,
	}

	return resp, nil
}

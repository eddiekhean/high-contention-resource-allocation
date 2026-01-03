package scheduler

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

type runtimeRequest struct {
	models.Request
	EnqueueTick int
}
type HybridConfig struct {
	Alpha float64 // priority weight
	Beta  float64 // waiting time weight
	Gamma float64 // fairness penalty
}

func computeScore(
	req runtimeRequest,
	now int,
	clientDebt map[int]float64,
	cfg HybridConfig,
) float64 {

	waiting := now - req.EnqueueTick

	return float64(req.Priority)*cfg.Alpha +
		float64(waiting)*cfg.Beta -
		clientDebt[req.ClientID]*cfg.Gamma
}

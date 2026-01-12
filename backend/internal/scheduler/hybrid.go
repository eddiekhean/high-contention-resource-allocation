package scheduler

import (
	"sort"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
)

type HybridStrategy struct {
	cfg HybridConfig
}

func NewHybridStrategy() *HybridStrategy {
	return &HybridStrategy{
		cfg: HybridConfig{Alpha: 10, Beta: 1, Gamma: 2},
	}
}

func (s *HybridStrategy) Name() string {
	return "hybrid"
}

func (s *HybridStrategy) Schedule(requests []models.Request) []Decision {
	var decisions []Decision
	var queue []runtimeRequest
	clientDebt := map[int]float64{}

	reqIdx := 0
	tick := 0

	for reqIdx < len(requests) || len(queue) > 0 {

		for reqIdx < len(requests) && requests[reqIdx].ArrivalAt <= tick {
			queue = append(queue, runtimeRequest{
				Request:     requests[reqIdx],
				EnqueueTick: tick,
			})
			reqIdx++
		}

		if len(queue) == 0 {
			tick++
			continue
		}

		sort.Slice(queue, func(i, j int) bool {
			return computeScore(queue[i], tick, clientDebt, s.cfg) >
				computeScore(queue[j], tick, clientDebt, s.cfg)
		})

		selected := queue[0]
		score := computeScore(selected, tick, clientDebt, s.cfg)

		decisions = append(decisions, Decision{
			Tick:    tick,
			Request: selected.Request,
			Score:   score,
		})

		clientDebt[selected.ClientID]++
		queue = queue[1:]
		tick++
	}

	return decisions
}

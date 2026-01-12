package scheduler

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

// Decision represents a scheduling decision
type Decision struct {
	Tick    int
	Request models.Request
	Score   float64
}

// Strategy defines the interface for different scheduling algorithms
type Strategy interface {
	Name() string
	Schedule(requests []models.Request) []Decision
}

// StrategyFactory handles the creation/retrieval of strategies
type StrategyFactory struct {
	strategies map[string]Strategy
}

func NewStrategyFactory() *StrategyFactory {
	f := &StrategyFactory{
		strategies: make(map[string]Strategy),
	}
	// Register default strategies
	f.Register(NewHybridStrategy())
	return f
}

func (f *StrategyFactory) Register(s Strategy) {
	f.strategies[s.Name()] = s
}

func (f *StrategyFactory) Get(name string) Strategy {
	return f.strategies[name]
}

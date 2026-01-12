package maze

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

type MazeSolver interface {
	Solve(
		start, end models.Point,
		graph map[models.Point]*models.Cell,
	) (models.SolveResult, bool)
}

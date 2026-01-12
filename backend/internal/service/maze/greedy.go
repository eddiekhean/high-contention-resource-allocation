package maze

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

type GreedySolver struct{}

func (s *GreedySolver) Solve(
	start, end models.Point,
	graph map[models.Point]*models.Cell,
) (models.SolveResult, bool) {

	open := []models.Point{start}

	parent := map[models.Point]*models.Point{
		start: nil,
	}

	var steps []models.SolveStep
	stepIdx := 0

	for len(open) > 0 {
		// ===== pick node with smallest heuristic =====
		bestIdx := 0
		for i := range open {
			if manhattan(open[i], end) < manhattan(open[bestIdx], end) {
				bestIdx = i
			}
		}

		curr := open[bestIdx]
		open = append(open[:bestIdx], open[bestIdx+1:]...)

		steps = append(steps, models.SolveStep{
			Step:  stepIdx,
			Type:  models.StepVisit,
			Point: curr,
		})
		stepIdx++

		if curr == end {
			path := reconstructPath(end, parent)
			for _, p := range path {
				steps = append(steps, models.SolveStep{
					Step:  stepIdx,
					Type:  models.StepPath,
					Point: p,
				})
				stepIdx++
			}
			return models.SolveResult{Path: path, Steps: steps}, true
		}

		cell := graph[curr]
		for _, next := range neighbors(curr, cell) {
			if _, visited := parent[next]; visited {
				continue
			}

			parent[next] = &curr
			open = append(open, next)

			steps = append(steps, models.SolveStep{
				Step:  stepIdx,
				Type:  models.StepEnqueue,
				Point: next,
			})
			stepIdx++
		}
	}

	return models.SolveResult{Steps: steps}, false
}

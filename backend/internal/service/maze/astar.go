package maze

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

type AStarSolver struct{}

func (s *AStarSolver) Solve(
	start, end models.Point,
	graph map[models.Point]*models.Cell,
) (models.SolveResult, bool) {

	open := []models.Point{start}

	parent := map[models.Point]*models.Point{
		start: nil,
	}

	gScore := map[models.Point]int{
		start: 0,
	}

	var steps []models.SolveStep
	stepIdx := 0

	for len(open) > 0 {
		// ===== pick node with smallest f = g + h =====
		bestIdx := 0
		for i := range open {
			f := gScore[open[i]] + manhattan(open[i], end)
			bestF := gScore[open[bestIdx]] + manhattan(open[bestIdx], end)
			if f < bestF {
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
			tentativeG := gScore[curr] + 1

			if g, ok := gScore[next]; ok && tentativeG >= g {
				continue
			}

			parent[next] = &curr
			gScore[next] = tentativeG
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

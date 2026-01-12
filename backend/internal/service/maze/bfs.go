package maze

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

type BFSSolver struct{}

func (s *BFSSolver) Solve(
	start, end models.Point,
	graph map[models.Point]*models.Cell,
) (models.SolveResult, bool) {

	queue := []models.Point{start}

	// parent dùng để reconstruct path + cũng là visited
	parent := map[models.Point]*models.Point{
		start: nil,
	}

	var steps []models.SolveStep
	stepIdx := 0

	for len(queue) > 0 {
		// ===== POP =====
		curr := queue[0]
		queue = queue[1:]

		steps = append(steps, models.SolveStep{
			Step:  stepIdx,
			Type:  models.StepVisit,
			Point: curr,
		})
		stepIdx++

		// ===== FOUND END =====
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

			return models.SolveResult{
				Path:  path,
				Steps: steps,
			}, true
		}

		// ===== EXPAND =====
		cell := graph[curr]
		for _, next := range neighbors(curr, cell) {
			if _, visited := parent[next]; visited {
				continue
			}

			parent[next] = &curr
			queue = append(queue, next)

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

package maze

import "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"

func (s *MazeService) buildGraph(
	req *models.MazeSolveRequest,
) map[models.Point]*models.Cell {

	graph := make(map[models.Point]*models.Cell, len(req.Cells))
	for i := range req.Cells {
		cell := &req.Cells[i]
		graph[models.Point{X: cell.X, Y: cell.Y}] = cell
	}
	return graph
}
func reconstructPath(
	end models.Point,
	parent map[models.Point]*models.Point,
) []models.Point {

	var path []models.Point
	for p := &end; p != nil; {
		path = append(path, *p)
		p = parent[*p]
	}

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func neighbors(p models.Point, cell *models.Cell) []models.Point {
	var res []models.Point

	if !cell.Walls.Top {
		res = append(res, models.Point{X: p.X, Y: p.Y - 1})
	}
	if !cell.Walls.Right {
		res = append(res, models.Point{X: p.X + 1, Y: p.Y})
	}
	if !cell.Walls.Bottom {
		res = append(res, models.Point{X: p.X, Y: p.Y + 1})
	}
	if !cell.Walls.Left {
		res = append(res, models.Point{X: p.X - 1, Y: p.Y})
	}

	return res
}
func manhattan(a, b models.Point) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

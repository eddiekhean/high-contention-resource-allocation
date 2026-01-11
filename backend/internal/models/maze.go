package models

type MazeConfig struct {
	Rows      int     `json:"rows"`
	Cols      int     `json:"cols"`
	LoopRatio float64 `json:"loop_ratio,omitempty"` // 0.0 to 1.0
	Seed      *int64  `json:"seed,omitempty"`       // Optional seed
}

type MazePoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type MazeCell struct {
	X     int     `json:"x"`
	Y     int     `json:"y"`
	Walls [4]bool `json:"walls"` // [top, right, bottom, left]
}

type MazeResponse struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Start MazePoint  `json:"start"`
	End   MazePoint  `json:"end"`
	Cells []MazeCell `json:"cells"`
}

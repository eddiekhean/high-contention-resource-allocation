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

type Walls struct {
	Top    bool `json:"top"`
	Right  bool `json:"right"`
	Bottom bool `json:"bottom"`
	Left   bool `json:"left"`
}

type MazeCell struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Walls Walls `json:"walls"`
}

type MazeResponse struct {
	Rows  int        `json:"rows"`
	Cols  int        `json:"cols"`
	Start MazePoint  `json:"start"`
	End   MazePoint  `json:"end"`
	Cells []MazeCell `json:"cells"`
}

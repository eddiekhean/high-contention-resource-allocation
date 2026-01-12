package models

type SolveStrategy string

const (
	StrategyBFS    SolveStrategy = "BFS"
	StrategyDFS    SolveStrategy = "DFS"
	StrategyAStar  SolveStrategy = "ASTAR"
	StrategyGreedy SolveStrategy = "GREEDY"
)

type StepType string

const (
	StepVisit   StepType = "VISIT"   // pop khỏi queue
	StepEnqueue StepType = "ENQUEUE" // push vào queue
	StepPath    StepType = "PATH"    // đường cuối
)

func (s SolveStrategy) IsValid() bool {
	switch s {
	case StrategyBFS, StrategyDFS, StrategyAStar, StrategyGreedy:
		return true
	}
	return false
}

type MazeConfig struct {
	Rows      int     `json:"rows" binding:"required,min=5,max=100"`
	Cols      int     `json:"cols" binding:"required,min=5,max=100"`
	LoopRatio float64 `json:"loop_ratio,omitempty" binding:"min=0,max=1"` // 0.0 to 1.0
	Seed      *int64  `json:"seed,omitempty"`                             // Optional seed
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

// MazeSolveRequest represents the API request body
type MazeSolveRequest struct {
	Rows     int           `json:"rows" binding:"required,min=5,max=100"`
	Cols     int           `json:"cols" binding:"required,min=5,max=100"`
	Start    Point         `json:"start" binding:"required"`
	End      Point         `json:"end" binding:"required"`
	Cells    []Cell        `json:"cells" binding:"required"`
	Strategy SolveStrategy `json:"strategy" binding:"required"`
}
type SolveResult struct {
	Path  []Point     // đường đi cuối
	Steps []SolveStep // timeline đầy đủ
}

type SolveStep struct {
	Step  int
	Type  StepType
	Point Point
}

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Cell struct {
	X     int   `json:"x"`
	Y     int   `json:"y"`
	Walls Walls `json:"walls"`
}

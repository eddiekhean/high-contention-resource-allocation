package models

type SimulateResponse struct {
	Simulation   Simulation `json:"simulation"`
	ArrivalOrder []ClientArrival
	Events       []Event `json:"events"`
}
type Event struct {
	Tick      int
	RequestID int
	ClientID  int
	Priority  int
	Score     float64
	Action    string // enqueue | allocated | wait | drop
}
type MazeMatchResponse struct {
	Matched  bool     `json:"matched"`
	Distance int      `json:"distance,omitempty"`
	Maze     *MazeDTO `json:"maze,omitempty"`
}

type MazeDTO struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

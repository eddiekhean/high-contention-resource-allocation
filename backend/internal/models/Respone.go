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

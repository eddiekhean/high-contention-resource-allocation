package models

type SimulateResponse struct {
	Simulation Simulation `json:"simulation"`
	Events     []Event    `json:"events"`
}
type Event struct {
	SimulationID string                 `json:"simulation_id"`
	Tick         int                    `json:"tick"`
	Type         string                 `json:"type"`
	Payload      map[string]interface{} `json:"payload"`
}

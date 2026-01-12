package models

type SimulationRequest struct {
	TotalClients  int    `json:"total_clients" binding:"required,gt=0" validate:"gt=0"`
	TotalVouchers int    `json:"total_vouchers" binding:"required,gt=0" validate:"gt=0"`
	Seed          int64  `json:"seed"`
	Policy        string `json:"policy" binding:"required,oneof=fifo priority lottery hybrid" validate:"oneof=fifo priority lottery hybrid"`
}
type Client struct {
	ID     int
	Class  string  // vip / paid / free
	Weight float64 // fairness weight
}
type Request struct {
	ID        int
	ClientID  int
	Priority  int // 1 = cao nhất
	ArrivalAt int // logical time (tick)
}
type RuntimeRequest struct {
	Request
	EnqueueTick int
	WaitingTime int
}
type ClientArrival struct {
	ClientID  int    `json:"client_id"`
	Class     string `json:"class"`
	FirstTick int    `json:"first_tick"`
}

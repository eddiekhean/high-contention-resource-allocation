package scheduler

import (
	"math/rand"
	"sort"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
)

func GenerateClients(
	seed int64,
	totalClients int,
) []models.Client {

	rng := rand.New(rand.NewSource(seed))

	clients := make([]models.Client, 0, totalClients)

	for i := 0; i < totalClients; i++ {
		roll := rng.Float64()

		client := models.Client{
			ID: i + 1,
		}

		// Distribution:
		// 10% VIP, 30% Paid, 60% Free
		switch {
		case roll < 0.10:
			client.Class = "vip"
			client.Weight = 1.5
		case roll < 0.40:
			client.Class = "paid"
			client.Weight = 1.0
		default:
			client.Class = "free"
			client.Weight = 0.7
		}

		clients = append(clients, client)
	}

	return clients
}
func priorityFromClass(class string) int {
	switch class {
	case "vip":
		return 1
	case "paid":
		return 2
	default:
		return 3
	}
}

// GenerateRequests sinh request từ danh sách client
// - mỗi client gửi 1–3 request
// - arrival có burst (Gaussian-like)
// - priority theo client class
func GenerateRequests(
	clients []models.Client,
	seed int64,
) []models.Request {

	rng := rand.New(rand.NewSource(seed + 1)) // +1 để khác GenerateClients

	var requests []models.Request
	reqID := 1

	for _, c := range clients {
		// số request mỗi client (heterogeneous)
		reqCount := rng.Intn(3) + 1 // 1–3

		for i := 0; i < reqCount; i++ {
			req := models.Request{
				ID:        reqID,
				ClientID:  c.ID,
				Priority:  priorityFromClass(c.Class),
				ArrivalAt: gaussianArrival(rng),
			}
			requests = append(requests, req)
			reqID++
		}
	}

	// sort theo arrival time (giả lập dòng thời gian)
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].ArrivalAt != requests[j].ArrivalAt {
			return requests[i].ArrivalAt < requests[j].ArrivalAt
		}
		return requests[i].ID < requests[j].ID
	})

	return requests
}
func gaussianArrival(rng *rand.Rand) int {
	// Gaussian centered around tick = 50
	v := int(rng.NormFloat64()*10 + 50)

	if v < 0 {
		return 0
	}
	return v
}

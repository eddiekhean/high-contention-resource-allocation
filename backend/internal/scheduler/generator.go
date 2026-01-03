package scheduler

import (
	"fmt"
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
	for _, c := range clients {
		fmt.Printf("ID=%d | class=%s | weight=%.1f\n",
			c.ID, c.Class, c.Weight,
		)
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
) ([]models.Request, []models.ClientArrival) {

	rng := rand.New(rand.NewSource(seed + 1))

	var requests []models.Request
	reqID := 1

	for _, c := range clients {
		reqCount := rng.Intn(3) + 1

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

	// sort theo arrival time
	sort.Slice(requests, func(i, j int) bool {
		if requests[i].ArrivalAt != requests[j].ArrivalAt {
			return requests[i].ArrivalAt < requests[j].ArrivalAt
		}
		return requests[i].ID < requests[j].ID
	})

	// === BUILD ARRIVAL ORDER ===
	clientMap := make(map[int]models.Client)
	for _, c := range clients {
		clientMap[c.ID] = c
	}

	seen := make(map[int]bool)
	var arrivalOrder []models.ClientArrival

	for _, r := range requests {
		if seen[r.ClientID] {
			continue
		}

		c := clientMap[r.ClientID]
		arrivalOrder = append(arrivalOrder, models.ClientArrival{
			ClientID:  c.ID,
			Class:     c.Class,
			FirstTick: r.ArrivalAt,
		})

		seen[r.ClientID] = true
	}

	return requests, arrivalOrder
}

func gaussianArrival(rng *rand.Rand) int {
	// Gaussian centered around tick = 50
	v := int(rng.NormFloat64()*10 + 50)

	if v < 0 {
		return 0
	}
	return v
}

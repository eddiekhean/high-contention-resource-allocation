package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type SlotStore struct {
	rdb *redis.Client
}

func NewSlotStore(rdb *redis.Client) *SlotStore {
	return &SlotStore{
		rdb: rdb,
	}
}

// redis key: simulation:{id}:slots
func slotKey(simulationID string) string {
	return fmt.Sprintf("simulation:%s:slots", simulationID)
}

// InitSlot khởi tạo số slot ban đầu cho 1 simulation
// Chỉ gọi 1 lần khi start simulation
func (s *SlotStore) InitSlot(
	ctx context.Context,
	simulationID string,
	slots int,
) error {
	key := slotKey(simulationID)

	ok, err := s.rdb.SetNX(ctx, key, slots, 0).Result()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("slot already initialized for simulation %s", simulationID)
	}

	return nil
}

// TryAcquire thử chiếm slot (atomic)
// return true nếu chiếm được
func (s *SlotStore) TryAcquire(
	ctx context.Context,
	simulationID string,
	n int,
) (bool, error) {
	key := slotKey(simulationID)

	val, err := s.rdb.DecrBy(ctx, key, int64(n)).Result()
	if err != nil {
		return false, err
	}

	if val < 0 {
		// rollback nếu thiếu slot
		_, _ = s.rdb.IncrBy(ctx, key, int64(n)).Result()
		return false, nil
	}

	return true, nil
}

// Release trả slot lại
func (s *SlotStore) Release(
	ctx context.Context,
	simulationID string,
	n int,
) error {
	key := slotKey(simulationID)
	return s.rdb.IncrBy(ctx, key, int64(n)).Err()
}

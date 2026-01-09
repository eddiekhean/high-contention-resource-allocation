package repository

import (
	"context"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
)

type ImageRepository interface {
	Create(ctx context.Context, img *domain.Image) error
	FindByID(ctx context.Context, id int64) (*domain.Image, error)
	FindByDHash(ctx context.Context, hash uint64) (*domain.Image, error)
	GetAll(ctx context.Context) ([]*domain.Image, error)

	// image similarity prefix-based (keeping for compatibility or optimization)
	FindByPrefix(
		ctx context.Context,
		prefix uint16,
		rangeSize int,
	) ([]*domain.Image, error)
}

package repository

import (
	"context"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
)

type ImageRepository interface {
	// basic
	Create(ctx context.Context, img *domain.Image) error
	FindByID(ctx context.Context, id int64) (*domain.Image, error)

	// image similarity
	FindByPrefix(
		ctx context.Context,
		prefix uint16,
		rangeSize int,
	) ([]*domain.Image, error)
}

package repository

import (
	"context"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
)

type ImageRepository interface {
	Create(ctx context.Context, img *domain.Image) error
	FindByID(ctx context.Context, id int64) (*domain.Image, error)
}

package repository

import (
	"context"
	"fmt"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/db"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
)

type PgImageRepository struct {
	db *db.Postgres
}

func NewPgImageRepository(db *db.Postgres) ImageRepository {
	return &PgImageRepository{db: db}
}
func (r *PgImageRepository) Create(
	ctx context.Context,
	img *domain.Image,
) error {
	return fmt.Errorf("Create not implemented")
}

func (r *PgImageRepository) FindByID(
	ctx context.Context,
	id int64,
) (*domain.Image, error) {
	return nil, fmt.Errorf("FindByID not implemented")
}

func (r *PgImageRepository) FindByPrefix(
	ctx context.Context,
	prefix uint16,
	rangeSize int,
) ([]*domain.Image, error) {

	low := int(prefix) - rangeSize
	if low < 0 {
		low = 0
	}
	high := int(prefix) + rangeSize

	rows, err := r.db.Query(ctx, `
		SELECT id, url, dhash
		FROM images
		WHERE dhash_prefix BETWEEN $1 AND $2
	`, low, high)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.Image
	for rows.Next() {
		var img domain.Image
		if err := rows.Scan(&img.ID, &img.URL, &img.DHash); err != nil {
			return nil, err
		}
		images = append(images, &img)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

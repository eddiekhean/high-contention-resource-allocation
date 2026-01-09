package repository

import (
	"context"

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
	err := r.db.QueryRow(ctx, `
		INSERT INTO images (url, dhash)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, img.URL, img.DHash).Scan(&img.ID, &img.CreatedAt)

	return err
}

func (r *PgImageRepository) FindByID(
	ctx context.Context,
	id int64,
) (*domain.Image, error) {
	var img domain.Image
	err := r.db.QueryRow(ctx, `
		SELECT id, url, dhash, created_at
		FROM images
		WHERE id = $1
	`, id).Scan(&img.ID, &img.URL, &img.DHash, &img.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *PgImageRepository) FindByDHash(ctx context.Context, hash uint64) (*domain.Image, error) {
	var img domain.Image
	err := r.db.QueryRow(ctx, `
		SELECT id, url, dhash, created_at
		FROM images
		WHERE dhash = $1
	`, hash).Scan(&img.ID, &img.URL, &img.DHash, &img.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *PgImageRepository) GetAll(ctx context.Context) ([]*domain.Image, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, url, dhash, created_at
		FROM images
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.Image
	for rows.Next() {
		var img domain.Image
		if err := rows.Scan(&img.ID, &img.URL, &img.DHash, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, &img)
	}
	return images, nil
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

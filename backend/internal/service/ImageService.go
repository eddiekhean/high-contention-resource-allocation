package service

import (
	"context"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/repository"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
)

type ImageService interface {
	CreateImage(ctx context.Context, img *domain.Image) error
	Match(ctx context.Context, dhashStr string) (*domain.Image, int, bool, error)
}

type imageService struct {
	repo repository.ImageRepository
	cfg  *config.ImageConfig
}

func NewImageService(repo repository.ImageRepository, cfg *config.ImageConfig) ImageService {
	return &imageService{repo: repo, cfg: cfg}
}

func (s *imageService) CreateImage(ctx context.Context, img *domain.Image) error {
	return s.repo.Create(ctx, img)
}
func (s *imageService) Match(
	ctx context.Context,
	dhashStr string,
) (*domain.Image, int, bool, error) {

	// 1. Parse dHash
	target, err := utils.ParseDHash(dhashStr)
	if err != nil {
		return nil, 0, false, err
	}

	// 2. Compute prefix (16 bit đầu)
	prefix := uint16(target >> 48)

	// 3. Query candidate images
	candidates, err := s.repo.FindByPrefix(
		ctx,
		prefix,
		s.cfg.PrefixRange,
	)
	if err != nil {
		return nil, 0, false, err
	}

	if len(candidates) == 0 {
		return nil, 0, false, nil
	}

	// 4. Compute Hamming distance
	bestDist := 64
	var best *domain.Image

	for _, img := range candidates {
		dist := utils.HammingDistance(target, img.DHash)
		if dist < bestDist {
			bestDist = dist
			best = img
		}
	}

	// 5. Check threshold
	if best != nil && bestDist <= s.cfg.MatchThreshold {
		return best, bestDist, true, nil
	}

	return nil, 0, false, nil
}

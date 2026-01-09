package service

import (
	"context"
	"fmt"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/client"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/repository"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/google/uuid"
)

type ImageService interface {
	CreateImage(ctx context.Context, img *domain.Image) error
	Match(ctx context.Context, dhashStr string) (*domain.Image, int, bool, error) // keeping for existing use
	MatchImage(ctx context.Context, dhash uint64) (*domain.Image, bool, error)
	GetUploadURL(ctx context.Context) (string, string, error)
	CommitImage(ctx context.Context, dhash uint64, s3Key string) (*domain.Image, error)
}

type imageService struct {
	repo     repository.ImageRepository
	cfg      *config.ImageConfig
	s3Client *client.S3Client
}

func NewImageService(repo repository.ImageRepository, cfg *config.ImageConfig, s3Client *client.S3Client) ImageService {
	return &imageService{repo: repo, cfg: cfg, s3Client: s3Client}
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

func (s *imageService) MatchImage(ctx context.Context, target uint64) (*domain.Image, bool, error) {
	// 1. Get all images from DB (stateless similarity search)
	// For production with many images, we should use a more efficient search (like VP-Tree or BK-Tree or SQL optimization)
	// but here we follow the requirement to compute Hamming distance in backend.
	images, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, false, err
	}

	bestDist := 64
	var best *domain.Image

	threshold := 5 // specified requirement

	for _, img := range images {
		dist := utils.HammingDistance(target, img.DHash)
		if dist < bestDist {
			bestDist = dist
			best = img
		}
	}

	if best != nil && bestDist <= threshold {
		// Generate presigned GET URL
		signedURL, err := s.s3Client.GetPresignedURL(ctx, best.URL, 10*time.Minute)
		if err != nil {
			return nil, false, err
		}

		// Return a copy with the signed URL
		result := *best
		result.URL = signedURL
		return &result, true, nil
	}

	return nil, false, nil
}

func (s *imageService) GetUploadURL(ctx context.Context) (string, string, error) {
	s3Key := fmt.Sprintf("images/%s.jpg", uuid.New().String())
	uploadURL, err := s.s3Client.GetPresignedPutURL(ctx, s3Key, 15*time.Minute)
	if err != nil {
		return "", "", err
	}
	return uploadURL, s3Key, nil
}

func (s *imageService) CommitImage(ctx context.Context, dhash uint64, s3Key string) (*domain.Image, error) {
	// 1. Verify object exists in S3
	exists, err := s.s3Client.Exists(ctx, s3Key)
	if err != nil || !exists {
		return nil, fmt.Errorf("object does not exist in S3 or check failed: %w", err)
	}

	// 2. Race condition check: re-verify similarity before insert
	match, matched, err := s.MatchImage(ctx, dhash)
	if err == nil && matched {
		return nil, fmt.Errorf("duplicate image detected during commit: id=%d", match.ID)
	}

	// 3. Insert into DB
	img := &domain.Image{
		URL:   s3Key,
		DHash: dhash,
	}

	if err := s.repo.Create(ctx, img); err != nil {
		return nil, err
	}

	return img, nil
}

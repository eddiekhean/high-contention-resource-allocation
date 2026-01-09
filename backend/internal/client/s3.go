package client

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	cfg "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

// NewS3Client khởi tạo S3 client từ app config
func NewS3Client(c *cfg.S3Config) (*S3Client, error) {
	if !c.Enabled {
		return nil, nil
	}

	if c.Addr == "" {
		return nil, errors.New("s3 region (addr) is empty")
	}

	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(c.Addr),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &S3Client{
		Client: client,
		Bucket: c.Bucket,
	}, nil
}

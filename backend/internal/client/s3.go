package client

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	cfg "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
)

type S3Client struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}

// NewS3Client khởi tạo S3 client từ app config
func NewS3Client(c *cfg.S3Config) (*S3Client, error) {
	if !c.Enabled {
		return nil, nil
	}

	if c.Addr == "" {
		return nil, errors.New("s3 region (addr) is empty")
	}

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(c.Addr),
	}

	if c.AccessKey != "" && c.SecretKey != "" {
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, ""),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if c.Endpoint != "" {
			o.BaseEndpoint = aws.String(c.Endpoint)
		}
	})

	return &S3Client{
		Client:        client,
		PresignClient: s3.NewPresignClient(client),
		Bucket:        c.Bucket,
	}, nil
}

// GetPresignedURL sinh URL để xem ảnh
func (s *S3Client) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	req, err := s.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// GetPresignedPutURL sinh URL để upload ảnh
func (s *S3Client) GetPresignedPutURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	req, err := s.PresignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiration))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// Exists kiểm tra file có tồn tại trên S3 không
func (s *S3Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, nil
		}
		// Some S3 implementations might return NotFound for HeadObject
		var nf *types.NotFound
		if errors.As(err, &nf) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Delete xóa file trên S3
func (s *S3Client) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	return err
}

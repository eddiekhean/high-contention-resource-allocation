package service

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/sirupsen/logrus"
)

type S3Service struct {
	logger *logrus.Logger
	cfg    *config.S3Config
	// repo   repository.ImageRepository
	client *s3.Client
}

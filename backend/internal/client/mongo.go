package client

import (
	"context"
	"fmt"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewMongoClient creates a MongoDB client from config.
// Returns nil, nil when URI is empty (Mongo disabled).
func NewMongoClient(cfg *config.MongoConfig) (*mongo.Client, error) {
	if cfg.URI == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(cfg.URI)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	return client, nil
}

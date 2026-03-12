package client

import (
	"context"
	"fmt"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool creates a pgx connection pool from config.
// Returns nil, nil when DSN is empty (Postgres disabled).
func NewPostgresPool(cfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	if cfg.DSN == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres DSN: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}

	return pool, nil
}

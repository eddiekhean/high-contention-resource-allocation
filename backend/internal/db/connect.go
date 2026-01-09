package db

import (
	"context"
	"fmt"
	"time"

	"os"

	cfg "github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

// NewPostgres tạo kết nối DB từ config
func NewPostgres(cfg *cfg.PostgresConfig) (*Postgres, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("postgres addr is empty")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Addr,
		cfg.DB,
	)

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	// config pool cơ bản (đủ cho đồ án)
	poolCfg.MaxConns = 10
	poolCfg.MinConns = 2
	poolCfg.MaxConnLifetime = 30 * time.Minute
	poolCfg.MaxConnIdleTime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	// Ping DB để chắc chắn kết nối OK
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return &Postgres{
		Pool: pool,
	}, nil
}

// Close đóng pool khi shutdown app
func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

// Exec thực thi query không trả về rows (INSERT, UPDATE, DELETE, etc.)
func (p *Postgres) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return p.Pool.Exec(ctx, sql, arguments...)
}

// Query thực thi query trả về nhiều rows
func (p *Postgres) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.Pool.Query(ctx, sql, args...)
}

// QueryRow thực thi query trả về 1 row duy nhất
func (p *Postgres) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.Pool.QueryRow(ctx, sql, args...)
}
func RunMigration(p *Postgres) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	sqlBytes, err := os.ReadFile("internal/db/001_create_images.sql")
	if err != nil {
		return err
	}

	_, err = p.Exec(ctx, string(sqlBytes))
	return err
}

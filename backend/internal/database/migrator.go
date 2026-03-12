// Package database provides auto-migration support using golang-migrate.
// SQL files are embedded into the binary so no external files are needed at runtime.
package database

import (
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// MigrateUp runs all pending UP migrations against the given pgx pool.
// Safe to call on every server startup — already-applied migrations are skipped.
func MigrateUp(pool *pgxpool.Pool) error {
	// Wrap pgxpool as database/sql for golang-migrate compatibility
	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("migrate: failed to create postgres driver: %w", err)
	}

	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("migrate: failed to load migration source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return fmt.Errorf("migrate: failed to init migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to run migrations: %w", err)
	}

	return nil
}

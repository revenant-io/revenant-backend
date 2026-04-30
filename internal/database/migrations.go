package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/revenantio/revenant-backend/internal/config"
)

func RunMigrations(db *sql.DB, cfg config.DatabaseConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		// Ignore if migrations directory doesn't exist yet
		if err == migrate.ErrNilVersion {
			return nil
		}
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

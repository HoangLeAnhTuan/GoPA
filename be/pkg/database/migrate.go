package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"gopa/migrations"
)

type MigrationDirection string

const (
	MigrationUp   MigrationDirection = "up"
	MigrationDown MigrationDirection = "down"
)

func Migrate(db *sqlx.DB, direction MigrationDirection, steps int) error {
	if steps < 1 {
		return fmt.Errorf("migration steps must be at least one")
	}

	source, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("create migration source: %w", err)
	}
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres migration driver: %w", err)
	}
	runner, err := migrate.NewWithInstance("iofs", source, "gopa", driver)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer runner.Close()

	switch direction {
	case MigrationUp:
		err = runner.Up()
	case MigrationDown:
		err = runner.Steps(-steps)
	default:
		return fmt.Errorf("unsupported migration direction %q", direction)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

package postgres

import (
	"context"
	"embed"
	"errors"

	"github.com/018bf/example/internal/configs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

func NewDatabase(config *configs.Config) (*sqlx.DB, error) {
	database, err := sqlx.Connect("postgres", config.Database.URI)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(config.Database.MaxOpenConnections)
	database.SetMaxIdleConns(config.Database.MaxIDLEConnections)
	return database, nil
}

type MigrateManager struct {
	database *sqlx.DB
	config   *configs.Config
}

func NewMigrateManager(database *sqlx.DB, config *configs.Config) *MigrateManager {
	return &MigrateManager{
		database: database,
		config:   config,
	}
}

func (m MigrateManager) Up(_ context.Context) error {
	source, err := iofs.New(MigrationsFS, "migrations")
	if err != nil {
		return err
	}
	instance, err := migrate.NewWithSourceInstance("iofs", source, m.config.Database.URI)
	if err != nil {
		return err
	}
	if err := instance.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

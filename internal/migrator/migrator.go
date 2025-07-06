package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"go-monolite/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var MigrationsFS embed.FS

const migrationsDir = "migrations"

var (
	up   = "up"
	down = "down"
)

type Migrator struct {
	moduleDriver source.Driver
	databaseName string
	db           *sql.DB
}

func MustGetNewMigrator(databaseName string, db *sql.DB) *Migrator {
	d, err := iofs.New(MigrationsFS, migrationsDir)
	if err != nil {
		panic(err)
	}
	return &Migrator{
		moduleDriver: d,
		databaseName: databaseName,
		db:           db,
	}
}

func (m *Migrator) Up() error {
	migrator, err := m.newMigrator(m.db)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	defer m.closeQuietly(migrator)

	if err = migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no new migrations to apply:" + migrate.ErrNoChange.Error())
			return nil
		}
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}

func (m *Migrator) Down() error {
	migrator, err := m.newMigrator(m.db)
	if err != nil {
		return fmt.Errorf("unable to create migration: %v", err)
	}

	defer m.closeQuietly(migrator)

	if err = migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("no new migrations to apply:" + migrate.ErrNoChange.Error())
			return nil
		}
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	return nil
}

func (m *Migrator) newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to create db instance: %w", err)
	}

	return migrate.NewWithInstance("migration_embeded_sql_files", m.moduleDriver, m.databaseName, driver)
}

func (m *Migrator) closeQuietly(migrator *migrate.Migrate) {
	moduleErr, dbErr := migrator.Close()
	if moduleErr != nil {
		logger.Warn(moduleErr, "error closing migration source")
	}
	if dbErr != nil {
		logger.Warn(dbErr, "error closing migration database")
	}
}

func (m *Migrator) Run(action string) (string, error) {
	if action != up && action != down {
		return "", fmt.Errorf("invalid action: %s, must be 'up' or 'down'", action)
	}

	switch action {
	case up:
		if err := m.Up(); err != nil {
			return "", fmt.Errorf("error applying migrations: %v", err)
		}
		return "migration has Up successfully", nil
	case down:
		if err := m.Down(); err != nil {
			return "", fmt.Errorf("error rolling back migrations: %v", err)
		}
		return "migration has Down successfully", nil
	default:
		return "", fmt.Errorf("unknown action: %s, use 'up' or 'down'", action)
	}
}

// base := "internal/module"

// err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		return err
// 	}

// 	if !info.IsDir() && strings.HasSuffix(path, ".sql") && strings.Contains(path, "/migrations/") {
// 		fmt.Println("Нашёл миграцию:", path)
// 	}

// 	return nil
// })

// if err != nil {
// 	panic(err)
// }

// --------------------------------------------

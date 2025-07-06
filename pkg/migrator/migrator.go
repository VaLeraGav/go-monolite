package migrator

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	migrationsDir string
	dbName        string
	db            *sql.DB
	migrate       *migrate.Migrate
	tempDir       string
}

func New(migrationsDir, dbName string, db *sql.DB) (*Migrator, error) {
	tempDir, err := os.MkdirTemp("", "merged-migrations-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}

	err = filepath.Walk(migrationsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		normPath := filepath.ToSlash(path)
		if !info.IsDir() && strings.HasSuffix(normPath, ".sql") && strings.Contains(normPath, "/migrations/") {
			dest := filepath.Join(tempDir, filepath.Base(path))
			return copyFile(path, dest)
		}
		return nil
	})
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to copy migration files: %w", err)
	}

	// Создаём виртуальную FS из временной папки
	migrationFS := os.DirFS(tempDir)
	sourceDriver, err := iofs.New(migrationFS, ".")
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create iofs source: %w", err)
	}

	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, dbName, dbDriver)
	if err != nil {
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("failed to create migrator instance: %w", err)
	}

	return &Migrator{
		migrationsDir: migrationsDir,
		dbName:        dbName,
		db:            db,
		migrate:       migrator,
		tempDir:       tempDir,
	}, nil
}

// Run performs Up and Down
func (m *Migrator) Run(action string) (string, error) {
	defer m.Close()

	var err error
	switch action {
	case "up":
		err = m.migrate.Up()
	case "down":
		err = m.migrate.Down()
	default:
		return "", fmt.Errorf("unsupported action %q", action)
	}

	if err != nil && err != migrate.ErrNoChange {
		return "", fmt.Errorf("error migrate: %w", err)
	}

	return fmt.Sprintf("Migration %q completed successfully", action), nil
}

func (m *Migrator) CreateMigration(modulePath, name string) error {
	absModulePath, err := m.validateModuleDir(modulePath)
	if err != nil {
		return err
	}

	migrationsPath, err := m.prepareMigrationsDir(absModulePath)
	if err != nil {
		return err
	}

	filenamePrefix := generateMigrationFilenamePrefix(name)
	upPath := filepath.Join(migrationsPath, filenamePrefix+".up.sql")
	downPath := filepath.Join(migrationsPath, filenamePrefix+".down.sql")

	if err := writeMigrationFile(upPath, "-- Write your UP migration here\n"); err != nil {
		return fmt.Errorf("failed to write up migration: %w", err)
	}

	if err := writeMigrationFile(downPath, "-- Write your DOWN migration here\n"); err != nil {
		return fmt.Errorf("failed to write down migration: %w", err)
	}

	fmt.Printf("Created migration:\n  %s\n  %s\n", upPath, downPath)
	return nil
}

func (m *Migrator) Version() error {
	version, dirty, err := m.migrate.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			return fmt.Errorf("migrations have not been applied yet, the version is missing: %w", err)
		}
		return fmt.Errorf("failed get Version: %w", err)
	}
	fmt.Printf("Current migration version: %d, dirty: %v\n", version, dirty)
	return nil
}

func (m *Migrator) Close() error {
	var errRemove error
	if m.tempDir != "" {
		errRemove = os.RemoveAll(m.tempDir)
	}

	srcErr, dbErr := m.migrate.Close()

	var combinedErr error
	if errRemove != nil {
		combinedErr = fmt.Errorf("remove temp dir error: %w", errRemove)
	}
	if srcErr != nil {
		combinedErr = fmt.Errorf("%w; migrate source close error: %v", combinedErr, srcErr)
	}
	if dbErr != nil {
		combinedErr = fmt.Errorf("%w; migrate database close error: %v", combinedErr, dbErr)
	}

	return combinedErr
}

func (m *Migrator) validateModuleDir(modulePath string) (string, error) {
	absModulePath := filepath.Join(m.migrationsDir, modulePath)

	info, err := os.Stat(absModulePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("module directory does not exist: %s", absModulePath)
		}
		return "", fmt.Errorf("failed to stat module directory: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("expected directory but found file: %s", absModulePath)
	}

	return absModulePath, nil
}

func (m *Migrator) prepareMigrationsDir(absModulePath string) (string, error) {
	migrationsPath := filepath.Join(absModulePath, "migrations")
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create migrations dir: %w", err)
	}
	return migrationsPath, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func generateMigrationFilenamePrefix(name string) string {
	timestamp := time.Now().Format("20060102_150405")
	safeName := strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("%s__%s__", timestamp, safeName)
}

func writeMigrationFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

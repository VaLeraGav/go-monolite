package testinit

import (
	"fmt"
	"go-monolite/internal/config"
	"go-monolite/internal/store"
	"go-monolite/pkg/logger"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

type Store any

type Handler interface {
	Init(r chi.Router)
}

func SetupStoreTest(t *testing.T) *store.Store {
	t.Helper()

	config := GetConfigs()
	testStore := NewStoreTest(t, config)
	logger.InitLogger(config.Env, "")

	return testStore
}

func GetConfigs() *config.Config {
	configPath := config.GetConfigPathFromTest(".env." + logger.EnvTest)
	configs := config.MustInit(configPath)

	return configs
}

func SetupTestServer[T Handler](t *testing.T, handler T) *httptest.Server {
	t.Helper()

	router := chi.NewRouter()
	handler.Init(router)

	server := httptest.NewServer(router)

	return server
}

func NewStoreTest(t *testing.T, config *config.Config) *store.Store {
	storeTest, err := store.NewPostgres(config)
	if err != nil {
		t.Fatalf("connectDb error: %v", err)
	}
	return storeTest
}

// очистка всех таблиц
func TruncateAllTables(db *sqlx.DB) error {
	var tableNames []string

	query := `
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
	`

	if err := db.Select(&tableNames, query); err != nil {
		return fmt.Errorf("failed to query table names: %w", err)
	}

	if len(tableNames) == 0 {
		return nil
	}

	tables := strings.Join(tableNames, ", ")
	truncateQuery := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables)

	if _, err := db.Exec(truncateQuery); err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}

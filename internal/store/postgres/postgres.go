package postgres

import (
	"context"
	"fmt"
	"go-monolite/internal/config"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // Импорт драйвера pgx
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL
)

func Connect(сfgEnv *config.Config) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connString := buildDbConnectUrl(сfgEnv)

	conn, err := sqlx.ConnectContext(ctx, "postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return conn, nil
}

func buildDbConnectUrl(сfgEnv *config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		сfgEnv.Db.Driver,
		сfgEnv.Db.User,
		сfgEnv.Db.Password,
		сfgEnv.Db.Host,
		сfgEnv.Db.ExternalPort,
		сfgEnv.Db.NameDb,
		сfgEnv.Db.Option,
	)
}

package store

import (
	"context"
	"errors"
	"go-monolite/internal/config"
	"go-monolite/internal/store/postgres"
	"go-monolite/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	Db *sqlx.DB
}

func NewPostgres(сfgEnv *config.Config) (*Store, error) {
	dbPgx, err := postgres.Connect(сfgEnv)
	if err != nil {
		return nil, err
	}
	return &Store{
		Db: dbPgx,
	}, nil
}

func (s *Store) Shutdown(ctx context.Context) error {
	logger.Info("closing database connection")
	return s.Db.Close()
}

func ContextError(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrTimeoutExceeded
	}
	if errors.Is(err, context.Canceled) {
		return ErrOperationCanceled
	}
	return err
}

type txKey struct{}

func WithTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

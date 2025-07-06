package storage

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"time"
)

type StorageRepository struct {
	store     *store.Store
	tableName string
}

func NewStorageRepository(store *store.Store) *StorageRepository {
	return &StorageRepository{
		store:     store,
		tableName: "storage",
	}
}

func (r *StorageRepository) Create(ctx context.Context, s *StorageEnt) (*uint, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			uuid, name, active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id
	`, r.tableName)

	now := time.Now()
	s.CreatedAt = now
	s.UpdatedAt = now

	var id uint
	if tx := store.GetTx(ctx); tx != nil {
		err := tx.QueryRowxContext(ctx, query,
			s.UUID,
			s.Name,
			s.Active,
			s.CreatedAt,
			s.UpdatedAt,
		).Scan(&id)
		if err != nil {
			return nil, store.ContextError(err)
		}
	} else {
		err := r.store.Db.QueryRowxContext(ctx, query,
			s.UUID,
			s.Name,
			s.Active,
			s.CreatedAt,
			s.UpdatedAt,
		).Scan(&id)
		if err != nil {
			return nil, store.ContextError(err)
		}
	}

	s.ID = id
	return &id, nil
}

func (r *StorageRepository) GetList(ctx context.Context) ([]StorageEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, name, active, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`, r.tableName)

	var storages []StorageEnt
	var err error
	if tx := store.GetTx(ctx); tx != nil {
		err = tx.SelectContext(ctx, &storages, query)
	} else {
		err = r.store.Db.SelectContext(ctx, &storages, query)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	if len(storages) == 0 {
		return nil, store.ErrNotFound
	}

	return storages, nil
}

func (r *StorageRepository) Update(ctx context.Context, s *StorageEnt) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, active = $2, updated_at = $3
		WHERE uuid = $4
	`, r.tableName)

	s.UpdatedAt = time.Now()

	var result sql.Result
	var err error
	if tx := store.GetTx(ctx); tx != nil {
		result, err = tx.ExecContext(ctx, query,
			s.Name,
			s.Active,
			s.UpdatedAt,
			s.UUID,
		)
	} else {
		result, err = r.store.Db.ExecContext(ctx, query,
			s.Name,
			s.Active,
			s.UpdatedAt,
			s.UUID,
		)
	}
	if err != nil {
		return store.ContextError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return store.ContextError(err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

func (r *StorageRepository) Delete(ctx context.Context, uuid string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = $1`, r.tableName)

	var result sql.Result
	var err error
	if tx := store.GetTx(ctx); tx != nil {
		result, err = tx.ExecContext(ctx, query, uuid)
	} else {
		result, err = r.store.Db.ExecContext(ctx, query, uuid)
	}
	if err != nil {
		return store.ContextError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return store.ContextError(err)
	}
	if rows == 0 {
		return store.ErrNotFound
	}

	return nil
}

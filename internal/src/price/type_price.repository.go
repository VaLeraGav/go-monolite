package price

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"time"
)

type TypePriceRepository struct {
	store     *store.Store
	tableName string
}

func NewTypePriceRepository(store *store.Store) *TypePriceRepository {
	return &TypePriceRepository{
		store:     store,
		tableName: "type_price",
	}
}

func (r *TypePriceRepository) Create(ctx context.Context, p *TypePriceEnt) (*uint, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			uuid, name, active, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id
	`, r.tableName)

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	var id uint
	err := r.store.Db.QueryRowxContext(ctx, query,
		p.UUID,
		p.Name,
		p.Active,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, store.ContextError(err)
	}

	p.ID = id
	return &id, nil
}

func (r *TypePriceRepository) GetList(ctx context.Context) ([]TypePriceEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, name, active, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`, r.tableName)

	var prices []TypePriceEnt
	err := r.store.Db.SelectContext(ctx, &prices, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	if len(prices) == 0 {
		return nil, store.ErrNotFound
	}

	return prices, nil
}

func (r *TypePriceRepository) Update(ctx context.Context, p *TypePriceEnt) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, active = $2, updated_at = $3
		WHERE uuid = $4
	`, r.tableName)

	p.UpdatedAt = time.Now()

	result, err := r.store.Db.ExecContext(ctx, query,
		p.Name,
		p.Active,
		p.UpdatedAt,
		p.UUID,
	)
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

func (r *TypePriceRepository) Delete(ctx context.Context, uuid string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = $1`, r.tableName)

	result, err := r.store.Db.ExecContext(ctx, query, uuid)
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

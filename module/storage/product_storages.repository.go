package storage

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductStoragesRepository struct {
	store     *store.Store
	tableName string
}

func NewProductStoragesRepository(store *store.Store) *ProductStoragesRepository {
	return &ProductStoragesRepository{
		store:     store,
		tableName: "product_storages",
	}
}

func (r *ProductStoragesRepository) GetByProductUUIDs(ctx context.Context, productUUIDs []uuid.UUID) ([]ProductStorageEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, product_uuid, storage_uuid, active, quantity, created_at, updated_at
			FROM %s
			WHERE product_uuid IN (?)
		ORDER BY created_at DESC
	`, r.tableName)

	query, args, err := sqlx.In(query, productUUIDs)
	if err != nil {
		return nil, store.ContextError(err)
	}

	query = r.store.Db.Rebind(query)

	var storages []ProductStorageEnt
	if tx := store.GetTx(ctx); tx != nil {
		err = tx.SelectContext(ctx, &storages, query, args...)
	} else {
		err = r.store.Db.SelectContext(ctx, &storages, query, args...)
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

func (r *ProductStoragesRepository) CreateBatch(ctx context.Context, records []ProductStorageEnt) error {
	if len(records) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			product_uuid, storage_uuid, active, quantity, created_at, updated_at
		) VALUES 
	`, r.tableName)

	args := []any{}
	now := time.Now()

	valueStrings := make([]string, 0, len(records))
	for i, rec := range records {
		rec.CreatedAt = now
		rec.UpdatedAt = now

		args = append(args,
			rec.ProductUUID,
			rec.StorageUUID,
			rec.Active,
			rec.Quantity,
			rec.CreatedAt,
			rec.UpdatedAt,
		)

		start := i*6 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)",
			start, start+1, start+2, start+3, start+4, start+5))
	}

	query += strings.Join(valueStrings, ", ")

	var err error
	if tx := store.GetTx(ctx); tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.store.Db.ExecContext(ctx, query, args...)
	}
	if err != nil {
		return store.ContextError(err)
	}

	return nil
}

func (r *ProductStoragesRepository) UpdateBatch(ctx context.Context, records []ProductStorageEnt) error {
	if len(records) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`
		UPDATE %s AS ps SET
			active = v.active,
			quantity = v.quantity,
			updated_at = v.updated_at
		FROM (VALUES
	`, r.tableName))

	now := time.Now().UTC()
	for i, rec := range records {
		builder.WriteString(fmt.Sprintf(
			"('%s', '%s', '%s', %d, TIMESTAMP '%s')",
			rec.ProductUUID.String(),
			rec.StorageUUID.String(),
			rec.Active,
			rec.Quantity,
			now.Format("2006-01-02 15:04:05.999999"),
		))

		if i < len(records)-1 {
			builder.WriteString(",\n")
		}
	}

	builder.WriteString(`
		) AS v(product_uuid, storage_uuid, active, quantity, updated_at)
		WHERE ps.product_uuid = v.product_uuid::uuid
		AND ps.storage_uuid = v.storage_uuid::uuid;
	`)

	query := builder.String()
	var err error
	if tx := store.GetTx(ctx); tx != nil {
		_, err = tx.ExecContext(ctx, query)
	} else {
		_, err = r.store.Db.ExecContext(ctx, query)
	}
	if err != nil {
		return fmt.Errorf("bulk update %s failed: %w", r.tableName, err)
	}

	return nil
}

func (r *ProductStoragesRepository) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = ANY($1)
	`, r.tableName)

	var result sql.Result
	var err error
	if tx := store.GetTx(ctx); tx != nil {
		result, err = tx.ExecContext(ctx, query, pq.Array(ids))
	} else {
		result, err = r.store.Db.ExecContext(ctx, query, pq.Array(ids))
	}
	if err != nil {
		return store.ContextError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return store.ContextError(err)
	}

	if int(rows) != len(ids) {
		return fmt.Errorf("deleted %d of %d %s: some ids not found", rows, len(ids), r.tableName)
	}

	return nil
}

package price

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

type ProductPricesRepository struct {
	store     *store.Store
	tableName string
}

func NewProductPricesRepository(store *store.Store) *ProductPricesRepository {
	return &ProductPricesRepository{
		store:     store,
		tableName: "product_prices",
	}
}

func (r *ProductPricesRepository) GetByProductUUIDs(ctx context.Context, productUUIDs []uuid.UUID) ([]ProductPriceEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, product_uuid, type_price_uuid, active, price, created_at, updated_at
			FROM %s
			WHERE product_uuid IN (?)
		ORDER BY created_at DESC
	`, r.tableName)

	query, args, err := sqlx.In(query, productUUIDs)
	if err != nil {
		return nil, store.ContextError(err)
	}

	query = r.store.Db.Rebind(query)

	var prices []ProductPriceEnt
	err = r.store.Db.SelectContext(ctx, &prices, query, args...)
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

func (r *ProductPricesRepository) CreateBatch(ctx context.Context, records []ProductPriceEnt) error {
	if len(records) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			product_uuid, type_price_uuid, active, price, created_at, updated_at
		) VALUES 
	`, r.tableName)

	args := []interface{}{}
	now := time.Now()

	valueStrings := make([]string, 0, len(records))
	for i, rec := range records {
		rec.CreatedAt = now
		rec.UpdatedAt = now

		args = append(args,
			rec.ProductUUID,
			rec.TypePriceUUID,
			rec.Active,
			rec.Price,
			rec.CreatedAt,
			rec.UpdatedAt,
		)

		start := i*6 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)",
			start, start+1, start+2, start+3, start+4, start+5))
	}

	query += strings.Join(valueStrings, ", ")

	_, err := r.store.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return store.ContextError(err)
	}

	return nil
}

func (r *ProductPricesRepository) UpdateBatch(ctx context.Context, records []ProductPriceEnt) error {
	if len(records) == 0 {
		return nil
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`
		UPDATE %s AS pp SET
			active = v.active,
			price = v.price,
			updated_at = v.updated_at
		FROM (VALUES
	`, r.tableName))

	now := time.Now().UTC()
	for i, rec := range records {
		builder.WriteString(fmt.Sprintf(
			"('%s', '%s', '%s', %f, TIMESTAMP '%s')",
			rec.ProductUUID.String(),
			rec.TypePriceUUID.String(),
			rec.Active,
			rec.Price,
			now.Format("2006-01-02 15:04:05.999999"),
		))

		if i < len(records)-1 {
			builder.WriteString(",\n")
		}
	}

	builder.WriteString(`
		) AS v(product_uuid, type_price_uuid, active, price, updated_at)
		WHERE pp.product_uuid = v.product_uuid::uuid
		AND pp.type_price_uuid = v.type_price_uuid::uuid;
	`)

	query := builder.String()
	_, err := r.store.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("bulk update %s failed: %w", r.tableName, err)
	}

	return nil
}

func (r *ProductPricesRepository) DeleteBatch(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE id = ANY($1)
	`, r.tableName)

	result, err := r.store.Db.ExecContext(ctx, query, pq.Array(ids))
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

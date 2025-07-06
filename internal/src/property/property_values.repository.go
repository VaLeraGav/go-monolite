package property

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PropertyValuesRepository struct {
	store     *store.Store
	tableName string
}

func NewPropertyValuesRepository(store *store.Store) *PropertyValuesRepository {
	return &PropertyValuesRepository{
		store:     store,
		tableName: "property_values",
	}
}

func (r *PropertyValuesRepository) CreateBatch(ctx context.Context, values []PropertyValueEnt) error {
	if len(values) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (key, slug, value, property_uuid)
		VALUES
	`, r.tableName)

	args := make([]any, 0, len(values)*4)
	valueStrings := make([]string, 0, len(values))

	for i, v := range values {
		start := i*4 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", start, start+1, start+2, start+3))
		args = append(args, v.Key, v.Slug, v.Value, v.PropertyUUID)
	}

	query += strings.Join(valueStrings, ", ")

	_, err := r.store.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return store.ContextError(err)
	}

	return nil
}

func (r *PropertyValuesRepository) GetList(ctx context.Context) ([]PropertyValueEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, key, slug, property_uuid, value
		FROM %s
	`, r.tableName)

	var properties []PropertyValueEnt
	err := r.store.Db.SelectContext(ctx, &properties, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	if len(properties) == 0 {
		return nil, store.ErrNotFound
	}

	return properties, nil
}

func (r *PropertyValuesRepository) GetByUUID(ctx context.Context, propertyUUID uuid.UUID) ([]PropertyValueEnt, error) {
	valuesQuery := fmt.Sprintf(`
		SELECT id, key, slug, property_uuid, value
		FROM %s
		WHERE property_uuid = $1
	`, r.tableName)

	var values []PropertyValueEnt
	err := r.store.Db.SelectContext(ctx, &values, valuesQuery, propertyUUID)
	if err != nil {
		return nil, store.ContextError(err)
	}
	if len(values) == 0 {
		return nil, store.ErrNotFound
	}
	return values, nil
}

func (r *PropertyValuesRepository) UpdateBatch(ctx context.Context, values []PropertyValueEnt) error {
	if len(values) == 0 {
		return nil
	}

	var (
		updateErrs []string
		now        = time.Now()
	)

	for _, v := range values {
		query := fmt.Sprintf(`
			UPDATE %s
			SET slug = $1, value = $2, updated_at = $3
			WHERE key = $4 AND property_uuid = $5
		`, r.tableName)

		result, err := r.store.Db.ExecContext(ctx, query,
			v.Slug,
			v.Value,
			now,
			v.Key,
			v.PropertyUUID,
		)
		if err != nil {
			updateErrs = append(updateErrs, fmt.Sprintf("key %s: %v", v.Key, err))
			continue
		}

		rows, err := result.RowsAffected()
		if err != nil {
			updateErrs = append(updateErrs, fmt.Sprintf("key %s: rowsAffected error: %v", v.Key, err))
			continue
		}

		if rows == 0 {
			updateErrs = append(updateErrs, fmt.Sprintf("key %s not found", v.Key))
		}
	}

	if len(updateErrs) > 0 {
		return fmt.Errorf("some updates failed:\n%s", strings.Join(updateErrs, "\n"))
	}

	return nil
}

func (r *PropertyValuesRepository) Delete(ctx context.Context, uuid uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE property_uuid = $1`, r.tableName)
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

func (r *PropertyValuesRepository) DeleteBatch(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE uuid = ANY($1)
	`, r.tableName)

	result, err := r.store.Db.ExecContext(ctx, query, pq.Array(keys))
	if err != nil {
		return store.ContextError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return store.ContextError(err)
	}

	if int(rows) != len(keys) {
		return fmt.Errorf("deleted %d of %d properties: some uuids not found", rows, len(keys))
	}

	return nil
}

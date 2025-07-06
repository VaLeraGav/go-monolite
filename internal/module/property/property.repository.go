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

type PropertyRepository struct {
	store     *store.Store
	tableName string
}

func NewPropertyRepository(store *store.Store) *PropertyRepository {
	return &PropertyRepository{
		store:     store,
		tableName: "property",
	}
}

func (r *PropertyRepository) CreateBatch(ctx context.Context, props []PropertyEnt) error {
	if len(props) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			uuid, slug, type, name, created_at, updated_at
		) VALUES 
	`, r.tableName)

	args := []interface{}{}
	now := time.Now()

	valueStrings := make([]string, 0, len(props))
	for i, p := range props {
		p.CreatedAt = now
		p.UpdatedAt = now

		args = append(args,
			p.UUID,
			p.Slug,
			p.Type,
			p.Name,
			p.CreatedAt,
			p.UpdatedAt,
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

func (r *PropertyRepository) GetList(ctx context.Context) ([]PropertyEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, slug, type, name, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`, r.tableName)

	var properties []PropertyEnt
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

func (r *PropertyRepository) GetByUUID(ctx context.Context, uuid string) (*PropertyEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, slug, type, name, created_at, updated_at
		FROM %s
		WHERE uuid = $1
	`, r.tableName)

	var property PropertyEnt
	err := r.store.Db.GetContext(ctx, &property, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	return &property, nil
}

func (r *PropertyRepository) UpdateBatch(ctx context.Context, props []PropertyEnt) error {
	if len(props) == 0 {
		return nil
	}

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf(`
	UPDATE %s AS p SET
		slug = v.slug,
		type = v.type,
		name = v.name,
		updated_at = v.updated_at
	FROM (VALUES
	`, r.tableName))

	// escape для безопасного построения запроса
	escape := func(s string) string {
		return strings.ReplaceAll(s, "'", "''")
	}

	now := time.Now().UTC()

	for i, p := range props {
		builder.WriteString(fmt.Sprintf(
			"('%s', '%s', '%s', '%s', TIMESTAMP '%s')",
			p.UUID.String(),
			escape(p.Slug),
			escape(p.Type),
			escape(p.Name),
			now.Format("2006-01-02 15:04:05.999999"),
		))

		if i != len(props)-1 {
			builder.WriteString(",\n")
		}
	}

	builder.WriteString(`
) AS v(uuid, slug, type, name, updated_at)
WHERE p.uuid = v.uuid::uuid;
`)

	query := builder.String()

	_, err := r.store.Db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("bulk update property failed: %w", err)
	}

	return nil
}

func (r *PropertyRepository) DeleteBatch(ctx context.Context, uuids []uuid.UUID) error {
	if len(uuids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
	DELETE FROM %s
		WHERE uuid = ANY($1)
	`, r.tableName)

	result, err := r.store.Db.ExecContext(ctx, query, pq.Array(uuids))
	if err != nil {
		return store.ContextError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return store.ContextError(err)
	}

	if int(rows) != len(uuids) {
		return fmt.Errorf("deleted %d of %d properties: some uuids not found", rows, len(uuids))
	}

	return nil
}

package product

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"time"
)

type Repository struct {
	store     *store.Store
	tableName string
}

func NewRepository(store *store.Store) *Repository {
	return &Repository{
		store:     store,
		tableName: "products",
	}
}

func (r *Repository) Create(ctx context.Context, p *ProductEnt) (*uint, error) {
	query := fmt.Sprintf(`
	INSERT INTO %s (
		uuid, name, unit, code, article, slug, active, step, brand_uuid, property,
		weight, width, length, height, volume, category_uuid, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
		$11, $12, $13, $14, $15, $16, $17, $18
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
		p.Unit,
		p.Code,
		p.Article,
		p.Slug,
		p.Active,
		p.Step,
		p.BrandUUID,
		p.Property,
		p.Weight,
		p.Width,
		p.Length,
		p.Height,
		p.Volume,
		p.CategoryUUID,
		p.CreatedAt,
		p.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, store.ContextError(err)
	}

	p.ID = id
	return &id, nil
}

func (r *Repository) GetByUUID(ctx context.Context, uuid string) (*ProductEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, name, unit, code, article, slug, active, step, brand_uuid, property,
			weight, width, length, height, volume, category_uuid, created_at, updated_at
		FROM %s
		WHERE uuid = $1
	`, r.tableName)

	var product ProductEnt
	err := r.store.Db.GetContext(ctx, &product, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	return &product, nil
}

func (r *Repository) Update(ctx context.Context, p *ProductEnt) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, unit = $2, code = $3, article = $4, slug = $5, active = $6,
			step = $7, brand_uuid = $8, property = $9, weight = $10, width = $11,
			length = $12, height = $13, volume = $14, category_uuid = $15, updated_at = $16
		WHERE uuid = $17
	`, r.tableName)

	p.UpdatedAt = time.Now()

	result, err := r.store.Db.ExecContext(ctx, query,
		p.Name,
		p.Unit,
		p.Code,
		p.Article,
		p.Slug,
		p.Active,
		p.Step,
		p.BrandUUID,
		p.Property,
		p.Weight,
		p.Width,
		p.Length,
		p.Height,
		p.Volume,
		p.CategoryUUID,
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

func (r *Repository) Delete(ctx context.Context, uuid string) error {
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

func (r *Repository) GetList(ctx context.Context) ([]ProductEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, name, unit, code, article, slug, active, step, brand_uuid, property,
			weight, width, length, height, volume, category_uuid, created_at, updated_at
		FROM %s
		ORDER BY created_at DESC
	`, r.tableName)

	var products []ProductEnt
	err := r.store.Db.SelectContext(ctx, &products, query)
	if err != nil {
		return nil, store.ContextError(err)
	}

	return products, nil
}

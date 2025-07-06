package category

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	store     *store.Store
	tableName string
}

func NewRepository(store *store.Store) *Repository {
	return &Repository{
		store:     store,
		tableName: "categories",
	}
}

func (r *Repository) CreateBatch(ctx context.Context, categories []CategoryEnt) error {
	if len(categories) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (
			uuid, name, slug, active, parent_uuid, created_at, updated_at
		) VALUES 
	`, r.tableName)

	args := []interface{}{}
	now := time.Now()

	valueStrings := make([]string, 0, len(categories))
	for i, c := range categories {
		c.CreatedAt = now
		c.UpdatedAt = now

		args = append(args,
			c.UUID,
			c.Name,
			c.Slug,
			c.Active,
			c.ParentUUID,
			c.CreatedAt,
			c.UpdatedAt,
		)

		start := i*7 + 1
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			start, start+1, start+2, start+3, start+4, start+5, start+6))
	}

	query += strings.Join(valueStrings, ", ")

	_, err := r.store.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return store.ContextError(err)
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, c *CategoryEnt) (*uint, error) {
	query := fmt.Sprintf(`
		INSERT INTO %s (
			uuid, name, slug, active, parent_uuid, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id
	`, r.tableName)

	now := time.Now()
	c.CreatedAt = now
	c.UpdatedAt = now

	var id uint
	err := r.store.Db.QueryRowxContext(ctx, query,
		c.UUID,
		c.Name,
		c.Slug,
		c.Active,
		c.ParentUUID,
		c.CreatedAt,
		c.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return nil, store.ContextError(err)
	}

	c.ID = id
	return &id, nil
}

func (r *Repository) GetByUUID(ctx context.Context, uuid string) (*CategoryEnt, error) {
	query := fmt.Sprintf(`
		SELECT id, uuid, name, slug, active, parent_uuid, created_at, updated_at
		FROM %s
		WHERE uuid = $1
	`, r.tableName)

	var category CategoryEnt
	err := r.store.Db.GetContext(ctx, &category, query, uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	return &category, nil
}

func (r *Repository) GetByUUIDs(ctx context.Context, uuids []string) ([]CategoryEnt, error) {
	if len(uuids) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf(`
		SELECT id, uuid, name, slug, active, parent_uuid, created_at, updated_at
		FROM %s
		WHERE uuid IN (?)
	`, r.tableName)

	query, args, err := sqlx.In(query, uuids)
	if err != nil {
		return nil, store.ContextError(err)
	}

	query = r.store.Db.Rebind(query)

	var categories []CategoryEnt
	err = r.store.Db.SelectContext(ctx, &categories, query, args...)
	if err != nil {
		return nil, store.ContextError(err)
	}

	return categories, nil
}

func (r *Repository) Update(ctx context.Context, c *CategoryEnt) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET name = $1, slug = $2, active = $3, parent_uuid = $4, updated_at = $5
		WHERE uuid = $6
	`, r.tableName)

	c.UpdatedAt = time.Now()

	result, err := r.store.Db.ExecContext(ctx, query,
		c.Name,
		c.Slug,
		c.Active,
		c.ParentUUID,
		c.UpdatedAt,
		c.UUID,
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

func (r *Repository) GetTree(ctx context.Context, rootUUID string) ([]*CategoryTree, error) {
	baseQuery := `
		WITH RECURSIVE descendants AS (
			SELECT
				id,
				uuid,
				name,
				slug,
				active,
				parent_uuid,
				created_at,
				updated_at,
				1 AS level
		FROM %[1]s
		WHERE %[2]s

			UNION ALL

			SELECT
				c.id,
				c.uuid,
				c.name,
				c.slug,
				c.active,
				c.parent_uuid,
				c.created_at,
				c.updated_at,
				d.level + 1
			FROM categories c
			INNER JOIN descendants d ON c.parent_uuid = d.uuid
		)
		SELECT *
		FROM descendants
		ORDER BY level, name
	`

	var (
		query string
		args  []any
	)

	if rootUUID == "" {
		query = fmt.Sprintf(baseQuery, r.tableName, "parent_uuid IS NULL")
	} else {
		query = fmt.Sprintf(baseQuery, r.tableName, "uuid = $1")
		args = []any{rootUUID}
	}

	var categories []*CategoryTree
	if err := r.store.Db.SelectContext(ctx, &categories, query, args...); err != nil {
		return nil, store.ContextError(err)
	}

	// Build the tree
	categoryMap := make(map[uuid.UUID]*CategoryTree, len(categories))
	var rootCategories []*CategoryTree

	for _, cat := range categories {
		categoryMap[cat.UUID] = cat
	}

	for _, cat := range categories {
		if cat.ParentUUID == nil {
			rootCategories = append(rootCategories, cat)
		} else if parent, ok := categoryMap[*cat.ParentUUID]; ok {
			parent.Children = append(parent.Children, cat)
		}
	}

	// Если передан rootUUID, но корневая категория не имеет потомков,
	// то она не попадёт в rootCategories (если её ParentUUID == nil).
	// Явно возвращаем её.
	if rootUUID != "" && len(rootCategories) == 0 {
		if root, ok := categoryMap[uuid.MustParse(rootUUID)]; ok {
			rootCategories = append(rootCategories, root)
		}
	}
	if rootUUID == "" || len(rootCategories) > 0 {
		return rootCategories, nil
	}
	return nil, store.ErrNotFound
}

// ================

// func InsertBatch[T any](
// 	ctx context.Context,
// 	db *sqlx.DB,
// 	table string,
// 	columns []string,
// 	items []T,
// 	toArgs func(T) []interface{},
// ) error {
// 	if len(items) == 0 {
// 		return nil
// 	}

// 	valueStrings := make([]string, 0, len(items))
// 	args := make([]interface{}, 0, len(items)*len(columns))

// 	for i, item := range items {
// 		argSlice := toArgs(item)
// 		if len(argSlice) != len(columns) {
// 			return fmt.Errorf("toArgs returned wrong number of args: got %d, want %d", len(argSlice), len(columns))
// 		}

// 		start = i*len(columns) + 1  // правильно: нумерация с 1
// 		valueStrings = append(valueStrings, buildPlaceholders(start, len(columns)))
// 		args = append(args, argSlice...)
// 	}

// 	query := fmt.Sprintf(
// 		"INSERT INTO %s (%s) VALUES %s",
// 		table,
// 		strings.Join(columns, ", "),
// 		strings.Join(valueStrings, ", "),
// 	)

// 	_, err := db.ExecContext(ctx, query, args...)
// 	return err
// }

// // buildPlaceholders генерирует ($N,$N+1,...,$N+count-1)
// func buildPlaceholders(start, count int) string {
// 	placeholders := make([]string, count)
// 	for i := 0; i < count; i++ {
// 		placeholders[i] = fmt.Sprintf("$%d", start+i)
// 	}
// 	return "(" + strings.Join(placeholders, ",") + ")"
// }

// type CategoryEnt struct {
// 	UUID       string
// 	Name       string
// 	Slug       string
// 	Active     bool
// 	ParentUUID *string
// 	CreatedAt  time.Time
// 	UpdatedAt  time.Time
// }

// func categoryToArgs(c CategoryEnt) []interface{} {
// 	now := time.Now()
// 	return []interface{}{
// 		c.UUID,
// 		c.Name,
// 		c.Slug,
// 		c.Active,
// 		c.ParentUUID,
// 		now,
// 		now,
// 	}
// }

// // где-то в коде:
// err := InsertBatch(ctx, db, "categories",
// 	[]string{"uuid", "name", "slug", "active", "parent_uuid", "created_at", "updated_at"},
// 	categoriesSlice,
// 	categoryToArgs,
// )
// if err != nil {
// 	// обработка ошибки
// }

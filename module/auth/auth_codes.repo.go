package auth

import (
	"context"
	"database/sql"
	"fmt"
	"go-monolite/internal/store"
	"time"
)

type AuthCodeRepository struct {
	store     *store.Store
	tableName string
}

func NewAuthCodeRepository(store *store.Store) *AuthCodeRepository {
	return &AuthCodeRepository{
		store:     store,
		tableName: "auth_codes",
	}
}

func (r *AuthCodeRepository) CountCodesLast24Hours(ctx context.Context, email, phone string) (int, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) FROM %s
		WHERE created_at >= NOW() - INTERVAL '24 hours'
		AND %%s
	`, r.tableName)

	var condition string
	var args []interface{}

	if email != "" {
		condition = "email = $1"
		args = append(args, email)
	} else {
		condition = "phone = $1"
		args = append(args, phone)
	}

	finalQuery := fmt.Sprintf(query, condition)

	var count int
	err := r.store.Db.GetContext(ctx, &count, finalQuery, args...)
	if err != nil {
		return 0, store.ContextError(err)
	}

	return count, nil
}

func (r *AuthCodeRepository) GetActiveCode(ctx context.Context, email, phone string) (*AuthCode, error) {
	query := fmt.Sprintf(`
		SELECT id, email, phone, code, expires_at, used, created_at
		FROM %s
		WHERE expires_at >= NOW() AND used = FALSE AND %%s
		ORDER BY created_at DESC
		LIMIT 1
	`, r.tableName)

	var condition string
	var args []interface{}

	if email != "" {
		condition = "email = $1"
		args = append(args, email)
	} else {
		condition = "phone = $1"
		args = append(args, phone)
	}

	finalQuery := fmt.Sprintf(query, condition)

	var code AuthCode
	err := r.store.Db.GetContext(ctx, &code, finalQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, store.ContextError(err)
	}

	return &code, nil
}

func (r *AuthCodeRepository) SaveCode(ctx context.Context, email, phone, code string, expiresAt time.Time) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (email, phone, code, expires_at)
		VALUES ($1, $2, $3, $4)
	`, r.tableName)

	_, err := r.store.Db.ExecContext(ctx, query,
		email,
		phone,
		code,
		expiresAt,
	)
	if err != nil {
		return store.ContextError(err)
	}

	return nil
}

package user

import (
	"context"
	"go-monolite/internal/store"
)

type Repository struct {
	store     *store.Store
	tableName string
}

func NewRepository(store *store.Store) *Repository {
	return &Repository{
		store:     store,
		tableName: "users",
	}
}

func (r *Repository) GetUser(ctx context.Context, id uint) error {
	return nil
}

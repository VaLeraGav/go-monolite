package auth

import (
	"go-monolite/internal/store"
)

type UserTokensRepository struct {
	store     *store.Store
	tableName string
}

func NewUserTokensRepository(store *store.Store) *UserTokensRepository {
	return &UserTokensRepository{
		store:     store,
		tableName: "user_tokens",
	}
}

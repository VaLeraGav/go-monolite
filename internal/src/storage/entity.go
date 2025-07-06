package storage

import (
	"time"

	"github.com/google/uuid"
)

type StorageEnt struct {
	ID        uint      `db:"id"`
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	Active    string    `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductStorageEnt struct {
	ID          uint      `db:"id"`
	ProductUUID uuid.UUID `db:"product_uuid"`
	StorageUUID uuid.UUID `db:"storage_uuid" `
	Active      string    `db:"active"`
	Quantity    int       `db:"quantity"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (e StorageEnt) ToResponse() StorageResponse {
	return StorageResponse{
		ID:     e.ID,
		UUID:   e.UUID,
		Name:   e.Name,
		Active: e.Active,
	}
}

package price

import (
	"time"

	"github.com/google/uuid"
)

type TypePriceEnt struct {
	ID        uint      `db:"id"`
	UUID      uuid.UUID `db:"uuid"`
	Name      string    `db:"name"`
	Active    string    `db:"active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ProductPriceEnt struct {
	ID            uint      `db:"id"`
	ProductUUID   uuid.UUID `db:"product_uuid"`
	TypePriceUUID uuid.UUID `db:"type_price_uuid"`
	Active        string    `db:"active"`
	Price         float64   `db:"price"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (e TypePriceEnt) ToResponse() TypePriceResponse {
	return TypePriceResponse{
		ID:     e.ID,
		UUID:   e.UUID,
		Name:   e.Name,
		Active: e.Active,
	}
}

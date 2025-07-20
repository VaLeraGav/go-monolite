package price

import (
	"go-monolite/pkg/validator"

	"github.com/google/uuid"
)

type UpsertRequest struct {
	General       *GeneralRequest   `json:"general,omitempty"`
	ProductPrices []ProductPriceDto `json:"data" validate:"required"`
}

type TypePriceResponse struct {
	ID     uint      `json:"id" example:"1"`
	UUID   uuid.UUID `json:"uuid" example:"a8098c1a-f86e-11da-bd1a-00112444be1e"`
	Name   string    `json:"name" example:"Розничная цена"`
	Active string    `json:"active" example:"true"`
}

type GeneralRequest struct {
	Prices []TypePriceRequest `json:"prices" validate:"required"`
}

type TypePriceRequest struct {
	UUID   uuid.UUID `json:"uuid" validate:"required"`
	Name   string    `json:"name" validate:"required"`
	Active string    `json:"active" validate:"required,oneof=Y N"`
}

type ProductPriceDto struct {
	ProductUUID   uuid.UUID             `json:"product_uuid" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProductPrices []ProductPriceItemDto `json:"prices" validate:"required"`
}

type ProductPriceItemDto struct {
	ProductUUID   uuid.UUID `swaggerignore:"true"`
	TypePriceUUID uuid.UUID `json:"type_price_uuid" validate:"required" example:"550e8400-e29b-41d4-a713-446655440000"`
	Active        string    `json:"active" validate:"required,oneof=Y N"  example:"Y"`
	Price         float64   `json:"price" validate:"gte=0" example:"200"`
}

type UpsertResponse struct {
	TypePrice    *TypePriceResponseDetails    `json:"type_price"`
	ProductPrice *ProductPriceResponseDetails `json:"product_price"`
}

type TypePriceResponseDetails struct {
	CountDeleted  int `json:"count_deleted"`
	CountInserted int `json:"count_inserted"`
	CountUpdated  int `json:"count_updated"`
}

type ProductPriceResponseDetails struct {
	CountInserted int `json:"count_inserted"`
	CountUpdated  int `json:"count_updated"`
}

func (r *UpsertRequest) Validate() error {
	if err := validator.Validate(r); err != nil {
		return err
	}

	for _, data := range r.ProductPrices {
		if err := data.Validate(); err != nil {
			return err
		}
		for _, productPrice := range data.ProductPrices {
			if err := productPrice.Validate(); err != nil {
				return err
			}
		}
	}
	if r.General != nil {
		for _, dto := range r.General.Prices {
			if err := dto.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (v TypePriceRequest) ToEntity() *TypePriceEnt {
	return &TypePriceEnt{
		UUID:   v.UUID,
		Name:   v.Name,
		Active: v.Active,
	}
}

func (d *TypePriceRequest) Validate() error {
	return validator.Validate(d)
}

func (d *ProductPriceDto) Validate() error {
	return validator.Validate(d)
}

func (d *ProductPriceItemDto) Validate() error {
	return validator.Validate(d)
}

func (d *ProductPriceItemDto) ToEntity() *ProductPriceEnt {
	return &ProductPriceEnt{
		ProductUUID:   d.ProductUUID,
		TypePriceUUID: d.TypePriceUUID,
		Active:        d.Active,
		Price:         d.Price,
	}
}

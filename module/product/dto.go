package product

import (
	"encoding/json"
	"go-monolite/pkg/validator"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type ProductDto struct {
	ID           uint            `json:"id" example:"1"`
	UUID         uuid.UUID       `json:"uuid" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name         string          `json:"name" validate:"required" example:"Корм для кошек"`
	Unit         *string         `json:"unit,omitempty" example:"шт"`
	Code         int             `json:"code" validate:"required,gt=0" example:"123456"`
	Article      *string         `json:"article,omitempty" example:"A-1234"`
	Active       string          `json:"active" validate:"required,oneof=Y N" example:"Y"`
	Step         *int            `json:"step,omitempty" validate:"omitempty,gt=0" example:"5"`
	BrandUUID    *uuid.UUID      `json:"brand_uuid,omitempty" example:"1d5f1d04-3d79-4b02-9245-c2e17f54cb10"`
	Property     json.RawMessage `json:"property,omitempty" example:"{\"color\": \"red\", \"material\": \"plastic\"}"`
	Weight       *float64        `json:"weight,omitempty" validate:"omitempty,numeric" example:"0.75"`
	Width        *float64        `json:"width,omitempty" validate:"omitempty,numeric" example:"20.5"`
	Length       *float64        `json:"length,omitempty" validate:"omitempty,numeric" example:"30.0"`
	Height       *float64        `json:"height,omitempty" validate:"omitempty,numeric" example:"15.0"`
	Volume       *float64        `json:"volume,omitempty" validate:"omitempty,numeric" example:"9.2"`
	CategoryUUID uuid.UUID       `json:"category_uuid" validate:"required" example:"7c9e6679-7425-40de-944b-e07fc1f90ae7"`
}

func (p *ProductDto) ToEntity() *ProductEnt {
	return &ProductEnt{
		ID:           p.ID,
		UUID:         p.UUID,
		Name:         p.Name,
		Unit:         p.Unit,
		Code:         p.Code,
		Article:      p.Article,
		Slug:         slug.Make(p.Name),
		Active:       p.Active,
		Step:         p.Step,
		BrandUUID:    p.BrandUUID,
		Property:     p.Property,
		Weight:       p.Weight,
		Width:        p.Width,
		Length:       p.Length,
		Height:       p.Height,
		Volume:       p.Volume,
		CategoryUUID: p.CategoryUUID,
	}
}

func (d *ProductDto) Validate() error {
	return validator.Validate(d)
}

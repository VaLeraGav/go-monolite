package product

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProductEnt struct {
	ID           uint            `db:"id" json:"id"`
	UUID         uuid.UUID       `db:"uuid" json:"uuid"`
	Name         string          `db:"name" json:"name"`
	Unit         *string         `db:"unit" json:"unit,omitempty"`
	Code         int             `db:"code" json:"code"`
	Article      *string         `db:"article" json:"article,omitempty"`
	Slug         string          `db:"slug" json:"slug"`
	Active       string          `db:"active" json:"active"`
	Step         *int            `db:"step" json:"step,omitempty"`
	BrandUUID    *uuid.UUID      `db:"brand_uuid" json:"brand_uuid,omitempty"`
	Property     json.RawMessage `db:"property" json:"property"`
	Weight       *float64        `db:"weight" json:"weight,omitempty"`
	Width        *float64        `db:"width" json:"width,omitempty"`
	Length       *float64        `db:"length" json:"length,omitempty"`
	Height       *float64        `db:"height" json:"height,omitempty"`
	Volume       *float64        `db:"volume" json:"volume,omitempty"`
	CategoryUUID uuid.UUID       `db:"category_uuid" json:"category_uuid"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time       `db:"updated_at" json:"updated_at"`
}

func (e *ProductEnt) PatchDto(request ProductDto) ProductDto {
	request.UUID = e.UUID
	if request.Name == "" {
		request.Name = e.Name
	}
	if request.Unit == nil {
		request.Unit = e.Unit
	}
	if request.Code == 0 {
		request.Code = e.Code
	}
	if request.Article == nil {
		request.Article = e.Article
	}
	if request.Active == "" {
		request.Active = e.Active
	}
	if request.Step == nil {
		request.Step = e.Step
	}
	if request.BrandUUID == nil {
		request.BrandUUID = e.BrandUUID
	}
	if request.Property == nil {
		request.Property = e.Property
	}
	if request.Weight == nil {
		request.Weight = e.Weight
	}
	if request.Width == nil {
		request.Width = e.Width
	}
	if request.Length == nil {
		request.Length = e.Length
	}
	if request.Height == nil {
		request.Height = e.Height
	}
	if request.Volume == nil {
		request.Volume = e.Volume
	}
	if request.CategoryUUID == uuid.Nil {
		request.CategoryUUID = e.CategoryUUID
	}
	return request
}

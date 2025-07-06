package property

import (
	"time"

	"github.com/google/uuid"
)

type PropertyValueEnt struct {
	ID           uint      `db:"id" json:"id"`
	Key          string    `db:"key" json:"key"` // не UUID так как 1с предает true false для bool значений
	Slug         string    `db:"slug" json:"slug"`
	PropertyUUID uuid.UUID `db:"property_uuid" json:"property_uuid"`
	Value        string    `db:"value" json:"value"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type PropertyEnt struct {
	ID        uint               `db:"id" json:"id"`
	UUID      uuid.UUID          `db:"uuid" json:"uuid"`
	Slug      string             `db:"slug" json:"slug"`
	Type      string             `db:"type" json:"type"`
	Name      string             `db:"name" json:"name"`
	CreatedAt time.Time          `db:"created_at" json:"created_at"`
	UpdatedAt time.Time          `db:"updated_at" json:"updated_at"`
	Values    []PropertyValueEnt `json:"values"`
}

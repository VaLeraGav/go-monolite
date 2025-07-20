package property

import (
	"encoding/json"
	"go-monolite/pkg/validator"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type PropResponse struct {
	Property       *PropertyResponse       `json:"property"`
	PropertyValues *PropertyValuesResponse `json:"property_values"`
}

type PropertyResponse struct {
	Deletes []PropertyEnt `json:"deletes,omitempty"`
	Inserts []PropertyEnt `json:"inserts,omitempty"`
	Updates []PropertyEnt `json:"updates,omitempty"`
}

type PropertyValuesResponse struct {
	Deletes []PropertyValueEnt `json:"deletes,omitempty"`
	Inserts []PropertyValueEnt `json:"inserts,omitempty"`
	Updates []PropertyValueEnt `json:"updates,omitempty"`
}

type PropertyValueDto struct {
	Key          string    `json:"key" validate:"required"`
	PropertyUUID uuid.UUID `json:"property_uuid" validate:"required"`
	Value        string    `json:"value" validate:"required"`
}

func (v PropertyValueDto) ToEntity() *PropertyValueEnt {
	return &PropertyValueEnt{
		Key:          v.Key,
		Slug:         slug.Make(v.Value),
		PropertyUUID: v.PropertyUUID,
		Value:        v.Value,
	}
}

func (d *PropertyValueDto) Validate() error {
	return validator.Validate(d)
}

type PropertyDto struct {
	UUID   uuid.UUID          `json:"uuid" validate:"required"`
	Type   string             `json:"type" validate:"required"`
	Name   string             `json:"name" validate:"required"`
	Values []PropertyValueDto `json:"values"`
}

func (v PropertyDto) ToEntity() *PropertyEnt {
	return &PropertyEnt{
		UUID: v.UUID,
		Slug: slug.Make(v.Name),
		Type: v.Type,
		Name: v.Name,
	}
}

func (d *PropertyDto) Validate() error {
	return validator.Validate(d)
}

func ValidatePropertyDtos(pdtos []PropertyDto) error {
	for _, pdto := range pdtos {
		if err := pdto.Validate(); err != nil {
			return err
		}
		for _, pvdto := range pdto.Values {
			if err := pvdto.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

func ParsePropertyDto(body []byte) ([]PropertyDto, error) {
	var properties []PropertyDto
	if err := json.Unmarshal(body, &properties); err != nil {
		return nil, err
	}

	for i := range properties {
		for j := range properties[i].Values {
			properties[i].Values[j].PropertyUUID = properties[i].UUID
		}
	}

	return properties, nil
}

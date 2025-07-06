package storage

import (
	"go-monolite/pkg/validator"

	"github.com/google/uuid"
)

type UpsertRequest struct {
	General         *GeneralRequest     `json:"general,omitempty"`
	ProductStorages []ProductStorageDto `json:"data" validate:"required"`
}

type GeneralRequest struct {
	Storages []StorageDto `json:"storages" validate:"required"`
}

type StorageDto struct {
	UUID   uuid.UUID `json:"uuid" validate:"required"`
	Name   string    `json:"name" validate:"required"`
	Active string    `json:"active" validate:"required,oneof=Y N"`
}

type StorageResponse struct {
	ID     uint      `json:"id" example:"1"`
	UUID   uuid.UUID `json:"uuid" validate:"required" example:"550e8400-e29b-41d4-a713-446655440000"`
	Name   string    `json:"name" validate:"required" example:"Backup Warehouse"`
	Active string    `json:"active" validate:"required,oneof=Y N" example:"Y"`
}

type ProductStorageDto struct {
	ProductUUID     uuid.UUID               `json:"product_uuid" validate:"required"`
	ProductStorages []ProductStorageItemDto `json:"storages" validate:"required"`
}

type ProductStorageItemDto struct {
	ProductUUID uuid.UUID `json:"product_uuid"`
	StorageUUID uuid.UUID `json:"storage_uuid" validate:"required"`
	Active      string    `json:"active" validate:"required,oneof=Y N"`
	Quantity    int       `json:"quantity" validate:"gte=0"`
}

type UpsertResponse struct {
	Storage        *StorageUpsertStatsResponse        `json:"storage"`
	ProductStorage *ProductStorageUpsertStatsResponse `json:"product_storage"`
}

type StorageUpsertStatsResponse struct {
	CountDeleted  int `json:"count_deleted"`
	CountInserted int `json:"count_inserted"`
	CountUpdated  int `json:"count_updated"`
}

type ProductStorageUpsertStatsResponse struct {
	CountInserted int `json:"count_inserted"`
	CountUpdated  int `json:"count_updated"`
}

func (r *UpsertRequest) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
	}

	if err := validator.Validate(r, tagMsg); err != nil {
		return err
	}

	for _, data := range r.ProductStorages {
		for _, entry := range data.ProductStorages {
			if err := entry.Validate(); err != nil {
				return err
			}
		}
		if err := data.Validate(); err != nil {
			return err
		}
	}

	if r.General != nil {
		for _, dto := range r.General.Storages {
			if err := dto.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *StorageDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"oneof":    "должно быть Y или N",
	}
	return validator.Validate(d, tagMsg)
}

func (d *ProductStorageDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
	}
	return validator.Validate(d, tagMsg)
}

func (d *ProductStorageItemDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"oneof":    "должно быть Y или N",
	}
	return validator.Validate(d, tagMsg)
}

func (v StorageDto) ToEntity() *StorageEnt {
	return &StorageEnt{
		UUID:   v.UUID,
		Name:   v.Name,
		Active: v.Active,
	}
}

func (d *ProductStorageItemDto) ToEntity() *ProductStorageEnt {
	return &ProductStorageEnt{
		ProductUUID: d.ProductUUID,
		StorageUUID: d.StorageUUID,
		Active:      d.Active,
		Quantity:    d.Quantity,
	}
}

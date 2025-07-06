package category

import (
	"go-monolite/pkg/validator"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CategoryRequest struct {
	UUID       uuid.UUID  `json:"uuid" validate:"required" example:"1"`
	Name       string     `json:"name" validate:"required"  example:"Категория a"`
	Active     string     `json:"active" validate:"required,oneof=Y N" example:"Y"`
	ParentUUID *uuid.UUID `json:"parent_uuid,omitempty"  example:"b3d8ef13-1234-4567-87ab-abcdef123456"`
}

type CategoryTreeResponse struct {
	CategoryResponse
	Level    int                     `json:"level" example:"1"`
	Children []*CategoryTreeResponse `json:"children,omitempty" swaggertype:"array,object"`
}

type CategoryResponse struct {
	ID         uint       `json:"id" example:"1"`
	UUID       uuid.UUID  `json:"uuid" example:"b3d8ef13-1234-4567-89ab-abcdef123456"`
	Slug       string     `db:"slug" example:"categoriia-a"`
	Name       string     `json:"name" example:"Категория a"`
	ParentUUID *uuid.UUID `json:"parent_uuid,omitempty" example:"b3d8ef13-1234-4567-87ab-abcdef123456"`
	Active     string     `json:"active" example:"Y"`
}

func (c CategoryRequest) ToEntity() CategoryEnt {
	return CategoryEnt{
		UUID:       c.UUID,
		Name:       c.Name,
		Slug:       slug.Make(c.Name),
		Active:     c.Active,
		ParentUUID: c.ParentUUID,
	}
}

func (c *CategoryRequest) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"oneof":    "может быть одним из: Y или N",
	}

	return validator.Validate(c, tagMsg)
}

func (c CategoryRequest) GetKey() string {
	return c.UUID.String()
}

func MapCategoryTreesToResponse(trees []*CategoryTree) []*CategoryTreeResponse {
	result := make([]*CategoryTreeResponse, 0, len(trees))
	for _, tree := range trees {
		result = append(result, mapCategoryTreeToResponse(tree))
	}
	return result
}

func mapCategoryTreeToResponse(tree *CategoryTree) *CategoryTreeResponse {
	if tree == nil {
		return nil
	}
	children := make([]*CategoryTreeResponse, 0, len(tree.Children))
	for _, child := range tree.Children {
		children = append(children, mapCategoryTreeToResponse(child))
	}
	return &CategoryTreeResponse{
		CategoryResponse: CategoryResponse{
			ID:         tree.ID,
			UUID:       tree.UUID,
			Slug:       tree.Slug,
			Name:       tree.Name,
			ParentUUID: tree.ParentUUID,
			Active:     tree.Active,
		},
		Level:    tree.Level,
		Children: children,
	}
}

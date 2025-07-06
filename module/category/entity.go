package category

import (
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

type CategoryTree struct {
	CategoryEnt
	Level    int             `db:"level"`
	Children []*CategoryTree `db:"children"`
}

type CategoryEnt struct {
	ID         uint       `db:"id"`
	UUID       uuid.UUID  `db:"uuid"`
	Slug       string     `db:"slug"`
	Name       string     `db:"name"`
	ParentUUID *uuid.UUID `db:"parent_uuid"`
	Active     string     `db:"active"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

func (e CategoryEnt) PatchDto(request *CategoryRequest) CategoryEnt {
	if request.Name != "" && e.Name != request.Name {
		e.Slug = slug.Make(request.Name)
		e.Name = request.Name
	}
	if request.Active != "" && e.Active != request.Active {
		e.Active = request.Active
	}
	if request.ParentUUID != nil && e.ParentUUID != request.ParentUUID {
		e.ParentUUID = request.ParentUUID
	}
	return e
}

func (e CategoryEnt) ToResponse() CategoryResponse {
	return CategoryResponse{
		ID:         e.ID,
		UUID:       e.UUID,
		Slug:       e.Slug,
		Name:       e.Name,
		ParentUUID: e.ParentUUID,
		Active:     e.Active,
	}
}

func (c CategoryEnt) GetKey() string {
	return c.UUID.String()
}

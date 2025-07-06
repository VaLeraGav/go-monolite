package category

import (
	"context"
	"errors"
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, reqs []CategoryRequest) ([]CategoryResponse, string, error) {
	if len(reqs) == 0 {
		return nil, "нет элементов для вставки", nil
	}

	uuids := make([]string, 0, len(reqs))
	for _, categoryReq := range reqs {
		err := categoryReq.Validate()
		if err != nil {
			return nil, "", err
		}
		uuids = append(uuids, categoryReq.UUID.String())
	}

	existing, err := s.repo.GetByUUIDs(ctx, uuids)
	if err != nil {
		return nil, "произошла ошибка при получении категорий", err
	}

	toInsert := helper.FilterNewEntities(reqs, existing)
	if len(toInsert) == 0 {
		return nil, "нет элементов для вставки", nil
	}

	err = s.repo.CreateBatch(ctx, toInsert)
	if err != nil {
		return nil, "произошла ошибка при создания категорий", err
	}

	insertedUUIDs := helper.CollectKeys(toInsert)

	categories, err := s.repo.GetByUUIDs(ctx, insertedUUIDs)
	if err != nil {
		return nil, "произошла ошибка при возвращении категорий", err
	}

	return helper.ToResponse(categories), "категории успешно создались", nil
}

func (s *Service) Update(ctx context.Context, req *CategoryRequest) (*CategoryResponse, string, error) {
	existing, err := s.repo.GetByUUID(ctx, req.UUID.String())
	if err != nil {
		return nil, "произошла ошибка при получении категории", err
	}
	if existing == nil {
		return nil, "категория не найдена", store.ErrNotFound
	}

	updated := existing.PatchDto(req)

	err = s.repo.Update(ctx, &updated)
	if err != nil {
		return nil, "произошла ошибка при обновлении категорий", err
	}

	afterUpdate, err := s.repo.GetByUUID(ctx, req.UUID.String())
	if err != nil {
		return nil, "произошла ошибка при возвращении категорий", err
	}

	response := afterUpdate.ToResponse()

	return &response, "", nil
}

func (s *Service) Delete(ctx context.Context, uuid string) (string, error) {
	_, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "категория не найдена", store.ErrNotFound
		}
		return "произошла ошибка при получении категории", err
	}
	err = s.repo.Delete(ctx, uuid)
	if err != nil {
		return "произошла ошибка при удалении категорий", err
	}

	return "успешно удалили", nil
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*CategoryResponse, string, error) {
	category, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, "категория не найдена", store.ErrNotFound
		}
		return nil, "произошла ошибка при получении категории", err
	}

	resp := category.ToResponse()

	return &resp, "", nil
}

func (s *Service) GetTree(ctx context.Context, rootUUID string) ([]*CategoryTreeResponse, string, error) {
	categoryTree, err := s.repo.GetTree(ctx, rootUUID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, "категория не найдена", store.ErrNotFound
		}
		return nil, "произошла ошибка при получении дерева категорий", err
	}

	response := MapCategoryTreesToResponse(categoryTree)

	return response, "", nil
}

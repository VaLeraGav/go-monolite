package product

import (
	"context"
	"errors"
	"go-monolite/internal/store"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, request *ProductDto) (*uint, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repo.Create(
		ctx,
		request.ToEntity(),
	)
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*ProductEnt, error) {
	products, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return products, nil
}

func (s *Service) Update(ctx context.Context, request *ProductDto) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repo.Update(
		ctx,
		request.ToEntity(),
	)
}

func (s *Service) Delete(ctx context.Context, uuid string) error {
	return s.repo.Delete(ctx, uuid)
}

func (s *Service) GetList(ctx context.Context) ([]ProductEnt, error) {
	return s.repo.GetList(ctx)
}

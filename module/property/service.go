package property

import (
	"context"
	"errors"
	"fmt"
	"go-monolite/internal/store"
	"go-monolite/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	propertyRepo       *PropertyRepository
	propertyValuesRepo *PropertyValuesRepository
}

func NewService(propertyRepo *PropertyRepository, propertyValuesRepo *PropertyValuesRepository) *Service {
	return &Service{propertyRepo, propertyValuesRepo}
}

func (s *Service) Upsert(ctx context.Context, dtos []PropertyDto) (*PropResponse, error) {
	if err := ValidatePropertyDtos(dtos); err != nil {
		return nil, err
	}

	propResp, err := s.upsertProperties(ctx, dtos)
	if err != nil {
		return nil, err
	}

	propValResp, err := s.upsertPropertyValues(ctx, dtos)
	if err != nil {
		return nil, err
	}

	return &PropResponse{
		Property:       propResp,
		PropertyValues: propValResp,
	}, nil
}

func (s *Service) upsertProperties(ctx context.Context, dtos []PropertyDto) (*PropertyResponse, error) {
	deletes, inserts, updates, err := s.preparePropertyDiff(ctx, dtos)
	if err != nil {
		return nil, err
	}

	if err := s.applyPropertyChanges(ctx, deletes, inserts, updates); err != nil {
		return nil, err
	}

	return &PropertyResponse{
		Deletes: deletes,
		Inserts: inserts,
		Updates: updates,
	}, nil
}

func (s *Service) upsertPropertyValues(ctx context.Context, dtos []PropertyDto) (*PropertyValuesResponse, error) {
	deletes, inserts, updates, err := s.preparePropertyValuesDiff(ctx, dtos)
	if err != nil {
		return nil, err
	}

	if err := s.applyPropertyValueChanges(ctx, deletes, inserts, updates); err != nil {
		return nil, err
	}

	return &PropertyValuesResponse{
		Deletes: deletes,
		Inserts: inserts,
		Updates: updates,
	}, nil
}

func (s *Service) preparePropertyDiff(ctx context.Context, dtos []PropertyDto) (deletes, inserts, updates []PropertyEnt, err error) {
	existing, err := s.propertyRepo.GetList(ctx)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, nil, nil, fmt.Errorf("ошибка при выполнении propertyRepo.GetList: %w", err)
	}

	desired := toPropertySlice(dtos)

	currentMap := toPropertyMap(existing)
	desiredMap := toPropertyMap(desired)

	deletes, inserts, updates = diffProperties(currentMap, desiredMap)
	return deletes, inserts, updates, nil
}

func (s *Service) applyPropertyChanges(ctx context.Context, deletes, inserts, updates []PropertyEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(deletes) > 0 {
		g.Go(func() error {
			uuids := make([]uuid.UUID, 0, len(deletes))
			for _, p := range deletes {
				uuids = append(uuids, p.UUID)
			}
			err := s.propertyRepo.DeleteBatch(ctx, uuids)
			logger.DebugCtx(ctx, "delete []uuids", "uuids", uuids)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyRepo.DeleteBatch: %w", err)
			}
			return nil
		})
	}

	if len(inserts) > 0 {
		g.Go(func() error {
			err := s.propertyRepo.CreateBatch(ctx, inserts)
			logger.DebugCtx(ctx, "inserts []PropertyEnt", "inserts", inserts)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyRepo.CreateBatch: %w", err)
			}
			return nil
		})
	}

	if len(updates) > 0 {
		g.Go(func() error {
			err := s.propertyRepo.UpdateBatch(ctx, updates)
			logger.DebugCtx(ctx, "updates []PropertyEnt", "updates", updates)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyRepo.UpdateBatch: %w", err)
			}
			return nil
		})
	}
	return g.Wait()
}

func (s *Service) preparePropertyValuesDiff(ctx context.Context, dtos []PropertyDto) (deletes, inserts, updates []PropertyValueEnt, err error) {
	var allValueDtos []PropertyValueDto
	for _, dto := range dtos {
		allValueDtos = append(allValueDtos, dto.Values...)
	}

	existing, err := s.propertyValuesRepo.GetList(ctx)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, nil, nil, fmt.Errorf("ошибка при выполнении propertyValuesRepo.GetList: %w", err)

	}

	desired := toPropertyValueSlice(allValueDtos)

	currentMap := toPropertyValueMap(existing)
	desiredMap := toPropertyValueMap(desired)

	// deletes нет, были удалены вместе с property cascade
	deletes, inserts, updates = diffPropertyValues(currentMap, desiredMap)
	return deletes, inserts, updates, nil
}

func (s *Service) applyPropertyValueChanges(ctx context.Context, deletes, inserts, updates []PropertyValueEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(deletes) > 0 {
		g.Go(func() error {
			keys := make([]string, 0, len(deletes))
			for _, p := range deletes {
				keys = append(keys, p.Key)
			}
			err := s.propertyValuesRepo.DeleteBatch(ctx, keys)
			logger.DebugCtx(ctx, "inserts []PropertyValueEnt", "inserts", inserts)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyValuesRepo.CreateBatch: %w", err)
			}
			return nil
		})
	}

	if len(inserts) > 0 {
		g.Go(func() error {
			err := s.propertyValuesRepo.CreateBatch(ctx, inserts)
			logger.DebugCtx(ctx, "inserts []PropertyValueEnt", "inserts", inserts)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyValuesRepo.CreateBatch: %w", err)
			}
			return nil
		})
	}

	if len(updates) > 0 {
		g.Go(func() error {
			err := s.propertyValuesRepo.UpdateBatch(ctx, updates)
			logger.DebugCtx(ctx, "updates []PropertyValueEnt", "updates", updates)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении propertyValuesRepo.UpdateBatch: %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

// TODO: так же передается поле "помечено на удаление", но решил пока игнорировать
func diffProperties(current, desired map[string]PropertyEnt) (deletes, inserts, updates []PropertyEnt) {
	for key, curr := range current {
		if _, ok := desired[key]; !ok {
			deletes = append(deletes, curr)
		}
	}
	for key, want := range desired {
		curr, ok := current[key]
		if ok {
			if !isPropertyEqual(curr, want) {
				updates = append(updates, want)
			}
		} else {
			inserts = append(inserts, want)
		}
	}
	return
}

func diffPropertyValues(current, desired map[string]PropertyValueEnt) (deletes, inserts, updates []PropertyValueEnt) {
	for key, curr := range current {
		if _, ok := desired[key]; !ok {
			deletes = append(deletes, curr)
		}
	}
	for key, want := range desired {
		curr, ok := current[key]
		if ok {
			if !isPropertyValueEqual(curr, want) {
				updates = append(updates, want)
			}
		} else {
			inserts = append(inserts, want)
		}
	}
	return
}

func isPropertyEqual(a, b PropertyEnt) bool {
	return a.Name == b.Name && a.Slug == b.Slug && a.Type == b.Type
}

func isPropertyValueEqual(a, b PropertyValueEnt) bool {
	return a.PropertyUUID == b.PropertyUUID && a.Slug == b.Slug && a.Value == b.Value
}

func toPropertyMap(list []PropertyEnt) map[string]PropertyEnt {
	result := make(map[string]PropertyEnt)
	for _, p := range list {
		result[p.UUID.String()] = p
	}
	return result
}

func toPropertyValueMap(list []PropertyValueEnt) map[string]PropertyValueEnt {
	result := make(map[string]PropertyValueEnt)
	for _, p := range list {
		result[p.Key] = p
	}
	return result
}

func toPropertyValueSlice(dtos []PropertyValueDto) []PropertyValueEnt {
	desired := make([]PropertyValueEnt, 0, len(dtos))
	for _, dto := range dtos {
		desired = append(desired, *dto.ToEntity())
	}
	return desired
}

func toPropertySlice(dtos []PropertyDto) []PropertyEnt {
	desired := make([]PropertyEnt, 0, len(dtos))
	for _, dto := range dtos {
		desired = append(desired, *dto.ToEntity())
	}
	return desired
}

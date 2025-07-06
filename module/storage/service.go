package storage

import (
	"context"
	"errors"
	"fmt"
	"go-monolite/internal/store"
	"go-monolite/pkg/helper"
	"go-monolite/pkg/logger"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	storageRepo        *StorageRepository
	productStorageRepo *ProductStoragesRepository
}

func NewService(storageRepo *StorageRepository, productStorageRepo *ProductStoragesRepository) *Service {
	return &Service{storageRepo, productStorageRepo}
}

func (s *Service) GetStorage(ctx context.Context) ([]StorageResponse, string, error) {
	existing, err := s.storageRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, "склады не найден", store.ErrNotFound
		}
		return nil, "произошла ошибка при получении складов", err
	}

	return helper.ToResponse(existing), "", nil
}

func (s *Service) Upsert(ctx context.Context, request UpsertRequest) (*UpsertResponse, string, error) {
	if err := request.Validate(); err != nil {
		return nil, "", err
	}

	tx, err := s.storageRepo.store.Db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txCtx := store.WithTx(ctx, tx)

	storageResponse, err := s.upsertStorages(txCtx, request)
	if err != nil {
		return nil, "Произошла ошибка при создании склада", err
	}

	productStorageUpsertResponse, err := s.upsertStoragesValue(txCtx, request.ProductStorages)
	if err != nil {
		return nil, "произошла ошибка при записи складов для товара", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &UpsertResponse{
		Storage:        storageResponse,
		ProductStorage: productStorageUpsertResponse,
	}, "", nil
}

func (s *Service) upsertStorages(ctx context.Context, request UpsertRequest) (*StorageUpsertStatsResponse, error) {
	deletes, inserts, updates, err := s.prepareStoragesDiff(ctx, request)
	if err != nil {
		return nil, err
	}

	if err := s.applyStorageChanges(ctx, deletes, inserts, updates); err != nil {
		return nil, err
	}

	return &StorageUpsertStatsResponse{
		CountDeleted:  len(deletes),
		CountInserted: len(inserts),
		CountUpdated:  len(updates),
	}, nil
}

func (s *Service) upsertStoragesValue(ctx context.Context, requestData []ProductStorageDto) (*ProductStorageUpsertStatsResponse, error) {
	inserts, updates, err := s.prepareStoragesValueDiff(ctx, requestData)
	if err != nil {
		return nil, err
	}

	inserts, err = s.filterValidStorageUUIDs(ctx, inserts)
	if err != nil {
		return nil, fmt.Errorf("не удалось отфильтровать inserts filterValidStorageUUIDs: %w", err)
	}

	if err := s.applyStorageValueChanges(ctx, inserts, updates); err != nil {
		return nil, err
	}

	return &ProductStorageUpsertStatsResponse{
		CountInserted: len(inserts),
		CountUpdated:  len(updates),
	}, nil
}

func (s *Service) filterValidStorageUUIDs(ctx context.Context, inserts []ProductStorageEnt) ([]ProductStorageEnt, error) {
	if len(inserts) == 0 {
		return inserts, nil
	}

	existing, err := s.storageRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("ошибка Storage не заполнен")
		}
		return nil, fmt.Errorf("ошибка при выполнении storageRepo.GetList: %w", err)
	}

	valid := make(map[uuid.UUID]struct{}, len(existing))
	for _, storage := range existing {
		valid[storage.UUID] = struct{}{}
	}

	var filtered []ProductStorageEnt
	for _, ins := range inserts {
		if _, ok := valid[ins.StorageUUID]; ok {
			filtered = append(filtered, ins)
		} else {
			logger.WarnCtx(ctx, errors.New("foreign key constraint violation avoided ProductStorageEnt"), "", "storage_uuid", ins.StorageUUID, "product_uuid", ins.ProductUUID)
		}
	}

	return filtered, nil
}

func (s *Service) prepareStoragesDiff(ctx context.Context, request UpsertRequest) (deletes, inserts, updates []StorageEnt, err error) {
	existing, err := s.storageRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			if request.General == nil || request.General.Storages == nil {
				return nil, nil, nil, fmt.Errorf("ошибка General должна быть заполнена")
			}
		} else {
			return nil, nil, nil, fmt.Errorf("ошибка при выполнении storageRepo.GetList: %w", err)
		}
	}

	desired := toStorageSlice(request.General.Storages)

	currentMap := toStoragesMap(existing)
	desiredMap := toStoragesMap(desired)

	deletes, inserts, updates = diffStorages(currentMap, desiredMap)
	return deletes, inserts, updates, nil
}

func (s *Service) applyStorageChanges(ctx context.Context, deletes, inserts, updates []StorageEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(deletes) > 0 {
		logger.DebugCtx(ctx, "delete uuid StorageEnt", "uuid", deletes)
		for _, d := range deletes {
			d := d
			g.Go(func() error {
				err := s.storageRepo.Delete(ctx, d.UUID.String())
				if err != nil {
					return fmt.Errorf("ошибка при удалении UUID %s: %w", d.UUID, err)
				}
				return nil
			})
		}
	}

	if len(inserts) > 0 {
		logger.DebugCtx(ctx, "inserts []StorageEnt", "inserts", inserts)
		for _, i := range inserts {
			i := i
			g.Go(func() error {
				_, err := s.storageRepo.Create(ctx, &i)
				if err != nil {
					return fmt.Errorf("ошибка при выполнении propertyRepo.CreateBatch: %w", err)
				}
				return nil
			})
		}
	}

	if len(updates) > 0 {
		logger.DebugCtx(ctx, "updates []StorageEnt", "updates", updates)
		for _, u := range updates {
			u := u
			g.Go(func() error {
				err := s.storageRepo.Update(ctx, &u)
				if err != nil {
					return fmt.Errorf("ошибка при выполнении propertyRepo.UpdateBatch: %w", err)
				}
				return nil
			})
		}
	}

	return g.Wait()
}

func (s *Service) applyStorageValueChanges(ctx context.Context, inserts, updates []ProductStorageEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(inserts) > 0 {
		g.Go(func() error {
			err := s.productStorageRepo.CreateBatch(ctx, inserts)
			logger.DebugCtx(ctx, "inserts []ProductStorageEnt", "inserts", inserts)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении productStorageRepo.CreateBatch: %w", err)
			}
			return nil
		})
	}

	if len(updates) > 0 {
		g.Go(func() error {
			err := s.productStorageRepo.UpdateBatch(ctx, updates)
			logger.DebugCtx(ctx, "updates []ProductStorageEnt", "updates", updates)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении productStorageRepo.UpdateBatch: %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (s *Service) prepareStoragesValueDiff(ctx context.Context, requestData []ProductStorageDto) (inserts, updates []ProductStorageEnt, err error) {
	desiredMap := toProductStorageDataMap(requestData)

	productUUIDs := helper.GetKeys(desiredMap)

	existing, err := s.productStorageRepo.GetByProductUUIDs(ctx, productUUIDs)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, nil, fmt.Errorf("ошибка при выполнении productStorageRepo.GetByProductUUIDs: %w", err)
	}

	currentMap := toProductStorageMap(existing)

	inserts, updates = diffProductStorages(currentMap, desiredMap)
	return inserts, updates, nil
}

func toStoragesMap(list []StorageEnt) map[uuid.UUID]StorageEnt {
	result := make(map[uuid.UUID]StorageEnt)
	for _, p := range list {
		result[p.UUID] = p
	}
	return result
}

func toProductStorageMap(list []ProductStorageEnt) map[uuid.UUID]map[uuid.UUID]ProductStorageEnt {
	result := make(map[uuid.UUID]map[uuid.UUID]ProductStorageEnt)
	for _, storage := range list {
		if _, ok := result[storage.ProductUUID]; !ok {
			result[storage.ProductUUID] = make(map[uuid.UUID]ProductStorageEnt)
		}
		result[storage.ProductUUID][storage.StorageUUID] = storage
	}
	return result
}

func toProductStorageDataMap(list []ProductStorageDto) map[uuid.UUID]map[uuid.UUID]ProductStorageEnt {
	result := make(map[uuid.UUID]map[uuid.UUID]ProductStorageEnt)
	for _, p := range list {
		prdEntMap := make(map[uuid.UUID]ProductStorageEnt)
		productUUID := p.ProductUUID
		for _, pl := range p.ProductStorages {
			ent := *pl.ToEntity()
			ent.ProductUUID = productUUID
			prdEntMap[pl.StorageUUID] = ent
		}
		result[p.ProductUUID] = prdEntMap
	}
	return result
}

func diffStorages(current, desired map[uuid.UUID]StorageEnt) (deletes, inserts, updates []StorageEnt) {
	for key, curr := range current {
		if _, ok := desired[key]; !ok {
			deletes = append(deletes, curr)
		}
	}

	for key, want := range desired {
		curr, ok := current[key]
		if ok {
			if !isStorageEqual(curr, want) {
				updates = append(updates, want)
			}
		} else {
			inserts = append(inserts, want)
		}
	}

	return
}

func isStorageEqual(a, b StorageEnt) bool {
	return a.Name == b.Name && a.Active == b.Active
}

func diffProductStorages(current, desired map[uuid.UUID]map[uuid.UUID]ProductStorageEnt) (inserts, updates []ProductStorageEnt) {
	for keyPrd, dsrStorageMap := range desired {
		currStorageMap, ok := current[keyPrd]
		if !ok {
			prd := helper.GetValues(dsrStorageMap)
			inserts = append(inserts, prd...)
			continue
		}

		// CASCADE удаление
		// for currKeyStorage, currStorage := range currStorageMap {
		// 	if _, ok := dsrStorageMap[currKeyStorage]; !ok {
		// 		deletes = append(deletes, currStorage)
		// 	}
		// }

		for dsrKeyStorage, dsrStorage := range dsrStorageMap {
			currStorage, ok := currStorageMap[dsrKeyStorage]
			if !ok {
				inserts = append(inserts, dsrStorage)
				continue
			}

			if !isProductStorageEqual(currStorage, dsrStorage) {
				updates = append(updates, dsrStorage)
			}
		}
	}
	return
}

func isProductStorageEqual(a, b ProductStorageEnt) bool {
	return a.Quantity == b.Quantity && a.Active == b.Active
}

func toStorageSlice(dts []StorageDto) []StorageEnt {
	desired := make([]StorageEnt, 0, len(dts))
	for _, dto := range dts {
		desired = append(desired, *dto.ToEntity())
	}
	return desired
}

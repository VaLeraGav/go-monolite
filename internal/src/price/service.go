package price

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
	typePriceRepo    *TypePriceRepository
	productPriceRepo *ProductPricesRepository
}

func NewService(typePriceRepo *TypePriceRepository, productPriceRepo *ProductPricesRepository) *Service {
	return &Service{typePriceRepo, productPriceRepo}
}

func (s *Service) GetTypePrice(ctx context.Context) ([]TypePriceResponse, string, error) {
	existing, err := s.typePriceRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, "типы цен не найден", store.ErrNotFound
		}
		return nil, "произошла ошибка при получении типы цен", err
	}

	return helper.ToResponse(existing), "", nil
}

func (s *Service) Upsert(ctx context.Context, request UpsertRequest) (*UpsertResponse, string, error) {
	if err := request.Validate(); err != nil {
		return nil, "", err
	}

	tx, err := s.typePriceRepo.store.Db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	txCtx := store.WithTx(ctx, tx)

	typePriceResponse, err := s.upsertTypePrices(txCtx, request)
	if err != nil {
		return nil, "произошла ошибка при создании типа цены", err
	}

	productPriceResponse, err := s.upsertPricesValue(txCtx, request.ProductPrices)
	if err != nil {
		return nil, "произошла ошибка при создании цен у товаров", err
	}

	if err := tx.Commit(); err != nil {
		return nil, "", fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &UpsertResponse{
		TypePrice:    typePriceResponse,
		ProductPrice: productPriceResponse,
	}, "", nil
}

func (s *Service) upsertTypePrices(ctx context.Context, request UpsertRequest) (*TypePriceResponseDetails, error) {
	deletes, inserts, updates, err := s.prepareTypePricesDiff(ctx, request)
	if err != nil {
		return nil, err
	}

	if err := s.applyTypePriceChanges(ctx, deletes, inserts, updates); err != nil {
		return nil, err
	}

	return &TypePriceResponseDetails{
		CountDeleted:  len(deletes),
		CountInserted: len(inserts),
		CountUpdated:  len(updates),
	}, nil
}

func (s *Service) upsertPricesValue(ctx context.Context, requestData []ProductPriceDto) (*ProductPriceResponseDetails, error) {
	inserts, updates, err := s.preparePricesValueDiff(ctx, requestData)
	if err != nil {
		return nil, err
	}

	inserts, err = s.filterValidPriceUUIDs(ctx, inserts)
	if err != nil {
		return nil, fmt.Errorf("не удалось отфильтровать inserts filterValidPriceUUIDs: %w", err)
	}

	if err := s.applyPriceValueChanges(ctx, inserts, updates); err != nil {
		return nil, err
	}

	return &ProductPriceResponseDetails{
		CountInserted: len(inserts),
		CountUpdated:  len(updates),
	}, nil
}

func (s *Service) filterValidPriceUUIDs(ctx context.Context, inserts []ProductPriceEnt) ([]ProductPriceEnt, error) {
	if len(inserts) == 0 {
		return inserts, nil
	}

	existing, err := s.typePriceRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, fmt.Errorf("ошибка Price не заполнен")
		}
		return nil, fmt.Errorf("ошибка при выполнении в filterValidPriceUUIDs typePriceRepo.GetList: %w", err)
	}

	valid := make(map[uuid.UUID]struct{}, len(existing))
	for _, price := range existing {
		valid[price.UUID] = struct{}{}
	}

	var filtered []ProductPriceEnt
	for _, ins := range inserts {
		if _, ok := valid[ins.TypePriceUUID]; ok {
			filtered = append(filtered, ins)
		} else {
			logger.WarnCtx(ctx, errors.New("foreign key constraint violation avoided ProductPriceEnt"), "", "price_uuid", ins.TypePriceUUID, "product_uuid", ins.ProductUUID)
		}
	}

	return filtered, nil
}

func (s *Service) prepareTypePricesDiff(ctx context.Context, request UpsertRequest) (deletes, inserts, updates []TypePriceEnt, err error) {
	existing, err := s.typePriceRepo.GetList(ctx)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			if request.General == nil || (request.General.Prices == nil) {
				return nil, nil, nil, fmt.Errorf("ошибка General должна быть заполнена")
			}
		} else {
			return nil, nil, nil, fmt.Errorf("ошибка при выполнении typePriceRepo.GetList: %w", err)
		}
	}

	desired := toTypePriceSlice(request.General.Prices)

	currentMap := toTypePricesMap(existing)
	desiredMap := toTypePricesMap(desired)

	deletes, inserts, updates = diffTypePrices(currentMap, desiredMap)
	return deletes, inserts, updates, nil
}

func (s *Service) applyTypePriceChanges(ctx context.Context, deletes, inserts, updates []TypePriceEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(deletes) > 0 {
		logger.DebugCtx(ctx, "delete uuid TypePriceEnt", "uuid", deletes)
		for _, d := range deletes {
			d := d
			g.Go(func() error {
				err := s.typePriceRepo.Delete(ctx, d.UUID.String())
				if err != nil {
					return fmt.Errorf("ошибка при удалении UUID %s: %w", d.UUID, err)
				}
				return nil
			})
		}
	}

	if len(inserts) > 0 {
		logger.DebugCtx(ctx, "inserts []TypePriceEnt", "inserts", inserts)
		for _, i := range inserts {
			i := i
			g.Go(func() error {
				_, err := s.typePriceRepo.Create(ctx, &i)
				if err != nil {
					return fmt.Errorf("ошибка при выполнении typePriceRepo.Create: %w", err)
				}
				return nil
			})
		}
	}

	if len(updates) > 0 {
		logger.DebugCtx(ctx, "updates []PriceEnt", "updates", updates)
		for _, u := range updates {
			u := u
			g.Go(func() error {
				err := s.typePriceRepo.Update(ctx, &u)
				if err != nil {
					return fmt.Errorf("ошибка при выполнении typePriceRepo.Update: %w", err)
				}
				return nil
			})
		}
	}

	return g.Wait()
}

func (s *Service) applyPriceValueChanges(ctx context.Context, inserts, updates []ProductPriceEnt) error {
	g, ctx := errgroup.WithContext(ctx)

	if len(inserts) > 0 {
		g.Go(func() error {
			err := s.productPriceRepo.CreateBatch(ctx, inserts)
			logger.DebugCtx(ctx, "inserts []ProductPriceEnt", "inserts", inserts)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении productPriceRepo.CreateBatch: %w", err)
			}
			return nil
		})
	}

	if len(updates) > 0 {
		g.Go(func() error {
			err := s.productPriceRepo.UpdateBatch(ctx, updates)
			logger.DebugCtx(ctx, "updates []ProductPriceEnt", "updates", updates)
			if err != nil {
				return fmt.Errorf("ошибка при выполнении productPriceRepo.UpdateBatch: %w", err)
			}
			return nil
		})
	}

	return g.Wait()
}

func (s *Service) preparePricesValueDiff(ctx context.Context, requestData []ProductPriceDto) (inserts, updates []ProductPriceEnt, err error) {
	desiredMap := toProductPriceDataMap(requestData)

	productUUIDs := helper.GetKeys(desiredMap)

	existing, err := s.productPriceRepo.GetByProductUUIDs(ctx, productUUIDs)
	if err != nil && !errors.Is(err, store.ErrNotFound) {
		return nil, nil, fmt.Errorf("ошибка при выполнении productPriceRepo.GetByProductUUIDs: %w", err)
	}

	currentMap := toProductPriceMap(existing)

	inserts, updates = diffProductPrices(currentMap, desiredMap)
	return inserts, updates, nil
}

func toTypePricesMap(list []TypePriceEnt) map[uuid.UUID]TypePriceEnt {
	result := make(map[uuid.UUID]TypePriceEnt)
	for _, p := range list {
		result[p.UUID] = p
	}
	return result
}

func toProductPriceMap(list []ProductPriceEnt) map[uuid.UUID]map[uuid.UUID]ProductPriceEnt {
	result := make(map[uuid.UUID]map[uuid.UUID]ProductPriceEnt)
	for _, price := range list {
		if _, ok := result[price.ProductUUID]; !ok {
			result[price.ProductUUID] = make(map[uuid.UUID]ProductPriceEnt)
		}
		result[price.ProductUUID][price.TypePriceUUID] = price
	}
	return result
}

func toProductPriceDataMap(list []ProductPriceDto) map[uuid.UUID]map[uuid.UUID]ProductPriceEnt {
	result := make(map[uuid.UUID]map[uuid.UUID]ProductPriceEnt)
	for _, p := range list {
		prdEntMap := make(map[uuid.UUID]ProductPriceEnt)
		productUUID := p.ProductUUID
		for _, pl := range p.ProductPrices {
			ent := *pl.ToEntity()
			ent.ProductUUID = productUUID
			prdEntMap[pl.TypePriceUUID] = ent
		}
		result[p.ProductUUID] = prdEntMap
	}
	return result
}

func diffTypePrices(current, desired map[uuid.UUID]TypePriceEnt) (deletes, inserts, updates []TypePriceEnt) {
	for key, curr := range current {
		if _, ok := desired[key]; !ok {
			deletes = append(deletes, curr)
		}
	}
	for key, want := range desired {
		curr, ok := current[key]
		if ok {
			if !isPriceEqual(curr, want) {
				updates = append(updates, want)
			}
		} else {
			inserts = append(inserts, want)
		}
	}
	return
}

func isPriceEqual(a, b TypePriceEnt) bool {
	return a.Name == b.Name && a.Active == b.Active
}

func diffProductPrices(current, desired map[uuid.UUID]map[uuid.UUID]ProductPriceEnt) (inserts, updates []ProductPriceEnt) {
	for keyPrd, dsrPriceMap := range desired {
		currPriceMap, ok := current[keyPrd]
		if !ok {
			prd := helper.GetValues(dsrPriceMap)
			inserts = append(inserts, prd...)
			continue
		}

		// CASCADE удаление
		// for currKeyPrice, currPrice := range currPriceMap {
		// 	if _, ok := dsrPriceMap[currKeyPrice]; !ok {
		// 		deletes = append(deletes, currPrice)
		// 	}
		// }

		for dsrKeyPrice, dsrPrice := range dsrPriceMap {
			currPrice, ok := currPriceMap[dsrKeyPrice]
			if !ok {
				inserts = append(inserts, dsrPrice)
				continue
			}

			if !isProductPriceEqual(currPrice, dsrPrice) {
				updates = append(updates, dsrPrice)
			}
		}
	}
	return
}

func isProductPriceEqual(a, b ProductPriceEnt) bool {
	return a.Price == b.Price && a.Active == b.Active
}

func toTypePriceSlice(dts []TypePriceRequest) []TypePriceEnt {
	desired := make([]TypePriceEnt, 0, len(dts))
	for _, dto := range dts {
		desired = append(desired, *dto.ToEntity())
	}
	return desired
}

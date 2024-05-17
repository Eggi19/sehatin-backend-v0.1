package usecases

import (
	"context"
	"fmt"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type StockHistoryUsecaseOpts struct {
	StockHistoryRepo repositories.StockHistoryRepository
}

type StockHistoryUsecase interface {
	GetStockHistoriesByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.StockHistoryParams) ([]entities.StockHistory, *entities.PaginationInfo, error)
}

type StockHistoryUsecaseImpl struct {
	StockHistoryRepo repositories.StockHistoryRepository
}

func NewStockHistoryUsecaseImpl(sthOpts *StockHistoryUsecaseOpts) StockHistoryUsecase {
	return &StockHistoryUsecaseImpl{
		StockHistoryRepo: sthOpts.StockHistoryRepo,
	}
}

func (u *StockHistoryUsecaseImpl) GetStockHistoriesByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.StockHistoryParams) ([]entities.StockHistory, *entities.PaginationInfo, error) {
	stockHistories, totalData, err := u.StockHistoryRepo.FindStockHistoriesByPharmacyId(ctx, pharmacyId, pharmacyManagerId, params)
	if err != nil {
		return nil, nil, err
	}

	for _, sh := range stockHistories {
		if sh.Pharmacy.PharmacyManager.Id != pharmacyManagerId {
			return nil, nil, custom_errors.Unauthorized(err, fmt.Sprintf(`pharmacy number %d not your own`, sh.Pharmacy.Id))
		}
	}

	totalPage := totalData / params.Limit
	if totalData%params.Limit > 0 {
		totalPage++
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return stockHistories, &pagination, nil
}

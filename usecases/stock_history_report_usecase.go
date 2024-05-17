package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type StockHistoryReportUsecaseOpts struct {
	StockHistoryReportRepo repositories.StockHistoryReportRepository
}

type StockHistoryReportUsecase interface {
	GetStockHistories(ctx context.Context, params entities.StockHistoryReportParams) ([]entities.StockHistoryReport, *entities.PaginationInfo, error)
}

type StockHistoryReportUsecaseImpl struct {
	StockHistoryReportRepository repositories.StockHistoryReportRepository
}

func NewStockHistoryReportUsecaseImpl(shuOpts *StockHistoryReportUsecaseOpts) StockHistoryReportUsecase {
	return &StockHistoryReportUsecaseImpl{
		StockHistoryReportRepository: shuOpts.StockHistoryReportRepo,
	}
}

func (u *StockHistoryReportUsecaseImpl) GetStockHistories(ctx context.Context, params entities.StockHistoryReportParams) ([]entities.StockHistoryReport, *entities.PaginationInfo, error) {
	stockReports, totalData, err := u.StockHistoryReportRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
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

	return stockReports, &pagination, nil
}

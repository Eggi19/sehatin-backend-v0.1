package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type SalesReportCategoryUsecaseOpts struct {
	SalesReportCategoryRepo repositories.SalesReportCategoryRepository
}

type SalesReportCategoryUsecase interface {
	GetStockHistories(ctx context.Context, params entities.SalesReportCategoryParams) ([]entities.SalesReportCategory, *entities.PaginationInfo, error)
}

type SalesReportCategoryUsecaseImpl struct {
	SalesReportCategoryRepository repositories.SalesReportCategoryRepository
}

func NewSalesReportCategoryUsecaseImpl(sruOpts *SalesReportCategoryUsecaseOpts) SalesReportCategoryUsecase {
	return &SalesReportCategoryUsecaseImpl{
		SalesReportCategoryRepository: sruOpts.SalesReportCategoryRepo,
	}
}

func (u *SalesReportCategoryUsecaseImpl) GetStockHistories(ctx context.Context, params entities.SalesReportCategoryParams) ([]entities.SalesReportCategory, *entities.PaginationInfo, error) {
	salesReportsCategory, totalData, err := u.SalesReportCategoryRepository.FindAll(ctx, params)
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

	return salesReportsCategory, &pagination, nil
}

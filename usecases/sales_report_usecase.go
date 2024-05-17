package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type SalesReportUsecaseOpts struct {
	SalesReportRepo repositories.SalesReportRepository
	CategoryRepo    repositories.CategoryRepository
}

type SalesReportUsecase interface {
	GetSalesReports(ctx context.Context, params entities.SalesReportParams) ([]entities.SalesReport, *entities.PaginationInfo, error)
}

type SalesReportUsecaseImpl struct {
	SalesReportRepository repositories.SalesReportRepository
	CategoryRepository    repositories.CategoryRepository
}

func NewSalesReportUsecaseImpl(sruOpts *SalesReportUsecaseOpts) SalesReportUsecase {
	return &SalesReportUsecaseImpl{
		SalesReportRepository: sruOpts.SalesReportRepo,
		CategoryRepository:    sruOpts.CategoryRepo,
	}
}

func (u *SalesReportUsecaseImpl) GetSalesReports(ctx context.Context, params entities.SalesReportParams) ([]entities.SalesReport, *entities.PaginationInfo, error) {
	salesReports, totalData, err := u.SalesReportRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(salesReports); i++ {
		categories, err := u.CategoryRepository.GetProductCategory(ctx, salesReports[i].PharmacyProduct.Product.Id)
		if err != nil {
			return nil, nil, err
		}
		salesReports[i].PharmacyProduct.Product.Categories = categories
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

	return salesReports, &pagination, nil
}

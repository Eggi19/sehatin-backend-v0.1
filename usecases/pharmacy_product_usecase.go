package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type PharmacyProductUsecaseOpts struct {
	PharmacyProductRepository repositories.PharmacyProductRepository
	PharmacyRepository        repositories.PharmacyRepository
	CategoryRepository        repositories.CategoryRepository
	ShippingMethodUsecase     ShippingMethodUsecase
	Transactor                repositories.Transactor
	ProductRepository         repositories.ProductRepository
	StockHistoryRepository    repositories.StockHistoryRepository
}

type PharmacyProductUsecase interface {
	GetNearestPharmacyProducts(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error)
	GetNearestPharmacyProductsWithTransaction(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error)
	GetProductDetail(ctx context.Context, req entities.PharmacyProductDetailParams) (*dtos.ProductDetail, error)
	GetProductDetailWithTransaction(ctx context.Context, req entities.PharmacyProductDetailParams) (*dtos.ProductDetail, error)
	CreatePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct, pharmacyManagerId int64) error
	UpdatePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct, pharmacyManagerId int64) error
	DeletePharmacyProduct(ctx context.Context, pharmacyProductId int64, pharmacyManagerId int64) error
	GetPharmacyProductsByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.PharmacyProductParams) ([]entities.PharmacyProduct, *entities.PaginationInfo, error)
	GetAllNearestPharmacyProducts(ctx context.Context, req entities.NearestPharmacyProductsParams) ([]dtos.ProductResponse, *entities.PaginationInfo, error)
	GetPharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error)
}

type PharmacyProductUsecaseImpl struct {
	PharmacyProductRepository repositories.PharmacyProductRepository
	PharmacyRepository        repositories.PharmacyRepository
	CategoryRepository        repositories.CategoryRepository
	ShippingMethodUsecase     ShippingMethodUsecase
	Transactor                repositories.Transactor
	ProductRepository         repositories.ProductRepository
	StockHistoryRepository    repositories.StockHistoryRepository
}

func NewPharmacyProductUsecaseImpl(productOpts *PharmacyProductUsecaseOpts) PharmacyProductUsecase {
	return &PharmacyProductUsecaseImpl{
		PharmacyProductRepository: productOpts.PharmacyProductRepository,
		PharmacyRepository:        productOpts.PharmacyRepository,
		CategoryRepository:        productOpts.CategoryRepository,
		ShippingMethodUsecase:     productOpts.ShippingMethodUsecase,
		Transactor:                productOpts.Transactor,
		ProductRepository:         productOpts.ProductRepository,
		StockHistoryRepository:    productOpts.StockHistoryRepository,
	}
}

func (u *PharmacyProductUsecaseImpl) GetNearestPharmacyProducts(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error) {
	pharmacy, err := u.PharmacyRepository.FindNearestPharmacy(ctx, req)
	if err != nil {
		return nil, err
	}

	if pharmacy == nil {
		products := []dtos.ProductResponse{}
		return products, nil
	}

	products, err := u.PharmacyProductRepository.GetPharmacyProduct(ctx, pharmacy.Id, pagination, req)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *PharmacyProductUsecaseImpl) GetNearestPharmacyProductsWithTransaction(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error) {
	var products []dtos.ProductResponse
	var err error

	_, err = u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		products, err = u.GetNearestPharmacyProducts(txCtx, req, pagination)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (u *PharmacyProductUsecaseImpl) GetProductDetail(ctx context.Context, req entities.PharmacyProductDetailParams) (*dtos.ProductDetail, error) {
	detail, err := u.PharmacyProductRepository.GetProductDetail(ctx, req.PharmacyProductId)
	if err != nil {
		return nil, err
	}

	pharmacy, err := u.PharmacyRepository.FindPharmacyByProduct(ctx, req.Coordinat)
	if err != nil {
		return nil, err
	}
	detail.Pharmacy = pharmacy
	for key, phar := range pharmacy {
		shippingMethodRes, err := u.ShippingMethodUsecase.GetShippingMethod(ctx, phar.Id)
		if err != nil {
			return nil, err
		}

		pharmacy[key].ShippingMethods = *shippingMethodRes
	}

	category, err := u.CategoryRepository.GetProductCategory(ctx, req.Coordinat.ProductId)
	if err != nil {
		return nil, err
	}
	detail.Category = dtos.ConvertToCategoryResponsesWithoutPagination(category)

	return detail, nil
}

func (u *PharmacyProductUsecaseImpl) GetProductDetailWithTransaction(ctx context.Context, req entities.PharmacyProductDetailParams) (*dtos.ProductDetail, error) {
	var detail *dtos.ProductDetail
	var err error

	_, err = u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		detail, err = u.GetProductDetail(txCtx, req)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return detail, nil
}

func (u *PharmacyProductUsecaseImpl) CreatePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct, pharmacyManagerId int64) error {
	_, err := u.ProductRepository.FindOneById(ctx, pp.Product.Id)
	if err != nil {
		return err
	}

	pharmacy, err := u.PharmacyRepository.FindOneById(ctx, pp.Pharmacy.Id)
	if err != nil {
		return err
	}

	if pharmacy.PharmacyManager.Id != pharmacyManagerId {
		return custom_errors.Forbidden()
	}

	foundedPP, _ := u.PharmacyProductRepository.FindOneByPharmacyAndProductId(ctx, pp.Pharmacy.Id, pp.Product.Id)
	if foundedPP != nil && foundedPP.Pharmacy.Id == pp.Pharmacy.Id {
		return custom_errors.BadRequest(err, "already registered")
	}

	ppId, err := u.PharmacyProductRepository.CreateOnePharmacyProduct(ctx, pp)
	if err != nil {
		return err
	}

	err = u.StockHistoryRepository.CreateOne(ctx, entities.StockHistory{
		PharmacyProduct: entities.PharmacyProduct{Id: *ppId},
		Pharmacy:        entities.Pharmacy{Id: pp.Pharmacy.Id},
		Quantity:        pp.TotalStock,
		Description:     "",
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyProductUsecaseImpl) UpdatePharmacyProduct(ctx context.Context, pp entities.PharmacyProduct, pharmacyManagerId int64) error {
	foundedPP, err := u.PharmacyProductRepository.GetOnePharmacyProduct(ctx, pp.Id)
	if err != nil {
		return err
	}

	pharmacy, err := u.PharmacyRepository.FindOneById(ctx, foundedPP.Pharmacy.Id)
	if err != nil {
		return err
	}

	if pharmacy.PharmacyManager.Id != pharmacyManagerId {
		return custom_errors.Forbidden()
	}

	stockHistoryStock := pp.TotalStock - foundedPP.TotalStock

	err = u.PharmacyProductRepository.UpdateOnePharmacyProduct(ctx, pp)
	if err != nil {
		return err
	}

	err = u.StockHistoryRepository.CreateOne(ctx, entities.StockHistory{
		PharmacyProduct: entities.PharmacyProduct{Id: foundedPP.Id},
		Pharmacy:        entities.Pharmacy{Id: pp.Pharmacy.Id},
		Quantity:        stockHistoryStock,
		Description:     "",
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyProductUsecaseImpl) DeletePharmacyProduct(ctx context.Context, pharmacyProductId int64, pharmacyManagerId int64) error {
	pp, err := u.PharmacyProductRepository.GetOnePharmacyProduct(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	pharmacy, err := u.PharmacyRepository.FindOneById(ctx, pp.Pharmacy.Id)
	if err != nil {
		return err
	}

	if pharmacy.PharmacyManager.Id != pharmacyManagerId {
		return custom_errors.Forbidden()
	}

	err = u.PharmacyProductRepository.DeleteOnePharmacyProduct(ctx, pharmacyProductId)
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyProductUsecaseImpl) GetPharmacyProductsByPharmacyId(ctx context.Context, pharmacyId, pharmacyManagerId int64, params entities.PharmacyProductParams) ([]entities.PharmacyProduct, *entities.PaginationInfo, error) {
	pharmacyProducts, totalData, err := u.PharmacyProductRepository.FindPharmacyProductsByPharmacyId(ctx, pharmacyId, pharmacyManagerId, params)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(pharmacyProducts); i++ {
		categories, err := u.CategoryRepository.GetProductCategory(ctx, pharmacyProducts[i].Product.Id)
		if err != nil {
			return nil, nil, err
		}
		pharmacyProducts[i].Product.Categories = categories
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

	return pharmacyProducts, &pagination, nil
}

func (u *PharmacyProductUsecaseImpl) GetAllNearestPharmacyProducts(ctx context.Context, params entities.NearestPharmacyProductsParams) ([]dtos.ProductResponse, *entities.PaginationInfo, error) {
	pharmacyProducts, totalData, err := u.PharmacyProductRepository.FindAllNearest(ctx, params)
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

	return pharmacyProducts, &pagination, nil
}

func (u *PharmacyProductUsecaseImpl) GetPharmacyProduct(ctx context.Context, id int64) (*entities.PharmacyProduct, error) {
	pp, err := u.PharmacyProductRepository.FindOnePharmacyProduct(ctx, id)
	if err != nil {
		return nil, err
	}

	categories := []entities.Category{}

	for i := 0; i < len(pp.Product.Categories); i++ {
		category, err := u.CategoryRepository.FindById(ctx, pp.Product.Categories[i].Id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *category)
	}

	pp.Product.Categories = categories

	return pp, nil
}

package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type ProductUsecaseOpts struct {
	ProductRepo         repositories.ProductRepository
	CategoryRepos       repositories.CategoryRepository
	ProductCategoryRepo repositories.ProductCategoryRepository
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product entities.Product) (int64, error)
	UpdateProduct(ctx context.Context, product entities.Product) error
	DeleteProduct(ctx context.Context, productId int64) error
	GetOneProduct(ctx context.Context, productId int64) (*entities.Product, error)
	GetAllProduct(ctx context.Context, params entities.ProductCategoryParams) ([]entities.Product, *entities.PaginationInfo, error)
}

type ProductUsecaseImpl struct {
	ProductRepository         repositories.ProductRepository
	CategoryRepository        repositories.CategoryRepository
	ProductCategoryRepository repositories.ProductCategoryRepository
}

func NewProductUsecaseImpl(pUpts *ProductUsecaseOpts) ProductUsecase {
	return &ProductUsecaseImpl{
		ProductRepository:         pUpts.ProductRepo,
		CategoryRepository:        pUpts.CategoryRepos,
		ProductCategoryRepository: pUpts.ProductCategoryRepo,
	}
}

func (u *ProductUsecaseImpl) CreateProduct(ctx context.Context, product entities.Product) (int64, error) {
	productId, err := u.ProductRepository.CreateOne(ctx, product)
	if err != nil {
		return 0, err
	}

	return *productId, nil
}

func (u *ProductUsecaseImpl) UpdateProduct(ctx context.Context, product entities.Product) error {
	err := u.ProductRepository.UpdateOne(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductUsecaseImpl) DeleteProduct(ctx context.Context, productId int64) error {
	err := u.ProductRepository.DeleteOne(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductUsecaseImpl) GetOneProduct(ctx context.Context, productId int64) (*entities.Product, error) {
	product, err := u.ProductRepository.FindOneById(ctx, productId)
	if err != nil {
		return nil, err
	}

	categories, err := u.CategoryRepository.GetProductCategory(ctx, productId)
	if err != nil {
		return nil, err
	}

	product.Categories = categories

	return product, nil
}

func (u *ProductUsecaseImpl) GetAllProduct(ctx context.Context, params entities.ProductCategoryParams) ([]entities.Product, *entities.PaginationInfo, error) {
	products, totalData, err := u.ProductRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(products); i++ {
		categories, err := u.CategoryRepository.GetProductCategory(ctx, products[i].Id)
		if err != nil {
			return nil, nil, err
		}
		products[i].Categories = categories
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

	return products, &pagination, nil
}

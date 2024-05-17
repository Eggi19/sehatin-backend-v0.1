package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type ProductCategoryUsecaseOpts struct {
	ProductCategoryRepo repositories.ProductCategoryRepository
}

type ProductCategoryUsecase interface {
	CreateProductCategory(ctx context.Context, pc entities.ProductCategory) error
	UpdateProductCategory(ctx context.Context, pc entities.ProductCategory) error
	DeleteProductCategoryByProductId(ctx context.Context, productId int64) error
	FindProductCategoryByProductId(ctx context.Context, productId int64) ([]entities.ProductCategory, error)
}

type ProductCategoryUsecaseImpl struct {
	ProductCategoryRepoitory repositories.ProductCategoryRepository
}

func NewProductCategoryUsecaseImpl(pcuOpts *ProductCategoryUsecaseOpts) ProductCategoryUsecase {
	return &ProductCategoryUsecaseImpl{
		ProductCategoryRepoitory: pcuOpts.ProductCategoryRepo,
	}
}

func (u *ProductCategoryUsecaseImpl) CreateProductCategory(ctx context.Context, pc entities.ProductCategory) error {
	err := u.ProductCategoryRepoitory.CreateOne(ctx, pc)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductCategoryUsecaseImpl) UpdateProductCategory(ctx context.Context, pc entities.ProductCategory) error {
	pcs, err := u.ProductCategoryRepoitory.FindAllByProductId(ctx, pc.ProductId)
	if err != nil {
		return err
	}

	for i := 0; i < len(pcs); i++ {
		pcs[i].CategoryId = pc.ProductId
		err := u.ProductCategoryRepoitory.UpdateOne(ctx, pc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ProductCategoryUsecaseImpl) DeleteProductCategoryByProductId(ctx context.Context, productId int64) error {
	err := u.ProductCategoryRepoitory.DeleteOneByProductId(ctx, productId)
	if err != nil {
		return err
	}

	return nil
}

func (u *ProductCategoryUsecaseImpl) FindProductCategoryByProductId(ctx context.Context, productId int64) ([]entities.ProductCategory, error) {
	pcs, err := u.ProductCategoryRepoitory.FindAllByProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	return pcs, nil
}

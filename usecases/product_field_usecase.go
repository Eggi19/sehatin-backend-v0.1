package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type ProductFieldUsecaseOpts struct {
	ProductFieldRepo repositories.ProductFieldRepository
}

type ProductFieldUsecase interface {
	GetAllForm(ctx context.Context) ([]entities.ProductForm, error)
	GetAllClassification(ctx context.Context) ([]entities.ProductClassification, error)
	GetAllManufacture(ctx context.Context) ([]entities.Manufacture, error)
}

type ProductFieldUsecaseImpl struct {
	ProductFieldRepository repositories.ProductFieldRepository
}

func NewProductFieldUsecaseImpl(pfuOpts *ProductFieldUsecaseOpts) ProductFieldUsecase {
	return &ProductFieldUsecaseImpl{
		ProductFieldRepository: pfuOpts.ProductFieldRepo,
	}
}

func (u *ProductFieldUsecaseImpl) GetAllForm(ctx context.Context) ([]entities.ProductForm, error) {
	form, err := u.ProductFieldRepository.FindAllForm(ctx)
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (u *ProductFieldUsecaseImpl) GetAllClassification(ctx context.Context) ([]entities.ProductClassification, error) {
	classification, err := u.ProductFieldRepository.FindAllClassification(ctx)
	if err != nil {
		return nil, err
	}

	return classification, nil
}

func (u *ProductFieldUsecaseImpl) GetAllManufacture(ctx context.Context) ([]entities.Manufacture, error) {
	manufacture, err := u.ProductFieldRepository.FindAllManufacture(ctx)
	if err != nil {
		return nil, err
	}

	return manufacture, nil
}

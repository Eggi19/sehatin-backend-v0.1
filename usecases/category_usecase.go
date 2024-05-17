package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type CategoryUsecaseOpts struct {
	CategoryRepo repositories.CategoryRepository
}

type CategoryUsecase interface {
	GetAllCategory(ctx context.Context, params entities.CategoryParams) ([]entities.Category, *entities.PaginationInfo, error)
	GetCategoryById(ctx context.Context, categoryId int64) (*entities.Category, error)
	CreateCategory(ctx context.Context, category entities.Category) error
	UpdateCategory(ctx context.Context, category entities.Category) error
	DeleteCategory(ctx context.Context, categoryId int64) error
}

type CategoryUsecaseImpl struct {
	CategoryRepository repositories.CategoryRepository
}

func NewCategoryUsecaseImpl(cuOpts *CategoryUsecaseOpts) CategoryUsecase {
	return &CategoryUsecaseImpl{
		CategoryRepository: cuOpts.CategoryRepo,
	}
}

func (u *CategoryUsecaseImpl) GetAllCategory(ctx context.Context, params entities.CategoryParams) ([]entities.Category, *entities.PaginationInfo, error) {
	categories, totalData, err := u.CategoryRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	var totalPage int
	if params.Limit != 0 && params.Page != 0 {
		totalPage = totalData / params.Limit
		if totalData%params.Limit > 0 {
			totalPage++
		}
	}

	pagination := entities.PaginationInfo{
		Page:      params.Page,
		Limit:     params.Limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}

	return categories, &pagination, nil
}

func (u *CategoryUsecaseImpl) GetCategoryById(ctx context.Context, categoryId int64) (*entities.Category, error) {
	category, err := u.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *CategoryUsecaseImpl) CreateCategory(ctx context.Context, category entities.Category) error {
	err := u.CategoryRepository.CreateOne(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (u *CategoryUsecaseImpl) UpdateCategory(ctx context.Context, category entities.Category) error {
	_, err := u.CategoryRepository.FindById(ctx, category.Id)
	if err != nil {
		return err
	}

	err = u.CategoryRepository.UpdateOne(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (u *CategoryUsecaseImpl) DeleteCategory(ctx context.Context, categoryId int64) error {
	_, err := u.CategoryRepository.FindById(ctx, categoryId)
	if err != nil {
		return err
	}

	err = u.CategoryRepository.DeleteOne(ctx, categoryId)
	if err != nil {
		return err
	}

	return nil
}

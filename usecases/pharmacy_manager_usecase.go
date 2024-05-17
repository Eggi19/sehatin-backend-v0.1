package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type PharmacyManagerOpts struct {
	PharmacyManagerRepo repositories.PharmacyManagerRepository
}

type PharmacyManagerUsecase interface {
	GetAllPharmacyManager(ctx context.Context, params entities.PharmacyManagerParams) ([]entities.PharmacyManager, *entities.PaginationInfo, error)
	GetPharmacyManagerById(ctx context.Context, pharmacyManagerId int64) (*entities.PharmacyManager, error)
	UpdatePharmacyManager(ctx context.Context, pharmacyManager entities.PharmacyManager) error
	DeletePharmacyManager(ctx context.Context, pharmacyManagerId int64) error
}

type PharmacyManagerUsecaseImpl struct {
	PharmacyManagerRepository repositories.PharmacyManagerRepository
}

func NewPharmacyManagerUsecaseImpl(pmuOpts *PharmacyManagerOpts) PharmacyManagerUsecase {
	return &PharmacyManagerUsecaseImpl{
		PharmacyManagerRepository: pmuOpts.PharmacyManagerRepo,
	}
}

func (u *PharmacyManagerUsecaseImpl) GetAllPharmacyManager(ctx context.Context, params entities.PharmacyManagerParams) ([]entities.PharmacyManager, *entities.PaginationInfo, error) {
	managers, totalData, err := u.PharmacyManagerRepository.FindAll(ctx, params)
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

	return managers, &pagination, nil
}

func (u *PharmacyManagerUsecaseImpl) GetPharmacyManagerById(ctx context.Context, pharmacyManagerId int64) (*entities.PharmacyManager, error) {
	pharmacyManager, err := u.PharmacyManagerRepository.FindOneById(ctx, pharmacyManagerId)
	if err != nil {
		return nil, err
	}

	return pharmacyManager, nil
}

func (u *PharmacyManagerUsecaseImpl) UpdatePharmacyManager(ctx context.Context, pharmacyManager entities.PharmacyManager) error {
	_, err := u.PharmacyManagerRepository.FindOneById(ctx, pharmacyManager.Id)
	if err != nil {
		return err
	}

	err = u.PharmacyManagerRepository.UpdateOne(ctx, pharmacyManager)
	if err != nil {
		return err
	}

	return nil
}

func (u *PharmacyManagerUsecaseImpl) DeletePharmacyManager(ctx context.Context, pharmacyManagerId int64) error {
	err := u.PharmacyManagerRepository.DeleteById(ctx, pharmacyManagerId)
	if err != nil {
		return err
	}

	return nil
}

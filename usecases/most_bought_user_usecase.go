package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type MostBoughtUserUsecaseOpts struct {
	MostBoughtUserRepo  repositories.MostBoughtUserRepository
	CategoryRepo        repositories.CategoryRepository
	PharmacyProductRepo repositories.PharmacyProductRepository
	PharmacyRepo        repositories.PharmacyRepository
}

type MostBoughtUserUsecase interface {
	GetMostBought(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error)
}

type MostBoughtUserUsecaseImpl struct {
	MostBoughtUserRepository  repositories.MostBoughtUserRepository
	CategoryRepository        repositories.CategoryRepository
	PharmacyProductRepository repositories.PharmacyProductRepository
	PharmacyRepository        repositories.PharmacyRepository
}

func NewMostBoughtUserUsecaseImpl(mbuOpts *MostBoughtUserUsecaseOpts) MostBoughtUserUsecase {
	return &MostBoughtUserUsecaseImpl{
		MostBoughtUserRepository:  mbuOpts.MostBoughtUserRepo,
		CategoryRepository:        mbuOpts.CategoryRepo,
		PharmacyProductRepository: mbuOpts.PharmacyProductRepo,
		PharmacyRepository:        mbuOpts.PharmacyRepo,
	}
}

func (u *MostBoughtUserUsecaseImpl) GetMostBought(ctx context.Context, req entities.NearestPharmacyParams, pagination entities.PaginationParams) ([]dtos.ProductResponse, error) {
	pharmacy, err := u.PharmacyRepository.FindNearestPharmacyMostBought(ctx, req)
	if err != nil {
		return nil, err
	}

	if pharmacy == nil {
		products := []dtos.ProductResponse{}
		return products, nil
	}

	mostBought, err := u.MostBoughtUserRepository.FindMostBought(ctx, pharmacy.Id, pagination, req)
	if err != nil {
		return nil, err
	}

	return mostBought, nil
}

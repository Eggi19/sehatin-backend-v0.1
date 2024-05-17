package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type SpecialistUsecaseOpts struct {
	SpecialistRepo repositories.SpecialistRepository
}

type SpecialistUsecase interface {
	GetAllSpecialist(ctx context.Context) ([]entities.DoctorSpecialist, error)
}

type SpecialistUsecaseImpl struct {
	SpecialistRepository repositories.SpecialistRepository
}

func NewSpecialistUsecaseImpl(spUseOpts *SpecialistUsecaseOpts) SpecialistUsecase {
	return &SpecialistUsecaseImpl{
		SpecialistRepository: spUseOpts.SpecialistRepo,
	}
}

func (u *SpecialistUsecaseImpl) GetAllSpecialist(ctx context.Context) ([]entities.DoctorSpecialist, error) {
	specialists, err := u.SpecialistRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return specialists, nil
}

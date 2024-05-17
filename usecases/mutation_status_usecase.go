package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type MutationSatusUsecaseOpts struct {
	MutationStatusRepo repositories.MutationStatusRepository
}

type MutationSatusUsecase interface {
	GetAll(ctx context.Context) ([]entities.MutationSatus, error)
	GetMutationStatusById(ctx context.Context, id int64) (*entities.MutationSatus, error)
}

type MutationSatusUsecaseImpl struct {
	MutationStatusRepository repositories.MutationStatusRepository
}

func NewMutationSatusUsecaseImpl(msuOpts *MutationSatusUsecaseOpts) MutationSatusUsecase {
	return &MutationSatusUsecaseImpl{
		MutationStatusRepository: msuOpts.MutationStatusRepo,
	}
}

func (u *MutationSatusUsecaseImpl) GetAll(ctx context.Context) ([]entities.MutationSatus, error) {
	mutationStatuses, err := u.MutationStatusRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return mutationStatuses, nil
}

func (u *MutationSatusUsecaseImpl) GetMutationStatusById(ctx context.Context, id int64) (*entities.MutationSatus, error) {
	mutationStatus, err := u.MutationStatusRepository.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	return mutationStatus, nil
}

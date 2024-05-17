package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type UserUsecaseOpts struct {
	UserRepo        repositories.UserRepository
	UserAddressRepo repositories.UserAddressRepository
	GenderRepo      repositories.GenderRepository
}

type UserUsecase interface {
	GetUserById(ctx context.Context, userId int64) (*entities.User, error)
	DeleteUser(ctx context.Context, userId int64) error
	GetAllUser(ctx context.Context, params entities.UserParams) ([]entities.User, *entities.PaginationInfo, error)
	UpdateUser(ctx context.Context, user entities.User) error
}

type UserUsecaseImpl struct {
	UserRepository        repositories.UserRepository
	UserAddressRepository repositories.UserAddressRepository
	GenderRepository      repositories.GenderRepository
}

func NewUserUsecaseImpl(uuOpts *UserUsecaseOpts) UserUsecase {
	return &UserUsecaseImpl{
		UserRepository:        uuOpts.UserRepo,
		UserAddressRepository: uuOpts.UserAddressRepo,
		GenderRepository:      uuOpts.GenderRepo,
	}
}

func (u *UserUsecaseImpl) GetUserById(ctx context.Context, userId int64) (*entities.User, error) {
	user, err := u.UserRepository.FindOneById(ctx, userId)
	if err != nil {
		return nil, err
	}

	addresses, err := u.UserAddressRepository.FindAllByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	user.Address = addresses

	return user, nil
}

func (u *UserUsecaseImpl) DeleteUser(ctx context.Context, userId int64) error {
	err := u.UserRepository.DeleteOne(ctx, userId)

	if err != nil {
		return err
	}

	return nil
}

func (u *UserUsecaseImpl) GetAllUser(ctx context.Context, params entities.UserParams) ([]entities.User, *entities.PaginationInfo, error) {
	users, totalData, err := u.UserRepository.FindAll(ctx, params)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(users)-1; i++ {
		addresses, err := u.UserAddressRepository.FindAllByUserId(ctx, users[i].Id)
		if err != nil {
			return nil, nil, err
		}
		users[i].Address = addresses
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

	return users, &pagination, nil
}

func (u *UserUsecaseImpl) UpdateUser(ctx context.Context, user entities.User) error {
	_, err := u.GenderRepository.FindById(ctx, user.Gender.Id)
	if err != nil {
		return err
	}

	err = u.UserRepository.UpdateOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

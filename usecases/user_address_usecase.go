package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
)

type UserAddressUsecaseOpts struct {
	UserAddressRepo repositories.UserAddressRepository
	UserRepo        repositories.UserRepository
}

type UserAddressUsecase interface {
	GetAddressById(ctx context.Context, addressId int64, userId int64) (*entities.UserAddress, error)
	CreateUserAddress(ctx context.Context, userAddress entities.UserAddress) error
	UpdateUserAddress(ctx context.Context, userAddress entities.UserAddress) error
	DeleteUserAddress(ctx context.Context, userAddressId, userId int64) error
}

type UserAddressUsecaseImpl struct {
	UserAddressRepository repositories.UserAddressRepository
	UserRepository        repositories.UserRepository
}

func NewUserAddressUsecaseImpl(uaOpts *UserAddressUsecaseOpts) UserAddressUsecase {
	return &UserAddressUsecaseImpl{
		UserAddressRepository: uaOpts.UserAddressRepo,
		UserRepository:        uaOpts.UserRepo,
	}
}

func (u *UserAddressUsecaseImpl) GetAddressById(ctx context.Context, addressId int64, userId int64) (*entities.UserAddress, error) {
	address, err := u.UserAddressRepository.FindById(ctx, addressId, userId)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (u *UserAddressUsecaseImpl) CreateUserAddress(ctx context.Context, userAddress entities.UserAddress) error {
	id, err := u.UserAddressRepository.CreateOne(ctx, userAddress)
	if err != nil {
		return err
	}

	err = u.UserAddressRepository.UpdateIsMainFalse(ctx, id, userAddress.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserAddressUsecaseImpl) UpdateUserAddress(ctx context.Context, userAddress entities.UserAddress) error {
	address, err := u.UserAddressRepository.FindById(ctx, userAddress.Id, userAddress.UserId)
	if err != nil {
		return err
	}

	if !userAddress.IsMain && address.IsMain {
		listAddress, err := u.UserAddressRepository.FindAllByUserId(ctx, userAddress.UserId)
		if err != nil {
			return err
		}
		err = u.UserAddressRepository.UpdateIsMain(ctx, listAddress[len(listAddress)-1].Id)
		if err != nil {
			return err
		}
		err = u.UserAddressRepository.UpdateOne(ctx, userAddress)
		if err != nil {
			return err
		}
		return nil
	}

	if userAddress.IsMain != address.IsMain {
		listAddress, err := u.UserAddressRepository.FindAllByUserId(ctx, userAddress.UserId)
		if err != nil {
			return err
		}
		for i := 0; i < len(listAddress); i++ {
			if listAddress[i].IsMain {
				err := u.UserAddressRepository.UpdateIsMain(ctx, listAddress[i].Id)
				if err != nil {
					return err
				}
				err = u.UserAddressRepository.UpdateOne(ctx, userAddress)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}

	err = u.UserAddressRepository.UpdateOne(ctx, userAddress)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserAddressUsecaseImpl) DeleteUserAddress(ctx context.Context, userAddressId, userId int64) error {
	userAddress, err := u.UserAddressRepository.FindById(ctx, userAddressId, userId)
	if err != nil {
		return err
	}

	if userAddress.IsMain {
		err = u.UserAddressRepository.DeleteOne(ctx, userAddressId, userId)
		if err != nil {
			return err
		}
		addresses, err := u.UserAddressRepository.FindAllByUserId(ctx, userId)
		if err != nil {
			return err
		}
		err = u.UserAddressRepository.UpdateIsMain(ctx, addresses[0].Id)
		if err != nil {
			return err
		}
		return nil
	}

	err = u.UserAddressRepository.DeleteOne(ctx, userAddressId, userId)
	if err != nil {
		return err
	}

	return nil
}

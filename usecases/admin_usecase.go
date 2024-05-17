package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type AdminUsecaseOpts struct {
	AdminRepository repositories.AdminRepository
	Transactor      repositories.Transactor
	HashAlgorithm   utils.Hasher
}

type AdminUsecase interface {
	CreateAdmin(ctx context.Context, admin entities.Admin) error
	CreateAdminWithTransaction(ctx context.Context, admin entities.Admin) error
	DeleteAdmin(ctx context.Context, id int64) error
	GetAdminById(ctx context.Context, id int64) (*entities.Admin, error)
	GetAllAdmin(ctx context.Context, params entities.AdminParams) ([]entities.Admin, *entities.PaginationInfo, error)
}

type AdminUsecaseImpl struct {
	AdminRepository repositories.AdminRepository
	Transactor      repositories.Transactor
	HashAlgorithm   utils.Hasher
}

func NewAdminUsecaseImpl(auOpts *AdminUsecaseOpts) AdminUsecase {
	return &AdminUsecaseImpl{
		AdminRepository: auOpts.AdminRepository,
		Transactor:      auOpts.Transactor,
		HashAlgorithm:   auOpts.HashAlgorithm,
	}
}

func (u *AdminUsecaseImpl) CreateAdmin(ctx context.Context, admin entities.Admin) error {
	pwd := admin.Password
	pwdHash, err := u.HashAlgorithm.HashPassword(pwd)
	if err != nil {
		return err
	}

	admin.Password = string(pwdHash)

	err = u.AdminRepository.CreateOne(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}

func (u *AdminUsecaseImpl) CreateAdminWithTransaction(ctx context.Context, admin entities.Admin) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.CreateAdmin(ctx, admin)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *AdminUsecaseImpl) DeleteAdmin(ctx context.Context, id int64) error {
	err := u.AdminRepository.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *AdminUsecaseImpl) GetAdminById(ctx context.Context, id int64) (*entities.Admin, error) {
	admin, err := u.AdminRepository.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	return admin, err
}

func (u *AdminUsecaseImpl) GetAllAdmin(ctx context.Context, params entities.AdminParams) ([]entities.Admin, *entities.PaginationInfo, error) {
	admins, totalData, err := u.AdminRepository.FindAll(ctx, params)
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

	return admins, &pagination, nil
}

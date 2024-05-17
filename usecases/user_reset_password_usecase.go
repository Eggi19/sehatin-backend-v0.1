package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type UserResetPasswordUsecaseOpts struct {
	UserResetPasswordRepo repositories.UserResetPasswordRepository
	UserRepo              repositories.UserRepository
	HashAlgorithm         utils.Hasher
	EmailSender           utils.EmailSender
	AuthTokenProvider     utils.AuthTokenProvider
	Transactor            repositories.Transactor
}

type UserResetPasswordUsecase interface {
	ChangePassword(ctx context.Context, userId int64, password string, newPassword string) error
	ChangePasswordWithTransaction(ctx context.Context, userId int64, password string, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ForgotPasswordWithTransaction(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email string, token string, newPassword string) error
	ResetPasswordWitTransaction(ctx context.Context, email string, token string, newPassword string) error
}

type UserResetPasswordUsecaseImpl struct {
	UserResetPasswordRepository repositories.UserResetPasswordRepository
	UserRepository              repositories.UserRepository
	HashAlgorithm               utils.Hasher
	EmailSender                 utils.EmailSender
	AuthTokenProvider           utils.AuthTokenProvider
	Transactor                  repositories.Transactor
}

func NewUserResetPasswordUsecaseImpl(uruOpts *UserResetPasswordUsecaseOpts) UserResetPasswordUsecase {
	return &UserResetPasswordUsecaseImpl{
		UserResetPasswordRepository: uruOpts.UserResetPasswordRepo,
		UserRepository:              uruOpts.UserRepo,
		HashAlgorithm:               uruOpts.HashAlgorithm,
		EmailSender:                 uruOpts.EmailSender,
		AuthTokenProvider:           uruOpts.AuthTokenProvider,
		Transactor:                  uruOpts.Transactor,
	}
}

func (u *UserResetPasswordUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	config, err := utils.ConfigInit()
	if err != nil {
		return err
	}

	user, err := u.UserRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return err
	}

	foundedReset, _ := u.UserResetPasswordRepository.FindOneByUserId(ctx, user.Id)
	if foundedReset != nil {
		u.UserResetPasswordRepository.DeleteOneById(ctx, foundedReset.Id)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = user.Id
	dataTokenMap[constants.UserEmail] = user.Email
	dataTokenMap[constants.Role] = constants.UserRole

	token, err := u.AuthTokenProvider.GenerateResetPasswordToken(dataTokenMap)
	if err != nil {
		return err
	}

	expTime := utils.SetExpire()

	message := fmt.Sprintf("%s%s", config.FrontendUrl, token)
	userReset := entities.UserResetPassword{
		UserId:    user.Id,
		Token:     token,
		ExpiredAt: expTime.ResetPasswordTokenExp,
	}

	_, err = u.UserResetPasswordRepository.CreateOne(ctx, userReset)
	if err != nil {
		return err
	}

	err = u.EmailSender.SendEmail(user.Email, message, expTime.ResetPasswordTokenExp.Format("2006-01-02T15:04:05Z07:00"))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserResetPasswordUsecaseImpl) ForgotPasswordWithTransaction(ctx context.Context, email string) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.ForgotPassword(ctx, email)
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

func (u *UserResetPasswordUsecaseImpl) ResetPassword(ctx context.Context, email string, token string, newPassword string) error {
	userReset, err := u.UserResetPasswordRepository.FindOneByToken(ctx, token)
	if err != nil {
		return err
	}

	if time.Now().Unix() > userReset.ExpiredAt.Unix() {
		return custom_errors.ExpiredResetPwdCode()
	}

	newPasswordHash, err := u.HashAlgorithm.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user, err := u.UserRepository.FindOneById(ctx, userReset.UserId)
	if err != nil {
		return err
	}

	if user.Email != email {
		return custom_errors.InvalidEmail()
	}

	err = u.UserRepository.UpdatePassword(ctx, user.Id, string(newPasswordHash))
	if err != nil {
		return err
	}

	err = u.UserResetPasswordRepository.DeleteOneById(ctx, userReset.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserResetPasswordUsecaseImpl) ResetPasswordWitTransaction(ctx context.Context, email string, token string, newPassword string) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.ResetPassword(ctx, email, token, newPassword)
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

func (u *UserResetPasswordUsecaseImpl) ChangePassword(ctx context.Context, userId int64, password string, newPassword string) error {
	user, err := u.UserRepository.FindPasswordById(ctx, userId)
	if err != nil {
		return err
	}

	true, err := u.HashAlgorithm.CheckPassword(password, []byte(user.Password.String))
	if err != nil {
		return custom_errors.BadRequest(err, constants.PasswordNotMatchErrMsg)
	}

	if true {
		newHashPassword, err := u.HashAlgorithm.HashPassword(newPassword)
		if err != nil {
			return err
		}

		err = u.UserRepository.UpdatePassword(ctx, userId, string(newHashPassword))
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func (u *UserResetPasswordUsecaseImpl) ChangePasswordWithTransaction(ctx context.Context, userId int64, password string, newPassword string) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.ChangePassword(ctx, userId, password, newPassword)
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

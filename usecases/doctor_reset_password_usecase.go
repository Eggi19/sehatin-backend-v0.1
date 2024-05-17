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

type DoctorResetPasswordUsecaseOpts struct {
	DoctorResetPasswordRepo repositories.DoctorResetPasswordRepository
	DoctorRepo              repositories.DoctorRepository
	HashAlgorithm           utils.Hasher
	EmailSender             utils.EmailSender
	AuthTokenProvider       utils.AuthTokenProvider
	Transactor              repositories.Transactor
}

type DoctorResetPasswordUsecase interface {
	ChangePassword(ctx context.Context, doctorId int64, password string, newPassword string) error
	ChangePasswordWithTransaction(ctx context.Context, doctorId int64, password string, newPassword string) error
	ForgotPassword(ctx context.Context, email string) error
	ForgotPasswordWithTransaction(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, email string, token string, newPassword string) error
	ResetPasswordWitTransaction(ctx context.Context, email string, token string, newPassword string) error
}

type DoctorResetPasswordUsecaseImpl struct {
	DoctorResetPasswordRepository repositories.DoctorResetPasswordRepository
	DoctorRepository              repositories.DoctorRepository
	HashAlgorithm                 utils.Hasher
	EmailSender                   utils.EmailSender
	AuthTokenProvider             utils.AuthTokenProvider
	Transactor                    repositories.Transactor
}

func NewDoctorResetPasswordUsecaseImpl(uruOpts *DoctorResetPasswordUsecaseOpts) DoctorResetPasswordUsecase {
	return &DoctorResetPasswordUsecaseImpl{
		DoctorResetPasswordRepository: uruOpts.DoctorResetPasswordRepo,
		DoctorRepository:              uruOpts.DoctorRepo,
		HashAlgorithm:                 uruOpts.HashAlgorithm,
		EmailSender:                   uruOpts.EmailSender,
		AuthTokenProvider:             uruOpts.AuthTokenProvider,
		Transactor:                    uruOpts.Transactor,
	}
}

func (u *DoctorResetPasswordUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	config, err := utils.ConfigInit()
	if err != nil {
		return err
	}

	doctor, err := u.DoctorRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return err
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = doctor.Id
	dataTokenMap[constants.UserEmail] = doctor.Email
	dataTokenMap[constants.Role] = constants.UserRole

	token, err := u.AuthTokenProvider.GenerateResetPasswordToken(dataTokenMap)
	if err != nil {
		return err
	}

	expTime := utils.SetExpire()

	message := fmt.Sprintf("%s%s", config.FrontendUrl, token)
	userReset := entities.UserResetPassword{
		UserId:    doctor.Id,
		Token:     token,
		ExpiredAt: expTime.ResetPasswordTokenExp,
	}

	_, err = u.DoctorResetPasswordRepository.CreateOne(ctx, userReset)
	if err != nil {
		return err
	}

	err = u.EmailSender.SendEmail(doctor.Email, message, expTime.ResetPasswordTokenExp.Format("2006-01-02T15:04:05Z07:00"))
	if err != nil {
		return err
	}

	return nil
}

func (u *DoctorResetPasswordUsecaseImpl) ForgotPasswordWithTransaction(ctx context.Context, email string) error {
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

func (u *DoctorResetPasswordUsecaseImpl) ResetPassword(ctx context.Context, email string, token string, newPassword string) error {
	userReset, err := u.DoctorResetPasswordRepository.FindOneByToken(ctx, token)
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

	doctor, err := u.DoctorRepository.FindOneById(ctx, userReset.UserId)
	if err != nil {
		return err
	}

	if doctor.Email != email {
		return custom_errors.InvalidEmail()
	}

	err = u.DoctorRepository.UpdatePassword(ctx, doctor.Id, string(newPasswordHash))
	if err != nil {
		return err
	}

	err = u.DoctorResetPasswordRepository.DeleteOneById(ctx, userReset.Id)
	if err != nil {
		return err
	}

	return nil
}

func (u *DoctorResetPasswordUsecaseImpl) ResetPasswordWitTransaction(ctx context.Context, email string, token string, newPassword string) error {
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

func (u *DoctorResetPasswordUsecaseImpl) ChangePassword(ctx context.Context, doctorId int64, password string, newPassword string) error {
	doctor, err := u.DoctorRepository.FindPasswordById(ctx, doctorId)
	if err != nil {
		return err
	}

	true, err := u.HashAlgorithm.CheckPassword(password, []byte(doctor.Password.String))
	if err != nil {
		return custom_errors.BadRequest(err, constants.PasswordNotMatchErrMsg)
	}

	if true {
		newHashPassword, err := u.HashAlgorithm.HashPassword(newPassword)
		if err != nil {
			return err
		}

		err = u.DoctorRepository.UpdatePassword(ctx, doctorId, string(newHashPassword))
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func (u *DoctorResetPasswordUsecaseImpl) ChangePasswordWithTransaction(ctx context.Context, doctorId int64, password string, newPassword string) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(ctx context.Context) (interface{}, error) {
		err := u.ChangePassword(ctx, doctorId, password, newPassword)
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

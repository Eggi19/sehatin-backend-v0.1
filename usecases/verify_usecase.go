package usecases

import (
	"context"
	"fmt"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type VerifyUsecaseOpts struct {
	UserRepo          repositories.UserRepository
	DoctorRepo        repositories.DoctorRepository
	HashAlgorithm     utils.Hasher
	Transactor        repositories.Transactor
	AuthTokenProvider utils.AuthTokenProvider
	EmailSender       utils.EmailSender
}

type VerifyUsecase interface {
	EmailVerification(ctx context.Context, req dtos.VerificationReq) error
	EmailVerificationWithTx(ctx context.Context, req dtos.VerificationReq) error
	ResendEmailVerification(ctx context.Context, req dtos.ResendVerificationReq) error
}

type VerifyUsecaseImpl struct {
	UserRepository    repositories.UserRepository
	DoctorRepository  repositories.DoctorRepository
	HashAlgorithm     utils.Hasher
	Transactor        repositories.Transactor
	AuthTokenProvider utils.AuthTokenProvider
	EmailSender       utils.EmailSender
}

func NewVerifyUsecaseImpl(verifyOpts *VerifyUsecaseOpts) *VerifyUsecaseImpl {
	return &VerifyUsecaseImpl{
		UserRepository:    verifyOpts.UserRepo,
		DoctorRepository:  verifyOpts.DoctorRepo,
		HashAlgorithm:     verifyOpts.HashAlgorithm,
		Transactor:        verifyOpts.Transactor,
		AuthTokenProvider: verifyOpts.AuthTokenProvider,
		EmailSender:       verifyOpts.EmailSender,
	}
}

func (u *VerifyUsecaseImpl) MatchUserPassword(ctx context.Context, email string, password string) error {
	user, err := u.UserRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return err
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(user.Password.String))
	if !isCorrectPassword {
		return custom_errors.Unauthorized(err, constants.InvalidPasswordErrMsg)
	}

	return nil
}

func (u *VerifyUsecaseImpl) MatchDoctorPassword(ctx context.Context, email string, password string) error {
	user, err := u.DoctorRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return err
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(user.Password.String))
	if !isCorrectPassword {
		return custom_errors.Unauthorized(err, constants.InvalidPasswordErrMsg)
	}

	return nil
}

func (u *VerifyUsecaseImpl) EmailVerification(ctx context.Context, req dtos.VerificationReq) error {
	jwtMap, err := u.AuthTokenProvider.ParseAndVerify(req.Token)
	if err != nil {
		return custom_errors.InvalidAuthToken()
	}

	data := jwtMap["data"]
	values, _ := data.(map[string]interface{})
	var email string
	var role string

	for key, value := range values {
		if key == constants.UserEmail {
			email = value.(string)
		}
		if key == constants.Role {
			role = value.(string)
		}
	}

	if role == constants.UserRole {
		err = u.MatchUserPassword(ctx, email, req.Password)
		if err != nil {
			return err
		}

		err := u.UserRepository.VerifyUser(ctx, email)
		if err != nil {
			return err
		}

		return nil
	}

	if role == constants.DoctorRole {
		err = u.MatchDoctorPassword(ctx, email, req.Password)
		if err != nil {
			return err
		}

		err := u.DoctorRepository.VerifyDoctor(ctx, email)
		if err != nil {
			return err
		}

		return nil
	}

	return custom_errors.InvalidRole()
}

func (u *VerifyUsecaseImpl) EmailVerificationWithTx(ctx context.Context, req dtos.VerificationReq) error {
	_, err := u.Transactor.WithinTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		err := u.EmailVerification(txCtx, req)
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

func (u *VerifyUsecaseImpl) checkRoleForSendVerify(ctx context.Context, req dtos.ResendVerificationReq) (string, error) {
	if req.Role == constants.UserRole {
		user, err := u.UserRepository.FindOneByEmail(ctx, req.Email)
		if err != nil {
			return "", err
		}
		if user.Email == "" {
			return "", custom_errors.EmailNotFound()
		}
		if user.IsVerified {
			return "", custom_errors.VerifiedEmail()
		}
		return constants.UserRole, nil
	}
	if req.Role == constants.DoctorRole {
		user, err := u.DoctorRepository.FindOneByEmail(ctx, req.Email)
		if err != nil {
			return "", err
		}
		if user.Email == "" {
			return "", custom_errors.EmailNotFound()
		}
		if user.IsVerified {
			return "", custom_errors.VerifiedEmail()
		}
		return constants.DoctorRole, nil
	}

	return "", custom_errors.InvalidRole()
}

func (u *VerifyUsecaseImpl) ResendEmailVerification(ctx context.Context, req dtos.ResendVerificationReq) error {
	config, err := utils.ConfigInit()
	if err != nil {
		return err
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap["userEmail"] = req.Email

	role, err := u.checkRoleForSendVerify(ctx, req)
	if err != nil {
		return err
	}
	dataTokenMap["role"] = role

	token, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("%s%s", config.FrontendUrl, token.AccessToken)

	err = u.EmailSender.SendEmail(req.Email, message, constants.VerificationEmailSubject)
	if err != nil {
		return custom_errors.SendEmail()
	}

	return nil
}

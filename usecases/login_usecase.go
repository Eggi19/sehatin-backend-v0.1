package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type LoginUsecaseOpts struct {
	UserRepo            repositories.UserRepository
	UserAddressRepo     repositories.UserAddressRepository
	DoctorRepo          repositories.DoctorRepository
	PharmacyManagerRepo repositories.PharmacyManagerRepository
	AdminRepo           repositories.AdminRepository
	HashAlgorithm       utils.Hasher
	AuthTokenProvider   utils.AuthTokenProvider
}

type LoginUsecase interface {
	LoginUser(ctx context.Context, email, password string) (*utils.JwtToken, *entities.User, error)
	LoginDoctor(ctx context.Context, email, password string) (*utils.JwtToken, *entities.Doctor, error)
	LoginPharmacyManager(ctx context.Context, email, password string) (*utils.JwtToken, *entities.PharmacyManager, error)
	LoginAdmin(ctx context.Context, email, password string) (*utils.JwtToken, *entities.Admin, error)
}

type LoginUsecaseImpl struct {
	UserRepository            repositories.UserRepository
	UserAddressRepository     repositories.UserAddressRepository
	DoctorRepository          repositories.DoctorRepository
	PharmacyManagerRepository repositories.PharmacyManagerRepository
	AdminRepository           repositories.AdminRepository
	AuthTokenProvider         utils.AuthTokenProvider
	HashAlgorithm             utils.Hasher
}

func NewLoginUsecaseImpl(loginOpts *LoginUsecaseOpts) LoginUsecase {
	return &LoginUsecaseImpl{
		UserRepository:            loginOpts.UserRepo,
		UserAddressRepository:     loginOpts.UserAddressRepo,
		DoctorRepository:          loginOpts.DoctorRepo,
		PharmacyManagerRepository: loginOpts.PharmacyManagerRepo,
		AdminRepository:           loginOpts.AdminRepo,
		AuthTokenProvider:         loginOpts.AuthTokenProvider,
		HashAlgorithm:             loginOpts.HashAlgorithm,
	}
}

func (u *LoginUsecaseImpl) LoginUser(ctx context.Context, email, password string) (*utils.JwtToken, *entities.User, error) {
	user, err := u.UserRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	addresses, err := u.UserAddressRepository.FindAllByUserId(ctx, user.Id)
	if err != nil {
		return nil, nil, err
	}
	user.Address = addresses

	if !user.IsVerified {
		return nil, nil, custom_errors.NotVerified()
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(user.Password.String))
	if !isCorrectPassword {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = user.Id
	dataTokenMap[constants.Role] = constants.UserRole

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := u.AuthTokenProvider.RefreshToken(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}

	return tokens, user, nil
}

func (u *LoginUsecaseImpl) LoginDoctor(ctx context.Context, email, password string) (*utils.JwtToken, *entities.Doctor, error) {
	doctor, err := u.DoctorRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	if !doctor.IsVerified {
		return nil, nil, custom_errors.NotVerified()
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(doctor.Password.String))
	if !isCorrectPassword {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = doctor.Id
	dataTokenMap[constants.Role] = constants.DoctorRole

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := u.AuthTokenProvider.RefreshToken(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}

	return tokens, doctor, nil
}

func (u *LoginUsecaseImpl) LoginPharmacyManager(ctx context.Context, email, password string) (*utils.JwtToken, *entities.PharmacyManager, error) {
	pharmacyManager, err := u.PharmacyManagerRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(pharmacyManager.Password))
	if !isCorrectPassword {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = pharmacyManager.Id
	dataTokenMap[constants.Role] = constants.PharmacyManagerRole

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := u.AuthTokenProvider.RefreshToken(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}

	return tokens, pharmacyManager, nil
}

func (u *LoginUsecaseImpl) LoginAdmin(ctx context.Context, email, password string) (*utils.JwtToken, *entities.Admin, error) {
	admin, err := u.AdminRepository.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	isCorrectPassword, err := u.HashAlgorithm.CheckPassword(password, []byte(admin.Password))
	if !isCorrectPassword {
		return nil, nil, custom_errors.Unauthorized(err, constants.InvalidCredentialsErrMsg)
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = admin.Id
	dataTokenMap[constants.Role] = constants.AdminRole

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := u.AuthTokenProvider.RefreshToken(dataTokenMap)
	if err != nil {
		return nil, nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken:  accessToken.AccessToken,
		RefreshToken: refreshToken.RefreshToken,
	}

	return tokens, admin, nil
}

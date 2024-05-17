package usecases

import (
	"context"
	"database/sql"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/repositories"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type OAuthUsecaseOpts struct {
	UserRepo          repositories.UserRepository
	UserAddressRepo   repositories.UserAddressRepository
	AuthTokenProvider utils.AuthTokenProvider
	GoogleSigner      utils.GoogleSigner
}

type OAuthUsecase interface {
	GoogleOauth(ctx context.Context, oauthData dtos.GoogleAuthRequest) (*utils.JwtToken, *entities.User, error)
}

type OAuthUsecaseImpl struct {
	UserRepository        repositories.UserRepository
	UserAddressRepository repositories.UserAddressRepository
	AuthTokenProvider     utils.AuthTokenProvider
	GoogleSigner          utils.GoogleSigner
}

func NewOAuthUsecaseImpl(oauthOpts *OAuthUsecaseOpts) *OAuthUsecaseImpl {
	return &OAuthUsecaseImpl{
		UserRepository:        oauthOpts.UserRepo,
		UserAddressRepository: oauthOpts.UserAddressRepo,
		AuthTokenProvider:     oauthOpts.AuthTokenProvider,
		GoogleSigner:          oauthOpts.GoogleSigner,
	}
}

func (u *OAuthUsecaseImpl) GoogleOauth(ctx context.Context, oauthData dtos.GoogleAuthRequest) (*utils.JwtToken, *entities.User, error) {
	if oauthData.AuthCode == nil {
		return nil, nil, custom_errors.Unauthorized(custom_errors.ErrNoGoogleAuthCode, constants.NoGoogleAuthCodeErrMsg)
	}

	googleToken, err := u.GoogleSigner.RetrieveToken(*oauthData.AuthCode)
	if err != nil {
		return nil, nil, custom_errors.InternalServerError(err)
	}

	userData, err := u.GoogleSigner.RetrieveUserData(*googleToken)
	if err != nil {
		return nil, nil, custom_errors.InternalServerError(err)
	}

	if !userData.VerifiedEmail {
		return nil, nil, custom_errors.Forbidden()
	}

	isValidPicture := userData.Picture != ""

	birthDate := sql.NullString{String: "", Valid: false}
	if userData.BirthDate != nil {
		birthDate.String = *userData.BirthDate
		birthDate.Valid = true
	}

	gender := &entities.Gender{}
	if userData.Gender != nil {
		if *userData.Gender == constants.Male {
			gender.Id = constants.MaleId
		}
		if *userData.Gender == constants.Female {
			gender.Id = constants.FemaleId
		}
	}

	user, err := u.UserRepository.FindOneByEmail(ctx, userData.Email)

	switch err.(type) {
	case *custom_errors.AppError:
		user, err = u.UserRepository.CreateOneUser(ctx, entities.User{
			Name:           userData.Name,
			Email:          userData.Email,
			ProfilePicture: sql.NullString{String: userData.Picture, Valid: isValidPicture},
			BirthDate:      birthDate,
			Gender:         gender,
			IsGoogle:       true,
			IsVerified:     true,
		})
		if err != nil {
			return nil, nil, custom_errors.InternalServerError(err)
		}

		user.Name = userData.Name
		user.Email = userData.Email
		user.ProfilePicture = sql.NullString{String: userData.Picture, Valid: isValidPicture}
		user.IsGoogle = true
		user.IsVerified = true
		user.IsOnline = false
	case error:
		return nil, nil, custom_errors.InternalServerError(err)
	}

	addresses, err := u.UserAddressRepository.FindAllByUserId(ctx, user.Id)
	if err != nil {
		return nil, nil, err
	}
	user.Address = addresses

	dataTokenMap := make(map[string]interface{})
	dataTokenMap[constants.Id] = user.Id
	dataTokenMap[constants.Role] = oauthData.Role

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

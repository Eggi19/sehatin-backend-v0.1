package usecases

import (
	"context"

	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
)

type RefreshTokenOpts struct {
	AuthTokenProvider utils.AuthTokenProvider
}

type RefreshTokenUsecase interface {
	RefreshToken(ctx context.Context, token string) (*utils.JwtToken, error)
}

type RefreshTokenImpl struct {
	AuthTokenProvider utils.AuthTokenProvider
}

func NewRefreshTokenImpl(refOpts *RefreshTokenOpts) *RefreshTokenImpl {
	return &RefreshTokenImpl{
		AuthTokenProvider: refOpts.AuthTokenProvider,
	}
}

func (u *RefreshTokenImpl) RefreshToken(ctx context.Context, token string) (*utils.JwtToken, error) {
	claims, err := u.AuthTokenProvider.ParseAndVerify(token)
	if err != nil {
		return nil, err
	}

	data := claims["data"]
	values, _ := data.(map[string]interface{})
	var id int64
	var role string

	for key, value := range values {
		if key == "id" {
			id = int64(id)
		}
		if key == "role" {
			role = value.(string)
		}
	}

	dataTokenMap := make(map[string]interface{})
	dataTokenMap["id"] = id
	dataTokenMap["role"] = role

	accessToken, err := u.AuthTokenProvider.CreateAndSign(dataTokenMap)
	if err != nil {
		return nil, err
	}

	tokens := &utils.JwtToken{
		AccessToken: accessToken.AccessToken,
	}

	return tokens, err
}

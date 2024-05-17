package utils

import (
	"errors"
	"strings"
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthTokenProvider interface {
	CreateAndSign(data map[string]interface{}) (*JwtToken, error)
	RefreshToken(data map[string]interface{}) (*JwtToken, error)
	ParseAndVerify(signed string) (jwt.MapClaims, error)
	IsAuthorized(ctx *gin.Context) (bool, *ClaimsData, error)
	GetToken(ctx *gin.Context) (string, error)
	ValidateAdminRoleJwt(ctx *gin.Context) error
	ValidateUserRoleJwt(ctx *gin.Context) error
	ValidatePharmacyManagerRoleJwt(ctx *gin.Context) error
	ValidateDoctorRoleJwt(ctx *gin.Context) error
	ValidateMultiRoleJwt(ctx *gin.Context, roles []string) error
	GenerateResetPasswordToken(data map[string]interface{}) (string, error)
}

type JwtProvider struct {
	config Config
}

type JwtToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ClaimsData struct {
	Id   int64
	Role string
}

func NewJwtProvider(config Config) AuthTokenProvider {
	return &JwtProvider{
		config: config,
	}
}

func (j *JwtProvider) CreateAndSign(data map[string]interface{}) (*JwtToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  j.config.Issuer,
		"exp":  time.Now().Add(time.Duration(j.config.ExpDurationHour) * time.Hour).Unix(),
		"iat":  time.Now(),
		"data": data,
	})

	signed, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		AccessToken:  signed,
		RefreshToken: "",
	}, nil
}

func (j *JwtProvider) RefreshToken(data map[string]interface{}) (*JwtToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  j.config.Issuer,
		"exp":  time.Now().Add(time.Duration(j.config.RefreshExpDuration) * time.Hour).Unix(),
		"iat":  time.Now(),
		"data": data,
	})

	signed, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		AccessToken:  "",
		RefreshToken: signed,
	}, nil
}

func (j *JwtProvider) ParseAndVerify(signed string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.SecretKey), nil
	}, jwt.WithIssuer(j.config.Issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if err.Error() == "token has invalid claims: token is expired" {
			return nil, custom_errors.BadRequest(err, constants.ExpiredTokenErrMsg)
		}
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, custom_errors.InvalidAuthToken()
}

func (j *JwtProvider) IsAuthorized(ctx *gin.Context) (bool, *ClaimsData, error) {
	token, err := j.GetToken(ctx)
	if err != nil {
		return false, nil, err
	}

	claims, err := j.ParseAndVerify(token)
	if err != nil {
		return false, nil, err
	}
	dataMap := claims["data"]
	data, _ := dataMap.(map[string]interface{})

	id := int64(data[constants.Id].(float64))
	role := data[constants.Role].(string)

	if id != 0 && role != "" {
		return true, &ClaimsData{Id: id, Role: role}, nil
	}

	return false, nil, custom_errors.InvalidAuthToken()
}

func (j *JwtProvider) ValidateAdminRoleJwt(ctx *gin.Context) error {
	authorized, data, err := j.IsAuthorized(ctx)
	if err != nil {
		return err
	}

	if authorized && data.Id != 0 && data.Role == constants.AdminRole {
		return nil
	}

	return custom_errors.Forbidden()
}

func (j *JwtProvider) ValidateUserRoleJwt(ctx *gin.Context) error {
	authorized, data, err := j.IsAuthorized(ctx)
	if err != nil {
		return err
	}

	if authorized && data.Id != 0 && data.Role == constants.UserRole {
		return nil
	}

	return custom_errors.Forbidden()
}

func (j *JwtProvider) ValidatePharmacyManagerRoleJwt(ctx *gin.Context) error {
	authorized, data, err := j.IsAuthorized(ctx)
	if err != nil {
		return err
	}

	if authorized && data.Id != 0 && data.Role == constants.PharmacyManagerRole {
		return nil
	}

	return custom_errors.Unauthorized(err, "invalid pharmacy manager token")
}

func (j *JwtProvider) ValidateDoctorRoleJwt(ctx *gin.Context) error {
	authorized, data, err := j.IsAuthorized(ctx)
	if err != nil {
		return err
	}

	if authorized && data.Id != 0 && data.Role == constants.DoctorRole {
		return nil
	}

	return custom_errors.Unauthorized(err, "invalid doctor token")
}

func (j *JwtProvider) ValidateMultiRoleJwt(ctx *gin.Context, roles []string) error {
	authorized, data, err := j.IsAuthorized(ctx)
	if err != nil {
		return err
	}

	for _, role := range roles {
		if authorized && data.Id != 0 && data.Role == role {
			return nil
		}
	}

	return custom_errors.Unauthorized(err, "invalid tokens")
}

func (j *JwtProvider) GetToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	t := strings.Fields(authHeader)
	if len(t) == 2 && t[0] == "Bearer" {
		authToken := t[1]
		return authToken, nil
	}

	return "", errors.New("token not found")
}

func (j *JwtProvider) GenerateResetPasswordToken(data map[string]interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":  j.config.Issuer,
		"exp":  time.Now().Add(time.Duration(j.config.ExpDurationHour) * time.Minute).Unix(),
		"iat":  time.Now(),
		"data": data,
	})

	signed, err := token.SignedString([]byte(j.config.SecretKey))
	if err != nil {
		return "", err
	}

	return signed, nil
}

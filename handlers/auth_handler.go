package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandlerOpts struct {
	LoginUsecase        usecases.LoginUsecase
	RegisterUsecase     usecases.RegisterUsecase
	VerifyUsecase       usecases.VerifyUsecase
	RefreshTokenUsecase usecases.RefreshTokenUsecase
	OAuthUsecase        usecases.OAuthUsecase
	AuthTokenProvider   utils.AuthTokenProvider
}

type AuthHandler struct {
	LoginUsecase        usecases.LoginUsecase
	RegisterUsecase     usecases.RegisterUsecase
	VerifyUsecase       usecases.VerifyUsecase
	RefreshTokenUsecase usecases.RefreshTokenUsecase
	OAuthUsecase        usecases.OAuthUsecase
	AuthTokenProvider   utils.AuthTokenProvider
}

func NewAuthHandler(ahOpts *AuthHandlerOpts) *AuthHandler {
	return &AuthHandler{
		LoginUsecase:        ahOpts.LoginUsecase,
		RegisterUsecase:     ahOpts.RegisterUsecase,
		VerifyUsecase:       ahOpts.VerifyUsecase,
		RefreshTokenUsecase: ahOpts.RefreshTokenUsecase,
		OAuthUsecase:        ahOpts.OAuthUsecase,
		AuthTokenProvider:   ahOpts.AuthTokenProvider,
	}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var payload dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}
	var (
		user            *entities.User
		doctor          *entities.Doctor
		pharmacyManager *entities.PharmacyManager
		admin           *entities.Admin
		token           *utils.JwtToken
		response        *dtos.LoginResponse
		err             error
	)

	switch payload.Role {
	case constants.UserRole:
		token, user, err = h.LoginUsecase.LoginUser(ctx, payload.Email, payload.Password)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
	case constants.DoctorRole:
		token, doctor, err = h.LoginUsecase.LoginDoctor(ctx, payload.Email, payload.Password)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
	case constants.PharmacyManagerRole:
		token, pharmacyManager, err = h.LoginUsecase.LoginPharmacyManager(ctx, payload.Email, payload.Password)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
	case constants.AdminRole:
		token, admin, err = h.LoginUsecase.LoginAdmin(ctx, payload.Email, payload.Password)
		if err != nil {
			_ = ctx.Error(err)
			return
		}
	default:
		err = custom_errors.BadRequest(err, constants.InvalidRoleErrMsg)
		_ = ctx.Error(err)
		return
	}

	availableRole := dtos.AvailableRole{
		User:            user,
		Doctor:          doctor,
		PharmacyManager: pharmacyManager,
		Admin:           admin,
	}

	response = dtos.ConvertToLoginResponse(payload.Role, &availableRole)
	expires := utils.SetExpire()

	response.Exp = expires.AccessTokenExp
	response.Tokens.AccessToken = token.AccessToken
	response.Tokens.RefreshToken = token.RefreshToken

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgLogin,
		Data:    response,
	})
}

func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
	var payload dtos.UserRegisterData
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	u := entities.User{
		Name:      payload.Name,
		Email:     payload.Email,
		Password:  *utils.StringToNullString(*payload.Password),
		BirthDate: *utils.StringToNullString(*payload.BirthDate),
		Gender:    &entities.Gender{Id: *payload.GenderId},
	}

	err := h.RegisterUsecase.RegisterUserWithTransaction(ctx, u)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgRegistered,
		Data:    nil,
	})
}

func (h *AuthHandler) Verification(ctx *gin.Context) {
	var payload dtos.VerificationReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err := h.VerifyUsecase.EmailVerificationWithTx(ctx, payload)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgVerified,
		Data:    nil,
	})
}

func (h *AuthHandler) ResendVerification(ctx *gin.Context) {
	var payload dtos.ResendVerificationReq
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err := h.VerifyUsecase.ResendEmailVerification(ctx, payload)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgSendEmail,
		Data:    nil,
	})
}

func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := h.AuthTokenProvider.GetToken(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	tokens, err := h.RefreshTokenUsecase.RefreshToken(ctx, refreshToken)
	if err != nil {
		ctx.Error(err)
		return
	}

	expires := utils.SetExpire()
	response := dtos.LoginResponse{
		User: nil,
		Exp:  expires.AccessTokenExp,
		Tokens: dtos.TokenResponse{
			AccessToken: tokens.AccessToken,
		},
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgRefreshCookie,
		Data:    response,
	})
}

func (h *AuthHandler) RegisterDoctor(ctx *gin.Context) {
	var payload dtos.DoctorRegisterData

	file, err := ctx.FormFile("certificate")
	if err != nil {
		ctx.Error(custom_errors.FileRequired())
		return
	}
	payload.Certificate = file

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err = h.RegisterUsecase.RegisterDoctorWithTransaction(ctx, payload)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgRegistered,
		Data:    nil,
	})
}

func (h *AuthHandler) RegisterPharmacyManager(ctx *gin.Context) {
	var payload dtos.PharmacyManagerData

	file, err := ctx.FormFile("logo")
	if err != nil {
		ctx.Error(custom_errors.FileRequired())
		return
	}
	payload.Logo = file

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	err = h.RegisterUsecase.RegisterPharmacyManager(ctx, payload)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgRegistered,
		Data:    nil,
	})
}

func (h *AuthHandler) GoogleOauth(ctx *gin.Context) {
	var payload dtos.GoogleAuthRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	token, user, err := h.OAuthUsecase.GoogleOauth(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	doctor := &entities.Doctor{
		Id:             user.Id,
		Name:           user.Name,
		Email:          user.Email,
		IsVerified:     user.IsVerified,
		IsGoogle:       user.IsGoogle,
		IsOnline:       user.IsOnline,
		ProfilePicture: user.ProfilePicture,
	}

	pharmacyManager := &entities.PharmacyManager{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Logo:  user.ProfilePicture.String,
	}

	admin := &entities.Admin{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	availableRole := dtos.AvailableRole{
		User:            user,
		Doctor:          doctor,
		PharmacyManager: pharmacyManager,
		Admin:           admin,
	}

	response := dtos.ConvertToLoginResponse(payload.Role, &availableRole)
	expires := utils.SetExpire()
	response.Exp = expires.AccessTokenExp
	response.Tokens.AccessToken = token.AccessToken
	response.Tokens.RefreshToken = token.RefreshToken

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgLogin,
		Data:    response,
	})
}

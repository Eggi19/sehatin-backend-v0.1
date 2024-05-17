package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type ResetPasswordHandlerOpts struct {
	UserResetPasswordUsecase   usecases.UserResetPasswordUsecase
	DoctorResetPasswordUsecase usecases.DoctorResetPasswordUsecase
	AuthTokenProvider          utils.AuthTokenProvider
}

type ResetPasswordHandler struct {
	UserResetPasswordUsecase   usecases.UserResetPasswordUsecase
	DoctorResetPasswordUsecase usecases.DoctorResetPasswordUsecase
	AuthTokenProvider          utils.AuthTokenProvider
}

func NewUserResetPasswordHandler(rpOpts *ResetPasswordHandler) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		UserResetPasswordUsecase:   rpOpts.UserResetPasswordUsecase,
		DoctorResetPasswordUsecase: rpOpts.DoctorResetPasswordUsecase,
		AuthTokenProvider:          rpOpts.AuthTokenProvider,
	}
}

func (h *ResetPasswordHandler) ForgotPassword(ctx *gin.Context) {
	var payload dtos.ForgotPasswordRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	switch payload.Role {
	case constants.UserRole:
		err := h.UserResetPasswordUsecase.ForgotPasswordWithTransaction(ctx, payload.Email)
		if err != nil {
			ctx.Error(err)
			return
		}
	case constants.DoctorRole:
		err := h.DoctorResetPasswordUsecase.ForgotPasswordWithTransaction(ctx, payload.Email)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgForgotPasswordSuccess,
		Data:    nil,
	})
}

func (h *ResetPasswordHandler) ResetPassword(ctx *gin.Context) {
	var payload dtos.ResetPasswordRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	claims, err := h.AuthTokenProvider.ParseAndVerify(payload.Token)
	if err != nil {
		ctx.Error(err)
		return
	}

	dataMap := claims["data"]
	data, _ := dataMap.(map[string]interface{})

	id := int64(data[constants.Id].(float64))
	role := data[constants.Role].(string)

	if id == 0 && role == "" {
		ctx.Error(custom_errors.Unauthorized(err, constants.ResponseMsgUnauthorized))
		return
	}

	datas := &utils.ClaimsData{Id: id, Role: role}

	switch datas.Role {
	case constants.UserRole:
		err := h.UserResetPasswordUsecase.ResetPasswordWitTransaction(ctx, payload.Email, payload.Token, payload.NewPassword)
		if err != nil {
			ctx.Error(err)
			return
		}
	case constants.DoctorRole:
		err := h.DoctorResetPasswordUsecase.ResetPasswordWitTransaction(ctx, payload.Email, payload.Token, payload.NewPassword)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgResetPasswordSuccess,
		Data:    nil,
	})
}

func (h *ResetPasswordHandler) ChangePassword(ctx *gin.Context) {
	var payload dtos.ChangePasswordRequest

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	switch datas.Role {
	case constants.UserRole:
		err = h.UserResetPasswordUsecase.ChangePasswordWithTransaction(ctx, datas.Id, payload.Password, payload.NewPassword)
		if err != nil {
			ctx.Error(err)
			return
		}
	case constants.DoctorRole:
		err = h.DoctorResetPasswordUsecase.ChangePasswordWithTransaction(ctx, datas.Id, payload.Password, payload.NewPassword)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgChangePasswordSuccess,
		Data:    nil,
	})
}

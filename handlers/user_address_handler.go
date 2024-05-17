package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type UserAddressHandlerOpts struct {
	UserAddressUsecase usecases.UserAddressUsecase
}

type UserAddressHandler struct {
	UserAddressUsecase usecases.UserAddressUsecase
}

func NewUserAddressHandler(uaOpts *UserAddressHandlerOpts) *UserAddressHandler {
	return &UserAddressHandler{UserAddressUsecase: uaOpts.UserAddressUsecase}
}

func (h *UserAddressHandler) CreateUserAddress(ctx *gin.Context) {

	var payload dtos.UserAddressCreateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	userAddress := entities.UserAddress{
		Id:          0,
		UserId:      int64(userId),
		City:        payload.City,
		CityId:      payload.CityId,
		Province:    payload.Province,
		Address:     payload.Address,
		District:    payload.District,
		SubDistrict: payload.SubDistrict,
		PostalCode:  payload.PostalCode,
		Longitude:   payload.Longitude,
		Latitude:    payload.Latitude,
		IsMain:      *payload.IsMain,
		Coordinate:  "",
	}

	err = h.UserAddressUsecase.CreateUserAddress(ctx, userAddress)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *UserAddressHandler) GetAddressById(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	addressId, err := utils.GetIdParamOrContext(ctx, "addressId")
	if err != nil {
		ctx.Error(err)
		return
	}

	address, err := h.UserAddressUsecase.GetAddressById(ctx, int64(addressId), int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToUserAddressResponse(address),
	})
}

func (h *UserAddressHandler) UpdateUserAddress(ctx *gin.Context) {
	var payload dtos.UserAddressUpdateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	addressId, err := utils.GetIdParamOrContext(ctx, "addressId")
	if err != nil {
		ctx.Error(err)
		return
	}

	userAddress := entities.UserAddress{
		Id:          int64(addressId),
		UserId:      int64(userId),
		City:        payload.City,
		CityId:      payload.CityId,
		Province:    payload.Province,
		Address:     payload.Address,
		District:    payload.District,
		SubDistrict: payload.SubDistrict,
		PostalCode:  payload.PostalCode,
		Longitude:   payload.Longitude,
		Latitude:    payload.Latitude,
		IsMain:      *payload.IsMain,
	}

	err = h.UserAddressUsecase.UpdateUserAddress(ctx, userAddress)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *UserAddressHandler) DeleteUserAddress(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	addressId, err := utils.GetIdParamOrContext(ctx, "addressId")
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.UserAddressUsecase.DeleteUserAddress(ctx, int64(addressId), int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

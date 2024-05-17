package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type ShippingMethodHandlerOpts struct {
	ShippingMethodUsecase usecases.ShippingMethodUsecase
}

type ShippingMethodHandler struct {
	ShippingMethodUsecase usecases.ShippingMethodUsecase
}

func NewShippingMethodHandler(smHandOpts *ShippingMethodHandlerOpts) *ShippingMethodHandler {
	return &ShippingMethodHandler{
		ShippingMethodUsecase: smHandOpts.ShippingMethodUsecase,
	}
}

func (h *ShippingMethodHandler) GetOfficialShippingCost(ctx *gin.Context) {
	var payload dtos.OfficialShippingFeeRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	cost, err := h.ShippingMethodUsecase.GetOfficialFee(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data: dtos.ShippingCostResponse{
			Cost: cost,
		},
	})
}

func (h *ShippingMethodHandler) GetNonOfficialShippingCost(ctx *gin.Context) {
	var payload dtos.NonOfficialShippingFeeRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	cost, err := h.ShippingMethodUsecase.GetNonOfficialFee(ctx, payload, userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ShippingCostResponse{
			Cost: cost,
		},
	})
}

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

type CartHandlerOpts struct {
	CartUsecase usecases.CartUsecase
}

type CartHandler struct {
	CartUsecase usecases.CartUsecase
}

func NewCartHandler(chOpts *CartHandlerOpts) *CartHandler {
	return &CartHandler{
		CartUsecase: chOpts.CartUsecase,
	}
}

func (h *CartHandler) CreateCartItem(ctx *gin.Context) {
	var payload dtos.CreateCartItemRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	cartEntity := entities.CartItem{
		UserId:            int64(userId),
		Quantity:          payload.Quantity,
		PharmacyProductId: payload.PharmacyProductId,
	}

	err = h.CartUsecase.CreateCartItem(ctx, cartEntity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *CartHandler) IncreaseCartItem(ctx *gin.Context) {
	var payload dtos.UpdateCartItemRequest
	payload.Quantity = 1

	cartId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	cartEntity := entities.CartItem{
		Id:       int64(cartId),
		Quantity: payload.Quantity,
	}

	err = h.CartUsecase.IncreaseCartItemQuantity(ctx, cartEntity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgIncreaseQuantity,
		Data:    nil,
	})
}

func (h *CartHandler) DecreaseCartItem(ctx *gin.Context) {
	var payload dtos.UpdateCartItemRequest
	payload.Quantity = 1

	cartId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	cartEntity := entities.CartItem{
		Id:       int64(cartId),
		Quantity: payload.Quantity,
	}

	err = h.CartUsecase.DecreaseCartItemQuantity(ctx, cartEntity)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDecreaseQuantity,
		Data:    nil,
	})
}

func (h *CartHandler) DeleteCartItem(ctx *gin.Context) {
	cartId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.CartUsecase.DeleteCartItem(ctx, int64(cartId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *CartHandler) GetAllCartItem(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	data, err := h.CartUsecase.GetUserCartItems(ctx, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    data,
	})
}

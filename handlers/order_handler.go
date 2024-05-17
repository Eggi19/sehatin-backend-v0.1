package handlers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandlerOpts struct {
	OrderUsecase usecases.OrderUsecase
}

type OrderHandler struct {
	OrderUsecase usecases.OrderUsecase
}

func NewOrderHandler(spHandOpts *OrderHandlerOpts) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: spHandOpts.OrderUsecase,
	}
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var payload dtos.OrderRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	order, err := h.OrderUsecase.CreateOrderWithTransaction(ctx, payload)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreateOrder,
		Data: dtos.CreateOrderResponse{
			PaymentDeadline: order.PaymentDeadline,
		},
	})
}

func (h *OrderHandler) GetAllOrderByUser(ctx *gin.Context) {
	var params entities.OrderParams

	params.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if ctx.Query("limit") == "" {
		params.Limit = constants.DefaultLimit
	}

	params.Page, _ = strconv.Atoi(ctx.Query("page"))
	if ctx.Query("page") == "" {
		params.Page = constants.DefaultPage
	}

	params.Status = ctx.Query("status")
	if ctx.Query("status") == "" {
		params.Status = constants.Pending
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, total, err := h.OrderUsecase.GetAllOrderByUser(ctx, int64(userId), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgGetOrder,
		Data:    dtos.GetAllOrdersResponse{
			PaginationInfo: dtos.PaginationResponse{
				Page:      params.Page,
				Limit:     params.Limit,
				TotalPage: int(math.Ceil(float64(total) / float64(params.Limit))),
				TotalData: total,
			},
			Data: result,
		},
	})
}

func (h *OrderHandler) GetAllOrderByPharmacyManager(ctx *gin.Context) {
	var params entities.OrderParams

	params.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if ctx.Query("limit") == "" {
		params.Limit = constants.DefaultLimit
	}

	params.Page, _ = strconv.Atoi(ctx.Query("page"))
	if ctx.Query("page") == "" {
		params.Page = constants.DefaultPage
	}

	params.Status = ctx.Query("status")
	if ctx.Query("status") == "" {
		params.Status = constants.Pending
	}

	params.PharmacyId, _ = strconv.Atoi(ctx.Query("pharmacyId"))
	if ctx.Query("pharmacyId") == "" {
		params.PharmacyId = 0
	}

	pharmacyManagerId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, total, err := h.OrderUsecase.GetAllOrderByPharmacyManager(ctx, int64(pharmacyManagerId), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgGetOrder,
		Data:    dtos.GetAllOrdersResponse{
			PaginationInfo: dtos.PaginationResponse{
				Page:      params.Page,
				Limit:     params.Limit,
				TotalPage: int(math.Ceil(float64(total) / float64(params.Limit))),
				TotalData: total,
			},
			Data: result,
		},
	})
}

func (h *OrderHandler) GetAllOrderByAdmin(ctx *gin.Context) {
	var params entities.OrderParams

	params.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if ctx.Query("limit") == "" {
		params.Limit = constants.DefaultLimit
	}

	params.Page, _ = strconv.Atoi(ctx.Query("page"))
	if ctx.Query("page") == "" {
		params.Page = constants.DefaultPage
	}

	params.Status = ctx.Query("status")
	if ctx.Query("status") == "" {
		params.Status = constants.Pending
	}

	params.PharmacyId, _ = strconv.Atoi(ctx.Query("pharmacyId"))
	if ctx.Query("pharmacyId") == "" {
		params.PharmacyId = 0
	}

	result, total, err := h.OrderUsecase.GetAllOrderByAdmin(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgGetOrder,
		Data: dtos.GetAllOrdersResponse{
			PaginationInfo: dtos.PaginationResponse{
				Page:      params.Page,
				Limit:     params.Limit,
				TotalPage: int(math.Ceil(float64(total) / float64(params.Limit))),
				TotalData: total,
			},
			Data: result,
		},
	})
}

func (h *OrderHandler) UpdateOrderStatusToProcessing(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.UpdateOrderStatusToProcessing(ctx, int64(orderId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) UpdateOrderStatusToShipped(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.OrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyManagerId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.UpdateOrderStatusToShipped(ctx, int64(orderId), int64(pharmacyManagerId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) UpdateOrderStatusToCompleted(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.OrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.UpdateOrderStatusToCompleted(ctx, int64(orderId), int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) UploadPaymentProof(ctx *gin.Context) {
	payload := dtos.UploadPaymentProofResponse{}

	file, err := ctx.FormFile("payment_proof")
	if err != nil {
		ctx.Error(custom_errors.FileRequired())
		return
	}
	payload.PaymentProof = file

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.UploadPaymentProof(ctx, payload, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) UpdateOrderStatusToCanceled(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.OrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.UpdateOrderStatusToCanceled(ctx, int64(orderId), int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) CancelOrderByAdmin(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.CancelOrderByAdmin(ctx, int64(orderId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) CancelOrderByPharmacyManager(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.OrderId)
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacyManagerId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.OrderUsecase.CancelOrderByPharmacyManager(ctx, int64(orderId), int64(pharmacyManagerId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *OrderHandler) GetOrderDetail(ctx *gin.Context) {
	orderId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := h.OrderUsecase.GetOrderDetail(ctx, int64(orderId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgGetOrder,
		Data:    result,
	})
}

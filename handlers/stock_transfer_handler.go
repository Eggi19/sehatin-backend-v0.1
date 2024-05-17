package handlers

import (
	"net/http"
	"strconv"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/entities"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/gin-gonic/gin"
)

type StockTransferHandlerOpts struct {
	StockTransferUsecase usecases.StockTransferUsecase
}

type StockTransferHandler struct {
	StockTransferUsecase usecases.StockTransferUsecase
}

func NewStockTransferHandler(sthOpts *StockTransferHandler) *StockTransferHandler {
	return &StockTransferHandler{StockTransferUsecase: sthOpts.StockTransferUsecase}
}

func (h *StockTransferHandler) CreateStockTransfer(ctx *gin.Context) {
	var payload dtos.StockTransferCreateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	stockTransfer := entities.StockTransfer{
		PharmacySender: entities.Pharmacy{
			Id:                        payload.PharmacySenderId,
			PharmacyManager:           entities.PharmacyManager{},
			PharmacyAddress:           entities.PharmacyAddress{},
			OfficialShippingMethod:    []entities.OfficialShippingMethod{},
			NonOfficialShippingMethod: []entities.NonOfficialShippingMethod{},
		},
		PharmacyReceiver: entities.Pharmacy{
			Id:                        payload.PharmacyReceiverId,
			PharmacyManager:           entities.PharmacyManager{},
			PharmacyAddress:           entities.PharmacyAddress{},
			OfficialShippingMethod:    []entities.OfficialShippingMethod{},
			NonOfficialShippingMethod: []entities.NonOfficialShippingMethod{},
		},
		MutationStatus: entities.MutationSatus{Id: constants.MutationPendingId},
		Product: entities.Product{
			Id:                    payload.ProductId,
			ProductForm:           entities.ProductForm{},
			ProductClassification: entities.ProductClassification{},
			Manufacture:           entities.Manufacture{},
			Categories:            []entities.Category{},
		},
		Quantity: payload.Quantity,
	}

	err := h.StockTransferUsecase.CreateStockRequest(ctx, stockTransfer)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *StockTransferHandler) GetAllStockTransfer(ctx *gin.Context) {
	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")

	limit := constants.DefaultLimit
	page := constants.DefaultPage

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
		if limit == 0 {
			ctx.Error(custom_errors.BadRequest(err, constants.ZeroLimitInputErrMsg))
		}
	}

	pageStr := ctx.Query("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
		if page == 0 {
			ctx.Error(custom_errors.BadRequest(err, constants.ZeroPageInputErrMsg))
		}
	}

	params := entities.StockTransferParams{
		SortBy: sortBy,
		Sort:   sort,
		Limit:  limit,
		Page:   page,
	}

	stockTransfers, pagination, err := h.StockTransferUsecase.GetAllStockTransfer(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToStockTransferResponses(stockTransfers, *pagination),
	})
}

func (h *StockTransferHandler) UpdateMutationStatus(ctx *gin.Context) {
	var payload dtos.MutationStatusIdRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	switch payload.MutationStatusId {
	case constants.MutationProcessedId:
		err := h.StockTransferUsecase.UpdateStatusProcessedWithTransaction(ctx, payload.StockTransferId, constants.MutationProcessedId)
		if err != nil {
			ctx.Error(err)
			return
		}
	case constants.MutationCanceledId:
		err := h.StockTransferUsecase.UpdateStatusCanceledWithTransaction(ctx, payload.StockTransferId, constants.MutationCanceledId)
		if err != nil {
			ctx.Error(err)
			return
		}
	default:
		err := h.StockTransferUsecase.UpdateStatusPendingWithTransaction(ctx, payload.StockTransferId, constants.MutationPendingId)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

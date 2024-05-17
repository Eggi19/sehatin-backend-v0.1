package handlers

import (
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

type StockHistoryHandlerOpts struct {
	StockHistoryUsecase usecases.StockHistoryUsecase
}

type StockHistoryHandler struct {
	StockHistoryUsecase usecases.StockHistoryUsecase
}

func NewStockHistoryHandler(sthOpts *StockHistoryHandlerOpts) *StockHistoryHandler {
	return &StockHistoryHandler{
		StockHistoryUsecase: sthOpts.StockHistoryUsecase,
	}
}

func (h *StockHistoryHandler) GetStockHistoriesByPharmacyId(ctx *gin.Context) {
	pharmacyId, err := utils.GetIdParamOrContext(ctx, "pharmacyId")
	if err != nil {
		ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	limit := constants.DefaultLimit
	page := constants.DefaultPage

	limitStr := ctx.Query("limit")
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
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
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
		if page == 0 {
			ctx.Error(custom_errors.BadRequest(err, constants.ZeroPageInputErrMsg))
		}
	}

	params := entities.StockHistoryParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	stockHistories, pagination, err := h.StockHistoryUsecase.GetStockHistoriesByPharmacyId(ctx, int64(pharmacyId), datas.Id, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToStockHistoryResponses(stockHistories, *pagination),
	})
}

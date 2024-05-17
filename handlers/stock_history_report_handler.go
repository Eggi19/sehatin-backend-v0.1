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

type StockHistoryReportHandlerOpts struct {
	StockHistoryReportUsecase usecases.StockHistoryReportUsecase
}

type StockHistoryReportHandler struct {
	StockHistoryReportUsecase usecases.StockHistoryReportUsecase
}

func NewStockHistoryReportHandler(shrOpts *StockHistoryReportHandlerOpts) *StockHistoryReportHandler {
	return &StockHistoryReportHandler{StockHistoryReportUsecase: shrOpts.StockHistoryReportUsecase}
}

func (h *StockHistoryReportHandler) GetStockHistoryReports(ctx *gin.Context) {
	var err error

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	limit := constants.DefaultLimit
	page := constants.DefaultPage
	pharmacyId := constants.DefaultId

	pharmacyIdStr := ctx.Query("pharmacyId")
	if pharmacyIdStr != "" {
		pharmacyId, err = strconv.Atoi(pharmacyIdStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
	}

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

	params := entities.StockHistoryReportParams{
		SortBy:     sortBy,
		Sort:       sort,
		Limit:      limit,
		Page:       page,
		Keyword:    keyword,
		PharmacyId: int64(pharmacyId),
	}

	stockReports, pagination, err := h.StockHistoryReportUsecase.GetStockHistories(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToStockHistoryReportResponses(stockReports, *pagination),
	})
}

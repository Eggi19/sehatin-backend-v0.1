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

type SalesReportCategoryHandlerOpts struct {
	SalesReportCategoryUsecase usecases.SalesReportCategoryUsecase
}

type SalesReportCategoryHandler struct {
	SalesReportCategoryUsecase usecases.SalesReportCategoryUsecase
}

func NewSalesReportCategoryHandler(srhOpts *SalesReportCategoryHandlerOpts) *SalesReportCategoryHandler {
	return &SalesReportCategoryHandler{SalesReportCategoryUsecase: srhOpts.SalesReportCategoryUsecase}
}

func (h *SalesReportCategoryHandler) GetSalesReportCategories(ctx *gin.Context) {
	var err error

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	limit := constants.DefaultLimit
	page := constants.DefaultPage
	categoryId := constants.DefaultId

	categoryIdStr := ctx.Query("pharmacyId")
	if categoryIdStr != "" {
		categoryId, err = strconv.Atoi(categoryIdStr)
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

	params := entities.SalesReportCategoryParams{
		SortBy:     sortBy,
		Sort:       sort,
		Limit:      limit,
		Page:       page,
		Keyword:    keyword,
		CategoryId: int64(categoryId),
	}

	salesReportCategories, paginaiton, err := h.SalesReportCategoryUsecase.GetStockHistories(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertSalesReportCategoryResponses(salesReportCategories, *paginaiton),
	})
}

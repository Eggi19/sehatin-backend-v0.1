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
	"github.com/gin-gonic/gin"
)

type MostBoughtUserHandlerOpts struct {
	MostBoughtUserUsecase usecases.MostBoughtUserUsecase
	SalesReportUsecase    usecases.SalesReportUsecase
}

type MostBoughtUserHandler struct {
	MostBoughtUserUsecase usecases.MostBoughtUserUsecase
	SalesReportUsecase    usecases.SalesReportUsecase
}

func NewMostBoughtUserHandler(mbhOpts *MostBoughtUserHandlerOpts) *MostBoughtUserHandler {
	return &MostBoughtUserHandler{
		MostBoughtUserUsecase: mbhOpts.MostBoughtUserUsecase,
		SalesReportUsecase:    mbhOpts.SalesReportUsecase,
	}
}

func (h *MostBoughtUserHandler) GetMostBought(ctx *gin.Context) {
	var req entities.NearestPharmacyParams
	var query entities.PaginationParams

	query.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if ctx.Query("limit") == "" {
		query.Limit = constants.DefaultLimit
	}

	query.Page, _ = strconv.Atoi(ctx.Query("page"))
	if ctx.Query("page") == "" {
		query.Page = constants.DefaultPage
	}

	req.Longitude, _ = strconv.ParseFloat(ctx.Query("longitude"), 64)
	req.Latitude, _ = strconv.ParseFloat(ctx.Query("latitude"), 64)
	if ctx.Query("latitude") == "" || ctx.Query("longitude") == "" {
		ctx.Error(custom_errors.BadRequest(nil, constants.CoordinateRequiredErrMsg))
		return
	}
	req.Radius, _ = strconv.Atoi(ctx.Query("radius"))
	if ctx.Query("radius") == "" {
		req.Radius = constants.DefaultRadius
	}

	products, err := h.MostBoughtUserUsecase.GetMostBought(ctx, req, query)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var total int = 0
	if len(products) != 0 {
		total = products[0].Total
	}

	data := dtos.GetProductResponse{
		PaginationInfo: *dtos.ConvertToPaginationResponse(entities.PaginationInfo{
			Page:      query.Page,
			Limit:     query.Limit,
			TotalPage: int(math.Ceil(float64(total) / float64(query.Limit))),
			TotalData: total,
		}),
		Products: products,
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    data,
	})
}

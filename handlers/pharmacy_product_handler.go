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

type PharmacyProductHandlerOpts struct {
	PharmacyProductUsecase usecases.PharmacyProductUsecase
}

type PharmacyProductHandler struct {
	PharmacyProductUsecase usecases.PharmacyProductUsecase
}

func NewPharmacyProductHandler(phOpts *PharmacyProductHandlerOpts) *PharmacyProductHandler {
	return &PharmacyProductHandler{
		PharmacyProductUsecase: phOpts.PharmacyProductUsecase,
	}
}

func (h *PharmacyProductHandler) GetNearestProducts(ctx *gin.Context) {
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

	req.CategoryId, _ = strconv.Atoi(ctx.Query("categoryId"))

	products, err := h.PharmacyProductUsecase.GetNearestPharmacyProductsWithTransaction(ctx, req, query)
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

func (h *PharmacyProductHandler) ProductDetail(ctx *gin.Context) {
	var request entities.PharmacyProductDetailParams
	productIdInt, _ := strconv.Atoi(ctx.Query("productId"))
	pharmacyProductIdInt, _ := strconv.Atoi(ctx.Query("pharmacyProductId"))

	request.Coordinat.Latitude, _ = strconv.ParseFloat(ctx.Query("latitude"), 64)
	request.Coordinat.Longitude, _ = strconv.ParseFloat(ctx.Query("longitude"), 64)
	request.Coordinat.Radius, _ = strconv.Atoi(ctx.Query("radius"))
	request.Coordinat.ProductId = int64(productIdInt)
	request.PharmacyProductId = int64(pharmacyProductIdInt)

	productDetail, err := h.PharmacyProductUsecase.GetProductDetail(ctx, request)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    productDetail,
	})
}

func (h *PharmacyProductHandler) CreatePharmacyProduct(ctx *gin.Context) {
	payload := dtos.PharmacyProductRequest{}

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	product := entities.Product{Id: payload.ProductId}
	pharmacy := entities.Pharmacy{Id: payload.PharmacyId}

	pp := entities.PharmacyProduct{
		Price:       payload.Price,
		TotalStock:  *payload.TotalStock,
		IsAvailable: *payload.IsAvailable,
		Product:     product,
		Pharmacy:    pharmacy,
	}

	err = h.PharmacyProductUsecase.CreatePharmacyProduct(ctx, pp, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *PharmacyProductHandler) UpdatePharmacyProduct(ctx *gin.Context) {
	payload := dtos.PharmacyProductRequest{}

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		return
	}

	pharmacyProductIdInt, _ := strconv.Atoi(ctx.Param("id"))

	product := entities.Product{
		Id: payload.ProductId,
	}
	pharmacy := entities.Pharmacy{Id: payload.PharmacyId}

	pp := entities.PharmacyProduct{
		Id:          int64(pharmacyProductIdInt),
		Price:       payload.Price,
		TotalStock:  *payload.TotalStock,
		IsAvailable: *payload.IsAvailable,
		Product:     product,
		Pharmacy:    pharmacy,
	}

	if *payload.TotalStock == 0 {
		pp.IsAvailable = false
	}

	err = h.PharmacyProductUsecase.UpdatePharmacyProduct(ctx, pp, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *PharmacyProductHandler) DeletePharmacyProduct(ctx *gin.Context) {
	pharmacyProductIdInt, _ := strconv.Atoi(ctx.Param("id"))

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		return
	}

	err = h.PharmacyProductUsecase.DeletePharmacyProduct(ctx, int64(pharmacyProductIdInt), datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *PharmacyProductHandler) GetPharmacyProductsByPharmacyId(ctx *gin.Context) {
	pharmacyId, err := utils.GetIdParamOrContext(ctx, "id")
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

	params := entities.PharmacyProductParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	pharmacyProducts, pagination, err := h.PharmacyProductUsecase.GetPharmacyProductsByPharmacyId(ctx, int64(pharmacyId), datas.Id, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyProductResponses(pharmacyProducts, *pagination),
	})
}

func (h *PharmacyProductHandler) GetAllNearestPharmacyProducts(ctx *gin.Context) {
	var params entities.NearestPharmacyProductsParams

	params.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	if ctx.Query("limit") == "" {
		params.Limit = constants.DefaultLimit
	}

	params.Page, _ = strconv.Atoi(ctx.Query("page"))
	if ctx.Query("page") == "" {
		params.Page = constants.DefaultPage
	}

	params.Longitude, _ = strconv.ParseFloat(ctx.Query("longitude"), 64)
	params.Latitude, _ = strconv.ParseFloat(ctx.Query("latitude"), 64)
	if ctx.Query("latitude") == "" || ctx.Query("longitude") == "" {
		ctx.Error(custom_errors.BadRequest(nil, constants.CoordinateRequiredErrMsg))
		return
	}
	params.Radius, _ = strconv.Atoi(ctx.Query("radius"))
	if ctx.Query("radius") == "" {
		params.Radius = constants.DefaultRadius
	}

	params.CategoryId, _ = strconv.Atoi(ctx.Query("categoryId"))

	params.SortBy = ctx.Query("sortBy")
	params.Sort = ctx.Query("sort")
	params.Keyword = ctx.Query("keyword")

	pharmacyProducts, paginationInfo, err := h.PharmacyProductUsecase.GetAllNearestPharmacyProducts(ctx, params)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data: dtos.GetProductResponse{
			PaginationInfo: dtos.PaginationResponse(*paginationInfo),
			Products:       pharmacyProducts,
		},
	})
}

func (h *PharmacyProductHandler) GetPharmacyProductById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	pp, err := h.PharmacyProductUsecase.GetPharmacyProduct(ctx, int64(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyProductItem(*pp),
	})
}

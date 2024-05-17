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

type ProductHandlerOpts struct {
	ProductUsecase         usecases.ProductUsecase
	ProductCategoryUsecase usecases.ProductCategoryUsecase
	UploadFile             utils.FileUploader
}

type ProductHandler struct {
	ProductUsecase         usecases.ProductUsecase
	ProductCategoryUsecase usecases.ProductCategoryUsecase
	UploadFile             utils.FileUploader
}

func NewProductHandler(phUpts *ProductHandlerOpts) *ProductHandler {
	return &ProductHandler{
		ProductUsecase:         phUpts.ProductUsecase,
		ProductCategoryUsecase: phUpts.ProductCategoryUsecase,
		UploadFile:             phUpts.UploadFile,
	}
}

func (h *ProductHandler) CreateProduct(ctx *gin.Context) {
	payload := dtos.ProductCreateRequest{
		Categories: []int64{},
	}

	getFile, err := ctx.FormFile("product_picture")
	if err != nil {
		ctx.Error(custom_errors.FileRequired())
		return
	}
	payload.ProductPicture = getFile

	if err = ctx.ShouldBind(&payload); err != nil {
		ctx.Error(err)
		return
	}

	var fileUrl *string
	var product entities.Product

	fileUrl, err = utils.GetFileUrl(ctx, payload.ProductPicture, "png")
	if err != nil {
		ctx.Error(err)
		return
	}

	product.Name = payload.Name
	product.GenericName = payload.GenericName
	product.Content = payload.Content
	product.Description = payload.Description
	product.UnitInPack = payload.UnitInPack
	product.SellingUnit = payload.SellingUnit
	product.Weight = payload.Weight
	product.Height = payload.Height
	product.Length = payload.Length
	product.Width = payload.Width
	product.SlugId = payload.SlugId
	product.ProductForm = entities.ProductForm{Id: payload.ProductFormId}
	product.ProductClassification = entities.ProductClassification{Id: payload.ProductClassificationId}
	product.Manufacture = entities.Manufacture{Id: payload.ManufactureId}
	product.ProductPicture = *fileUrl

	var productId int64
	productId = 0

	productId, err = h.ProductUsecase.CreateProduct(ctx, product)
	if err != nil {
		ctx.Error(err)
		return
	}

	for i := 0; i < len(payload.Categories); i++ {
		productCategory := entities.ProductCategory{
			Id:         0,
			ProductId:  productId,
			CategoryId: 0,
		}
		productCategory.CategoryId = payload.Categories[i]
		productCategory.ProductId = productId
		err := h.ProductCategoryUsecase.CreateProductCategory(ctx, productCategory)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *ProductHandler) UpdateProduct(ctx *gin.Context) {
	payload := dtos.ProductCreateRequest{
		Categories: []int64{},
	}

	getFile, _ := ctx.FormFile("product_picture")
	payload.ProductPicture = getFile

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.Error(err)
		return
	}

	productId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var product entities.Product

	if payload.ProductPicture != nil {
		fileUrl, err := utils.GetFileUrl(ctx, payload.ProductPicture, "png")
		if err != nil {
			ctx.Error(err)
			return
		}
		product.ProductPicture = *fileUrl
	} else {
		product.ProductPicture = ""
	}

	product.Id = int64(productId)
	product.Name = payload.Name
	product.GenericName = payload.GenericName
	product.Content = payload.Content
	product.Description = payload.Description
	product.UnitInPack = payload.UnitInPack
	product.SellingUnit = payload.SellingUnit
	product.Weight = payload.Weight
	product.Height = payload.Height
	product.Length = payload.Length
	product.Width = payload.Width
	product.SlugId = payload.SlugId
	product.ProductForm = entities.ProductForm{Id: payload.ProductFormId}
	product.ProductClassification = entities.ProductClassification{Id: payload.ProductClassificationId}
	product.Manufacture = entities.Manufacture{Id: payload.ManufactureId}

	err = h.ProductUsecase.UpdateProduct(ctx, product)
	if err != nil {
		ctx.Error(err)
		return
	}

	producCategory := entities.ProductCategory{
		ProductId:  0,
		CategoryId: 0,
	}

	for i := 0; i < len(payload.Categories); i++ {
		producCategory.ProductId = int64(productId)
		producCategory.CategoryId = payload.Categories[i]
		err := h.ProductCategoryUsecase.UpdateProductCategory(ctx, producCategory)
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

func (h *ProductHandler) DeleteProduct(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("id"))

	err := h.ProductUsecase.DeleteProduct(ctx, int64(productId))
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.ProductCategoryUsecase.DeleteProductCategoryByProductId(ctx, int64(productId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *ProductHandler) GetProductById(ctx *gin.Context) {
	productId, _ := strconv.Atoi(ctx.Param("id"))

	product, err := h.ProductUsecase.GetOneProduct(ctx, int64(productId))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToProductResponse(*product),
	})
}

func (h *ProductHandler) GetAllProduct(ctx *gin.Context) {
	var err error

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

	params := entities.ProductCategoryParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	products, pagination, err := h.ProductUsecase.GetAllProduct(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToProductResponses(products, *pagination),
	})
}

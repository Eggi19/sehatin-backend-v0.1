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

type CategoryHandlerOpts struct {
	CategoryUsecase usecases.CategoryUsecase
}

type CategoryHandler struct {
	CategoryUsecase usecases.CategoryUsecase
}

func NewCategoryHandler(chOpts *CategoryHandlerOpts) *CategoryHandler {
	return &CategoryHandler{
		CategoryUsecase: chOpts.CategoryUsecase,
	}
}

func (h *CategoryHandler) GetAllCategory(ctx *gin.Context) {
	var err error
	params := entities.CategoryParams{}

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	limit := constants.DefaultLimit
	page := constants.DefaultPage

	limitStr, limitExist := ctx.GetQuery("limit")
	if limitExist && limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
		if limit == 0 {
			ctx.Error(custom_errors.BadRequest(err, constants.ZeroLimitInputErrMsg))
		}
	}
	if !limitExist || limitStr == "" {
		limit = 0
	}

	pageStr, pageExist := ctx.GetQuery("page")
	if pageExist && pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
		if page == 0 {
			ctx.Error(custom_errors.BadRequest(err, constants.ZeroPageInputErrMsg))
		}
	}
	if !pageExist || pageStr == "" {
		page = 0
	}

	params.Page = page
	params.Limit = limit
	params.SortBy = sortBy
	params.Sort = sort
	params.Keyword = keyword

	categories, pagination, err := h.CategoryUsecase.GetAllCategory(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToCategoryResponses(categories, *pagination),
	})
}

func (h *CategoryHandler) GetCategoryById(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Param("id"))

	category, err := h.CategoryUsecase.GetCategoryById(ctx, int64(categoryId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToCategoryResponse(*category),
	})
}

func (h *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var payload dtos.CategoryRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	category := entities.Category{
		Name: payload.Name,
	}

	err := h.CategoryUsecase.CreateCategory(ctx, category)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var payload dtos.CategoryRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	categoryId, _ := strconv.Atoi(ctx.Param("id"))

	category := entities.Category{
		Id:   int64(categoryId),
		Name: payload.Name,
	}

	err := h.CategoryUsecase.UpdateCategory(ctx, category)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	categoryId, _ := strconv.Atoi(ctx.Param("id"))

	err := h.CategoryUsecase.DeleteCategory(ctx, int64(categoryId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

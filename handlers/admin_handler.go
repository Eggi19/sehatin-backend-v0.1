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

type AdminHandlerOpts struct {
	AdminUsecase usecases.AdminUsecase
}

type AdminHandler struct {
	AdminUsecase usecases.AdminUsecase
}

func NewAdminHandler(ahOpts *AdminHandlerOpts) *AdminHandler {
	return &AdminHandler{AdminUsecase: ahOpts.AdminUsecase}
}

func (h *AdminHandler) CreateAdmin(ctx *gin.Context) {
	var payload dtos.AdminRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	admin := entities.Admin{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}

	err := h.AdminUsecase.CreateAdminWithTransaction(ctx, admin)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *AdminHandler) DeleteAdmin(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := h.AdminUsecase.DeleteAdmin(ctx, int64(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *AdminHandler) GetAdminById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	admin, err := h.AdminUsecase.GetAdminById(ctx, int64(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToAdminResponse(admin),
	})
}

func (h *AdminHandler) GetAllAdmin(ctx *gin.Context) {
	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

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

	params := entities.AdminParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	admins, pagination, err := h.AdminUsecase.GetAllAdmin(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToAdminResponses(admins, *pagination),
	})
}

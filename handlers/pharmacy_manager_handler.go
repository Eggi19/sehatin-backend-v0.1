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

type PharmacyManagerHandlerOpts struct {
	PharmacyManagerUsecase usecases.PharmacyManagerUsecase
	UploadFile             utils.FileUploader
}

type PharmacyManagerHandler struct {
	PharmacyManagerUsecase usecases.PharmacyManagerUsecase
	UploadFile             utils.FileUploader
}

func NewPharmacyManagerHandler(pmhOpts *PharmacyManagerHandlerOpts) *PharmacyManagerHandler {
	return &PharmacyManagerHandler{
		PharmacyManagerUsecase: pmhOpts.PharmacyManagerUsecase,
		UploadFile:             pmhOpts.UploadFile,
	}
}

func (h *PharmacyManagerHandler) GetAllPharmacyManager(ctx *gin.Context) {
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

	params := entities.PharmacyManagerParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	managers, pagination, err := h.PharmacyManagerUsecase.GetAllPharmacyManager(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyManagerResponses(managers, *pagination),
	})
}

func (h *PharmacyManagerHandler) GetPharmacyManagerById(ctx *gin.Context) {
	pmId, _ := strconv.Atoi(ctx.Param("id"))

	pharmacyManager, err := h.PharmacyManagerUsecase.GetPharmacyManagerById(ctx, int64(pmId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyManagerResponse(pharmacyManager),
	})
}

func (h *PharmacyManagerHandler) UpdatePharmacyManager(ctx *gin.Context) {
	var payload dtos.PharmacyMangerUpdateRequest

	getFile, _ := ctx.FormFile("logo")
	payload.Logo = getFile

	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.Error(err)
		return
	}

	pmId, _ := strconv.Atoi(ctx.Param("id"))

	var pharmacyManager entities.PharmacyManager

	if payload.Logo != nil {
		fileUrl, err := utils.GetFileUrl(ctx, payload.Logo, "png")
		if err != nil {
			ctx.Error(err)
			return
		}
		pharmacyManager.Logo = *fileUrl
	} else {
		pharmacyManager.Logo = ""
	}

	pharmacyManager.Id = int64(pmId)
	pharmacyManager.Name = payload.Name
	pharmacyManager.PhoneNumber = payload.PhoneNumber

	err := h.PharmacyManagerUsecase.UpdatePharmacyManager(ctx, pharmacyManager)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *PharmacyManagerHandler) DeletePharmacyMaagerById(ctx *gin.Context) {
	pmId, _ := strconv.Atoi(ctx.Param("id"))

	err := h.PharmacyManagerUsecase.DeletePharmacyManager(ctx, int64(pmId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

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

type DoctorHandlerOpts struct {
	DoctorUsecase usecases.DoctorUsecase
	UploadFile    utils.FileUploader
}

type DoctorHandler struct {
	DoctorUsecase usecases.DoctorUsecase
	UploadFile    utils.FileUploader
}

func NewDoctorHandler(doctorOpts *DoctorHandlerOpts) *DoctorHandler {
	return &DoctorHandler{
		DoctorUsecase: doctorOpts.DoctorUsecase,
		UploadFile:    doctorOpts.UploadFile,
	}
}

func (h *DoctorHandler) GetDoctorById(ctx *gin.Context) {
	doctorId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	doctor, err := h.DoctorUsecase.GetDoctorById(ctx, int64(doctorId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToDoctorResponse(*doctor),
	})
}

func (h *DoctorHandler) UpdateDoctor(ctx *gin.Context) {
	var payload dtos.DoctorUpdateRequest

	getFile, _ := ctx.FormFile("profile_picture")
	payload.ProfilePicture = getFile

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	doctorId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	var doctor entities.Doctor

	if payload.Certificate != nil {
		fileUrl, err := utils.GetFileUrl(ctx, payload.Certificate, "pdf")
		if err != nil {
			ctx.Error(err)
			return
		}
		doctor.Certificate = *utils.StringToNullString(*fileUrl)
	}

	if payload.ProfilePicture != nil {
		fileUrl, err := utils.GetFileUrl(ctx, payload.ProfilePicture, "png")
		if err != nil {
			ctx.Error(err)
			return
		}
		doctor.ProfilePicture = *utils.StringToNullString(*fileUrl)
	}

	doctor.Id = int64(doctorId)
	doctor.Name = payload.Name
	doctor.Fee = *utils.Int64ToNullInt64(int64(payload.Fee))
	doctor.WorkStartYear = *utils.Int64ToNullInt64(int64(payload.WorkStartYear))
	doctor.Specialist = &entities.DoctorSpecialist{Id: *utils.Int64ToNullInt64(payload.SpecialistId)}

	err = h.DoctorUsecase.UpdateDoctor(ctx, doctor)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    dtos.ConvertToDoctorResponse(doctor),
	})
}

func (h *DoctorHandler) DeleteDoctor(ctx *gin.Context) {
	doctorId, _ := strconv.Atoi(ctx.Param("id"))

	err := h.DoctorUsecase.DeleteDoctor(ctx, int64(doctorId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *DoctorHandler) GetAllDoctor(ctx *gin.Context) {
	var err error

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	specialistId := constants.DefaultId
	limit := constants.DefaultLimit
	page := constants.DefaultPage

	specialistIdStr := ctx.Query("specialistId")
	if specialistIdStr != "" {
		specialistId, err = strconv.Atoi(specialistIdStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
			return
		}
	}

	isOnline := ctx.Query("is-online")
	if isOnline != "" {
		_, err = strconv.ParseBool(isOnline)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidBooleanInputErrMsg))
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

	params := entities.DoctorParams{
		SpecialistId: int64(specialistId),
		IsOnline:     isOnline,
		SortBy:       sortBy,
		Sort:         sort,
		Limit:        limit,
		Page:         page,
		Keyword:      keyword,
	}

	var doctors []entities.Doctor
	var pagination *entities.PaginationInfo

	var datas *utils.ClaimsData
	datas = nil

	datas, _ = utils.GetDataFromContext(ctx)
	if datas != nil && datas.Role == constants.AdminRole {
		doctors, pagination, err = h.DoctorUsecase.GetAllDoctor(ctx, params, false)
		if err != nil {
			ctx.Error(err)
			return
		}
	} else {
		doctors, pagination, err = h.DoctorUsecase.GetAllDoctor(ctx, params, true)
		if err != nil {
			ctx.Error(err)
			return
		}
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToDoctorResponses(doctors, *pagination),
	})
}

func (h *DoctorHandler) ToggleDoctorIsOnline(ctx *gin.Context) {
	doctorId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.DoctorUsecase.ToggleDoctorIsOnline(ctx, int64(doctorId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

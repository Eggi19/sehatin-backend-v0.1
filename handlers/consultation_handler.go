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

type ConsultationHandlerOpts struct {
	ConsultationUsecase usecases.ConsultationUsecase
}

type ConsultationHandler struct {
	ConsultationUsecase usecases.ConsultationUsecase
}

func NewConsultationHandler(chOpts *ConsultationHandlerOpts) *ConsultationHandler {
	return &ConsultationHandler{
		ConsultationUsecase: chOpts.ConsultationUsecase,
	}
}

func (h *ConsultationHandler) GetAllConsultationByUser(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	params := entities.ConsultationParams{}

	status := ctx.Query("status")

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

	params.Page = page
	params.Limit = limit
	params.Status = status

	consultations, pagination, err := h.ConsultationUsecase.GetAllConsultationByUser(ctx, int64(userId), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToConsultationResponses(consultations, *pagination),
	})
}

func (h *ConsultationHandler) GetAllConsultationByDoctor(ctx *gin.Context) {
	doctorId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	params := entities.ConsultationParams{}

	status := ctx.Query("status")

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

	params.Page = page
	params.Limit = limit
	params.Status = status

	consultations, pagination, err := h.ConsultationUsecase.GetAllConsultationByDoctor(ctx, int64(doctorId), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToConsultationResponses(consultations, *pagination),
	})
}

func (h *ConsultationHandler) GetConsultationById(ctx *gin.Context) {
	consultationId, _ := strconv.Atoi(ctx.Param("id"))

	consultation, err := h.ConsultationUsecase.GetConsultationById(ctx, int64(consultationId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToConsultationResponse(*consultation),
	})
}

func (h *ConsultationHandler) CreateConsultation(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	var payload dtos.ConsultationRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	consultation := entities.Consultation{
		User:             entities.User{Id: int64(userId), Gender: &entities.Gender{}},
		Doctor:           entities.Doctor{Id: payload.DoctorId, Specialist: &entities.DoctorSpecialist{}},
		PatientGender:    entities.Gender{Id: payload.PatientGenderId},
		PatientName:      payload.PatientName,
		PatientBirthDate: payload.PatientBirthDate,
	}

	newC, err := h.ConsultationUsecase.CreateConsultation(ctx, consultation)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    dtos.ConvertToConsultationResponse(*newC),
	})
}

func (h *ConsultationHandler) EndConsultation(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, "")
	if err != nil {
		ctx.Error(err)
		return
	}

	paramId := ctx.Param("id")
	consultationId, err := strconv.Atoi(paramId)
	if err != nil {
		_ = ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
		return
	}

	err = h.ConsultationUsecase.EndConsultation(ctx, int64(consultationId), int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *ConsultationHandler) CreateChat(ctx *gin.Context) {
	var payload dtos.ChatRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	consultationId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	chat := entities.Chat{
		ConsultationId: int64(consultationId),
		IsFromUser:     datas.Role == constants.UserRole,
		Content:        payload.Content,
		Type:           payload.Type,
	}

	err = h.ConsultationUsecase.CreateChat(ctx, chat, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *ConsultationHandler) CreateChatFile(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(custom_errors.FileRequired())
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	consultationId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	chat := entities.Chat{
		ConsultationId: int64(consultationId),
		IsFromUser:     datas.Role == constants.UserRole,
		Type:           "file",
		File:           file,
	}

	err = h.ConsultationUsecase.CreateChat(ctx, chat, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *ConsultationHandler) CreatePrescription(ctx *gin.Context) {
	var payload dtos.PrescriptionRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	consultationId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	prescriptionData := entities.PrescriptionData{
		ConsultationId: int64(consultationId),
		Quantities:     payload.Quantities,
	}

	products := []entities.Product{}

	for i := 0; i < len(payload.Products); i++ {
		products = append(products, entities.Product{
			Id: payload.Products[i],
		})
	}

	prescriptionData.Products = products

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	fileUrl, err := h.ConsultationUsecase.CreatePrescription(ctx, prescriptionData, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data: dtos.PrescriptionUrlResponse{
			PrescriptionUrl: fileUrl,
		},
	})
}

func (h *ConsultationHandler) CreateCertificate(ctx *gin.Context) {
	var payload dtos.CertificateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	consultationId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	certificateData := entities.CertificateData{
		ConsultationId: int64(consultationId),
		Diagnosis:      payload.Diagnosis,
		StartDate:      payload.StartDate,
		EndDate:        payload.EndDate,
		PatientAge:     payload.PatientAge,
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	fileUrl, err := h.ConsultationUsecase.CreateCertificate(ctx, certificateData, datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data: dtos.CertificateUrlResponse{
			CertificateUrl: fileUrl,
		},
	})
}

func (h *ConsultationHandler) AddPrescriptionToCart(ctx *gin.Context) {
	consultationId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	datas, err := utils.GetDataFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.ConsultationUsecase.AddPrescriptionToCart(ctx, int64(consultationId), datas.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgAddedPrescriptionToCart,
		Data:    nil,
	})
}

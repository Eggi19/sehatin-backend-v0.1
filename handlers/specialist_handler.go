package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/gin-gonic/gin"
)

type SpecialistHandlerOpts struct {
	SpecialistUsecase usecases.SpecialistUsecase
}

type SpecialistHandler struct {
	SpecialistUsecase usecases.SpecialistUsecase
}

func NewSpecialistHandler(spHandOpts *SpecialistHandlerOpts) *SpecialistHandler {
	return &SpecialistHandler{
		SpecialistUsecase: spHandOpts.SpecialistUsecase,
	}
}

func (h *SpecialistHandler) GetAllSpecialist(ctx *gin.Context) {
	specialists, err := h.SpecialistUsecase.GetAllSpecialist(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToSpecialistResponses(specialists),
	})
}

package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

type MutationSatusHandlerOpts struct {
	MutationSatusUsecase usecases.MutationSatusUsecase
}

type MutationStatusHandler struct {
	MutationSatusUsecase usecases.MutationSatusUsecase
}

func NewMutationStatusHandler(msh *MutationSatusHandlerOpts) *MutationStatusHandler {
	return &MutationStatusHandler{MutationSatusUsecase: msh.MutationSatusUsecase}
}

func (h *MutationStatusHandler) GetAllMutationStatus(ctx *gin.Context) {
	mutationStatuses, err := h.MutationSatusUsecase.GetAll(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToMutationStatusResponses(mutationStatuses),
	})
}

func (h *MutationStatusHandler) GetOneMutationStatus(ctx *gin.Context) {
	id, err := utils.GetIdParamOrContext(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	mutationStatus, err := h.MutationSatusUsecase.GetMutationStatusById(ctx, int64(id))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToMutationStatusResponse(mutationStatus),
	})
}

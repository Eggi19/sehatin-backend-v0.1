package handlers

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/gin-gonic/gin"
)

type ProductFieldHandlerOpts struct {
	ProductFieldUsecase usecases.ProductFieldUsecase
}

type ProductFieldHandler struct {
	ProductFieldUsecase usecases.ProductFieldUsecase
}

func NewProductFieldHandler(pfhOpts *ProductFieldHandlerOpts) *ProductFieldHandler {
	return &ProductFieldHandler{ProductFieldUsecase: pfhOpts.ProductFieldUsecase}
}

func (h *ProductFieldHandler) GetAllForm(ctx *gin.Context) {
	forms, err := h.ProductFieldUsecase.GetAllForm(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToProductFieldResponses(forms, nil, nil),
	})
}

func (h *ProductFieldHandler) GetAllClassification(ctx *gin.Context) {
	classifications, err := h.ProductFieldUsecase.GetAllClassification(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToProductFieldResponses(nil, classifications, nil),
	})
}

func (h *ProductFieldHandler) GetAllManufacture(ctx *gin.Context) {
	manufactures, err := h.ProductFieldUsecase.GetAllManufacture(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToProductFieldResponses(nil, nil, manufactures),
	})

}

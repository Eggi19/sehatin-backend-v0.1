package handlers

import (
	"net/http"
	"strconv"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/usecases"
	"github.com/gin-gonic/gin"
)

type LocationHandlerOpts struct {
	LocationUsecase usecases.LocationUsecase
}

type LocationHandler struct {
	LocationUsecase usecases.LocationUsecase
}

func NewLocationHandler(lhOpts *LocationHandlerOpts) *LocationHandler {
	return &LocationHandler{lhOpts.LocationUsecase}
}

func (h *LocationHandler) GetAllProvinces(ctx *gin.Context) {
	provinces, err := h.LocationUsecase.GetAllProvinces(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    provinces,
	})
}

func (h *LocationHandler) GetCitiesByProvinceId(ctx *gin.Context) {
	provinceId, _ := strconv.Atoi(ctx.Param("id"))

	cities, err := h.LocationUsecase.GetCitiesByProvinceId(ctx, int16(provinceId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    cities,
	})
}

func (h *LocationHandler) GetDistrictsByCityId(ctx *gin.Context) {
	cityId, _ := strconv.Atoi(ctx.Param("id"))

	districts, err := h.LocationUsecase.GetDistrictsByCityId(ctx, int16(cityId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    districts,
	})
}

func (h *LocationHandler) GetSubDistrictsByDistrictId(ctx *gin.Context) {
	districtId, _ := strconv.Atoi(ctx.Param("id"))

	subDistricts, err := h.LocationUsecase.GetSubDistrictsByDistrictId(ctx, int16(districtId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    subDistricts,
	})
}

func (h *LocationHandler) ReverseCoordinate(ctx *gin.Context) {
	latitude := ctx.Query("lat")
	longitude := ctx.Query("lon")

	res, err := h.LocationUsecase.ReverseCoordinate(ctx, latitude, longitude)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    res,
	})
}

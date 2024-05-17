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

type PharmacyHandlerOpts struct {
	PharmacyUsecase usecases.PharmacyUsecase
}

type PharmacyHandler struct {
	PharmacyUsecase usecases.PharmacyUsecase
}

func NewPharmacyHandler(phOpts *PharmacyHandlerOpts) *PharmacyHandler {
	return &PharmacyHandler{
		PharmacyUsecase: phOpts.PharmacyUsecase,
	}
}

func (h *PharmacyHandler) GetPharmacyById(ctx *gin.Context) {
	pharmacyId, _ := strconv.Atoi(ctx.Param("id"))

	pharmacy, err := h.PharmacyUsecase.GetPharmacyById(ctx, int64(pharmacyId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyResponse(pharmacy),
	})
}

func (h *PharmacyHandler) CreatePharmacy(ctx *gin.Context) {
	var payload dtos.PharmacyCreateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	pharmacyManagerId, err := utils.GetIdParamOrContext(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	pharmacy := entities.Pharmacy{
		Name:                    payload.Name,
		OperationalHour:         payload.OperationalHour,
		OperationalDay:          payload.OperationalDay,
		PharmacistName:          payload.PharmacistName,
		PharmacistLicenseNumber: payload.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   payload.PharmacistPhoneNumber,
		PharmacyAddress: entities.PharmacyAddress{
			City:        payload.City,
			Province:    payload.Province,
			Address:     payload.Address,
			District:    payload.District,
			SubDistrict: payload.SubDistrict,
			PostalCode:  payload.PostalCode,
			Longitude:   payload.Longitude,
			Latitude:    payload.Latitude,
		},
		PharmacyManager: entities.PharmacyManager{Id: int64(pharmacyManagerId)},
	}

	officialShipping := []entities.OfficialShippingMethod{}
	nonOfficialShipping := []entities.NonOfficialShippingMethod{}

	for _, official := range payload.OfficialShippingId {
		officialShipping = append(officialShipping, entities.OfficialShippingMethod{Id: official})
	}

	for _, nonOfficial := range payload.NonOfficialShippingId {
		nonOfficialShipping = append(nonOfficialShipping, entities.NonOfficialShippingMethod{Id: nonOfficial})
	}

	pharmacy.OfficialShippingMethod = officialShipping
	pharmacy.NonOfficialShippingMethod = nonOfficialShipping

	err = h.PharmacyUsecase.CreatePharmacy(ctx, pharmacy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, dtos.ResponseMessage{
		Message: constants.ResponseMsgCreated,
		Data:    nil,
	})
}

func (h *PharmacyHandler) UpdatePharmacy(ctx *gin.Context) {
	var payload dtos.PharmacyUpdateRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.Error(err)
		return
	}

	pId, _ := strconv.Atoi(ctx.Param("id"))

	pharmacy := entities.Pharmacy{
		Id:                      int64(pId),
		Name:                    payload.Name,
		OperationalHour:         payload.OperationalHour,
		OperationalDay:          payload.OperationalDay,
		PharmacistName:          payload.PharmacistName,
		PharmacistLicenseNumber: payload.PharmacistLicenseNumber,
		PharmacistPhoneNumber:   payload.PharmacistPhoneNumber,
		PharmacyAddress: entities.PharmacyAddress{
			Id:          0,
			PharmacyId:  int64(pId),
			City:        payload.City,
			Province:    payload.Province,
			Address:     payload.Address,
			District:    payload.District,
			SubDistrict: payload.SubDistrict,
			PostalCode:  payload.PostalCode,
			Longitude:   payload.Longitude,
			Latitude:    payload.Latitude,
		},
	}

	officialShipping := []entities.OfficialShippingMethod{}
	nonOfficialShipping := []entities.NonOfficialShippingMethod{}

	for _, official := range payload.OfficialShippingId {
		officialShipping = append(officialShipping, entities.OfficialShippingMethod{Id: official})
	}

	for _, nonOfficial := range payload.NonOfficialShippingId {
		nonOfficialShipping = append(nonOfficialShipping, entities.NonOfficialShippingMethod{Id: nonOfficial})
	}

	pharmacy.OfficialShippingMethod = officialShipping
	pharmacy.NonOfficialShippingMethod = nonOfficialShipping

	err := h.PharmacyUsecase.UpdatePharmacy(ctx, pharmacy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    nil,
	})
}

func (h *PharmacyHandler) DeletePharmacyById(ctx *gin.Context) {
	pId, _ := strconv.Atoi(ctx.Param("id"))

	err := h.PharmacyUsecase.DeletePharmacyById(ctx, int64(pId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *PharmacyHandler) GetAllPharmacyByPharmacyManager(ctx *gin.Context) {
	pharmacyManagerId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

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

	params := entities.PharmacyParams{
		SortBy:  sortBy,
		Sort:    sort,
		Limit:   limit,
		Page:    page,
		Keyword: keyword,
	}

	pharmacies, pagination, err := h.PharmacyUsecase.GetAllPharmacyByPharmacyManagerId(ctx, int64(pharmacyManagerId), params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToPharmacyResponses(pharmacies, *pagination),
	})
}

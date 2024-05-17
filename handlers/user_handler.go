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

type UserHandlerOpts struct {
	UserUsecase usecases.UserUsecase
	UploadFile  utils.FileUploader
}

type UserHandler struct {
	UserUsecase usecases.UserUsecase
	UploadFile  utils.FileUploader
}

func NewUserHandler(uhOpts *UserHandlerOpts) *UserHandler {
	return &UserHandler{
		UserUsecase: uhOpts.UserUsecase,
		UploadFile:  uhOpts.UploadFile,
	}
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.UserUsecase.GetUserById(ctx, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToUserResponse(*user),
	})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = h.UserUsecase.DeleteUser(ctx, int64(userId))
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgDeleted,
		Data:    nil,
	})
}

func (h *UserHandler) GetAllUser(ctx *gin.Context) {
	var err error

	sortBy := ctx.Query("sortBy")
	sort := ctx.Query("sort")
	keyword := ctx.Query("keyword")

	genderId := constants.DefaultId
	limit := constants.DefaultLimit
	page := constants.DefaultPage

	genderIdStr := ctx.Query("gender-id")
	if genderIdStr != "" {
		genderId, err = strconv.Atoi(genderIdStr)
		if err != nil {
			ctx.Error(custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg))
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

	params := entities.UserParams{
		GenderId: int64(genderId),
		SortBy:   sortBy,
		Sort:     sort,
		Limit:    limit,
		Page:     page,
		Keyword:  keyword,
	}

	users, pagination, err := h.UserUsecase.GetAllUser(ctx, params)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgOK,
		Data:    dtos.ConvertToUserResponses(users, *pagination),
	})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var payload dtos.UserUpdateRequest

	getFile, _ := ctx.FormFile("profile_picture")
	payload.ProfilePicture = getFile

	if err := ctx.ShouldBind(&payload); err != nil {
		_ = ctx.Error(err)
		return
	}

	userId, err := utils.GetIdParamOrContext(ctx, constants.Id)
	if err != nil {
		ctx.Error(err)
		return
	}

	var user entities.User

	if payload.ProfilePicture != nil {
		fileUrl, err := utils.GetFileUrl(ctx, payload.ProfilePicture, "png")
		if err != nil {
			ctx.Error(err)
			return
		}
		user.ProfilePicture = *utils.StringToNullString(*fileUrl)
	}

	user.Id = int64(userId)
	user.Name = payload.Name
	user.Gender = &entities.Gender{Id: payload.GenderId, Name: ""}
	user.Address = []entities.UserAddress{}
	user.BirthDate = *utils.StringToNullString(payload.BirthDate)

	err = h.UserUsecase.UpdateUser(ctx, user)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, dtos.ResponseMessage{
		Message: constants.ResponseMsgUpdated,
		Data:    dtos.ConvertToUserResponse(user),
	})
}

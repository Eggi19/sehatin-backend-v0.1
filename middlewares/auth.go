package middlewares

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/tsanaativa/sehatin-backend-v0.1/utils"
	"github.com/gin-gonic/gin"
)

func SetPublic() func(*gin.Context) {
	return func(ctx *gin.Context) {
		data := &utils.ClaimsData{
			Id:   0,
			Role: "public",
		}
		ctx.Set("data", data)
		ctx.Next()
	}
}

func JwtAuthMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, data, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil && data == nil {
			if err.Error() == custom_errors.TokenExpired().Error() {
				ctx.AbortWithStatusJSON(http.StatusForbidden, dtos.ErrResponse{
					Message: constants.ExpiredTokenErrMsg,
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Set("data", data)
		ctx.Next()
	}
}

func JwtAdminAuthMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, _, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		err = utils.NewJwtProvider(config).ValidateAdminRoleJwt(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Next()
	}
}

func JwtUserAuthMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, _, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		err = utils.NewJwtProvider(config).ValidateUserRoleJwt(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Next()
	}
}

func JwtPharmacyManagerMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, _, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		err = utils.NewJwtProvider(config).ValidatePharmacyManagerRoleJwt(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Next()
	}
}

func JwtDoctorMiddleware(config utils.Config) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, _, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		err = utils.NewJwtProvider(config).ValidateDoctorRoleJwt(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Next()
	}
}

func JwtMultiRoleMiddleware(config utils.Config, roles []string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		authorized, _, err := utils.NewJwtProvider(config).IsAuthorized(ctx)
		if !authorized && err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		err = utils.NewJwtProvider(config).ValidateMultiRoleJwt(ctx, roles)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dtos.ErrResponse{
				Message: constants.ResponseMsgUnauthorized,
			})
			return
		}
		ctx.Next()
	}
}

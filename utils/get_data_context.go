package utils

import (
	"strconv"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/custom_errors"
	"github.com/gin-gonic/gin"
)

func GetDataFromContext(ctx *gin.Context) (*ClaimsData, error) {
	datas, exists := ctx.Get("data")
	if !exists {
		return nil, custom_errors.ContextNotFound()
	}

	return datas.(*ClaimsData), nil
}

func GetIdParamOrContext(ctx *gin.Context, key string) (int, error) {
	valueStr, exists := ctx.Params.Get(key)
	if exists {
		valueInt, err := strconv.Atoi(valueStr)
		if err != nil {
			return 0, custom_errors.BadRequest(err, constants.InvalidIntegerInputErrMsg)
		}
		return valueInt, nil
	}

	datas, err := GetDataFromContext(ctx)
	if err != nil {
		return 0, custom_errors.Forbidden()
	}
	return int(datas.Id), nil
}

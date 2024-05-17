package middlewares

import (
	"net/http"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/tsanaativa/sehatin-backend-v0.1/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestId(c *gin.Context) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.ErrResponse{
			Message: constants.ResponseMsgErrorInternalServer,
		})
		return
	}

	c.Set(constants.RequestId, uuid)
	c.Next()
}

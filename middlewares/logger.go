package middlewares

import (
	"time"

	"github.com/tsanaativa/sehatin-backend-v0.1/constants"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger(log *logrus.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()

		requestId, exist := c.Get(constants.RequestId)
		if !exist {
			requestId = ""
		}

		entry := log.WithFields(logrus.Fields{
			"request_id":  requestId,
			"latency":     time.Since(start),
			"method":      c.Request.Method,
			"status_code": statusCode,
			"path":        path,
		})

		if statusCode >= 400 && statusCode <= 599 {
			firstErr := c.Errors
			entry.WithField("error", firstErr).Error(firstErr.String())
			return
		}

		entry.Info("request processed")
	}

}

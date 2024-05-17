package middlewares

import (
	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
	} else {
		c.Next()
	}
}

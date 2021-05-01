package middleware

import (
	"github.com/gin-gonic/gin"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: log request
		c.Next()
	}
}

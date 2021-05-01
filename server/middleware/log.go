package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Log(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: log request
		c.Next()
	}
}

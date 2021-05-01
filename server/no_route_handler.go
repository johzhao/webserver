package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func NoRouteHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		logger.Error("no route for request",
			zap.String("method", request.Method),
			zap.String("url", request.URL.String()))

		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("no route for request with method: %s, url: %s", request.Method, request.URL.String()),
			"data":    nil,
		})
	}
}

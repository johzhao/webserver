package middleware

import (
	"fmt"
	"net/http"
	"webserver/errors"
	"webserver/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NoRouteHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		request := c.Request
		logger.
			With(logging.ContextField(c)...).
			Error("no route for request",
				zap.String("method", request.Method),
				zap.String("url", request.URL.String()),
			)

		msg := fmt.Sprintf("no route for request with method: %s, url: %s", request.Method, request.URL.String())
		c.JSON(http.StatusNotFound, gin.H{
			"code":    errors.ErrRouteNotFound,
			"message": msg,
			"data":    nil,
		})
	}
}

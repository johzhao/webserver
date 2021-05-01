package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    -1,
			"message": fmt.Sprintf("no route for request with method: %s, url: %s", c.Request.Method, c.Request.URL.String()),
			"data":    nil,
		})
	}
}

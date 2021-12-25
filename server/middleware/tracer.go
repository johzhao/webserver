package middleware

import (
	"strings"
	"webserver/logging"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

func Tracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracerID := c.Request.Header.Get(logging.TracerIDKey)
		if len(tracerID) == 0 {
			tracerID = strings.ReplaceAll(uuid.NewV4().String(), "-", "")
		}

		c.Set(logging.TracerIDKey, tracerID)
	}
}

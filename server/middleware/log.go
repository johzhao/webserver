package middleware

import (
	"time"
	"webserver/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Log(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		beginTime := time.Now()

		c.Next()

		endTime := time.Now()
		elapsed := endTime.Sub(beginTime)
		elapsedInMicroSecond := float64(elapsed / time.Microsecond)
		elapsedInMillSecond := elapsedInMicroSecond / 1000.0

		logger.
			With(logging.ContextField(c)...).
			Info("handle request",
				zap.String("method", c.Request.Method),
				zap.String("url", c.Request.URL.String()),
				zap.Float64("cost(ms)", elapsedInMillSecond),
			)
	}
}

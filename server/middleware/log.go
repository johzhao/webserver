package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func Log(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: log request
		beginTime := time.Now()
		c.Next()
		endTime := time.Now()
		elapsed := endTime.Sub(beginTime)
		elapsedInMicroSecond := float64(elapsed / time.Microsecond)
		elapsedInMillSecond := elapsedInMicroSecond / 1000.0
		logger.Info("handle request",
			zap.String("method", c.Request.Method),
			zap.String("url", c.Request.URL.String()),
			zap.Float64("cost in ms", elapsedInMillSecond))
	}
}

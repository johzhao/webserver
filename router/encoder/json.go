package encoder

import (
	"net/http"
	"webserver/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewJSONResponseEncoder(logger *zap.Logger) ResponseEncoder {
	return &JSONResponseEncoder{
		logger: logger,
	}
}

type JSONResponseEncoder struct {
	logger *zap.Logger
}

func (e *JSONResponseEncoder) ResponseWithData(ctx *gin.Context, data interface{}) error {
	ctx.JSON(http.StatusOK, data)
	return nil
}

func (e *JSONResponseEncoder) ResponseWithError(ctx *gin.Context, err error) {
	e.logger.
		With(logging.ContextField(ctx)...).
		Error("request failed",
			zap.String("method", ctx.Request.Method),
			zap.String("url", ctx.Request.URL.String()),
			zap.Error(err),
		)

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}

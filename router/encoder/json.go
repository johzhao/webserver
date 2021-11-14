package encoder

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewJsonResponseEncoder() ResponseEncoder {
	return &JsonResponseEncoder{}
}

type JsonResponseEncoder struct {
}

func (e *JsonResponseEncoder) ResponseWithData(ctx *gin.Context, data interface{}) error {
	responseBody := make(map[string]interface{}, 3)
	responseBody["code"] = 0
	responseBody["message"] = "success"
	if data != nil {
		responseBody["data"] = data
	}

	ctx.JSON(http.StatusOK, responseBody)

	return nil
}

func (e *JsonResponseEncoder) ResponseWithError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    -1,
		"message": err.Error(),
	})
}

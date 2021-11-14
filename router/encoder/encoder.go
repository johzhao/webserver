package encoder

import "github.com/gin-gonic/gin"

type ResponseEncoder interface {
	ResponseWithData(ctx *gin.Context, data interface{}) error

	ResponseWithError(ctx *gin.Context, err error)
}

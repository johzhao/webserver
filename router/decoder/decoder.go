package decoder

import "github.com/gin-gonic/gin"

type Decoder func(ctx *gin.Context) (interface{}, error)

type RequestDecoder interface {
	DecodeRequest(ctx *gin.Context) (interface{}, error)
}

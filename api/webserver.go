package api

import (
	"github.com/gin-gonic/gin"
)

type ServerHandlerFunc func(ctx *gin.Context) (resp interface{}, err error)

type WebServer interface {
	SetupServer() error
	SetupRoute(httpMethod string, relativePath string, handler ServerHandlerFunc)
	RunServer() error
}

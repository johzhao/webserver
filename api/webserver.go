package api

import (
	"github.com/gin-gonic/gin"
	"webserver/server"
)

type ServerHandlerFunc func(ctx *gin.Context) (resp interface{}, err error)

type WebServer interface {
	SetupServer() error
	AddRoute(conf *server.RouteConfig)
	RunServer() error
	StopServer(err error)
}

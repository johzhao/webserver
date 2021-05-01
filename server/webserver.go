package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webserver/api"
	"webserver/server/middleware"
	"webserver/user"
)

func NewWebServer(logger *zap.Logger, userController user.Controller) api.WebServer {
	return &webServer{
		logger:         logger,
		userController: userController,
	}
}

type webServer struct {
	engine *gin.Engine

	logger         *zap.Logger
	userController user.Controller
}

func (s *webServer) SetupServer() error {
	engine := gin.New()
	engine.Use(
		middleware.Log(s.logger),
		middleware.Recovery(s.logger),
	)
	s.engine = engine

	s.userController.SetupRoute(s)

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.engine.NoRoute(NoRouteHandler(s.logger))

	return nil
}

func (s webServer) SetupRoute(httpMethod string, relativePath string, handler api.ServerHandlerFunc) {
	s.engine.Handle(httpMethod, relativePath, DefaultJSONEncode(handler))
}

func DefaultJSONEncode(handler api.ServerHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := handler(c)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    -1,
				"message": err.Error(),
				"data":    nil,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data":    data,
			})
		}
	}
}

func (s webServer) RunServer() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.engine,
	}
	return serveGracefulShutdownServer(srv)
}

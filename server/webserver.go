package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webserver/controller"
	"webserver/server/middleware"
	"webserver/transport"
)

type WebServer interface {
	SetupServer() error
	RunServer() error
	StopServer(err error)
}

func NewWebServer(logger *zap.Logger, userController controller.Controller) WebServer {
	return &webServer{
		logger:         logger,
		userController: userController,
	}
}

type webServer struct {
	engine         *gin.Engine
	srv            *http.Server
	logger         *zap.Logger
	userController controller.Controller
}

func (s *webServer) SetupServer() error {
	engine := gin.New()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	engine.Use(
		cors.New(corsConfig),
		middleware.Log(s.logger),
		middleware.Recovery(s.logger),
	)
	s.engine = engine

	for _, routeConfig := range s.userController.Routes() {
		s.engine.Handle(routeConfig.Method, routeConfig.Path, transport.MakeJsonRouteHandler(routeConfig))
	}

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	s.engine.NoRoute(NoRouteHandler(s.logger))

	return nil
}

func (s *webServer) RunServer() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.engine,
	}
	s.srv = srv
	return s.srv.ListenAndServe()
}

//goland:noinspection GoUnusedParameter
func (s *webServer) StopServer(err error) {
	s.logger.Info("stop server")
	_ = s.srv.Shutdown(context.Background())
	s.srv = nil
}

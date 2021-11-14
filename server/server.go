package server

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webserver/router"
	"webserver/server/middleware"
)

func NewServer(logger *zap.Logger) Server {
	return Server{
		logger: logger,
	}
}

type Server struct {
	logger *zap.Logger
	engine *gin.Engine
	srv    *http.Server
}

func (s *Server) SetupServer() error {
	engine := gin.New()
	engine.Use(
		cors.New(s.corsConfig()),
		middleware.Log(s.logger),
		middleware.Recovery(s.logger),
	)
	engine.NoRoute(middleware.NoRouteHandler(s.logger))

	s.engine = engine

	return nil
}

func (s *Server) corsConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true

	return corsConfig
}

func (s *Server) RunServer() error {
	srv := &http.Server{
		Addr:    ":8080", // TODO: get from config
		Handler: s.engine,
	}
	s.srv = srv
	return s.srv.ListenAndServe()
}

//goland:noinspection GoUnusedParameter
func (s *Server) StopServer(err error) {
	s.logger.Info("stop server")
	_ = s.srv.Shutdown(context.Background())
	s.srv = nil
}

func (s *Server) HandleRouter(router router.Router) {
	s.engine.Handle(router.Method(), router.Path(), router.HandleRequest)
}

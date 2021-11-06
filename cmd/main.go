package main

import (
	"go.uber.org/zap"
	"os"
	"webserver/controller"
	"webserver/database/repository"
	"webserver/logger"
	"webserver/server"
	"webserver/service"
	tracerCreator "webserver/tracing/creator"
	"webserver/utility"
)

func main() {
	zapLogger := logger.SetupLogger()
	tracer, err := tracerCreator.NewTracer("webserver", "", zapLogger)
	if err != nil {
		os.Exit(1)
	}
	defer tracer.Close()

	userRepository := repository.NewUserRepository(zapLogger)
	userService := service.NewUserService(userRepository, zapLogger)
	userController := controller.NewUserController(userService, zapLogger)

	webServer := server.NewWebServer(zapLogger, userController)
	if err := webServer.SetupServer(); err != nil {
		zapLogger.Info("setup server failed", zap.Error(err))
		os.Exit(1)
	}

	zapLogger.Info("start server")

	g := utility.MakeGroup()
	g.Add(webServer.RunServer, webServer.StopServer)

	if err := g.Run(); err != nil {
		zapLogger.Info("run failed", zap.Error(err))
		os.Exit(1)
	}
}

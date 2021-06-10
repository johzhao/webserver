package main

import (
	"go.uber.org/zap"
	"os"
	"webserver/logger"
	"webserver/server"
	tracerCreator "webserver/tracing/creator"
	"webserver/user"
)

func main() {
	zapLogger := logger.SetupLogger()
	tracer, err := tracerCreator.NewTracer("webserver", "", zapLogger)
	if err != nil {
		os.Exit(1)
	}
	defer tracer.Close()

	userRepository := user.NewUserRepository(zapLogger)
	userService := user.NewUserService(userRepository, zapLogger)
	userController := user.NewUserController(userService, zapLogger)

	webServer := server.NewWebServer(zapLogger, userController)
	if err := webServer.SetupServer(); err != nil {
		zapLogger.Info("setup server failed", zap.Error(err))
		os.Exit(1)
	}

	zapLogger.Info("start server")

	g := server.MakeGroup()
	g.Add(webServer.RunServer, webServer.StopServer)

	if err := g.Run(); err != nil {
		zapLogger.Info("run failed", zap.Error(err))
		os.Exit(1)
	}
}

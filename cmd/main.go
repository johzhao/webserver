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
	"webserver/transport"
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
	svr := server.NewServer(zapLogger)

	if err := svr.SetupServer(); err != nil {
		zapLogger.Info("setup server failed", zap.Error(err))
		os.Exit(1)
	}

	transport.SetupRouters(svr, userController)

	zapLogger.Info("start server")

	g := utility.MakeGroup()
	g.Add(svr.RunServer, svr.StopServer)

	if err := g.Run(); err != nil {
		zapLogger.Info("run failed", zap.Error(err))
		os.Exit(1)
	}
}

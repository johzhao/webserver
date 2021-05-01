package main

import (
	"os"
	logger2 "webserver/logger"
	"webserver/server"
	"webserver/user"
)

func main() {
	zapLogger := logger2.SetupLogger()

	userRepository := user.NewUserRepository(zapLogger)
	userService := user.NewUserService(userRepository, zapLogger)
	userController := user.NewUserController(userService, zapLogger)

	webServer := server.NewWebServer(zapLogger, userController)
	if err := webServer.SetupServer(); err != nil {
		os.Exit(1)
	}

	zapLogger.Info("start server")

	if err := webServer.RunServer(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"
	"webserver/server"
	"webserver/user"
)

func main() {
	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	webServer := server.NewServer(userController)
	if err := webServer.SetupServer(); err != nil {
		os.Exit(1)
	}

	if err := webServer.RunServer(); err != nil {
		os.Exit(1)
	}
}

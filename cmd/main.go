package main

import (
	"github.com/gin-gonic/gin"
	"webserver/user"
)

func main() {
	userRepository := user.NewUserRepository()
	userService := user.NewUserService(userRepository)
	userController := user.NewUserController(userService)

	r := gin.Default()
	r.POST("/users", userController.CreateUser)
	r.PUT("/users/:user_id", userController.UpdateUser)
	r.GET("/users/:user_id", userController.GetUser)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	_ = r.Run()
}

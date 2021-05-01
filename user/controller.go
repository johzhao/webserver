package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/api"
	"webserver/api/user"
	"webserver/api/user/command"
)

func NewUserController(service user.Service) Controller {
	return Controller{
		userService: service,
	}
}

type Controller struct {
	userService user.Service
}

func (c Controller) SetupRoute(webServer api.WebServer) {
	webServer.SetupRoute(http.MethodPost, "/users", c.CreateUser)
	webServer.SetupRoute(http.MethodPut, "/users/:user_id", c.UpdateUser)
	webServer.SetupRoute(http.MethodGet, "/users/:user_id", c.GetUser)
}

func (c Controller) CreateUser(ctx *gin.Context) (interface{}, error) {
	cmd := command.CreateUserCommand{} // XXX: make the cmd
	return c.userService.CreateUser(ctx, cmd)
}

func (c Controller) UpdateUser(ctx *gin.Context) (interface{}, error) {
	userID := ctx.Params.ByName("user_id")
	cmd := command.UpdateUserCommand{
		UserID: userID,
	} // XXX: make the cmd
	err := c.userService.UpdateUser(ctx, cmd)
	return nil, err
}

func (c Controller) GetUser(ctx *gin.Context) (interface{}, error) {
	userID := ctx.Params.ByName("user_id")
	return c.userService.GetUser(ctx, userID)
}

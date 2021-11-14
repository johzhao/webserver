package transport

import (
	"net/http"
	"webserver/controller"
	"webserver/model/command"
	"webserver/model/query"
	"webserver/router"
)

func userRouters(userController controller.User) []router.Router {
	return []router.Router{
		router.NewJsonRouter(http.MethodPost, "/users", command.CreateUserCommand{}, userController.CreateUser),
		router.NewJsonRouter(http.MethodPut, "/users/:user_id", command.UpdateUserCommand{}, userController.UpdateUser),
		router.NewJsonRouter(http.MethodGet, "/users", query.GetUserQuery{}, userController.GetUser),
		router.NewJsonRouter(http.MethodGet, "/users/fail", nil, userController.FailedTest),
	}
}

package transport

import (
	"go.uber.org/zap"
	"net/http"
	"webserver/controller"
	"webserver/model/command"
	"webserver/model/query"
	"webserver/router"
)

func userRouters(logger *zap.Logger, userController controller.User) []router.Router {
	return []router.Router{
		router.NewJSONRouter(logger, http.MethodPost, "/users", command.CreateUserCommand{}, userController.CreateUser),
		router.NewJSONRouter(logger, http.MethodPut, "/users/:user_id", command.UpdateUserCommand{}, userController.UpdateUser),
		router.NewJSONRouter(logger, http.MethodGet, "/users", query.GetUserQuery{}, userController.GetUser),
		router.NewJSONRouter(logger, http.MethodGet, "/users/fail", nil, userController.FailedTest),
	}
}

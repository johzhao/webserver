package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/api"
	"webserver/api/user"
	"webserver/api/user/command"
	"webserver/api/user/query"
	"webserver/server"
)

func NewUserController(service user.Service, logger *zap.Logger) Controller {
	return Controller{
		userService: service,
		logger:      logger,
	}
}

type Controller struct {
	userService user.Service
	logger      *zap.Logger
}

func (c Controller) SetupRoute(webServer api.WebServer) {
	webServer.AddRoute(&server.RouteConfig{
		Path:          "/users",
		RequestObject: command.CreateUserCommand{},
		Handler:       c.CreateUser,
	})
	webServer.AddRoute(&server.RouteConfig{
		Method:        "PUT",
		Path:          "/users/:user_id",
		RequestObject: command.UpdateUserCommand{},
		Handler:       c.UpdateUser,
	})
	webServer.AddRoute(&server.RouteConfig{
		Method:        "GET",
		Path:          "/users",
		RequestObject: query.GetUserQuery{},
		Handler:       c.GetUser,
	})
}

//goland:noinspection GoUnusedParameter
func (c Controller) CreateUser(ctx context.Context, req interface{}) (interface{}, error) {
	cmd, ok := req.(command.CreateUserCommand)
	if !ok {
		return nil, fmt.Errorf("invalid request struct")
	}

	//return c.userService.CreateUser(ctx, cmd)
	return map[string]interface{}{
		"name": cmd.Name,
		"age":  cmd.Age,
	}, nil
}

//goland:noinspection GoUnusedParameter
func (c Controller) UpdateUser(ctx context.Context, req interface{}) (interface{}, error) {
	cmd, ok := req.(command.UpdateUserCommand)
	if !ok {
		return nil, fmt.Errorf("invalid request struct")
	}

	//err := c.userService.UpdateUser(ctx, cmd)
	return map[string]interface{}{
		"user_id": cmd.UserID,
	}, nil
}

//goland:noinspection GoUnusedParameter
func (c Controller) GetUser(ctx context.Context, req interface{}) (interface{}, error) {
	q, ok := req.(query.GetUserQuery)
	if !ok {
		return nil, fmt.Errorf("invalid request struct")
	}

	return map[string]interface{}{
		"user_ids": q.UserIDs,
	}, nil
}

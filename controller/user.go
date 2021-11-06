package controller

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/model/command"
	"webserver/model/query"
	"webserver/service"
	"webserver/transport"
)

func NewUserController(userService service.UserService, logger *zap.Logger) Controller {
	return Controller{
		userService: userService,
		logger:      logger,
	}
}

type Controller struct {
	userService service.UserService
	logger      *zap.Logger
}

func (c Controller) Routes() []*transport.JsonRouteConfig {
	return []*transport.JsonRouteConfig{
		{
			Method:        "POST",
			Path:          "/users",
			RequestObject: command.CreateUserCommand{},
			Handler:       c.CreateUser,
		},
		{
			Method:        "PUT",
			Path:          "/users/:user_id",
			RequestObject: command.UpdateUserCommand{},
			Handler:       c.UpdateUser,
		},
		{
			Method:        "GET",
			Path:          "/users",
			RequestObject: query.GetUserQuery{},
			Handler:       c.GetUser,
		},
	}
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

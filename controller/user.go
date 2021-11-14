package controller

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/model/command"
	"webserver/service"
)

func NewUserController(userService service.UserService, logger *zap.Logger) User {
	return User{
		userService: userService,
		logger:      logger,
	}
}

type User struct {
	userService service.UserService
	logger      *zap.Logger
}

//goland:noinspection GoUnusedParameter
func (c User) CreateUser(ctx context.Context, req interface{}) (interface{}, error) {
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
func (c User) UpdateUser(ctx context.Context, req interface{}) (interface{}, error) {
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
func (c User) GetUser(ctx context.Context, req interface{}) (interface{}, error) {
	//q, ok := req.(query.GetUserQuery)
	//if !ok {
	//	return nil, fmt.Errorf("invalid request struct")
	//}
	//
	//return map[string]interface{}{
	//	"user_ids": q.UserIDs,
	//}, nil
	return nil, nil
}

//goland:noinspection GoUnusedParameter
func (c User) FailedTest(ctx context.Context, req interface{}) (interface{}, error) {
	return nil, fmt.Errorf("something wrong here")
}

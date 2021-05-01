package user

import (
	"context"
	"webserver/api/user/command"
	"webserver/api/user/dto"
)

type Service interface {
	CreateUser(ctx context.Context, cmd command.CreateUserCommand) (string, error)
	UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error

	GetUser(ctx context.Context, userID string) (*dto.User, error)
}

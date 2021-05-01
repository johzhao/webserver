package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/user/model"
)

func NewUserRepository(logger *zap.Logger) Repository {
	return Repository{
		logger: logger,
	}
}

type Repository struct {
	logger *zap.Logger
}

func (r Repository) Save(ctx context.Context, user *model.User) (string, error) {
	return "", fmt.Errorf("need implement")
}

func (r Repository) GetUser(ctx context.Context, userID string) (*model.User, error) {
	return &model.User{
		ID:   "100",
		Name: "Username",
		Age:  40,
	}, nil
}

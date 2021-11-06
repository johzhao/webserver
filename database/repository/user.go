package repository

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/model/do"
)

func NewUserRepository(logger *zap.Logger) Repository {
	return Repository{
		logger: logger,
	}
}

type Repository struct {
	logger *zap.Logger
}

//goland:noinspection GoUnusedParameter
func (r Repository) Save(ctx context.Context, user *do.User) (string, error) {
	return "", fmt.Errorf("need implement")
}

//goland:noinspection GoUnusedParameter
func (r Repository) GetUser(ctx context.Context, userID string) (*do.User, error) {
	return &do.User{
		ID:   "100",
		Name: "Username",
		Age:  40,
	}, nil
}

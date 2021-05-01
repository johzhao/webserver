package user

import (
	"context"
	"fmt"
	"webserver/user/model"
)

func NewUserRepository() Repository {
	return Repository{}
}

type Repository struct {
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

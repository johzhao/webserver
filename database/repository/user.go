package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"webserver/model/do"
)

type UserRepository interface {
	Save(ctx context.Context, user *do.User) (string, error)
	GetUser(ctx context.Context, userID string) (*do.User, error)
}

func NewUserRepository(logger *zap.Logger, queryer sqlx.QueryerContext, execer sqlx.ExecerContext) UserRepository {
	return &userRepository{
		logger:  logger,
		queryer: queryer,
		execer:  execer,
	}
}

type userRepository struct {
	logger  *zap.Logger
	queryer sqlx.QueryerContext
	execer  sqlx.ExecerContext
}

//goland:noinspection GoUnusedParameter
func (r *userRepository) Save(ctx context.Context, user *do.User) (string, error) {
	return "", fmt.Errorf("need implement")
}

//goland:noinspection GoUnusedParameter
func (r *userRepository) GetUser(ctx context.Context, userID string) (*do.User, error) {
	return &do.User{
		ID:   "100",
		Name: "Username",
		Age:  40,
	}, nil
}

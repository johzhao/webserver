package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (r *userRepository) Save(ctx context.Context, user *do.User) (string, error) {
	// TODO: generate user ID when empty
	sql := "INSERT INTO `users` (`id`, `name`, `age`) VALUES (?, ?, ?)"
	_, err := r.execer.ExecContext(ctx, sql, user.ID, user.Name, user.Age)
	if err != nil {
		return "", errors.WithMessage(err, "save user failed")
	}

	return user.ID, nil
}

func (r *userRepository) GetUser(ctx context.Context, userID string) (*do.User, error) {
	sql := "SELECT * FROM `users` WHERE `id`=?"
	var result do.User
	if err := r.queryer.QueryRowxContext(ctx, sql, userID).Scan(&result); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("get user by id %s failed", userID))
	}

	return &result, nil
}

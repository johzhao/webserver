package service

import (
	"context"
	"go.uber.org/zap"
	"webserver/database"
	"webserver/errors"
	"webserver/model"
	"webserver/model/command"
	"webserver/model/do"
	"webserver/model/dto"
)

type UserService interface {
	CreateUser(ctx context.Context, cmd command.CreateUserCommand) (string, error)
	UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error

	GetUser(ctx context.Context, userID string) (*dto.User, error)
}

func NewUserService(logger *zap.Logger, db database.Database) UserService {
	return &userService{
		db:     db,
		logger: logger,
	}
}

type userService struct {
	db     database.Database
	logger *zap.Logger
}

func (s *userService) CreateUser(ctx context.Context, cmd command.CreateUserCommand) (string, error) {
	if err := cmd.Validation(); err != nil {
		return "", errors.Wrap(err, "request was invalid")
	}

	db, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = db.Rollback() }()

	userToCreate := do.User{} // XXX: create the user object from cmd
	_ = userToCreate

	return db.UserRepository().Save(ctx, &userToCreate)
}

func (s *userService) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error {
	if err := cmd.Validation(); err != nil {
		return errors.Wrap(err, "request was invalid")
	}

	db, err := s.db.BeginTransaction(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = db.Rollback() }()

	userToUpdate, err := db.UserRepository().GetUser(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// TODO: update the user by cmd

	if _, err := db.UserRepository().Save(ctx, userToUpdate); err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUser(ctx context.Context, userID string) (*dto.User, error) {
	if len(userID) == 0 {
		err := errors.ErrorBadRequest.New("invalid user id")
		return nil, errors.AddErrorContext(err, map[string]interface{}{
			"userID": userID,
		})
	}

	resultUser, err := s.db.UserRepository().GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return model.AssembleDTOUser(resultUser), nil
}

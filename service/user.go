package service

import (
	"context"
	"go.uber.org/zap"
	"webserver/database/repository"
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

func NewUserService(repository repository.Repository, logger *zap.Logger) UserService {
	return &userService{
		repository: repository,
		logger:     logger,
	}
}

type userService struct {
	repository repository.Repository
	logger     *zap.Logger
}

func (s userService) CreateUser(ctx context.Context, cmd command.CreateUserCommand) (string, error) {
	if err := cmd.Validation(); err != nil {
		return "", errors.Wrap(err, "request was invalid")
	}

	userToCreate := do.User{} // XXX: create the user object from cmd

	return s.repository.Save(ctx, &userToCreate)
}

func (s userService) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error {
	if err := cmd.Validation(); err != nil {
		return errors.Wrap(err, "request was invalid")
	}

	userToUpdate, err := s.repository.GetUser(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	// XXX: update the user by cmd

	//if _, err := s.repository.Save(ctx, userToUpdate); err != nil {
	//	return err
	//}
	_ = userToUpdate

	return nil
}

func (s userService) GetUser(ctx context.Context, userID string) (*dto.User, error) {
	if len(userID) == 0 {
		err := errors.ErrorBadRequest.New("invalid user id")
		return nil, errors.AddErrorContext(err, map[string]interface{}{
			"userID": userID,
		})
	}

	resultUser, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return model.AssembleDTOUser(resultUser), nil
}

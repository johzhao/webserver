package user

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"webserver/api/user"
	"webserver/api/user/command"
	"webserver/api/user/dto"
	"webserver/user/model"
)

func NewUserService(repository Repository, logger *zap.Logger) user.Service {
	return &Service{
		repository: repository,
		logger:     logger,
	}
}

type Service struct {
	repository Repository
	logger     *zap.Logger
}

func (s Service) CreateUser(ctx context.Context, cmd command.CreateUserCommand) (string, error) {
	if err := cmd.Validation(); err != nil {
		return "", err
	}

	userToCreate := model.User{} // XXX: create the user object from cmd

	return s.repository.Save(ctx, &userToCreate)
}

func (s Service) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error {
	if err := cmd.Validation(); err != nil {
		return err
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

func (s Service) GetUser(ctx context.Context, userID string) (*dto.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("invalid user id")
	}

	resultUser, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return assembleDTOUser(resultUser), nil
}

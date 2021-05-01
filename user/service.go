package user

import (
	"context"
	"go.uber.org/zap"
	"webserver/api/user"
	"webserver/api/user/command"
	"webserver/api/user/dto"
	"webserver/errors"
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
		return "", errors.Wrap(err, "request was invalid")
	}

	userToCreate := model.User{} // XXX: create the user object from cmd

	return s.repository.Save(ctx, &userToCreate)
}

func (s Service) UpdateUser(ctx context.Context, cmd command.UpdateUserCommand) error {
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

func (s Service) GetUser(ctx context.Context, userID string) (*dto.User, error) {
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

	return assembleDTOUser(resultUser), nil
}

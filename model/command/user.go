package command

import (
	"webserver/errors"
)

type CreateUserCommand struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c CreateUserCommand) Validation() error {
	if len(c.Name) == 0 {
		err := errors.ErrorBadRequest.New("invalid user name")
		return errors.AddErrorContext(err, map[string]interface{}{
			"parameters": c,
		})
	}
	if c.Age <= 0 || c.Age >= 120 {
		err := errors.ErrorBadRequest.New("invalid user age")
		return errors.AddErrorContext(err, map[string]interface{}{
			"parameters": c,
		})
	}

	return nil
}

type UpdateUserCommand struct {
	UserID string `path:"user_id"`
}

func (c UpdateUserCommand) Validation() error {
	if len(c.UserID) == 0 {
		return errors.ErrorBadRequest.New("invalid user id")
	}

	return nil
}

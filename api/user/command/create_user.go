package command

import (
	"webserver/errors"
)

type CreateUserCommand struct {
	Name string
	Age  int
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

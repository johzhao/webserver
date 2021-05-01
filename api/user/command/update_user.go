package command

import "webserver/errors"

type UpdateUserCommand struct {
	UserID string
}

func (c UpdateUserCommand) Validation() error {
	if len(c.UserID) == 0 {
		return errors.ErrorBadRequest.New("invalid user id")
	}

	return nil
}

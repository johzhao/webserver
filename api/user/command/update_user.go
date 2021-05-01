package command

type UpdateUserCommand struct {
	UserID string
}

func (c UpdateUserCommand) Validation() error {
	return nil
}

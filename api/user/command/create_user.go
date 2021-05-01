package command

import "fmt"

type CreateUserCommand struct {
	Name string
	Age  int
}

func (c CreateUserCommand) Validation() error {
	if len(c.Name) == 0 {
		return fmt.Errorf("invalid user name")
	}
	if c.Age <= 0 || c.Age >= 120 {
		return fmt.Errorf("user age must great than 0 and less than 120")
	}

	return nil
}

package user

import (
	"webserver/api/user/dto"
	"webserver/user/model"
)

func assembleDTOUser(user *model.User) *dto.User {
	result := dto.User{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}
	return &result
}

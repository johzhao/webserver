package model

import (
	"webserver/model/do"
	"webserver/model/dto"
)

func AssembleDTOUser(user *do.User) *dto.User {
	result := dto.User{
		ID:   user.ID,
		Name: user.Name,
		Age:  user.Age,
	}
	return &result
}

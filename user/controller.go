package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webserver/api/user"
	"webserver/api/user/command"
)

func NewUserController(service user.Service) Controller {
	return Controller{
		userService: service,
	}
}

type Controller struct {
	userService user.Service
}

func (c Controller) CreateUser(ctx *gin.Context) {
	cmd := command.CreateUserCommand{} // XXX: make the cmd
	userID, err := c.userService.CreateUser(ctx, cmd)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"user_id": userID,
			},
		})
	}
}

func (c Controller) UpdateUser(ctx *gin.Context) {
	userID := ctx.Params.ByName("user_id")
	cmd := command.UpdateUserCommand{
		UserID: userID,
	} // XXX: make the cmd
	err := c.userService.UpdateUser(ctx, cmd)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    nil,
		})
	}
}

func (c Controller) GetUser(ctx *gin.Context) {
	userID := ctx.Params.ByName("user_id")
	dtoUser, err := c.userService.GetUser(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    -1,
			"message": err.Error(),
			"data":    nil,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"user": dtoUser,
			},
		})
	}
}

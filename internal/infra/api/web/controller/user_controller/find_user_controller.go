package user_controller

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/rest_err"
	"fullcycle-auction_go/internal/usecase/user_usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase user_usecase.UserUseCaseInterface
}

func NewUserController(userUseCase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}

func (u *UserController) FindUserById(c *gin.Context) {
	userId := c.Param("userId")
	fmt.Println(userId)

	userData, err := u.userUseCase.FindUserById(context.Background(), userId)
	fmt.Println(userData)
	if err != nil {
		errRest := rest_err.ConvertError(err)
		c.JSON(errRest.Code, errRest)
		return
	}

	c.JSON(http.StatusOK, userData)
}

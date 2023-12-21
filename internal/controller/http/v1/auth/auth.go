package auth

import (
	"context"
	"errors"
	"meetup/internal/service/hash"
	"meetup/internal/service/response"
	"net/http"

	"github.com/gin-gonic/gin"

	"meetup/internal/auth"
	"meetup/internal/service/request"
)

type Controller struct {
	user User
	auth Auth
}

func NewController(user User, auth Auth) *Controller {
	return &Controller{user, auth}
}

func (ac Controller) SignIn(c *gin.Context) {
	var data auth.SignIn

	if er := request.BindFunc(c, &data); er != nil {
		response.RespondError(c, er)
		return
	}

	ctx := context.Background()

	userDetail, er := ac.user.GetUserByUsername(ctx, data.Username)
	if er != nil {
		response.RespondError(c, er)
		return
	}
	isValidPassword := hash.CheckPasswordHash(data.Password, *userDetail.Password)
	if !isValidPassword {
		c.JSON(http.StatusOK, response.Errors{
			Error:  errors.New("invalid password!").Error(),
			Status: false,
		})
		return
	}

	var generateTokenData auth.GenerateToken

	generateTokenData.Username = userDetail.Username
	generateTokenData.UserId = userDetail.Id
	generateTokenData.Role = userDetail.Role

	token, err := ac.auth.GenerateToken(ctx, generateTokenData)

	if err != nil {
		c.JSON(http.StatusOK, response.Errors{
			Error:  err.Error(),
			Status: false,
		})

		return
	}

	response.Respond(c, map[string]string{
		"token": token,
	})
}

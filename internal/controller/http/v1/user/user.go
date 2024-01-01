package user

import (
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/user"
	user_repo "meetup/internal/repository/postgres/user"
	"meetup/internal/service/hash"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	user User
}

func NewController(user User) *Controller {
	return &Controller{user}
}

func (uc Controller) CreateUser(c *gin.Context) {
	var data user_repo.CreateUserRequest

	if err := request.BindFunc(c, &data, "Username", "Password", "Role"); err != nil {
		response.RespondError(c, err)
		return
	}
	password, err := hash.HashPassword(*data.Password)
	if err != nil {
		response.RespondError(c, &pkg.Error{
			Err:    pkg.WrapError(err, "creating user hash password"),
			Status: http.StatusInternalServerError,
		})
		return
	}

	data.Password = &password
	detail, er := uc.user.UserCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetUserById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := uc.user.UserGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetUserList(c *gin.Context) {
	var filter user_repo.Filter
	fieldErrors := make([]pkg.FieldError, 0)

	limit, err := request.GetQuery(c, reflect.Int, "limit")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := limit.(*int); ok {
		filter.Limit = value
	}

	offset, err := request.GetQuery(c, reflect.Int, "offset")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := offset.(*int); ok {
		filter.Offset = value
	}

	role, err := request.GetQuery(c, reflect.String, "role")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := role.(*string); ok {
		filter.Role = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := uc.user.UserGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (uc Controller) UpdateUser(c *gin.Context) {
	var data user_repo.UpdateUserRequest
	if err := request.BindFunc(c, &data, "Username", "Role"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := uc.user.UserUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) DeleteUser(c *gin.Context) {

	Id := c.Param("id")

	er := uc.user.UserDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) ResetUserPassword(c *gin.Context) {
	var data user.UpdatePasswordRequest

	if err := request.BindFunc(c, &data, "Id", "NewPassword"); err != nil {
		response.RespondError(c, err)
		return
	}

	er := uc.user.UserUpdatePassword(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

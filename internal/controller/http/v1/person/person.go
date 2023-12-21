package person

import (
	"meetup/internal/pkg"
	person_repo "meetup/internal/repository/postgres/person"
	"meetup/internal/service/file"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	person Person
}

func NewController(person Person) *Controller {
	return &Controller{person}
}

func (uc Controller) CreatePerson(c *gin.Context) {
	var data person_repo.CreatePersonRequest

	if err := request.BindFunc(c, &data, "FullName"); err != nil {
		response.RespondError(c, err)
		return
	}

	avatarLink, err := file.NewService().Upload(c, data.Avatar, "avatar")
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.AvatarLink = avatarLink

	detail, er := uc.person.PersonCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetPersonById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := uc.person.PersonGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetPersonList(c *gin.Context) {
	var filter person_repo.Filter
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

	search, err := request.GetQuery(c, reflect.String, "search")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := search.(*string); ok {
		filter.Search = value
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := uc.person.PersonGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (uc Controller) UpdatePerson(c *gin.Context) {
	var data person_repo.UpdatePersonRequest
	if err := request.BindFunc(c, &data, "FullName"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := uc.person.PersonUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) DeletePerson(c *gin.Context) {

	Id := c.Param("id")

	er := uc.person.PersonDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

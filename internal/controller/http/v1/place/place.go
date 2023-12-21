package place

import (
	"meetup/internal/pkg"
	place_repo "meetup/internal/repository/postgres/place"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	place Place
}

func NewController(place Place) *Controller {
	return &Controller{place}
}

func (pc Controller) CreatePlace(c *gin.Context) {
	var data place_repo.CreatePlaceRequest

	if err := request.BindFunc(c, &data, "Name"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := pc.place.PlaceCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (pc Controller) GetPlaceById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := pc.place.PlaceGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (pc Controller) GetPlaceList(c *gin.Context) {
	var filter place_repo.Filter
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

	list, count, er := pc.place.PlaceGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (pc Controller) UpdatePlace(c *gin.Context) {
	var data place_repo.UpdatePlaceRequest
	if err := request.BindFunc(c, &data, "Name"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := pc.place.PlaceUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (pc Controller) DeletePlace(c *gin.Context) {

	Id := c.Param("id")

	er := pc.place.PlaceDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

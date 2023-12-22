package meeting_place_place

import (
	"meetup/internal/pkg"
	meeting_place_repo "meetup/internal/repository/postgres/meeting_place"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	meeting_place MeetingPlace
}

func NewController(meeting_place MeetingPlace) *Controller {
	return &Controller{meeting_place}
}

func (mc Controller) CreateMeetingPlace(c *gin.Context) {
	var data meeting_place_repo.CreateMeetingPlaceRequest

	if err := request.BindFunc(c, &data, "MeetingId", "PersonId", "PlaceId"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := mc.meeting_place.MeetingPlaceCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mc Controller) GetMeetingPlaceById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := mc.meeting_place.MeetingPlaceGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mc Controller) GetMeetingPlaceList(c *gin.Context) {
	var filter meeting_place_repo.Filter
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

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := mc.meeting_place.MeetingPlaceGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (mc Controller) GetMeetingPlaceListWithName(c *gin.Context) {
	meeting_id := c.Param("meeting_id")

	list, count, er := mc.meeting_place.MeetingPlaceList(c, meeting_id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (mc Controller) UpdateMeetingPlace(c *gin.Context) {
	var data meeting_place_repo.UpdateMeetingPlaceRequest
	if err := request.BindFunc(c, &data); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := mc.meeting_place.MeetingPlaceUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (mc Controller) DeleteMeetingPlace(c *gin.Context) {

	Id := c.Param("id")

	er := mc.meeting_place.MeetingPlaceDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

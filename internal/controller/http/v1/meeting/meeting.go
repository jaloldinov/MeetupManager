package meeting

import (
	"meetup/internal/pkg"
	"meetup/internal/pkg/config"
	meeting_repo "meetup/internal/repository/postgres/meeting"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	meeting Meeting
}

func NewController(meeting Meeting) *Controller {
	return &Controller{meeting}
}

func (mc Controller) CreateMeeting(c *gin.Context) {
	var data meeting_repo.CreateMeetingRequest

	if err := request.BindFunc(c, &data, "Title", "StartTime", "EndTime"); err != nil {
		response.RespondError(c, err)
		return
	}

	detail, er := mc.meeting.MeetingCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mc Controller) GetMeetingById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := mc.meeting.MeetingGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (mc Controller) GetMeetingList(c *gin.Context) {
	var filter meeting_repo.Filter
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

	lang, err := request.GetQuery(c, reflect.String, "lang")
	if err != nil {
		fieldErrors = append(fieldErrors, *err)
	} else if value, ok := lang.(*string); ok {
		filter.Lang = value
	}

	if lang == nil {
		filter.Lang = &config.GetConf().DefaultLang
	}

	if len(fieldErrors) > 0 {
		response.RespondError(c, &pkg.Error{
			Err:    errors.New("invalid query"),
			Fields: fieldErrors,
			Status: http.StatusBadRequest,
		})
		return
	}

	list, count, er := mc.meeting.MeetingGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (mc Controller) UpdateMeeting(c *gin.Context) {
	var data meeting_repo.UpdateMeetingRequest
	if err := request.BindFunc(c, &data, "Title", "MeetingTime"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := mc.meeting.MeetingUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (mc Controller) DeleteMeeting(c *gin.Context) {

	Id := c.Param("id")

	er := mc.meeting.MeetingDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

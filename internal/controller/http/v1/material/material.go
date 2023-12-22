package material

import (
	"meetup/internal/pkg"
	"meetup/internal/pkg/config"
	material_repo "meetup/internal/repository/postgres/material"
	"meetup/internal/service/file"
	"meetup/internal/service/request"
	"meetup/internal/service/response"

	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	material Material
}

func NewController(material Material) *Controller {
	return &Controller{material}
}

func (uc Controller) CreateMaterial(c *gin.Context) {
	var data material_repo.CreateMaterialRequest
	data.File = make(map[string]string)

	if err := request.BindFunc(c, &data, "Title"); err != nil {
		response.RespondError(c, err)
		return
	}

	uz_file, err := file.NewService().Upload(c, data.UzFile, "materials/uz_file")
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.File["uz"] = uz_file

	en_file, err := file.NewService().Upload(c, data.EnFile, "materials/en_file")
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.File["en"] = en_file

	ru_file, err := file.NewService().Upload(c, data.RuFile, "materials/ru_file")
	if err != nil {
		response.RespondError(c, err)
		return
	}
	data.File["ru"] = ru_file

	detail, er := uc.material.MaterialCreate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetMaterialById(c *gin.Context) {
	idParam := c.Param("id")

	detail, er := uc.material.MaterialGetById(c, idParam)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, detail)
}

func (uc Controller) GetMaterialList(c *gin.Context) {
	var filter material_repo.Filter
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

	list, count, er := uc.material.MaterialGetAll(c, filter)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.Respond(c, map[string]interface{}{
		"results": list,
		"count":   count,
	})
}

func (uc Controller) UpdateMaterial(c *gin.Context) {
	var data material_repo.UpdateMaterialRequest
	if err := request.BindFunc(c, &data, "FullName"); err != nil {
		response.RespondError(c, err)
		return
	}

	data.Id = c.Param("id")

	er := uc.material.MaterialUpdate(c, data)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

func (uc Controller) DeleteMaterial(c *gin.Context) {

	Id := c.Param("id")

	er := uc.material.MaterialDelete(c, Id)
	if er != nil {
		response.RespondError(c, er)
		return
	}

	response.RespondNoData(c)
}

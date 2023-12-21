package response

import (
	"meetup/internal/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Respond(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"status":  true,
		"data":    data,
	})
}

func RespondNoData(c *gin.Context) {
	c.JSON(http.StatusOK, StatusOk{
		Message: "ok",
		Status:  true,
	})
}

func RespondError(c *gin.Context, err *pkg.Error) {
	c.JSON(err.Status, Errors{
		Error:  err.Err.Error(),
		Fields: err.Fields,
		Status: false,
	})
}

type Errors struct {
	Error  string           `json:"error"`
	Fields []pkg.FieldError `json:"fields"`
	Status bool             `json:"status"`
}

type StatusOk struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

type StatusBadRequest struct {
	Error string `json:"error"`
}

type CustomError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

package request

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"

	"meetup/internal/pkg"
)

func BindFunc(c *gin.Context, data interface{}, requiredFields ...string) *pkg.Error {
	r := c.ShouldBind(data)
	if r != nil {
		return &pkg.Error{
			Err:    r,
			Status: http.StatusInternalServerError,
		}
	}

	err := validateStruct(data, requiredFields...)

	return err
}

func validateStruct(s interface{}, requiredFields ...string) *pkg.Error {
	structVal := reflect.Value{}
	if reflect.Indirect(reflect.ValueOf(s)).Kind() == reflect.Struct {
		structVal = reflect.Indirect(reflect.ValueOf(s))
	} else {
		return &pkg.Error{
			Err:    errors.New("input object should be struct"),
			Status: http.StatusBadRequest,
		}
	}

	errFields := make([]pkg.FieldError, 0)

	structType := reflect.Indirect(reflect.ValueOf(s)).Type()
	fieldNum := structVal.NumField()

	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()
		if !isSet {
			log.Print(isSet, fieldName, reflect.ValueOf(field))
			for _, f := range requiredFields {

				if f == fieldName {
					errFields = append(errFields, pkg.FieldError{
						Err:   errors.New("field is required!"),
						Field: fieldName,
					})
				}
			}
		}
	}
	if len(errFields) > 0 {
		return &pkg.Error{
			Err:    errors.New("required fields"),
			Fields: errFields,
			Status: http.StatusBadRequest,
		}
	}
	return nil
}

// func GetContext(c *gin.Context) (context.Context, *pkg.Error) {
// 	headerToken := c.Request.Header["Authorization"]
// 	var token string
// 	if len(headerToken) > 0{
// 		splitToken := strings.Split(headerToken[0], " ")
// 		if len(splitToken) == 2 {
// 			 token = splitToken[1]
// 		}
// 	}

// 	tokenData, err := // need parse token and get role, user_id

// 	ctx := context.WithValue(context.Background(), "role", "ADMIN")
// 	ctx = context.WithValue(ctx, "userId", "uuid from token")

// 	return ctx
// }

func GetTokenFromHeader(c *gin.Context) (string, error) {
	headerToken := c.Request.Header["Authorization"]
	if len(headerToken) > 0 {
		splitToken := strings.Split(headerToken[0], " ")
		if len(splitToken) == 2 {
			return splitToken[1], nil
		}

		return "", errors.New("Invalid token!")
	}

	return "", errors.New("Token is not found!")
}

func GetQuery(c *gin.Context, queryType reflect.Kind, query string) (interface{}, *pkg.FieldError) {
	switch queryType {
	case reflect.String:
		if len(c.Request.URL.Query()[query]) > 0 {
			return &c.Request.URL.Query()[query][0], nil
		}
	case reflect.Int:
		if len(c.Request.URL.Query()[query]) > 0 {
			queryInt, err := strconv.Atoi(c.Request.URL.Query()[query][0])
			if err != nil {
				return nil, &pkg.FieldError{
					Err:   errors.New(fmt.Sprintf("%s must be number", query)),
					Field: query,
				}
			}

			return &queryInt, nil
		}
	case reflect.Bool:
		if len(c.Request.URL.Query()[query]) > 0 {
			queryBool, err := strconv.ParseBool(c.Request.URL.Query()[query][0])
			if err != nil {
				return nil, &pkg.FieldError{
					Err:   errors.New(fmt.Sprintf("%s must be boolean", query)),
					Field: query,
				}
			}

			return &queryBool, nil
		}
	case 10:
		if len(c.Request.URL.Query()[query]) > 0 {
			q := c.Request.URL.Query()[query][0]
			err := checkSQLInjection(q)
			if err != nil {
				return nil, &pkg.FieldError{
					Err:   err,
					Field: query,
				}
			}
			return &q, nil
		}
	}

	return nil, nil
}

func checkSQLInjection(text string) error {
	containsOr := strings.Contains(text, " or ")
	if containsOr {
		return errors.New("sql injection detected!")
	}
	containsDropTable := strings.Contains(text, "drop table")
	if containsDropTable {
		return errors.New("sql injection detected!")
	}
	containsDropDatabase := strings.Contains(text, "drop database")
	if containsDropDatabase {
		return errors.New("sql injection detected!")
	}
	containsSelectTable := strings.Contains(text, "select ")
	if containsSelectTable {
		return errors.New("sql injection detected!")
	}
	containsComment := strings.Contains(text, " -- ")
	if containsComment {
		return errors.New("sql injection detected!")
	}
	match, _ := regexp.MatchString("([^a-zA-Z0-9]+)=([^a-zA-Z0-9]+)", text)
	if match {
		return errors.New("sql injection detected!")
	}
	return nil
}

func GetParam(c *gin.Context, paramType reflect.Kind, param string) (interface{}, *pkg.FieldError) {
	switch paramType {
	case reflect.String:
		if len(c.Param(param)) > 0 {
			return c.Param(param), nil
		}
	case reflect.Int:
		if len(c.Param(param)) > 0 {
			_, err := strconv.Atoi(c.Param(param))
			if err != nil {
				return nil, &pkg.FieldError{
					Err:   errors.New(fmt.Sprintf("%s must be number", param)),
					Field: param,
				}
			}
			return c.Param(param), nil
		}
	case 11:
		//uuid check
		if len(c.Param(param)) > 0 {
			_, err := uuid.Parse(c.Param(param))
			if err != nil {
				return nil, &pkg.FieldError{
					Err:   errors.New(fmt.Sprintf("%s must be number", param)),
					Field: param,
				}
			}
			return c.Param(param), nil
		}
	}

	return nil, nil
}

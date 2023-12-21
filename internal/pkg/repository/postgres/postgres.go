package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"meetup/internal/pkg"
	"net/http"
	"reflect"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type CtxData struct {
	UserId string
	Role   string
	Lang   string
}

type Database struct {
	*bun.DB
	DefaultLang   string
	ServerBaseUrl string
}

func New(DBUsername, DBPassword, DBPort, DBName string) *Database {
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", DBUsername, DBPassword, DBPort, DBName)
	fmt.Println(dsn)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	return &Database{
		DB: db,
	}
}

func (d Database) DeleteRow(ctx context.Context, table, role, id string) *pkg.Error {
	dataCtx, er := d.CheckCtx(ctx)
	if er != nil {
		return er
	}

	if dataCtx.Role == "" {
		return &pkg.Error{
			Err:    errors.New("role in context is required"),
			Status: http.StatusInternalServerError,
		}
	}

	if dataCtx.Role != role {
		return &pkg.Error{
			Err:    errors.New(fmt.Sprintf("you have not permission to delete from table: %s", table)),
			Status: http.StatusInternalServerError,
		}
	}

	_, err := d.NewUpdate().
		Table(table).
		Where("deleted_at is null AND id = ?", id).
		Set("deleted_at = ?", time.Now()).
		Set("deleted_by = ?", dataCtx.UserId).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	return nil
}

func (d Database) ValidateStruct(s interface{}, requiredFields ...string) *pkg.Error {
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

func (d Database) CheckCtx(ctx context.Context) (CtxData, *pkg.Error) {
	fieldErrors := make([]pkg.FieldError, 0)
	userId, ok := ctx.Value("user_id").(string)
	if !ok {
		fieldErrors = append(fieldErrors, pkg.FieldError{
			Err:   errors.New("missing field in ctx"),
			Field: "user_id",
		})
	}
	ctxRole, ok := ctx.Value("role").(string)
	if !ok {
		fieldErrors = append(fieldErrors, pkg.FieldError{
			Err:   errors.New("missing field in ctx"),
			Field: "role",
		})
	}

	ctxLang, ok := ctx.Value("lang").(string)
	if !ok {
		fieldErrors = append(fieldErrors, pkg.FieldError{
			Err:   errors.New("missing field in ctx"),
			Field: "role",
		})
	}

	if len(fieldErrors) > 0 {
		return CtxData{}, &pkg.Error{
			Err:    errors.New("missing fields in context"),
			Fields: fieldErrors,
			Status: http.StatusInternalServerError,
		}
	}

	return CtxData{
		UserId: userId,
		Role:   ctxRole,
		Lang:   ctxLang,
	}, nil
}

func (d Database) GetLang(ctx context.Context) string {
	return d.DefaultLang
}

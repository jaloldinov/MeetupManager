package person

import (
	"context"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/person"
)

type Person interface {
	PersonCreate(ctx context.Context, data person.CreatePersonRequest) (person.CreatePersonResponse, *pkg.Error)
	PersonGetById(ctx context.Context, id string) (person.GetPersonResponse, *pkg.Error)
	PersonGetAll(ctx context.Context, filter person.Filter) ([]person.GetPersonListResponse, int, *pkg.Error)
	PersonUpdate(ctx context.Context, data person.UpdatePersonRequest) *pkg.Error
	PersonDelete(ctx context.Context, id string) *pkg.Error
}

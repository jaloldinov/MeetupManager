package person

import (
	"context"
	"errors"
	"meetup/internal/entity"
	"meetup/internal/pkg"
	"meetup/internal/pkg/repository/postgres"
	"net/http"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) PersonCreate(ctx context.Context, request CreatePersonRequest) (CreatePersonResponse, *pkg.Error) {
	var detail CreatePersonResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreatePersonResponse{}, er
	}

	detail.FullName = *request.FullName
	detail.Avatar = request.AvatarLink
	detail.Info = request.Info
	detail.Status = request.Status
	detail.CreatedBy = &dataCtx.UserId
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Returning("id").Exec(ctx)

	if err != nil {
		return CreatePersonResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating person"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) PersonGetById(ctx context.Context, id string) (GetPersonResponse, *pkg.Error) {
	var person GetPersonResponse

	err := r.NewSelect().Model(&person).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetPersonResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting person get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return person, nil
}

func (r Repository) PersonGetAll(ctx context.Context, filter Filter) ([]GetPersonListResponse, int, *pkg.Error) {
	var list []GetPersonListResponse

	q := r.NewSelect().Model(&list)
	q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
		query.Where("deleted_at is null")
		return query
	})
	if filter.Limit != nil {
		q.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		q.Offset(*filter.Offset)
	}

	if filter.Search != nil {
		q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
			query.Where("lower(full_name) like lower(?)", "%"+*filter.Search+"%")
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting person list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) PersonUpdate(ctx context.Context, request UpdatePersonRequest) *pkg.Error {
	var detail entity.Person
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating person get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.FullName != nil {
		detail.FullName = *request.FullName
	}

	if request.AvatarLink != nil {
		detail.Avatar = *request.AvatarLink
	}

	if request.Info != nil {
		detail.Info = request.Info
	}
	if request.Status != nil {
		detail.Status = *request.Status
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating person"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) PersonDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	result, err := r.NewUpdate().
		Table("persons").
		Where("deleted_at is null AND id = ?", id).
		Set("deleted_at = ?, deleted_by = ?", time.Now(), dataCtx.UserId).
		Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    errors.New("delete row error, updating"),
			Status: http.StatusInternalServerError,
		}
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return &pkg.Error{
			Err:    errors.New("no matching ID found"),
			Status: http.StatusNotFound,
		}
	}

	return nil
}

package place

import (
	"context"
	"errors"
	"meetup/internal/entity"
	"meetup/internal/pkg"
	"meetup/internal/pkg/repository/postgres"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Repository struct {
	*postgres.Database
}

func NewRepository(postgresDB *postgres.Database) *Repository {
	return &Repository{postgresDB}
}

func (r Repository) PlaceCreate(ctx context.Context, request CreatePlaceRequest) (CreatePlaceResponse, *pkg.Error) {
	var detail CreatePlaceResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreatePlaceResponse{}, er
	}

	detail.Id = uuid.NewString()
	detail.Name = *request.Name
	detail.Code = request.Code
	detail.CreatedBy = &dataCtx.UserId
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreatePlaceResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating place"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) PlaceGetById(ctx context.Context, id string) (GetPlaceResponse, *pkg.Error) {
	var place GetPlaceResponse

	err := r.NewSelect().Model(&place).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetPlaceResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting place get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return place, nil
}

func (r Repository) PlaceGetAll(ctx context.Context, filter Filter) ([]GetPlaceListResponse, int, *pkg.Error) {
	var list []GetPlaceListResponse

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
			query.Where("lower(name) like lower(?)", "%"+*filter.Search+"%")
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting place list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) PlaceUpdate(ctx context.Context, request UpdatePlaceRequest) *pkg.Error {
	var detail entity.Place
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating place get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Name != nil {
		detail.Name = *request.Name
	}

	if request.Code != nil {
		detail.Code = *request.Code
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating place"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) PlaceDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	result, err := r.NewUpdate().
		Table("places").
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

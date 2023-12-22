package material

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

func (r Repository) MaterialCreate(ctx context.Context, request CreateMaterialRequest) (CreateMaterialResponse, *pkg.Error) {
	var detail CreateMaterialResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateMaterialResponse{}, er
	}

	detail.Index = request.Index
	detail.Title = *request.Title
	detail.Content = request.Content
	detail.File = request.File
	detail.CreatedBy = &dataCtx.UserId
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Returning("id").Exec(ctx)

	if err != nil {
		return CreateMaterialResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating material"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) MaterialGetById(ctx context.Context, id string) (GetMaterialResponse, *pkg.Error) {
	var material GetMaterialResponse

	err := r.NewSelect().Model(&material).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetMaterialResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting material get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return material, nil
}

func (r Repository) MaterialGetAll(ctx context.Context, filter Filter) ([]GetMaterialListResponse, int, *pkg.Error) {
	var list []GetMaterialListResponse

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
			Err:    pkg.WrapError(err, "selecting material list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) MaterialUpdate(ctx context.Context, request UpdateMaterialRequest) *pkg.Error {
	var detail entity.Material
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating material get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Index != nil {
		detail.Index = *request.Index
	}
	if request.Title != nil {
		detail.Title = *request.Title
	}
	if request.Content != nil {
		detail.Content = *request.Content
	}
	if request.File != nil {
		detail.File = request.File
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating material"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) MaterialDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	result, err := r.NewUpdate().
		Table("materials").
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

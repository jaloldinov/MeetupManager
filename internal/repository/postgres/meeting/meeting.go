package meeting

import (
	"context"
	"errors"
	"fmt"
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

func (r Repository) MeetingCreate(ctx context.Context, request CreateMeetingRequest) (CreateMeetingResponse, *pkg.Error) {
	var detail CreateMeetingResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateMeetingResponse{}, er
	}

	detail.Title = request.Title
	detail.Description = request.Description
	detail.StartTime = *request.StartTime
	detail.EndTime = *request.EndTime
	detail.CreatedBy = &dataCtx.UserId
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Exec(ctx)

	if err != nil {
		return CreateMeetingResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating meeting"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) MeetingGetById(ctx context.Context, id string) (GetMeetingResponse, *pkg.Error) {
	var meeting GetMeetingResponse

	err := r.NewSelect().Model(&meeting).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetMeetingResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting meeting get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return meeting, nil
}

func (r Repository) MeetingGetAll(ctx context.Context, filter Filter) ([]GetMeetingListResponse, int, *pkg.Error) {
	var list []GetMeetingListResponse

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
			query.Where(fmt.Sprintf(" lower(title->>'%s') similar to lower('%s')", *filter.Lang, "%"+*filter.Search+"%"))
			return query
		})
	}

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting meeting list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) MeetingUpdate(ctx context.Context, request UpdateMeetingRequest) *pkg.Error {
	var detail entity.Meeting
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating meeting get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.Title != nil {
		detail.Title = request.Title
	}

	if request.Description != nil {
		detail.Description = request.Description
	}

	if request.StartTime != nil {
		detail.StartTime = request.StartTime
	}
	if request.EndTime != nil {
		detail.EndTime = request.EndTime
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating meeting"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) MeetingDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	result, err := r.NewUpdate().
		Table("meetings").
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

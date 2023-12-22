package meeting_place

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

func (r Repository) MeetingPlaceCreate(ctx context.Context, request CreateMeetingPlaceRequest) (CreateMeetingPlaceResponse, *pkg.Error) {
	var detail CreateMeetingPlaceResponse
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return CreateMeetingPlaceResponse{}, er
	}

	detail.MeetingId = request.MeetingId
	detail.PersonId = request.PersonId
	detail.PlaceId = request.PlaceId
	detail.CreatedBy = &dataCtx.UserId
	detail.CreatedAt = time.Now()

	_, err := r.NewInsert().Model(&detail).Returning("id").Exec(ctx)

	if err != nil {
		return CreateMeetingPlaceResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "creating meeting_place"),
			Status: http.StatusInternalServerError,
		}
	}

	return detail, nil
}

func (r Repository) MeetingPlaceGetById(ctx context.Context, id string) (GetMeetingPlaceResponse, *pkg.Error) {
	var meeting_place GetMeetingPlaceResponse

	err := r.NewSelect().Model(&meeting_place).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return GetMeetingPlaceResponse{}, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting meeting_place get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	return meeting_place, nil
}

func (r Repository) MeetingPlaceGetAll(ctx context.Context, filter Filter) ([]GetMeetingPlaceListResponse, int, *pkg.Error) {
	var list []GetMeetingPlaceListResponse

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

	// if filter.Role != nil {
	// 	q.WhereGroup(" and ", func(query *bun.SelectQuery) *bun.SelectQuery {
	// 		query.Where("role = ? ", filter.Role)
	// 		return query
	// 	})
	// }

	q.Order("created_at desc")

	count, err := q.ScanAndCount(ctx)
	if err != nil {
		return nil, 0, &pkg.Error{
			Err:    pkg.WrapError(err, "selecting meeting_place list"),
			Status: http.StatusInternalServerError,
		}
	}
	return list, count, nil
}

func (r Repository) MeetingPlaceList(ctx context.Context, meetingID string) ([]MeetingPlaceResponse, int, *pkg.Error) {
	var list []MeetingPlaceResponse

	query := fmt.Sprintf(`
		SELECT
			mp.person_id,
			p.full_name AS person_name,
			pl.id as place_id,
			pl.name AS place_name
		FROM
			meeting_places mp
		INNER JOIN persons p ON mp.person_id = p.id
		INNER JOIN places pl ON mp.place_id = pl.id
		WHERE 
			mp.meeting_id = %v
	`, meetingID)

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, 0, &pkg.Error{Status: 500, Err: err}
	}
	defer rows.Close()

	for rows.Next() {
		var meetingPlace MeetingPlaceResponse
		if err := rows.Scan(&meetingPlace.PersonId, &meetingPlace.PersonName, &meetingPlace.PlaceId, &meetingPlace.PlaceName); err != nil {
			return nil, 0, &pkg.Error{Status: 500, Err: err}
		}
		list = append(list, meetingPlace)
	}

	var count int
	err = r.DB.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(1) FROM meeting_places WHERE meeting_id = %v", meetingID)).Scan(&count)
	if err != nil {
		return nil, 0, &pkg.Error{Status: 500, Err: err}
	}

	return list, count, nil
}

func (r Repository) MeetingPlaceUpdate(ctx context.Context, request UpdateMeetingPlaceRequest) *pkg.Error {
	var detail entity.MeetingPlace
	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}

	err := r.NewSelect().Model(&detail).Where("id = ?", &request.Id).Scan(ctx)
	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating meeting_place get by id"),
			Status: http.StatusInternalServerError,
		}
	}

	if request.MeetingId != nil {
		detail.MeetingId = *request.MeetingId
	}
	if request.PersonId != nil {
		detail.PersonId = *request.PersonId
	}
	if request.PlaceId != nil {
		detail.PlaceId = *request.PlaceId
	}

	date := time.Now()
	detail.UpdatedAt = &date
	detail.UpdatedBy = &dataCtx.UserId

	_, err = r.NewUpdate().Model(&detail).Where("id = ?", detail.Id).Exec(ctx)

	if err != nil {
		return &pkg.Error{
			Err:    pkg.WrapError(err, "updating meeting_place"),
			Status: http.StatusInternalServerError,
		}
	}
	return nil
}

func (r Repository) MeetingPlaceDelete(ctx context.Context, id string) *pkg.Error {

	dataCtx, er := r.CheckCtx(ctx)
	if er != nil {
		return er
	}
	result, err := r.NewUpdate().
		Table("meeting_places").
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

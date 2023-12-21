package meeting

import (
	"context"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/meeting"
)

type Meeting interface {
	MeetingCreate(ctx context.Context, data meeting.CreateMeetingRequest) (meeting.CreateMeetingResponse, *pkg.Error)
	MeetingGetById(ctx context.Context, id string) (meeting.GetMeetingResponse, *pkg.Error)
	MeetingGetAll(ctx context.Context, filter meeting.Filter) ([]meeting.GetMeetingListResponse, int, *pkg.Error)
	MeetingUpdate(ctx context.Context, data meeting.UpdateMeetingRequest) *pkg.Error
	MeetingDelete(ctx context.Context, id string) *pkg.Error
}

package meeting_place

import (
	"context"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/meeting_place"
)

type MeetingPlace interface {
	MeetingPlaceCreate(ctx context.Context, data meeting_place.CreateMeetingPlaceRequest) (meeting_place.CreateMeetingPlaceResponse, *pkg.Error)
	MeetingPlaceGetById(ctx context.Context, id string) (meeting_place.GetMeetingPlaceResponse, *pkg.Error)
	MeetingPlaceGetAll(ctx context.Context, filter meeting_place.Filter) ([]meeting_place.GetMeetingPlaceListResponse, int, *pkg.Error)
	MeetingPlaceUpdate(ctx context.Context, data meeting_place.UpdateMeetingPlaceRequest) *pkg.Error
	MeetingPlaceDelete(ctx context.Context, id string) *pkg.Error
	MeetingPlaceList(ctx context.Context, meetingID string) ([]meeting_place.MeetingPlaceResponse, int, *pkg.Error)

	//for websocket
	GetDetailByPlaceID(ctx context.Context, meetingID string) (*meeting_place.MeetingDetail, *pkg.Error)
}

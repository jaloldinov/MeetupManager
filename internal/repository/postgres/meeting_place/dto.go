package meeting_place

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
}

type CreateMeetingPlaceRequest struct {
	MeetingId int `json:"meeting_id" form:"meeting_id"`
	PersonId  int `json:"person_id" form:"person_id"`
	PlaceId   int `json:"place_id" form:"place_id"`
}

type CreateMeetingPlaceResponse struct {
	bun.BaseModel `bun:"table:meeting_places"`

	Id        *int      `json:"id" bun:"id,pk"`
	MeetingId int       `json:"meeting_id" bun:"meeting_id"`
	PersonId  int       `json:"person_id" bun:"person_id"`
	PlaceId   int       `json:"place_id" bun:"place_id"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy *string   `json:"-"`
}

type GetMeetingPlaceResponse struct {
	bun.BaseModel `bun:"table:meeting_places"`

	Id        int `json:"id" bun:"id,pk"`
	MeetingId int `json:"meeting_id" bun:"meeting_id"`
	PersonId  int `json:"person_id" bun:"person_id"`
	PlaceId   int `json:"place_id" bun:"place_id"`
}

type GetMeetingPlaceListResponse struct {
	bun.BaseModel `bun:"table:meeting_places"`

	Id        int `json:"id" bun:"id,pk"`
	MeetingId int `json:"meeting_id" bun:"meeting_id"`
	PersonId  int `json:"person_id" bun:"person_id"`
	PlaceId   int `json:"place_id" bun:"place_id"`
}

type MeetingPlaceResponse struct {
	PersonId   int    `json:"person_id"`
	PersonName string `json:"full_name"`
	PlaceId    int    `json:"place_id"`
	PlaceName  string `json:"place_name"`
}

type UpdateMeetingPlaceRequest struct {
	Id        string `json:"id" bun:"id,pk"`
	MeetingId *int   `json:"meeting_id" form:"meeting_id"`
	PersonId  *int   `json:"person_id" form:"person_id"`
	PlaceId   *int   `json:"place_id" form:"place_id"`
}

type DetailMeetingPlaceResponse struct {
	bun.BaseModel `bun:"table:meeting_places"`
}

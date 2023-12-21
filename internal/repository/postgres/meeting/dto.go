package meeting

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Search *string
	Lang   *string
}

type CreateMeetingRequest struct {
	Id          string             `json:"id" form:"id,pk"`
	Title       *map[string]string `json:"title" form:"title"`
	Description map[string]string  `json:"description" form:"description"`
	MeetingTime time.Time          `json:"meeting_time" form:"meeting_time"`
}

type CreateMeetingResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	MeetingTime time.Time          `json:"meeting_time" bun:"meeting_time"`
	CreatedAt   time.Time          `json:"-" bun:"created_at"`
	CreatedBy   *string            `json:"-"`
}

type GetMeetingResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	MeetingTime time.Time          `json:"meeting_time" bun:"meeting_time"`
	CreatedAt   time.Time          `json:"-" bun:"created_at"`
}

type GetMeetingListResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	MeetingTime time.Time          `json:"meeting_time" bun:"meeting_time"`
}

type UpdateMeetingRequest struct {
	Id          string             `json:"id" form:"id,pk"`
	Title       *map[string]string `json:"title" form:"title"`
	Description map[string]string  `json:"description" form:"description"`
	MeetingTime *time.Time         `json:"meeting_time" form:"meeting_time"`
}

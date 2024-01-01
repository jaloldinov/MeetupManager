package meeting

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit     *int
	Offset    *int
	Search    *string
	Lang      *string
	StartTime *time.Time
	EndTime   *time.Time
}

type CreateMeetingRequest struct {
	Title       *map[string]string `json:"title" form:"title"`
	Description map[string]string  `json:"description" form:"description"`
	StartTime   *time.Time         `json:"start_time" bun:"start_time"`
	EndTime     *time.Time         `json:"end_time" bun:"end_time"`
}

type CreateMeetingResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          *int               `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	StartTime   time.Time          `json:"start_time" bun:"start_time"`
	EndTime     time.Time          `json:"end_time" bun:"end_time"`
	CreatedAt   time.Time          `json:"-" bun:"created_at"`
	CreatedBy   *string            `json:"-"`
}

type GetMeetingResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	StartTime   time.Time          `json:"start_time" bun:"start_time"`
	EndTime     time.Time          `json:"end_time" bun:"end_time"`
	CreatedAt   time.Time          `json:"-" bun:"created_at"`
}

type GetMeetingListResponse struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	StartTime   time.Time          `json:"start_time" bun:"start_time"`
	EndTime     time.Time          `json:"end_time" bun:"end_time"`
}

type UpdateMeetingRequest struct {
	Id          string             `json:"id" form:"id,pk"`
	Title       *map[string]string `json:"title" form:"title"`
	Description map[string]string  `json:"description" form:"description"`
	StartTime   *time.Time         `json:"start_time" bun:"start_time"`
	EndTime     *time.Time         `json:"end_time" bun:"end_time"`
}

package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Meeting struct {
	bun.BaseModel `bun:"table:meetings"`

	Id          string             `json:"id" bun:"id,pk"`
	Title       *map[string]string `json:"title" bun:"title"`
	Description map[string]string  `json:"description" bun:"description"`
	MeetingTime time.Time          `json:"meeting_time" bun:"meeting_time"`
	CreatedAt   time.Time          `json:"created_at" bun:"created_at"`
	CreatedBy   *string            `json:"created_by" bun:"created_by"`
	UpdatedAt   *time.Time         `json:"updated_at" bun:"updated_at"`
	UpdatedBy   *string            `json:"updated_by" bun:"updated_by"`
	DeletedAt   *time.Time         `json:"deleted_at" bun:"deleted_at"`
	DeletedBy   *string            `json:"deleted_by" bun:"deleted_by"`
}

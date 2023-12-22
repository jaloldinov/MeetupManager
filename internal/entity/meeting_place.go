package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type MeetingPlace struct {
	bun.BaseModel `bun:"table:meeting_places"`

	Id        int        `json:"id" bun:"id,pk"`
	MeetingId int        `json:"meeting_id" bun:"meeting_id"`
	PersonId  int        `json:"person_id" bun:"person_id"`
	PlaceId   int        `json:"place_id" bun:"place_id"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Place struct {
	bun.BaseModel `bun:"table:users"`

	Id        string     `json:"id" bun:"id,pk"`
	Name      string     `json:"name" bun:"name"`
	Code      string     `json:"code" bun:"code"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

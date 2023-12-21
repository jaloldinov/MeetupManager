package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Person struct {
	bun.BaseModel `bun:"table:persons"`

	Id        string            `json:"id" bun:"id,pk"`
	FullName  string            `json:"full_name" bun:"full_name"`
	Avatar    string            `json:"avatar" bun:"avatar"`
	Info      map[string]string `json:"info" bun:"info"`
	Status    bool              `json:"status" bun:"status"`
	CreatedAt time.Time         `json:"created_at" bun:"created_at"`
	CreatedBy *string           `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time        `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string           `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time        `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string           `json:"deleted_by" bun:"deleted_by"`
}

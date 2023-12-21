package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Id        string     `json:"id" bun:"id,pk"`
	Username  string     `json:"username" bun:"username"`
	Password  string     `json:"password,omitempty" bun:"password"`
	Role      string     `json:"role" bun:"role"`
	CreatedAt time.Time  `json:"created_at" bun:"created_at"`
	CreatedBy *string    `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string    `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string    `json:"deleted_by" bun:"deleted_by"`
}

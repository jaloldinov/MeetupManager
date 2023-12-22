package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type Material struct {
	bun.BaseModel `bun:"table:materials"`

	Id        string            `json:"id" bun:"id,pk"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	File      map[string]string `json:"file" bun:"file"`
	Index     int               `json:"index" bun:"index"`
	CreatedAt time.Time         `json:"created_at" bun:"created_at"`
	CreatedBy *string           `json:"created_by" bun:"created_by"`
	UpdatedAt *time.Time        `json:"updated_at" bun:"updated_at"`
	UpdatedBy *string           `json:"updated_by" bun:"updated_by"`
	DeletedAt *time.Time        `json:"deleted_at" bun:"deleted_at"`
	DeletedBy *string           `json:"deleted_by" bun:"deleted_by"`
}

package place

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Search *string
}

type CreatePlaceRequest struct {
	Name *string `json:"name" form:"name"`
	Code string  `json:"code" form:"code"`
}

type CreatePlaceResponse struct {
	bun.BaseModel `bun:"table:places"`

	Id        *int      `json:"id" bun:"id,pk"`
	Name      string    `json:"name" bun:"name"`
	Code      string    `json:"code" bun:"code"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy *string   `json:"-"`
}

type GetPlaceResponse struct {
	bun.BaseModel `bun:"table:places"`

	Id        string    `json:"id" bun:"id,pk"`
	Name      string    `json:"name" bun:"name"`
	Code      string    `json:"code" bun:"code"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
}

type GetPlaceListResponse struct {
	bun.BaseModel `bun:"table:places"`

	Id   string `json:"id" bun:"id,pk"`
	Name string `json:"name" bun:"name"`
	Code string `json:"code" bun:"code"`
}

type UpdatePlaceRequest struct {
	Id   string  `json:"id" bun:"id"`
	Name *string `json:"name" form:"name"`
	Code *string `json:"code" form:"code"`
}

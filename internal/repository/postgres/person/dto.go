package person

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Search *string
}

type CreatePersonRequest struct {
	FullName   *string               `json:"full_name" form:"full_name"`
	Avatar     *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink string                `json:"-" form:"-"`
	Info       map[string]string     `json:"info" form:"info"`
	Status     bool                  `json:"status" form:"status"`
}

type CreatePersonResponse struct {
	bun.BaseModel `bun:"table:persons"`

	Id        *int              `json:"id" bun:"id"`
	FullName  string            `json:"full_name" bun:"full_name"`
	Avatar    string            `json:"avatar" bun:"avatar"`
	Info      map[string]string `json:"info" bun:"info"`
	Status    bool              `json:"status" bun:"status"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
	CreatedBy *string           `json:"-"`
}

type GetPersonResponse struct {
	bun.BaseModel `bun:"table:persons"`

	Id        int               `json:"id" bun:"id"`
	FullName  string            `json:"full_name" bun:"full_name"`
	Avatar    string            `json:"avatar" bun:"avatar"`
	Info      map[string]string `json:"info" bun:"info"`
	Status    bool              `json:"status" bun:"status"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
}

type GetPersonListResponse struct {
	bun.BaseModel `bun:"table:persons"`

	Id       int    `json:"id" bun:"id"`
	FullName string `json:"full_name" bun:"full_name"`
	Avatar   string `json:"avatar" bun:"avatar"`
	Status   bool   `json:"status" bun:"status"`
}

type UpdatePersonRequest struct {
	Id         string                `json:"id" bun:"id"`
	FullName   *string               `json:"full_name" form:"full_name"`
	Avatar     *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink *string               `json:"-" form:"-"`
	Info       map[string]string     `json:"info" form:"info"`
	Status     *bool                 `json:"status" form:"status"`
}

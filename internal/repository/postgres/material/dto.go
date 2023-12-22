package material

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Search *string
	Lang   *string
}

type CreateMaterialRequest struct {
	Index   int                   `json:"index" form:"index"`
	Title   *map[string]string    `json:"title" form:"title"`
	Content map[string]string     `json:"content" form:"content"`
	UzFile  *multipart.FileHeader `json:"-" form:"uz_file"`
	EnFile  *multipart.FileHeader `json:"-" form:"en_file"`
	RuFile  *multipart.FileHeader `json:"-" form:"ru_file"`
	File    map[string]string     `json:"-" form:"-"`
}

type CreateMaterialResponse struct {
	bun.BaseModel `bun:"table:materials"`

	Id        *int              `json:"id" bun:"id"`
	Index     int               `json:"index" bun:"index"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	File      map[string]string `json:"file" bun:"file"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
	CreatedBy *string           `json:"-"`
}

type GetMaterialResponse struct {
	bun.BaseModel `bun:"table:materials"`

	Id        int               `json:"id" bun:"id"`
	Index     int               `json:"index" bun:"index"`
	Title     map[string]string `json:"title" bun:"title"`
	Content   map[string]string `json:"content" bun:"content"`
	File      map[string]string `json:"file" bun:"file"`
	CreatedAt time.Time         `json:"-" bun:"created_at"`
}

type GetMaterialListResponse struct {
	bun.BaseModel `bun:"table:materials"`

	Id    int               `json:"id" bun:"id"`
	Index int               `json:"index" bun:"index"`
	Title map[string]string `json:"title" bun:"title"`
	// Content map[string]string  `json:"content" bun:"content"`
	File map[string]string `json:"file" bun:"file"`
}

type UpdateMaterialRequest struct {
	Id      string                `json:"id" bun:"id"`
	Index   *int                  `json:"index" form:"index"`
	Title   *map[string]string    `json:"title" form:"title"`
	Content *map[string]string    `json:"content" form:"content"`
	UzFile  *multipart.FileHeader `json:"-" form:"uz_file"`
	EnFile  *multipart.FileHeader `json:"-" form:"en_file"`
	RuFile  *multipart.FileHeader `json:"-" form:"ru_file"`
	File    map[string]string     `json:"-" bun:"-"`
}

package user

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Role   *string
}

type CreateUserRequest struct {
	Username *string `json:"username" form:"username"`
	Password *string `json:"password" form:"password"`
	Role     *string `json:"role" form:"role"`
}

type CreateUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id        string    `json:"id" bun:"id"`
	Username  string    `json:"username" bun:"username"`
	Password  *string   `json:"-" bun:"password"`
	Role      string    `json:"role" bun:"role"`
	CreatedAt time.Time `json:"-" bun:"created_at"`
	CreatedBy *string   `json:"-"`
}

type GetUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id       string `json:"id" bun:"id"`
	Username string `json:"username" bun:"username"`
	Role     string `json:"role" bun:"role"`
}

type GetUserListResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id       string `json:"id" bun:"id"`
	Username string `json:"username" bun:"username"`
	Role     string `json:"role" bun:"role"`
}

type UpdateUserRequest struct {
	Id       string  `json:"id" bun:"id"`
	Username *string `json:"username" form:"username"`
	Role     *string `json:"role" form:"role"`
}

type UpdatePasswordRequest struct {
	Id          *string `json:"id" form:"id"`
	NewPassword *string `json:"new_password" form:"new_password"`
}

type DetailUserResponse struct {
	bun.BaseModel `bun:"table:users"`

	Id       string  `json:"id" bun:"id"`
	Username string  `json:"username" bun:"username"`
	Password *string `json:"-" bun:"password"`
	Role     string  `json:"-" bun:"role"`
}

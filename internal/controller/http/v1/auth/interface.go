package auth

import (
	"context"
	"meetup/internal/auth"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/user"
)

type Auth interface {
	GenerateToken(ctx context.Context, data auth.GenerateToken) (string, error)
}

type User interface {
	GetUserByUsername(ctx context.Context, username string) (user.DetailUserResponse, *pkg.Error)
}

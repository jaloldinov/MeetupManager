package material

import (
	"context"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/material"
)

type Material interface {
	MaterialCreate(ctx context.Context, data material.CreateMaterialRequest) (material.CreateMaterialResponse, *pkg.Error)
	MaterialGetById(ctx context.Context, id string) (material.GetMaterialResponse, *pkg.Error)
	MaterialGetAll(ctx context.Context, filter material.Filter) ([]material.GetMaterialListResponse, int, *pkg.Error)
	MaterialUpdate(ctx context.Context, data material.UpdateMaterialRequest) *pkg.Error
	MaterialDelete(ctx context.Context, id string) *pkg.Error
}

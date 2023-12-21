package place

import (
	"context"
	"meetup/internal/pkg"
	"meetup/internal/repository/postgres/place"
)

type Place interface {
	PlaceCreate(ctx context.Context, data place.CreatePlaceRequest) (place.CreatePlaceResponse, *pkg.Error)
	PlaceGetById(ctx context.Context, id string) (place.GetPlaceResponse, *pkg.Error)
	PlaceGetAll(ctx context.Context, filter place.Filter) ([]place.GetPlaceListResponse, int, *pkg.Error)
	PlaceUpdate(ctx context.Context, data place.UpdatePlaceRequest) *pkg.Error
	PlaceDelete(ctx context.Context, id string) *pkg.Error
}

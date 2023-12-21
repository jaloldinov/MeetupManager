package locale

import (
	"context"
	"encoding/json"
	"meetup/internal/pkg"
	"net/http"
	"os"
	"strings"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) ReadJSONFile(ctx context.Context, fileLink string) (map[string]string, error) {
	content, err := os.ReadFile(fileLink)
	if err != nil {
		return nil, err
	}

	// Now let's un marshall the data into `payload`
	var payload map[string]string
	err = json.Unmarshal(content, &payload)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (s Service) GetLocale(ctx context.Context, locales string) (map[string]map[string]string, *pkg.Error) {
	splitLocales := strings.Split(locales, ",")
	localeJSON := make(map[string]map[string]string, len(splitLocales))

	for _, l := range splitLocales {
		jsonData, err := s.ReadJSONFile(ctx, "./locale/"+l+".json")
		if err != nil {
			return nil, &pkg.Error{
				Err:    pkg.WrapError(err, "get locale"),
				Status: http.StatusInternalServerError,
			}
		}

		localeJSON[l] = jsonData
	}

	return localeJSON, nil
}

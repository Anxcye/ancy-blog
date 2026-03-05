// File: tmdb.go
// Purpose: TMDB service wrapper for handler layer.
// Module: backend/internal/service, business layer.
// Related: service/tmdb package, handler/admin.

package service

import (
	"encoding/json"
	"fmt"

	"github.com/anxcye/ancy-blog/backend/internal/service/tmdb"
)

type TMDBService struct {
	integrationService *IntegrationService
}

type TMDBMetadata struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"posterPath"`
	ReleaseDate string  `json:"releaseDate"`
	VoteAverage float64 `json:"voteAverage"`
}

func NewTMDBService(integrationService *IntegrationService) *TMDBService {
	return &TMDBService{
		integrationService: integrationService,
	}
}

func (s *TMDBService) GetMetadata(mediaType, id string) (*TMDBMetadata, error) {
	provider, ok := s.integrationService.GetIntegrationProviderForRuntime("tmdb")
	if !ok || !provider.Enabled {
		return nil, fmt.Errorf("TMDB provider not configured or disabled")
	}

	var config map[string]interface{}
	if err := json.Unmarshal(provider.ConfigJSON, &config); err != nil {
		return nil, fmt.Errorf("invalid TMDB config: %w", err)
	}

	apiKey := ""
	if key, ok := config["api_key"].(string); ok {
		apiKey = key
	}

	client := tmdb.New(apiKey)

	if mediaType == "movie" {
		movie, err := client.GetMovie(id)
		if err != nil {
			return nil, err
		}
		return &TMDBMetadata{
			ID:          movie.ID,
			Title:       movie.Title,
			Overview:    movie.Overview,
			PosterPath:  movie.PosterPath,
			ReleaseDate: movie.ReleaseDate,
			VoteAverage: movie.VoteAverage,
		}, nil
	}

	tv, err := client.GetTV(id)
	if err != nil {
		return nil, err
	}
	return &TMDBMetadata{
		ID:          tv.ID,
		Title:       tv.Name,
		Overview:    tv.Overview,
		PosterPath:  tv.PosterPath,
		ReleaseDate: tv.FirstAirDate,
		VoteAverage: tv.VoteAverage,
	}, nil
}

// File: tmdb.go
// Purpose: TMDB API client for fetching movie/TV metadata.
// Module: backend/internal/service/tmdb, integration layer.
// Related: handler/admin/integration.go, config.

package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Service struct {
	apiKey     string
	httpClient *http.Client
}

type MovieDetail struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
	VoteAverage float64 `json:"vote_average"`
}

type TVDetail struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	PosterPath   string  `json:"poster_path"`
	FirstAirDate string  `json:"first_air_date"`
	VoteAverage  float64 `json:"vote_average"`
}

func New(apiKey string) *Service {
	return &Service{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *Service) GetMovie(id string) (*MovieDetail, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%s?language=zh-CN", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
	}

	var movie MovieDetail
	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Service) GetTV(id string) (*TVDetail, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}
	url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%s?language=zh-CN", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("TMDB API returned status %d", resp.StatusCode)
	}

	var tv TVDetail
	if err := json.NewDecoder(resp.Body).Decode(&tv); err != nil {
		return nil, err
	}
	return &tv, nil
}

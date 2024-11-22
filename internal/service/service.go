package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/iurikman/songs/internal/config"
	"github.com/iurikman/songs/internal/models"
)

type Service struct {
	db db
}

func NewService(db db) *Service {
	return &Service{
		db: db,
	}
}

type db interface {
	CreateSong(ctx context.Context, song models.Song) (*models.Song, error)
}

func (s *Service) CreateSong(ctx context.Context, song models.Song) (*models.Song, error) {
	cfg := config.NewConfig()

	u, err := url.Parse(cfg.APIURL)
	if err != nil {
		return nil, fmt.Errorf("rl.Parse(cfg.APIURL): %w", err)
	}

	query := u.Query()
	query.Set("groupe", song.Group)
	query.Set("song", song.Name)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil) err: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != http.StatusOK {
		switch {
		case resp.StatusCode == http.StatusBadRequest:
			return nil, models.ErrBadRequest
		case resp.StatusCode == http.StatusInternalServerError:
			return nil, fmt.Errorf("internal server error: %w", err)
		}
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("resp.Body.ReadAll() err: %w", err)
	}

	var songDetail models.SongDetail

	if err := json.Unmarshal(respBody, &songDetail); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(respBody, &songDetail): %w", err)
	}

	createdSong, err := s.db.CreateSong(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("s.db.createSong(ctx, song) err: %w", err)
	}

	createdSong.ReleaseDate = songDetail.ReleaseDate
	createdSong.Text = songDetail.Text
	createdSong.Link = songDetail.Link

	return createdSong, nil
}

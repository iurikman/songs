package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/iurikman/songs/internal/models"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
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
	GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error)
	GetText(ctx context.Context, id uuid.UUID) (*string, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error)
}

func (s *Service) CreateSong(ctx context.Context, song models.Song) (*models.Song, error) {
	songDetail, err := getDetails(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("getDetails(ctx, song) err: %w", err)
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	createdSong, err := s.db.CreateSong(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("s.db.createSong(ctx, song) err: %w", err)
	}

	return createdSong, nil
}

func (s *Service) GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error) {
	songs, err := s.db.GetSongs(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.db.createSong(ctx, params) err: %w", err)
	}

	return songs, nil
}

func (s *Service) GetText(ctx context.Context, id uuid.UUID, verse int) (*string, error) {
	text, err := s.db.GetText(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("s.db.getText(ctx, id, params) err: %w", err)
	}

	splittedText := strings.Split(*text, "\n\n")

	if verse < 1 || verse > len(splittedText) {
		return nil, models.ErrVerseIsNotValid
	}

	textOfVerse := splittedText[verse-1]

	return &textOfVerse, nil
}

func (s *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	if err := s.db.DeleteSong(ctx, id); err != nil {
		return fmt.Errorf("s.db.deleteSong(ctx, id) err: %w", err)
	}

	return nil
}

func (s *Service) UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error) {
	updatedSong, err := s.db.UpdateSong(ctx, id, song)
	if err != nil {
		return nil, fmt.Errorf("s.db.UpdateSong(ctx, id, song) err: %w", err)
	}

	return updatedSong, nil
}

func getDetails(ctx context.Context, song models.Song) (*models.SongDetail, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	u, err := url.Parse(os.Getenv("API_URL") + os.Getenv("API_PORT"))
	if err != nil {
		return nil, fmt.Errorf("url.Parse(os.Getenv(\"API_URL\")) err: %w", err)
	}

	query := u.Query()
	query.Set("group", song.Group)
	query.Set("song", song.Name)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil) err: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

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

	return &songDetail, nil
}

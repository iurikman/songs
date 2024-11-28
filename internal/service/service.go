package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/songs/internal/models"
	log "github.com/sirupsen/logrus"
)

type Service struct {
	db          db
	songDetails songDetailsClient
}

func NewService(db db, songDetailsServer songDetailsClient) *Service {
	log.Debug("Initializing new service")

	return &Service{
		db:          db,
		songDetails: songDetailsServer,
	}
}

type songDetailsClient interface {
	Get(ctx context.Context, song models.Song) (*models.Song, error)
}

type db interface {
	CreateSong(ctx context.Context, song models.Song) (*models.Song, error)
	GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error)
	GetText(ctx context.Context, id uuid.UUID, verse int) (*string, error)
	DeleteSong(ctx context.Context, id uuid.UUID) error
	UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error)
}

func (s *Service) CreateSong(ctx context.Context, song models.Song) (*models.Song, error) {
	log.Debugf("Creating song, getting song details from songdetails: %+v", song)

	songWithDetails, err := s.songDetails.Get(ctx, song)
	if err != nil {
		return nil, fmt.Errorf("getDetails(ctx, song) err: %w", err)
	}

	log.Debug("Details retrieved and assigned to song, creating new song")

	createdSong, err := s.db.CreateSong(ctx, *songWithDetails)
	if err != nil {
		return nil, fmt.Errorf("s.db.createSong(ctx, song) err: %w", err)
	}

	log.Infof("Song successfully created: %+v", createdSong)

	return createdSong, nil
}

func (s *Service) GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error) {
	log.Debugf("Retrieving songs with params: %+v", params)

	songs, err := s.db.GetSongs(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("s.db.createSong(ctx, params) err: %w", err)
	}

	log.Infof("Successfully retrieved %d songs", len(songs))

	return songs, nil
}

func (s *Service) GetText(ctx context.Context, id uuid.UUID, verse int) (*string, error) {
	log.Debugf("Retrieving text for song ID: %s, verse: %d", id, verse)

	textOfVerse, err := s.db.GetText(ctx, id, verse)
	if err != nil {
		return nil, fmt.Errorf("s.db.getText(ctx, id, verse) err: %w", err)
	}

	return textOfVerse, nil
}

func (s *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	log.Debugf("Deleting song with ID: %s", id)

	if err := s.db.DeleteSong(ctx, id); err != nil {
		return fmt.Errorf("s.db.deleteSong(ctx, id) err: %w", err)
	}

	log.Infof("Song successfully deleted with ID: %s", id)

	return nil
}

func (s *Service) UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error) {
	log.Debugf("Updating song with ID: %s", id)

	updatedSong, err := s.db.UpdateSong(ctx, id, song)
	if err != nil {
		return nil, fmt.Errorf("s.db.UpdateSong(ctx, id, song) err: %w", err)
	}

	log.Infof("Song successfully updated: %+v", updatedSong)

	return updatedSong, nil
}

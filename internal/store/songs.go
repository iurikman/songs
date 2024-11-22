package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iurikman/songs/internal/models"
	"github.com/jackc/pgx/v5"
)

func (p *Postgres) CreateSong(ctx context.Context, song models.Song) (*models.Song, error) {
	query := `	INSERT INTO songs (id, release_date, name, music_group, text, link)
				VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING id, release_date, name, group, text, link
				`

	createdSong := new(models.Song)

	err := p.db.QueryRow(
		ctx,
		query,
		uuid.New(),
		song.ReleaseDate,
		song.Name,
		song.Group,
		song.Text,
		song.Link,
	).Scan(
		&createdSong.ID,
		&createdSong.ReleaseDate,
		&createdSong.Name,
		&createdSong.Group,
		&createdSong.Text,
		&createdSong.Link,
	)
	if err != nil {
		return nil, fmt.Errorf("creating song err: %w", err)
	}

	return createdSong, nil
}

func (p *Postgres) UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) error {
	query := `	UPDATE songs SET release_date = $2, name = $3, music_group = $4, text = $5, link = $6
	            WHERE id = $1
				RETURNING id, release_date, name, group, text, link
				`

	updatedSong := new(models.Song)

	err := p.db.QueryRow(
		ctx,
		query,
		id,
		song.ReleaseDate,
		song.Name,
		song.Group,
		song.Text,
		song.Link,
	).Scan(
		&updatedSong.ID,
		&updatedSong.ReleaseDate,
		&updatedSong.Name,
		&updatedSong.Group,
		&updatedSong.Text,
		&updatedSong.Link,
	)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return models.ErrSongNotFound
	case err != nil:
		return fmt.Errorf("updating song err: %w", err)
	}

	return nil
}

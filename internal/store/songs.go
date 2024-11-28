package store

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/iurikman/songs/internal/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) CreateSong(ctx context.Context, song models.Song) (*models.Song, error) {
	query := `	INSERT INTO songs (id, release_date, name, music_group, text, link, deleted)
				VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id, release_date, name, music_group, text, link, deleted
				`

	createdSong := new(models.Song)

	err := p.db.QueryRow(
		ctx,
		query,
		song.ID,
		song.ReleaseDate,
		song.Name,
		song.Group,
		song.Text,
		song.Link,
		song.Deleted,
	).Scan(
		&createdSong.ID,
		&createdSong.ReleaseDate,
		&createdSong.Name,
		&createdSong.Group,
		&createdSong.Text,
		&createdSong.Link,
		&createdSong.Deleted,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
			return nil, models.ErrDuplicateSong
		case err != nil:
			return nil, fmt.Errorf("creating song err: %w", err)
		}
	}

	return createdSong, nil
}

func (p *Postgres) GetSongs(ctx context.Context, params models.Params) ([]*models.Song, error) {
	songs := make([]*models.Song, 0, 1)

	query := `
				SELECT id, release_date, name, music_group, text, link
				FROM songs
				WHERE deleted=false
			`

	if params.Filter != "" {
		query += fmt.Sprintf(" and name LIKE '%%%s%%'", params.Filter)
	}

	if params.Sorting != "" {
		query += " ORDER BY " + params.Sorting
		if params.Descending {
			query += " DESC"
		}
	}

	query += fmt.Sprintf(" OFFSET %d LIMIT %d", params.Offset, params.Limit)

	rows, err := p.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("getting songs err: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		song := new(models.Song)

		err := rows.Scan(
			&song.ID,
			&song.ReleaseDate,
			&song.Name,
			&song.Group,
			&song.Text,
			&song.Link,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning song err: %w", err)
		}

		songs = append(songs, song)
	}

	return songs, nil
}

func (p *Postgres) GetText(ctx context.Context, id uuid.UUID, verse int) (*string, error) {
	text := ""

	query := `
				SELECT text
				FROM songs
				WHERE id = $1 and deleted=false
			`

	err := p.db.QueryRow(ctx, query, id).Scan(&text)
	if err != nil {
		return nil, fmt.Errorf("p.db.QueryRow(ctx, query, id).Scan(&text) err: %w", err)
	}

	splittedText := strings.Split(text, "\n\n")

	if verse < 1 || verse > len(splittedText) {
		return nil, models.ErrVerseIsNotValid
	}

	textOfVerse := splittedText[verse-1]

	return &textOfVerse, nil
}

func (p *Postgres) DeleteSong(ctx context.Context, id uuid.UUID) error {
	query := `
				UPDATE songs SET deleted = true WHERE id = $1 and deleted = false
			`

	result, err := p.db.Exec(ctx, query, id)

	switch {
	case result.RowsAffected() == 0:
		return models.ErrSongNotFound
	case err != nil:
		return fmt.Errorf("deleting song error: %w", err)
	}

	return nil
}

func (p *Postgres) UpdateSong(ctx context.Context, id uuid.UUID, song models.Song) (*models.Song, error) {
	query := `	UPDATE songs SET release_date = $2, name = $3, music_group = $4, text = $5, link = $6
	            WHERE id = $1
				RETURNING id, release_date, name, music_group, text, link
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
		return nil, models.ErrSongNotFound
	case err != nil:
		return nil, fmt.Errorf("updating song err: %w", err)
	}

	return updatedSong, nil
}

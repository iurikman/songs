package songdetails

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iurikman/songs/internal/models"
)

type SongDetails struct {
	host string
}

func NewSongDetails(host string) *SongDetails {
	return &SongDetails{
		host: host,
	}
}

func (s *SongDetails) Get(ctx context.Context, song models.Song) (*models.Song, error) {
	reqURLstring := s.host + "/info" + "?song=" + song.Name + "&group=" + song.Group

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURLstring, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext(ctx, \"GET\", reqURLstring, nil) err: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.Do(req) err: %w", err)
	}
	defer resp.Body.Close()

	var songDetails models.SongDetails

	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		return nil, fmt.Errorf("json.Decode() err: %w", err)
	}

	songWithDetails := &models.Song{
		ID:          song.ID,
		ReleaseDate: songDetails.ReleaseDate,
		Text:        songDetails.Text,
		Link:        songDetails.Link,
		Name:        song.Name,
		Group:       song.Group,
	}

	return songWithDetails, nil
}

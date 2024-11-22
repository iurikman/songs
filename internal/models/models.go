package models

import "github.com/google/uuid"

type Song struct {
	ID          uuid.UUID `json:"id"`
	ReleaseDate string    `json:"releaseDate"`
	Name        string    `json:"name"`
	Group       string    `json:"musicGroup"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

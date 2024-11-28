package models

import "github.com/google/uuid"

type Song struct {
	ID          uuid.UUID `json:"id"`
	ReleaseDate string    `json:"releaseDate"`
	Name        string    `json:"name"`
	Group       string    `json:"musicGroup"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	Deleted     bool      `json:"deleted"`
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type Params struct {
	Offset     int    `schema:"offset"`
	Limit      int    `schema:"limit"`
	Sorting    string `schema:"sorting"`
	Descending bool   `schema:"descending"`
	Filter     string `schema:"filter"`
}

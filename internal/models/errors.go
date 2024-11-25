package models

import "errors"

var (
	ErrSongNotFound    = errors.New("song not found")
	ErrVerseIsNotValid = errors.New("verse is not valid")
	ErrDuplicateSong   = errors.New("duplicate song")
)

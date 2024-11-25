package models

import "errors"

var (
	ErrSongNotFound    = errors.New("song not found")
	ErrBadRequest      = errors.New("bad request")
	ErrVerseIsNotValid = errors.New("verse is not valid")
)

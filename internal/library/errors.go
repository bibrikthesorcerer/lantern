package library

import "errors"

var (
	ErrTrackNotFound = errors.New("track not found")
	ErrAlbumNotFound = errors.New("album not found")
)

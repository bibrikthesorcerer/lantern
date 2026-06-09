package library

import (
	"time"

	"go.senan.xyz/taglib"
)

type AlbumCover struct {
	taglib.ImageDesc
	Data []byte // Raw picture data.
}

type Track struct {
	ID       uint16      `json:"id"`
	Path     string      `json:"path"`
	Title    string      `json:"title"`
	Artist   string      `json:"artist"`
	Album    string      `json:"album"`
	Track    uint16      `json:"track"`
	Year     uint16      `json:"-"`
	Filename string      `json:"filename"`
	ModTime  time.Time   `json:"modtime"`
	Size     int64       `json:"size"`
	Cover    *AlbumCover `json:"-"`
}

type Album struct {
	ID          uint16  `json:"id"`
	Year        uint16  `json:"year"`
	Title       string  `json:"title"`
	AlbumArtist string  `json:"artist"`
	TotalTracks uint16  `json:"total_tracks"`
	Tracks      []Track `json:"tracks"`
}

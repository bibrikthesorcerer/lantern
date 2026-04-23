package main

import (
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/dhowden/tag"
)

type MusicServer struct {
	library map[uint16]Track
	albums  map[string]Album
	dir     string
}

type Track struct {
	ID       uint16       `json:"id"`
	Path     string       `json:"path"`
	Title    string       `json:"title"`
	Artist   string       `json:"artist"`
	Album    string       `json:"album"`
	Track    uint16       `json:"track"`
	Year     uint16       `json:"-"`
	Filename string       `json:"filename"`
	ModTime  time.Time    `json:"modtime"`
	Cover    *tag.Picture `json:"-"`
}

type Album struct {
	ID          uint16  `json:"id"`
	Year        uint16  `json:"year"`
	Title       string  `json:"title"`
	AlbumArtist string  `json:"artist"`
	TotalTracks uint16  `json:"total_tracks"`
	Tracks      []Track `json:"tracks"`
}

func NewServer(dir string) (*MusicServer, error) {
	library, albums, err := scanLibrary(dir)
	if err != nil {
		return nil, err
	}
	return &MusicServer{library: library, albums: albums, dir: dir}, nil
}

func (s *MusicServer) ListenAndServe() error {
	log.Infof("Launching server on :8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal("Server shutdown: %s", err)
	return err
}

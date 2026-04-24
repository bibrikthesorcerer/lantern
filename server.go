package main

import (
	"fmt"
	"net/http"
	"time"

	clog "github.com/charmbracelet/log"
	"github.com/dhowden/tag"
)

type MusicServer struct {
	library map[uint16]Track
	albums  map[string]Album
	conf    *Config
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

func NewServer(c *Config) (*MusicServer, error) {
	library, albums, err := scanLibrary(c.Dir)
	if err != nil {
		return nil, fmt.Errorf("library scan: %w", err)
	}
	return &MusicServer{library: library, albums: albums, conf: c}, nil
}

func (s *MusicServer) ListenAndServe() {
	addrString := fmt.Sprintf(":%d", s.conf.Port)
	clog.Infof("Launching server on %s", addrString)
	err := http.ListenAndServe(addrString, nil)
	clog.Fatalf("Server shutdown: %s", err)
}

package main

import (
	"fmt"
	"net/http"

	"github.com/bibrikthesorcerer/lantern/internal/library"
	"github.com/bibrikthesorcerer/lantern/internal/middleware"

	clog "github.com/charmbracelet/log"
)

type MusicServer struct {
	repo *library.LibraryRepository
	conf *Config
}

func NewServer(c *Config) (*MusicServer, error) {
	repo, err := library.NewTrackRepository(c.DBPath)
	if err != nil {
		return nil, fmt.Errorf("library setup: %w", err)
	}
	return &MusicServer{repo: repo, conf: c}, nil
}

func (s *MusicServer) RunServer() error {
	router := SetUpRouting(s)

	addrString := fmt.Sprintf(":%d", s.conf.Port)
	srv := &http.Server{
		Addr:    addrString,
		Handler: middleware.Logging(router),
	}

	clog.Infof("Launching server on %s", addrString)
	return fmt.Errorf("server shutdown: %w", srv.ListenAndServe())
}

package web

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/bibrikthesorcerer/lantern/internal/config"
	"github.com/bibrikthesorcerer/lantern/internal/library"
	"github.com/bibrikthesorcerer/lantern/internal/web/middleware"

	clog "github.com/charmbracelet/log"
)

type MusicServer struct {
	Repo        *library.LibraryRepository
	Conf        *config.Config
	StaticFiles fs.FS
}

func NewServer(c *config.Config, staticFiles fs.FS) (*MusicServer, error) {
	repo, err := library.NewTrackRepository(c.DBPath)
	if err != nil {
		return nil, fmt.Errorf("library setup: %w", err)
	}
	return &MusicServer{Repo: repo, Conf: c, StaticFiles: staticFiles}, nil
}

func (s *MusicServer) RunServer() error {
	router := SetUpRouting(s)

	addrString := fmt.Sprintf(":%d", s.Conf.Port)
	srv := &http.Server{
		Addr:    addrString,
		Handler: middleware.Logging(router),
	}

	clog.Infof("Launching server on %s", addrString)
	return fmt.Errorf("server shutdown: %w", srv.ListenAndServe())
}

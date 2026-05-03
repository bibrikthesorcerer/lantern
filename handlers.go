package main

import (
	"cmp"
	"embed"
	"encoding/json"
	"html/template"
	"maps"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/charmbracelet/log"
	clog "github.com/charmbracelet/log"
)

//go:embed static/*
var staticFiles embed.FS

func SetUpRouting(s *MusicServer) *http.ServeMux {
	router := http.NewServeMux()

	fs := http.FileServer(http.FS(staticFiles))
	router.Handle("GET /static/", fs)
	router.HandleFunc("GET /{$}", s.handleHome)
	router.HandleFunc("GET /tracks", s.handleTrackList)
	router.HandleFunc("GET /albums", s.handleAlbumList)
	router.HandleFunc("GET /api/stream/{id}", s.handleStream)
	router.HandleFunc("GET /api/tracks", s.handleTracks)
	router.HandleFunc("GET /api/cover/{id}", s.handleCover)
	router.HandleFunc("GET /api/albums", s.handleAlbums)
	router.HandleFunc("GET /api/tracks/{id}/download", s.handleTrackDownload)
	router.HandleFunc("GET /api/search", s.handleSearch)

	return router
}

func (s *MusicServer) handleHome(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"static/base.html",
		"static/nav.html",
		"static/player.html",
		"static/home.html",
	}
	tmpl, err := template.ParseFS(staticFiles, files...)
	if err != nil {
		clog.Errorf("couldn't parse templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		clog.Errorf("couldn't execute templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MusicServer) handleStream(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Errorf("handleStream: %s", err)
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}
	track, ok := s.library[uint16(id)]

	if !ok {
		log.Infof("404 Response for track #%d", id)
		http.NotFound(w, r)
		return
	}

	f, err := os.Open(track.Path)
	if err != nil {
		log.Errorf("handleStream: %s", err)
		http.Error(w, "cannot open file", http.StatusInternalServerError)
	}
	defer f.Close()
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, track.Filename, track.ModTime, f)
}

func (s *MusicServer) handleTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tracks := slices.SortedFunc(maps.Values(s.library), func(a, b Track) int {
		return cmp.Compare(a.ID, b.ID)
	})
	err := json.NewEncoder(w).Encode(tracks)
	if err != nil {
		log.Errorf("handleTracks: %s", err)
		http.Error(w, "cannot serialize Tracks objects", http.StatusInternalServerError)
	}
}

func (s *MusicServer) handleTrackList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"static/base.html",
		"static/nav.html",
		"static/player.html",
		"static/tracks_page.html",
	}
	tmpl, err := template.ParseFS(staticFiles, files...)
	if err != nil {
		clog.Errorf("couldn't parse templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		clog.Errorf("couldn't execute templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MusicServer) handleCover(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Errorf("handleStream: %s", err)
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}

	track, ok := s.library[uint16(id)]
	if !ok {
		http.NotFound(w, r)
		return
	}
	if track.Cover == nil {
		//TODO: serve a default placeholder image
		http.ServeFile(w, r, "static/placeholder.png")
		return
	}
	w.Header().Set("Content-Type", track.Cover.MIMEType)
	w.Write(track.Cover.Data)
}

func (s *MusicServer) handleAlbums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	albums := slices.SortedFunc(maps.Values(s.albums), func(a, b Album) int {
		return cmp.Compare(a.ID, b.ID)
	})
	for _, a := range albums {
		slices.SortFunc(a.Tracks, func(a, b Track) int { return cmp.Compare(a.Track, b.Track) })
	}
	err := json.NewEncoder(w).Encode(albums)
	if err != nil {
		log.Errorf("handleAlbums: %s", err)
		http.Error(w, "cannot serialize Tracks objects", http.StatusInternalServerError)
	}
}

func (s *MusicServer) handleAlbumList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"static/base.html",
		"static/nav.html",
		"static/player.html",
		"static/albums_page.html",
	}
	tmpl, err := template.ParseFS(staticFiles, files...)
	if err != nil {
		clog.Errorf("couldn't parse templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		clog.Errorf("couldn't execute templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *MusicServer) handleTrackDownload(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Errorf("handleStream: %s", err)
		http.Error(w, "ID must be an integer", http.StatusBadRequest)
		return
	}
	track := s.library[uint16(id)]
	w.Header().Set("Content-Disposition", "attachment; filename="+track.Path)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeFile(w, r, track.Path)
}

func (s *MusicServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	tracks := s.search(query)

	err := json.NewEncoder(w).Encode(tracks)
	if err != nil {
		log.Errorf("handleSearch: %s", err)
		http.Error(w, "cannot serialize Tracks objects", http.StatusInternalServerError)
	}
}

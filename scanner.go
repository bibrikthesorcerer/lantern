package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	clog "github.com/charmbracelet/log"
	"github.com/dhowden/tag"
)

var audioTypes = []string{"mp3", "wav", "ogg"} //TODO: flac ParseTrack error
var trackIndex uint16 = 0
var albumIndex uint16 = 0

func isSupportedAudio(name string) bool {
	_, ext, found := strings.Cut(name, ".")
	if !found {
		return false
	}

	return slices.Contains(audioTypes, ext)
}

func scanLibrary(dir string) (map[uint16]Track, map[string]Album, error) {
	library := make(map[uint16]Track)
	albums := make(map[string]Album)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !isSupportedAudio(d.Name()) {
			return nil
		}

		track, err := parseTrack(path, d)
		if err != nil {
			log.Printf("skipping %s: %v", path, err)
			return nil
		}

		library[track.ID] = track

		al, ok := albums[track.Album]
		if !ok {
			clog.Debugf("indexing album: %s, found in track %s", track.Album, track.Path)
			albums[track.Album] = Album{
				ID:          albumIndex,
				Title:       track.Album,
				Year:        track.Year,
				AlbumArtist: track.Artist,
				Tracks:      []Track{track},
			}
			albumIndex++
			return nil
		}
		al.Tracks = append(al.Tracks, track)
		albums[track.Album] = al
		return nil
	})

	return library, albums, err
}

func parseTrack(path string, entry os.DirEntry) (Track, error) {
	tr, err := entry.Info()
	if err != nil {
		log.Errorf("parseTrack: %s", err)
		return Track{}, err
	}

	// read metadata
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("parseTrack: %s", err)
		return Track{}, err
	}
	defer f.Close()

	m, err := tag.ReadFrom(f)
	if err != nil {
		log.Errorf("parseTrack: %s", err)
		return Track{}, err
	}

	trackNum, _ := m.Track()

	res := Track{
		ID:       trackIndex,
		Path:     path,
		Filename: tr.Name(),
		ModTime:  tr.ModTime(),
		Title:    m.Title(),
		Artist:   m.Artist(),
		Album:    m.Album(),
		Track:    uint16(trackNum),
		Year:     uint16(m.Year()),
		Cover:    m.Picture(),
	}
	trackIndex++
	return res, nil
}

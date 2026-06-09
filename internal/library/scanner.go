package library

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"go.senan.xyz/taglib"
)

var audioTypes = []string{"mp3", "wav", "ogg", "flac"}
var trackIndex uint16 = 0

func isSupportedAudio(name string) bool {
	_, ext, found := strings.Cut(name, ".")
	if !found {
		return false
	}

	return slices.Contains(audioTypes, ext)
}

func ScanLibrary(dir string, repoInsertFn func(Track) error) error {
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

		err = repoInsertFn(track)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}

		return nil
	})

	return err
}

func parseTrack(path string, entry os.DirEntry) (Track, error) {
	tr, err := entry.Info()
	if err != nil {
		return Track{}, fmt.Errorf("get file info: %w", err)
	}

	image, err := taglib.ReadImage(path)
	if err != nil {
		return Track{}, fmt.Errorf("image read: %w", err)
	}

	props, err := taglib.ReadProperties(path)
	if err != nil {
		return Track{}, fmt.Errorf("props read: %w", err)
	}
	cover := AlbumCover{Data: image}
	cover.ImageDesc = props.Images[0]

	m, err := taglib.ReadTags(path)
	if err != nil {
		return Track{}, fmt.Errorf("tag reading: %w", err)
	}

	rawTrack := m[taglib.TrackNumber][0]
	trackStr := strings.Split(rawTrack, "/")[0]
	trackNum, _ := strconv.ParseInt(strings.TrimSpace(trackStr), 10, 16)
	year, _ := strconv.ParseUint(strings.TrimSpace(m[taglib.Date][0]), 10, 16) //NOTE: wrap every metadata extraction in Trim?

	res := Track{
		ID:          trackIndex,
		Path:        path,
		Filename:    tr.Name(),
		ModTime:     tr.ModTime(),
		Size:        tr.Size(),
		Title:       m[taglib.Title][0],
		Artist:      m[taglib.Artist][0],
		Album:       m[taglib.Album][0],
		AlbumArtist: m[taglib.AlbumArtist][0],
		TrackNum:    uint16(trackNum),
		Year:        uint16(year),
		Cover:       &cover,
	}
	trackIndex++
	return res, nil
}

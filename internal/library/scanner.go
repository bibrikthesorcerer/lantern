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

type metadata map[string][]string

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

func firstInTag(m metadata, key string) string {
	v, ok := m[key]
	if !ok || len(v) == 0 {
		return ""
	}
	return strings.TrimSpace(v[0])
}

func extractTrackNum(m metadata) uint16 {
	if v := firstInTag(m, taglib.TrackNumber); v != "" {
		trackStr, _, _ := strings.Cut(v, "/")
		trNum, err := strconv.ParseUint(trackStr, 10, 16)
		if err != nil {
			return 0
		}
		return uint16(trNum)
	}
	return 0
}

func extractYear(m metadata) uint16 {
	v := firstInTag(m, taglib.Date)

	if len(v) >= 4 {
		v = v[:4]
	}

	year, err := strconv.ParseUint(v, 10, 16)
	if err != nil {
		return 0
	}

	return uint16(year)
}

func extractAlbumArtist(m metadata) string {
	if v := firstInTag(m, taglib.AlbumArtist); v != "" {
		return v
	}

	return firstInTag(m, taglib.Artist)
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

	res := Track{
		ID:          trackIndex,
		Path:        path,
		Filename:    tr.Name(),
		ModTime:     tr.ModTime(),
		Size:        tr.Size(),
		Title:       firstInTag(m, taglib.Title),
		Artist:      firstInTag(m, taglib.Artist),
		Album:       firstInTag(m, taglib.Album),
		AlbumArtist: extractAlbumArtist(m),
		TrackNum:    extractTrackNum(m),
		Year:        extractYear(m),
		Cover:       &cover,
	}
	trackIndex++
	return res, nil
}

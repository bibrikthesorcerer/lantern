package library

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.senan.xyz/taglib"
)

func TestAlbumArtistFallback(t *testing.T) {
	m := metadata{
		taglib.Artist: []string{"Pink Floyd"},
	}
	require.Equal(t, extractAlbumArtist(m), "Pink Floyd")

	m = metadata{
		taglib.Artist: []string{"Pink Floyd"},
	}
	require.Equal(t, extractAlbumArtist(m), "Pink Floyd")
}

func TestYearParsing(t *testing.T) {
	m := metadata{
		taglib.Date: []string{"1969"},
	}
	require.Equal(t, extractYear(m), uint16(1969))

	m = metadata{
		taglib.Date: []string{"1969-10-10"},
	}
	require.Equal(t, extractYear(m), uint16(1969))

	m = metadata{
		taglib.Date: []string{"abcedfg-30-234"},
	}
	require.Equal(t, extractYear(m), uint16(0))
}

func TestTrackNumParsing(t *testing.T) {
	// id3
	m := metadata{
		taglib.TrackNumber: []string{"01"},
	}
	require.Equal(t, extractTrackNum(m), uint16(1))

	// vorbis
	m = metadata{
		taglib.TrackNumber: []string{"01/12"},
	}
	require.Equal(t, extractTrackNum(m), uint16(1))

	// iTunes
	m = metadata{
		taglib.TrackNumber: []string{"01", "12"},
	}
	require.Equal(t, extractTrackNum(m), uint16(1))
}

func TestTagsExtractionWithFallback(t *testing.T) {
	m := metadata{
		taglib.Title:  []string{"03", "Rhubarb"},
		taglib.Artist: []string{"Aphex Twin"},
		taglib.Album:  []string{"       Selected Ambient Works Volume II        "},
	}

	require.Equal(t, firstInTag(m, taglib.Title), "03")
	require.Equal(t, firstInTag(m, taglib.Artist), "Aphex Twin")
	require.Equal(t, firstInTag(m, taglib.Album), "Selected Ambient Works Volume II")

	m = metadata{}
	require.Equal(t, firstInTag(m, taglib.Title), "")
}

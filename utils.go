package main

import (
	"maps"
	"slices"
	"strings"
)

func (s *MusicServer) search(query string) []Track {
	query = strings.ToLower(query)
	res := make([]Track, 0)
	tracks := slices.Collect(maps.Values(s.library))
	for _, tr := range tracks {
		if strings.Contains(strings.ToLower(tr.Title), query) ||
			strings.Contains(strings.ToLower(tr.Artist), query) ||
			strings.Contains(strings.ToLower(tr.Album), query) {
			res = append(res, tr)
		}
	}
	return res
}

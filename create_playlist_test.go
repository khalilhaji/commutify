package main

import (
	"strconv"
	"testing"

	"github.com/khalilhaji/commutify/spotify"
)

type MockSpotifyClient struct {
	name   string
	tracks []spotify.Song
	length int
}

func (m *MockSpotifyClient) UserTracks() ([]spotify.Song, error) {
	res := make([]spotify.Song, m.length)

	for i := range res {
		res[i] = spotify.Song{Duration: 1000, ID: strconv.Itoa(i), Popularity: i % 100}
	}

	return res, nil
}

func (m *MockSpotifyClient) CreatePlaylist(name string, tracks []spotify.Song) error {
	m.name = name
	m.tracks = tracks
	return nil
}

func TestCreatePlaylist(t *testing.T) {
	mock := &MockSpotifyClient{length: 1000}

	tracks, err := mock.UserTracks()

	err = createPlaylist("test", 1000, mock, func(_ []spotify.Song, _ int) []spotify.Song {
		return tracks
	})
	if err != nil {
		t.Error(err)
	}

	if mock.name != "test" {
		t.Error("Error setting playlist name")
	}

	for i, v := range mock.tracks {
		if v.ID != tracks[i].ID {
			t.Error("tracks do not match")
		}
	}
}

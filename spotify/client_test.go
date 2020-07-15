package spotify

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserTracks(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		assert.Equal("/me/tracks", r.URL.String())
		w.Header().Set("Content-Type", "application/json")
		json, err := ioutil.ReadFile("sample_data/tracks.json")
		assert.NoError(err)
		_, err = w.Write(json)
		assert.NoError(err)
	}))
	defer ts.Close()

	cli := client{&http.Client{}, "test_userid", ts.URL}
	songs, err := cli.UserTracks()
	assert.NoError(err)
	assert.Equal([]Song{{Duration: 137040, ID: "2jpDioAB9tlYXMdXDK3BGl", Popularity: 19}}, songs)
}

func TestUserTracksUnavailable(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		assert.Equal("/me/tracks", r.URL.String())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Server error"))
	}))
	defer ts.Close()

	cli := client{&http.Client{}, "test_userid", ts.URL}

	_, err := cli.UserTracks()
	assert.Error(err)
}

func TestCreatePlaylist(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.String() {
		case "/users/test_userid/playlists":
			createReq := &struct {
				Name, Description string
			}{}
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			err = json.Unmarshal(body, createReq)
			assert.NoError(err)
			assert.Equal("A New Playlist", createReq.Name)
			assert.Equal("Created by commutify", createReq.Description)
			json, err := ioutil.ReadFile("sample_data/create_playlist.json")
			assert.NoError(err)
			_, err = w.Write(json)
			assert.NoError(err)
		case "/playlists/7d2D2S200NyUE5KYs80PwO/tracks":
			addReq := &struct {
				Uris []string
			}{}
			body, err := ioutil.ReadAll(r.Body)
			assert.NoError(err)
			err = json.Unmarshal(body, addReq)
			assert.NoError(err)
			assert.Equal([]string{"spotify:track:songid1", "spotify:track:songid2", "spotify:track:songid3"}, addReq.Uris)
			json, err := ioutil.ReadFile("sample_data/add_to_playlist.json")
			assert.NoError(err)
			w.WriteHeader(http.StatusCreated)
			_, err = w.Write(json)
			assert.NoError(err)
		}
	}))
	defer ts.Close()

	cli := client{&http.Client{}, "test_userid", ts.URL}

	err := cli.CreatePlaylist("A New Playlist",
		[]Song{
			{ID: "songid1"},
			{ID: "songid2"},
			{ID: "songid3"},
		})
	assert.NoError(err)
}

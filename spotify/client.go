// Package spotify provides authorization flow and playlist creation
package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type clientConfig struct {
	ClientID, ClientSecret string
}

// Client provides functionality for playlist operations offered by the spotify API.
type Client struct {
	apiClient *http.Client
	userID    string
}

// NewSpotifyClient creates a new spotify client by receiving an authorization token from the spotify api.
func NewSpotifyClient() (*Client, error) {

	rawConf, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	conf := &clientConfig{}
	err = json.Unmarshal(rawConf, conf)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	var config = oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint:     spotify.Endpoint,
		RedirectURL:  "http://localhost:8080",
		Scopes:       []string{"playlist-modify-private", "playlist-read-private", "user-library-read"},
	}

	codeChannel := make(chan string)
	var wg sync.WaitGroup
	wg.Add(1)
	go receiveToken(ctx, codeChannel, &wg)

	url := config.AuthCodeURL("state")
	exec.Command("open", url).Start()

	fmt.Println("Authorizing")

	code := <-codeChannel

	tok, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client := config.Client(ctx, tok)
	wg.Wait()
	resp, err := client.Get("https://api.spotify.com/v1/me")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var parsedBody map[string]interface{}
	json.Unmarshal(body, &parsedBody)
	uid, ok := parsedBody["id"].(string)
	if !ok {
		log.Fatal("could not retrieve user id")
	}

	return &Client{client, uid}, nil
}

func receiveToken(ctx context.Context, ch chan string, wg *sync.WaitGroup) {
	s := http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer s.Shutdown(ctx)

		values, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			close(ch)
			return
		}

		code, ok := values["code"]
		if !ok || len(code) != 1 {
			close(ch)
			return
		}

		ch <- code[0]

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, err = w.Write([]byte("<script>window.close();</script>"))
		w.(http.Flusher).Flush()
		if err != nil {
			log.Fatal(err)
		}
	})

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	wg.Done()
}

// Song represents a single track with an id and a length in milliseconds
type Song struct {
	Duration int    `json:"duration_ms"`
	ID       string `json:"id"`
}

// UserTracks gets a list of all the songs in the spotify user's library.
func (c *Client) UserTracks() ([]Song, error) {
	resp, err := c.apiClient.Get("https://api.spotify.com/v1/me/tracks")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	type apiResponse struct {
		Songs []struct {
			Track Song
		} `json:"items"`
		Next string
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsedResp apiResponse
	json.Unmarshal(body, &parsedResp)

	res := make([]Song, len(parsedResp.Songs))
	for i, v := range parsedResp.Songs {
		res[i] = v.Track
	}

	for parsedResp.Next != "" {
		resp, err = c.apiClient.Get(parsedResp.Next)

		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		parsedResp = apiResponse{}

		body, err = ioutil.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		json.Unmarshal(body, &parsedResp)

		for _, v := range parsedResp.Songs {
			res = append(res, v.Track)
		}
	}

	return res, nil
}

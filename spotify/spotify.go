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

type spotifyClient struct {
	apiClient *http.Client
}

func NewSpotifyClient() (*spotifyClient, error) {

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
		Scopes:       []string{"playlist-modify-private", "playlist-read-private"},
	}

	codeChannel := make(chan string)

	var wg sync.WaitGroup
	// wg.Add(1)

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
	return &spotifyClient{client}, nil
}

func (sc *spotifyClient) GetProfile() (*http.Response, error) {
	return sc.apiClient.Get("https://api.spotify.com/v1/me")
}

func receiveToken(ctx context.Context, ch chan string, wg *sync.WaitGroup) {
	s := http.Server{Addr: ":8080"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// defer s.Shutdown(ctx)
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

		fmt.Println("Code:", code)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write([]byte("<script>window.close();</script>"))
		if err != nil {
			log.Fatal(err)
		}
		ch <- code[0]
	})

	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	// wg.Done()
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/khalilhaji/commutify/maps"

	"github.com/khalilhaji/commutify/spotify"
)

// SpotifyClient represents the functionality provided by the client in the spotify package
type SpotifyClient interface {
	UserTracks() ([]spotify.Song, error)

	CreatePlaylist(string, []spotify.Song) error
}

// Stores tokens necessary for maps api
type mapConfig struct {
	GeocodioToken   string
	CitymapperToken string
}

func main() {
	rawConf, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal(err)
	}
	var conf mapConfig
	if err := json.Unmarshal(rawConf, &conf); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enter start address:")
	stdinReader := bufio.NewReader(os.Stdin)
	start, err := stdinReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Enter end address:")
	end, err := stdinReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	startCoord, err := maps.GetCoordinates(start, conf.GeocodioToken)
	if err != nil {
		log.Fatal(err)
	}

	endCoord, err := maps.GetCoordinates(end, conf.GeocodioToken)
	if err != nil {
		log.Fatal(err)
	}

	duration, err := maps.GetTime(startCoord, endCoord, conf.CitymapperToken)
	if err != nil || duration == 0 {
		fmt.Println(err)
		fmt.Println("Error retrieving transit directions, enter duration manually in minutes:")
		_, err = fmt.Scanf("%d", &duration)
	}

	if err != nil {
		log.Fatal(err)
	}

	client, err := spotify.NewSpotifyClient()

	if err != nil || client == nil {
		log.Fatal(err)
	}

	fmt.Println("Enter a name for the new playlist:")

	name, err := stdinReader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	if err := createPlaylist(name, duration*60, client, knapsack); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Playlist created successfully, check your spotify library")
}

// Creates a spotify playlist with the given name and length through the given client interface object
func createPlaylist(name string, length int, client SpotifyClient, trackProc func([]spotify.Song, int) []spotify.Song) (err error) {
	tracks, err := client.UserTracks()
	if err != nil {
		return err
	}

	shuffleTracks(&tracks)

	finalTracks := trackProc(tracks, length)

	if err := client.CreatePlaylist(name, finalTracks); err != nil {
		return err
	}
	return
}

// Shuffles a list of songs to ensure that the same songs aren't returned for each query
func shuffleTracks(tracks *[]spotify.Song) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*tracks), func(i, j int) {
		(*tracks)[i], (*tracks)[j] = (*tracks)[j], (*tracks)[i]
	})
}

// Implementation of the solution to the knapsack problem.
// The popularity of each track is its value and the duration is its weight.
func knapsack(songs []spotify.Song, duration int) []spotify.Song {
	dp := make([][]int, len(songs)+1)

	for i := range dp {
		dp[i] = make([]int, duration+1)
	}

	for i := 1; i < len(dp); i++ {
		for j := 0; j < duration+1; j++ {
			if (songs[i-1].Duration / 1000) > j {
				dp[i][j] = dp[i-1][j]
			} else {
				prev := dp[i-1][j]
				withCurr := songs[i-1].Popularity + dp[i-1][duration-(songs[i-1].Duration/1000)]
				if prev > withCurr {
					dp[i][j] = prev
				} else {
					dp[i][j] = withCurr
				}
			}
		}
	}

	var res []spotify.Song

	maxVal := dp[len(songs)][duration]
	w := duration
	for i := len(dp) - 1; i > 0; i-- {
		if dp[i][w] == dp[i-1][w] {
			continue
		} else {
			res = append(res, songs[i-1])
			maxVal = maxVal - songs[i-1].Popularity
			w = w - (songs[i-1].Duration / 1000)
		}
	}
	return res
}

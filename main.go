package main

import (
	"fmt"
	"log"

	"github.com/khalilhaji/commutify/spotify"
)

func main() {

	client, err := spotify.NewSpotifyClient()

	if err != nil || client == nil {
		log.Fatal(err)
	}
	fmt.Println("Authorization successful!")

	fmt.Println("Grabbing user songs")

	tracks, err := client.UserTracks()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(tracks), "tracks acquired")
}

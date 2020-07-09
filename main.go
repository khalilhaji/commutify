package main

import (
	"fmt"
	"log"

	"github.com/khalilhaji/commutify/spotify"
)

func main() {

	client, err := spotify.NewSpotifyClient()

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.GetProfile()
	if resp.StatusCode == 200 {
		fmt.Println("Authorization successful!")
	}

}

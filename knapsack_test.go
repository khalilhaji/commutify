package main

import (
	"testing"

	"github.com/khalilhaji/commutify/spotify"
)

func TestAllIncluded(t *testing.T) {
	songs := []spotify.Song{
		spotify.Song{Duration: 180000, ID: "1", Popularity: 50},
		spotify.Song{Duration: 120000, ID: "2", Popularity: 70},
		spotify.Song{Duration: 360000, ID: "3", Popularity: 90},
	}

	res := knapsack(songs, 660)

	if len(res) != 3 {
		t.Errorf("Length should be 3, got: %d", len(res))
	}

	sum := 0
	for _, v := range res {
		sum += v.Duration
	}
	if sum != 660000 {
		t.Errorf("Duration should be 660, got %d", sum)
	}
}

func TestSome(t *testing.T) {
	songs := []spotify.Song{
		spotify.Song{Duration: 180000, ID: "1", Popularity: 50},
		spotify.Song{Duration: 120000, ID: "2", Popularity: 70},
		spotify.Song{Duration: 360000, ID: "3", Popularity: 90},
		spotify.Song{Duration: 660000, ID: "4", Popularity: 100},
	}

	res := knapsack(songs, 660)

	if len(res) != 3 {
		t.Errorf("Length should be 3, got: %d", len(res))
	}
}

func TestLong(t *testing.T) {
	songs := []spotify.Song{
		spotify.Song{Duration: 100000, ID: "1", Popularity: 50},
		spotify.Song{Duration: 100000, ID: "2", Popularity: 60},
		spotify.Song{Duration: 100000, ID: "3", Popularity: 70},
		spotify.Song{Duration: 100000, ID: "4", Popularity: 80},
	}

	res := knapsack(songs, 300)

	if len(res) != 3 {
		t.Errorf("Length should be 3, got: %d", len(res))
	}

	totalValue := 0
	for _, v := range res {
		totalValue += v.Popularity
	}

	if totalValue != 210 {
		t.Errorf("Total value should be 210, got: %d", totalValue)
	}
}

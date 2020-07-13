package maps

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Coordinate represent coordinates on the map
type Coordinate struct {
	Lat, Lng float64
}

// QueryString escapes the coordinates to use in a url query.
func (c Coordinate) QueryString() string {
	return url.QueryEscape(fmt.Sprintf("%v, %v", c.Lat, c.Lng))
}

// GetCoordinates from an address using geocodio API
func GetCoordinates(address, token string) (Coordinate, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.geocod.io/v1.6/geocode?q=%s&api_key=%s&limit=1", url.QueryEscape(address), token))
	if err != nil {
		return Coordinate{}, err
	} else if resp.StatusCode != 200 {
		return Coordinate{}, fmt.Errorf("Geocodio returned status code: %d", resp.StatusCode)
	}

	type GeocodeResponse struct {
		Results []struct {
			Location Coordinate `json:"location"`
		} `json:"results"`
	}

	var respUnmarshal GeocodeResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Coordinate{}, err
	}

	if err := json.Unmarshal(body, &respUnmarshal); err != nil {
		return Coordinate{}, err
	}

	return respUnmarshal.Results[0].Location, nil
}

// GetTime from start coordinate to end in minutes.
func GetTime(start, end Coordinate, token string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("https://developer.citymapper.com/api/1/traveltime/?startcoord=%s&endcoord=%s&key=%s",
		start.QueryString(),
		end.QueryString(),
		token))
	if err != nil {
		return 0, err
	} else if resp.StatusCode != 200 {
		return 0, fmt.Errorf("Citymapper returned status code: %d", resp.StatusCode)
	}

	type travelTimeResponse struct {
		Time int `json:"travel_time_minutes"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	var ttime travelTimeResponse
	if err := json.Unmarshal(body, &ttime); err != nil {
		return 0, err
	}
	return ttime.Time, nil
}

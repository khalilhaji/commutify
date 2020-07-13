# commutify
A tool that creates a Spotify playlist exactly the length of your commute

![Go](https://github.com/khalilhaji/commutify/workflows/Go/badge.svg)


# Setup

Create `config.json` and include your spotify `clientID` and `clientSecret` along with `geocodioToken` and `citymapperToken`

# Todos

## Spotify API
* [x] Create spotify app
* [x] Get time from stdin
* [x] Authenticate user
* [x] Grab liked songs from user
* [x] Create song list for new playlist with length within tolerance

## ~~Google Maps~~ Citymapper API
* [x] Accept address from stdin
* [x] Feed time into playlist generator

## Future Improvements
* [ ] Account for pausing during boarding/transferring between train and bus
* [ ] Prompt user for route to determine time
* [ ] Adjust BPM based on morning/evening commute and transportation method/speed

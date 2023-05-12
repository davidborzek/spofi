package views

import (
	"log"

	"github.com/davidborzek/spofi/pkg/rofi"
)

func addQueueError(err error) {
	rofi.Error("Failed to add the track to the queue. Try again.")
	log.Println(err)
}

func playTrackError(err error) {
	rofi.Error("Failed to play the track. Try again.")
	log.Println(err)
}

func playAlbumError(err error) {
	rofi.Error("Failed to play the album. Try again.")
	log.Println(err)
}

func playPlaylistError(err error) {
	rofi.Error("Failed to play the playlist. Try again.")
	log.Println(err)
}

func getAlbumError(err error) {
	rofi.Error("Failed to get the album. Try again.")
	log.Println(err)
}

func getAlbumsError(err error) {
	rofi.Error("Failed to get albums. Try again.")
	log.Println(err)
}

func getPlaylistsError(err error) {
	rofi.Error("Failed to get playlists. Try again.")
	log.Println(err)
}

func getTracksError(err error) {
	rofi.Error("Failed to get tracks. Try again.")
	log.Println(err)
}

func selectDeviceError(err error) {
	rofi.Error("Failed to select the device. Try again.")
	log.Println(err)
}

func getDevicesError(err error) {
	rofi.Error("Failed to get available devices. Try again.")
	log.Println(err)
}

func noDevicesFoundError() {
	rofi.Error("No devices found.")
}

func getPlayerStateError(err error) {
	rofi.Error("Failed to get player status. Try again.")
	log.Println(err)
}

func playPauseError(err error) {
	rofi.Error("Failed to pause/resume. Try again.")
	log.Println(err)
}

func skipTrackError(err error) {
	rofi.Error("Failed to skip track. Try again.")
	log.Println(err)
}

func previousTrackError(err error) {
	rofi.Error("Failed to go to previous track. Try again.")
	log.Println(err)
}

func updatePlayerError(err error) {
	rofi.Error("Failed to update player. Try again.")
	log.Println(err)
}

func getQueueError(err error) {
	rofi.Error("Failed to get queue. Try again.")
	log.Println(err)
}

func getRecentlyPlayedTracksError(err error) {
	rofi.Error("Failed to get recently played tracks. Try again.")
	log.Println(err)
}

func searchError(err error) {
	rofi.Error("Failed to search. Try again.")
	log.Println(err)
}

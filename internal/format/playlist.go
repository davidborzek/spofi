package format

import (
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

// FormatPlaylistRows formats each playlist as row for rofi.
func FormatPlaylistRows(playlists []spotify.BasePlaylist, icon string) []rofi.Row {
	data := make([][]string, len(playlists))
	for i, playlist := range playlists {
		data[i] = []string{
			playlist.Name,
		}
	}

	rawRows := BuildRows(data, 30)
	rows := make([]rofi.Row, len(playlists))
	for i, rawRow := range rawRows {
		rows[i] = rofi.Row{
			Title: FormatIcon(icon, rawRow),
			Value: playlists[i].URI,
		}
	}

	return rows
}

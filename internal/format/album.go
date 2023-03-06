package format

import (
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

// FormatAlbumRows formats each album as row for rofi.
func FormatAlbumRows(albums []spotify.Album, icon string) []rofi.Row {
	data := make([][]string, len(albums))
	for i, album := range albums {
		data[i] = []string{
			album.Name,
			album.Artists[0].Name,
		}
	}

	rawRows := BuildRows(data, 30)
	rows := make([]rofi.Row, len(albums))
	for i, rawRow := range rawRows {
		rows[i] = rofi.Row{
			Title: FormatIcon(icon, rawRow),
			Value: albums[i].URI,
		}
	}

	return rows
}

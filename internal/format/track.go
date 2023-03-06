package format

import (
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

// FormatTrackRows formats each track as row for rofi.
func FormatTrackRows(tracks []spotify.Track, icon string) []rofi.Row {
	data := make([][]string, len(tracks))
	for i, track := range tracks {
		data[i] = []string{
			track.Name,
			track.Artists[0].Name,
		}
	}

	rawRows := BuildRows(data, 30)
	rows := make([]rofi.Row, len(tracks))
	for i, rawRow := range rawRows {
		rows[i] = rofi.Row{
			Title: FormatIcon(icon, rawRow),
			Value: tracks[i].URI,
		}
	}

	return rows
}

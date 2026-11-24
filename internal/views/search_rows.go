package views

import (
	"fmt"
	"strings"

	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

const (
	// rowActionTrack plays a single track without a context.
	rowActionTrack = "track"
	// rowActionContext plays a track inside its context
	// (playlist, album, liked songs, ...).
	rowActionContext = "context"

	// maxContextMatches limits how many matches of the current
	// context are shown, so the list does not get flooded.
	maxContextMatches = 10

	// likedSongsContextURI is the playback context of the
	// liked songs collection.
	likedSongsContextURI = "spotify:collection:tracks"
)

// searchAction describes what should happen when a row is selected.
type searchAction struct {
	// Kind is either rowActionTrack or rowActionContext.
	Kind string
	// Context is the context uri to play the track in.
	Context string
	// Track is the uri of the track.
	Track string
}

// encodeActionResult packs the action into a rofi row value.
func encodeActionResult(action searchAction) string {
	return strings.Join([]string{action.Kind, action.Context, action.Track}, "|")
}

// decodeActionResult unpacks a rofi row value into an action.
func decodeActionResult(value string) searchAction {
	parts := strings.SplitN(value, "|", 3)
	if len(parts) != 3 {
		return searchAction{
			Kind:  rowActionTrack,
			Track: value,
		}
	}

	return searchAction{
		Kind:    parts[0],
		Context: parts[1],
		Track:   parts[2],
	}
}

// searchResults holds the grouped results of a search.
type searchResults struct {
	// tracks are the plain search results of the track search.
	tracks []spotify.Track
	// context are the matches inside the context that is
	// currently playing (playlist or album).
	context []spotify.Track
	// contextURI is the uri of that context.
	contextURI string
	// contextLabel describes the context, e.g. "playlist: Chill".
	contextLabel string
	// favorites are the search results that are saved
	// in the liked songs of the user.
	favorites []spotify.Track
}

// buildSearchRows groups the search results by context and
// returns them as rofi rows. The order is:
// current context, liked songs, single track, album and
// finally the remaining search results.
func buildSearchRows(results searchResults, icon string) []rofi.Row {
	rows := make([]rofi.Row, 0, len(results.tracks)+len(results.context))

	seen := make(map[string]bool)

	add := func(track spotify.Track, label string, action searchAction) {
		if track.URI == "" {
			return
		}

		rows = append(rows, rofi.Row{
			Title: format.FormatIcon(
				icon,
				fmt.Sprintf("%s  (%s)", formatTrackTitle(track), label),
			),
			Value: encodeActionResult(action),
		})
	}

	for i, track := range uniqueTracks(results.context) {
		if i >= maxContextMatches {
			break
		}
		seen[track.URI] = true
		add(track, results.contextLabel, searchAction{
			Kind:    rowActionContext,
			Context: results.contextURI,
			Track:   track.URI,
		})
	}

	for i, track := range uniqueTracks(results.favorites) {
		if i >= maxContextMatches {
			break
		}
		add(track, "liked songs", searchAction{
			Kind:    rowActionContext,
			Context: likedSongsContextURI,
			Track:   track.URI,
		})
	}

	if len(results.tracks) > 0 {
		top := results.tracks[0]
		seen[top.URI] = true
		add(top, "single track", searchAction{
			Kind:  rowActionTrack,
			Track: top.URI,
		})

		if top.Album.URI != "" {
			add(top, fmt.Sprintf("album: %s", top.Album.Name), searchAction{
				Kind:    rowActionContext,
				Context: top.Album.URI,
				Track:   top.URI,
			})
		}
	}

	for _, track := range uniqueTracks(results.tracks) {
		if seen[track.URI] {
			continue
		}
		add(track, "single track", searchAction{
			Kind:  rowActionTrack,
			Track: track.URI,
		})
	}

	return rows
}

// uniqueTracks removes duplicate tracks (by uri)
// while keeping the order.
func uniqueTracks(tracks []spotify.Track) []spotify.Track {
	seen := make(map[string]bool, len(tracks))
	unique := make([]spotify.Track, 0, len(tracks))
	for _, track := range tracks {
		if track.URI == "" || seen[track.URI] {
			continue
		}
		seen[track.URI] = true
		unique = append(unique, track)
	}

	return unique
}

// formatTrackTitle formats a track as "name | artist".
func formatTrackTitle(track spotify.Track) string {
	name := format.Truncate(track.Name, 34)

	if len(track.Artists) == 0 {
		return name
	}

	return fmt.Sprintf("%s | %s", name, format.Truncate(track.Artists[0].Name, 22))
}

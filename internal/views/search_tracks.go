package views

import (
	"fmt"
	"log"
	"strings"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

const (
	// playlistScanLimit is the maximum amount of tracks that are
	// scanned to find matches in the currently playing playlist.
	playlistScanLimit = 2000
)

type searchTracksView struct {
	rofi rofi.App
	app  *app.App

	parent View

	query string
}

func NewSearchTrackView(app *app.App) *searchTracksView {
	var msg = ""
	if app.Config.ShowKeybindings {
		msg = format.FormatKeybindings(
			format.Keybinding{
				Key:         app.Config.Keybindings.ToggleSearchType,
				Description: "Toggle search type",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.AddToQueue,
				Description: "Add to queue",
			},
		)
	}

	title := format.FormatIcon(
		app.Config.Icons.Track,
		"Tracks",
	)

	r := rofi.App{
		Prompt: title,
		Keybindings: []string{
			app.Config.Keybindings.AddToQueue,
			app.Config.Keybindings.ToggleSearchType,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
		Message:    msg,
	}

	view := &searchTracksView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *searchTracksView) SetParent(parent View) {
	view.parent = parent
}

func (view *searchTracksView) SetQuery(query string) {
	view.query = query
}

func (view *searchTracksView) Show(payload ...interface{}) {
	if view.query == "" {
		return
	}

	rows, err := view.search()
	if err != nil {
		searchError(err)
		return
	}

	view.rofi.Rows = rows

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.KeyEvent:
		switch evt.Key {
		case view.app.Config.Keybindings.AddToQueue:
			action := decodeActionResult(evt.Selection.Value)
			err := view.app.Player.AddQueue(action.Track)
			if err != nil {
				addQueueError(err)
			}
			view.Show()
		case view.app.Config.Keybindings.ToggleSearchType:
			albumSearch := NewSearchAlbumsView(view.app)
			albumSearch.SetParent(view.parent)
			albumSearch.SetQuery(view.query)
			albumSearch.Show()
		}
	case rofi.SelectedEvent:
		view.play(decodeActionResult(evt.Selection.Value))
	}
}

// play plays the track of the given action. Context actions
// start (or jump into) the context of the track, so the
// following tracks keep playing.
func (view *searchTracksView) play(action searchAction) {
	if action.Track == "" {
		return
	}

	if action.Kind != rowActionContext || action.Context == "" {
		if err := view.app.Player.PlayTrack(action.Track); err != nil {
			playTrackError(err)
		}
		return
	}

	if err := view.app.Player.PlayContext(action.Context, action.Track); err != nil {
		// Not every context is playable (e.g. the liked songs
		// collection). Fall back to playing the track alone.
		if fallbackErr := view.app.Player.PlayTrack(action.Track); fallbackErr != nil {
			playTrackError(fallbackErr)
		}
	}
}

func (view *searchTracksView) search() ([]rofi.Row, error) {
	response, err := view.app.SpotifyClient.Search(view.query, "track")
	if err != nil {
		return nil, err
	}

	results := searchResults{
		tracks: response.Tracks.Items,
	}

	results.context, results.contextURI, results.contextLabel = view.searchCurrentContext()
	results.favorites = view.searchFavorites(results.tracks)

	return buildSearchRows(results, view.app.Config.Icons.Track), nil
}

// searchCurrentContext searches the context that is currently
// playing (playlist or album) for the query.
func (view *searchTracksView) searchCurrentContext() ([]spotify.Track, string, string) {
	player, err := view.app.SpotifyClient.GetPlayer()
	if err != nil || player == nil || player.Context == nil {
		return nil, "", ""
	}

	context := player.Context
	id := spotify.URIToID(context.URI)
	if id == "" {
		return nil, "", ""
	}

	switch context.Type {
	case "playlist":
		label := "current playlist"
		if playlist, err := view.app.SpotifyClient.GetPlaylist(id); err == nil && playlist.Name != "" {
			label = fmt.Sprintf("playlist: %s", playlist.Name)
		}

		tracks, err := view.getPlaylistTracks(id)
		if err != nil {
			return nil, "", ""
		}

		return filterTracks(tracks, view.query), context.URI, label
	case "album":
		album, err := view.app.SpotifyClient.GetAlbum(id)
		if err != nil {
			return nil, "", ""
		}

		return filterTracks(
			album.Tracks.Items,
			view.query,
		), context.URI, fmt.Sprintf("album: %s", album.Name)
	}

	return nil, "", ""
}

// getPlaylistTracks fetches the tracks of a playlist until
// playlistScanLimit is reached.
func (view *searchTracksView) getPlaylistTracks(playlistID string) ([]spotify.Track, error) {
	const pageSize = 50

	tracks := make([]spotify.Track, 0, pageSize)
	offset := 0
	for offset < playlistScanLimit {
		page, err := view.app.SpotifyClient.GetPlaylistTracks(playlistID, pageSize, offset)
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, page.Tracks()...)
		offset += len(page.Items)

		if len(page.Items) < pageSize || offset >= page.Total {
			break
		}
	}

	return tracks, nil
}

// searchFavorites returns the search results that are saved
// in the liked songs of the user.
func (view *searchTracksView) searchFavorites(tracks []spotify.Track) []spotify.Track {
	ids := make([]string, 0, len(tracks))
	indexes := make([]int, 0, len(tracks))
	for i, track := range tracks {
		if track.ID == "" {
			continue
		}
		ids = append(ids, track.ID)
		indexes = append(indexes, i)
	}

	if len(ids) == 0 {
		return nil
	}

	saved, err := view.app.SpotifyClient.AreTracksSaved(ids)
	if err != nil {
		return nil
	}

	favorites := make([]spotify.Track, 0, len(saved))
	for i, isSaved := range saved {
		if isSaved && i < len(indexes) {
			favorites = append(favorites, tracks[indexes[i]])
		}
	}

	return favorites
}

// filterTracks returns the tracks that match the query
// by name or artist.
func filterTracks(tracks []spotify.Track, query string) []spotify.Track {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}

	matches := make([]spotify.Track, 0)
	for _, track := range tracks {
		if strings.Contains(strings.ToLower(track.Name), q) {
			matches = append(matches, track)
			continue
		}

		for _, artist := range track.Artists {
			if strings.Contains(strings.ToLower(artist.Name), q) {
				matches = append(matches, track)
				break
			}
		}
	}

	return matches
}

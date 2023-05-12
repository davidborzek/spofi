package views

import (
	"fmt"
	"log"
	"os"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

const (
	devicesViewID        = "devices_view"
	likedTracksViewID    = "liked_tracks_view"
	playerViewID         = "player_view"
	playlistsViewID      = "playlists_view"
	queueViewID          = "queue_view"
	recentlyPlayedViewID = "recently_played_view"
	savedAlbumsViewID    = "saved_albums_view"
	searchViewID         = "search_view"
)

type mainView struct {
	app    *app.App
	parent View

	rofi rofi.App

	devicesView        View
	likedTracksView    View
	playerView         View
	playlistsView      View
	queueView          View
	recentlyPlayedView View
	savedAlbumsView    View
	searchTracksView   *searchTracksView
	searchView         View
}

func NewMainView(app *app.App) View {
	devicesViewTitle := format.FormatIcon(
		app.Config.Icons.Device,
		"Devices",
	)

	playerViewTitle := format.FormatIcon(
		app.Config.Icons.Player,
		"Player",
	)

	likedTracksViewTitle := format.FormatIcon(
		app.Config.Icons.LikedTracks,
		"Liked Tracks",
	)

	playlistsViewTitle := format.FormatIcon(
		app.Config.Icons.Playlist,
		"Playlists",
	)

	queueViewTitle := format.FormatIcon(
		app.Config.Icons.Queue,
		"Queue",
	)

	recentlyPlayedViewTitle := format.FormatIcon(
		app.Config.Icons.RecentlyPlayed,
		"Recently Played",
	)

	savedAlbumsViewTitle := format.FormatIcon(
		app.Config.Icons.Album,
		"Albums",
	)

	searchViewTitle := format.FormatIcon(
		app.Config.Icons.Search,
		"Search",
	)

	r := rofi.App{
		IgnoreCase: true,
		KBCustom: []string{
			app.Config.Keybindings.TogglePauseResume,
			app.Config.Keybindings.NextTrack,
			app.Config.Keybindings.PreviousTrack,
			app.Config.Keybindings.ToggleRepeat,
			app.Config.Keybindings.ToggleShuffle,
		},
		Rows: []rofi.Row{
			{
				Title: playerViewTitle,
				Value: playerViewID,
			},
			{
				Title: searchViewTitle,
				Value: searchViewID,
			},
			{
				Title: likedTracksViewTitle,
				Value: likedTracksViewID,
			},
			{
				Title: savedAlbumsViewTitle,
				Value: savedAlbumsViewID,
			},
			{
				Title: playlistsViewTitle,
				Value: playlistsViewID,
			},
			{
				Title: queueViewTitle,
				Value: queueViewID,
			},
			{
				Title: recentlyPlayedViewTitle,
				Value: recentlyPlayedViewID,
			},
			{
				Title: devicesViewTitle,
				Value: devicesViewID,
			},
		},
	}

	view := &mainView{
		rofi:               r,
		app:                app,
		devicesView:        NewDevicesView(app, devicesViewTitle),
		likedTracksView:    NewLikedTracksView(app, likedTracksViewTitle),
		playerView:         NewPlayerView(app, playerViewTitle),
		playlistsView:      NewPlaylistsView(app, playlistsViewTitle),
		queueView:          NewQueueView(app, queueViewTitle),
		recentlyPlayedView: NewRecentlyPlayedView(app, recentlyPlayedViewTitle),
		savedAlbumsView:    NewSavedAlbumsView(app, savedAlbumsViewTitle),
		searchTracksView:   NewSearchTrackView(app),
		searchView:         NewSearchView(app, searchViewTitle),
	}

	view.devicesView.SetParent(view)
	view.likedTracksView.SetParent(view)
	view.playerView.SetParent(view)
	view.playlistsView.SetParent(view)
	view.queueView.SetParent(view)
	view.recentlyPlayedView.SetParent(view)
	view.savedAlbumsView.SetParent(view)
	view.searchTracksView.SetParent(view)
	view.searchView.SetParent(view)

	return view
}

func (view *mainView) buildPlayerMessage() string {
	player, err := view.app.SpotifyClient.GetPlayer()
	if err != nil {
		getPlayerStateError(err)
		os.Exit(1)
	}

	currentlyPlaying := "Nothing is currently playing."
	if player != nil {
		status := view.app.Config.Icons.Pause
		if player.IsPlaying {
			status = view.app.Config.Icons.Play
		}

		shuffle := view.app.Config.Icons.ShuffleOff
		if player.ShuffleState {
			shuffle = view.app.Config.Icons.ShuffleOn
		}

		repeat := view.app.Config.Icons.RepeatOff
		if player.RepeatState == spotify.RepeatContext {
			repeat = view.app.Config.Icons.RepeatContext
		} else if player.RepeatState == spotify.RepeatTrack {
			repeat = view.app.Config.Icons.RepeatTrack
		}

		title := format.FormatTitle(player.Item.Name, player.Item.Artists[0].Name)
		currentlyPlaying = fmt.Sprintf(
			"%s | %s | %s %s",
			status,
			title,
			shuffle,
			repeat,
		)
	}

	return currentlyPlaying
}

func (view *mainView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.KBCustom1 {
		if err := view.app.Player.PlayPause(); err != nil {
			playPauseError(err)
		}
		view.Show()
		return
	}

	if code == rofi.KBCustom2 {
		if err := view.app.Player.Next(); err != nil {
			skipTrackError(err)
		}
		view.Show()
		return
	}

	if code == rofi.KBCustom3 {
		if err := view.app.Player.Previous(); err != nil {
			previousTrackError(err)
		}
		view.Show()
		return
	}

	if code == rofi.KBCustom4 {
		if err := view.app.Player.ToggleRepeat(); err != nil {
			updatePlayerError(err)
		}
		view.Show()
		return
	}

	if code == rofi.KBCustom5 {
		if err := view.app.Player.ToggleShuffle(); err != nil {
			updatePlayerError(err)
		}
		view.Show()
		return
	}

	if code > 0 {
		return
	}

	switch selection.Value {
	case playerViewID:
		view.playerView.Show()
	case devicesViewID:
		view.devicesView.Show()
	case searchViewID:
		view.searchView.Show()
	case likedTracksViewID:
		view.likedTracksView.Show()
	case queueViewID:
		view.queueView.Show()
	case recentlyPlayedViewID:
		view.recentlyPlayedView.Show()
	case savedAlbumsViewID:
		view.savedAlbumsView.Show()
	case playlistsViewID:
		view.playlistsView.Show()
	default:
		view.searchTracksView.SetQuery(selection.Title)
		view.searchTracksView.Show()
	}
}

func (view *mainView) Show(payload ...interface{}) {
	msg := view.buildPlayerMessage()
	view.rofi.Prompt = msg

	selection, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(selection, code)
}

func (view *mainView) SetParent(parent View) {
	view.parent = parent
}

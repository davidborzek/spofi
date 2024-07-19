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
	playerViewID         = "player_view"
	likedTracksViewID    = "liked_tracks_view"
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
	searchView         View
	likedTracksView    View
	queueView          View
	recentlyPlayedView View
	searchTracksView   *searchTracksView
	playerView         View
	savedAlbumsView    View
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

	var msg = ""
	if app.Config.ShowKeybindings {
		msg = format.FormatKeybindings(
			format.Keybinding{
				Key:         app.Config.Keybindings.TogglePauseResume,
				Description: "Toggle pause/resume",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.NextTrack,
				Description: "Next track",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.PreviousTrack,
				Description: "Previous track",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.ToggleRepeat,
				Description: "Toggle repeat",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.ToggleShuffle,
				Description: "Toggle shuffle",
			},
		)
	}

	r := rofi.App{
		IgnoreCase: true,
		Keybindings: []string{
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
		Message: msg,
	}

	view := &mainView{
		rofi:               r,
		app:                app,
		devicesView:        NewDevicesView(app, devicesViewTitle),
		searchView:         NewSearchView(app, searchViewTitle),
		likedTracksView:    NewLikedTracksView(app, likedTracksViewTitle),
		queueView:          NewQueueView(app, queueViewTitle),
		recentlyPlayedView: NewRecentlyPlayedView(app, recentlyPlayedViewTitle),
		searchTracksView:   NewSearchTrackView(app),
		playerView:         NewPlayerView(app, playerViewTitle),
		savedAlbumsView:    NewSavedAlbumsView(app, savedAlbumsViewTitle),
	}

	view.playerView.SetParent(view)
	view.devicesView.SetParent(view)
	view.searchView.SetParent(view)
	view.likedTracksView.SetParent(view)
	view.queueView.SetParent(view)
	view.recentlyPlayedView.SetParent(view)
	view.searchTracksView.SetParent(view)
	view.savedAlbumsView.SetParent(view)

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

func (view *mainView) Show(payload ...interface{}) {
	msg := view.buildPlayerMessage()
	view.rofi.Prompt = msg

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.KeyEvent:
		switch evt.Key {
		case view.app.Config.Keybindings.TogglePauseResume:
			if err := view.app.Player.PlayPause(); err != nil {
				playPauseError(err)
			}
		case view.app.Config.Keybindings.NextTrack:
			if err := view.app.Player.Next(); err != nil {
				skipTrackError(err)
			}
		case view.app.Config.Keybindings.PreviousTrack:
			if err := view.app.Player.Previous(); err != nil {
				previousTrackError(err)
			}
		case view.app.Config.Keybindings.ToggleRepeat:
			if err := view.app.Player.ToggleRepeat(); err != nil {
				updatePlayerError(err)
			}
		case view.app.Config.Keybindings.ToggleShuffle:
			if err := view.app.Player.ToggleShuffle(); err != nil {
				updatePlayerError(err)
			}
		}

		view.Show()
	case rofi.SelectedEvent:
		switch evt.Selection.Value {
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
		default:
			view.searchTracksView.SetQuery(evt.Selection.Title)
			view.searchTracksView.Show()
		}
	}
}

func (view *mainView) SetParent(parent View) {
	view.parent = parent
}

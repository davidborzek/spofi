package views

import (
	"fmt"
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

const (
	playerTogglePauseAction = "player_toggle_pause"
	playerNextAction        = "player_next"
	playerPreviousAction    = "player_previous"
	playerToggleShuffle     = "player_toggle_shuffle"
	playerToggleRepeat      = "player_toggle_repeat"
)

type playerView struct {
	rofi   rofi.App
	app    *app.App
	parent View
}

func NewPlayerView(app *app.App, title string) View {
	r := rofi.App{
		Prompt:       title,
		ShowBack:     true,
		NoCustom:     true,
		IgnoreCase:   true,
		RenderMarkup: true,
	}

	view := &playerView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *playerView) togglePauseResume() error {
	player, err := view.app.SpotifyClient.GetPlayer()

	if err == nil {
		if player.IsPlaying {
			err = view.app.SpotifyClient.Pause(
				view.app.Config.Device.ID,
			)
		} else {
			err = view.app.SpotifyClient.Play(
				view.app.Config.Device.ID,
			)
		}
	}

	return err
}

func (view *playerView) toggleShuffleState() error {
	player, err := view.app.SpotifyClient.GetPlayer()
	if err != nil {
		return err
	}

	if player != nil {
		return view.app.SpotifyClient.SetShuffleState(
			view.app.Config.Device.ID,
			!player.ShuffleState,
		)
	}

	return nil
}

func (view *playerView) toggleRepeatState() error {
	player, err := view.app.SpotifyClient.GetPlayer()
	if err != nil {
		return err
	}

	if player != nil {
		s := spotify.RepeatOff
		if player.RepeatState == spotify.RepeatOff {
			s = spotify.RepeatContext
		}
		if player.RepeatState == spotify.RepeatContext {
			s = spotify.RepeatTrack
		}

		return view.app.SpotifyClient.SetRepeatMode(
			view.app.Config.Device.ID,
			s,
		)
	}

	return nil
}

func (view *playerView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape {
		view.parent.Show()
		return
	}

	if code > 0 {
		return
	}

	if selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	var err error

	switch selection.Value {
	case playerTogglePauseAction:
		err = view.togglePauseResume()
	case playerNextAction:
		err = view.app.SpotifyClient.Next(view.app.Config.Device.ID)
	case playerPreviousAction:
		err = view.app.SpotifyClient.Previous(view.app.Config.Device.ID)
	case playerToggleShuffle:
		err = view.toggleShuffleState()
	case playerToggleRepeat:
		err = view.toggleRepeatState()
	}

	if err != nil {
		updatePlayerError(err)
	}

	view.Show()
}

func (view *playerView) Show(payload ...interface{}) {
	player, err := view.app.SpotifyClient.GetPlayer()
	if err != nil {
		getPlayerStateError(err)
		view.parent.Show()
		return
	}

	playPauseKey := format.FormatIcon(view.app.Config.Icons.Player, "Nothing is currently playing.")
	toggleShuffleKey := format.FormatIcon(view.app.Config.Icons.ShuffleOn, "Shuffle")
	toggleRepeatKey := format.FormatIcon(view.app.Config.Icons.RepeatContext, "Repeat")

	if player != nil {
		playPauseIcon := view.app.Config.Icons.Play
		if player.IsPlaying {
			playPauseIcon = view.app.Config.Icons.Pause
		}

		playPauseKey = fmt.Sprintf(
			"%s | %s %s | %s | %s/%s",
			playPauseIcon,
			view.app.Config.Icons.Track,
			player.Item.Name,
			player.Item.Artists[0].Name,
			format.FormatTime(player.ProgressMs),
			format.FormatTime(player.Item.DurationMs),
		)

		if player.RepeatState == "off" {
			toggleRepeatKey = format.FormatIcon(view.app.Config.Icons.RepeatOff, "Repeat <u>off</u> context track")
		} else if player.RepeatState == "context" {
			toggleRepeatKey = format.FormatIcon(view.app.Config.Icons.RepeatContext, "Repeat off <u>context</u> track")
		} else if player.RepeatState == "track" {
			toggleRepeatKey = format.FormatIcon(view.app.Config.Icons.RepeatTrack, "Repeat off context <u>track</u>")
		}

		if player.ShuffleState {
			toggleShuffleKey = format.FormatIcon(view.app.Config.Icons.ShuffleOn, "Shuffle <u>true</u> false")
		} else {
			toggleShuffleKey = format.FormatIcon(view.app.Config.Icons.ShuffleOff, "Shuffle true <u>false</u>")
		}
	}

	view.rofi.Rows = []rofi.Row{
		{
			Title: playPauseKey,
			Value: playerTogglePauseAction,
		},
		{
			Title: format.FormatIcon(view.app.Config.Icons.Next, "Next"),
			Value: playerNextAction,
		},
		{
			Title: format.FormatIcon(view.app.Config.Icons.Previous, "Previous"),
			Value: playerPreviousAction,
		},
		{
			Title: toggleShuffleKey,
			Value: playerToggleShuffle,
		},
		{
			Title: toggleRepeatKey,
			Value: playerToggleRepeat,
		},
	}

	selection, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(selection, code)
}

func (view *playerView) SetParent(parent View) {
	view.parent = parent
}

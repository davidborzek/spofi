package views

import (
	"fmt"
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
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

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.SelectedEvent:
		var err error

		switch evt.Selection.Value {
		case playerTogglePauseAction:
			err = view.app.Player.PlayPause()
		case playerNextAction:
			err = view.app.Player.Next()
		case playerPreviousAction:
			err = view.app.Player.Previous()
		case playerToggleShuffle:
			err = view.app.Player.ToggleShuffle()
		case playerToggleRepeat:
			err = view.app.Player.ToggleRepeat()
		}

		if err != nil {
			updatePlayerError(err)
		}

		view.Show()
	}
}

func (view *playerView) SetParent(parent View) {
	view.parent = parent
}

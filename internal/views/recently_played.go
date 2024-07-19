package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

type recentlyPlayedView struct {
	rofi rofi.App
	app  *app.App

	parent View
}

func NewRecentlyPlayedView(app *app.App, title string) View {
	var msg = ""
	if app.Config.ShowKeybindings {
		msg = format.FormatKeybindings(
			format.Keybinding{
				Key:         app.Config.Keybindings.AddToQueue,
				Description: "Add to queue",
			},
		)
	}

	r := rofi.App{
		Prompt: title,
		Keybindings: []string{
			app.Config.Keybindings.AddToQueue,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
		Message:    msg,
	}

	view := &recentlyPlayedView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *recentlyPlayedView) getRecentlyPlayedTracks() ([]rofi.Row, error) {
	result, err := view.app.SpotifyClient.GetRecentlyPlayedTracks()
	if err != nil {
		return nil, err
	}

	tracks := make([]spotify.Track, len(result.Items))
	for i, item := range result.Items {
		tracks[i] = item.Track
	}

	rows := format.FormatTrackRows(
		tracks,
		view.app.Config.Icons.Track,
	)
	return rows, nil
}

func (view *recentlyPlayedView) Show(payload ...interface{}) {
	rows, err := view.getRecentlyPlayedTracks()
	if err != nil {
		getRecentlyPlayedTracksError(err)
		return
	}

	if len(rows) == 0 {
		rofi.Error("No recently played tracks.")
		view.parent.Show()
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
			err := view.app.Player.AddQueue(evt.Selection.Value)
			if err != nil {
				addQueueError(err)
			}
		}

		view.Show()
	case rofi.SelectedEvent:
		err := view.app.Player.PlayTrack(evt.Selection.Value)
		if err != nil {
			playTrackError(err)
		}
	}

}

func (view *recentlyPlayedView) SetParent(parent View) {
	view.parent = parent
}

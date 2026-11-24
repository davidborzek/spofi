package views

import (
	"fmt"
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
)

type likedTracksView struct {
	rofi rofi.App
	app  *app.App

	parent View

	title string
}

func NewLikedTracksView(app *app.App, title string) View {
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

	view := &likedTracksView{
		rofi:  r,
		app:   app,
		title: title,
	}

	return view
}

func (view *likedTracksView) getTracks() ([]rofi.Row, error) {
	tracks, err := view.app.SpotifyClient.GetAllLikedTracks()
	if err != nil {
		return nil, err
	}

	return format.FormatTrackRows(
		tracks,
		view.app.Config.Icons.Track,
	), nil
}

func (view *likedTracksView) addToQueue(uri string) {
	err := view.app.Player.AddQueue(uri)

	if err != nil {
		addQueueError(err)
	}
}

func (view *likedTracksView) Show(payload ...interface{}) {
	rows, err := view.getTracks()
	if err != nil {
		getTracksError(err)
		return
	}

	view.rofi.Prompt = fmt.Sprintf("%s (%d)", view.title, len(rows))
	view.rofi.Rows = rows

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.KeyEvent:
		if evt.Key == view.app.Config.Keybindings.AddToQueue {
			view.addToQueue(evt.Selection.Value)
		}

		view.Show()
	case rofi.SelectedEvent:
		err := view.app.Player.PlayTrack(evt.Selection.Value)
		if err != nil {
			playTrackError(err)
		}
	}
}

func (view *likedTracksView) SetParent(parent View) {
	view.parent = parent
}

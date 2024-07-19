package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

type albumView struct {
	rofi rofi.App
	app  *app.App

	album *spotify.AlbumWithTracks

	parent View
}

func NewAlbumView(app *app.App) View {
	var msg = ""
	if app.Config.ShowKeybindings {
		msg = format.FormatKeybindings(
			format.Keybinding{
				Key:         app.Config.Keybindings.PlayAlbum,
				Description: "Play album",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.AddToQueue,
				Description: "Add to queue",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.PlayTrack,
				Description: "Play track",
			},
		)
	}

	r := rofi.App{
		Keybindings: []string{
			app.Config.Keybindings.PlayAlbum,
			app.Config.Keybindings.AddToQueue,
			app.Config.Keybindings.PlayTrack,
		},
		NoCustom:   true,
		IgnoreCase: true,
		ShowBack:   true,
		Message:    msg,
	}

	view := &albumView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *albumView) addToQueue(uri string) {
	err := view.app.Player.AddQueue(uri)
	if err != nil {
		addQueueError(err)
	}

	view.Show()
}

func (view *albumView) playTrack(uri string) {
	err := view.app.Player.PlayTrack(uri)
	if err != nil {
		playTrackError(err)
	}
}

func (view *albumView) playAlbum(uri ...string) {
	err := view.app.Player.PlayContext(
		view.album.URI,
		uri...,
	)

	if err != nil {
		playAlbumError(err)
		return
	}
}

func (view *albumView) setPrompt() {
	view.rofi.Prompt = format.FormatTitle(view.album.Name, view.album.Artists[0].Name)
}

func (view *albumView) setRows() {
	view.rofi.Rows = format.FormatTrackRows(
		view.album.Tracks.Items,
		view.app.Config.Icons.Track,
	)
}

func (view *albumView) Show(payload ...interface{}) {
	if len(payload) > 0 {
		switch t := payload[0].(type) {
		case spotify.AlbumWithTracks:
			view.album = &t
		case string:
			res, err := view.app.SpotifyClient.GetAlbum(t)
			if err != nil {
				getAlbumError(err)
				return
			}
			view.album = res
		}
	}

	if view.album == nil {
		return
	}

	view.setPrompt()
	view.setRows()

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.KeyEvent:
		switch evt.Key {
		case view.app.Config.Keybindings.PlayAlbum:
			view.playAlbum()
		case view.app.Config.Keybindings.AddToQueue:
			view.addToQueue(evt.Selection.Value)
		case view.app.Config.Keybindings.PlayTrack:
			view.playTrack(evt.Selection.Value)
		}
	case rofi.SelectedEvent:
		view.playAlbum(evt.Selection.Value)
	}
}

func (view *albumView) SetParent(parent View) {
	view.parent = parent
}

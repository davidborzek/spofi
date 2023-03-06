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
	r := rofi.App{
		KBCustom: []string{
			app.Config.Keybindings.PlayAlbum,
			app.Config.Keybindings.AddToQueue,
			app.Config.Keybindings.PlayTrack,
		},
		NoCustom:   true,
		IgnoreCase: true,
		ShowBack:   true,
	}

	view := &albumView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *albumView) addToQueue(uri string) {
	err := view.app.SpotifyClient.AddQueue(
		uri,
		view.app.Config.Device.ID,
	)

	if err != nil {
		addQueueError(err)
	}

	view.Show()
}

func (view *albumView) playTrack(uri string) {
	err := view.app.SpotifyClient.PlayTrack(
		uri,
		view.app.Config.Device.ID,
	)

	if err != nil {
		playTrackError(err)
	}
}

func (view *albumView) playAlbum(uri ...string) {
	err := view.app.SpotifyClient.PlayContext(
		view.album.URI,
		view.app.Config.Device.ID,
		uri...,
	)

	if err != nil {
		playAlbumError(err)
		return
	}
}

func (view *albumView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape || selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	if code == rofi.KBCustom1 {
		view.playAlbum()
		return
	}

	if code == rofi.KBCustom2 {
		view.addToQueue(selection.Value)
		return
	}

	if code == rofi.KBCustom3 {
		view.playTrack(selection.Value)
		return
	}

	if code > 0 {
		return
	}

	view.playAlbum(selection.Value)
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

	result, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(result, code)

}

func (view *albumView) SetParent(parent View) {
	view.parent = parent
}

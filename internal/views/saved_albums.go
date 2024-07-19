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
	savedAlbumsViewLimit = 10
)

type savedAlbumsView struct {
	rofi rofi.App
	app  *app.App

	parent View
	albums *spotify.SavedAlbumResponse

	title      string
	page       int
	totalPages int

	albumView View
}

func NewSavedAlbumsView(app *app.App, title string) View {
	var msg = ""
	if app.Config.ShowKeybindings {
		msg = format.FormatKeybindings(
			format.Keybinding{
				Key:         app.Config.Keybindings.NextPage,
				Description: "Next page",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.PreviousPage,
				Description: "Previous page",
			},
			format.Keybinding{
				Key:         app.Config.Keybindings.PlayAlbum,
				Description: "Play album",
			},
		)
	}

	r := rofi.App{
		Prompt: title,
		Keybindings: []string{
			app.Config.Keybindings.NextPage,
			app.Config.Keybindings.PreviousPage,
			app.Config.Keybindings.PlayAlbum,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
		Message:    msg,
	}

	view := &savedAlbumsView{
		rofi:      r,
		app:       app,
		page:      1,
		title:     title,
		albumView: NewAlbumView(app),
	}

	view.albumView.SetParent(view)

	return view
}

func (view *savedAlbumsView) getAlbums() ([]rofi.Row, error) {
	currentOffset := (view.page - 1) * savedAlbumsViewLimit

	result, err := view.app.SpotifyClient.GetSavedAlbums(savedAlbumsViewLimit, currentOffset)
	if err != nil {
		return nil, err
	}

	view.albums = result

	view.totalPages = result.Total / savedAlbumsViewLimit

	albums := make([]spotify.Album, len(result.Items))
	for i, item := range result.Items {

		albums[i] = item.Album.Album
	}

	rows := format.FormatAlbumRows(
		albums,
		view.app.Config.Icons.Album,
	)
	return rows, nil
}

func (view *savedAlbumsView) Show(payload ...interface{}) {
	rows, err := view.getAlbums()
	if err != nil {
		getAlbumsError(err)
		return
	}

	view.rofi.Prompt = fmt.Sprintf("%s %d/%d", view.title, view.page, view.totalPages)
	view.rofi.Rows = rows

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.page = 1
		view.parent.Show()
	case rofi.KeyEvent:
		switch evt.Key {
		case view.app.Config.Keybindings.NextPage:
			if view.page < view.totalPages {
				view.page += 1
			}

			view.Show()
		case view.app.Config.Keybindings.PreviousPage:
			if view.page > 1 {
				view.page -= 1
			}

			view.Show()
		case view.app.Config.Keybindings.PlayAlbum:
			err := view.app.Player.PlayContext(evt.Selection.Value)
			if err != nil {
				playAlbumError(err)
			}
		}
	case rofi.SelectedEvent:
		for _, a := range view.albums.Items {
			if a.Album.URI == evt.Selection.Value {
				view.albumView.Show(a.Album)
				return
			}
		}
	}
}

func (view *savedAlbumsView) SetParent(parent View) {
	view.parent = parent
}

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
	r := rofi.App{
		Prompt: title,
		KBCustom: []string{
			app.Config.Keybindings.NextPage,
			app.Config.Keybindings.PreviousPage,
			app.Config.Keybindings.PlayAlbum,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
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

func (view *savedAlbumsView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape {
		view.page = 1
		view.parent.Show()
		return
	}

	if code == rofi.KBCustom1 {
		if view.page < view.totalPages {
			view.page += 1
		}

		view.Show()
		return
	}

	if code == rofi.KBCustom2 {
		if view.page > 1 {
			view.page -= 1
		}

		view.Show()
		return
	}

	if code == rofi.KBCustom3 {
		err := view.app.SpotifyClient.PlayContext(
			selection.Value,
			view.app.Config.Device.ID,
		)

		if err != nil {
			playAlbumError(err)
		}

		return
	}

	if code > 0 {
		return
	}

	if selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	for _, a := range view.albums.Items {
		if a.Album.URI == selection.Value {
			view.albumView.Show(a.Album)
			return
		}
	}
}

func (view *savedAlbumsView) Show(payload ...interface{}) {
	rows, err := view.getAlbums()
	if err != nil {
		getAlbumsError(err)
		return
	}

	view.rofi.Prompt = fmt.Sprintf("%s %d/%d", view.title, view.page, view.totalPages)
	view.rofi.Rows = rows

	result, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(result, code)
}

func (view *savedAlbumsView) SetParent(parent View) {
	view.parent = parent
}

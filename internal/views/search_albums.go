package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
	"github.com/davidborzek/spofi/pkg/spotify"
)

type searchAlbumsView struct {
	rofi rofi.App
	app  *app.App

	parent View

	query string
}

func NewSearchAlbumsView(app *app.App) *searchAlbumsView {
	title := format.FormatIcon(
		app.Config.Icons.Album,
		"Albums",
	)

	r := rofi.App{
		Prompt: title,
		KBCustom: []string{
			app.Config.Keybindings.ToggleSearchType,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
	}

	return &searchAlbumsView{
		rofi: r,
		app:  app,
	}
}

func (view *searchAlbumsView) SetParent(parent View) {
	view.parent = parent
}

func (view *searchAlbumsView) SetQuery(query string) {
	view.query = query
}

func (view *searchAlbumsView) Show(payload ...interface{}) {
	if view.query == "" {
		return
	}

	if err := view.search(); err != nil {
		rofi.Error("Could not search albums.")
		log.Println(err)
		return
	}

	selection, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(selection, code)
}

func (view *searchAlbumsView) search() error {
	response, err := view.app.SpotifyClient.Search(view.query, "album")
	if err != nil {
		return err
	}

	view.rofi.Rows = format.FormatAlbumRows(
		response.Albums.Items,
		view.app.Config.Icons.Album,
	)
	return nil
}

func (view *searchAlbumsView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape || selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	if code == rofi.KBCustom1 {
		trackSearch := NewSearchTrackView(view.app)
		trackSearch.SetParent(view.parent)
		trackSearch.SetQuery(view.query)
		trackSearch.Show()
		return
	}

	if code > 0 {
		return
	}

	album := NewAlbumView(view.app)
	album.SetParent(view)
	album.Show(spotify.URIToID(selection.Value))
}

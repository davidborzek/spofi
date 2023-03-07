package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
)

type searchTracksView struct {
	rofi rofi.App
	app  *app.App

	parent View

	query string
}

func NewSearchTrackView(app *app.App) *searchTracksView {
	title := format.FormatIcon(
		app.Config.Icons.Track,
		"Tracks",
	)

	r := rofi.App{
		Prompt: title,
		KBCustom: []string{
			app.Config.Keybindings.AddToQueue,
			app.Config.Keybindings.ToggleSearchType,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
	}

	view := &searchTracksView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *searchTracksView) SetParent(parent View) {
	view.parent = parent
}

func (view *searchTracksView) SetQuery(query string) {
	view.query = query
}

func (view *searchTracksView) Show(payload ...interface{}) {
	if view.query == "" {
		return
	}

	if err := view.search(); err != nil {
		searchError(err)
		return
	}

	selection, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(selection, code)
}

func (view *searchTracksView) search() error {
	response, err := view.app.SpotifyClient.Search(view.query, "track")
	if err != nil {
		return err
	}

	view.rofi.Rows = format.FormatTrackRows(
		response.Tracks.Items,
		view.app.Config.Icons.Track,
	)
	return nil
}

func (view *searchTracksView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape || selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	if code == rofi.KBCustom1 {
		err := view.app.Player.AddQueue(selection.Value)
		if err != nil {
			addQueueError(err)
		}

		view.Show()
		return
	}

	if code == rofi.KBCustom2 {
		albumSearch := NewSearchAlbumsView(view.app)
		albumSearch.SetParent(view.parent)
		albumSearch.SetQuery(view.query)
		albumSearch.Show()
		return
	}

	if code > 0 {
		return
	}

	err := view.app.Player.PlayTrack(selection.Value)
	if err != nil {
		playTrackError(err)
		return
	}
}

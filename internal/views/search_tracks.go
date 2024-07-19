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
		Keybindings: []string{
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
			view.Show()
		case view.app.Config.Keybindings.ToggleSearchType:
			albumSearch := NewSearchAlbumsView(view.app)
			albumSearch.SetParent(view.parent)
			albumSearch.SetQuery(view.query)
			albumSearch.Show()
		}
	case rofi.SelectedEvent:
		err := view.app.Player.PlayTrack(evt.Selection.Value)
		if err != nil {
			playTrackError(err)
		}
	}
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

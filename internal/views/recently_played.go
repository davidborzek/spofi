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
	r := rofi.App{
		Prompt:     title,
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
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

func (view *recentlyPlayedView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape {
		view.parent.Show()
		return
	}

	if code > 0 {
		return
	}

	if selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	view.Show()
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

	result, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(result, code)

}

func (view *recentlyPlayedView) SetParent(parent View) {
	view.parent = parent
}

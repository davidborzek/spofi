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
	likeTracksViewLimit = 10
)

type likedTracksView struct {
	rofi rofi.App
	app  *app.App

	parent View

	title      string
	page       int
	totalPages int
}

func NewLikedTracksView(app *app.App, title string) View {
	r := rofi.App{
		Prompt: title,
		Keybindings: []string{
			app.Config.Keybindings.NextPage,
			app.Config.Keybindings.PreviousPage,
			app.Config.Keybindings.AddToQueue,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
	}

	view := &likedTracksView{
		rofi:  r,
		app:   app,
		page:  1,
		title: title,
	}

	return view
}

func (view *likedTracksView) getTracks() ([]rofi.Row, error) {
	currentOffset := (view.page - 1) * likeTracksViewLimit

	result, err := view.app.SpotifyClient.GetLikedTracks(likeTracksViewLimit, currentOffset)
	if err != nil {
		return nil, err
	}

	view.totalPages = result.Total / likeTracksViewLimit

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
		case view.app.Config.Keybindings.PreviousPage:
			if view.page > 1 {
				view.page -= 1
			}
		case view.app.Config.Keybindings.AddToQueue:
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

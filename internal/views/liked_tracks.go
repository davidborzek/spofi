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
		KBCustom: []string{
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
	err := view.app.SpotifyClient.AddQueue(
		uri,
		view.app.Config.Device.ID,
	)

	if err != nil {
		addQueueError(err)
		log.Println(err)
	}

	view.Show()
}

func (view *likedTracksView) handleSelection(selection *rofi.Row, code int) {
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
		view.addToQueue(selection.Value)
		return
	}

	if code > 0 {
		return
	}

	if selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	err := view.app.SpotifyClient.PlayTrack(
		selection.Value,
		view.app.Config.Device.ID,
	)

	if err != nil {
		playTrackError(err)
		log.Println(err)
		return
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

	result, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(result, code)
}

func (view *likedTracksView) SetParent(parent View) {
	view.parent = parent
}

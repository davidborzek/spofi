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
	playlistsViewLimit = 10
)

type playlistsView struct {
	rofi rofi.App
	app  *app.App

	parent    View
	playlists *spotify.PlaylistsResponse

	title      string
	page       int
	totalPages int
}

func NewPlaylistsView(app *app.App, title string) View {
	r := rofi.App{
		Prompt: title,
		KBCustom: []string{
			app.Config.Keybindings.NextPage,
			app.Config.Keybindings.PreviousPage,
			app.Config.Keybindings.PlayPlaylist,
		},
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
	}

	view := &playlistsView{
		rofi:  r,
		app:   app,
		page:  1,
		title: title,
	}

	return view
}

func (view *playlistsView) getPlaylists() ([]rofi.Row, error) {
	currentOffset := (view.page - 1) * playlistsViewLimit

	result, err := view.app.SpotifyClient.GetUsersPlaylists(playlistsViewLimit, currentOffset)
	if err != nil {
		return nil, err
	}

	view.playlists = result

	view.totalPages = result.Total / playlistsViewLimit

	rows := format.FormatPlaylistRows(
		result.Items,
		view.app.Config.Icons.Playlist,
	)
	return rows, nil
}

func (view *playlistsView) handleSelection(selection *rofi.Row, code int) {
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
		err := view.app.Player.PlayContext(selection.Value)
		if err != nil {
			playPlaylistError(err)
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
}

func (view *playlistsView) Show(payload ...interface{}) {
	rows, err := view.getPlaylists()
	if err != nil {
		getPlaylistsError(err)
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

func (view *playlistsView) SetParent(parent View) {
	view.parent = parent
}

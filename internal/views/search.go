package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/pkg/rofi"
)

type searchView struct {
	rofi rofi.App
	app  *app.App

	parent View

	filter string
}

func NewSearchView(app *app.App, title string) View {
	r := rofi.App{
		Prompt:   title,
		ShowBack: true,
	}

	view := &searchView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *searchView) handleSelection(selection *rofi.Row, code int) {
	if code == rofi.Escape || selection.Title == rofi.Back {
		view.parent.Show()
		return
	}

	if code > 0 {
		return
	}

	if selection.Title == "" {
		rofi.Error("Search cannot be empty.")
		view.Show()
		return
	}

	tracks := NewSearchTrackView(view.app)
	tracks.SetParent(view)
	tracks.SetQuery(selection.Title)
	tracks.Show()
}

func (view *searchView) Show(payload ...interface{}) {
	view.rofi.Filter = view.filter

	selection, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.filter = selection.Title

	view.handleSelection(selection, code)
}

func (view *searchView) SetParent(parent View) {
	view.parent = parent
}

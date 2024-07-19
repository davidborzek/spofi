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

func (view *searchView) Show(payload ...interface{}) {
	view.rofi.Filter = view.filter

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt := evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.SelectedEvent:
		view.filter = evt.Selection.Title

		if evt.Selection.Title == "" {
			rofi.Error("Search cannot be empty.")
			view.Show()
			return
		}

		tracks := NewSearchTrackView(view.app)
		tracks.SetParent(view)
		tracks.SetQuery(evt.Selection.Title)
		tracks.Show()
	}
}

func (view *searchView) SetParent(parent View) {
	view.parent = parent
}

package views

import (
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/format"
	"github.com/davidborzek/spofi/pkg/rofi"
)

type queueView struct {
	rofi rofi.App
	app  *app.App

	parent View
}

func NewQueueView(app *app.App, title string) View {
	r := rofi.App{
		Prompt:     title,
		ShowBack:   true,
		NoCustom:   true,
		IgnoreCase: true,
	}

	view := &queueView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *queueView) getQueue() ([]rofi.Row, error) {
	result, err := view.app.SpotifyClient.GetQueue()
	if err != nil {
		return nil, err
	}

	rows := format.FormatTrackRows(
		result.Queue,
		view.app.Config.Icons.Track,
	)
	return rows, nil
}

func (view *queueView) Show(payload ...interface{}) {
	rows, err := view.getQueue()
	if err != nil {
		getQueueError(err)
		return
	}

	if len(rows) == 0 {
		rofi.Error("Queue is empty.")
		view.parent.Show()
		return
	}

	view.rofi.Rows = rows

	evt, err := view.rofi.Run()
	if err != nil {
		log.Fatalln(err.Error())
	}

	switch evt.(type) {
	case rofi.BackEvent, rofi.CancelledEvent:
		view.parent.Show()
	case rofi.SelectedEvent:
		view.Show()
	}
}

func (view *queueView) SetParent(parent View) {
	view.parent = parent
}

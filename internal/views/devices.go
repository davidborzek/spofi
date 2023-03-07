package views

import (
	"fmt"
	"log"

	"github.com/davidborzek/spofi/internal/app"
	"github.com/davidborzek/spofi/internal/config"
	"github.com/davidborzek/spofi/pkg/rofi"
)

type devicesView struct {
	rofi rofi.App
	app  *app.App

	parent View
}

func NewDevicesView(app *app.App, title string) View {
	r := rofi.App{
		Prompt:     title,
		NoCustom:   true,
		IgnoreCase: true,
		ShowBack:   true,
	}

	view := &devicesView{
		rofi: r,
		app:  app,
	}

	return view
}

func (view *devicesView) getDevices() ([]rofi.Row, error) {
	result, err := view.app.SpotifyClient.GetDevices()
	if err != nil {
		return nil, err
	}

	rows := make([]rofi.Row, len(result.Devices))

	for i, device := range result.Devices {
		rows[i] = rofi.Row{
			Title: device.Name,
			Value: device.ID,
		}
	}

	return rows, nil
}

func (view *devicesView) getCurrentDevice() string {
	var msg string
	if view.app.Config.Device.Name != "" {
		msg = fmt.Sprintf(
			"Current device: %s",
			view.app.Config.Device.Name,
		)
	} else {
		msg = "No device selected"
	}
	return msg
}

func (view *devicesView) handleSelection(selection *rofi.Row, code int) {
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

	view.app.Config.Device = config.SpotifyDevice{
		ID:   selection.Value,
		Name: selection.Title,
	}

	if err := view.app.Config.Write(); err != nil {
		selectDeviceError(err)
		return
	}

	view.app.Player.SetDevice(selection.Value)
	view.Show()
}

func (view *devicesView) Show(payload ...interface{}) {
	rows, err := view.getDevices()
	if err != nil {
		getDevicesError(err)
		return
	}

	if len(rows) == 0 {
		noDevicesFoundError()
		view.parent.Show()
		return
	}

	msg := view.getCurrentDevice()

	view.rofi.Message = msg
	view.rofi.Rows = rows

	result, code, err := view.rofi.Show()
	if err != nil {
		log.Fatalln(err.Error())
	}

	view.handleSelection(result, code)

}

func (view *devicesView) SetParent(parent View) {
	view.parent = parent
}

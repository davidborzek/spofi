package rofi

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

const (
	// The exit status of rofi when a entry was selected.
	statusSelected = 0
	// The exit status of rofi when the selection was cancelled.
	statusCancelled = 1
	// The exit status of rofi when the first custom keybinding was pressed.
	statusKbCustom = 10
)

var (
	customTheme = ""
)

// SetCustomTheme globally sets a custom rofi theme
// for all rofi views.
func SetCustomTheme(theme string) {
	customTheme = theme
}

// Event represents a rofi event.
type Event interface{}

// SelectedEvent appears when the user selected a entry.
type SelectedEvent struct {
	Selection Row
}

// CancelledEvent appears when the user cancels a selection.
type CancelledEvent struct{}

// BackEvent appears when the user selects the back option.
type BackEvent struct{}

// KeyEvent appears whe the user presses a custom keybinding.
type KeyEvent struct {
	Selection Row
	Key       string
}

// Row represents the entries of a rofi menu.
type Row struct {
	Title string
	Value string
}

// App represent a rofi app.
type App struct {
	// Prompt is the prompt of the rofi menu.
	Prompt string
	// Message is the message of the rofi menu.
	Message string
	// Filter is the filter of the rofi menu.
	Filter string
	// IgnoreCase defines the case-insensitivity of the search.
	IgnoreCase bool
	// NoCustom disables custom selection.
	NoCustom bool
	// RenderMarkup enables markup rendering.
	RenderMarkup bool
	// ShowBack shows the back (..) option.
	ShowBack bool
	// Keybindings are the custom keybindings of the rofi menu.
	Keybindings []string
	// Rows are the rows of the rofi menu.
	Rows []Row

	previousSelection int
}

func (a *App) parseArgs() []string {
	args := []string{
		"-dmenu",
	}

	if customTheme != "" {
		args = append(args, "-theme")
		args = append(args, customTheme)
	}

	if a.Prompt != "" {
		args = append(args, "-p")
		args = append(args, a.Prompt)
	}

	if a.Message != "" {
		args = append(args, "-mesg")
		args = append(args, a.Message)
	}

	if a.Filter != "" {
		args = append(args, "-filter")
		args = append(args, a.Filter)
	}

	if a.IgnoreCase {
		args = append(args, "-i")
	}

	if a.NoCustom {
		args = append(args, "-no-custom")
	}

	if a.RenderMarkup {
		args = append(args, "-markup-rows")
	}

	for i, key := range a.Keybindings {
		args = append(args, fmt.Sprintf("-kb-custom-%d", i+1))
		args = append(args, key)
	}

	selected := a.previousSelection
	// Skip back button and select next entry
	// when entries are available.
	if a.ShowBack && len(a.Rows) > 0 {
		selected++
	}

	args = append(args, "-selected-row")
	args = append(args, strconv.Itoa(selected))

	return args
}

// findSelection searches a row in the known app rows.
// If no row was found then a new
// row with only `Title` will be returned.
func (a *App) findSelection(title string) (Row, int) {
	var selection *Row
	index := 0
	for i, entry := range a.Rows {
		if entry.Title == title {
			selection = &entry
			index = i
			break
		}
	}

	if selection == nil {
		selection = &Row{
			Title: title,
		}
	}

	return *selection, index
}

// Run runs the rofi menu and returns a Event.
func (a *App) Run() (Event, error) {
	args := a.parseArgs()

	cmd := exec.Command("rofi", args...)
	buf := bytes.NewBufferString("")

	if a.ShowBack {
		fmt.Fprintln(buf, "..")
	}

	for _, entry := range a.Rows {
		fmt.Fprintln(buf, entry.Title)
	}

	cmd.Stdin = buf
	out, err := cmd.CombinedOutput()

	status := 0
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			status = exiterr.ExitCode()
		} else {
			return nil, err
		}
	}

	key := strings.TrimSpace(string(out))
	selection, index := a.findSelection(key)
	a.previousSelection = index

	if status == statusSelected {
		if a.ShowBack && selection.Title == ".." && index == 0 {
			return BackEvent{}, nil
		}

		return SelectedEvent{
			Selection: selection,
		}, nil
	}

	if status == statusCancelled {
		return CancelledEvent{}, nil
	}

	if status >= statusKbCustom {
		return KeyEvent{
			Selection: selection,
			Key:       a.Keybindings[status-statusKbCustom],
		}, nil
	}

	return nil, fmt.Errorf("received invalid rofi status: %d", status)
}

// Error displays a rofi error view
// with a given message.
func Error(msg string) error {
	args := []string{
		"-e", msg,
	}

	if customTheme != "" {
		args = append(args, "-theme")
		args = append(args, customTheme)
	}

	cmd := exec.Command("rofi", args...)
	_, err := cmd.CombinedOutput()
	return err
}

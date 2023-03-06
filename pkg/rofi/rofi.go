package rofi

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type App struct {
	Prompt       string
	Message      string
	Filter       string
	IgnoreCase   bool
	NoCustom     bool
	ShowBack     bool
	RenderMarkup bool
	KBCustom     []string
	Rows         []Row

	previousSelection int
}

const (
	command = "rofi"

	Escape = 1

	KBCustom1  = 10
	KBCustom2  = 11
	KBCustom3  = 12
	KBCustom4  = 13
	KBCustom5  = 14
	KBCustom6  = 15
	KBCustom7  = 16
	KBCustom8  = 17
	KBCustom9  = 18
	KBCustom10 = 19
	KBCustom11 = 20
	KBCustom12 = 21
	KBCustom13 = 22
	KBCustom14 = 23
	KBCustom15 = 24
	KBCustom16 = 25
	KBCustom17 = 26
	KBCustom18 = 27
	KBCustom19 = 28

	Back = ".."
)

var (
	customTheme = ""
)

// Row represents the entries of a rofi menu.
type Row struct {
	Title string
	Value string
}

// SetCustomTheme globally sets a custom rofi theme
// for all rofi views.
func SetCustomTheme(theme string) {
	customTheme = theme
}

// parseArgs is an internal implementation to
// parse the args from the struct into cli arguments
// for rofi.
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

	for i, key := range a.KBCustom {
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
func (a *App) findSelection(title string) (*Row, int) {
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

	return selection, index
}

// Show a new rofi window with a given configuration
// and returns the result.
func (a *App) Show() (*Row, int, error) {
	args := a.parseArgs()

	cmd := exec.Command(command, args...)
	stdinBuffer := bytes.NewBufferString("")

	if a.ShowBack {
		fmt.Fprintln(stdinBuffer, Back)
	}

	for _, entry := range a.Rows {
		fmt.Fprintln(stdinBuffer, entry.Title)
	}

	cmd.Stdin = stdinBuffer
	out, err := cmd.CombinedOutput()

	exitCode := 0
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			exitCode = exiterr.ExitCode()
		} else {
			return nil, 0, err
		}
	}

	key := strings.TrimSpace(string(out))
	selection, index := a.findSelection(key)
	a.previousSelection = index

	return selection, exitCode, nil
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

	cmd := exec.Command(command, args...)
	_, err := cmd.CombinedOutput()
	return err
}

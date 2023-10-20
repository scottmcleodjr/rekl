package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Level int

const (
	levelInfo  Level = iota // levelInfo is for information messages written to the event view
	levelError              // levelError is for error messages written to the event view
)

// tui is the application's tview terminal UI.
type tui struct {
	eventView  *tview.TextView
	inputField *tview.InputField
	app        *tview.Application
}

// newTUI returns a new tui.
func newTUI() *tui {
	eventView := tview.NewTextView()
	eventView.SetDynamicColors(true).SetBorder(true)
	eventView.Write([]byte(welcomeText)) // No timestamp on the welcome text

	inputField := tview.NewInputField().SetLabel("Input:")
	inputField.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// This helps keep things lined up correctly if you resize the window
		eventView.ScrollToEnd()
		return x, y, width, height
	})

	inputForm := tview.NewForm().AddFormItem(inputField)
	inputForm.SetBorder(true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(eventView, 0, 1, false).
		AddItem(inputForm, 5, 0, true)
	app := tview.NewApplication().SetRoot(flex, true)

	return &tui{
		eventView:  eventView,
		inputField: inputField,
		app:        app,
	}
}

// writeToEventView writes messages to the event view.
// writeToEventView prepends a UTC timestamp and a short line
// prefix indicating if the message is information or an error.
func (t *tui) writeToEventView(level Level, message string) {
	utcTime := time.Now().UTC()
	var linePrefix string
	switch level {
	case levelInfo:
		linePrefix = fmt.Sprintf("%02d:%02d [green::b]>[-::-]", utcTime.Hour(), utcTime.Minute())
	case levelError:
		linePrefix = fmt.Sprintf("%02d:%02d [red::b]> Error:[-::-]", utcTime.Hour(), utcTime.Minute())
	}
	line := fmt.Sprintf("%s %s\n", linePrefix, message)
	t.eventView.Write([]byte(line))
}

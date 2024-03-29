package tui

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Level int

const (
	LevelInfo  Level = iota // LevelInfo is for information messages written to the event view
	LevelError              // LevelError is for error messages written to the event view
)

// TUI is the application's tview terminal UI.
type TUI struct {
	eventView  *tview.TextView
	inputField *tview.InputField
	app        *tview.Application
}

// New returns a new tui.
func New() *TUI {
	eventView := tview.NewTextView()
	eventView.SetDynamicColors(true).SetBorder(true)

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

	return &TUI{
		eventView:  eventView,
		inputField: inputField,
		app:        app,
	}
}

// WriteEvent writes messages to the event view.
// WriteEvent prepends a UTC timestamp and a short line
// prefix indicating if the message is information or an error.
func (t *TUI) WriteEvent(level Level, message string) {
	utcTime := time.Now().UTC()
	var linePrefix string
	switch level {
	case LevelInfo:
		linePrefix = fmt.Sprintf("%02d:%02d [green::b]>[-::-]", utcTime.Hour(), utcTime.Minute())
	case LevelError:
		linePrefix = fmt.Sprintf("%02d:%02d [red::b]> Error:[-::-]", utcTime.Hour(), utcTime.Minute())
	}
	line := fmt.Sprintf("%s %s\n", linePrefix, message)
	t.eventView.Write([]byte(line))
}

// ClearEvents removes all text from the event view.
func (t *TUI) ClearEvents() {
	t.eventView.Clear()
}

// InputText returns the current content of the input field.
func (t *TUI) InputText() string {
	return t.inputField.GetText()
}

// ClearInputText clears the current content of the input field.
func (t *TUI) ClearInputText() {
	t.inputField.SetText("")
}

// SetInputCapture sets the capture function for key events in the input field.
func (t *TUI) SetInputCapture(captureFunc func(capture *tcell.EventKey) *tcell.EventKey) {
	t.inputField.SetInputCapture(captureFunc)
}

// RunApp starts the application and the main event loop.
func (t *TUI) RunApp() error {
	return t.app.Run()
}

// StopApp stops the application, causing RunApp() to return.
func (t *TUI) StopApp() {
	t.app.Stop()
}

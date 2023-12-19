package handler

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/tui"
)

// UserInterface is an interface of the TUI methods used by InputHandler.
type UserInterface interface {
	WriteToEventView(level tui.Level, message string)
	ClearEventView()
	InputFieldText() string
	ClearInputField()
	StopApp()
}

// InputHandler processes user input to the TUI.  InputHandler accepts a Keyer, TUI,
// and Config and returns a function to use as the TUI input field capture function.
func InputHandler(keyer *cwkeyer.Keyer, ui UserInterface, cfg *config.Config) func(*tcell.EventKey) *tcell.EventKey {
	return func(capture *tcell.EventKey) *tcell.EventKey {

		if hotkeyHandler(capture, keyer, ui, cfg) {
			return nil // Don't return the capture for hotkeys
		}

		if commandHandler(capture, keyer, ui, cfg) {
			return capture
		}

		if capture.Key() == tcell.KeyEnter {
			message := strings.ToUpper(ui.InputFieldText())
			err := keyer.QueueMessage(message)
			if err != nil {
				ui.WriteToEventView(tui.LevelError, err.Error())
				return capture
			}
			ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("Sending: %s", message))
			ui.ClearInputField()
		}

		return capture
	}
}

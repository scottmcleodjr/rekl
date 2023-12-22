package handler_test

import (
	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/rekl/tui"
)

var (
	enterKey    = tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
	escKey      = tcell.NewEventKey(tcell.KeyESC, 0, tcell.ModNone)
	upKey       = tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
	downKey     = tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
	message1Key = tcell.NewEventKey(tcell.KeyRune, '!', tcell.ModNone)
	message2Key = tcell.NewEventKey(tcell.KeyRune, '@', tcell.ModNone)
	message3Key = tcell.NewEventKey(tcell.KeyRune, '#', tcell.ModNone)
	message4Key = tcell.NewEventKey(tcell.KeyRune, '$', tcell.ModNone)
	message5Key = tcell.NewEventKey(tcell.KeyRune, '%', tcell.ModNone)
	message6Key = tcell.NewEventKey(tcell.KeyRune, '^', tcell.ModNone)
	message7Key = tcell.NewEventKey(tcell.KeyRune, '&', tcell.ModNone)
	message8Key = tcell.NewEventKey(tcell.KeyRune, '*', tcell.ModNone)
	message9Key = tcell.NewEventKey(tcell.KeyRune, '(', tcell.ModNone)
	message0Key = tcell.NewEventKey(tcell.KeyRune, ')', tcell.ModNone)
)

// testUI is a stub implementation of UserInterface.
type testUI struct {
	events         []string
	inputFieldText string
	stopped        bool
}

func (ui *testUI) WriteEvent(level tui.Level, message string) {
	ui.events = append(ui.events, message)
}

func (ui *testUI) ClearEvents() {
	ui.events = []string{}
}

func (ui *testUI) InputText() string {
	return ui.inputFieldText
}

func (ui *testUI) ClearInputText() {
	ui.inputFieldText = ""
}

func (ui *testUI) StopApp() {
	ui.stopped = true
}

func (ui *testUI) lastEvent() string {
	return ui.events[len(ui.events)-1]
}

// testKey is a no-op implementation of cwkeyer.Key.
type testKey struct{}

func (ts testKey) Down() error {
	return nil
}

func (ts testKey) Up() error {
	return nil
}

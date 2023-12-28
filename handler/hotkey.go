package handler

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/tui"
)

func hotkeyHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui UserInterface, cfg *config.Config) bool {
	switch capture.Key() {
	case tcell.KeyUp:
		incrementSpeed(ui, cfg)
		return true
	case tcell.KeyDown:
		decrementSpeed(ui, cfg)
		return true
	case tcell.KeyESC:
		stopCW(keyer, ui)
		return true
	}

	position := strings.IndexRune(")!@#$%^&*(", capture.Rune())
	if ui.InputText() == "" && position != -1 {
		sendMessage(keyer, ui, cfg, position)
		return true
	}

	return false
}

func incrementSpeed(ui UserInterface, cfg *config.Config) {
	err := cfg.IncrementSpeed()
	if err != nil {
		ui.WriteEvent(tui.LevelError, err.Error())
	}
	ui.WriteEvent(tui.LevelInfo, fmt.Sprintf("The CW speed is %d WPM.", cfg.Speed()))
}

func decrementSpeed(ui UserInterface, cfg *config.Config) {
	err := cfg.DecrementSpeed()
	if err != nil {
		ui.WriteEvent(tui.LevelError, err.Error())
	}
	ui.WriteEvent(tui.LevelInfo, fmt.Sprintf("The CW speed is %d WPM.", cfg.Speed()))
}

func stopCW(keyer *cwkeyer.Keyer, ui UserInterface) {
	keyer.DrainSendQueue()
	ui.WriteEvent(tui.LevelInfo, "All messages stopped.")
}

func sendMessage(keyer *cwkeyer.Keyer, ui UserInterface, cfg *config.Config, position int) {
	message, err := cfg.Message(position)
	if err != nil {
		ui.WriteEvent(tui.LevelError, err.Error())
		return
	}

	err = keyer.QueueMessage(message)
	if err != nil {
		ui.WriteEvent(tui.LevelError, err.Error())
		return
	}
	ui.WriteEvent(tui.LevelInfo, fmt.Sprintf("Sending: %s", message))
}

package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/tui"
)

func commandHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui UserInterface, cfg *config.Config) bool {
	if !(capture.Key() == tcell.KeyEnter && strings.HasPrefix(ui.InputFieldText(), "\\")) {
		return false
	}

	splitInput := strings.SplitN(ui.InputFieldText(), " ", 2)
	command := splitInput[0]
	var commandArg string
	if len(splitInput) > 1 {
		commandArg = splitInput[1]
	}

	switch command {
	case "\\speed":
		handleSpeedCommand(ui, cfg, commandArg)
	case "\\1":
		handleMessageSetCommand(ui, cfg, 1, commandArg)
	case "\\2":
		handleMessageSetCommand(ui, cfg, 2, commandArg)
	case "\\3":
		handleMessageSetCommand(ui, cfg, 3, commandArg)
	case "\\4":
		handleMessageSetCommand(ui, cfg, 4, commandArg)
	case "\\5":
		handleMessageSetCommand(ui, cfg, 5, commandArg)
	case "\\6":
		handleMessageSetCommand(ui, cfg, 6, commandArg)
	case "\\7":
		handleMessageSetCommand(ui, cfg, 7, commandArg)
	case "\\8":
		handleMessageSetCommand(ui, cfg, 8, commandArg)
	case "\\9":
		handleMessageSetCommand(ui, cfg, 9, commandArg)
	case "\\0":
		handleMessageSetCommand(ui, cfg, 0, commandArg)
	case "\\config":
		ui.WriteToEventView(tui.LevelInfo, cfg.String())
		ui.ClearInputField()
	case "\\help":
		ui.WriteToEventView(tui.LevelInfo, config.HelpText)
		ui.ClearInputField()
	case "\\clear":
		ui.ClearEventView()
		ui.ClearInputField()
	case "\\quit":
		ui.StopApp()
	default:
		ui.WriteToEventView(tui.LevelError, "unknown Command")
	}

	return true
}

func handleSpeedCommand(ui UserInterface, cfg *config.Config, arg string) {
	if arg != "" {
		newSpeed, err := strconv.Atoi(arg)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, "unable to parse speed argument")
			return
		}

		err = cfg.SetSpeed(newSpeed)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
		}
	}

	ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("The CW speed is %d WPM.", cfg.Speed()))
	ui.ClearInputField()
}

func handleMessageSetCommand(ui UserInterface, cfg *config.Config, position int, arg string) {
	err := cfg.SetMessage(position, arg)
	if err != nil {
		ui.WriteToEventView(tui.LevelError, err.Error())
		return
	}

	// Fetch it back from config so we get any formatting changes
	// Ignore err because we just set this message, will be nil
	message, _ := cfg.Message(position)
	ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("Saved message %d: %s", position, message))
	ui.ClearInputField()
}

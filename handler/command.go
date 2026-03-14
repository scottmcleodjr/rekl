package handler

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
	if !(capture.Key() == tcell.KeyEnter && strings.HasPrefix(ui.InputText(), "\\")) {
		return false
	}

	splitInput := strings.SplitN(ui.InputText(), " ", 2)
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
	case "\\messages":
		handleMessagesCommand(ui, cfg)
		ui.ClearInputText()
	case "\\help":
		ui.WriteEvent(tui.LevelInfo, config.HelpText)
		ui.ClearInputText()
	case "\\clear":
		ui.ClearEvents()
		ui.ClearInputText()
	case "\\quit":
		ui.StopApp()
	default:
		ui.WriteEvent(tui.LevelError, "unknown Command")
	}

	return true
}

func handleSpeedCommand(ui UserInterface, cfg *config.Config, arg string) {
	if arg != "" {
		newSpeed, err := strconv.Atoi(arg)
		if err != nil {
			ui.WriteEvent(tui.LevelError, "unable to parse speed argument")
			return
		}

		err = cfg.Speed.Set(newSpeed)
		if err != nil {
			ui.WriteEvent(tui.LevelError, err.Error())
		}
	}

	ui.WriteEvent(tui.LevelInfo, fmt.Sprintf("The CW speed is %d WPM.", cfg.Speed.WPM()))
	ui.ClearInputText()
}

func handleMessageSetCommand(ui UserInterface, cfg *config.Config, position int, arg string) {
	message, err := cfg.Messages.SetAt(position, arg)
	if err != nil {
		ui.WriteEvent(tui.LevelError, err.Error())
		return
	}
	ui.WriteEvent(tui.LevelInfo, fmt.Sprintf("Saved message %d: %s", position, message))
	ui.ClearInputText()
}

func handleMessagesCommand(ui UserInterface, cfg *config.Config) {
	var sb strings.Builder
	messages := cfg.Messages.All()
	fmt.Fprint(&sb, "Messages:\n")
	for i := 1; i <= 10; i++ {
		position := i % 10 // Put 0 last like on a keyboard
		fmt.Fprintf(&sb, "    %d: %s", position, messages[position])

		// Add a newline after all but the last line
		if position != 0 {
			fmt.Fprintln(&sb)
		}
	}
	ui.WriteEvent(tui.LevelInfo, sb.String())
}

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/tui"
)

type inputHandler func(*tcell.EventKey, *cwkeyer.Keyer, *tui.TUI, *config.Config) (*tcell.EventKey, bool)

// inputHandlers is a slice of functions in a specific
// order that handles input from the tui InputField.
var inputHandlers = []inputHandler{
	speedIncrementHandler,
	speedDecrementHandler,
	speedSetHandler,
	speedHandler,
	messageSetHandler,
	messageSendHandler,
	configHandler,
	helpHandler,
	quitHandler,
	clearHandler,
	stopHandler,
	unknownCommandHandler, // This needs to be second to last
	cwHandler,             // This needs to be last
}

func speedIncrementHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyUp {
		err := cfg.IncrementSpeed()
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
		}
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.Speed()))
		return capture, true
	}

	return capture, false
}

func speedDecrementHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyDown {
		err := cfg.DecrementSpeed()
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
		}
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.Speed()))
		return capture, true
	}

	return capture, false
}

func speedHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\speed" {
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("The speed is currently %d WPM.", cfg.Speed()))
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func speedSetHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	var speedRegex = regexp.MustCompile(`^\\speed\s(\d+)$`)
	speedMatch := speedRegex.FindStringSubmatch(ui.InputFieldText())
	if capture.Key() == tcell.KeyEnter && speedMatch != nil {
		newSpeed, err := strconv.Atoi(speedMatch[1])
		if err != nil {
			// This should be unreachable with the regex test above
			ui.WriteToEventView(tui.LevelError, "Unable to parse speed input.")
			ui.ClearInputField()
			return capture, true
		}
		err = cfg.SetSpeed(newSpeed)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
		}
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.Speed()))
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func messageSetHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	var messageSetRegex = regexp.MustCompile(`^\\(\d)=(.*)$`)
	messageSetMatch := messageSetRegex.FindStringSubmatch(ui.InputFieldText())
	if capture.Key() == tcell.KeyEnter && messageSetMatch != nil {
		messageNumber, err := strconv.Atoi(messageSetMatch[1])
		if err != nil {
			// This should be unreachable with the regex test above
			ui.WriteToEventView(tui.LevelError, err.Error())
			return capture, true
		}
		message := messageSetMatch[2]
		err = cfg.SetMessage(messageNumber, message)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
			return capture, true
		}

		// Fetch it back from config so we get any formatting changes
		formattedMessage, err := cfg.Message(messageNumber)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
			return capture, true
		}
		ui.WriteToEventView(
			tui.LevelInfo,
			fmt.Sprintf("Saved message %d: %s", messageNumber, formattedMessage),
		)
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func messageSendHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	position := strings.IndexRune(")!@#$%^&*(", capture.Rune())
	// Check the input field text so it won't just start
	// sending if you try to use one of these in a message
	if position != -1 && ui.InputFieldText() == "" {
		message, err := cfg.Message(position)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
			return nil, true
		}

		err = keyer.QueueMessage(message)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
			return nil, true
		}
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("Sending: %s", message))
		return nil, true
	}

	return capture, false
}

func configHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\config" {
		ui.WriteToEventView(tui.LevelInfo, cfg.String())
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func helpHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\help" {
		ui.WriteToEventView(tui.LevelInfo, config.HelpText)
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func quitHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\quit" {
		ui.StopApp()
	}

	return capture, false
}

func clearHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\clear" {
		ui.ClearEventView()
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

func stopHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEsc ||
		(capture.Key() == tcell.KeyEnter && ui.InputFieldText() == "\\stop") {
		keyer.DrainSendQueue()
		ui.WriteToEventView(tui.LevelInfo, "All messages stopped.")
		if ui.InputFieldText() == "\\stop" {
			ui.ClearInputField()
		}
		return capture, true
	}

	return capture, false
}

func unknownCommandHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && strings.HasPrefix(ui.InputFieldText(), "\\") {
		ui.WriteToEventView(tui.LevelError, "unknown Command")
		return capture, true
	}

	return capture, false
}

func cwHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, ui *tui.TUI, cfg *config.Config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter {
		message := strings.ToUpper(ui.InputFieldText())
		err := keyer.QueueMessage(message)
		if err != nil {
			ui.WriteToEventView(tui.LevelError, err.Error())
			return capture, true
		}
		ui.WriteToEventView(tui.LevelInfo, fmt.Sprintf("Sending: %s", message))
		ui.ClearInputField()
		return capture, true
	}

	return capture, false
}

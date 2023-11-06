package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
)

type inputHandler func(*tcell.EventKey, *cwkeyer.Keyer, *tui, *config) (*tcell.EventKey, bool)

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

func speedIncrementHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyUp {
		err := cfg.incrementSpeed()
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
		}
		tui.writeToEventView(levelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
		return capture, true
	}

	return capture, false
}

func speedDecrementHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyDown {
		err := cfg.decrementSpeed()
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
		}
		tui.writeToEventView(levelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
		return capture, true
	}

	return capture, false
}

func speedHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\speed" {
		tui.writeToEventView(levelInfo, fmt.Sprintf("The speed is currently %d WPM.", cfg.speed))
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func speedSetHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	var speedRegex = regexp.MustCompile(`^\\speed\s(\d+)$`)
	speedMatch := speedRegex.FindStringSubmatch(tui.inputField.GetText())
	if capture.Key() == tcell.KeyEnter && speedMatch != nil {
		newSpeed, err := strconv.Atoi(speedMatch[1])
		if err != nil {
			// This should be unreachable with the regex test above
			tui.writeToEventView(levelError, "Unable to parse speed input.")
			tui.inputField.SetText("")
			return capture, true
		}
		err = cfg.setSpeed(newSpeed)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
		}
		tui.writeToEventView(levelInfo, fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func messageSetHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	var messageSetRegex = regexp.MustCompile(`^\\(\d)=(.*)$`)
	messageSetMatch := messageSetRegex.FindStringSubmatch(tui.inputField.GetText())
	if capture.Key() == tcell.KeyEnter && messageSetMatch != nil {
		messageNumber, err := strconv.Atoi(messageSetMatch[1])
		if err != nil {
			// This should be unreachable with the regex test above
			tui.writeToEventView(levelError, err.Error())
			return capture, true
		}
		message := messageSetMatch[2]
		err = cfg.setMessage(messageNumber, message)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
			return capture, true
		}

		// Fetch it back from config so we get any formatting changes
		formattedMessage, err := cfg.message(messageNumber)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
			return capture, true
		}
		tui.writeToEventView(
			levelInfo,
			fmt.Sprintf("Saved message %d: %s", messageNumber, formattedMessage),
		)
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func messageSendHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	position := strings.IndexRune(")!@#$%^&*(", capture.Rune())
	// Check the input field text so it won't just start
	// sending if you try to use one of these in a message
	if position != -1 && tui.inputField.GetText() == "" {
		message, err := cfg.message(position)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
			return nil, true
		}

		err = keyer.QueueMessage(message)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
			return nil, true
		}
		tui.writeToEventView(levelInfo, fmt.Sprintf("Sending: %s", message))
		return nil, true
	}

	return capture, false
}

func configHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\config" {
		tui.writeToEventView(levelInfo, cfg.string())
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func helpHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\help" {
		tui.writeToEventView(levelInfo, helpText)
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func quitHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\quit" {
		tui.app.Stop()
	}

	return capture, false
}

func clearHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\clear" {
		tui.eventView.Clear()
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

func stopHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEsc ||
		(capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\stop") {
		keyer.DrainSendQueue()
		tui.writeToEventView(levelInfo, "All messages stopped.")
		if tui.inputField.GetText() == "\\stop" {
			tui.inputField.SetText("")
		}
		return capture, true
	}

	return capture, false
}

func unknownCommandHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter && strings.HasPrefix(tui.inputField.GetText(), "\\") {
		tui.writeToEventView(levelError, "unknown Command")
		return capture, true
	}

	return capture, false
}

func cwHandler(capture *tcell.EventKey, keyer *cwkeyer.Keyer, tui *tui, cfg *config) (*tcell.EventKey, bool) {

	if capture.Key() == tcell.KeyEnter {
		message := strings.ToUpper(tui.inputField.GetText())
		err := keyer.QueueMessage(message)
		if err != nil {
			tui.writeToEventView(levelError, err.Error())
			return capture, true
		}
		tui.writeToEventView(levelInfo, fmt.Sprintf("Sending: %s", message))
		tui.inputField.SetText("")
		return capture, true
	}

	return capture, false
}

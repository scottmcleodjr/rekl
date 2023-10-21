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
	helpHandler,
	quitHandler,
	clearHandler,
	stopHandler,
	cwHandler,
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

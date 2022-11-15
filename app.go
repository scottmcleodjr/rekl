package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/scottmcleodjr/cwkeyer"
)

const (
	initSpeed   = 18 // Moderate, what I use normally
	minSpeed    = 5  // Very very slow
	maxSpeed    = 50 // Very very fast
	welcomeText = `[::b]Welcome to the K3GDS REKL[::-]

[::i]Written by Scott K3GDS, 2022
v0.1.0[::-]

Enter "\help" for a list of supported commands.

`
	helpText = `
A command should be entered as input with no additional text on the line.
A hotkey can be used at any time without submitting the input field.
Any other inputs will be sent as CW if all characters are sendable.
    "\help"       COMMAND    Display this help text
    "\quit"       COMMAND    Exit the program
    "\clear"      COMMAND    Clear the display
    "\speed"      COMMAND    Display the current WPM speed
    "\speed N"    COMMAND    Set the CW speed to N WPM
    [Up Arrow[]    HOTKEY     Increment the CW speed by 1 WPM   
    [Down Arrow[]  HOTKEY     Decrement the CW speed by 1 WPM
    "\stop"       COMMAND    Stop sending CW immediately
    [ESC[]         HOTKEY     Stop sending CW immediately
`
)

// speedRegex identifies input that is used to set the CW speed.
var speedRegex = regexp.MustCompile(`^\\speed\s(\d+)$`)

// config is a cwkeyer.SpeedProvider for the cwkeyer.Keyer.
type config struct {
	speed int
}

func (cfg *config) Speed() int {
	return cfg.speed
}

// tui is the terminal ui tree and contains fields
// for components we need access to elsewhere.
type tui struct {
	eventView  *tview.TextView
	inputField *tview.InputField
	app        *tview.Application
}

// newTUI returns a tui with all the fields correctly initialized.
func newTUI() *tui {
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

	return &tui{
		eventView:  eventView,
		inputField: inputField,
		app:        app,
	}
}

// writeStringToEventView writes a string to the eventView.
// writeStringToEventView prepends a UTC timestamp and a green arrow
// and appends a newline.
func (t *tui) writeStringToEventView(s string) {
	utcTime := time.Now().UTC()
	line := fmt.Sprintf("%02d:%02d [green::b]>[-::-] %s\n", utcTime.Hour(), utcTime.Minute(), s)
	t.eventView.Write([]byte(line))
}

// inputHandler contains the logic for parsing input into the inputField.
// It will be set as the inputCapture for the inputField.
func inputHandler(capture *tcell.EventKey, keyer cwkeyer.Keyer, tui *tui, cfg *config) *tcell.EventKey {

	/* SPEED CONTROLS */

	if capture.Key() == tcell.KeyUp {
		if cfg.speed < maxSpeed {
			cfg.speed++
			tui.writeStringToEventView(fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
			return capture
		}
		tui.writeStringToEventView(fmt.Sprintf("The speed is already at the maximum of %d WPM.", maxSpeed))
		return capture
	}

	if capture.Key() == tcell.KeyDown {
		if cfg.speed > minSpeed {
			cfg.speed--
			tui.writeStringToEventView(fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
			return capture
		}
		tui.writeStringToEventView(fmt.Sprintf("The speed is already at the minimum of %d WPM.", minSpeed))
		return capture
	}

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\speed" {
		tui.writeStringToEventView(fmt.Sprintf("The speed is currently %d WPM.", cfg.speed))
		tui.inputField.SetText("")
		return capture
	}

	speedMatch := speedRegex.FindStringSubmatch(tui.inputField.GetText())
	if capture.Key() == tcell.KeyEnter && speedMatch != nil {
		newSpeed, err := strconv.Atoi(speedMatch[1])
		if err != nil {
			tui.writeStringToEventView("[red::]ERROR:[-::] Unable to parse speed input.")
			tui.inputField.SetText("")
			return capture
		}
		if newSpeed < minSpeed || newSpeed > maxSpeed {
			tui.writeStringToEventView("New speed is not in acceptable range.")
			tui.inputField.SetText("")
			return capture
		}
		cfg.speed = newSpeed
		tui.writeStringToEventView(fmt.Sprintf("The speed is now %d WPM.", cfg.speed))
		tui.inputField.SetText("")
		return capture
	}

	/* USAGE CONTROLS */

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\help" {
		tui.writeStringToEventView(helpText)
		tui.inputField.SetText("")
		return capture
	}

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\quit" {
		tui.app.Stop()
	}

	if capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\clear" {
		tui.eventView.Clear()
		tui.inputField.SetText("")
		return capture
	}

	/* SEND AND STOP */

	if capture.Key() == tcell.KeyEsc ||
		(capture.Key() == tcell.KeyEnter && tui.inputField.GetText() == "\\stop") {
		keyer.DrainSendQueue()
		tui.writeStringToEventView("All messages stopped.")
		if tui.inputField.GetText() == "\\stop" {
			tui.inputField.SetText("")
		}
		return capture
	}

	if capture.Key() == tcell.KeyEnter {
		message := strings.ToUpper(tui.inputField.GetText())
		err := keyer.QueueMessage(message)
		if err != nil {
			tui.writeStringToEventView(fmt.Sprintf("[red::]ERROR:[-::] %s", err.Error()))
			return capture
		}
		tui.writeStringToEventView(fmt.Sprintf("[orange::]SENDING:[-::] %s", message))
		tui.inputField.SetText("")
		return capture
	}

	// Could do input sanitization here?
	return capture
}

func main() {

	beep := flag.Bool("beep", false, "If the REKL should use a Beep Key instead of a Serial DTR Key")
	portName := flag.String("port", "/tty/USB0", "Serial port name for Serial DTR Key")
	flag.Parse()

	var key cwkeyer.Key
	var err error
	if *beep {
		key, err = cwkeyer.NewBeepKey(700, 48000, 1200) // Using the suggested values
	} else {
		key, err = cwkeyer.NewSerialDTRKey(*portName, 115200) // Using the suggested value
	}
	if err != nil {
		log.Fatal(fmt.Sprintf("unable to create key: %s", err))
	}

	cfg := config{speed: initSpeed}
	keyer := cwkeyer.New(&cfg, key)

	tui := newTUI()
	tui.eventView.Write([]byte(welcomeText)) // No timestamp on the welcome text
	tui.inputField.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {
		return inputHandler(capture, keyer, tui, &cfg)
	})

	go func() {
		for {
			err := keyer.ProcessSendQueue(false)
			if err != nil {
				tui.writeStringToEventView(fmt.Sprintf("[red::]ERROR:[-::] %s", err.Error()))
			}
		}
	}()

	err = tui.app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

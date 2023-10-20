package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
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

// config is a cwkeyer.SpeedProvider for the cwkeyer.Keyer.
type config struct {
	speed int
}

// newConfig returns a new Config.
func newConfig() *config {
	return &config{speed: initSpeed}
}

// Speed returns the current CW WPM speed.  Speed is
// exported for the cwkeyer.SpeedProvider interface.
func (cfg *config) Speed() int {
	return cfg.speed
}

// setSpeed sets the current CW WPM speed.
func (cfg *config) setSpeed(speed int) error {
	if speed < minSpeed {
		return fmt.Errorf("new speed is below minimum of %d", minSpeed)
	}
	if speed > maxSpeed {
		return fmt.Errorf("new speed is above maximum of %d", maxSpeed)
	}
	cfg.speed = speed
	return nil
}

// incrementSpeed raises the current CW WPM speed by one.
func (cfg *config) incrementSpeed() error {
	return cfg.setSpeed(cfg.Speed() + 1)
}

// decrementSpeed lowers the current CW WPM speed by one.
func (cfg *config) decrementSpeed() error {
	return cfg.setSpeed(cfg.Speed() - 1)
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
		log.Fatalf("unable to create key: %s", err)
	}

	rekl := newConfig()
	keyer := cwkeyer.New(rekl, key)
	tui := newTUI()

	tui.inputField.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {
		for _, handler := range inputHandlers {
			_, fired := handler(capture, &keyer, tui, rekl)
			if fired {
				return capture
			}
		}
		return capture
	})

	go func() {
		for {
			err := keyer.ProcessSendQueue(false)
			if err != nil {
				tui.writeToEventView(levelError, err.Error())
			}
		}
	}()

	err = tui.app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

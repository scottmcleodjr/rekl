package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
)

const (
	initSpeed   = 18 // Moderate, what I use normally
	minSpeed    = 5  // Very very slow
	maxSpeed    = 50 // Very very fast
	welcomeText = `[::b]Welcome to the K3GDS REKL[::-]

[::i]Written by Scott K3GDS
v0.1.3[::-]

Enter "\help" for a list of supported commands.

`
	helpText = `
A command should be entered as input with no additional text on the line.
A hotkey can be used at any time without submitting the input field.
Any other inputs will be sent as CW if all characters are sendable.

    "\help"       COMMAND    Display this help text
    "\quit"       COMMAND    Exit the program
    "\clear"      COMMAND    Clear the display
	"\config"     COMMAND    Display the current REKL configurations
    "\speed"      COMMAND    Display the current WPM speed
    "\speed N"    COMMAND    Set the CW speed to N WPM
    [Up Arrow]    HOTKEY     Increment the CW speed by 1 WPM
    [Down Arrow]  HOTKEY     Decrement the CW speed by 1 WPM
	"\N=..."      COMMAND    Save a message at memory position N
	[Shift+N]     HOTKEY     Send the message at memory position N
    "\stop"       COMMAND    Stop sending CW immediately
    [ESC[]         HOTKEY     Stop sending CW immediately
`
)

// config is a cwkeyer.SpeedProvider for the cwkeyer.Keyer.
type config struct {
	speed    int
	messages [10]string
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

// message returns the message at position N or an empty
// string if that message is not set.
func (cfg *config) message(n int) (string, error) {
	if n < 0 || n > 9 {
		return "", errors.New("message number out of range")
	}
	return cfg.messages[n], nil
}

// setMessage sets the message at position N to the string
// message argument.
func (cfg *config) setMessage(n int, message string) error {
	if n < 0 || n > 9 {
		return errors.New("message number out of range")
	}
	message = strings.ToUpper(strings.TrimSpace(message))
	for _, r := range message {
		if !cwkeyer.IsKeyable(r) {
			return fmt.Errorf("message contains unsupported rune %c", r)
		}
	}
	cfg.messages[n] = message
	return nil
}

// string returns the current configurations as a multiline string.
func (cfg *config) string() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nSpeed: %d WPM\n", cfg.speed))
	sb.WriteString("Messages:\n")
	for i := 1; i <= 10; i++ {
		position := i % 10                  // Put 0 last like on a keyboard
		message, _ := cfg.message(position) // Error is not reachable here
		sb.WriteString(fmt.Sprintf("    %d: %s\n", position, message))
	}
	return sb.String()
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

	cfg := newConfig()
	keyer := cwkeyer.New(cfg, key)
	tui := newTUI()

	tui.inputField.SetInputCapture(func(capture *tcell.EventKey) *tcell.EventKey {
		for _, handler := range inputHandlers {
			captureOut, fired := handler(capture, &keyer, tui, cfg)
			if fired {
				return captureOut
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

package config

import (
	"fmt"
	"strings"
)

const (
	WelcomeText = `[::b]Welcome to the K3GDS REKL[::-]

[::i]Written by Scott K3GDS
v0.2.2[::-]

Enter "\help" for a list of supported commands.

`
	HelpText = `
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
	"\N ..."      COMMAND    Save a message at memory position N
	[Shift+N]     HOTKEY     Send the message at memory position N
    [ESC[]         HOTKEY     Stop sending CW immediately
`
)

// Config holds current configuration state for the REKL application.
type Config struct {
	Speed    *Speed
	Messages *Messages
}

// New returns a new Config.
func New() *Config {
	return &Config{
		Speed:    NewSpeed(),
		Messages: NewMessages(),
	}
}

// String returns the current configuration as a multiline String.
func (cfg *Config) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nSpeed: %d WPM\n", cfg.Speed.WPM()))
	sb.WriteString("Messages:\n")
	for i := 1; i <= 10; i++ {
		position := i % 10                      // Put 0 last like on a keyboard
		message, _ := cfg.Messages.At(position) // Error is not reachable here
		sb.WriteString(fmt.Sprintf("    %d: %s\n", position, message))
	}
	return sb.String()
}

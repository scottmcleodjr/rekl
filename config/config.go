package config

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
    "\speed"      COMMAND    Display the current WPM speed
    "\speed N"    COMMAND    Set the CW speed to N WPM
    [Up Arrow]    HOTKEY     Increment the CW speed by 1 WPM
    [Down Arrow]  HOTKEY     Decrement the CW speed by 1 WPM
	"\N ..."      COMMAND    Save a message at memory position N
	"\messages"   COMMAND    Display the current saved messages
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

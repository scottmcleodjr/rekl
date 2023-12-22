package config_test

import (
	"testing"

	"github.com/scottmcleodjr/rekl/config"
)

func TestSetSpeed(t *testing.T) {
	tests := []struct {
		input       int
		speedWanted int
		errorWanted bool
	}{
		{input: 15, speedWanted: 15, errorWanted: false},
		{input: config.MinSpeed, speedWanted: config.MinSpeed, errorWanted: false},
		{input: config.MaxSpeed, speedWanted: config.MaxSpeed, errorWanted: false},
		{input: 1, speedWanted: config.InitSpeed, errorWanted: true},
		{input: 150, speedWanted: config.InitSpeed, errorWanted: true},
	}

	for _, test := range tests {
		cfg := config.New()
		err := cfg.SetSpeed(test.input)

		speed := cfg.Speed()
		if speed != test.speedWanted {
			t.Errorf("got %d, want %d for speed", speed, test.speedWanted)
		}
		if test.errorWanted && (err == nil) {
			t.Error("got nil, want error after setting speed")
		}
		if !test.errorWanted && (err != nil) {
			t.Error("got error, want nil after setting speed")
		}
	}
}

func TestIncrementSpeed(t *testing.T) {
	cfg := config.New()
	expectedSpeed := config.InitSpeed

	for {
		err := cfg.IncrementSpeed()
		expectedSpeed++
		if err != nil {
			break
		}

		speed := cfg.Speed()
		if speed != expectedSpeed {
			t.Errorf("got %d, want %d for speed after IncrementSpeed", speed, expectedSpeed)
		}
	}

	// We broke from loop: err was not nil, speed should be at Max
	speed := cfg.Speed()
	if speed != config.MaxSpeed {
		t.Errorf("got %d, want %d for speed after IncrementSpeed returned error", speed, config.MaxSpeed)
	}
}

func TestDecrementSpeed(t *testing.T) {
	cfg := config.New()
	expectedSpeed := config.InitSpeed

	for {
		err := cfg.DecrementSpeed()
		expectedSpeed--
		if err != nil {
			break
		}

		speed := cfg.Speed()
		if speed != expectedSpeed {
			t.Errorf("got %d, want %d for speed after DecrementSpeed", speed, expectedSpeed)
		}
	}

	// We broke from loop: err was not nil, speed should be at Min
	speed := cfg.Speed()
	if speed != config.MinSpeed {
		t.Errorf("got %d, want %d for speed after DecrementSpeed returned error", speed, config.MinSpeed)
	}
}

func TestMessage(t *testing.T) {
	// Because TestSetMessage already covers this code well,
	// just checking for errors on out of range positions here.
	tests := []struct {
		position    int
		errorWanted bool
	}{
		{position: 0, errorWanted: false},
		{position: 5, errorWanted: false},
		{position: 9, errorWanted: false},
		{position: -1, errorWanted: true},
		{position: 10, errorWanted: true},
		{position: 42, errorWanted: true},
	}

	for _, test := range tests {
		cfg := config.New()
		_, err := cfg.Message(test.position)
		if test.errorWanted && (err == nil) {
			t.Errorf("got nil, want error getting message %d", test.position)
		}
		if !test.errorWanted && (err != nil) {
			t.Errorf("got error, want nil getting message %d", test.position)
		}
	}
}

func TestSetMessage(t *testing.T) {
	inputs := []struct {
		position    int
		message     string
		errorWanted bool
	}{
		{position: 0, message: "CQ CQ K3GDS", errorWanted: false}, // Valid message
		{position: 1, message: "   5NN TU  ", errorWanted: false}, // Valid with spaces
		{position: 3, message: "lower case?", errorWanted: false}, // Valid with lowercase
		{position: 3, message: " newmessage", errorWanted: false}, // Overwrite position 3
		{position: 6, message: "invalid $  ", errorWanted: true},  // Invalid char
		{position: 7, message: "invalid %  ", errorWanted: true},  // Invalid char
		{position: -1, message: "badposition", errorWanted: true}, // Bad position
		{position: 10, message: "badposition", errorWanted: true}, // Bad position
	}

	cfg := config.New()

	// Submit all the test messages and verify returned error
	for _, input := range inputs {
		err := cfg.SetMessage(input.position, input.message)
		if input.errorWanted && (err == nil) {
			t.Errorf("got nil, want error after setting message %d", input.position)
		}
		if !input.errorWanted && (err != nil) {
			t.Errorf("got error, want nil after setting message %d", input.position)
		}
	}

	// Check all the message positions for the expected value
	expectedMessages := map[int]string{}
	expectedMessages[0] = "CQ CQ K3GDS"
	expectedMessages[1] = "5NN TU"
	expectedMessages[3] = "NEWMESSAGE"

	for position := 0; position < 10; position++ {
		got, err := cfg.Message(position)
		want := expectedMessages[position]
		if got != want {
			t.Errorf("got %q, want %q after getting test messages", got, want)
		}
		if err != nil {
			t.Errorf("got error, want nil getting message %d", position)
		}
	}
}

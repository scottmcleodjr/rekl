package config_test

import (
	"testing"

	"github.com/scottmcleodjr/rekl/config"
)

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

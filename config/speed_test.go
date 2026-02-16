package config_test

import (
	"testing"

	"github.com/scottmcleodjr/rekl/config"
)

func TestSet(t *testing.T) {
	tests := []struct {
		input       int
		speedWanted int
		errorWanted bool
	}{
		{input: 15, speedWanted: 15, errorWanted: false},
		{input: config.MinWPM, speedWanted: config.MinWPM, errorWanted: false},
		{input: config.MaxWPM, speedWanted: config.MaxWPM, errorWanted: false},
		{input: 1, speedWanted: config.InitialWPM, errorWanted: true},
		{input: 150, speedWanted: config.InitialWPM, errorWanted: true},
	}

	for _, test := range tests {
		s := config.NewSpeed()
		err := s.Set(test.input)

		speed := s.WPM()
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

func TestIncrement(t *testing.T) {
	s := config.NewSpeed()
	expectedSpeed := config.InitialWPM

	for {
		err := s.Increment()
		expectedSpeed++
		if err != nil {
			break
		}

		speed := s.WPM()
		if speed != expectedSpeed {
			t.Errorf("got %d, want %d for speed after IncrementSpeed", speed, expectedSpeed)
		}
	}

	// We broke from loop: err was not nil, speed should be at Max
	speed := s.WPM()
	if speed != config.MaxWPM {
		t.Errorf("got %d, want %d for speed after IncrementSpeed returned error", speed, config.MaxWPM)
	}
}

func TestDecrement(t *testing.T) {
	s := config.NewSpeed()
	expectedSpeed := config.InitialWPM

	for {
		err := s.Decrement()
		expectedSpeed--
		if err != nil {
			break
		}

		speed := s.WPM()
		if speed != expectedSpeed {
			t.Errorf("got %d, want %d for speed after DecrementSpeed", speed, expectedSpeed)
		}
	}

	// We broke from loop: err was not nil, speed should be at Min
	speed := s.WPM()
	if speed != config.MinWPM {
		t.Errorf("got %d, want %d for speed after DecrementSpeed returned error", speed, config.MinWPM)
	}
}

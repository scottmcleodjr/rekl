package config

import (
	"fmt"
	"sync"
)

const (
	InitialWPM = 18 // Moderate, what I use normally
	MinWPM     = 5  // Very very slow
	MaxWPM     = 50 // Very very fast
)

// Speed holds current CW WPM.
// Speed is also the cwkeyer.SpeedProvider for the cwkeyer.Keyer.
type Speed struct {
	mu  sync.Mutex
	wpm int
}

// NewSpeed returns a new Speed.
func NewSpeed() *Speed {
	return &Speed{wpm: InitialWPM}
}

// WPM returns the current CW WPM.
func (s *Speed) WPM() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.wpm
}

// Speed wraps WPM for the cwkeyer.SpeedProvider interface.
func (s *Speed) Speed() int {
	return s.WPM()
}

// Set sets the current CW WPM.
func (s *Speed) Set(wpm int) error {
	if wpm < MinWPM {
		return fmt.Errorf("new speed is below minimum of %d", MinWPM)
	}
	if wpm > MaxWPM {
		return fmt.Errorf("new speed is above maximum of %d", MaxWPM)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wpm = wpm
	return nil
}

// Increment raises the current CW WPM by one.
func (s *Speed) Increment() error {
	return s.Set(s.WPM() + 1)
}

// Decrement lowers the current CW WPM by one.
func (s *Speed) Decrement() error {
	return s.Set(s.WPM() - 1)
}

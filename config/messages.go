package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/scottmcleodjr/cwkeyer"
)

const Capacity = 10 // Number of message slots (mapped to keyboard 0-9)

// Messages manages stored CW messages in numbered slots (0-9).
type Messages struct {
	mu    sync.Mutex
	slots [Capacity]string
}

// NewMessages returns a new Messages
func NewMessages() *Messages {
	return &Messages{}
}

// At returns the message at the given slot position (0-9).
// If an error occurs, At returns an empty string and the error.
func (m *Messages) At(position int) (string, error) {
	if err := validatePosition(position); err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.slots[position], nil
}

// SetAt stores a message at the given slot position (0-9).
// SetAt returns the formatted message. If an error occurs,
// SetAt returns an empty string and the error.
func (m *Messages) SetAt(position int, message string) (string, error) {
	if err := validatePosition(position); err != nil {
		return "", err
	}

	// Format and validate
	message = strings.ToUpper(strings.TrimSpace(message))
	for _, r := range message {
		if !cwkeyer.IsKeyable(r) {
			return "", fmt.Errorf("message contains unkeyable character: %c", r)
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.slots[position] = message
	return message, nil
}

// All returns a copy of all message slots for serialization.
func (m *Messages) All() [Capacity]string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.slots
}

func validatePosition(pos int) error {
	if pos < 0 || pos >= Capacity {
		return errors.New("message position out of range [0-9]")
	}
	return nil
}

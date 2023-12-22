package handler_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/handler"
)

func TestSpeedIncrement(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	// Increase speed from initial speed
	inputHandler(upKey)
	if cfg.Speed() != config.InitSpeed+1 {
		t.Errorf("got %d, want %d when incrementing from InitSpeed", cfg.Speed(), config.InitSpeed+1)
	}

	// Does not increase from limit
	cfg.SetSpeed(config.MaxSpeed)
	inputHandler(upKey)
	if cfg.Speed() != config.MaxSpeed {
		t.Errorf("got %d, want %d when incrementing from MaxSpeed", cfg.Speed(), config.MaxSpeed)
	}
}

func TestSpeedDecrement(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	// Decrease speed from initial speed
	inputHandler(downKey)
	if cfg.Speed() != config.InitSpeed-1 {
		t.Errorf("got %d, want %d when decrementing from InitSpeed", cfg.Speed(), config.InitSpeed-1)
	}

	// Does not decrease from limit
	cfg.SetSpeed(config.MinSpeed)
	inputHandler(downKey)
	if cfg.Speed() != config.MinSpeed {
		t.Errorf("got %d, want %d when decrementing from MaxSpeed", cfg.Speed(), config.MinSpeed)
	}
}

func TestStopCW(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	keyer.QueueMessage("Test message")
	inputHandler(escKey)
	if !keyer.SendQueueIsEmpty() {
		t.Error("send queue not empty after stop hotkey")
	}

}

func TestMessageSend(t *testing.T) {
	tests := []struct {
		position int
		key      *tcell.EventKey
	}{
		{position: 1, key: message1Key},
		{position: 2, key: message2Key},
		{position: 3, key: message3Key},
		{position: 4, key: message4Key},
		{position: 5, key: message5Key},
		{position: 6, key: message6Key},
		{position: 7, key: message7Key},
		{position: 8, key: message8Key},
		{position: 9, key: message9Key},
		{position: 0, key: message0Key},
	}

	for _, test := range tests {
		cfg := config.New()
		keyer := cwkeyer.New(cfg, testKey{})
		ui := &testUI{}
		inputHandler := handler.InputHandler(keyer, ui, cfg)

		cfg.SetMessage(test.position, "Test message")
		inputHandler(test.key)
		if keyer.SendQueueIsEmpty() {
			t.Errorf("send queue empty after sending message at position %d", test.position)
		}
	}
}

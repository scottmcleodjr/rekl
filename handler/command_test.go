package handler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/handler"
	"github.com/scottmcleodjr/rekl/tui"
)

func TestSpeedCommand(t *testing.T) {
	tests := []struct {
		input int
		want  int
	}{
		{input: 30, want: 30},                // Valid speed
		{input: 1, want: config.InitialWPM},  // Too low
		{input: 99, want: config.InitialWPM}, // Too high
	}

	for _, test := range tests {
		cfg := config.New()
		keyer := cwkeyer.New(cfg.Speed, testKey{})
		ui := &testUI{}
		inputHandler := handler.InputHandler(keyer, ui, cfg)

		ui.inputFieldText = fmt.Sprintf("\\speed %d", test.input)
		inputHandler(enterKey)
		if cfg.Speed.WPM() != test.want {
			t.Errorf("got %d, want %d after setting speed to %d", cfg.Speed.WPM(), test.want, test.input)
		}
	}
}

func TestMessageSet(t *testing.T) {
	tests := []struct {
		input        string
		wantPosition int
		wantMessage  string
	}{
		{input: "\\1 cq de k3gds", wantPosition: 1, wantMessage: "CQ DE K3GDS"}, // Check formatting
		{input: "\\2 CQ CQ %$#", wantPosition: 2, wantMessage: ""},              // Check invalid chars
		{input: "\\3 ", wantPosition: 3, wantMessage: ""},                       // Empty with space
		{input: "\\4", wantPosition: 4, wantMessage: ""},                        // Empty without space
		{input: "\\5 TEST MESSAGE", wantPosition: 5, wantMessage: "TEST MESSAGE"},
		{input: "\\6 TEST MESSAGE", wantPosition: 6, wantMessage: "TEST MESSAGE"},
		{input: "\\7 TEST MESSAGE", wantPosition: 7, wantMessage: "TEST MESSAGE"},
		{input: "\\8 TEST MESSAGE", wantPosition: 8, wantMessage: "TEST MESSAGE"},
		{input: "\\9 TEST MESSAGE", wantPosition: 9, wantMessage: "TEST MESSAGE"},
		{input: "\\0 TEST MESSAGE", wantPosition: 0, wantMessage: "TEST MESSAGE"},
	}

	for _, test := range tests {
		cfg := config.New()
		keyer := cwkeyer.New(cfg.Speed, testKey{})
		ui := &testUI{}
		inputHandler := handler.InputHandler(keyer, ui, cfg)

		ui.inputFieldText = test.input
		inputHandler(enterKey)
		message, _ := cfg.Messages.At(test.wantPosition)
		if message != test.wantMessage {
			t.Errorf("got %q, want %q for input %q", message, test.wantMessage, test.input)
		}
	}
}

func TestMessages(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg.Speed, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	cfg.Messages.SetAt(5, "Message Text")

	ui.inputFieldText = "\\messages"
	inputHandler(enterKey)

	got := ui.lastEvent()
	wantContains := []string{
		"Messages:",
		"5: MESSAGE TEXT",
	}

	for _, want := range wantContains {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func TestHelp(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg.Speed, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	ui.inputFieldText = "\\help"
	inputHandler(enterKey)
	lastEvent := ui.lastEvent()
	if lastEvent != config.HelpText {
		t.Errorf("got event %q, want %q for help", lastEvent, config.HelpText)
	}
}

func TestClear(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg.Speed, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	ui.WriteEvent(tui.LevelInfo, "First event")
	ui.WriteEvent(tui.LevelInfo, "Second event")

	ui.inputFieldText = "\\clear"
	inputHandler(enterKey)
	if len(ui.events) != 0 {
		t.Error("ui events not empty after clear command")
	}
}

func TestQuit(t *testing.T) {
	cfg := config.New()
	keyer := cwkeyer.New(cfg.Speed, testKey{})
	ui := &testUI{}
	inputHandler := handler.InputHandler(keyer, ui, cfg)

	ui.inputFieldText = "\\quit"
	inputHandler(enterKey)
	if !ui.stopped {
		t.Error("ui not stopped after quit command")
	}
}

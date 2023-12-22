package handler_test

import (
	"testing"

	"github.com/scottmcleodjr/cwkeyer"
	"github.com/scottmcleodjr/rekl/config"
	"github.com/scottmcleodjr/rekl/handler"
)

func TestCWSend(t *testing.T) {
	tests := []struct {
		input  string
		wantCW bool
	}{
		{input: "CQ CQ CQ DE K3GDS", wantCW: true}, // Valid for CW
		{input: "5NN @ TU", wantCW: false},         // Invalid character
	}

	for _, test := range tests {
		cfg := config.New()
		keyer := cwkeyer.New(cfg, testKey{})
		ui := &testUI{}
		inputHandler := handler.InputHandler(keyer, ui, cfg)

		ui.inputFieldText = test.input
		inputHandler(enterKey)
		if keyer.SendQueueIsEmpty() && test.wantCW {
			t.Errorf("got empty keyer send queue, want items in queue for input %q", test.input)
		}
		if !keyer.SendQueueIsEmpty() && !test.wantCW {
			t.Errorf("got items in keyer send queue, want empty queue for input %q", test.input)
		}
	}
}

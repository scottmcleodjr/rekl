package config_test

import (
	"testing"

	"github.com/scottmcleodjr/rekl/config"
)

func TestMessages_At(t *testing.T) {
	tests := []struct {
		name     string
		position int
		setup    func(*config.Messages)
		want     string
		wantErr  bool
	}{
		{
			name:     "valid empty slot",
			position: 0,
			want:     "",
			wantErr:  false,
		},
		{
			name:     "valid slot with message",
			position: 5,
			setup: func(m *config.Messages) {
				m.SetAt(5, "CQ CQ TEST")
			},
			want:    "CQ CQ TEST",
			wantErr: false,
		},
		{
			name:     "position too low",
			position: -1,
			wantErr:  true,
		},
		{
			name:     "position too high",
			position: 10,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := config.NewMessages()
			if tt.setup != nil {
				tt.setup(m)
			}

			got, err := m.At(tt.position)
			if (err != nil) != tt.wantErr {
				t.Errorf("At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("At() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMessages_SetAt(t *testing.T) {
	tests := []struct {
		name     string
		position int
		message  string
		want     string
		wantErr  bool
	}{
		{
			name:     "valid simple message",
			position: 0,
			message:  "CQ CQ DE K3GDS",
			want:     "CQ CQ DE K3GDS",
			wantErr:  false,
		},
		{
			name:     "valid message with trim whitespace",
			position: 2,
			message:  " 5NN TU ",
			want:     "5NN TU",
			wantErr:  false,
		},
		{
			name:     "valid message with lowercase",
			position: 3,
			message:  "hello world",
			want:     "HELLO WORLD",
			wantErr:  false,
		},
		{
			name:     "valid message with / and ?",
			position: 5,
			message:  "K3GDS/6?",
			want:     "K3GDS/6?",
			wantErr:  false,
		},
		{
			name:     "invalid message with %",
			position: 0,
			message:  "TEST%",
			wantErr:  true,
		},
		{
			name:     "invalid message with &",
			position: 0,
			message:  "TEST&TEST",
			wantErr:  true,
		},
		{
			name:     "invalid position too low",
			position: -1,
			message:  "TEST",
			wantErr:  true,
		},
		{
			name:     "invalid position too high",
			position: 10,
			message:  "TEST",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := config.NewMessages()
			got, err := m.SetAt(tt.position, tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("SetAt() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMessages_All(t *testing.T) {
	m := config.NewMessages()
	m.SetAt(0, "FIRST")
	m.SetAt(5, "MIDDLE")
	m.SetAt(9, "LAST")

	want := [10]string{}
	want[0] = "FIRST"
	want[5] = "MIDDLE"
	want[9] = "LAST"

	if got := m.All(); got != want {
		t.Errorf("All(), got %q, want %q", got, want)
	}
}

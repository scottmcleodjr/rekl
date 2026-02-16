package config_test

import (
	"testing"

	"github.com/scottmcleodjr/rekl/config"
)

func TestSpeed_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   int
		want    int
		wantErr bool
	}{
		{name: "valid WPM", input: 15, want: 15, wantErr: false},
		{name: "minimum WPM", input: config.MinWPM, want: config.MinWPM, wantErr: false},
		{name: "maximum WPM", input: config.MaxWPM, want: config.MaxWPM, wantErr: false},
		{name: "below minimum", input: 1, want: config.InitialWPM, wantErr: true},
		{name: "above maximum", input: 150, want: config.InitialWPM, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := config.NewSpeed()
			err := s.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got := s.WPM(); got != tt.want {
				t.Errorf("WPM() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestSpeed_Increment(t *testing.T) {
	s := config.NewSpeed()

	// Increment from initial to max
	for {
		startWPM := s.WPM()

		err := s.Increment()
		if err != nil {
			t.Fatalf("at %d WPM, Increment() error = %v, want nil", startWPM, err)
		}

		want := startWPM + 1
		got := s.WPM()
		if got != want {
			t.Errorf("after Increment(), WPM() = %d, want %d", got, want)
		}

		if got == config.MaxWPM {
			break
		}
	}

	// Incrementing past max should error and not change speed
	err := s.Increment()
	if err == nil {
		t.Errorf("at %d WPM, Increment() error = nil, want error", config.MaxWPM)
	}

	if got := s.WPM(); got != config.MaxWPM {
		t.Errorf("after Increment(), WPM() = %d, want %d (unchanged value after error)", got, config.MaxWPM)
	}
}

func TestSpeed_Decrement(t *testing.T) {
	s := config.NewSpeed()

	// Decrement from initial to min
	for {
		startWPM := s.WPM()

		err := s.Decrement()
		if err != nil {
			t.Fatalf("at %d WPM, Decrement() error = %v, want nil", startWPM, err)
		}

		want := startWPM - 1
		got := s.WPM()
		if got != want {
			t.Errorf("after Decrement(), WPM() = %d, want %d", got, want)
		}

		if got == config.MinWPM {
			break
		}
	}

	// Decrementing past min should error and not change speed
	err := s.Decrement()
	if err == nil {
		t.Errorf("at %d WPM, Decrement() error = nil, want error", config.MinWPM)
	}

	if got := s.WPM(); got != config.MinWPM {
		t.Errorf("after Decrement(), WPM() = %d, want %d (unchanged value after error)", got, config.MinWPM)
	}
}

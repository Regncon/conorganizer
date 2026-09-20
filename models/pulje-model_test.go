package models

import "testing"

func TestPuljeLabel_ReturnsHumanReadableNames(t *testing.T) {
	tests := []struct {
		pulje Pulje
		want  string
	}{
		{PuljeFredagKveld, "Fredag kveld"},
		{PuljeLordagMorgen, "Lørdag morgen"},
		{PuljeLordagKveld, "Lørdag kveld"},
		{PuljeSondagMorgen, "Søndag morgen"},
		{Pulje("unknown"), "unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.pulje), func(t *testing.T) {
			if got := tt.pulje.Label(); got != tt.want {
				t.Fatalf("Pulje.Label() = %q, want %q", got, tt.want)
			}
		})
	}
}

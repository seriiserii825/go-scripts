package scripts

import "testing"

func TestRgbToHex(t *testing.T) {
	tests := []struct {
		in     string
		want   string
		wantOk bool
	}{
		{"rgb(255, 87, 51)", "#ff5733", true},
		{" RGB(255,87,51)\n", "#ff5733", true},
		{"255,87,51", "#ff5733", true},
		{"255 87 51", "#ff5733", true},
		{"rgba(255, 87, 51, 0.5)", "#ff573380", true},
		{"rgb(255 87 51 / 50%)", "#ff573380", true},
		{"rgba(255, 87, 51, 1)", "#ff5733", true},
		{"rgb(0, 0, 0)", "#000000", true},
		{"rgb(256, 0, 0)", "", false},
		{"rgb(-1, 0, 0)", "", false},
		{"rgb(255, 87)", "", false},
		{"rgba(255, 87, 51, 2)", "", false},
		{"#ff5733", "", false},
		{"", "", false},
	}
	for _, tt := range tests {
		got, err := rgbToHex(tt.in)
		if (err == nil) != tt.wantOk || got != tt.want {
			t.Errorf("rgbToHex(%q) = %q, %v, want %q", tt.in, got, err, tt.want)
		}
	}
}

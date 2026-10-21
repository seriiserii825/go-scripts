package scripts

import "testing"

func TestParseHex(t *testing.T) {
	tests := []struct {
		in     string
		want   [3]int
		wantOk bool
	}{
		{"#3366cc", [3]int{0x33, 0x66, 0xcc}, true},
		{" #FF5733\n", [3]int{0xff, 0x57, 0x33}, true},
		{"ff5733", [3]int{0xff, 0x57, 0x33}, true},
		{"#fff", [3]int{}, false},
		{"#gg5733", [3]int{}, false},
		{"", [3]int{}, false},
	}
	for _, tt := range tests {
		got, err := parseHex(tt.in)
		if (err == nil) != tt.wantOk || got != tt.want {
			t.Errorf("parseHex(%q) = %v, %v", tt.in, got, err)
		}
	}
}

func TestShade(t *testing.T) {
	tests := []struct {
		in      string
		lighten bool
		percent int
		want    string
	}{
		{"#3366cc", true, 20, "#5c85d6"},
		{"#3366cc", false, 20, "#2952a3"},
		{"#000000", true, 50, "#808080"},
		{"#ffffff", false, 50, "#7f7f7f"},
		{"#123456", true, 100, "#ffffff"},
		{"#123456", false, 100, "#000000"},
	}
	for _, tt := range tests {
		rgb, _ := parseHex(tt.in)
		if got := formatHex(shade(rgb, tt.lighten, tt.percent)); got != tt.want {
			t.Errorf("shade(%s, %v, %d) = %s, want %s", tt.in, tt.lighten, tt.percent, got, tt.want)
		}
	}
}

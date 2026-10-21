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
		{"#fff", [3]int{0xff, 0xff, 0xff}, true},
		{"#36c", [3]int{0x33, 0x66, 0xcc}, true},
		{"#ffff", [3]int{}, false},
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

func TestBuildPalette(t *testing.T) {
	rgb, _ := parseHex("#f0badc")
	steps := buildPalette(rgb)
	if len(steps) != 10 {
		t.Fatalf("len(steps) = %d, want 10", len(steps))
	}

	want := map[int]string{
		10:  "#907084",
		20:  "#a8829a",
		30:  "#c095b0",
		40:  "#d8a7c6",
		50:  "#f0badc",
		60:  "#f2c1e0",
		70:  "#f3c8e3",
		80:  "#f5cfe7",
		90:  "#f6d6ea",
		100: "#f8ddee",
	}
	for _, step := range steps {
		if got := formatHex(step.rgb); got != want[step.index] {
			t.Errorf("index %d = %s, want %s", step.index, got, want[step.index])
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct{ in, want int }{
		{0, 0}, {5, 5}, {-5, 5}, {-40, 40},
	}
	for _, tt := range tests {
		if got := abs(tt.in); got != tt.want {
			t.Errorf("abs(%d) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

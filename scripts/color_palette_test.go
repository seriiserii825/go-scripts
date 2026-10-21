package scripts

import "testing"

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

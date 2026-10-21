package scripts

import "testing"

func TestParsePair(t *testing.T) {
	tests := []struct {
		in     string
		a, b   int
		wantOk bool
	}{
		{"1920,1080", 1920, 1080, true},
		{" 1920 , 1080 ", 1920, 1080, true},
		{"1920", 0, 0, false},
		{"a,1080", 0, 0, false},
		{"0,1080", 0, 0, false},
		{"1920,1080,5", 0, 0, false},
	}
	for _, tt := range tests {
		a, b, err := parsePair(tt.in)
		if (err == nil) != tt.wantOk || a != tt.a || b != tt.b {
			t.Errorf("parsePair(%q) = %d, %d, %v", tt.in, a, b, err)
		}
	}
}

func TestAspectRatio(t *testing.T) {
	tests := map[[2]int]string{
		{1920, 1080}: "1.78/1",
		{1080, 1920}: "1/1.78",
		{800, 400}:   "2/1",
		{500, 500}:   "1/1",
	}
	for in, want := range tests {
		if got := aspectRatio(in[0], in[1]); got != want {
			t.Errorf("aspectRatio(%d, %d) = %q, want %q", in[0], in[1], got, want)
		}
	}
}

package scripts

import "testing"

func TestFontWeight(t *testing.T) {
	tests := map[string]string{
		"Roboto-Regular":     "normal",
		"Roboto-Thin":        "100",
		"Roboto-ExtraLight":  "200",
		"Roboto-Light":       "300",
		"Roboto-Medium":      "500",
		"Roboto-SemiBold":    "600",
		"Roboto-DemiBold":    "600",
		"Roboto-Bold":        "700",
		"Roboto-ExtraBold":   "800",
		"Roboto-Black":       "900",
		"Roboto-HeavyItalic": "900",
		"Roboto-LightItalic": "300",
	}
	for name, want := range tests {
		if got := fontWeight(name); got != want {
			t.Errorf("fontWeight(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestFontFamily(t *testing.T) {
	tests := map[string]string{
		"Roboto-Bold":   "Roboto",
		"OpenSans-Bold": "OpenSans",
		"roboto":        "Roboto",
		"монтсеррат-x":  "Монтсеррат",
	}
	for name, want := range tests {
		if got := fontFamily(name); got != want {
			t.Errorf("fontFamily(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestFontsCSS(t *testing.T) {
	files := []string{"Roboto-Bold.woff", "Roboto-Bold.woff2", "Roboto-Italic.woff2", "fonts.css"}
	want := `@font-face {
  font-family: 'Roboto';
  src: url('assets/fonts/Roboto-Bold.woff2') format('woff2'),
       url('assets/fonts/Roboto-Bold.woff') format('woff');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
`
	if got := fontsCSS(files, "assets/fonts"); got != want {
		t.Errorf("fontsCSS =\n%s\nwant\n%s", got, want)
	}

	onlyWoff2 := fontsCSS([]string{"A-Italic.woff2", "B-Bold.woff2"}, "fonts")
	wantWoff2 := `@font-face {
  font-family: 'A';
  src: url('fonts/A-Italic.woff2') format('woff2');
  font-weight: normal;
  font-style: italic;
  font-display: swap;
}

@font-face {
  font-family: 'B';
  src: url('fonts/B-Bold.woff2') format('woff2');
  font-weight: 700;
  font-style: normal;
  font-display: swap;
}
`
	if onlyWoff2 != wantWoff2 {
		t.Errorf("fontsCSS(only woff2) =\n%s\nwant\n%s", onlyWoff2, wantWoff2)
	}
}

package scripts

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/tools"
)

// RgbToHex берёт rgb-цвет из буфера обмена (или спрашивает), переводит его
// в hex и копирует результат в буфер.
func RgbToHex() error {
	if err := tools.Require(tools.Clipboard(), tools.Notify); err != nil {
		return err
	}
	hex, err := askRgb()
	if err != nil {
		return err
	}
	rgb, _ := parseHex(hex[:7])
	fmt.Println(swatch(rgb), hex)
	return clipboard.Write(hex)
}

// askRgb берёт цвет из буфера обмена, а если там не rgb — спрашивает, пока
// не введут правильно. Возвращает готовый hex.
func askRgb() (string, error) {
	// пустой буфер — тоже ошибка (xclip выходит с 1), просто спрашиваем
	if text, err := clipboard.Read(); err == nil {
		if hex, err := rgbToHex(text); err == nil {
			return hex, nil
		}
	}
	fmt.Println(color.Wrap(color.Yellow, "No rgb color in the clipboard."))
	for {
		s, err := input.Line("Rgb color (rgb(255, 87, 51) or 255,87,51): ")
		if err != nil {
			return "", err
		}
		hex, err := rgbToHex(s)
		if err == nil {
			return hex, nil
		}
		fmt.Println(color.Wrap(color.Red, "Invalid input. "+err.Error()+"."))
	}
}

// rgbFunc снимает обёртку rgb(...)/rgba(...), если она есть.
var rgbFunc = regexp.MustCompile(`(?i)^rgba?\((.*)\)$`)

// rgbToHex разбирает "rgb(255, 87, 51)", "rgba(255, 87, 51, 0.5)",
// "255,87,51" и "255 87 51 / 50%". С альфой меньше 1 возвращает #RRGGBBAA,
// иначе #RRGGBB.
func rgbToHex(s string) (string, error) {
	errBad := errors.New("not an rgb color, use rgb(R, G, B)")
	s = strings.TrimSpace(s)
	if m := rgbFunc.FindStringSubmatch(s); m != nil {
		s = m[1]
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == ',' || r == '/' || r == ' ' || r == '\t'
	})
	if len(parts) != 3 && len(parts) != 4 {
		return "", errBad
	}
	var rgb [3]int
	for i := range rgb {
		n, err := strconv.Atoi(parts[i])
		if err != nil || n < 0 || n > 255 {
			return "", errBad
		}
		rgb[i] = n
	}
	hex := formatHex(rgb)
	if len(parts) == 4 {
		alpha, err := parseAlpha(parts[3])
		if err != nil {
			return "", errBad
		}
		if alpha < 255 {
			hex += fmt.Sprintf("%02x", alpha)
		}
	}
	return hex, nil
}

// parseAlpha переводит альфу "0.5" или "50%" в 0..255.
func parseAlpha(s string) (int, error) {
	percent := strings.HasSuffix(s, "%")
	f, err := strconv.ParseFloat(strings.TrimSuffix(s, "%"), 64)
	if err != nil {
		return 0, err
	}
	if percent {
		f /= 100
	}
	if f < 0 || f > 1 {
		return 0, errors.New("alpha out of range")
	}
	return int(math.Round(f * 255)), nil
}

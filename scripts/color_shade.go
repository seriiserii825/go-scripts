package scripts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/tools"
)

// ColorShade берёт hex-цвет из буфера обмена (или спрашивает), осветляет или затемняет его на
// заданный процент и копирует результат в буфер: #3366cc, l, 20 -> #5c85d6.
func ColorShade() error {
	if err := tools.Require(tools.Clipboard(), tools.Notify); err != nil {
		return err
	}
	rgb, err := askHex()
	if err != nil {
		return err
	}
	fmt.Println("Color:", swatch(rgb), formatHex(rgb))

	lighten, err := askLighten()
	if err != nil {
		return err
	}
	percent, err := askPercent()
	if err != nil {
		return err
	}
	shaded := shade(rgb, lighten, percent)
	fmt.Println("Result:", swatch(shaded), formatHex(shaded))
	return clipboard.Write(formatHex(shaded))
}

// askHex берёт цвет из буфера обмена, а если там не hex — спрашивает, пока
// не введут правильно.
func askHex() ([3]int, error) {
	// пустой буфер — тоже ошибка (xclip выходит с 1), просто спрашиваем
	if text, err := clipboard.Read(); err == nil {
		if rgb, err := parseHex(text); err == nil {
			return rgb, nil
		}
	}
	fmt.Println(color.Wrap(color.Yellow, "No hex color in the clipboard."))
	for {
		s, err := input.Line("Hex color (#RRGGBB): ")
		if err != nil {
			return [3]int{}, err
		}
		rgb, err := parseHex(s)
		if err == nil {
			return rgb, nil
		}
		fmt.Println(color.Wrap(color.Red, "Invalid input. "+err.Error()+"."))
	}
}

// askLighten спрашивает l/d, пока не введут правильно: true — осветлить.
func askLighten() (bool, error) {
	for {
		s, err := input.Line("Lighten or darken? (l/d): ")
		if err != nil {
			return false, err
		}
		switch strings.ToLower(s) {
		case "l", "lighten":
			return true, nil
		case "d", "darken":
			return false, nil
		}
		fmt.Println(color.Wrap(color.Red, "Invalid input. Enter l or d."))
	}
}

func askPercent() (int, error) {
	for {
		n, err := input.PositiveInt("Percent (1-100): ")
		if err != nil {
			return 0, err
		}
		if n <= 100 {
			return n, nil
		}
		fmt.Println(color.Wrap(color.Red, "Invalid input. Max is 100."))
	}
}

// parseHex разбирает "#RRGGBB" (решётка необязательна, регистр любой).
func parseHex(s string) ([3]int, error) {
	hex := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(hex) != 6 {
		return [3]int{}, errors.New("not a hex color, use #RRGGBB")
	}
	var rgb [3]int
	for i := range rgb {
		n, err := strconv.ParseUint(hex[i*2:i*2+2], 16, 8)
		if err != nil {
			return [3]int{}, errors.New("not a hex color, use #RRGGBB")
		}
		rgb[i] = int(n)
	}
	return rgb, nil
}

func formatHex(rgb [3]int) string {
	return fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])
}

// shade смешивает цвет с белым (lighten) или чёрным (darken) на percent%:
// 100% lighten — всегда #ffffff, 100% darken — #000000.
func shade(rgb [3]int, lighten bool, percent int) [3]int {
	var out [3]int
	for i, c := range rgb {
		target := 0
		if lighten {
			target = 255
		}
		// +50 — округление до ближайшего при делении на 100
		out[i] = c + ((target-c)*percent+sign(target-c)*50)/100
	}
	return out
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}

// swatch — квадратик этого цвета в терминале (truecolor).
func swatch(rgb [3]int) string {
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm    %s", rgb[0], rgb[1], rgb[2], color.Reset)
}

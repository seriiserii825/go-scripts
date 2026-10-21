package scripts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
	"github.com/seriiserii825/go-libs/tools"
)

// ColorShade берёт hex-цвет из буфера обмена (или спрашивает), строит на его
// основе палитру из 10 оттенков с шагом 10% (сам цвет — в середине, на
// 50%), даёт выбрать один через fzf и копирует выбранный цвет в буфер.
func ColorShade() error {
	if err := tools.Require(tools.Fzf, tools.FzfFilter, tools.Clipboard(), tools.Notify); err != nil {
		return err
	}
	rgb, err := askHex()
	if err != nil {
		return err
	}
	fmt.Println("Base color:", swatch(rgb), formatHex(rgb))

	steps := buildPalette(rgb)
	lines := make([]string, len(steps))
	for i, step := range steps {
		lines[i] = fmt.Sprintf("%s %3d%%  %s", swatch(step.rgb), step.index, formatHex(step.rgb))
	}

	out, code, err := fzf.Run(lines, []string{"--ansi", "--height", "50%", "--reverse", "--no-mouse"})
	if err != nil {
		return err
	}
	if code != 0 {
		msg.Warning("No color selected.")
		return nil
	}
	selected := strings.TrimRight(out, "\n")
	hex := selected[strings.LastIndex(selected, "#"):]
	selectedRgb, err := parseHex(hex)
	if err != nil {
		return err
	}
	fmt.Println("Selected:", swatch(selectedRgb), hex)
	return clipboard.Write(hex)
}

type paletteStep struct {
	index int
	rgb   [3]int
}

// buildPalette строит 10 оттенков (10%..100%) от rgb: 50% — сам цвет,
// каждые 10 пунктов разницы с 50 — это shade на такой же процент (40% и
// 60% — ±10%, 10% и 100% — ±40%/50%).
func buildPalette(rgb [3]int) []paletteStep {
	steps := make([]paletteStep, 0, 10)
	for index := 10; index <= 100; index += 10 {
		diff := index - 50
		if diff == 0 {
			steps = append(steps, paletteStep{index, rgb})
			continue
		}
		steps = append(steps, paletteStep{index, shade(rgb, diff > 0, abs(diff))})
	}
	return steps
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
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
		s, err := input.Line("Hex color (#RRGGBB or #RGB): ")
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

// parseHex разбирает "#RRGGBB" и короткий "#RGB" (#fff -> #ffffff);
// решётка необязательна, регистр любой.
func parseHex(s string) ([3]int, error) {
	hex := strings.TrimPrefix(strings.TrimSpace(s), "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
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

package scripts

import (
	"fmt"
	"strings"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/tools"
)

// ColorPalette берёт hex-цвет из буфера обмена (или спрашивает) и строит
// палитру из 10 оттенков с шагом 10%: сам цвет стоит на 50%, индексы
// ниже — темнее, выше — светлее (40% -> затемнение на 10%, 60% ->
// осветление на 10%, 10% -> затемнение на 40%, 100% -> осветление на
// 50%, и т.д.). Выводит все оттенки и копирует CSS-переменные в буфер.
func ColorPalette() error {
	if err := tools.Require(tools.Clipboard(), tools.Notify); err != nil {
		return err
	}
	rgb, err := askHex()
	if err != nil {
		return err
	}
	fmt.Println("Base color:", swatch(rgb), formatHex(rgb))

	steps := buildPalette(rgb)
	for _, step := range steps {
		fmt.Printf("%3d%%  %s %s\n", step.index, swatch(step.rgb), formatHex(step.rgb))
	}
	css := formatPaletteCSS(steps)
	fmt.Println(css)
	return clipboard.Write(css)
}

type paletteStep struct {
	index int
	rgb   [3]int
}

// buildPalette строит 10 оттенков (10%..100%) от rgb: 50% — сам цвет,
// каждые 10 пунктов разницы с 50 — это shade на такой же процент.
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

// formatPaletteCSS формирует CSS custom properties для копирования в буфер.
func formatPaletteCSS(steps []paletteStep) string {
	var b strings.Builder
	for _, step := range steps {
		fmt.Fprintf(&b, "--color-%d: %s;\n", step.index, formatHex(step.rgb))
	}
	return strings.TrimRight(b.String(), "\n")
}

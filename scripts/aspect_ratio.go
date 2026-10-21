package scripts

import (
	"fmt"
	"math"
	"strconv"

	"github.com/seriiserii825/go-libs/clipboard"
)

// AspectRatio спрашивает ширину и высоту и копирует в буфер обмена
// CSS-свойство: "aspect-ratio: 1.78/1;".
func AspectRatio() error {
	if err := require(toolClipboard(), toolNotify); err != nil {
		return err
	}
	width, height, err := askPair("Enter width and height w,h: ")
	if err != nil {
		return err
	}
	css := fmt.Sprintf("aspect-ratio: %s;", aspectRatio(width, height))
	fmt.Println(css)
	return clipboard.Write(css)
}

// aspectRatio — большая сторона к меньшей, округлено до 2 знаков:
// 1920x1080 -> "1.78/1", 1080x1920 -> "1/1.78".
func aspectRatio(width, height int) string {
	switch {
	case width > height:
		return round2(float64(width)/float64(height)) + "/1"
	case width < height:
		return "1/" + round2(float64(height)/float64(width))
	default:
		return "1/1"
	}
}

// round2 округляет до 2 знаков без лишних нулей: 2 -> "2", 1.777 -> "1.78".
func round2(f float64) string {
	return strconv.FormatFloat(math.Round(f*100)/100, 'f', -1, 64)
}

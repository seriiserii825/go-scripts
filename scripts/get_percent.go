package scripts

import (
	"fmt"

	"github.com/seriiserii825/go-libs/clipboard"
)

// GetPercent считает, сколько процентов от полной ширины составляет ширина,
// и копирует результат ("30.00%") в буфер обмена.
func GetPercent() error {
	if err := require(toolClipboard(), toolNotify); err != nil {
		return err
	}
	fullWidth, width, err := askPair("Enter full width and width f,w: ")
	if err != nil {
		return err
	}
	percent := fmt.Sprintf("%.2f%%", float64(width)/float64(fullWidth)*100)
	fmt.Println(percent)
	return clipboard.Write(percent)
}

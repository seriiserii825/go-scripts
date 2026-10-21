package scripts

import (
	"fmt"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/input"
)

// CalcBlocks считает, сколько процентов ширины занимает блок, и копирует
// результат ("30.00%") в буфер обмена с уведомлением.
func CalcBlocks() error {
	fullWidth, err := input.PositiveInt("Enter full width: ")
	if err != nil {
		return err
	}
	blockWidth, err := input.PositiveInt("Enter block width: ")
	if err != nil {
		return err
	}
	percent := float64(blockWidth) / float64(fullWidth) * 100
	fmt.Printf("Block width is %.2f%% of %d width.\n", percent, fullWidth)
	return clipboard.Write(fmt.Sprintf("%.2f%%", percent))
}

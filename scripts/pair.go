package scripts

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/input"
)

// askPair спрашивает два положительных числа через запятую ("1920,1080"),
// пока не введут правильно.
func askPair(prompt string) (int, int, error) {
	for {
		s, err := input.Line(prompt)
		if err != nil {
			return 0, 0, err
		}
		a, b, err := parsePair(s)
		if err == nil {
			return a, b, nil
		}
		fmt.Println(color.Wrap(color.Red, "Invalid input. "+err.Error()))
	}
}

func parsePair(s string) (int, int, error) {
	first, second, found := strings.Cut(s, ",")
	if !found {
		return 0, 0, errors.New("Enter two numbers separated by comma.")
	}
	a, errA := strconv.Atoi(strings.TrimSpace(first))
	b, errB := strconv.Atoi(strings.TrimSpace(second))
	if errA != nil || errB != nil {
		return 0, 0, errors.New("Both values must be numbers.")
	}
	if a <= 0 || b <= 0 {
		return 0, 0, errors.New("Both values must be greater than 0.")
	}
	return a, b, nil
}

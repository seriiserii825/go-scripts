package main

import (
	"fmt"

	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/msg"
)

func main() {
	errorMessage := "An error occurred"
	msg.Success("Hello, World!")
	msg.Error(errorMessage)

	s, ok, err := fzf.One([]string{"Restore from Downloads", "Backup"})
	fmt.Println(s, ok, err)
}

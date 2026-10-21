package main

import (
	"fmt"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/menu"
	"github.com/seriiserii825/go-libs/msg"
	"github.com/seriiserii825/go-libs/table"
)

func main() {
	errorMessage := "An error occurred"
	msg.Success("Hello, World!")
	msg.Error(errorMessage)

	s, ok, err := fzf.One([]string{"Restore from Downloads", "Backup"})
	fmt.Println(s, ok, err)

	table.Show("Сайты", []string{"#", "Name", "Status"}, [][]string{
		{"0", "Restore from Downloads", "ok"},
		{"1", "Бэкап", "fail"},
		{"2", "Deploy", "ok"},
	}, map[int]string{1: color.Red})

	var items []string
	for i := range 30 {
		items = append(items, fmt.Sprintf("пункт-%02d", i))
	}
	table.Grid(items, color.Blue)

	n, ok, err := menu.SelectFzfMenu(
		[]string{color.Wrap(color.Green, "Deploy"), "Backup"},
		"Что делаем?",
		color.Wrap(color.Red, "Exit"),
	)
	fmt.Println(n, ok, err)
}

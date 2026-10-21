package main

import (
	"errors"
	"io"
	"os"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/menu"
	"github.com/seriiserii825/go-libs/msg"

	"github.com/seriiserii825/go-scripts/scripts"
)

type item struct {
	label string
	run   func() error
}

func notPorted() error {
	msg.Warning("Not ported to Go yet")
	return nil
}

var items = []item{
	{color.Wrap(color.Green, "Create project"), notPorted},
	{color.Wrap(color.Blue, "Convert fonts"), notPorted},
	{color.Wrap(color.Green, "Cyr to lat"), notPorted},
	{color.Wrap(color.Blue, "Bitbucket change url"), notPorted},
	{color.Wrap(color.Green, "Open github or bitbucket"), notPorted},
	{color.Wrap(color.Yellow, "Generate password"), notPorted},
	{color.Wrap(color.Green, "Aspect Ratio"), scripts.AspectRatio},
	{color.Wrap(color.Blue, "Sort files"), notPorted},
	{color.Wrap(color.Green, "Git mirror"), notPorted},
	{color.Wrap(color.Green, "Get percent from width"), scripts.GetPercent},
}

func main() {
	labels := make([]string, len(items))
	for i, it := range items {
		labels[i] = it.label
	}

	n, ok, err := menu.SelectFzfMenu(labels, "Scripts", color.Wrap(color.Red, "Exit"))
	if err != nil {
		msg.Error(err.Error())
		os.Exit(1)
	}
	if !ok || n == 0 {
		return
	}

	if err := items[n-1].run(); err != nil && !errors.Is(err, io.EOF) {
		msg.Error(err.Error())
		os.Exit(1)
	}
}

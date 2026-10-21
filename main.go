package main

import (
	"errors"
	"io"
	"os"

	"github.com/seriiserii825/go-libs/cli"
	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/msg"

	"github.com/seriiserii825/go-scripts/scripts"
)

func notPorted() error {
	msg.Warning("Not ported to Go yet")
	return nil
}

var app = cli.App{
	Name: "go-scripts",
	Commands: []cli.Command{
		{
			Name:  "create-project",
			Short: "Create project",
			Long:  "Not ported to Go yet.",
			Color: color.Green,
			Run:   notPorted,
		},
		{
			Name:  "convert-fonts",
			Short: "Convert fonts",
			Long:  "Not ported to Go yet.",
			Color: color.Blue,
			Run:   notPorted,
		},
		{
			Name:  "cyr-to-lat",
			Short: "Cyr to lat",
			Long: `Pick a file in the current folder with fzf, strip its ID3 tags
(id3convert -s, skipped if not installed) and rename it from Cyrillic
to Latin:

  "Песня Щука.mp3" -> "Pesnya SCHuka.mp3"

Fails if a file with the new name already exists.`,
			Color: color.Green,
			Run:   scripts.CyrToLat,
		},
		{
			Name:  "generate-password",
			Short: "Generate password",
			Long: `Ask the password length and whether to include symbols (y/n),
generate 10 random passwords (crypto/rand), pick one with fzf and copy
it to the clipboard. The password is not shown in the notification.`,
			Color: color.Yellow,
			Run:   scripts.GeneratePassword,
		},
		{
			Name:  "aspect-ratio",
			Short: "Aspect Ratio",
			Long: `Ask width and height separated by comma and copy the CSS property
to the clipboard:

  1920,1080 -> aspect-ratio: 1.78/1;
  1080,1920 -> aspect-ratio: 1/1.78;`,
			Color: color.Green,
			Run:   scripts.AspectRatio,
		},
		{
			Name:  "sort-files",
			Short: "Sort files",
			Long: `Copy every file from the current folder and its subfolders into
folders by extension. Originals stay in place:

  sub/photo.jpg -> jpg/photo.jpg
  doc.pdf       -> pdf/doc.pdf

Skipped: hidden files and folders (.git, .env), files without
extension, files already in their extension folder. Files with the
same name overwrite each other. JPG and jpg go to different folders.

After copying, empty subfolders are removed. Asks for confirmation
(y/n) first.`,
			Color: color.Blue,
			Run:   scripts.SortFiles,
		},
		{
			Name:  "get-percent",
			Short: "Get percent from width",
			Long: `Ask full width and width separated by comma and copy the percent
to the clipboard:

  1000,300 -> 30.00%`,
			Color: color.Green,
			Run:   scripts.GetPercent,
		},
	},
}

func main() {
	if err := app.Run(os.Args[1:]); err != nil && !errors.Is(err, io.EOF) {
		msg.Error(err.Error())
		os.Exit(1)
	}
}

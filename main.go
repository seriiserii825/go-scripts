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

var app = cli.App{
	Name: "go-scripts",
	Commands: []cli.Command{
		{
			Name:  "create-project",
			Short: "Create project",
			Long: `Create a new WordPress theme in the current folder from the starter
~/Documents/Sites/wp-projects/bs-vite:

  - copies the starter into ./<project-name> (without .git)
  - replaces bs-vite with <project-name> in all files
  - replaces bs_vite with <project_name> (dashes -> underscores)

  my-site -> ./my-site, "bs-vite" -> "my-site", "bs_vite" -> "my_site"

Fails if ./<project-name> already exists or the starter is missing.`,
			Color: color.Green,
			Run:   scripts.CreateProject,
		},
		{
			Name:  "convert-fonts",
			Short: "Convert fonts",
			Long: `Run from the folder with font files (FontName-Weight.ttf):

  - converts every .ttf to .woff (ttf2woff) and .woff2 (woff2_compress),
    the .ttf is removed only if both conversions succeed
  - asks the relative path to fonts (default: assets/fonts)
  - writes fonts.css with @font-face for every font and copies it
    to the clipboard

Weight and style come from the file name: Thin 100, ExtraLight 200,
Light 300, Regular normal, Medium 500, SemiBold/DemiBold 600, Bold 700,
ExtraBold 800, Black/Heavy 900, Italic -> font-style: italic.
font-family is the part before the first "-": OpenSans-Bold -> OpenSans.

Needs: npm install -g ttf2woff, sudo pacman -S woff2.`,
			Color: color.Blue,
			Run:   scripts.ConvertFonts,
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
	if len(os.Args) == 1 { // без аргументов — fzf-меню
		if err := scripts.CheckMenuTools(); err != nil {
			msg.Error(err.Error())
			os.Exit(1)
		}
	}
	if err := app.Run(os.Args[1:]); err != nil && !errors.Is(err, io.EOF) {
		msg.Error(err.Error())
		os.Exit(1)
	}
}

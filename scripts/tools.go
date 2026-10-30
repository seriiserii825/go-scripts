package scripts

import "github.com/seriiserii825/go-libs/tools"

// Программы, которые нужны только отдельным скриптам; общие (fzf,
// fzf-word-filter, notify-send, буфер обмена) — в go-libs/tools.
var (
	toolWoff2      = tools.Tool{Bin: "woff2_compress", Install: "sudo pacman -S woff2"}
	toolTTF2Woff   = tools.Tool{Bin: "ttf2woff", Install: "npm install -g ttf2woff"}
	toolPyftsubset = tools.Tool{
		Bin:     "pyftsubset",
		Install: "sudo pacman -S python-fonttools python-brotli",
	}
	toolID3 = tools.Tool{
		Bin:      "id3convert",
		Install:  "sudo pacman -S id3lib",
		Optional: true,
		Without:  "ID3 tags are not stripped",
	}
)

// CheckMenuTools — проверка для fzf-меню (go-scripts без аргументов).
func CheckMenuTools() error {
	return tools.Require(tools.Fzf, tools.FzfFilter)
}

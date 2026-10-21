package scripts

import (
	"fmt"
	"os"
	"strings"

	"github.com/seriiserii825/go-libs/command"
	"github.com/seriiserii825/go-libs/msg"
)

// tool — внешняя программа, которую вызывает скрипт.
type tool struct {
	bin      string
	install  string
	optional bool   // без неё скрипт работает, только хуже
	without  string // что будет без optional-программы
}

var (
	toolFzf       = tool{bin: "fzf", install: "sudo pacman -S fzf"}
	toolWoff2     = tool{bin: "woff2_compress", install: "sudo pacman -S woff2"}
	toolTTF2Woff  = tool{bin: "ttf2woff", install: "npm install -g ttf2woff"}
	toolFzfFilter = tool{
		bin:      "fzf-word-filter",
		install:  "go install github.com/seriiserii825/go-libs/cmd/fzf-word-filter@latest (or copy the binary to ~/.local/bin)",
		optional: true,
		without:  "fzf uses its own search",
	}
	toolNotify = tool{
		bin:      "notify-send",
		install:  "sudo pacman -S libnotify",
		optional: true,
		without:  "no desktop notifications",
	}
	toolID3 = tool{
		bin:      "id3convert",
		install:  "sudo pacman -S id3lib",
		optional: true,
		without:  "ID3 tags are not stripped",
	}
)

// toolClipboard — wl-copy в Wayland, xclip в X11 (как в пакете clipboard).
func toolClipboard() tool {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return tool{bin: "wl-copy", install: "sudo pacman -S wl-clipboard"}
	}
	return tool{bin: "xclip", install: "sudo pacman -S xclip"}
}

// require проверяет, что программы установлены. Нет обязательных — ошибка
// со списком команд установки; нет необязательных — только предупреждение.
func require(tools ...tool) error {
	var missing []string
	for _, t := range tools {
		if command.Exists(t.bin) {
			continue
		}
		if t.optional {
			msg.Warning(fmt.Sprintf("%s is not installed (%s). Install: %s", t.bin, t.without, t.install))
			continue
		}
		missing = append(missing, fmt.Sprintf("  %s: %s", t.bin, t.install))
	}
	if len(missing) > 0 {
		return fmt.Errorf("required programs are not installed:\n%s", strings.Join(missing, "\n"))
	}
	return nil
}

// CheckMenuTools — проверка для fzf-меню (go-scripts без аргументов).
func CheckMenuTools() error {
	return require(toolFzf, toolFzfFilter)
}

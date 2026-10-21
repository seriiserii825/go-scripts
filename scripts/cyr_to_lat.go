package scripts

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/command"
	"github.com/seriiserii825/go-libs/fzf"
	"github.com/seriiserii825/go-libs/msg"
)

var cyrToLatReplacer = strings.NewReplacer(
	"Ё", "YO", "Й", "I", "Ц", "TS", "У", "U", "К", "K", "Е", "E", "Н", "N",
	"Г", "G", "Ш", "SH", "Щ", "SCH", "З", "Z", "Х", "H", "Ъ", "I", "Ф", "F",
	"Ы", "I", "В", "V", "А", "A", "П", "P", "Р", "R", "О", "O", "Л", "L",
	"Д", "D", "Ж", "ZH", "Э", "E", "Я", "Ya", "Ч", "CH", "С", "S", "М", "M",
	"И", "I", "Т", "T", "Ь", "I", "Б", "B", "Ю", "YU",
	"ё", "yo", "й", "i", "ц", "ts", "у", "u", "к", "k", "е", "e", "н", "n",
	"г", "g", "ш", "sh", "щ", "sch", "з", "z", "х", "h", "ъ", "i", "ф", "f",
	"ы", "i", "в", "v", "а", "a", "п", "p", "р", "r", "о", "o", "л", "l",
	"д", "d", "ж", "zh", "э", "e", "я", "ya", "ч", "ch", "с", "s", "м", "m",
	"и", "i", "т", "t", "ь", "i", "б", "b", "ю", "yu",
)

// CyrToLat даёт выбрать файл в текущей папке, убирает из него ID3-теги
// (id3convert -s) и переименовывает кириллицу в латиницу.
func CyrToLat() error {
	if err := require(toolFzf, toolFzfFilter, toolID3); err != nil {
		return err
	}
	files, err := visibleFiles()
	if err != nil {
		return err
	}

	filename, ok, err := fzf.One(files)
	if err != nil {
		return err
	}
	if !ok {
		msg.Warning("No file selected")
		return nil
	}
	fmt.Println(color.Wrap(color.Green, filename))

	if command.Exists("id3convert") {
		if err := command.Run("id3convert", "-s", filename); err != nil {
			msg.Warning(err.Error())
		}
	}

	newName := cyrToLatReplacer.Replace(filename)
	if newName != filename {
		if _, err := os.Stat(newName); err == nil {
			return fmt.Errorf("%s already exists", newName)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		if err := os.Rename(filename, newName); err != nil {
			return err
		}
	}

	files, err = visibleFiles()
	if err != nil {
		return err
	}
	fmt.Println(color.Wrap(color.Yellow, "Changed files===================="))
	for _, f := range files {
		fmt.Println(color.Wrap(color.Blue, f))
	}
	return nil
}

// visibleFiles — имена в текущей папке без скрытых (с точки).
func visibleFiles() ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

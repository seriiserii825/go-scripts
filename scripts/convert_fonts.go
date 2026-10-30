package scripts

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/seriiserii825/go-libs/clipboard"
	"github.com/seriiserii825/go-libs/color"
	"github.com/seriiserii825/go-libs/command"
	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
	"github.com/seriiserii825/go-libs/tools"
)

// subsetUnicodes — что остаётся в шрифте после обрезки: латиница (с
// расширенной A/B), комбинируемые диакритики, кириллица, типографская
// пунктуация, валюты, №, ™.
const subsetUnicodes = "U+0000-024F,U+0300-036F,U+0400-04FF,U+2000-206F,U+20A0-20CF,U+2116,U+2122"

// ConvertFonts по желанию обрезает шрифты из текущей папки до латиницы и
// кириллицы, конвертирует .ttf/.otf в .woff и .woff2 (исходник удаляется
// только если обе конвертации прошли), создаёт fonts.css с @font-face для
// каждого шрифта и копирует его в буфер обмена.
func ConvertFonts() error {
	start := time.Now()

	files, err := visibleFiles()
	if err != nil {
		return err
	}
	src := append(withExt(files, ".ttf"), withExt(files, ".otf")...)
	if len(src) == 0 && len(withExt(files, ".woff")) == 0 && len(withExt(files, ".woff2")) == 0 {
		printFontsHelp()
		return nil
	}

	subset, err := input.Confirm("Subset fonts to Latin + Cyrillic (drops icons and other scripts)?")
	if err != nil {
		return err
	}

	need := []tools.Tool{tools.Clipboard(), tools.Notify}
	if len(src) > 0 {
		need = append(need, toolTTF2Woff, toolWoff2)
	}
	if subset {
		need = append(need, toolPyftsubset)
	}
	if err := tools.Require(need...); err != nil {
		return err
	}

	if subset {
		// Обрезаем до конвертации: тогда .woff/.woff2 получаются уже из
		// обрезанного исходника, а готовые .woff/.woff2 обрезаются на месте.
		fonts := slices.Concat(src, withExt(files, ".woff"), withExt(files, ".woff2"))
		for _, f := range fonts {
			if err := subsetFont(f); err != nil {
				msg.Warning(fmt.Sprintf("%s not subset: %v", f, err))
			}
		}
	}

	if len(src) > 0 {
		for _, f := range src {
			// Битый шрифт не должен останавливать остальные: пропускаем его
			// и оставляем исходник на месте.
			if err := convertFont(f); err != nil {
				msg.Warning(fmt.Sprintf("%s skipped: %v", f, err))
				continue
			}
			if err := os.Remove(f); err != nil {
				return err
			}
		}
		if files, err = visibleFiles(); err != nil {
			return err
		}
	}

	relPath, err := input.String("Enter relative path to fonts folder (default: assets/fonts): ", true)
	if err != nil {
		return err
	}
	if relPath == "" {
		relPath = "assets/fonts"
	}

	css := fontsCSS(files, strings.TrimSuffix(relPath, "/"))
	if err := os.WriteFile("fonts.css", []byte(css), 0o644); err != nil {
		return err
	}
	if err := clipboard.Write(css); err != nil {
		return err
	}
	msg.Success("fonts.css created and copied to clipboard")
	fmt.Printf("--- %.2f seconds ---\n", time.Since(start).Seconds())
	return nil
}

// subsetFont обрезает шрифт на месте до subsetUnicodes, сохраняя формат
// (sfnt, woff или woff2). Пишет во временный файл, чтобы при ошибке
// pyftsubset исходник остался целым.
func subsetFont(f string) error {
	before, err := os.Stat(f)
	if err != nil {
		return err
	}
	tmp := f + ".subset"
	// MERG, meta, trak — таблицы Apple, которые браузеры не читают, а
	// pyftsubset не умеет обрезать и выкидывает с WARNING; выкидываем явно.
	args := []string{f, "--unicodes=" + subsetUnicodes, "--layout-features=*",
		"--drop-tables+=MERG,meta,trak", "--output-file=" + tmp}
	if flavor := strings.TrimPrefix(filepath.Ext(f), "."); flavor == "woff" || flavor == "woff2" {
		args = append(args, "--flavor="+flavor)
	}
	if err := command.Run("pyftsubset", args...); err != nil {
		os.Remove(tmp)
		return err
	}
	after, err := os.Stat(tmp)
	if err != nil {
		return err
	}
	if err := os.Rename(tmp, f); err != nil {
		return err
	}
	msg.Info(fmt.Sprintf("%s: %d KB -> %d KB", f, before.Size()/1024, after.Size()/1024))
	return nil
}

// convertFont создаёт рядом с X.ttf/X.otf файлы X.woff и X.woff2; если одна
// из конвертаций не удалась, уже созданные файлы удаляются. ttf2woff и
// woff2_compress работают с любым sfnt, в том числе с CFF-шрифтами .otf.
func convertFont(f string) error {
	base := strings.TrimSuffix(f, filepath.Ext(f))
	if err := command.Run("ttf2woff", f, base+".woff"); err != nil {
		os.Remove(base + ".woff")
		return err
	}
	if err := command.Run("woff2_compress", f); err != nil {
		os.Remove(base + ".woff")
		os.Remove(base + ".woff2")
		return err
	}
	return nil
}

// fontsCSS собирает @font-face для каждого шрифта, у которого есть .woff
// или .woff2; в src идут те форматы, что есть в папке.
func fontsCSS(files []string, relPath string) string {
	var names []string
	for _, f := range files {
		ext := filepath.Ext(f)
		if ext != ".woff" && ext != ".woff2" {
			continue
		}
		if name := strings.TrimSuffix(f, ext); !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	slices.Sort(names)

	var css strings.Builder
	for _, name := range names {
		var src []string
		if slices.Contains(files, name+".woff2") {
			src = append(src, fmt.Sprintf("url('%s/%s.woff2') format('woff2')", relPath, name))
		}
		if slices.Contains(files, name+".woff") {
			src = append(src, fmt.Sprintf("url('%s/%s.woff') format('woff')", relPath, name))
		}
		if css.Len() > 0 {
			css.WriteString("\n")
		}
		fmt.Fprintf(&css, `@font-face {
  font-family: '%s';
  src: %s;
  font-weight: %s;
  font-style: %s;
  font-display: swap;
}
`, fontFamily(name), strings.Join(src, ",\n       "), fontWeight(name), fontStyle(name))
	}
	return css.String()
}

// fontFamily — часть имени до первого "-" с заглавной первой буквой:
// "OpenSans-Bold" -> "OpenSans".
func fontFamily(name string) string {
	family, _, _ := strings.Cut(name, "-")
	r, size := utf8.DecodeRuneInString(family)
	return string(unicode.ToUpper(r)) + family[size:]
}

// fontWeight определяет вес по имени файла; порядок проверок как в Python:
// более поздняя перекрывает раннюю ("SemiBold" сначала 700, потом 600).
func fontWeight(name string) string {
	n := strings.ToLower(name)
	weight := "normal"
	if strings.Contains(n, "extralight") {
		weight = "200"
	} else if strings.Contains(n, "light") {
		weight = "300"
	}
	if strings.Contains(n, "extrabold") {
		weight = "800"
	} else if strings.Contains(n, "bold") {
		weight = "700"
	}
	if strings.Contains(n, "thin") {
		weight = "100"
	}
	if strings.Contains(n, "medium") {
		weight = "500"
	}
	if strings.Contains(n, "semibold") || strings.Contains(n, "demibold") {
		weight = "600"
	}
	if strings.Contains(n, "black") || strings.Contains(n, "heavy") {
		weight = "900"
	}
	return weight
}

func fontStyle(name string) string {
	if strings.Contains(strings.ToLower(name), "italic") {
		return "italic"
	}
	return "normal"
}

func withExt(files []string, ext string) []string {
	var out []string
	for _, f := range files {
		if strings.HasSuffix(f, ext) {
			out = append(out, f)
		}
	}
	return out
}

func printFontsHelp() {
	fmt.Println(color.Wrap(color.Red, "[ERROR]") + " No .ttf, .otf, .woff or .woff2 files found in the current directory.")
	fmt.Println("Expected format: FontName-Weight.ttf, FontName-Weight.otf or FontName-Weight.woff")
	fmt.Println("\nExamples:")
	for _, f := range []string{"Roboto-Regular.ttf", "Roboto-Bold.ttf", "Roboto-Italic.ttf", "Roboto-Light.ttf", "Roboto-ExtraBold.ttf", "Roboto-SemiBold.ttf"} {
		fmt.Println("  " + color.Wrap(color.Green, f))
	}
	fmt.Println("\nRun the script from the folder containing your font files.")
}

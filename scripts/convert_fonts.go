package scripts

import (
	"fmt"
	"os"
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

// ConvertFonts конвертирует .ttf из текущей папки в .woff и .woff2 (ttf
// удаляется только если обе конвертации прошли), создаёт fonts.css с
// @font-face для каждого шрифта и копирует его в буфер обмена.
func ConvertFonts() error {
	start := time.Now()

	files, err := visibleFiles()
	if err != nil {
		return err
	}
	ttf := withExt(files, ".ttf")
	if len(ttf) == 0 && len(withExt(files, ".woff")) == 0 && len(withExt(files, ".woff2")) == 0 {
		printFontsHelp()
		return nil
	}

	need := []tools.Tool{tools.Clipboard(), tools.Notify}
	if len(ttf) > 0 {
		need = append(need, toolTTF2Woff, toolWoff2)
	}
	if err := tools.Require(need...); err != nil {
		return err
	}

	if len(ttf) > 0 {
		for _, f := range ttf {
			// Битый шрифт не должен останавливать остальные: пропускаем его
			// и оставляем .ttf на месте.
			if err := convertTTF(f); err != nil {
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

// convertTTF создаёт рядом с X.ttf файлы X.woff и X.woff2; если одна из
// конвертаций не удалась, уже созданные файлы удаляются.
func convertTTF(f string) error {
	base := strings.TrimSuffix(f, ".ttf")
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

// fontsCSS собирает @font-face для каждого шрифта. Основа — .woff, если их
// нет — .woff2; в src идут те форматы, что есть в папке.
func fontsCSS(files []string, relPath string) string {
	ext := ".woff"
	base := withExt(files, ext)
	if len(base) == 0 {
		ext = ".woff2"
		base = withExt(files, ext)
	}

	var css strings.Builder
	for _, f := range base {
		name := strings.TrimSuffix(f, ext)
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
	fmt.Println(color.Wrap(color.Red, "[ERROR]") + " No .ttf, .woff or .woff2 files found in the current directory.")
	fmt.Println("Expected format: FontName-Weight.ttf or FontName-Weight.woff")
	fmt.Println("\nExamples:")
	for _, f := range []string{"Roboto-Regular.ttf", "Roboto-Bold.ttf", "Roboto-Italic.ttf", "Roboto-Light.ttf", "Roboto-ExtraBold.ttf", "Roboto-SemiBold.ttf"} {
		fmt.Println("  " + color.Wrap(color.Green, f))
	}
	fmt.Println("\nRun the script from the folder containing your font files.")
}

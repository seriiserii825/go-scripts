package scripts

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
)

// SortFiles копирует все файлы текущей папки и подпапок в папки по
// расширению (photo.jpg -> jpg/photo.jpg), затем удаляет пустые подпапки.
// Скрытые файлы и папки (.git и т.п.) не трогаются.
func SortFiles() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}
	ok, err := input.Bool(fmt.Sprintf("current folder is: '%s', do you want to continue (y/n)? ", root))
	if err != nil || !ok {
		return err
	}

	copied, err := sortFiles(root)
	if err != nil {
		return err
	}
	msg.Success(fmt.Sprintf("Copied %d files", copied))
	return removeEmptyFolders(root)
}

func sortFiles(root string) (int, error) {
	// Сначала собираем список, потом копируем: иначе обход зашёл бы в только
	// что созданные папки.
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path != root && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() && filepath.Ext(d.Name()) != "" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	copied := 0
	for _, src := range files {
		folder := filepath.Join(root, filepath.Ext(src)[1:])
		dst := filepath.Join(folder, filepath.Base(src))
		if dst == src {
			continue // уже лежит в своей папке
		}
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			if err := os.Mkdir(folder, 0o755); err != nil {
				return copied, err
			}
			fmt.Printf("Folder %s created.\n", filepath.Base(folder))
		}
		if err := copyFile(src, dst); err != nil {
			return copied, err
		}
		copied++
	}
	return copied, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// removeEmptyFolders удаляет пустые подпапки, начиная с самых глубоких:
// так удаляется и папка, в которой были только пустые папки.
func removeEmptyFolders(root string) error {
	var dirs []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() || path == root {
			return nil
		}
		if strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		dirs = append(dirs, path)
		return nil
	})
	if err != nil {
		return err
	}

	for _, dir := range slices.Backward(dirs) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			if err := os.Remove(dir); err != nil {
				return err
			}
			fmt.Println("Removed empty folder:", dir)
		}
	}
	return nil
}

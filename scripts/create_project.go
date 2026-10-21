package scripts

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/seriiserii825/go-libs/input"
	"github.com/seriiserii825/go-libs/msg"
)

const (
	starterName           = "bs-vite"
	starterNameUnderscore = "bs_vite"
)

// CreateProject создаёт в текущей папке новый проект из шаблона
// ~/Documents/Sites/wp-projects/bs-vite: копирует его без .git и заменяет
// в файлах bs-vite на имя проекта, bs_vite — на имя с подчёркиваниями.
func CreateProject() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	starter := filepath.Join(home, "Documents/Sites/wp-projects", starterName)
	fmt.Println("PROJECT_STARTER:", starter)

	name, err := input.String("Enter project name: ", false)
	if err != nil {
		return err
	}
	if strings.ContainsRune(name, os.PathSeparator) || name == "." || name == ".." {
		return fmt.Errorf("invalid project name %q", name)
	}

	if err := createProject(starter, name, name); err != nil {
		return err
	}
	underscore := strings.ReplaceAll(name, "-", "_")
	fmt.Printf("replaced %s with %s\n", starterName, name)
	fmt.Printf("replaced %s with %s\n", starterNameUnderscore, underscore)
	msg.Success("Project " + name + " created")
	return nil
}

// createProject копирует starter в dst (без .git), заменяя в содержимом
// файлов имя шаблона на name.
func createProject(starter, dst, name string) error {
	if _, err := os.Stat(starter); err != nil {
		return fmt.Errorf("project starter not found: %w", err)
	}
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("directory %s already exists", dst)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	replacer := strings.NewReplacer(
		starterName, name,
		starterNameUnderscore, strings.ReplaceAll(name, "-", "_"),
	)

	return filepath.WalkDir(starter, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}
		rel, err := filepath.Rel(starter, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.Mkdir(target, info.Mode().Perm())
		}
		if !d.Type().IsRegular() {
			return fmt.Errorf("%s: unsupported file type (symlink?)", path)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if bytes.Contains(data, []byte(starterName)) || bytes.Contains(data, []byte(starterNameUnderscore)) {
			data = []byte(replacer.Replace(string(data)))
		}
		return os.WriteFile(target, data, info.Mode().Perm())
	})
}

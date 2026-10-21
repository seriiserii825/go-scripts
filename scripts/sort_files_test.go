package scripts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSortFiles(t *testing.T) {
	root := t.TempDir()
	write := func(rel, content string) {
		t.Helper()
		path := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("a.jpg", "A")
	write("sub/b.jpg", "B")
	write("sub/deep/c.pdf", "C")
	write("jpg/old.jpg", "OLD") // папка по расширению уже есть
	write("README", "no ext")
	write(".git/config.txt", "hidden")
	write(".env.txt", "hidden file")
	if err := os.MkdirAll(filepath.Join(root, "empty/nested"), 0o755); err != nil {
		t.Fatal(err)
	}

	copied, err := sortFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	if copied != 3 {
		t.Errorf("copied = %d, want 3", copied)
	}

	for rel, want := range map[string]string{
		"jpg/a.jpg":      "A",
		"jpg/b.jpg":      "B",
		"jpg/old.jpg":    "OLD",
		"pdf/c.pdf":      "C",
		"a.jpg":          "A", // оригиналы остаются
		"sub/deep/c.pdf": "C",
	} {
		got, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q, %v; want %q", rel, got, err, want)
		}
	}
	for _, rel := range []string{"txt", "README/README"} {
		if _, err := os.Stat(filepath.Join(root, rel)); err == nil {
			t.Errorf("%s should not exist", rel)
		}
	}

	if err := removeEmptyFolders(root); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "empty")); err == nil {
		t.Error("empty/ should be removed together with empty/nested/")
	}
	if _, err := os.Stat(filepath.Join(root, "sub")); err != nil {
		t.Error("sub/ is not empty and should stay")
	}
}

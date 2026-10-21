package scripts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateProject(t *testing.T) {
	starter := filepath.Join(t.TempDir(), "bs-vite")
	files := map[string]string{
		"style.css":       "Theme Name: bs-vite\n",
		"inc/setup.php":   "function bs_vite_setup() {} // bs-vite",
		"assets/logo.png": "\x89PNG binary",
		".git/config":     "[core]",
		".env":            "PORT=3000",
	}
	for rel, content := range files {
		path := filepath.Join(starter, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	dst := filepath.Join(t.TempDir(), "my-site")
	if err := createProject(starter, dst, "my-site"); err != nil {
		t.Fatal(err)
	}

	for rel, want := range map[string]string{
		"style.css":       "Theme Name: my-site\n",
		"inc/setup.php":   "function my_site_setup() {} // my-site",
		"assets/logo.png": "\x89PNG binary",
		".env":            "PORT=3000",
	} {
		got, err := os.ReadFile(filepath.Join(dst, rel))
		if err != nil || string(got) != want {
			t.Errorf("%s = %q, %v; want %q", rel, got, err, want)
		}
	}
	if _, err := os.Stat(filepath.Join(dst, ".git")); err == nil {
		t.Error(".git should not be copied")
	}

	if err := createProject(starter, dst, "my-site"); err == nil {
		t.Error("second createProject into existing dir: want error")
	}
	if err := createProject(filepath.Join(t.TempDir(), "nope"), filepath.Join(t.TempDir(), "x"), "x"); err == nil {
		t.Error("missing starter: want error")
	}
	if err := createProject(starter, filepath.Join(starter, "inner"), "inner"); err == nil {
		t.Error("dst inside starter: want error")
	}
	if _, err := os.Stat(filepath.Join(starter, "inner")); err == nil {
		t.Error("dst inside starter: nothing should be created")
	}
}

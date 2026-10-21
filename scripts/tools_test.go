package scripts

import (
	"strings"
	"testing"
)

func TestRequire(t *testing.T) {
	if err := require(tool{bin: "sh"}); err != nil {
		t.Errorf("require(sh) = %v; want nil", err)
	}
	if err := require(tool{bin: "no-such-xyz", optional: true}); err != nil {
		t.Errorf("optional missing program: err = %v; want nil", err)
	}

	err := require(
		tool{bin: "sh"},
		tool{bin: "no-such-a", install: "pacman -S a"},
		tool{bin: "no-such-b", install: "npm i -g b"},
	)
	if err == nil {
		t.Fatal("missing required programs: want error")
	}
	for _, want := range []string{"no-such-a: pacman -S a", "no-such-b: npm i -g b"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error has no %q:\n%v", want, err)
		}
	}
	if strings.Contains(err.Error(), "sh:") {
		t.Errorf("installed program in error:\n%v", err)
	}
}

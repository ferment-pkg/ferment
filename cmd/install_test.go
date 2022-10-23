package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func TestParseFpkg(t *testing.T) {
	home := os.Getenv("PWD")
	pkg := parseFpkg(fmt.Sprintf("%s/../Barrells/qemu.fpkg", home))
	if pkg.name != "qemu" {
		t.Fatalf("expected qemu, got %s", pkg.name)
	}
	t.Log(pkg.caveats)
}
func TestExtractFerment(t *testing.T) {
	home := os.Getenv("PWD")
	fs, err := extractFerment(fmt.Sprintf("%s/../test/test.ferment", home))
	if err != nil {
		t.Fatal(err)
	}
	fi, err := afero.ReadDir(fs, "")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fi {
		t.Log(f.Name())
	}
	f, err := fs.Open("test.fpkg")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) < 0 {
		t.Fatal("no bytes read")
	}
}

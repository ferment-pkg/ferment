package cmd

import (
	"fmt"
	"io"
	"net/http"
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
func TestHeadHttp(t *testing.T) {
	r, err := http.Head("https://api.fermentpkg.tech/")
	if err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	if r.StatusCode < 200 || r.StatusCode > 299 {
		t.Fatalf("expected 200-299, got %d", r.StatusCode)
	}
	t.Log(r)

}

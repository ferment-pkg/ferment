package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestParseFpkg(t *testing.T) {
	home := os.Getenv("PWD")
	pkg := parseFpkg(fmt.Sprintf("%s/../Barrells/qemu.fpkg", home))
	if pkg.name != "qemu" {
		t.Fatalf("expected qemu, got %s", pkg.name)
	}
	t.Log(pkg.caveats)
}

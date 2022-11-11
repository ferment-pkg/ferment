package cmd

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/theckman/yacspin"
)

func TestParseFpkg(t *testing.T) {
	home := os.Getenv("PWD")
	pkg := parseFpkg(fmt.Sprintf("%s/../test/test.fpkg", home))
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
	if len(b) == 0 {
		t.Fatal("no bytes read")
	}
}
func TestBuildFromSource(t *testing.T) {
	spinner, err := yacspin.New(yacspin.Config{
		Frequency:     100 * time.Millisecond,
		CharSet:       yacspin.CharSets[9],
		Suffix:        " Building from source",
		StopCharacter: "✓",
	})
	if err != nil {
		t.Fatal(err)
	}
	var l, _ = os.Executable()
	var location = l[:len(l)-len("/ferment")]
	err = os.MkdirAll(fmt.Sprintf("%s/Installed/test", location), 0755)
	if err != nil {
		t.Fatal(err)
	}
	spinner.Start()
	runFpkgCommand("test", "test", `echo hello`, "Build", spinner)
	spinner.Stop()
}
func TestTestFromSource(t *testing.T) {
	spinner, err := yacspin.New(yacspin.Config{
		Frequency:     100 * time.Millisecond,
		CharSet:       yacspin.CharSets[9],
		Suffix:        " Testing",
		StopCharacter: "✓",
	})
	if err != nil {
		t.Fatal(err)
	}
	var l, _ = os.Executable()
	var location = l[:len(l)-len("/ferment")]
	err = os.MkdirAll(fmt.Sprintf("%s/Installed/test", location), 0755)
	if err != nil {
		t.Fatal(err)
	}
	spinner.Start()
	runFpkgCommand("test", "test", `@1=echo e
	match $1 == "e"`, "Test", spinner)
	spinner.Stop()
}
func TestFileMan(t *testing.T) {
	spinner, err := yacspin.New(yacspin.Config{
		Frequency:     100 * time.Millisecond,
		CharSet:       yacspin.CharSets[9],
		Suffix:        " Testing",
		StopCharacter: "✓",
	})
	if err != nil {
		t.Fatal(err)
	}
	var l, _ = os.Executable()
	var location = l[:len(l)-len("/ferment")]
	err = os.MkdirAll(fmt.Sprintf("%s/Installed/test", location), 0755)
	if err != nil {
		t.Fatal(err)
	}
	spinner.Start()
	runFpkgCommand("test", "test", `fileMan write /tmp/test.txt Hello World\n`, "Test", spinner)
	spinner.Stop()
	//check if /tmp/test.txt has the text "Hello World"
	f, err := os.Open("/tmp/test.txt")
	if err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "Hello World\n" {
		t.Fatal("fileMan failed actual:", string(b))
	}

}

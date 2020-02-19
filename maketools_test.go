package maketools

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileExists(t *testing.T) {
	if !FileExists("maketools_test.go") {
		t.Fatalf("File should exist")
	}
	if FileExists("file-that-doesnt-exist") {
		t.Fatalf("File should not exist")
	}
}

func TestExpandUserHome(t *testing.T) {
	actual := ExpandUserHome("~/foo")
	expected := filepath.Join("/home", os.Getenv("USER"), "foo")
	if actual != expected {
		t.Fatalf("Bad user home expansion: %s != %s", actual, expected)
	}
}

func TestFindMakefile(t *testing.T) {
	actual := FindMakefile()
	if actual != "Makefile" {
		t.Fatalf("Found bad makefile: %s", actual)
	}
}

func TestIncludedFiles(t *testing.T) {
	source := `
include Makefile README.md
include LICENSE.txt
	`
	actual, err := IncludedFiles(source)
	expected := []string{"Makefile", "README.md", "LICENSE.txt"}
	if err != nil {
		t.Fatalf("Error parsing included files: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Bad included files: %#v != %#v", actual, expected)
	}
}

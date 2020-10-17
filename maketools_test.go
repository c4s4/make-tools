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

func TestFindMakefile(t *testing.T) {
	actual := FindMakefile()
	if actual != "Makefile" {
		t.Fatalf("Found bad makefile: %s", actual)
	}
}

func TestExpandUserHome(t *testing.T) {
	actual := ExpandUserHome("~/foo")
	expected := filepath.Join("/home", os.Getenv("USER"), "foo")
	if actual != expected {
		t.Fatalf("Bad user home expansion: %s != %s", actual, expected)
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

func TestParseMakefile(t *testing.T) {
	source := `
target:
	echo "target"

target2: # comment
	echo "target2"

target3: target2 # other comment
	echo "target3"

target4: target3
	echo "target4"

target5:target4 target2#another comment
	echo "target5"
`
	expected := []HelpLine{
		{Name: "target", Description: "", Dependencies: nil},
		{Name: "target2", Description: "comment", Dependencies: nil},
		{Name: "target3", Description: "other comment", Dependencies: []string{"target2"}},
		{Name: "target4", Description: "", Dependencies: []string{"target3"}},
		{Name: "target5", Description: "another comment", Dependencies: []string{"target4", "target2"}},
	}
	actual, err := ParseMakefile(source)
	if err != nil {
		t.Fatalf("Error parsing makefile: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("HelpLine not as expected: %v", actual)
	}
}

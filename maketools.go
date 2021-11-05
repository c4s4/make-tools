package maketools

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Filenames is ordered list of makefile names
var Filenames = []string{"GNUmakefile", "makefile", "Makefile"}

// HelpLineRegexp is the regexp to catch target lines
var HelpLineRegexp = regexp.MustCompile(`(?m)^([\w-]+):[\t ]*([^#\n]+)?[\t ]*(#[\t ]*(.*))?$`)

// IncludedRegexp is the regexp to catch included makefiles
var IncludedRegexp = regexp.MustCompile(`(?m)^-?include\s+(.*)$`)

// HelpLine holds information about a line of help
type HelpLine struct {
	Name         string
	Description  string
	Dependencies []string
}

// HelpLineSorter sorts HelpLine by name
type HelpLineSorter []HelpLine

func (a HelpLineSorter) Len() int           { return len(a) }
func (a HelpLineSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a HelpLineSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

// FileExists tells if given file exists:
// - file: the name of the file to test
// Return: a boolean that tells if file exists
func FileExists(file string) bool {
	if stat, err := os.Stat(file); err == nil && !stat.IsDir() {
		return true
	}
	return false
}

// FindMakefile looks for makefile in current directory
// Return: found makefile name
func FindMakefile() string {
	for _, name := range Filenames {
		if FileExists(name) {
			return name
		}
	}
	return ""
}

// ExpandUserHome expand path starting with "~/":
// - path: the path to expand
// Return: expanded path
func ExpandUserHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		currentUser, err := user.Current()
		if err != nil {
			panic("could not get current user")
		}
		home := currentUser.HomeDir
		path = filepath.Join(home, path[2:])
	}
	return path
}

// IncludedFiles extracts included makefiles
// - source: makefile source
// Return: included makefiles
func IncludedFiles(source string) ([]string, error) {
	var included []string
	lines := IncludedRegexp.FindAllStringSubmatch(source, -1)
	for _, line := range lines {
		globs := strings.Split(line[1], " ")
		for _, glob := range globs {
			glob = ExpandUserHome(glob)
			files, err := filepath.Glob(glob)
			if err != nil {
				return nil, err
			}
			included = append(included, files...)
		}
	}
	return included, nil
}

// ReadFile and exit on error
// - filename: name of the file to read
// Return: file content as a string
func ReadFile(filename string) string {
	bytes, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		fmt.Printf("Error reading makefile %s: %v", filename, err)
		os.Exit(1)
	}
	return string(bytes)
}

// ParseMakefile parses passed makefile
// - source: the makefile source
// - recursive: tells if we should parse recursively (defaults to true)
// Return: HelpLine list and error if any
func ParseMakefile(source string, recursive bool) ([]HelpLine, error) {
	result := HelpLineRegexp.FindAllStringSubmatch(source, -1)
	var help []HelpLine
	for _, line := range result {
		dependencies := strings.Split(strings.TrimSpace(line[2]), " ")
		if len(dependencies) == 1 && dependencies[0] == "" {
			dependencies = nil
		}
		helpLine := HelpLine{
			Name:         line[1],
			Description:  line[4],
			Dependencies: dependencies,
		}
		help = append(help, helpLine)
	}
	if recursive {
		filenames, err := IncludedFiles(source)
		if err != nil {
			return nil, fmt.Errorf("parsing included makefile: %v", err)
		}
		for _, filename := range filenames {
			included := ExpandUserHome(filename)
			helpsIncluded, err := ParseMakefile(ReadFile(included), recursive)
			if err != nil {
				return nil, fmt.Errorf("parsing included makefile: %v", err)
			}
			help = append(help, helpsIncluded...)
		}
	}
	sort.Sort(HelpLineSorter(help))
	return help, nil
}

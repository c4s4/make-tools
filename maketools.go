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
var HelpLineRegexp = regexp.MustCompile(`(?m)^([\w-_]+):(\s+([\w\s]+))?(\s+#\s+(.*))?$`)

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
		user, _ := user.Current()
		home := user.HomeDir
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
			for _, file := range files {
				included = append(included, file)
			}
		}
	}
	return included, nil
}

// ParseMakefile parses passed makefile
// - filename: name of found makefile
// Return: HelpLine list and error if any
func ParseMakefile(filename string) ([]HelpLine, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading makefile %s: %v", filename, err)
	}
	source := string(bytes)
	result := HelpLineRegexp.FindAllStringSubmatch(source, -1)
	var help []HelpLine
	for _, line := range result {
		dependencies := strings.Split(strings.TrimSpace(line[3]), " ")
		if len(dependencies) == 1 && dependencies[0] == "" {
			dependencies = nil
		}
		helpLine := HelpLine{
			Name:         line[1],
			Description:  line[5],
			Dependencies: dependencies,
		}
		help = append(help, helpLine)
	}
	filenames, err := IncludedFiles(source)
	if err != nil {
		return nil, fmt.Errorf("parsing included makefile: %v", err)
	}
	for _, filename := range filenames {
		included := ExpandUserHome(filename)
		helpsIncluded, err := ParseMakefile(included)
		if err != nil {
			return nil, fmt.Errorf("parsing included makefile: %v", err)
		}
		for _, helpIncluded := range helpsIncluded {
			help = append(help, helpIncluded)
		}
	}
	sort.Sort(HelpLineSorter(help))
	return help, nil
}

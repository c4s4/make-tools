package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	maketools "github.com/c4s4/make-tools"
)

// Version as printed with -version option
var Version = "UNKNOWN"

const (
	// Help as printed with -help option
	Help = `make-targets [-help] [-version]
Print makefile target list:
-help       To print this help
-version    To print version`
)

// ParseCommandLine parses command line and returns:
// - help: a boolean that tells if we print help
// - version: a boolean that tells if we print version
func ParseCommandLine() (*bool, *bool) {
	help := flag.Bool("help", false, "Print help")
	version := flag.Bool("version", false, "Print version")
	flag.Parse()
	return help, version
}

// TargetListFormatter formats help lines to print
func TargetListFormatter(help []maketools.HelpLine) string {
	targets := make([]string, len(help))
	for i, helpLine := range help {
		targets[i] = helpLine.Name
	}
	sort.Strings(targets)
	return strings.Join(targets, " ")
}

// Error prints an error and exits
// - err: the error to print
func Error(err error) {
	println(fmt.Sprintf("ERROR: %v", err))
	os.Exit(1)
}

func main() {
	help, version := ParseCommandLine()
	if *help {
		fmt.Println(Help)
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	filename := maketools.FindMakefile()
	if filename == "" {
		println("No makefile found")
		os.Exit(1)
	}
	helpLines, err := maketools.ParseMakefile(maketools.ReadFile(filename))
	if err != nil {
		Error(err)
	}
	fmt.Println(TargetListFormatter(helpLines))
}

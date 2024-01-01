package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	maketools "github.com/c4s4/make-tools"
)

// Version as printed with -version option
var Version = "UNKNOWN"

const (
	// Help as printed with -help option
	Help = `make-help [-help] [-version]
Print makefile help:
-help       To print this help
-version    To print version
-root       To parse root makefile only
-mute       Don't print targets without comment`
)

// ParseCommandLine parses command line and returns:
// - help: tells if we print help
// - version: tells if we print version
// - root: tells if we parse recursively
// - mute: tells not to print targets without comment
func ParseCommandLine() (help, version, root, mute *bool) {
	help = flag.Bool("help", false, "Print help")
	version = flag.Bool("version", false, "Print version")
	root = flag.Bool("root", false, "Parse root makefile only")
	mute = flag.Bool("mute", false, "Don't print targets without comment")
	flag.Parse()
	return
}

// HelpLineFormatter formats help lines to print
func HelpLineFormatter(help []maketools.HelpLine, mute bool) string {
	indent := 0
	for _, helpLine := range help {
		if indent < len(helpLine.Name) {
			indent = len(helpLine.Name)
		}
	}
	var lines []string
	for _, helpLine := range help {
		if helpLine.Description == "" && mute {
			continue
		}
		spaces := indent - len(helpLine.Name)
		line := "\033[93m" + helpLine.Name + "\033[0m" + strings.Repeat(" ", spaces)
		if helpLine.Description != "" {
			line += " " + helpLine.Description
		}
		if len(helpLine.Dependencies) > 0 {
			dependencies := strings.Join(helpLine.Dependencies, ", ")
			line += " [" + dependencies + "]"
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Error prints an error and exits
// - err: the error to print
func Error(err error) {
	println(fmt.Sprintf("ERROR: %v", err))
	os.Exit(1)
}

func main() {
	help, version, root, mute := ParseCommandLine()
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
	helpLines, err := maketools.ParseMakefile(maketools.ReadFile(filename), !*root)
	if err != nil {
		Error(err)
	}
	fmt.Println(HelpLineFormatter(helpLines, *mute))
}

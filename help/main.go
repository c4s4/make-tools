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

// HelpLineFormatter formats help lines to print
func HelpLineFormatter(help []maketools.HelpLine) string {
	indent := 0
	for _, helpLine := range help {
		if indent < len(helpLine.Name) {
			indent = len(helpLine.Name)
		}
	}
	var lines []string
	for _, helpLine := range help {
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
	fmt.Println(HelpLineFormatter(helpLines))
}

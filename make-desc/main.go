package main

import (
	"flag"
	"fmt"
	"os"

	maketools "github.com/c4s4/make-tools"
)

// Version as printed with -version option
var Version = "UNKNOWN"

const (
	// Help as printed with -help option
	Help = `make-desc [-help] [-version] target
Print makefile help:
-help       To print this help
-version    To print version
target      Name of the target to describe`
)

// ParseCommandLine parses command line and returns:
// - help: tells if we print help
// - version: tells if we print version
// - string: name of the target to describe
func ParseCommandLine() (help, version *bool, target string) {
	help = flag.Bool("help", false, "Print help")
	version = flag.Bool("version", false, "Print version")
	flag.Parse()
	args := flag.Args()
	if len(args) == 1 {
		target = args[0]
	}
	return
}

// Error prints an error and exits
// - err: the error to print
func Error(err error) {
	println(fmt.Sprintf("ERROR: %v", err))
	os.Exit(1)
}

// GetDescription returns description of given target:
// - helpLines: list of help lines for the makefile
// - target: name of the target to describe
// Return: description of target or ""
func GetDescription(helpLines []maketools.HelpLine, target string) string {
	for _, helpLine := range helpLines {
		if helpLine.Name == target {
			return helpLine.Description
		}
	}
	return ""
}

func main() {
	help, version, target := ParseCommandLine()
	if *help {
		fmt.Println(Help)
		os.Exit(0)
	}
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if target == "" {
		println("You must pass target on command line")
		os.Exit(1)
	}
	filename := maketools.FindMakefile()
	if filename == "" {
		println("No makefile found")
		os.Exit(1)
	}
	helpLines, err := maketools.ParseMakefile(maketools.ReadFile(filename), true)
	if err != nil {
		Error(err)
	}
	fmt.Println(GetDescription(helpLines, target))
}

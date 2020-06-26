# Make Tools

Tools to pimp your Makefiles.

## Installation

### Unix users (Linux, BSDs and MacOSX)

Unix users may download and install latest *make-tools* release with command:

```bash
sh -c "$(curl https://sweetohm.net/dist/make-tools/install)"
```

If *curl* is not installed on you system, you might run:

```bash
sh -c "$(wget -O - https://sweetohm.net/dist/make-tools/install)"
```

**Note:** Some directories are protected, even as *root*, on **MacOSX** (since *El Capitan* release), thus you can't install *project* in */usr/bin* for instance.

### Binary package

Otherwise, you can download latest binary archive at <https://github.com/c4s4/make-tools/releases>. Unzip the archive, put the binary of your platform somewhere in your *PATH* and rename it *project*.

## Usage

There are two tools in this toolbox:

### make-targets

This tool lists targets available in current makefile and included ones recursively. This is called to perform Bash completion.

To enable Bash target completion on make, source following file:

```
complete -W "\`make-targets\`" make
```

### make-help

This tool prints help on targets in current makefile and included ones recursively. Thus to get help about makefile in current directory, you might type:

```
$ make-help
build   Build binary [clean]
clean   Clean generated files and test cache
fmt     Format Go source code
help    Print help on Makefile
install [build]
run     Run make help [build]
```

You might also print help calling following target *help*:

```
.PHONY: help
help: # Print help on Makefile
	@make-help
```

## make-desc

Describe given target. If makefile defines a *build* target, you might describe it with:

```
$ make-desc build
Build binary
```

This prints the target description on command line.

*Enjoy!*

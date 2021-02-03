# Parent Makefiles https://github.com/c4s4/make

include ~/.make/Golang.mk

go-build: build
go-binaries: binaries
test: go-test # Run unit tests
release: go-release # Perform release (you must pass VERSION=X.Y.Z on command line)

build: clean # Build binary for current platform
	$(title)
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/ ./...
	@cd $(BUILD_DIR); \
	for file in *; do \
		mv $$file make-$$file; \
	done

binaries: clean # Build binaries for all platforms
	$(title)
	@mkdir -p $(BUILD_DIR)/bin
	@gox -ldflags "-X main.Version=$(VERSION) -s -f" -osarch '$(GOOSARCH)' -output=$(BUILD_DIR)/bin/make-{{.Dir}}-{{.OS}}-{{.Arch}} ./...

# Parent Makefiles https://github.com/c4s4/make

include ~/.make/Golang.mk

go-build: clean # Build binary
	$(title)
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/ ./...
	@cd $(BUILD_DIR); \
	for file in *; do \
		mv $$file make-$$file; \
	done

go-binaries: clean # Build binaries
	$(title)
	@mkdir -p $(BUILD_DIR)/bin
	@gox -ldflags "-X main.Version=$(VERSION) -s -f" -osarch '$(GOOSARCH)' -output=$(BUILD_DIR)/bin/make-{{.Dir}}-{{.OS}}-{{.Arch}} ./...

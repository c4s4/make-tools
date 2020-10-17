# Parent makefile https://github.com/c4s4/make

include ~/.make/color.mk
include ~/.make/help.mk
include ~/.make/git.mk

BUILD_DIR = build
VERSION = 1.1.2

.DEFAULT_GOAL :=
default: fmt clean test build

clean: # Clean generated files and test cache
	$(title)
	@rm -rf $(BUILD_DIR)
	@go clean -testcache

fmt: # Format Go source code
	$(title)
	@go fmt ./...

test: # Run tests
	$(title)
	@go test -cover ./...

.PHONY: build
build: clean # Build binary
	$(title)
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/make-help help/main.go
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/make-targets targets/main.go
	@go build -ldflags "-X main.Version=$(VERSION) -s -f" -o $(BUILD_DIR)/make-desc desc/main.go

binaries: clean # Build binaries
	$(title)
	@mkdir -p $(BUILD_DIR)/bin
	@gox -ldflags "-X main.Version=$(VERSION) -s -f" -osarch '!darwin/386' -output=$(BUILD_DIR)/bin/make-{{.Dir}}-{{.OS}}-{{.Arch}} ./...

install: build # Install binaries in GOPATH
	$(title)
	@cp $(BUILD_DIR)/make-* $$GOPATH/bin/

deploy: binaries # Deploy binaries on server
	$(title)
	@scp install $(BUILD_DIR)/bin/* casa@sweetohm.net:/home/web/dist/make-tools/

documentation: # Generate documentation
	$(title)
	@mkdir -p $(BUILD_DIR)
	@cp LICENSE.txt $(BUILD_DIR)
	@md2pdf -o $(BUILD_DIR)/README.pdf README.md

archive: binaries documentation # Build distribution archive
	$(title)
	@mkdir -p $(BUILD_DIR)/make-tools
	@mv $(BUILD_DIR)/bin $(BUILD_DIR)/make-tools
	@mv $(BUILD_DIR)/README.pdf $(BUILD_DIR)/LICENSE.txt $(BUILD_DIR)/make-tools
	@cd $(BUILD_DIR) && tar cvf make-tools-$(VERSION).tar make-tools/ && gzip make-tools-$(VERSION).tar

tag: # Tag project
	$(title)
	@git diff-index --quiet HEAD -- || (echo "$(RED)ERROR$(END) There are uncommitted changes" && exit 1)
	@test `git rev-parse --abbrev-ref HEAD)` = 'master' || (echo "$(RED)ERROR$(END) You are not on branch master" && exit 1)
	@git tag -a $(VERSION) -m  "Release $(TAG)"
	@git push origin $(TAG)

release: test deploy archive tag # Perform a release
	@echo "$(GRE)OK$(EBD) Release done!"

run: build # Run make tools
	$(title)
	@$(BUILD_DIR)/make-help
	@$(BUILD_DIR)/make-targets
	@$(BUILD_DIR)/make-desc build

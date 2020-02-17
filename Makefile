include ~/.make/*.mk

BUILD_DIR=build

clean: # Clean generated files and test cache
	@rm -rf $(BUILD_DIR)
	@go clean -testcache

fmt: # Format Go source code
	@go fmt ./...

.PHONY: build
build: clean # Build binary
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-X main.Version=$(COMMIT) -s -f" -o $(BUILD_DIR)/make-help help/main.go
	@go build -ldflags "-X main.Version=$(COMMIT) -s -f" -o $(BUILD_DIR)/make-targets targets/main.go

binaries: clean # Build binaries
	@echo "$(YEL)Building binaries...$(END)"
	@mkdir -p $(BUILD_DIR)/bin
	@gox -ldflags "-s -f" -output=$(BUILD_DIR)/bin/{{.Dir}}-{{.OS}}-{{.Arch}} ./...
	@cd $(BUILD_DIR)/bin && for file in *; do mv "$$file" "make-$$file"; done

install: build # Install binaries in GOPATH
	@echo "$(YEL)Installing binaries in GOPATH$(END)"
	@cp $(BUILD_DIR)/make-* $$GOPATH/bin/

deploy: binaries # Deploy binaries on server
	@echo "$(YEL)Deploying binaries on server...$(END)"
	@scp install $(BUILD_DIR)/bin/* casa@sweetohm.net:/home/web/dist/make-tools/

documentation: # Generate documentation
	@echo "$(YEL)Generating documentation$(END)"
	@mkdir -p $(BUILD_DIR)
	@cp LICENSE.txt $(BUILD_DIR)
	@md2pdf -o $(BUILD_DIR)/README.pdf README.md

archive: binaries documentation # Build distribution archive
	@echo "$(YEL)Building distribution archive$(END)"
	@mkdir -p $(BUILD_DIR)/make-tools
	@mv $(BUILD_DIR)/bin $(BUILD_DIR)/make-tools
	@mv $(BUILD_DIR)/README.pdf $(BUILD_DIR)/LICENSE.txt $(BUILD_DIR)/make-tools
	@cd $(BUILD_DIR) && tar cvf make-tools.tar make-tools/ && gzip make-tools.tar

run: build # Run make tools
	@echo "$(YEL)Running make tools$(END)"
	@$(BUILD_DIR)/make-help
	@$(BUILD_DIR)/make-targets

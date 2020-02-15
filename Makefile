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
	@rename s/help/make-help/ $(BUILD_DIR)/bin/*
	@rename s/targets/make-targets/ $(BUILD_DIR)/bin/*

install: build # Install binaries in GOPATH
	@cp $(BUILD_DIR)/make-* $$GOPATH/bin/

deploy: binaries # Deploy binaries on server
	@scp install $(BUILD_DIR)/bin/* casa@sweetohm.net:/home/web/dist/make-tools/

run: build # Run make help
	@$(BUILD_DIR)/make-help
	@$(BUILD_DIR)/make-targets

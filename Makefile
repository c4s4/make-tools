include *.mk

BUILD_DIR=build

clean: # Clean generated files and test cache
	@rm -rf $(BUILD_DIR)
	@go clean -testcache

fmt: # Format Go source code
	@go fmt ./...

.PHONY: build
build: clean # Build binary
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags "-s -f" -o $(BUILD_DIR)/make-help help/main.go
	@go build -ldflags "-s -f" -o $(BUILD_DIR)/make-targets targets/main.go

run: build # Run make help
	@$(BUILD_DIR)/make-help
	@$(BUILD_DIR)/make-targets

install: build
	@cp $(BUILD_DIR)/make-* $$GOPATH/bin/

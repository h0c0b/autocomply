.PHONY: all darwin_amd64 darwin_arm64 linux_amd64 linux_arm64 clean docs

BINARY_NAME=autocomply
SOURCE_FILE=./main.go
BUILD_DIR=build
DOC_SRC_DIR=./output

all: clean darwin_amd64 darwin_arm64 linux_amd64 linux_arm64

darwin_amd64:
	@echo "Building for Darwin (AMD64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 $(SOURCE_FILE)

darwin_arm64:
	@echo "Building for Darwin (ARM64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 $(SOURCE_FILE)

linux_amd64:
	@echo "Building for Linux (AMD64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 $(SOURCE_FILE)

linux_arm64:
	@echo "Building for Linux (ARM64)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 $(SOURCE_FILE)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

docs:
	@echo "Generating PDF documents..."
	@if [ -d "$(DOC_SRC_DIR)" ] && [ -n "$(wildcard $(DOC_SRC_DIR)/*.md)" ]; then \
		for file in $(DOC_SRC_DIR)/*.md; do \
			pandoc "$$file" -o "$${file%.md}.pdf"; \
		done; \
	else \
		echo "No Markdown files found in $(DOC_SRC_DIR)"; \
	fi
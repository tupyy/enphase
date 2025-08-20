# Simple Makefile for enphase-cli - Linux only

BINARY_NAME=enphase-cli
BUILD_DIR=./bin

# Build for Linux x86_64
.PHONY: build
build:
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run the app
.PHONY: run
run: build
	$(BUILD_DIR)/$(BINARY_NAME)
# Simple Makefile for enphase-cli - Linux only

BINARY_NAME=enphase-cli
BUILD_DIR=./bin
IMAGE_NAME=quay.io/ctupangiu/enphase
TAG?=latest

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

# Build container image
.PHONY: container-build
container-build:
	podman build -t $(IMAGE_NAME):$(TAG) -f Containerfile .

# Push container image to registry
.PHONY: container-push
container-push: container-build
	podman push $(IMAGE_NAME):$(TAG)

# Build and push container image
.PHONY: container-release
container-release: container-push
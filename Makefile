# Define the Go command and output binary name
GO=go
GOFMT=gofmt
BIN_NAME := claudio
TARGET := $(BIN_NAME)

.PHONY: all build run clean test help build-linux-amd64 build-darwin-amd64 build-darwin-arm64 build-all setup

# Default target: builds the application
all: build

## build: build the Go application
build: $(TARGET)

$(TARGET): cmd/aicli/main.go
	$(GO) build -o $(TARGET) ./cmd/aicli/main.go

## build-linux-amd64: build for Linux x64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_NAME)-linux-amd64 ./cmd/aicli

## build-darwin-amd64: build for macOS Intel
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BIN_NAME)-darwin-amd64 ./cmd/aicli

## build-darwin-arm64: build for macOS Silicon
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -o $(BIN_NAME)-darwin-arm64 ./cmd/aicli

## build-all: build for all platforms
build-all: build-linux-amd64 build-darwin-amd64 build-darwin-arm64

## run: build and run the application
run: build
	./$(TARGET)

## clean: remove the generated executable
clean:
	rm -f $(TARGET)

## test: run Go tests
test:
	$(GO) test ./...

## setup: install development tools and configure git hooks
setup:
	@which lefthook > /dev/null 2>&1 || (echo "lefthook not found. Install with: brew install lefthook" && exit 1)
	@which gitleaks > /dev/null 2>&1 || (echo "gitleaks not found. Install with: brew install gitleaks" && exit 1)
	lefthook install
	@echo "Git hooks configured successfully."

## help: display this help message
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

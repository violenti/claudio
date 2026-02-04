# Define the Go command and output binary name
GO=go
GOFMT=gofmt
BIN_NAME := claudio
TARGET := $(BIN_NAME)

.PHONY: all build run clean test help

# Default target: builds the application
all: build

## build: build the Go application
build: $(TARGET)

$(TARGET): cmd/aicli/main.go
	$(GO) build -o $(TARGET) ./cmd/aicli/main.go

## run: build and run the application
run: build
	./$(TARGET)

## clean: remove the generated executable
clean:
	rm -f $(TARGET)

## test: run Go tests
test:
	$(GO) test ./...

## help: display this help message
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

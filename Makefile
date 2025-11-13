# Go MCP Git Server Makefile

# Variables
BINARY_NAME=go-mcp-git
BINARY_PATH=./bin/$(BINARY_NAME)
MAIN_PATH=.
GO_FILES=$(shell find . -name "*.go" -type f)

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

# Default target
.PHONY: all
all: clean test build

# Build the binary
.PHONY: build
build: $(BINARY_PATH)

$(BINARY_PATH): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	go build $(BUILD_FLAGS) $(LDFLAGS) -o $(BINARY_PATH) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p bin
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	golangci-lint run

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the server (for development)
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH) --verbose

# Install the binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	go install $(BUILD_FLAGS) $(LDFLAGS) $(MAIN_PATH)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, test, and build"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  deps         - Install dependencies"
	@echo "  run          - Build and run the server"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  help         - Show this help"

.PHONY: build install clean test

BINARY_NAME=go-start
BUILD_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/go-start
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install ./cmd/go-start
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

lint:
	@echo "Linting code..."
	@golangci-lint run ./...

deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

example:
	@echo "Creating example project..."
	@./$(BUILD_DIR)/$(BINARY_NAME) create example-api
	@cd example-api && go mod tidy

help:
	@echo "Available targets:"
	@echo "  build          - Build the CLI tool"
	@echo "  install        - Install the CLI tool to GOPATH/bin"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  deps           - Install dependencies"
	@echo "  example        - Create an example project"
	@echo "  doctor         - Diagnose local env and workspace"
	@echo "  help           - Show this help message"

doctor:
	@$(MAKE) build >/dev/null
	@./bin/go-start doctor

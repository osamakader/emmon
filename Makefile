# Makefile for Embedded Linux Monitor (emmon)

# Variables
BINARY_NAME=emmon
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
all: clean build

# Build the application
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for multiple platforms
build-all: build-linux build-arm64 build-arm

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 .

# Build for ARM64 (Raspberry Pi 3/4, etc.)
build-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 .

# Build for ARM32 (Raspberry Pi 1/2, etc.)
build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-linux-arm .

# Build for macOS
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 .

# Build for Windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with race detection
test-race:
	$(GOTEST) -race -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run

# Install linter
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run the web interface
run-web:
	./$(BINARY_NAME) web

# Run the terminal interface
run-terminal:
	./$(BINARY_NAME) terminal

# Run with debug logging
run-web-debug:
	./$(BINARY_NAME) web --log-level debug

# Create release packages
release: build-all
	mkdir -p releases
	tar -czf releases/$(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64
	tar -czf releases/$(BINARY_NAME)-linux-arm64.tar.gz $(BINARY_NAME)-linux-arm64
	tar -czf releases/$(BINARY_NAME)-linux-arm.tar.gz $(BINARY_NAME)-linux-arm
	zip releases/$(BINARY_NAME)-darwin-amd64.zip $(BINARY_NAME)-darwin-amd64
	zip releases/$(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe

# Install to system
install: build
	sudo cp $(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)

# Uninstall from system
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Create systemd service file
install-service: install
	@echo "Installing systemd service..."
	@sudo cp scripts/$(BINARY_NAME).service /etc/systemd/system/
	@echo "Enabling and starting service..."
	@sudo systemctl daemon-reload
	@sudo systemctl enable $(BINARY_NAME)
	@sudo systemctl start $(BINARY_NAME)
	@echo "Service installed and started. Check status with: sudo systemctl status $(BINARY_NAME)"

# Remove systemd service
uninstall-service:
	@echo "Stopping and disabling service..."
	@sudo systemctl stop $(BINARY_NAME) || true
	@sudo systemctl disable $(BINARY_NAME) || true
	@sudo rm -f /etc/systemd/system/$(BINARY_NAME).service
	@sudo systemctl daemon-reload
	@echo "Service removed."

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for all platforms"
	@echo "  build-linux    - Build for Linux AMD64"
	@echo "  build-arm64    - Build for Linux ARM64"
	@echo "  build-arm      - Build for Linux ARM32"
	@echo "  build-mac      - Build for macOS"
	@echo "  build-windows  - Build for Windows"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  deps           - Install dependencies"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  install-lint   - Install linter"
	@echo "  run-web        - Run web interface"
	@echo "  run-terminal   - Run terminal interface"
	@echo "  run-web-debug  - Run web interface with debug logging"
	@echo "  release        - Create release packages"
	@echo "  install        - Install to system"
	@echo "  uninstall      - Uninstall from system"
	@echo "  install-service - Install as systemd service"
	@echo "  uninstall-service - Remove systemd service"
	@echo "  help           - Show this help"

.PHONY: all build build-all build-linux build-arm64 build-arm build-mac build-windows clean test test-race test-coverage deps fmt lint install-lint run-web run-terminal run-web-debug release install uninstall install-service uninstall-service help 
.PHONY: build run clean install test fmt vet

# Binary name
BINARY_NAME=gotui

# Build the application
build:
	go build -o $(BINARY_NAME) .

# Run the application
run: build
	./$(BINARY_NAME)

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the binary to GOPATH/bin
install:
	go install .

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o $(BINARY_NAME)-windows-amd64.exe .

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the application"
	@echo "  run        - Build and run the application"
	@echo "  clean      - Remove build artifacts"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  install    - Install the binary to GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  check      - Run fmt, vet, and test"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  help       - Show this help message"

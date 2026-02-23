# DBZ Makefile

# Variables
BINARY_NAME=dbz
INSTALL_PATH=/usr/local/bin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build clean test deps install uninstall release help

# Default target
all: test build

## help: Show this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the binary for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go
	@echo "✅ Build complete: ./$(BINARY_NAME)"

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "✅ Clean complete"

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

## install: Install binary to system
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installation complete"
	@echo "Run 'dbz --help' to get started"

## uninstall: Remove binary from system
uninstall:
	@echo "Removing $(BINARY_NAME) from $(INSTALL_PATH)..."
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Uninstall complete"

## run: Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

## dev: Run with development settings
dev:
	@echo "Running in development mode..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go
	./$(BINARY_NAME) --help

## fmt: Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

## lint: Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

## vet: Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

## release: Build for all platforms
release: clean
	@echo "Building releases for all platforms..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		$(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-$${platform%/*}-$${platform#*/} main.go; \
		echo "✅ Built for $$platform"; \
	done
	@echo "📦 Release builds complete in dist/"

## docker-build: Build using Docker
docker-build:
	@echo "Building with Docker..."
	docker run --rm -v $(PWD):/workspace -w /workspace golang:1.21 \
		make build

## install-script: Make install script executable
install-script:
	chmod +x scripts/install.sh

## serve-install: Serve install script locally for testing
serve-install:
	@echo "Serving install script at http://localhost:8080/install.sh"
	@echo "Test with: curl -fsSL http://localhost:8080/install.sh | bash"
	python3 -m http.server 8080 --directory scripts

# Development helpers
## create-postgres: Create PostgreSQL for testing
create-postgres:
	./$(BINARY_NAME) create postgres --port 5432

## create-mysql: Create MySQL for testing
create-mysql:
	./$(BINARY_NAME) create mysql --port 3306

## list: List all databases
list:
	./$(BINARY_NAME) list

## cleanup: Clean up test databases
cleanup:
	./$(BINARY_NAME) delete postgres || true
	./$(BINARY_NAME) delete mysql || true
	rm -f *.db

# CI/CD helpers
## ci-test: Run tests for CI
ci-test:
	$(GOTEST) -race -coverprofile=coverage.out ./...

## ci-build: Build for CI
ci-build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) main.go

## coverage: Generate coverage report
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Version info
version:
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"

# Default help target
.DEFAULT_GOAL := help
# Go parameters
BINARY_NAME=eth-parser
MAIN_PACKAGE=./cmd/main.go

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
BUILD_FLAGS=-v

# Test flags
TEST_FLAGS=-v
COVERAGE_FLAGS=-coverprofile=coverage.out

# Directories
CMD_DIR=./cmd
INTERNAL_DIR=./internal
PKG_DIR=./pkg

.PHONY: all build clean test coverage deps run help

all: clean deps build test

build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PACKAGE)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f coverage.out

test:
	$(GOTEST) $(TEST_FLAGS) ./...

coverage:
	$(GOTEST) $(COVERAGE_FLAGS) ./...
	go tool cover -html=coverage.out

coverage-func:
	$(GOTEST) $(COVERAGE_FLAGS) ./...
	go tool cover -func=coverage.out

deps:
	$(GOMOD) download
	$(GOMOD) tidy

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PACKAGE)
	./$(BINARY_NAME)

test-api:
	$(GOTEST) $(TEST_FLAGS) ./internal/api

test-parser:
	$(GOTEST) $(TEST_FLAGS) ./internal/parser

test-storage:
	$(GOTEST) $(TEST_FLAGS) ./internal/storage

test-rpc:
	$(GOTEST) $(TEST_FLAGS) ./internal/rpc

lint:
	golangci-lint run

help:
	@echo "Available commands:"
	@echo "  make all          - Clean, download dependencies, build, and test"
	@echo "  make build        - Build the application"
	@echo "  make clean        - Clean build files"
	@echo "  make test         - Run all tests"
	@echo "  make coverage     - Generate test coverage report and open in browser"
	@echo "  make coverage-func- Show test coverage report in terminal"
	@echo "  make deps         - Download dependencies"
	@echo "  make run          - Build and run the application"
	@echo "  make test-api     - Run API tests"
	@echo "  make test-parser  - Run parser tests"
	@echo "  make test-storage - Run storage tests"
	@echo "  make test-rpc     - Run RPC tests"
	@echo "  make lint         - Run linter"
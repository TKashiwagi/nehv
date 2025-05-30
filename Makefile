# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=configure
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build clean test deps tidy help

all: clean build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) ./cmd

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)
	rm -f .nehv_configure_history
	rm -f *.tmp *.temp *.log

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

deps:
	$(GOGET) -v ./...

tidy:
	$(GOMOD) tidy

help:
	@echo "Available targets:"
	@echo "  all            - Clean and build the project"
	@echo "  build          - Build the project"
	@echo "  clean          - Remove binary files and temporary files"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy up dependencies"
	@echo "  help           - Show this help message" 
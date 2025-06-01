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

# Detect WSL environment
IS_WSL := $(shell uname -a | grep -i microsoft > /dev/null && echo true || echo false)

.PHONY: all build build-linux clean test deps tidy help

all: clean build

build: build-linux
ifeq ($(IS_WSL),false)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) ./cmd
endif

build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o /tmp/$(BINARY_NAME) ./cmd
	mv /tmp/$(BINARY_NAME) .

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
	@echo "  build          - Build Windows and Linux binaries"
	@echo "  build-linux    - Build Linux binary only"
	@echo "  clean          - Remove binary files and temporary files"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy up dependencies"
	@echo "  help           - Show this help message" 
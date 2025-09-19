# Makefile for Logic MCP Server

.PHONY: build test docker clean help run-stdio run-http test-integration deps lint format

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=logic-mcp
BINARY_PATH=./cmd/server

# Docker parameters
DOCKER_IMAGE=logic-mcp
DOCKER_TAG=latest

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	$(GOBUILD) -o $(BINARY_NAME) -v $(BINARY_PATH)

test: ## Run unit tests
	$(GOTEST) -v ./internal/...

test-integration: ## Run integration tests (requires SWI-Prolog)
	@echo "Running integration tests..."
	@if ! command -v swipl > /dev/null; then \
		echo "Error: SWI-Prolog (swipl) not found. Please install SWI-Prolog to run integration tests."; \
		echo "  Ubuntu/Debian: sudo apt-get install swi-prolog"; \
		echo "  macOS: brew install swi-prolog"; \
		echo "  Alpine: apk add swi-prolog"; \
		exit 1; \
	fi
	$(GOTEST) -v ./test/... -tags=integration

test-all: test test-integration ## Run all tests

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run-stdio: docker-build ## Run Docker container in STDIO mode
	docker run --rm -i $(DOCKER_IMAGE):$(DOCKER_TAG) -mode stdio

docker-run-http: docker-build ## Run Docker container in HTTP mode
	docker run --rm -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG) -mode http -port 8080

run-stdio: build ## Run in STDIO mode
	./$(BINARY_NAME) -mode stdio

run-http: build ## Run in HTTP mode
	./$(BINARY_NAME) -mode http -port 8080

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

# Development helpers
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@if ! command -v swipl > /dev/null; then \
		echo "Warning: SWI-Prolog not found. Please install it for full functionality."; \
		echo "  Ubuntu/Debian: sudo apt-get install swi-prolog"; \
		echo "  macOS: brew install swi-prolog"; \
		echo "  Alpine: apk add swi-prolog"; \
	fi
	$(GOMOD) download

example-basic: build ## Run basic Prolog examples
	@echo "Loading basic examples into Logic MCP..."
	@echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"prolog_load_facts","arguments":{"facts":"animal(cat).\\nanimal(dog).\\nhas_fur(cat).\\nhas_fur(dog).\\nmammal(X) :- animal(X), has_fur(X)."}}}' | ./$(BINARY_NAME) -mode stdio

example-family: build ## Run family tree example
	@echo "Loading family tree example into Logic MCP..."
	@echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"prolog_query","arguments":{"query":"parent(john, bob)."}}}' | ./$(BINARY_NAME) -mode stdio

lint: ## Run linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

format: ## Format code
	go fmt ./...

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Testing specific components
test-prolog: ## Test only prolog engine
	$(GOTEST) -v ./internal/prolog/...

test-mcp: ## Test only MCP server
	$(GOTEST) -v ./internal/mcp/...

test-tools: ## Test only logic tools
	$(GOTEST) -v ./internal/tools/...

# Quick development cycle
dev: clean build test ## Clean, build and test in one command

# Check if everything is ready for production
check: deps lint test test-integration ## Full check before deployment
	@echo "✅ All checks passed! Ready for deployment."
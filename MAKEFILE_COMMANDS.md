# Makefile Commands Reference

This document describes all available Makefile commands for the Logic MCP Server project.

## Basic Commands

### `make help`
Shows all available commands with descriptions.

### `make build`
Builds the binary executable (`logic-mcp`) from source code.

### `make clean`
Removes build artifacts and the binary file.

### `make deps`
Downloads and tidies Go module dependencies.

## Testing Commands

### `make test`
Runs unit tests for all internal packages.

### `make test-integration`
Runs integration tests that require SWI-Prolog to be installed.
**Requires:** SWI-Prolog (`swipl` command available)

### `make test-all`
Runs both unit and integration tests.

### Component-specific tests:
- `make test-prolog` - Tests only the Prolog engine
- `make test-mcp` - Tests only the MCP server  
- `make test-tools` - Tests only the logic tools

## Running the Server

### `make run-stdio`
Builds and runs the server in STDIO mode (for MCP clients like VSCode Copilot).

### `make run-http`
Builds and runs the server in HTTP mode on port 8080.

## Docker Commands

### `make docker-build`
Builds the Docker image with tag `logic-mcp:latest`.

### `make docker-run-stdio`
Builds and runs the Docker container in STDIO mode.

### `make docker-run-http`
Builds and runs the Docker container in HTTP mode, exposing port 8080.

## Development Commands

### `make dev-setup`
Sets up the development environment by checking for SWI-Prolog and downloading dependencies.

### `make dev`
Quick development cycle: clean, build, and test in one command.

### `make format`
Formats all Go source code using `go fmt`.

### `make lint`
Runs golangci-lint if available, otherwise shows installation instructions.

### `make install-tools`
Installs development tools like golangci-lint.

## Example Commands

### `make example-basic`
Loads basic Prolog examples (animals, mammals) into the MCP server.

### `make example-family`
Demonstrates family tree relationships using the MCP server.

## Quality Assurance

### `make check`
Runs a comprehensive check including:
- Dependency management
- Linting
- Unit tests  
- Integration tests

This is the command to run before deployment to ensure everything is working correctly.

## Example Usage

```bash
# Complete development setup
make dev-setup
make install-tools

# Development cycle
make dev           # Build and test
make lint          # Check code quality
make format        # Format code

# Testing
make test-all      # Run all tests
make check         # Full pre-deployment check

# Running
make run-stdio     # For MCP clients
make run-http      # For web integration

# Docker deployment
make docker-build
make docker-run-stdio
```

## Requirements

- **Go 1.23+**: Required for building
- **SWI-Prolog**: Required for integration tests and runtime
- **Docker**: Required for containerized deployment
- **golangci-lint**: Optional, for code linting

## Installation Instructions

### SWI-Prolog
- **macOS**: `brew install swi-prolog`
- **Ubuntu/Debian**: `sudo apt-get install swi-prolog`
- **Alpine**: `apk add swi-prolog`

### golangci-lint
```bash
make install-tools
# or manually:
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```
# Logic MCP Server

MCP (Model Context Protocol) server that enhances Language Learning Models with powerful logic reasoning capabilities using SWI-Prolog. This server provides tools for evaluating logic problems, defining rules, and validating logical statements.

## Features

- 🧠 **Logic Engine**: Powered by SWI-Prolog for robust logical reasoning
- 🔧 **MCP Tools**: Comprehensive set of tools for logic operations
- 📚 **Examples**: Rich collection of Prolog examples and tutorials  
- ✅ **Validation**: Verbose input validation and error handling
- 🌐 **Dual Interface**: Support for both STDIO and HTTP communication
- 🐳 **Docker Ready**: Complete containerization with SWI-Prolog
- 🧪 **Tested**: Comprehensive integration tests

## Quick Start

### Using Docker (Recommended)

```bash
# Build the Docker image
make docker-build

# Run in STDIO mode (for MCP clients like VSCode Copilot)
make docker-run-stdio

# Run in HTTP mode (for web integration)
make docker-run-http
```

### Local Development

```bash
# Setup development environment
make dev-setup

# Build the binary
make build

# Run tests
make test

# Run in STDIO mode
make run-stdio

# Run in HTTP mode  
make run-http
```

### VSCode Copilot Integration

1. Build the Docker image: `make docker-build`
2. Configure MCP in your VSCode settings using the provided config in `.vscode/mcp-config.json`

## Available Tools

The Logic MCP Server provides the following tools:

### `prolog_query`
Execute Prolog queries and return results.

**Example:**
```json
{
  "name": "prolog_query", 
  "arguments": {
    "query": "member(X, [1,2,3])."
  }
}
```

### `prolog_load_facts`
Load Prolog facts and rules into the knowledge base.

**Example:**
```json
{
  "name": "prolog_load_facts",
  "arguments": {
    "facts": "parent(tom, bob).\nparent(bob, pat).\nancestor(X, Y) :- parent(X, Y).\nancestor(X, Y) :- parent(X, Z), ancestor(Z, Y)."
  }
}
```

### `prolog_validate_syntax`
Validate Prolog syntax without executing.

**Example:**
```json
{
  "name": "prolog_validate_syntax",
  "arguments": {
    "code": "animal(cat).\nmammal(X) :- animal(X), has_fur(X)."
  }
}
```

### `prolog_clear_kb`
Clear the Prolog knowledge base.

**Example:**
```json
{
  "name": "prolog_clear_kb",
  "arguments": {}
}
```

### `prolog_solve_problem`
Solve complex logic problems by loading facts/rules and executing queries.

**Example:**
```json
{
  "name": "prolog_solve_problem",
  "arguments": {
    "problem_description": "Family relationships puzzle",
    "facts_and_rules": "parent(john, mary).\nparent(mary, alice).\ngrandparent(X, Z) :- parent(X, Y), parent(Y, Z).",
    "queries": ["grandparent(john, alice)."]
  }
}
```

### `prolog_explain_solution`
Get step-by-step explanations of Prolog solutions.

**Example:**
```json
{
  "name": "prolog_explain_solution", 
  "arguments": {
    "query": "ancestor(john, alice).",
    "facts": "parent(john, mary).\nparent(mary, alice).\nancestor(X, Y) :- parent(X, Y).\nancestor(X, Y) :- parent(X, Z), ancestor(Z, Y)."
  }
}
```

## Examples

The server includes comprehensive examples in the `examples/` directory:

- **`basic_examples.pl`**: Fundamental Prolog concepts (facts, rules, queries)
- **`family_tree.pl`**: Complex family relationship modeling  
- **`logic_puzzles.pl`**: Advanced puzzles (Einstein's riddle, map coloring, Sudoku, N-Queens)

### Running Examples

```bash
# Load and test basic examples
make example-basic

# Load family tree example
make example-family
```

## API Reference

### STDIO Mode
The server communicates via JSON-RPC 2.0 over standard input/output:

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "prolog_query",
    "arguments": {
      "query": "mammal(cat)."
    }
  }
}
```

### HTTP Mode  
REST API available at `http://localhost:8080/mcp`:

```bash
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"prolog_query","arguments":{"query":"true."}}}'
```

Health check endpoint: `GET http://localhost:8080/health`

## Development

### Project Structure
```
logic-mcp/
├── cmd/server/          # Main server application
├── internal/
│   ├── mcp/            # MCP protocol implementation  
│   ├── prolog/         # SWI-Prolog wrapper
│   └── tools/          # Logic tools implementation
├── examples/           # Prolog examples and tutorials
├── test/              # Integration tests
├── .vscode/           # VSCode configuration
├── Dockerfile         # Container definition
└── Makefile          # Build automation
```

### Building with Go

```bash
# Build server
make build

# Run all tests  
make test-all

# Clean artifacts
make clean
```

### Requirements

- Go 1.23+
- SWI-Prolog (for local development)
- Docker (for containerized deployment)

### Running Tests

```bash
# Unit tests
make test

# Integration tests (requires SWI-Prolog)
make test-integration

# All tests
make test-all

# Test specific components
make test-prolog    # Test Prolog engine only
make test-mcp       # Test MCP server only  
make test-tools     # Test logic tools only
```

## Error Handling

The server provides verbose error messages for common issues:

- **Syntax Errors**: Detailed validation with suggestions
- **Runtime Errors**: Clear execution error descriptions  
- **Timeout Handling**: Configurable query timeouts
- **Resource Management**: Automatic cleanup of Prolog processes

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-logic`)
3. Commit your changes (`git commit -am 'Add amazing logic feature'`)
4. Push to the branch (`git push origin feature/amazing-logic`)
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Prolog Documentation](https://www.swi-prolog.org/pldoc/doc_for?object=manual) 
- 🤝 [MCP Protocol Specification](https://modelcontextprotocol.io/)
- 🐛 [Issue Tracker](https://github.com/tomasz-sikora/logic-mcp/issues)

---

**Made with ❤️ for the AI and Logic Programming community**

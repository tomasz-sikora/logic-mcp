package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tomasz-sikora/logic-mcp/internal/prolog"
	"github.com/tomasz-sikora/logic-mcp/internal/tools"
)

func main() {
	var (
		mode = flag.String("mode", "stdio", "Server mode: stdio or http")
		port = flag.String("port", "8080", "HTTP server port (when mode=http)")
	)
	flag.Parse()

	// Initialize Prolog engine
	prologEngine, err := prolog.NewEngine()
	if err != nil {
		log.Fatalf("Failed to initialize Prolog engine: %v", err)
	}
	defer prologEngine.Close()

	// Create MCP server with official SDK
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "logic-mcp",
		Version: "v1.0.0",
	}, nil)

	// Initialize logic tools
	logicTools := tools.NewLogicTools(prologEngine)

	// Add all tools to the server
	if err := logicTools.RegisterTools(server); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}

	// Start server based on mode
	switch *mode {
	case "stdio":
		log.Println("Starting MCP server in STDIO mode...")
		err := server.Run(context.Background(), &mcp.StdioTransport{})
		if err != nil {
			log.Fatalf("STDIO server error: %v", err)
		}
		log.Println("server.Run() completed without error")
	case "http":
		log.Printf("Starting MCP server in HTTP mode on port %s...", *port)
		// For HTTP mode, we'll need to implement a custom HTTP transport
		// The official SDK primarily supports STDIO mode by default
		log.Fatalf("HTTP mode not yet implemented with official SDK - use stdio mode")
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'stdio' or 'http'\n", *mode)
		os.Exit(1)
	}

	log.Println("Server stopped")
}

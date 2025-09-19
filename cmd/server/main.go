package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
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

	// Create function to build per-session servers with isolated engines
	createSessionServer := func() (*mcp.Server, error) {
		// Create isolated Prolog engine for this session
		prologEngine, err := prolog.NewEngine()
		if err != nil {
			return nil, fmt.Errorf("failed to initialize session Prolog engine: %v", err)
		}

		// Create MCP server for this session
		server := mcp.NewServer(&mcp.Implementation{
			Name:    "logic-mcp",
			Version: "v1.0.0",
		}, nil)

		// Initialize logic tools with session engine
		logicTools := tools.NewLogicTools(prologEngine)

		// Add all tools to the session server
		if err := logicTools.RegisterTools(server); err != nil {
			return nil, fmt.Errorf("failed to register tools: %v", err)
		}

		return server, nil
	}

	// Start server based on mode
	switch *mode {
	case "stdio":
		log.Println("Starting MCP server in STDIO mode...")
		// Create dedicated server for STDIO mode (single session)
		server, err := createSessionServer()
		if err != nil {
			log.Fatalf("Failed to create STDIO server: %v", err)
		}
		err = server.Run(context.Background(), &mcp.StdioTransport{})
		if err != nil {
			log.Fatalf("STDIO server error: %v", err)
		}
		log.Println("server.Run() completed without error")
	case "http":
		log.Printf("Starting MCP server in HTTP mode on port %s...", *port)

		// Create StreamableHTTPHandler with per-session server creation
		handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
			// Create isolated server for each session
			server, err := createSessionServer()
			if err != nil {
				log.Printf("Failed to create session server: %v", err)
				return nil // This will result in a 400 Bad Request
			}
			log.Printf("Created new session server for request from %s", req.RemoteAddr)
			return server
		}, &mcp.StreamableHTTPOptions{
			JSONResponse: true, // Use JSON responses for better debugging
			Stateless:    true, // Enable stateless mode for easier HTTP testing
		})

		addr := fmt.Sprintf(":%s", *port)
		log.Printf("MCP HTTP server listening on %s", addr)
		if err := http.ListenAndServe(addr, handler); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'stdio' or 'http'\n", *mode)
		os.Exit(1)
	}

	log.Println("Server stopped")
}

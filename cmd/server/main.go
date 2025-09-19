package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tomasz-sikora/logic-mcp/internal/mcp"
	"github.com/tomasz-sikora/logic-mcp/internal/prolog"
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

	// Initialize MCP server
	mcpServer, err := mcp.NewServer(prologEngine)
	if err != nil {
		log.Fatalf("Failed to initialize MCP server: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	// Start server based on mode
	switch *mode {
	case "stdio":
		log.Println("Starting MCP server in STDIO mode...")
		if err := mcpServer.ServeSTDIO(ctx); err != nil {
			log.Fatalf("STDIO server error: %v", err)
		}
	case "http":
		log.Printf("Starting MCP server in HTTP mode on port %s...", *port)
		if err := mcpServer.ServeHTTP(ctx, *port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid mode: %s. Use 'stdio' or 'http'\n", *mode)
		os.Exit(1)
	}

	log.Println("Server stopped")
}

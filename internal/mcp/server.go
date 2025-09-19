package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/tomasz-sikora/logic-mcp/internal/prolog"
	"github.com/tomasz-sikora/logic-mcp/internal/tools"
)

// Server represents the MCP server
type Server struct {
	prologEngine *prolog.Engine
	tools        *tools.LogicTools
}

// MCPRequest represents a generic MCP request
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse represents a generic MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// NewServer creates a new MCP server instance
func NewServer(prologEngine *prolog.Engine) (*Server, error) {
	if prologEngine == nil {
		return nil, fmt.Errorf("prolog engine cannot be nil")
	}

	logicTools := tools.NewLogicTools(prologEngine)

	return &Server{
		prologEngine: prologEngine,
		tools:        logicTools,
	}, nil
}

// ServeSTDIO starts the MCP server in STDIO mode
func (s *Server) ServeSTDIO(ctx context.Context) error {
	// Read from stdin and write to stdout
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			var request MCPRequest
			if err := decoder.Decode(&request); err != nil {
				// Send error response
				response := MCPResponse{
					JSONRPC: "2.0",
					ID:      request.ID,
					Error: &MCPError{
						Code:    -32700,
						Message: "Parse error",
						Data:    err.Error(),
					},
				}
				encoder.Encode(response)
				continue
			}

			response := s.handleRequest(ctx, &request)
			if err := encoder.Encode(response); err != nil {
				return fmt.Errorf("failed to encode response: %w", err)
			}
		}
	}
}

// ServeHTTP starts the MCP server in HTTP mode
func (s *Server) ServeHTTP(ctx context.Context, port string) error {
	router := mux.NewRouter()

	router.HandleFunc("/mcp", s.handleHTTPRequest).Methods("POST")
	router.HandleFunc("/health", s.handleHealth).Methods("GET")

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	return server.ListenAndServe()
}

// handleHTTPRequest handles HTTP requests
func (s *Server) handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response := MCPResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error: &MCPError{
				Code:    -32700,
				Message: "Parse error",
				Data:    err.Error(),
			},
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := s.handleRequest(r.Context(), &request)
	json.NewEncoder(w).Encode(response)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "logic-mcp",
	})
}

// handleRequest processes MCP requests
func (s *Server) handleRequest(ctx context.Context, request *MCPRequest) *MCPResponse {
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(ctx, request)
	case "resources/list":
		return s.handleResourcesList(request)
	case "resources/read":
		return s.handleResourcesRead(ctx, request)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: "Method not found",
				Data:    fmt.Sprintf("Unknown method: %s", request.Method),
			},
		}
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(request *MCPRequest) *MCPResponse {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": true,
			},
			"resources": map[string]interface{}{
				"subscribe":   true,
				"listChanged": true,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "logic-mcp",
			"version": "1.0.0",
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleToolsList handles the tools/list request
func (s *Server) handleToolsList(request *MCPRequest) *MCPResponse {
	tools := s.tools.GetToolDefinitions()

	result := map[string]interface{}{
		"tools": tools,
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleToolsCall handles the tools/call request
func (s *Server) handleToolsCall(ctx context.Context, request *MCPRequest) *MCPResponse {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Expected object with 'name' and 'arguments'",
			},
		}
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Missing or invalid 'name' field",
			},
		}
	}

	arguments, ok := params["arguments"].(map[string]interface{})
	if !ok {
		arguments = make(map[string]interface{})
	}

	result, err := s.tools.CallTool(ctx, toolName, arguments)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: "Tool execution error",
				Data:    err.Error(),
			},
		}
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleResourcesList handles the resources/list request
func (s *Server) handleResourcesList(request *MCPRequest) *MCPResponse {
	resources := []map[string]interface{}{
		{
			"uri":         "prolog://examples/basic",
			"name":        "Basic Prolog Examples",
			"description": "Basic Prolog predicates and examples",
			"mimeType":    "text/prolog",
		},
		{
			"uri":         "prolog://examples/logic-puzzles",
			"name":        "Logic Puzzles",
			"description": "Examples of logic puzzles solved with Prolog",
			"mimeType":    "text/prolog",
		},
		{
			"uri":         "prolog://examples/family-tree",
			"name":        "Family Tree Example",
			"description": "Family relationships modeling in Prolog",
			"mimeType":    "text/prolog",
		},
	}

	result := map[string]interface{}{
		"resources": resources,
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleResourcesRead handles the resources/read request
func (s *Server) handleResourcesRead(ctx context.Context, request *MCPRequest) *MCPResponse {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Expected object with 'uri'",
			},
		}
	}

	uri, ok := params["uri"].(string)
	if !ok {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid params",
				Data:    "Missing or invalid 'uri' field",
			},
		}
	}

	content, err := s.getResourceContent(uri)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: "Resource not found",
				Data:    err.Error(),
			},
		}
	}

	result := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"uri":      uri,
				"mimeType": "text/prolog",
				"text":     content,
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// getResourceContent returns the content for a given resource URI
func (s *Server) getResourceContent(uri string) (string, error) {
	switch {
	case strings.HasSuffix(uri, "basic"):
		return getBasicPrologExamples(), nil
	case strings.HasSuffix(uri, "logic-puzzles"):
		return getLogicPuzzleExamples(), nil
	case strings.HasSuffix(uri, "family-tree"):
		return getFamilyTreeExample(), nil
	default:
		return "", fmt.Errorf("unknown resource URI: %s", uri)
	}
}

// getBasicPrologExamples returns basic Prolog examples
func getBasicPrologExamples() string {
	return `% Basic Prolog Facts and Rules

% Facts about animals
animal(dog).
animal(cat).
animal(bird).
animal(fish).

% Facts about properties
has_fur(dog).
has_fur(cat).
has_feathers(bird).
has_scales(fish).

% Rules
mammal(X) :- animal(X), has_fur(X).
can_fly(X) :- animal(X), has_feathers(X).

% Example queries:
% ?- mammal(dog).     % Should return true
% ?- mammal(bird).    % Should return false
% ?- can_fly(bird).   % Should return true
% ?- mammal(X).       % Should return X = dog; X = cat`
}

// getLogicPuzzleExamples returns logic puzzle examples
func getLogicPuzzleExamples() string {
	return `% Logic Puzzle: Einstein's Riddle (simplified version)

% Houses numbered 1-3
house(1). house(2). house(3).

% Colors
color(red). color(blue). color(green).

% Pets
pet(dog). pet(cat). pet(bird).

% The puzzle solution predicate
solve_puzzle(Houses) :-
    % Houses = [house(Num1, Color1, Pet1), house(Num2, Color2, Pet2), house(Num3, Color3, Pet3)]
    Houses = [house(1, C1, P1), house(2, C2, P2), house(3, C3, P3)],
    
    % All colors must be different
    permutation([red, blue, green], [C1, C2, C3]),
    
    % All pets must be different  
    permutation([dog, cat, bird], [P1, P2, P3]),
    
    % Constraints
    member(house(1, red, _), Houses),     % House 1 is red
    member(house(_, blue, dog), Houses),  % Blue house has a dog
    member(house(3, _, cat), Houses).     % House 3 has a cat

% Helper predicate
next_to(X, Y, List) :-
    append(_, [X, Y | _], List);
    append(_, [Y, X | _], List).

% Example query:
% ?- solve_puzzle(Solution).`
}

// getFamilyTreeExample returns family tree example
func getFamilyTreeExample() string {
	return `% Family Tree Example

% Basic facts about people
person(john).
person(mary).
person(bob).
person(alice).
person(charlie).
person(diana).

% Parent relationships
parent(john, bob).
parent(mary, bob).
parent(bob, alice).
parent(bob, charlie).
parent(alice, diana).

% Gender facts
male(john).
male(bob).
male(charlie).
female(mary).
female(alice).
female(diana).

% Rules for family relationships
father(X, Y) :- parent(X, Y), male(X).
mother(X, Y) :- parent(X, Y), female(X).

child(X, Y) :- parent(Y, X).
son(X, Y) :- child(X, Y), male(X).
daughter(X, Y) :- child(X, Y), female(X).

grandparent(X, Z) :- parent(X, Y), parent(Y, Z).
grandfather(X, Z) :- grandparent(X, Z), male(X).
grandmother(X, Z) :- grandparent(X, Z), female(X).

sibling(X, Y) :- parent(Z, X), parent(Z, Y), X \= Y.
brother(X, Y) :- sibling(X, Y), male(X).
sister(X, Y) :- sibling(X, Y), female(X).

% Example queries:
% ?- father(john, bob).      % Is John the father of Bob?
% ?- mother(X, bob).         % Who is Bob's mother?
% ?- grandparent(john, X).   % Who are John's grandchildren?
% ?- sibling(alice, charlie). % Are Alice and Charlie siblings?`
}

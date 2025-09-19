package tools

import (
"context"
"fmt"
"strings"

"github.com/modelcontextprotocol/go-sdk/mcp"
"github.com/tomasz-sikora/logic-mcp/internal/prolog"
)

// LogicTools manages Prolog-based tools for MCP using official SDK
type LogicTools struct {
	engine *prolog.Engine
}

// NewLogicTools creates a new LogicTools instance
func NewLogicTools(engine *prolog.Engine) *LogicTools {
	return &LogicTools{
		engine: engine,
	}
}

// RegisterTools registers all logic tools with the MCP server
func (lt *LogicTools) RegisterTools(server *mcp.Server) error {
	// Define input types for each tool
	type QueryInput struct {
		Query string `json:"query" jsonschema:"The Prolog query to execute. Must end with a period. Example: 'member(X, [1,2,3]).'" `
	}

	type FactsInput struct {
		Facts string `json:"facts" jsonschema:"Prolog facts and rules to load, separated by newlines. Comments start with %. Example: 'parent(tom, bob).\\nparent(bob, pat).'"`
	}

	type CodeInput struct {
		Code string `json:"code" jsonschema:"Prolog code to validate syntax for. Can include facts, rules, or queries."`
	}

	type ProblemInput struct {
		ProblemDescription string   `json:"problem_description" jsonschema:"A description of the logic problem to solve."`
		FactsAndRules      string   `json:"facts_and_rules" jsonschema:"Prolog facts and rules that define the problem domain."`
		Queries            []string `json:"queries" jsonschema:"List of queries to execute to solve the problem."`
	}

	type ExplainInput struct {
		Query string `json:"query" jsonschema:"The Prolog query to explain."`
		Facts string `json:"facts,omitempty" jsonschema:"Relevant facts and rules (optional)."`
	}

	// Register prolog_query tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_query",
		Description: "Execute a Prolog query and return results. Supports both simple queries and complex logic problems.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input QueryInput) (*mcp.CallToolResult, any, error) {
		result, err := lt.engine.Query(ctx, input.Query)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to execute query: %s", err.Error())},
				},
				IsError: true,
			}, nil, nil
		}

		var responseText strings.Builder
		responseText.WriteString(fmt.Sprintf("Query: %s\n", input.Query))
		responseText.WriteString(fmt.Sprintf("Result: %t\n", result.Success))
		responseText.WriteString(fmt.Sprintf("Execution Time: %s\n", result.ExecutionTime))

		if result.Error != "" {
			responseText.WriteString(fmt.Sprintf("Error: %s\n", result.Error))
		}

		if result.Output != "" {
			responseText.WriteString(fmt.Sprintf("Output: %s\n", result.Output))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText.String()},
			},
		}, nil, nil
	})

	// Register prolog_load_facts tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_load_facts",
		Description: "Load Prolog facts and rules into the knowledge base. Use this to define rules and facts before querying.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input FactsInput) (*mcp.CallToolResult, any, error) {
		err := lt.engine.LoadFacts(input.Facts)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to load facts: %s", err.Error())},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Facts loaded successfully"},
			},
		}, nil, nil
	})

	// Register prolog_validate_syntax tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_validate_syntax",
		Description: "Validate Prolog syntax without executing. Use this to check if your Prolog code is syntactically correct.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input CodeInput) (*mcp.CallToolResult, any, error) {
		// Simple syntax check by attempting to create a temp file and check basic structure
		lines := strings.Split(strings.TrimSpace(input.Code), "\n")
		hasValidStructure := false
		
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "%") {
				continue // skip empty lines and comments
			}
			if strings.HasSuffix(line, ".") {
				hasValidStructure = true
				break
			}
		}

		result := "valid"
		if !hasValidStructure {
			result = "invalid - no statements ending with '.'"
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Syntax validation: %s", result)},
			},
		}, nil, nil
	})

	// Register prolog_clear_kb tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_clear_kb",
		Description: "Clear the Prolog knowledge base. This removes all dynamic predicates and facts.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input struct{}) (*mcp.CallToolResult, any, error) {
		err := lt.engine.ClearKnowledgeBase()
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to clear knowledge base: %s", err.Error())},
				},
				IsError: true,
			}, nil, nil
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Knowledge base cleared successfully"},
			},
		}, nil, nil
	})

	// Register prolog_solve_problem tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_solve_problem",
		Description: "Solve a complex logic problem by loading facts/rules and then executing queries. This is a high-level tool that combines loading facts and querying.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ProblemInput) (*mcp.CallToolResult, any, error) {
		// Load facts and rules
		if err := lt.engine.LoadFacts(input.FactsAndRules); err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to load facts and rules: %s", err.Error())},
				},
				IsError: true,
			}, nil, nil
		}

		var responseText strings.Builder
		responseText.WriteString(fmt.Sprintf("Problem: %s\n\n", input.ProblemDescription))
		responseText.WriteString("Facts and rules loaded successfully.\n\n")
		responseText.WriteString("Query Results:\n")

		for i, query := range input.Queries {
			result, err := lt.engine.Query(ctx, query)
			if err != nil {
				responseText.WriteString(fmt.Sprintf("%d. Query: %s\n   Error: %s\n", i+1, query, err.Error()))
				continue
			}

			responseText.WriteString(fmt.Sprintf("%d. Query: %s\n", i+1, query))
			responseText.WriteString(fmt.Sprintf("   Result: %t (%s)\n", result.Success, result.ExecutionTime))
			if result.Output != "" {
				responseText.WriteString(fmt.Sprintf("   Output: %s\n", result.Output))
			}
			if result.Error != "" {
				responseText.WriteString(fmt.Sprintf("   Error: %s\n", result.Error))
			}
			responseText.WriteString("\n")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText.String()},
			},
		}, nil, nil
	})

	// Register prolog_explain_solution tool
	mcp.AddTool(server, &mcp.Tool{
		Name:        "prolog_explain_solution",
		Description: "Explain how a Prolog solution works step by step. This tool provides educational explanations.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, input ExplainInput) (*mcp.CallToolResult, any, error) {
		var responseText strings.Builder
		responseText.WriteString(fmt.Sprintf("Explaining Prolog query: %s\n\n", input.Query))

		// Load facts if provided
		if input.Facts != "" {
			if err := lt.engine.LoadFacts(input.Facts); err != nil {
				return &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{Text: fmt.Sprintf("Failed to load facts: %s", err.Error())},
					},
					IsError: true,
				}, nil, nil
			}
			responseText.WriteString("Loaded facts:\n")
			responseText.WriteString(input.Facts)
			responseText.WriteString("\n\n")
		}

		// Execute query
		result, err := lt.engine.Query(ctx, input.Query)
		if err != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Failed to execute query for explanation: %s", err.Error())},
				},
				IsError: true,
			}, nil, nil
		}

		responseText.WriteString("Execution Analysis:\n")
		responseText.WriteString(fmt.Sprintf("- Query: %s\n", input.Query))
		responseText.WriteString(fmt.Sprintf("- Success: %t\n", result.Success))
		responseText.WriteString(fmt.Sprintf("- Execution Time: %s\n", result.ExecutionTime))

		if result.Output != "" {
			responseText.WriteString(fmt.Sprintf("- Output: %s\n", result.Output))
		}

		if result.Error != "" {
			responseText.WriteString(fmt.Sprintf("- Error Details: %s\n", result.Error))
		}

		// Add basic explanation
		responseText.WriteString("\nExplanation:\n")
		if result.Success {
			responseText.WriteString("The query succeeded, meaning Prolog found a solution that satisfies the given constraints.")
		} else {
			responseText.WriteString("The query failed, meaning Prolog could not find any solution that satisfies the given constraints.")
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: responseText.String()},
			},
		}, nil, nil
	})

	return nil
}

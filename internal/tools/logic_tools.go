package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/tomasz-sikora/logic-mcp/internal/prolog"
)

// ToolDefinition represents an MCP tool definition
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	Content []map[string]interface{} `json:"content"`
	IsError bool                     `json:"isError,omitempty"`
}

// LogicTools manages Prolog-based tools for MCP
type LogicTools struct {
	engine *prolog.Engine
}

// NewLogicTools creates a new LogicTools instance
func NewLogicTools(engine *prolog.Engine) *LogicTools {
	return &LogicTools{
		engine: engine,
	}
}

// GetToolDefinitions returns all available tool definitions
func (lt *LogicTools) GetToolDefinitions() []ToolDefinition {
	return []ToolDefinition{
		{
			Name:        "prolog_query",
			Description: "Execute a Prolog query and return results. Supports both simple queries and complex logic problems.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The Prolog query to execute. Must end with a period. Example: 'member(X, [1,2,3]).'",
					},
				},
				"required": []string{"query"},
			},
		},
		{
			Name:        "prolog_load_facts",
			Description: "Load Prolog facts and rules into the knowledge base. Use this to define rules and facts before querying.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"facts": map[string]interface{}{
						"type":        "string",
						"description": "Prolog facts and rules to load, separated by newlines. Comments start with %. Example: 'parent(tom, bob).\\nparent(bob, pat).'",
					},
				},
				"required": []string{"facts"},
			},
		},
		{
			Name:        "prolog_validate_syntax",
			Description: "Validate Prolog syntax without executing. Use this to check if your Prolog code is syntactically correct.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"code": map[string]interface{}{
						"type":        "string",
						"description": "Prolog code to validate syntax for. Can include facts, rules, or queries.",
					},
				},
				"required": []string{"code"},
			},
		},
		{
			Name:        "prolog_clear_kb",
			Description: "Clear the Prolog knowledge base. This removes all dynamic predicates and facts.",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "prolog_solve_problem",
			Description: "Solve a complex logic problem by loading facts/rules and then executing queries. This is a high-level tool that combines loading facts and querying.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"problem_description": map[string]interface{}{
						"type":        "string",
						"description": "A description of the logic problem to solve.",
					},
					"facts_and_rules": map[string]interface{}{
						"type":        "string",
						"description": "Prolog facts and rules that define the problem domain.",
					},
					"queries": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
						"description": "List of queries to execute to solve the problem.",
					},
				},
				"required": []string{"problem_description", "facts_and_rules", "queries"},
			},
		},
		{
			Name:        "prolog_explain_solution",
			Description: "Explain how a Prolog solution works step by step. This tool provides educational explanations.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{
						"type":        "string",
						"description": "The Prolog query to explain.",
					},
					"facts": map[string]interface{}{
						"type":        "string",
						"description": "Relevant facts and rules (optional).",
					},
				},
				"required": []string{"query"},
			},
		},
	}
}

// CallTool executes a tool with the given name and arguments
func (lt *LogicTools) CallTool(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error) {
	switch name {
	case "prolog_query":
		return lt.handleQuery(ctx, args)
	case "prolog_load_facts":
		return lt.handleLoadFacts(args)
	case "prolog_validate_syntax":
		return lt.handleValidateSyntax(args)
	case "prolog_clear_kb":
		return lt.handleClearKB()
	case "prolog_solve_problem":
		return lt.handleSolveProblem(ctx, args)
	case "prolog_explain_solution":
		return lt.handleExplainSolution(ctx, args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// handleQuery handles prolog_query tool calls
func (lt *LogicTools) handleQuery(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	query, ok := args["query"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'query' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	result, err := lt.engine.Query(ctx, query)
	if err != nil {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Failed to execute query: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	var responseText strings.Builder
	responseText.WriteString(fmt.Sprintf("Query: %s\n", query))
	responseText.WriteString(fmt.Sprintf("Result: %t\n", result.Success))
	responseText.WriteString(fmt.Sprintf("Execution Time: %s\n", result.ExecutionTime))

	if result.Error != "" {
		responseText.WriteString(fmt.Sprintf("Error: %s\n", result.Error))
	}

	if result.Output != "" {
		responseText.WriteString(fmt.Sprintf("Output: %s\n", result.Output))
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": responseText.String(),
			},
		},
	}, nil
}

// handleLoadFacts handles prolog_load_facts tool calls
func (lt *LogicTools) handleLoadFacts(args map[string]interface{}) (*ToolResult, error) {
	facts, ok := args["facts"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'facts' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	err := lt.engine.LoadFacts(facts)
	if err != nil {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Failed to load facts: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": "Facts loaded successfully!",
			},
		},
	}, nil
}

// handleValidateSyntax handles prolog_validate_syntax tool calls
func (lt *LogicTools) handleValidateSyntax(args map[string]interface{}) (*ToolResult, error) {
	code, ok := args["code"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'code' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	err := lt.engine.ValidateQuery(code)
	if err != nil {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Syntax validation failed: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": "Syntax is valid!",
			},
		},
	}, nil
}

// handleClearKB handles prolog_clear_kb tool calls
func (lt *LogicTools) handleClearKB() (*ToolResult, error) {
	err := lt.engine.ClearKnowledgeBase()
	if err != nil {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Failed to clear knowledge base: %s", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": "Knowledge base cleared successfully!",
			},
		},
	}, nil
}

// handleSolveProblem handles prolog_solve_problem tool calls
func (lt *LogicTools) handleSolveProblem(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	description, ok := args["problem_description"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'problem_description' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	factsAndRules, ok := args["facts_and_rules"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'facts_and_rules' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	queries, ok := args["queries"].([]interface{})
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'queries' parameter must be an array",
				},
			},
			IsError: true,
		}, nil
	}

	var responseText strings.Builder
	responseText.WriteString(fmt.Sprintf("🧩 Solving Problem: %s\n\n", description))

	// Clear knowledge base first
	if err := lt.engine.ClearKnowledgeBase(); err != nil {
		responseText.WriteString(fmt.Sprintf("⚠️ Warning: Failed to clear knowledge base: %s\n", err.Error()))
	}

	// Load facts and rules
	responseText.WriteString("📚 Loading facts and rules...\n")
	if err := lt.engine.LoadFacts(factsAndRules); err != nil {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("❌ Failed to load facts and rules: %s\n", err.Error()),
				},
			},
			IsError: true,
		}, nil
	}
	responseText.WriteString("✅ Facts and rules loaded successfully!\n\n")

	// Execute queries
	responseText.WriteString("🔍 Executing queries:\n")
	for i, queryInterface := range queries {
		query, ok := queryInterface.(string)
		if !ok {
			responseText.WriteString(fmt.Sprintf("❌ Query %d: Invalid query type (must be string)\n", i+1))
			continue
		}

		result, err := lt.engine.Query(ctx, query)
		if err != nil {
			responseText.WriteString(fmt.Sprintf("❌ Query %d (%s): Failed - %s\n", i+1, query, err.Error()))
			continue
		}

		status := "❌ Failed"
		if result.Success {
			status = "✅ Success"
		}
		responseText.WriteString(fmt.Sprintf("%s Query %d: %s\n", status, i+1, query))

		if result.Error != "" {
			responseText.WriteString(fmt.Sprintf("   Error: %s\n", result.Error))
		}
		if result.Output != "" {
			responseText.WriteString(fmt.Sprintf("   Output: %s\n", result.Output))
		}
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": responseText.String(),
			},
		},
	}, nil
}

// handleExplainSolution handles prolog_explain_solution tool calls
func (lt *LogicTools) handleExplainSolution(ctx context.Context, args map[string]interface{}) (*ToolResult, error) {
	query, ok := args["query"].(string)
	if !ok {
		return &ToolResult{
			Content: []map[string]interface{}{
				{
					"type": "text",
					"text": "Error: 'query' parameter must be a string",
				},
			},
			IsError: true,
		}, nil
	}

	facts, _ := args["facts"].(string)

	var responseText strings.Builder
	responseText.WriteString(fmt.Sprintf("🔬 Explaining Prolog Solution: %s\n\n", query))

	// If facts are provided, load them first
	if facts != "" {
		if err := lt.engine.LoadFacts(facts); err != nil {
			responseText.WriteString(fmt.Sprintf("⚠️ Warning: Failed to load provided facts: %s\n", err.Error()))
		} else {
			responseText.WriteString("📚 Loaded provided facts and rules\n\n")
		}
	}

	// Execute the query
	result, err := lt.engine.Query(ctx, query)
	if err != nil {
		responseText.WriteString(fmt.Sprintf("❌ Failed to execute query: %s\n", err.Error()))
	} else {
		responseText.WriteString("📝 Query Execution:\n")
		responseText.WriteString(fmt.Sprintf("   Query: %s\n", query))
		responseText.WriteString(fmt.Sprintf("   Result: %t\n", result.Success))
		responseText.WriteString(fmt.Sprintf("   Execution Time: %s\n", result.ExecutionTime))

		if result.Output != "" {
			responseText.WriteString(fmt.Sprintf("   Output: %s\n", result.Output))
		}

		if result.Error != "" {
			responseText.WriteString(fmt.Sprintf("   Error: %s\n", result.Error))
		}

		responseText.WriteString("\n💡 Explanation:\n")
		if result.Success {
			responseText.WriteString("The query succeeded, meaning Prolog was able to prove the goal using the loaded facts and rules through logical inference.\n")
		} else {
			responseText.WriteString("The query failed, meaning Prolog could not prove the goal with the available facts and rules. This could be because:\n")
			responseText.WriteString("- The goal is not derivable from the current knowledge base\n")
			responseText.WriteString("- Required facts or rules are missing\n")
			responseText.WriteString("- There's a logical inconsistency\n")
		}
	}

	return &ToolResult{
		Content: []map[string]interface{}{
			{
				"type": "text",
				"text": responseText.String(),
			},
		},
	}, nil
}
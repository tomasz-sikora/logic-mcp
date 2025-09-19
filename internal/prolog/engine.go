package prolog

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// QueryResult represents the result of a Prolog query
type QueryResult struct {
	Success       bool             `json:"success"`
	Solutions     []map[string]any `json:"solutions,omitempty"`
	Output        string           `json:"output,omitempty"`
	Error         string           `json:"error,omitempty"`
	ExecutionTime time.Duration    `json:"execution_time"`
}

// Engine manages SWI-Prolog execution
type Engine struct {
	tempFiles []string
	mutex     sync.Mutex
	closed    bool
	facts     []string // Store loaded facts
}

// NewEngine creates a new Prolog engine instance
func NewEngine() (*Engine, error) {
	// Check if SWI-Prolog is available
	if _, err := exec.LookPath("swipl"); err != nil {
		return nil, fmt.Errorf("SWI-Prolog not found: %w", err)
	}

	engine := &Engine{
		facts: make([]string, 0),
	}

	return engine, nil
}

// Query executes a Prolog query and returns the result
func (e *Engine) Query(ctx context.Context, query string) (*QueryResult, error) {
	startTime := time.Now()

	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.closed {
		return nil, fmt.Errorf("engine is closed")
	}

	// Validate query
	if strings.TrimSpace(query) == "" {
		return &QueryResult{
			Success:       false,
			Error:         "Empty query provided",
			ExecutionTime: time.Since(startTime),
		}, nil
	}

	// Ensure query ends with a period
	if !strings.HasSuffix(strings.TrimSpace(query), ".") {
		query = strings.TrimSpace(query) + "."
	}

	// Execute query using batch mode
	result, err := e.executeQueryBatch(ctx, query)
	if err != nil {
		return &QueryResult{
			Success:       false,
			Error:         err.Error(),
			ExecutionTime: time.Since(startTime),
		}, nil
	}

	result.ExecutionTime = time.Since(startTime)
	return result, nil
}

// executeQueryBatch executes a query in batch mode
func (e *Engine) executeQueryBatch(ctx context.Context, query string) (*QueryResult, error) {
	// Create temporary file for the query
	tempFile, err := e.createTempFile("query.pl")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile)

	// Write facts and query to file
	content := strings.Join(e.facts, "\n")
	if content != "" {
		content += "\n"
	}

	// Create a goal that will test the query and print result
	testGoal := fmt.Sprintf(`
main :-
    (   (%s) ->
        write('SUCCESS: true')
    ;   write('SUCCESS: false')
    ),
    nl,
    halt.
`, strings.TrimSuffix(query, "."))

	content += testGoal

	if err := ioutil.WriteFile(tempFile, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write query file: %w", err)
	}

	// Execute SWI-Prolog with the file
	cmd := exec.CommandContext(ctx, "swipl", "-q", "-g", "main", "-t", "halt", tempFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return &QueryResult{
			Success: false,
			Output:  string(output),
			Error:   fmt.Sprintf("execution failed: %v", err),
		}, nil
	}

	// Parse output
	outputStr := string(output)
	success := strings.Contains(outputStr, "SUCCESS: true")

	return &QueryResult{
		Success: success,
		Output:  outputStr,
	}, nil
}

// LoadFacts loads Prolog facts and rules into the knowledge base
func (e *Engine) LoadFacts(facts string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.closed {
		return fmt.Errorf("engine is closed")
	}

	// Parse facts line by line
	lines := strings.Split(facts, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "%") {
			// Ensure line ends with period
			if !strings.HasSuffix(line, ".") {
				line += "."
			}
			e.facts = append(e.facts, line)
		}
	}

	return nil
}

// ValidateQuery validates Prolog syntax without executing
func (e *Engine) ValidateQuery(query string) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return fmt.Errorf("empty query")
	}

	// Basic syntax validation
	if !strings.HasSuffix(query, ".") {
		return fmt.Errorf("query must end with a period")
	}

	// Check balanced parentheses
	if err := e.validateBalanced(query, '(', ')'); err != nil {
		return fmt.Errorf("unbalanced parentheses: %w", err)
	}

	// Check balanced brackets
	if err := e.validateBalanced(query, '[', ']'); err != nil {
		return fmt.Errorf("unbalanced brackets: %w", err)
	}

	return nil
}

// validateBalanced checks if brackets/parentheses are balanced
func (e *Engine) validateBalanced(s string, open, close rune) error {
	count := 0
	for _, ch := range s {
		if ch == open {
			count++
		} else if ch == close {
			count--
			if count < 0 {
				return fmt.Errorf("unmatched closing %c", close)
			}
		}
	}
	if count != 0 {
		return fmt.Errorf("unmatched opening %c", open)
	}
	return nil
}

// ClearKnowledgeBase clears all loaded facts and rules
func (e *Engine) ClearKnowledgeBase() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.closed {
		return fmt.Errorf("engine is closed")
	}

	e.facts = make([]string, 0)
	return nil
}

// Close cleans up the engine resources
func (e *Engine) Close() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.closed {
		return nil
	}

	// Clean up temporary files
	for _, file := range e.tempFiles {
		os.Remove(file)
	}
	e.tempFiles = nil

	e.closed = true
	return nil
}

// createTempFile creates a temporary file with the given name
func (e *Engine) createTempFile(name string) (string, error) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("logic_mcp_%d_%s", time.Now().UnixNano(), name))

	e.tempFiles = append(e.tempFiles, tempFile)
	return tempFile, nil
}

// GetLoadedFacts returns currently loaded facts (for debugging)
func (e *Engine) GetLoadedFacts() []string {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	result := make([]string, len(e.facts))
	copy(result, e.facts)
	return result
}

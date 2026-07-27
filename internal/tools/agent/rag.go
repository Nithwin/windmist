package agent

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Nithwin/WindMist/internal/rag"
	"github.com/Nithwin/WindMist/internal/tools"
)

// SemanticSearchTool allows the agent to search the codebase semantically.
type SemanticSearchTool struct {
	searcher *rag.Searcher
}

// NewSemanticSearchTool creates a new SemanticSearchTool.
func NewSemanticSearchTool(searcher *rag.Searcher) *SemanticSearchTool {
	return &SemanticSearchTool{
		searcher: searcher,
	}
}

// Definition returns the tool's schema.
func (t *SemanticSearchTool) Definition() tools.Definition {
	return tools.Definition{
		Name:        "semantic_search",
		Description: "Searches the codebase semantically using vector embeddings. Useful for finding concepts, functionality, or related code when you don't know the exact file path or keywords. Only use this if you need to explore the codebase conceptually.",
		Category:    tools.CategorySearch,
		Permission:  tools.PermReadOnly,
		Parameters: []tools.Parameter{
			{
				Name:        "query",
				Type:        "string",
				Description: "The search query (e.g. 'user authentication logic', 'database connection setup')",
				Required:    true,
			},
			{
				Name:        "top_k",
				Type:        "integer",
				Description: "Number of results to return (default 5, max 15)",
				Required:    false,
			},
		},
	}
}

// Run executes the semantic search.
func (t *SemanticSearchTool) Run(ctx context.Context, call tools.Call) tools.Result {
	start := time.Now()

	query, ok := call.Args["query"].(string)
	if !ok || query == "" {
		return tools.Result{Error: fmt.Errorf("query is required")}
	}

	topK := 5
	if k, ok := call.Args["top_k"].(float64); ok {
		topK = int(k)
	}
	if topK > 15 {
		topK = 15
	}
	if topK < 1 {
		topK = 5
	}

	results, err := t.searcher.Search(query, topK)
	if err != nil {
		return tools.Result{Error: fmt.Errorf("search failed: %w", err)}
	}

	if len(results) == 0 {
		return tools.Result{
			Output:   "No relevant code found.",
			Duration: time.Since(start),
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d relevant code chunks for '%s':\n\n", len(results), query))

	var filesRead []string
	seenFiles := make(map[string]bool)

	for i, res := range results {
		sb.WriteString(fmt.Sprintf("--- Result %d (Similarity: %.2f) ---\n", i+1, res.Similarity))
		sb.WriteString(fmt.Sprintf("File: %s (Lines %d-%d)\n", res.Chunk.FilePath, res.Chunk.StartLine, res.Chunk.EndLine))
		sb.WriteString("```\n")
		sb.WriteString(res.Chunk.Content)
		sb.WriteString("\n```\n\n")

		if !seenFiles[res.Chunk.FilePath] {
			seenFiles[res.Chunk.FilePath] = true
			filesRead = append(filesRead, res.Chunk.FilePath)
		}
	}

	return tools.Result{
		Output:    sb.String(),
		Duration:  time.Since(start),
		FilesRead: filesRead,
	}
}

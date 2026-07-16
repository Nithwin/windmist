package editing

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestSearchText(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	file1 := filepath.Join(tempDir, "provider.go")
	content1 := `package provider

type Provider interface {
	Generate() error
}

func newProvider() Provider {
	return nil
}`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatalf("failed to create provider.go: %v", err)
	}

	file2 := filepath.Join(tempDir, "chat.go")
	content2 := `package chat

func run(provider Provider) {
	// do something
}`
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatalf("failed to create chat.go: %v", err)
	}

	opts := SearchOptions{
		Root:  tempDir,
		Query: "Provider",
	}

	results, err := Search(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Fatalf("expected matches across 2 files, got %d: %+v", len(results), results)
	}

	// Verify grouping by file
	fileMatches := make(map[string][]Match)
	for _, r := range results {
		fileMatches[filepath.Base(r.Path)] = r.Matches
	}

	if len(fileMatches["provider.go"]) != 4 {
		t.Errorf("expected 4 matches in provider.go, got %d", len(fileMatches["provider.go"]))
	}
	if len(fileMatches["chat.go"]) != 2 {
		t.Errorf("expected 2 matches in chat.go, got %d", len(fileMatches["chat.go"]))
	}
}

func TestSearchCaseSensitivity(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	path := filepath.Join(tempDir, "test.go")
	content := `Provider provider PROVIDER`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test.go: %v", err)
	}

	// Case insensitive (default)
	resInsensitive, err := Search(ctx, SearchOptions{Root: tempDir, Query: "provider"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resInsensitive) != 1 || len(resInsensitive[0].Matches) != 3 {
		t.Errorf("expected 3 matches case-insensitive, got %+v", resInsensitive)
	}

	// Case sensitive
	resSensitive, err := Search(ctx, SearchOptions{Root: tempDir, Query: "provider", CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resSensitive) != 1 || len(resSensitive[0].Matches) != 1 {
		t.Errorf("expected 1 match case-sensitive, got %+v", resSensitive)
	}
}

func TestSearchMultilineQuery(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	path := filepath.Join(tempDir, "test.go")
	content := "func f1() error {\n\treturn nil\n}\n\nfunc f2() error {\n\treturn nil\n}\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test.go: %v", err)
	}

	results, err := Search(ctx, SearchOptions{
		Root:  tempDir,
		Query: "error {\n\treturn nil\n}",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || len(results[0].Matches) != 2 {
		t.Fatalf("expected 2 multi-line matches, got %+v", results)
	}
	if results[0].Matches[0].Line != 1 || results[0].Matches[1].Line != 5 {
		t.Errorf("expected matches starting on lines 1 and 5, got %+v", results[0].Matches)
	}
}

func TestSearchWholeWord(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	path := filepath.Join(tempDir, "test.txt")
	content := `provider provider.Generate() providers provider_func`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test.txt: %v", err)
	}

	results, err := Search(ctx, SearchOptions{Root: tempDir, Query: "provider", WholeWord: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || len(results[0].Matches) != 2 {
		t.Fatalf("expected 2 whole-word matches ('provider' and 'provider' in 'provider.Generate()'), got %+v", results)
	}
}

func TestSearchRegex(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	path := filepath.Join(tempDir, "test.go")
	content := `provider.Generate()
provider.Close()
provider.Other`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test.go: %v", err)
	}

	results, err := Search(ctx, SearchOptions{
		Root:  tempDir,
		Query: `provider\..*\(\)`,
		Type:  SearchRegex,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || len(results[0].Matches) != 2 {
		t.Fatalf("expected 2 regex matches, got %+v", results)
	}
}

func TestSearchStream(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	if err := os.WriteFile(filepath.Join(tempDir, "stream1.go"), []byte("package main\n// provider"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "stream2.go"), []byte("package main\n// provider"), 0644); err != nil {
		t.Fatal(err)
	}

	var streamedResults []SearchResult
	err := SearchStream(ctx, SearchOptions{Root: tempDir, Query: "provider"}, func(res SearchResult) error {
		streamedResults = append(streamedResults, res)
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected stream error: %v", err)
	}

	if len(streamedResults) != 2 {
		t.Fatalf("expected 2 streamed results immediately, got %d", len(streamedResults))
	}
}

func TestSearchFileName(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	if err := os.WriteFile(filepath.Join(tempDir, "provider.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "other.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}

	results, err := Search(ctx, SearchOptions{
		Root:  tempDir,
		Query: "provider",
		Type:  SearchFileName,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 1 || filepath.Base(results[0].Path) != "provider.go" {
		t.Fatalf("expected 1 filename match for provider.go, got %+v", results)
	}
}

func TestSearchBinarySkip(t *testing.T) {
	tempDir := t.TempDir()
	ctx := context.Background()

	binaryPath := filepath.Join(tempDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 'p', 'r', 'o', 'v', 'i', 'd', 'e', 'r'}
	if err := os.WriteFile(binaryPath, binaryContent, 0644); err != nil {
		t.Fatal(err)
	}

	results, err := Search(ctx, SearchOptions{Root: tempDir, Query: "provider"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("expected binary file to be skipped, got %+v", results)
	}
}

func TestSearchToolRun(t *testing.T) {
	tempDir := t.TempDir()
	tool := NewSearchTool()
	ctx := context.Background()

	if err := os.WriteFile(filepath.Join(tempDir, "app.go"), []byte("func main() { fmt.Println() }"), 0644); err != nil {
		t.Fatal(err)
	}

	call := tools.Call{
		Name: "search_text",
		Args: map[string]any{
			"query": "fmt.Println",
			"path":  tempDir,
		},
	}

	result := tool.Run(ctx, call)
	if result.Error != nil {
		t.Fatalf("unexpected error: %v", result.Error)
	}

	results, ok := result.Output.(Results)
	if !ok {
		t.Fatalf("expected Results output, got %T", result.Output)
	}

	if len(results) != 1 || len(results[0].Matches) != 1 {
		t.Fatalf("expected 1 match via tool run, got %+v", results)
	}
}

func TestSearchDefinition(t *testing.T) {
	tool := NewSearchTool()
	def := tool.Definition()

	if def.Name != "search_text" {
		t.Errorf("unexpected name: %s", def.Name)
	}
	if def.Description == "" {
		t.Error("empty description")
	}
	if len(def.Parameters) < 1 {
		t.Fatalf("unexpected parameters count: %d", len(def.Parameters))
	}
	if def.Parameters[0].Name != "query" || !def.Parameters[0].Required {
		t.Error("expected query parameter to be required")
	}
}

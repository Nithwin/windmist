package rag

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Indexer manages building the RAG index for a project.
type Indexer struct {
	store    *DocumentStore
	embedder *TFIDFEmbedder
}

// NewIndexer creates a new Indexer.
func NewIndexer(store *DocumentStore, embedder *TFIDFEmbedder) *Indexer {
	return &Indexer{
		store:    store,
		embedder: embedder,
	}
}

// IndexProject scans a directory, chunks all valid source files,
// trains the embedder, and stores the embedded chunks.
func (i *Indexer) IndexProject(rootDir string) (int, error) {
	// 1. Collect all valid files
	var files []string
	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			name := d.Name()
			// Skip hidden and common ignored directories
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "dist" || name == "build" {
				return filepath.SkipDir
			}
			return nil
		}

		// Only index code/text files. Very basic filter for now.
		ext := filepath.Ext(path)
		switch ext {
		case ".go", ".md", ".txt", ".js", ".ts", ".py", ".rs", ".c", ".cpp", ".h", ".json", ".yaml", ".yml":
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("scan project: %w", err)
	}

	// 2. Read and chunk all files
	var allChunks []Chunk
	var documents []string // Just the text content for building vocabulary

	cfg := DefaultChunkConfig()

	for _, path := range files {
		contentBytes, err := os.ReadFile(path)
		if err != nil {
			continue // Skip unreadable files
		}
		content := string(contentBytes)

		// Get relative path for cleaner storage
		relPath, _ := filepath.Rel(rootDir, path)
		if relPath == "" {
			relPath = path
		}

		chunks := ChunkFile(relPath, content, cfg)
		for _, c := range chunks {
			allChunks = append(allChunks, c)
			documents = append(documents, c.Content)
		}
	}

	if len(allChunks) == 0 {
		return 0, nil
	}

	// 3. Train the embedder on the corpus
	i.embedder.BuildVocabulary(documents)

	// 4. Clear old index
	if err := i.store.Clear(); err != nil {
		return 0, fmt.Errorf("clear store: %w", err)
	}

	// 5. Embed and store chunks
	indexedCount := 0
	for _, chunk := range allChunks {
		vec := i.embedder.Embed(chunk.Content)
		if vec == nil {
			continue
		}

		if err := i.store.InsertChunk(chunk, vec); err != nil {
			// Log error but continue
			fmt.Printf("Warning: failed to insert chunk for %s: %v\n", chunk.FilePath, err)
			continue
		}
		indexedCount++
	}

	return indexedCount, nil
}

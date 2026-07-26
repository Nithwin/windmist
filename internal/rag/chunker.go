package rag

import (
	"bufio"
	"strings"
)

// Chunk represents a section of a source file.
type Chunk struct {
	FilePath  string
	StartLine int
	EndLine   int
	Content   string
}

// ChunkConfig controls how files are split into chunks.
type ChunkConfig struct {
	// MaxChunkLines is the maximum number of lines per chunk.
	MaxChunkLines int
	// OverlapLines is how many lines overlap between adjacent chunks.
	OverlapLines int
}

// DefaultChunkConfig returns sensible defaults for code chunking.
func DefaultChunkConfig() ChunkConfig {
	return ChunkConfig{
		MaxChunkLines: 40,
		OverlapLines:  5,
	}
}

// ChunkFile splits a file's content into overlapping chunks.
// It uses a simple line-based sliding window approach that respects
// blank-line boundaries (tries to split at natural breaks).
func ChunkFile(filePath, content string, cfg ChunkConfig) []Chunk {
	if cfg.MaxChunkLines <= 0 {
		cfg.MaxChunkLines = 40
	}
	if cfg.OverlapLines < 0 {
		cfg.OverlapLines = 0
	}

	lines := splitLines(content)
	if len(lines) == 0 {
		return nil
	}

	// If the whole file fits in one chunk, return it as-is.
	if len(lines) <= cfg.MaxChunkLines {
		return []Chunk{
			{
				FilePath:  filePath,
				StartLine: 1,
				EndLine:   len(lines),
				Content:   content,
			},
		}
	}

	var chunks []Chunk
	start := 0
	for start < len(lines) {
		end := start + cfg.MaxChunkLines
		if end > len(lines) {
			end = len(lines)
		}

		// Try to find a natural break point (blank line) near the end
		// to avoid splitting mid-function.
		bestBreak := end
		if end < len(lines) {
			for i := end - 1; i > start+cfg.MaxChunkLines/2; i-- {
				if strings.TrimSpace(lines[i]) == "" {
					bestBreak = i + 1
					break
				}
			}
			end = bestBreak
		}

		chunkContent := strings.Join(lines[start:end], "\n")
		chunks = append(chunks, Chunk{
			FilePath:  filePath,
			StartLine: start + 1, // 1-indexed
			EndLine:   end,
			Content:   chunkContent,
		})

		// Advance with overlap
		step := end - start - cfg.OverlapLines
		if step < 1 {
			step = 1
		}
		start += step
	}

	return chunks
}

// splitLines splits text into lines, preserving empty lines.
func splitLines(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

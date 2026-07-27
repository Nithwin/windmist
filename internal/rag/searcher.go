package rag

import (
	"sort"
)

// SearchResult represents a matched chunk with its similarity score.
type SearchResult struct {
	Chunk      Chunk
	Similarity float32
}

// Searcher handles semantic queries against the document store.
type Searcher struct {
	store    *DocumentStore
	embedder *TFIDFEmbedder
}

// NewSearcher creates a new Searcher.
func NewSearcher(store *DocumentStore, embedder *TFIDFEmbedder) *Searcher {
	return &Searcher{
		store:    store,
		embedder: embedder,
	}
}

// Search finds the top-K most similar chunks to the query string.
func (s *Searcher) Search(query string, topK int) ([]SearchResult, error) {
	if topK <= 0 {
		topK = 5
	}

	queryVector := s.embedder.Embed(query)
	if queryVector == nil {
		return nil, nil // Empty embedder or query
	}

	allChunks, err := s.store.GetAllChunks()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, ic := range allChunks {
		chunkVector := DecodeVector(ic.Vector)
		if chunkVector == nil {
			continue
		}

		similarity := CosineSimilarity(queryVector, chunkVector)

		// Only consider positive similarities above a small threshold
		if similarity > 0.05 {
			results = append(results, SearchResult{
				Chunk: Chunk{
					FilePath:  ic.FilePath,
					StartLine: ic.StartLine,
					EndLine:   ic.EndLine,
					Content:   ic.Content,
				},
				Similarity: similarity,
			})
		}
	}

	// Sort descending by similarity
	sort.Slice(results, func(i, j int) bool {
		return results[i].Similarity > results[j].Similarity
	})

	// Take top K
	if len(results) > topK {
		results = results[:topK]
	}

	return results, nil
}

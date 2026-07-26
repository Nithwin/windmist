package rag

import (
	"math"
	"strings"
	"unicode"
)

// TFIDFEmbedder generates TF-IDF-based embeddings entirely in pure Go.
// No API calls needed — works fully offline.
type TFIDFEmbedder struct {
	// vocabulary maps tokens to their dimension index.
	vocabulary map[string]int
	// idf stores inverse document frequency for each token.
	idf map[string]float64
	// dimensions is the embedding vector size (vocabulary size, capped).
	dimensions int
	// maxDimensions caps the vector size for memory efficiency.
	maxDimensions int
}

// NewTFIDFEmbedder creates a new TF-IDF embedder.
func NewTFIDFEmbedder(maxDimensions int) *TFIDFEmbedder {
	if maxDimensions <= 0 {
		maxDimensions = 512
	}
	return &TFIDFEmbedder{
		vocabulary:    make(map[string]int),
		idf:           make(map[string]float64),
		maxDimensions: maxDimensions,
	}
}

// BuildVocabulary builds a vocabulary from a corpus of documents.
// Each document is a string of text (e.g., a code chunk).
// This must be called before Embed().
func (e *TFIDFEmbedder) BuildVocabulary(documents []string) {
	docFreq := make(map[string]int)
	allTokens := make(map[string]bool)

	for _, doc := range documents {
		tokens := tokenize(doc)
		seen := make(map[string]bool)
		for _, tok := range tokens {
			allTokens[tok] = true
			if !seen[tok] {
				docFreq[tok]++
				seen[tok] = true
			}
		}
	}

	// Build vocabulary — pick the top tokens by document frequency.
	// This acts as a natural feature selection for the most relevant terms.
	type tokenFreq struct {
		token string
		freq  int
	}
	ranked := make([]tokenFreq, 0, len(allTokens))
	for tok := range allTokens {
		ranked = append(ranked, tokenFreq{tok, docFreq[tok]})
	}

	// Sort by frequency (descending), but skip tokens that appear in
	// too many documents (>80%) as they're not discriminative.
	totalDocs := len(documents)
	filtered := make([]tokenFreq, 0, len(ranked))
	for _, tf := range ranked {
		ratio := float64(tf.freq) / float64(totalDocs)
		if ratio < 0.8 && tf.freq > 1 {
			filtered = append(filtered, tf)
		}
	}

	// Sort by frequency (descending) using a simple selection sort
	// for the top maxDimensions entries.
	dim := e.maxDimensions
	if dim > len(filtered) {
		dim = len(filtered)
	}
	for i := 0; i < dim; i++ {
		maxIdx := i
		for j := i + 1; j < len(filtered); j++ {
			if filtered[j].freq > filtered[maxIdx].freq {
				maxIdx = j
			}
		}
		filtered[i], filtered[maxIdx] = filtered[maxIdx], filtered[i]
	}

	e.vocabulary = make(map[string]int, dim)
	for i := 0; i < dim; i++ {
		e.vocabulary[filtered[i].token] = i
	}
	e.dimensions = dim

	// Compute IDF for each token in vocabulary
	e.idf = make(map[string]float64, dim)
	for tok := range e.vocabulary {
		df := docFreq[tok]
		if df == 0 {
			df = 1
		}
		e.idf[tok] = math.Log(float64(totalDocs+1) / float64(df+1))
	}
}

// Embed generates a TF-IDF vector for the given text.
// The vector dimensions correspond to the vocabulary built via BuildVocabulary.
func (e *TFIDFEmbedder) Embed(text string) Vector {
	if e.dimensions == 0 {
		return nil
	}

	tokens := tokenize(text)
	if len(tokens) == 0 {
		return make(Vector, e.dimensions)
	}

	// Compute term frequency
	tf := make(map[string]int)
	for _, tok := range tokens {
		tf[tok]++
	}

	// Build TF-IDF vector
	vec := make(Vector, e.dimensions)
	for tok, count := range tf {
		idx, ok := e.vocabulary[tok]
		if !ok {
			continue
		}
		// TF: normalized by document length
		termFreq := float64(count) / float64(len(tokens))
		// IDF from pre-computed values
		idf := e.idf[tok]
		if idf == 0 {
			idf = 1
		}
		vec[idx] = float32(termFreq * idf)
	}

	return Normalize(vec)
}

// Dimensions returns the number of dimensions in the embeddings.
func (e *TFIDFEmbedder) Dimensions() int {
	return e.dimensions
}

// tokenize splits text into code-aware tokens.
// It handles camelCase, snake_case, and common programming constructs.
func tokenize(text string) []string {
	var tokens []string
	text = strings.ToLower(text)

	// Split on non-alphanumeric boundaries
	var current strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
		} else {
			if current.Len() > 0 {
				tok := current.String()
				if len(tok) > 1 && !isStopWord(tok) {
					tokens = append(tokens, tok)
				}
				current.Reset()
			}
		}
	}
	if current.Len() > 0 {
		tok := current.String()
		if len(tok) > 1 && !isStopWord(tok) {
			tokens = append(tokens, tok)
		}
	}

	// Also split camelCase tokens
	expanded := make([]string, 0, len(tokens)*2)
	for _, tok := range tokens {
		expanded = append(expanded, tok)
		parts := splitCamelCase(tok)
		if len(parts) > 1 {
			for _, p := range parts {
				if len(p) > 1 {
					expanded = append(expanded, p)
				}
			}
		}
	}

	return expanded
}

// splitCamelCase splits "camelCase" into ["camel", "case"].
func splitCamelCase(s string) []string {
	var parts []string
	var current strings.Builder
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			if current.Len() > 0 {
				parts = append(parts, strings.ToLower(current.String()))
				current.Reset()
			}
		}
		current.WriteRune(r)
	}
	if current.Len() > 0 {
		parts = append(parts, strings.ToLower(current.String()))
	}
	return parts
}

// isStopWord returns true for common programming/English stop words.
func isStopWord(w string) bool {
	stops := map[string]bool{
		"the": true, "is": true, "at": true, "in": true, "on": true,
		"to": true, "of": true, "an": true, "if": true, "or": true,
		"it": true, "be": true, "as": true, "do": true, "no": true,
		"so": true, "we": true, "he": true, "by": true, "up": true,
		"my": true, "me": true, "am": true, "go": true,
	}
	return stops[w]
}

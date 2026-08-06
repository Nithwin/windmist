package rag

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndexProject_DoesNotSkipRootDot(t *testing.T) {
	dir := t.TempDir()
	var b strings.Builder
	for i := 0; i < 40; i++ {
		fmt.Fprintf(&b, "package main\nfunc Helper%d() {\n\tprintln(\"hello world from helper %d\")\n}\n", i, i)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(b.String()), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".git", "config"), []byte("ignored"), 0644); err != nil {
		t.Fatal(err)
	}

	dbPath := filepath.Join(dir, "rag.db")
	store, err := NewDocumentStore(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	embedder := NewTFIDFEmbedder(64)
	indexer := NewIndexer(store, embedder)

	// Reproduce the historical bug: indexing with relative "." as root.
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(origWD) })

	count, err := indexer.IndexProject(".")
	if err != nil {
		t.Fatalf("IndexProject: %v", err)
	}
	if count == 0 {
		t.Fatal("expected at least one indexed chunk; root '.' was likely skipped")
	}
}

func TestBuildVocabulary_Deterministic(t *testing.T) {
	docs := []string{
		"alpha beta gamma delta epsilon",
		"alpha beta gamma delta zeta",
		"alpha beta gamma theta iota",
		"alpha beta unique word here",
	}
	e1 := NewTFIDFEmbedder(8)
	e2 := NewTFIDFEmbedder(8)
	e1.BuildVocabulary(docs)
	e2.BuildVocabulary(docs)

	v1 := e1.Embed("alpha beta gamma")
	v2 := e2.Embed("alpha beta gamma")
	if len(v1) != len(v2) {
		t.Fatalf("dimension mismatch: %d vs %d", len(v1), len(v2))
	}
	for i := range v1 {
		if v1[i] != v2[i] {
			t.Fatalf("non-deterministic embedding at index %d: %v vs %v", i, v1[i], v2[i])
		}
	}
}

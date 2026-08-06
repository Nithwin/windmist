package rag

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// DocumentStore manages the storage and retrieval of RAG chunks and vectors.
type DocumentStore struct {
	db *sqlx.DB
}

// IndexedChunk represents a chunk retrieved from the database.
type IndexedChunk struct {
	ID        int       `db:"id"`
	FilePath  string    `db:"file_path"`
	StartLine int       `db:"start_line"`
	EndLine   int       `db:"end_line"`
	Content   string    `db:"content"`
	Vector    []byte    `db:"vector"`
	IndexedAt time.Time `db:"indexed_at"`
}

var ragSchema = `
CREATE TABLE IF NOT EXISTS rag_chunks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT NOT NULL,
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    content TEXT NOT NULL,
    vector BLOB NOT NULL,
    indexed_at DATETIME NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rag_chunks_file_path ON rag_chunks(file_path);
`

// NewDocumentStore creates or opens the RAG SQLite database.
func NewDocumentStore(dbPath string) (*DocumentStore, error) {
	db, err := sqlx.Connect("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("connect rag db: %w", err)
	}

	if _, err := db.Exec(ragSchema); err != nil {
		return nil, fmt.Errorf("apply rag schema: %w", err)
	}

	return &DocumentStore{db: db}, nil
}

// Close closes the database connection.
func (s *DocumentStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Clear removes all existing chunks for a fresh re-index.
func (s *DocumentStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM rag_chunks")
	return err
}

// InsertChunk stores a single chunk and its embedding vector.
func (s *DocumentStore) InsertChunk(chunk Chunk, vector Vector) error {
	query := `
		INSERT INTO rag_chunks (file_path, start_line, end_line, content, vector, indexed_at)
		VALUES (:file_path, :start_line, :end_line, :content, :vector, :indexed_at)
	`

	record := IndexedChunk{
		FilePath:  chunk.FilePath,
		StartLine: chunk.StartLine,
		EndLine:   chunk.EndLine,
		Content:   chunk.Content,
		Vector:    EncodeVector(vector),
		IndexedAt: time.Now(),
	}

	_, err := s.db.NamedExec(query, record)
	return err
}

// GetAllChunks retrieves all chunks for similarity comparison.
func (s *DocumentStore) GetAllChunks() ([]IndexedChunk, error) {
	var chunks []IndexedChunk
	err := s.db.Select(&chunks, "SELECT * FROM rag_chunks")
	return chunks, err
}

// GetChunksByFile retrieves all chunks for a specific file.
func (s *DocumentStore) GetChunksByFile(filePath string) ([]IndexedChunk, error) {
	var chunks []IndexedChunk
	err := s.db.Select(&chunks, "SELECT * FROM rag_chunks WHERE file_path = ? ORDER BY start_line ASC", filePath)
	return chunks, err
}

// UpdateChunkVector replaces the embedding for an existing chunk so query
// vectors stay aligned after vocabulary rebuilds.
func (s *DocumentStore) UpdateChunkVector(id int, vector Vector) error {
	_, err := s.db.Exec("UPDATE rag_chunks SET vector = ? WHERE id = ?", EncodeVector(vector), id)
	return err
}

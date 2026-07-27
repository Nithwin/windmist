package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var schema = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    title TEXT,
    project_path TEXT,
    provider TEXT,
    model TEXT,
    agent_mode TEXT,
    token_count INTEGER DEFAULT 0,
    cost_estimate REAL DEFAULT 0.0,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    role TEXT,
    content TEXT,
    tool_calls TEXT,
    tool_results TEXT,
    token_count INTEGER DEFAULT 0,
    created_at DATETIME
);

CREATE TABLE IF NOT EXISTS file_changes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id TEXT REFERENCES sessions(id) ON DELETE CASCADE,
    message_id INTEGER REFERENCES messages(id) ON DELETE CASCADE,
    batch_id TEXT DEFAULT '',
    file_path TEXT,
    change_type TEXT,
    before_content TEXT,
    after_content TEXT,
    undone BOOLEAN DEFAULT 0,
    created_at DATETIME
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_messages_session_id ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_file_changes_session_id ON file_changes(session_id);
CREATE INDEX IF NOT EXISTS idx_file_changes_message_id ON file_changes(message_id);
CREATE INDEX IF NOT EXISTS idx_file_changes_batch_id ON file_changes(batch_id);
`

type Store struct {
	db *sqlx.DB
}

// NewStore initializes the SQLite database at ~/.windmist/sessions.db
func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home dir: %w", err)
	}

	windmistDir := filepath.Join(home, ".windmist")
	if err := os.MkdirAll(windmistDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	dbPath := filepath.Join(windmistDir, "sessions.db")

	// Enable foreign keys
	db, err := sqlx.Connect("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Apply schema
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// NewStoreForTest creates a new Store with a specific path for testing
func NewStoreForTest(dbPath string) (*Store, error) {
	// Enable foreign keys
	db, err := sqlx.Connect("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Apply schema
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to apply schema: %w", err)
	}

	return &Store{db: db}, nil
}

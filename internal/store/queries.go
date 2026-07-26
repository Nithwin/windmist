package store

import (
	"fmt"
	"time"
)

// CreateSession creates a new session in the database
func (s *Store) CreateSession(session *Session) error {
	session.CreatedAt = time.Now()
	session.UpdatedAt = session.CreatedAt

	query := `
		INSERT INTO sessions (id, title, project_path, provider, model, agent_mode, token_count, cost_estimate, created_at, updated_at)
		VALUES (:id, :title, :project_path, :provider, :model, :agent_mode, :token_count, :cost_estimate, :created_at, :updated_at)
	`
	_, err := s.db.NamedExec(query, session)
	return err
}

// GetSession retrieves a session by ID
func (s *Store) GetSession(id string) (*Session, error) {
	var session Session
	err := s.db.Get(&session, "SELECT * FROM sessions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessionsByProject gets all sessions for a specific project
func (s *Store) ListSessionsByProject(projectPath string) ([]Session, error) {
	var sessions []Session
	err := s.db.Select(&sessions, "SELECT * FROM sessions WHERE project_path = ? ORDER BY updated_at DESC", projectPath)
	return sessions, err
}

// UpdateSession updates the metadata of a session
func (s *Store) UpdateSession(session *Session) error {
	session.UpdatedAt = time.Now()
	query := `
		UPDATE sessions 
		SET title = :title, provider = :provider, model = :model, agent_mode = :agent_mode, token_count = :token_count, cost_estimate = :cost_estimate, updated_at = :updated_at
		WHERE id = :id
	`
	_, err := s.db.NamedExec(query, session)
	return err
}

// SaveMessage stores a new message and returns its ID
func (s *Store) SaveMessage(msg *Message) error {
	msg.CreatedAt = time.Now()

	query := `
		INSERT INTO messages (session_id, role, content, tool_calls, tool_results, token_count, created_at)
		VALUES (:session_id, :role, :content, :tool_calls, :tool_results, :token_count, :created_at)
	`
	res, err := s.db.NamedExec(query, msg)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		msg.ID = int(id)
	}

	// Update the session's updated_at timestamp
	_, _ = s.db.Exec("UPDATE sessions SET updated_at = ? WHERE id = ?", msg.CreatedAt, msg.SessionID)

	return nil
}

// GetMessagesBySession gets all messages for a session, ordered by creation time
func (s *Store) GetMessagesBySession(sessionID string) ([]Message, error) {
	var messages []Message
	err := s.db.Select(&messages, "SELECT * FROM messages WHERE session_id = ? ORDER BY id ASC", sessionID)
	return messages, err
}

// SaveFileChange logs a file change for undo/redo
func (s *Store) SaveFileChange(change *FileChange) error {
	change.CreatedAt = time.Now()

	query := `
		INSERT INTO file_changes (session_id, message_id, batch_id, file_path, change_type, before_content, after_content, undone, created_at)
		VALUES (:session_id, :message_id, :batch_id, :file_path, :change_type, :before_content, :after_content, :undone, :created_at)
	`
	res, err := s.db.NamedExec(query, change)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		change.ID = int(id)
	}

	return nil
}

// GetFileChangesBySession retrieves all file changes in a session
func (s *Store) GetFileChangesBySession(sessionID string) ([]FileChange, error) {
	var changes []FileChange
	err := s.db.Select(&changes, "SELECT * FROM file_changes WHERE session_id = ? ORDER BY id ASC", sessionID)
	return changes, err
}

// GetLastFileChange gets the most recent file change for a session
func (s *Store) GetLastFileChange(sessionID string) (*FileChange, error) {
	var change FileChange
	err := s.db.Get(&change, "SELECT * FROM file_changes WHERE session_id = ? ORDER BY id DESC LIMIT 1", sessionID)
	if err != nil {
		return nil, err
	}
	return &change, nil
}

// GetLastBatchForUndo gets the most recent batch of changes that haven't been undone
func (s *Store) GetLastBatchForUndo(sessionID string) ([]FileChange, error) {
	var batchID string
	err := s.db.Get(&batchID, "SELECT batch_id FROM file_changes WHERE session_id = ? AND undone = 0 ORDER BY id DESC LIMIT 1", sessionID)
	if err != nil {
		return nil, err
	}

	var changes []FileChange
	err = s.db.Select(&changes, "SELECT * FROM file_changes WHERE session_id = ? AND batch_id = ? ORDER BY id DESC", sessionID, batchID)
	return changes, err
}

// GetNextBatchForRedo gets the oldest batch of changes that are currently undone
func (s *Store) GetNextBatchForRedo(sessionID string) ([]FileChange, error) {
	var batchID string
	err := s.db.Get(&batchID, "SELECT batch_id FROM file_changes WHERE session_id = ? AND undone = 1 ORDER BY id ASC LIMIT 1", sessionID)
	if err != nil {
		return nil, err
	}

	var changes []FileChange
	err = s.db.Select(&changes, "SELECT * FROM file_changes WHERE session_id = ? AND batch_id = ? ORDER BY id ASC", sessionID, batchID)
	return changes, err
}

// SetBatchUndoneState updates the undone status of a batch
func (s *Store) SetBatchUndoneState(sessionID string, batchID string, undone bool) error {
	_, err := s.db.Exec("UPDATE file_changes SET undone = ? WHERE session_id = ? AND batch_id = ?", undone, sessionID, batchID)
	return err
}

// ClearRedoHistory removes all file changes that are currently undone for a session
func (s *Store) ClearRedoHistory(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM file_changes WHERE session_id = ? AND undone = 1", sessionID)
	return err
}

// DeleteSession completely deletes a session and all cascading data
func (s *Store) DeleteSession(id string) error {
	res, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

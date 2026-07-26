package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStoreIntegration(t *testing.T) {
	// Temporarily mock user home dir to a temp directory
	tempHome, err := os.MkdirTemp("", "windmist_home_*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempHome)

	// Since NewStore relies on os.UserHomeDir(), we just mock the db path directly for testing
	windmistDir := filepath.Join(tempHome, ".windmist")
	os.MkdirAll(windmistDir, 0755)

	dbPath := filepath.Join(windmistDir, "sessions.db")

	store, err := NewStoreForTest(dbPath)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	defer store.Close()

	// Test CreateSession
	session := &Session{
		ID:          "sess_123",
		Title:       "Test Session",
		ProjectPath: "/home/user/project",
		Provider:    "openai",
		Model:       "gpt-4",
		AgentMode:   "build",
	}

	err = store.CreateSession(session)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Test GetSession
	retrieved, err := store.GetSession("sess_123")
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}
	if retrieved.Title != "Test Session" {
		t.Fatalf("expected title 'Test Session', got %s", retrieved.Title)
	}

	// Test SaveMessage
	msg := &Message{
		SessionID: "sess_123",
		Role:      "user",
		Content:   "Hello world",
	}
	err = store.SaveMessage(msg)
	if err != nil {
		t.Fatalf("failed to save message: %v", err)
	}
	if msg.ID == 0 {
		t.Fatal("expected message ID to be set")
	}

	// Test GetMessages
	messages, err := store.GetMessagesBySession("sess_123")
	if err != nil {
		t.Fatalf("failed to get messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "Hello world" {
		t.Fatalf("expected message content 'Hello world', got %s", messages[0].Content)
	}

	// Test DeleteSession (should cascade and delete messages too, assuming SQLite foreign keys are enabled)
	err = store.DeleteSession("sess_123")
	if err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	// Verify messages are deleted
	messages, _ = store.GetMessagesBySession("sess_123")
	if len(messages) != 0 {
		t.Fatalf("expected 0 messages after cascade delete, got %d", len(messages))
	}
}

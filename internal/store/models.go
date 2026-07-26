package store

import (
	"time"
)

type Session struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	ProjectPath  string    `db:"project_path"`
	Provider     string    `db:"provider"`
	Model        string    `db:"model"`
	AgentMode    string    `db:"agent_mode"`
	TokenCount   int       `db:"token_count"`
	CostEstimate float64   `db:"cost_estimate"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Message struct {
	ID          int       `db:"id"`
	SessionID   string    `db:"session_id"`
	Role        string    `db:"role"` // user, assistant, tool, system
	Content     string    `db:"content"`
	ToolCalls   string    `db:"tool_calls"`   // JSON encoded
	ToolResults string    `db:"tool_results"` // JSON encoded
	TokenCount  int       `db:"token_count"`
	CreatedAt   time.Time `db:"created_at"`
}

type FileChange struct {
	ID            int       `db:"id"`
	SessionID     string    `db:"session_id"`
	MessageID     int       `db:"message_id"`
	FilePath      string    `db:"file_path"`
	ChangeType    string    `db:"change_type"` // create, edit, delete
	BeforeContent string    `db:"before_content"`
	AfterContent  string    `db:"after_content"`
	CreatedAt     time.Time `db:"created_at"`
}

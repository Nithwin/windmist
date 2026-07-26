package defaults

import (
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/editing"
	"github.com/Nithwin/WindMist/internal/tools/filesystem"
	"github.com/Nithwin/WindMist/internal/tools/system"
)

// RegisterAll registers the small basic tool set used by the beginner version.
func RegisterAll(m *tools.Manager, _ system.ApprovalCallback, _ *config.Config) {
	if m == nil {
		return
	}

	// Basic file tools.
	m.Register(filesystem.NewReadTool())
	m.Register(filesystem.NewWriteTool())
	m.Register(filesystem.NewListTool())

	// Keep only one simple editing tool so the project stays easy to follow.
	m.Register(editing.NewReplaceTextTool())
}

package defaults

import (
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/editing"
	"github.com/Nithwin/WindMist/internal/tools/filesystem"
)

// RegisterAll registers all built-in filesystem and editing tools onto the manager.
func RegisterAll(m *tools.Manager) {
	if m == nil {
		return
	}

	// Filesystem tools
	m.Register(filesystem.NewReadTool())
	m.Register(filesystem.NewWriteTool())
	m.Register(filesystem.NewDeleteTool())
	m.Register(filesystem.NewListTool())
	m.Register(filesystem.NewRenameTool())
	m.Register(filesystem.NewAppendTool())
	m.Register(filesystem.NewCreateTool())
	m.Register(filesystem.NewInfoTool())
	m.Register(filesystem.NewExistsTool())

	// Editing tools
	m.Register(editing.NewReplaceTextTool())
	m.Register(editing.NewReplaceRangeTool())
	m.Register(editing.NewDeleteRangeTool())
	m.Register(editing.NewReadContextTool())
	m.Register(editing.NewInsertTextTool())
	m.Register(editing.NewSearchTool())
}

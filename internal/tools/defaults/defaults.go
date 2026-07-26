package defaults

import (
	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/tools"
	"github.com/Nithwin/WindMist/internal/tools/agent"
	"github.com/Nithwin/WindMist/internal/tools/editing"
	"github.com/Nithwin/WindMist/internal/tools/filesystem"
	"github.com/Nithwin/WindMist/internal/tools/system"
	"github.com/Nithwin/WindMist/internal/tools/web"
)

// RegisterAll registers all built-in filesystem and editing tools onto the manager.
func RegisterAll(m *tools.Manager, approvalCb system.ApprovalCallback, cfg *config.Config) {
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
	m.Register(filesystem.NewGlobTool())
	m.Register(filesystem.NewGrepTool())

	// Editing tools
	m.Register(editing.NewReplaceTextTool())
	m.Register(editing.NewReplaceRangeTool())
	m.Register(editing.NewDeleteRangeTool())
	m.Register(editing.NewReadContextTool())
	m.Register(editing.NewInsertTextTool())
	m.Register(editing.NewSearchTool())
	m.Register(editing.NewPatchTool())

	// System tools
	m.Register(system.NewCommandTool(approvalCb))
	m.Register(system.NewGitTool(approvalCb))

	// Web tools
	m.Register(web.NewWebSearchTool())
	m.Register(web.NewFetchTool())

	// Agent tools
	m.Register(agent.NewTodoTool())
	if cfg != nil {
		m.Register(agent.NewSubAgentTool(cfg))
	}
}

package tools

type Manager struct {
	tools map[string]Tool
}

func NewManager() *Manager {
	return &Manager{
		tools: make(map[string]Tool),
	}
}

func (m *Manager) Register(tool Tool) {
	def := tool.Definition()
	m.tools[def.Name] = tool
}

func (m *Manager) Get(name string) (Tool, bool) {
	tool, ok := m.tools[name]
	return tool, ok
}

func (m *Manager) List() []Tool {
	list := make([]Tool, 0, len(m.tools))

	for _, tool := range m.tools {
		list = append(list, tool)
	}

	return list
}

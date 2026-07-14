package ui

import "github.com/charmbracelet/glamour"

type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
}

func NewMarkdownRenderer() (*MarkdownRenderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(0),
	)

	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		renderer: r,
	}, nil
}

func (m *MarkdownRenderer) Render(text string) string {
	if m == nil || m.renderer == nil {
		return text
	}

	out, err := m.renderer.Render(text)
	if err != nil {
		return text
	}

	return out
}
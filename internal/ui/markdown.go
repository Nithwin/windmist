package ui

import "github.com/charmbracelet/glamour"

// windmistStyle is a minimal, clean Glamour style for WindMist.
// Plain white text, bold headings, bordered code blocks, no flashy colors.
var windmistStyle = []byte(`{
	"document": {
		"margin": 0
	},
	"block_quote": {
		"indent": 2,
		"color": "248"
	},
	"paragraph": {
		"color": "252"
	},
	"list": {
		"level_indent": 2,
		"color": "252"
	},
	"heading": {
		"bold": true,
		"color": "255"
	},
	"h1": {
		"bold": true,
		"color": "255"
	},
	"h2": {
		"bold": true,
		"color": "255"
	},
	"h3": {
		"bold": true,
		"color": "255"
	},
	"h4": {
		"bold": true,
		"color": "255"
	},
	"h5": {
		"bold": true,
		"color": "255"
	},
	"h6": {
		"bold": true,
		"color": "255"
	},
	"text": {
		"color": "252"
	},
	"strikethrough": {
		"crossed_out": true
	},
	"emph": {
		"italic": true,
		"color": "252"
	},
	"strong": {
		"bold": true,
		"color": "255"
	},
	"hr": {
		"color": "240",
		"format": "\n─────────────────────────────────────────\n"
	},
	"item": {
		"color": "252",
		"block_prefix": "• "
	},
	"enumeration": {
		"color": "252"
	},
	"task": {
		"ticked": "[✓] ",
		"unticked": "[ ] "
	},
	"link": {
		"color": "39",
		"underline": true
	},
	"link_text": {
		"color": "39",
		"bold": true
	},
	"image": {
		"color": "252"
	},
	"image_text": {
		"color": "252"
	},
	"code": {
		"color": "203",
		"background_color": "236"
	},
	"code_block": {
		"color": "252",
		"margin": 1,
		"chroma": {
			"text": {
				"color": "#d6d6d6"
			},
			"comment": {
				"color": "#767676"
			},
			"keyword": {
				"color": "#00afff",
				"bold": true
			},
			"literal": {
				"color": "#ffaf5f"
			},
			"name": {
				"color": "#d6d6d6"
			},
			"literal_string": {
				"color": "#87d787"
			},
			"literal_number": {
				"color": "#ffaf5f"
			},
			"name_function": {
				"color": "#00afff"
			},
			"name_class": {
				"color": "#00afff",
				"bold": true
			},
			"name_builtin": {
				"color": "#5fd7ff"
			},
			"operator": {
				"color": "#d6d6d6"
			},
			"punctuation": {
				"color": "#d6d6d6"
			},
			"generic_inserted": {
				"color": "#87d787"
			},
			"generic_deleted": {
				"color": "#ff5f5f"
			}
		}
	},
	"table": {
		"color": "252"
	},
	"definition_list": {},
	"definition_term": {
		"bold": true
	},
	"definition_description": {
		"block_prefix": " "
	},
	"html_block": {},
	"html_span": {}
}`)

type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
}

func NewMarkdownRenderer() (*MarkdownRenderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(windmistStyle),
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

func (m *MarkdownRenderer) RenderWithWidth(text string, width int) string {
	if m == nil {
		return text
	}
	if width <= 0 {
		width = 80
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithStylesFromJSONBytes(windmistStyle),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return m.Render(text)
	}
	out, err := r.Render(text)
	if err != nil {
		return m.Render(text)
	}
	return out
}

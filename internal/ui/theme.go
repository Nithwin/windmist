package ui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

type ThemeColors struct {
	Background string `yaml:"background"`
	Foreground string `yaml:"foreground"`
	Accent     string `yaml:"accent"`
	Success    string `yaml:"success"`
	Error      string `yaml:"error"`
	Warning    string `yaml:"warning"`
	Info       string `yaml:"info"`
	Muted      string `yaml:"muted"`
	Border     string `yaml:"border"`
	Selection  string `yaml:"selection"`
}

type Theme struct {
	Name   string      `yaml:"name"`
	Colors ThemeColors `yaml:"colors"`
}

var BuiltinThemes = map[string]Theme{
	"catppuccin": {
		Name: "catppuccin",
		Colors: ThemeColors{
			Background: "#1e1e2e",
			Foreground: "#cdd6f4",
			Accent:     "#cba6f7",
			Success:    "#a6e3a1",
			Error:      "#f38ba8",
			Warning:    "#f9e2af",
			Info:       "#89b4fa",
			Muted:      "#6c7086",
			Border:     "#313244",
			Selection:  "#313244",
		},
	},
	"dracula": {
		Name: "dracula",
		Colors: ThemeColors{
			Background: "#282a36",
			Foreground: "#f8f8f2",
			Accent:     "#bd93f9",
			Success:    "#50fa7b",
			Error:      "#ff5555",
			Warning:    "#f1fa8c",
			Info:       "#8be9fd",
			Muted:      "#6272a4",
			Border:     "#44475a",
			Selection:  "#44475a",
		},
	},
	"gruvbox": {
		Name: "gruvbox",
		Colors: ThemeColors{
			Background: "#282828",
			Foreground: "#ebdbb2",
			Accent:     "#d3869b",
			Success:    "#b8bb26",
			Error:      "#fb4934",
			Warning:    "#fabd2f",
			Info:       "#83a598",
			Muted:      "#928374",
			Border:     "#504945",
			Selection:  "#504945",
		},
	},
	"nord": {
		Name: "nord",
		Colors: ThemeColors{
			Background: "#2e3440",
			Foreground: "#d8dee9",
			Accent:     "#b48ead",
			Success:    "#a3be8c",
			Error:      "#bf616a",
			Warning:    "#ebcb8b",
			Info:       "#81a1c1",
			Muted:      "#4c566a",
			Border:     "#3b4252",
			Selection:  "#3b4252",
		},
	},
	"tokyo-night": {
		Name: "tokyo-night",
		Colors: ThemeColors{
			Background: "#1a1b26",
			Foreground: "#c0caf5",
			Accent:     "#bb9af7",
			Success:    "#9ece6a",
			Error:      "#f7768e",
			Warning:    "#e0af68",
			Info:       "#7aa2f7",
			Muted:      "#565f89",
			Border:     "#292e42",
			Selection:  "#292e42",
		},
	},
	"solarized": {
		Name: "solarized",
		Colors: ThemeColors{
			Background: "#002b36",
			Foreground: "#839496",
			Accent:     "#6c71c4",
			Success:    "#859900",
			Error:      "#dc322f",
			Warning:    "#b58900",
			Info:       "#268bd2",
			Muted:      "#586e75",
			Border:     "#073642",
			Selection:  "#073642",
		},
	},
	"monokai": {
		Name: "monokai",
		Colors: ThemeColors{
			Background: "#272822",
			Foreground: "#f8f8f2",
			Accent:     "#ae81ff",
			Success:    "#a6e22e",
			Error:      "#f92672",
			Warning:    "#e6db74",
			Info:       "#66d9ef",
			Muted:      "#75715e",
			Border:     "#3e3d32",
			Selection:  "#3e3d32",
		},
	},
	"windmist": {
		Name: "windmist",
		Colors: ThemeColors{
			Background: "#161122", // Deep Premium Amethyst
			Foreground: "#F8FAFC", // Crisp white
			Accent:     "#00F0FF", // Neon Cyan (Main)
			Success:    "#00FF9D", // Neon Mint
			Error:      "#FF006A", // Neon Crimson
			Warning:    "#FFB800", // Neon Gold
			Info:       "#D946EF", // Vibrant Fuchsia
			Muted:      "#6B6282", // Muted purple/gray
			Border:     "#2D243F", // Deep purple border
			Selection:  "#2D243F",
		},
	},
}

var CurrentThemeName = "windmist"

func ApplyTheme(t Theme) {
	CurrentThemeName = t.Name

	Purple = lipgloss.Color(t.Colors.Accent)
	PurpleDark = lipgloss.Color(t.Colors.Accent)
	PurpleDim = lipgloss.Color(t.Colors.Accent)
	Cyan = lipgloss.Color(t.Colors.Info)
	Green = lipgloss.Color(t.Colors.Success)
	Amber = lipgloss.Color(t.Colors.Warning)
	Red = lipgloss.Color(t.Colors.Error)
	Muted = lipgloss.Color(t.Colors.Muted)
	MutedLight = lipgloss.Color(t.Colors.Foreground) 
	Surface = lipgloss.Color(t.Colors.Background)
	White = lipgloss.Color(t.Colors.Foreground)
	Border = lipgloss.Color(t.Colors.Border)
	Selection = lipgloss.Color(t.Colors.Selection)

	UpdateStyles()
}

func LoadTheme(name string, customDir string) error {
	if name == "" {
		name = "windmist"
	}

	// First check built-in themes
	if t, ok := BuiltinThemes[name]; ok {
		ApplyTheme(t)
		return nil
	}

	// Try loading from customDir
	themePath := filepath.Join(customDir, name+".yaml")
	data, err := os.ReadFile(themePath)
	if err != nil {
		ApplyTheme(BuiltinThemes["windmist"])
		return fmt.Errorf("theme %s not found: %w", name, err)
	}

	var t Theme
	if err := yaml.Unmarshal(data, &t); err != nil {
		return fmt.Errorf("failed to parse theme %s: %w", name, err)
	}

	ApplyTheme(t)
	return nil
}

func AvailableThemes() []string {
	var themes []string
	for k := range BuiltinThemes {
		themes = append(themes, k)
	}
	return themes
}

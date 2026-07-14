package editing

import (
	"fmt"
	"strings"
)

// ToUnifiedDiff formats the entire patch into a human-readable unified-style diff preview.
// This allows users and interactive prompts to inspect exactly what lines will be removed and added before applying modifications.
func (p *Patch) ToUnifiedDiff() string {
	if p == nil || len(p.Files) == 0 {
		return ""
	}

	var sb strings.Builder
	for i, fp := range p.Files {
		if len(fp.Operations) == 0 {
			continue
		}

		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("--- %s\n+++ %s\n", fp.Path, fp.Path))

		for _, op := range fp.Operations {
			switch op.Type {
			case Replace:
				sb.WriteString(fmt.Sprintf("@@ line %d @@\n", op.StartLine))
				if op.OldText != "" {
					oldLines := strings.Split(strings.TrimSuffix(op.OldText, "\n"), "\n")
					for _, line := range oldLines {
						sb.WriteString(fmt.Sprintf("-%s\n", line))
					}
				}
				if op.NewText != "" {
					newLines := strings.Split(strings.TrimSuffix(op.NewText, "\n"), "\n")
					for _, line := range newLines {
						sb.WriteString(fmt.Sprintf("+%s\n", line))
					}
				}

			case Insert:
				sb.WriteString(fmt.Sprintf("@@ line %d @@\n", op.StartLine))
				if op.NewText != "" {
					newLines := strings.Split(strings.TrimSuffix(op.NewText, "\n"), "\n")
					for _, line := range newLines {
						sb.WriteString(fmt.Sprintf("+%s\n", line))
					}
				}

			case Delete:
				sb.WriteString(fmt.Sprintf("@@ line %d @@\n", op.StartLine))
				if op.OldText != "" {
					oldLines := strings.Split(strings.TrimSuffix(op.OldText, "\n"), "\n")
					for _, line := range oldLines {
						sb.WriteString(fmt.Sprintf("-%s\n", line))
					}
				}
			}
		}
	}

	return sb.String()
}

package editing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestToUnifiedDiffGolden(t *testing.T) {
	patch := NewPatch(
		Operation{
			Type:      Replace,
			File:      "server.go",
			StartLine: 42,
			EndLine:   42,
			OldText:   "Provider",
			NewText:   "Model",
		},
		Operation{
			Type:      Insert,
			File:      "server.go",
			StartLine: 60,
			EndLine:   62,
			NewText:   "func health() bool {\n\treturn true\n}",
		},
		Operation{
			Type:      Delete,
			File:      "config.go",
			StartLine: 15,
			EndLine:   15,
			OldText:   "timeout=10",
		},
	)

	diff := patch.ToUnifiedDiff()

	expectedBytes, err := os.ReadFile(filepath.Join("testdata", "expected.diff"))
	if err != nil {
		t.Fatalf("failed to read golden file: %v", err)
	}

	expected := strings.TrimSpace(string(expectedBytes))
	actual := strings.TrimSpace(diff)

	if actual != expected {
		t.Errorf("ToUnifiedDiff output does not match golden file:\n\n"+
			"=== Expected ===\n%s"+
			"\n\n=== Actual ===\n%s",
			expected, actual)
	}
}

func TestToUnifiedDiffEmpty(t *testing.T) {
	var p *Patch
	if diff := p.ToUnifiedDiff(); diff != "" {
		t.Errorf("expected empty string for nil patch, got %q", diff)
	}

	p = NewPatch()
	if diff := p.ToUnifiedDiff(); diff != "" {
		t.Errorf("expected empty string for empty patch, got %q", diff)
	}
}

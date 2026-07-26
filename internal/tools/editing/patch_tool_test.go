package editing

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Nithwin/WindMist/internal/tools"
)

func TestPatchTool(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "windmist_patch_tool_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(filePath, []byte("line 1\nline 2\nline 3\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Change working directory for patch to work correctly
	oldCwd, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldCwd)

	patchStr := `--- test.txt
+++ test.txt
@@ -1,3 +1,3 @@
 line 1
-line 2
+line 2 edited
 line 3
`

	tool := NewPatchTool()
	res := tool.Run(context.Background(), tools.Call{
		Args: map[string]interface{}{
			"diff": patchStr,
		},
	})

	if res.Error != nil {
		t.Fatalf("unexpected error: %v", res.Error)
	}

	content, err := os.ReadFile("test.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(content), "line 2 edited") {
		t.Fatalf("patch was not applied correctly, got content: %s", string(content))
	}
}

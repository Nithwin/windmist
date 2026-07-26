import re

with open('internal/chat/commands.go', 'r') as f:
    content = f.read()

# The file contains the Registry declaration up to line ~150, then specific funcs.
# Let's just find the functions and move them.

def extract_func(name):
    global content
    pattern = rf"func {name}\(.*?\) .*?{{.*?^}}"
    match = re.search(pattern, content, re.MULTILINE | re.DOTALL)
    if match:
        func_body = match.group(0)
        # Remove from content
        content = content.replace(func_body, "")
        return func_body
    return ""

commands_ai = ["selectProviderCmd", "selectModelCmd", "selectModeCmd", "selectSubagentCmd", "setAPIKeyCmd"]
commands_mcp = ["selectMCPCmd", "mcpEnvPromptChain"]
commands_ui = ["selectThemeCmd"]
commands_session = ["selectSessionCmd"]

def write_file(filename, funcs, imports):
    body = "\n\n".join([extract_func(f) for f in funcs])
    if body.strip():
        with open(filename, 'w') as f:
            f.write(f"package chat\n\nimport (\n{imports}\n)\n\n{body}\n")

write_file('internal/chat/commands_ai.go', commands_ai, '\t"fmt"\n\t"strings"\n\n\t"github.com/Nithwin/WindMist/internal/config"\n\t"github.com/Nithwin/WindMist/internal/ui/selector"\n\ttea "github.com/charmbracelet/bubbletea"')
write_file('internal/chat/commands_mcp.go', commands_mcp, '\t"fmt"\n\t"strings"\n\t"strconv"\n\n\t"github.com/Nithwin/WindMist/internal/mcp"\n\t"github.com/Nithwin/WindMist/internal/ui/selector"\n\ttea "github.com/charmbracelet/bubbletea"')
write_file('internal/chat/commands_ui.go', commands_ui, '\t"github.com/Nithwin/WindMist/internal/ui"\n\t"github.com/Nithwin/WindMist/internal/ui/selector"\n\ttea "github.com/charmbracelet/bubbletea"')
write_file('internal/chat/commands_session.go', commands_session, '\t"fmt"\n\t"os"\n\t"strings"\n\t"time"\n\n\t"github.com/Nithwin/WindMist/internal/ui/selector"\n\ttea "github.com/charmbracelet/bubbletea"')

with open('internal/chat/commands.go', 'w') as f:
    f.write(content)


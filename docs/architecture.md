# WindMist Architecture & Engineering Specification (`v1.0.0`)

> **Design Philosophy:** WindMist adheres strictly to **Clean Architecture with Domain-Driven Modules**. Every capability is isolated behind explicit Go interfaces (`ai.Provider`, `tools.Tool`), making the platform resilient, testable, and maintainable without architectural bloat or circular dependency cycles.

---

## 1. System Overview

WindMist is an autonomous AI software engineering agent running directly inside your terminal. It is built 100% in high-performance **Go**, operating locally on your filesystem to inspect code, edit exact ranges or full files, execute tools, and engage in multi-turn reasoning loops.

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│                               USER TERMINAL (CLI)                                │
└────────────────────────────────────────┬─────────────────────────────────────────┘
                                         │
                                         ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                             GO ENGINE (internal/)                                │
│                                                                                  │
│  ┌──────────────────┐    ┌──────────────────┐    ┌────────────────────────────┐  │
│  │     Cobra CLI    │    │  Bubble Tea TUI  │    │       Config / Viper       │  │
│  │ (cmd/windmist/*) │    │ (internal/chat/*)│    │       (~/.windmist/)       │  │
│  └────────┬─────────┘    └────────┬─────────┘    └─────────────┬──────────────┘  │
│           │                       │                            │                 │
│           └───────────────────────┼────────────────────────────┘                 │
│                                   ▼                                              │
│                  ┌──────────────────────────────────┐                            │
│                  │           Agent Loop             │                            │
│                  │      (internal/agent/*)          │                            │
│                  └────────────────┬─────────────────┘                            │
│                                   │                                              │
│         ┌─────────────────────────┴─────────────────────────┐                    │
│         ▼                                                   ▼                    │
│  ┌──────────────────────────────────┐       ┌───────────────────────────────┐    │
│  │        Atomic Tool Engine        │       │       Provider Router         │    │
│  │     (internal/tools/defaults)    │       │     (internal/providers/*)    │    │
│  └────────────────┬─────────────────┘       └───────────────┬───────────────┘    │
│                   │                                         │                    │
│         ┌─────────┴─────────┐                               │                    │
│         ▼                   ▼                               ▼                    │
│  ┌─────────────┐     ┌─────────────┐               ┌─────────────────┐           │
│  │ Filesystem  │     │   Editing   │               │ Native Gemini   │           │
│  │ (read/write)│     │(replace/rng)│               │  Tool Calling   │           │
│  └─────────────┘     └─────────────┘               └─────────────────┘           │
└──────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Domain-Driven Modules (`internal/`)

Go uses `internal/` to guarantee that core engine mechanisms cannot be imported directly by external applications. Every subdirectory under `internal/` represents a self-contained domain module with explicit interfaces.

### `internal/agent`
The stateless autonomous reasoning engine.
- **`Agent` & `runLoop`:** Orchestrates multi-turn interactions between the user, the LLM provider, and tool executions.
- **`toolDefinitions()`:** Translates registered `tools.Tool` items into provider-agnostic schemas (`ai.ToolDefinition`).
- **Turn Tracking & Usage Metrics:** Accumulates prompt and candidate tokens (`ai.Usage`) across turns while enforcing `agent.Config.MaxTurns` (`ErrMaxTurnsExceeded`).

### `internal/tools` & Subpackages
The toolbox that allows the agent to inspect and mutate the workspace. Every tool implements the `tools.Tool` interface:
```go
type Tool interface {
    Definition() Definition
    Run(ctx context.Context, call Call) Result
}
```
- **`internal/tools/filesystem`:** Atomic filesystem operations (`read`, `write`, `append`, `delete`, `rename`, `list`, `create`, `info`, `exists`).
- **`internal/tools/editing`:** Precision code editing and context inspection (`replace_text`, `replace_range`, `delete_range`, `read_context`, `insert_text`, `search_text`).
- **`internal/tools/defaults`:** Helper package (`RegisterAll`) that registers all 15 built-in tools cleanly without circular dependency cycles.

### `internal/providers/gemini`
The native LLM gateway.
- Implements `ai.Provider` (`Generate`, `Stream`).
- **Schema Translation (`translateTools`):** Maps Go `ai.ToolParameter` definitions directly into Gemini OpenAPI `OBJECT` schemas with typed property maps (`STRING`, `INTEGER`, `BOOLEAN`, `ARRAY`, `OBJECT`).
- **Turn & Function Mapping (`translateMessages`):** Converts multi-turn `ai.Message` history into Gemini `Content`/`Part` arrays, separating user prompts, model `FunctionCall` requests, and `FunctionResponse` tool execution outputs.

### `internal/chat` & `internal/ui`
The rich interactive terminal user interface (`windmist`).
- Powered by **Bubble Tea** and **Lip Gloss** (`model.go`, `view.go`, `update.go`).
- Renders GitHub-flavored Markdown tables and code blocks dynamically via `ui.MarkdownRenderer` and `glamour`.
- Streams multi-turn agent responses cleanly to responsive terminal viewport bubbles.

### `cmd/windmist`
The command-line entrypoint (`main.go` & `cmd/`).
- **`windmist chat <prompt>`:** Runs a synchronous, multi-turn `Agent.Run(ctx, prompt)` execution and prints final resolved output.
- **`windmist` (no args):** Launches the interactive Bubble Tea TUI.
- **`windmist set <key> <value>` / `windmist show`:** Manages local configuration via `viper` (`~/.windmist/config.yaml`).
- **`windmist version`:** Displays current semantic release version (`v1.0.0`).

---

## 3. Autonomous Agent Execution Loop

When a prompt is dispatched to `Agent.Run(ctx, prompt)`, the agent runs our stateless iterative loop (`runLoop` inside `loop.go`):

```text
[User Prompt] -> [appendUser] -> [Provider.Generate with Tools]
                                            │
                                            ▼
                                  ┌──────────────────┐
                                  │ Model Candidate  │
                                  └─────────┬────────┘
                                            │
                           ┌────────────────┴────────────────┐
                           ▼                                 ▼
                 [No ToolCalls Found]               [ToolCalls Present]
                           │                                 │
                           ▼                                 ▼
                     Return Result               For each Call in ToolCalls:
                    (Done / Success)             1. appendAssistant(Call)
                                                 2. executor.execute(Call)
                                                 3. appendToolResults(Output)
                                                             │
                                                             └───────► Loop to Turn + 1
```

1. **State Isolation:** A fresh `[]ai.Message` slice is initialized for each `Run` invocation, preventing state leakage across concurrent requests or successive TUI turns.
2. **Execution & Dispatch:** When the model returns one or more `ToolCalls`, `executor.execute()` looks up the tool in `tools.Manager` (`filesystem` or `editing`), runs `Run(ctx, call)`, and captures the output or error.
3. **Structured Recovery:** The tool outputs (`ai.ToolResult`) are appended to the conversation history (`appendToolResults`), and the model is invoked again with updated context so it can verify its edits or self-correct any errors until it completes the task or reaches `MaxTurns`.

---

## 4. Current Directory Structure (`v1.0.0`)

WindMist strictly adheres to Go's `internal/` package boundary conventions:

```text
windmist/
├── cmd/
│   └── windmist/             # Application entry point (main.go)
├── internal/                 # Private domain-driven Go engine modules
│   ├── agent/                # Stateless multi-turn reasoning and loop execution
│   ├── chat/                 # Bubble Tea terminal controllers and viewports
│   ├── config/               # Configuration loading (~/.windmist/config.yaml)
│   ├── providers/
│   │   ├── ai/               # Provider-agnostic interfaces (Provider, Message, ToolDefinition)
│   │   └── gemini/           # Native Gemini v1beta schema and tool calling adapter
│   ├── tools/
│   │   ├── defaults/         # Tool registry (RegisterAll)
│   │   ├── editing/          # Precision code editing tools (replace, insert, search)
│   │   └── filesystem/       # Atomic filesystem operations (read, write, list)
│   └── ui/                   # Lip Gloss and Glamour markdown styling helpers
├── docs/                     # Architectural deep-dives and developer guides
├── README.md                 # Project overview & installation
├── CONTRIBUTING.md           # Workflow & coding standards
├── SECURITY.md               # Threat models & vulnerability reporting
└── ROADMAP.md                # Development phases
```

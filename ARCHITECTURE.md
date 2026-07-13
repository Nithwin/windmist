# WindMist Architecture & Engineering Specification

> **Design Philosophy:** WindMist adheres strictly to **Clean Architecture with Domain-Driven Modules**, rather than traditional layered MVC. Every capability is isolated behind strict interfaces, making the platform resilient, testable, and infinitely extensible without fear of architectural drift.

---

## 1. System Overview

WindMist is a hybrid CLI-first autonomous AI software engineering platform. It is engineered to operate directly on the user's local filesystem and git repository, leveraging remote or local AI models to reason, plan, edit code, execute tests, and recover from runtime failures.

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
│  │   CLI / UI       │    │     Commands     │    │       Config / Viper       │  │
│  │  (Cobra/Bubble)  │    │  (chat/build/fix)│    │       (~/.windmist/)       │  │
│  └────────┬─────────┘    └────────┬─────────┘    └─────────────┬──────────────┘  │
│           │                       │                            │                 │
│           └───────────────────────┼────────────────────────────┘                 │
│                                   ▼                                              │
│                  ┌──────────────────────────────────┐                            │
│                  │           Agent Engine           │                            │
│                  │  (Planner, Loop, Recovery State) │                            │
│                  └────────────────┬─────────────────┘                            │
│                                   │                                              │
│         ┌─────────────────────────┼─────────────────────────┐                    │
│         ▼                         ▼                         ▼                    │
│  ┌─────────────┐          ┌───────────────┐         ┌──────────────┐             │
│  │ Repository  │          │ Tool Executor │         │  Memory Store│             │
│  │ (AST / Git) │          │ (File/Shell)  │         │   (SQLite)   │             │
│  └─────────────┘          └───────┬───────┘         └──────────────┘             │
│                                   │                                              │
│                                   ▼                                              │
│                         ┌───────────────────┐                                    │
│                         │  Provider Router  │                                    │
│                         │ (Gemini/OpenAI/etc│                                    │
│                         └─────────┬─────────┘                                    │
└───────────────────────────────────┼──────────────────────────────────────────────┘
                                    │ HTTP / gRPC / JSON
                                    ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         PYTHON AI SERVICE (python/)                              │
│                                                                                  │
│  ┌──────────────────┐    ┌──────────────────┐    ┌────────────────────────────┐  │
│  │    Embeddings    │    │   RAG Engine     │    │   Vector Store / Eval      │  │
│  │  (SentenceTrans) │    │  (Semantic Ret.) │    │       (ChromaDB)           │  │
│  └──────────────────┘    └──────────────────┘    └────────────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Decoupled Language Strategy: Go vs. Python

To compete directly with enterprise-grade tools like *Claude Code* and *Codex CLI*, WindMist divides labor based on strict performance and ecosystem realities.

### Go: High-Performance Operating System Layer
Go `1.25+` owns everything that requires **low latency, native concurrency, and direct OS interaction**:
- **Instantaneous Startup:** Zero bytecode overhead or Python virtual environment latency.
- **Concurrent File & Symbol Walking:** High-throughput traversal of massive codebases using Goroutines and worker pools (`internal/repository`).
- **Git Integration:** Direct local repository mutation, diff generation, and commit atomic safety (`internal/git`).
- **Tool Sandbox Execution:** Spawning, monitoring, and capturing `stdout`/`stderr` from terminal commands safely (`internal/executor`).
- **Memory & Session State:** Fast SQLite transactions for prompt context and user history (`internal/memory`).

### Python: Machine Learning & Intelligence Layer
Python `3.13+` owns **AI-heavy computational tasks** where Go lacks native ecosystem depth:
- **Embeddings & Transformers:** Generating vector representations for code snippets using `sentence-transformers` (`python/embeddings`).
- **Vector Search & RAG:** Local index management via `ChromaDB` or `Qdrant` for semantic code retrieval (`python/rag`).
- **Evaluation & Benchmarking:** Automated scoring of agent completions against standard code benchmarks (`python/evaluation`).

### Inter-Process Communication (IPC)
The Go engine communicates with the Python service exclusively via a localized **FastAPI HTTP/JSON server or gRPC**. Go never imports Python C-bindings. If the user only runs a basic `chat` command without vector search, the Python service does not even need to be running, preserving system memory.

---

## 3. Domain-Driven Modules (`internal/`)

Go uses `internal/` to guarantee that core engine mechanisms cannot be imported directly by external Go repositories. Every subdirectory under `internal/` represents a self-contained domain module with explicit interfaces.

### `internal/agent`
The central autonomous reasoning brain.
- **`Planner` Interface:** Decomposes complex user requests (e.g., *"Refactor the authentication middleware to use JWT"*) into directed acyclic dependency graphs (DAGs) of execution steps.
- **`AgentLoop`:** Executes steps, monitors tool outputs, manages token budgets, and invokes self-correction workflows if a sub-step fails.

### `internal/tools`
The toolbox that allows the agent to interact with the world. Every tool implements the `Tool` interface:
```go
type Tool interface {
    Name() string
    Description() string
    InputSchema() json.RawMessage
    Execute(ctx context.Context, input json.RawMessage) (*ToolOutput, error)
}
```
- Built-in Tools: `ReadFileTool`, `WriteFileTool`, `SearchTool`, `ShellTool`, `GitTool`, `TestTool`, `LintTool`.

### `internal/providers`
The unified LLM gateway.
- Isolates vendor-specific SDKs behind a singular `Provider` streaming interface:
```go
type Provider interface {
    CompleteStream(ctx context.Context, req *CompletionRequest) (<-chan *CompletionChunk, error)
    Models() []ModelInfo
}
```
- Supported implementations: **Gemini**, **OpenAI**, **Anthropic**, **Groq**, **Ollama**, **Azure OpenAI**.

### `internal/repository`
Code awareness and symbol tracking.
- Parses Go, Python, TypeScript, and Rust ASTs to extract structural symbols, class definitions, and function signatures without needing full compiler server (LSP) overhead.
- Respects `.gitignore` and `.windmistignore` patterns natively.

### `internal/memory` & Storage Schema
Local persistence powered by **SQLite** (`~/.windmist/history.db` & `~/.windmist/memory.db`).
- **Sessions Table:** Tracks multi-turn conversation histories, user prompts, and agent reasoning chains.
- **Workspace Cache Table:** Stores file hashes, modification timestamps, and indexing metadata to prevent redundant embedding calculations.

---

## 4. Autonomous Agent Execution Loop

When a user triggers an actionable task via `windmist build` or `windmist fix`, the Agent follows our standardized **Reason -> Plan -> Execute -> Verify** lifecycle:

```text
[User Command] -> [Agent Context Builder]
                         │
                         ▼
             ┌───────────────────────┐
             │  Generate Plan (DAG)  │
             └───────────┬───────────┘
                         │
                         ▼
        ┌─────────► [Step Selection] ◄────────┐
        │                │                    │
        │                ▼                    │
        │       [Select Appropriate Tool]     │
        │                │                    │
        │                ▼                    │
        │       [Execute Sandbox Tool]        │
        │                │                    │
        │                ▼                    │
 [Failure/Error]   [Verify Result Output]     │
        │                │                    │
        └────────────────┼────────────────────┘
                         │ (Success)
                         ▼
              [Final Verification & Commit]
```

1. **Context Construction:** Retrieves workspace symbols from `internal/repository` and relevant semantic snippets from `python/rag`.
2. **Task Planning:** The `Planner` outputs a strict, sequential step plan.
3. **Tool Dispatch:** For each step, the Agent constructs tool arguments and requests confirmation from the user if the action is deemed potentially destructive (e.g., `ShellTool` executing `rm -rf` or `git reset`).
4. **Self-Correction & Recovery:** If `TestTool` or `LintTool` returns a non-zero exit code or error output, the Agent captures `stderr`, feeds it back into the prompt context, and initiates a recovery retry loop (up to `max_retries`).
5. **Atomic Finalization:** Once all verification tests pass, `GitTool` stages mutated files and generates a semantic commit summary.

---

## 5. Extensibility & Future Scaling

By locking our architecture into Domain-Driven Modules today, WindMist is positioned for seamless enterprise expansion:
- **Model Context Protocol (MCP):** Our tool framework directly maps to MCP definitions, allowing external third-party servers to inject custom tools dynamically into the agent loop.
- **Go SDK (`sdk/`):** Allows IDE plugin developers (VS Code, Neovim, JetBrains) to embed WindMist's agent engine directly inside custom workflows without relying on terminal CLI wrappers.

# WindMist ── AI Software Engineer in Your Terminal

<div align="center">

![WindMist Banner](https://img.shields.io/badge/WindMist-AI%20Software%20Engineer-6366f1?style=for-the-badge&logo=go&logoColor=white)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Python Version](https://img.shields.io/badge/Python-3.13+-3776AB?style=for-the-badge&logo=python&logoColor=white)](https://python.org)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=for-the-badge)](CONTRIBUTING.md)

**Not just another ChatGPT wrapper.**  
WindMist is an autonomous, lightning-fast AI Software Engineer running directly inside your terminal, engineered from the ground up for speed, deep repository awareness, and multi-step reasoning.

[Key Features](#-key-features) • [Why WindMist?](#-why-windmist) • [Architecture](#-high-level-architecture) • [Getting Started](#-getting-started) • [Roadmap](ROADMAP.md) • [Contributing](CONTRIBUTING.md)

</div>

---

## ⚡ The Vision

When you type:

```bash
windmist build authentication
```

WindMist doesn't just output a snippet of code and leave you to copy-paste it. **It acts as a true engineering partner:**

1. **Understands your repository** by indexing local code structures, ASTs, symbols, and dependencies.
2. **Plans the implementation** step-by-step through an autonomous planning loop.
3. **Edits multiple files** safely across your workspace.
4. **Runs tests and linters** to verify its own work.
5. **Fixes errors automatically** if tests fail or compiler bugs arise.
6. **Commits atomic changes** with clear, descriptive Git commit messages.
7. **Asks for confirmation** whenever sensitive or high-impact actions are proposed.

Our mission is to build the premier open-source terminal AI engineer to compete directly with tools like *Claude Code*, *Gemini CLI*, and *OpenAI Codex CLI* — while giving developers complete sovereignty over their models and infrastructure.

---

## 🚀 Key Features

- **🧠 Multi-Model & Provider Agnostic:** Plug in your choice of **Gemini 2.5 / 3.0**, **OpenAI (GPT-4o)**, **Anthropic (Claude 3.5/3.7 Sonnet)**, **Groq**, **Azure OpenAI**, or completely local models via **Ollama**.
- **⚡ Hybrid Go + Python Engine:** The best of both worlds — Go handles blazing-fast CLI interactions, concurrent file walking, Git operations, and tool loops; a decoupled Python microservice handles deep RAG, embeddings, vector indexing, and local evaluation.
- **🏗️ Clean Architecture with Domain-Driven Modules:** Designed for longevity, extensibility, and maintainability. Every capability is encapsulated behind strict interfaces.
- **🛠️ Autonomous Tool Ecosystem:** Built-in specialized tools for file reading/writing, shell command execution, semantic search, Git workflows, web browsing, testing, and linting.
- **💾 Local First & Privacy Focused:** Uses localized **SQLite** (`~/.windmist/`) for persistent session history, agent state, and workspace memory without shipping your private code to third-party databases.

---

## 📐 High-Level Architecture

```text
                                  User
                                   │
                                   ▼
                            WindMist CLI (Go)
                                   │
                     ┌─────────────┼──────────────┐
                     │             │              │
                     ▼             ▼              ▼
                 Commands      AI Engine      Config
                     │             │              │
                     └──────┬──────┴──────┬───────┘
                            │
                            ▼
                      Repository Engine
                            │
                     ┌──────┼────────┐
                     ▼      ▼        ▼
              File System Git    Memory
                            │
                            ▼
                      Tool Execution
                            │
                            ▼
                      AI Providers
            (Gemini / OpenAI / Anthropic / Groq / Ollama)
                            │
                            ▼ [Decoupled via HTTP/JSON]
                      Python AI Service
                 (Embeddings / RAG / Evaluation)
```

---

## 💡 Why Go + Python?

Most AI tools force a compromise: pure Python tools suffer from slow startup times, heavy dependency footprints, and sluggish terminal UX, while pure Go tools struggle to tap into the cutting-edge ecosystem of machine learning and vector indexing libraries. 

**WindMist eliminates this compromise through strict separation of responsibilities:**

### Go Owns Core Performance & UX
Everything that must be instantaneous, concurrent, and close to the operating system lives in Go (`internal/`):
- CLI parsing (**Cobra**) & Rich Terminal UI (**Bubble Tea + Lip Gloss**)
- Concurrent file walking, symbol extraction, and AST inspection
- Autonomous agent loop, prompt building, and tool execution
- Git integration and local workspace mutation
- Multi-provider HTTP/streaming integration

### Python Owns Deep AI Workloads
Complex machine learning workloads are isolated in a lightweight, decoupled **FastAPI** service (`python/`):
- Local & remote text embeddings (**Sentence Transformers**)
- Advanced RAG & semantic code retrieval
- Local vector storage indexing (**ChromaDB / Qdrant**)
- Future model fine-tuning, training, and automated evaluation pipelines

### Clean Decoupling
Go and Python communicate exclusively over clean HTTP/gRPC interfaces or structured IPC. This ensures that WindMist starts in milliseconds, works reliably without heavy Python environments for core tasks, and scales cleanly when advanced AI features are engaged.

---

## 🛠️ Core Commands (Planned)

| Command | Description |
| :--- | :--- |
| `windmist chat` | Start an interactive, context-aware AI pairing session in your terminal. |
| `windmist build <prompt>` | Give high-level instructions and let WindMist autonomously plan, code, and test. |
| `windmist review [PR/Branch]` | Perform a comprehensive architectural and code quality review of recent changes. |
| `windmist fix` | Automatically diagnose and resolve failing compiler errors or unit tests. |
| `windmist doctor` | Verify system health, API provider connectivity, SQLite state, and local dependencies. |
| `windmist auth` | Securely configure and store API keys for AI providers. |

---

## 🗂️ Project Directory Structure (Planned)

We strictly adhere to Go's `internal/` package conventions to guarantee encapsulation and maintain a clean boundary between private engine mechanics and external SDKs.

```text
windmist/
├── cmd/
│   └── windmist/          # Main application entry point (main.go)
├── internal/              # Private domain-driven modules
│   ├── cli/               # Cobra commands and interactive terminal controllers
│   ├── commands/          # Subcommand handlers (chat, review, fix, etc.)
│   ├── config/            # Viper configuration and environment management
│   ├── agent/             # Core autonomous reasoning and state recovery loop
│   ├── planner/           # Multi-step task decomposition and graph planning
│   ├── tools/             # Atomic tool implementations (File, Git, Shell, Search)
│   ├── providers/         # Unified LLM interfaces (Gemini, OpenAI, Anthropic, Groq)
│   ├── git/               # Git repository inspection and mutation utilities
│   ├── repository/        # AST parsing, file walking, and symbol indexing
│   ├── memory/            # SQLite storage interface for sessions and long-term memory
│   ├── session/           # Active workspace state and multi-turn context manager
│   ├── executor/          # Safe shell execution and sandboxing controls
│   ├── parser/            # Code parsing and response formatting tools
│   ├── prompt/            # System instructions and dynamic prompt builders
│   ├── streaming/         # Server-Sent Events (SSE) and terminal streaming helpers
│   ├── telemetry/         # Optional, privacy-first diagnostic metrics
│   └── util/              # Shared concurrent primitives and helpers
├── sdk/                   # Public Go SDK for extending WindMist programmatically
├── plugins/               # External tool and Model Context Protocol (MCP) integrations
├── api/                   # Interface definitions and OpenAPI/proto specs
├── python/                # Decoupled AI microservice
│   ├── embeddings/        # Sentence transformer embedding generators
│   ├── rag/               # Retrieval-Augmented Generation retrieval engines
│   ├── evaluation/        # Benchmark and quality evaluation suites
│   └── server/            # FastAPI server entry point
├── examples/              # Sample workflows and configuration templates
├── docs/                  # Architectural deep-dives and user guides
├── scripts/               # Build, installation, and setup automation scripts
├── tests/                 # End-to-end integration and system test suites
├── .github/               # CI/CD workflows and community templates
├── Makefile               # Universal build and development automation
├── go.mod                 # Go module definition
└── README.md              # Project overview
```

---

## 🗺️ Roadmap & Phases

We are currently in **Phase 1: Foundation & Planning**. Check out [ROADMAP.md](ROADMAP.md) for full phase breakdowns:

- [x] **Phase 0:** Product Vision, Architecture Design & Planning (`We are here`)
- [ ] **Phase 1:** Foundation (CLI setup, Config, Provider Interfaces, Logging)
- [ ] **Phase 2:** Interactive Chat (Terminal UI, Streaming, Markdown/Syntax Rendering)
- [ ] **Phase 3:** Repository Awareness (AST/File Indexing, Git Integration, Ignore Patterns)
- [ ] **Phase 4:** Autonomous Agent (Planner, Tool Execution Loop, Multi-Step Reasoning)
- [ ] **Phase 5:** Advanced AI Service (Python RAG, Embeddings, Semantic Search, SQLite Memory)
- [ ] **Phase 6:** Production & Ecosystem (MCP Support, Plugin Engine, Telemetry, Cross-Platform Packaging)

---

## 🤝 Contributing & Code of Conduct

We are designing WindMist like a high-growth startup engineering team: high standards, clean architecture, thorough documentation, and rigorous testing.

- **Contributing Guide:** Please read our [CONTRIBUTING.md](CONTRIBUTING.md) to understand our architectural principles, Go/Python integration guidelines, and pull request workflow.
- **Code of Conduct:** We are committed to fostering an inclusive and professional community. Read our [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

---

## 📄 License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <sub>Built with ❤️ for the next generation of autonomous engineering in the terminal.</sub>
</div>

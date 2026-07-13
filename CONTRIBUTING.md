# Contributing to WindMist

First off, thank you for considering contributing to **WindMist**! 🎉

We are building an autonomous, high-performance AI Software Engineer right inside the terminal. We treat this project with the rigor and design standards of a top-tier startup product. Whether you are fixing bugs, improving documentation, adding new AI providers, or designing new tools, your contributions are invaluable.

This guide details our engineering philosophy, architecture rules, and workflow so you can jump right in.

---

## 🏛️ Core Architectural Principles

Before writing any code, it is critical to understand our fundamental design decision:

> **We follow Clean Architecture with Domain-Driven Modules, NOT a traditional layered MVC approach.**

Every feature must encapsulate its own logic, state, and dependencies behind clean, well-defined interfaces. This ensures WindMist can scale to support dozens of AI providers, custom tools, and complex plugins without spaghetti dependencies or architectural rot.

### 1. Separation of Responsibilities: Go vs. Python

We strictly separate our **Go** engine (`internal/`) from our **Python** AI microservice (`python/`):

| Layer | Language | Responsibilities | Why? |
| :--- | :--- | :--- | :--- |
| **Core Engine (`internal/`)** | **Go 1.25+** | CLI commands, terminal UI, concurrent file walking, AST parsing, Git operations, tool execution loop, agent reasoning loop, SQLite storage, and provider HTTP routing. | Go provides instant startup times, high concurrency, low memory overhead, and native cross-platform binaries without virtual environments. |
| **AI Service (`python/`)** | **Python 3.13+** | Local/remote embeddings (`Sentence Transformers`), advanced RAG pipelines, vector indexing (`ChromaDB`), evaluation benchmarks, and future model fine-tuning. | Python possesses an unmatched machine learning ecosystem that Go cannot replicate natively. |

#### ⚠️ Communication Rule
Go and Python must **never** tightly couple or share memory space. They communicate cleanly over decoupled HTTP (`FastAPI`) or gRPC endpoints. Go only starts or calls the Python service when heavy AI/RAG tasks are requested.

### 2. The `internal/` Boundary
In Go, packages inside `internal/` cannot be imported by external applications.
- **Always put WindMist engine modules inside `internal/`** (`internal/agent`, `internal/tools`, `internal/providers`, etc.).
- Only put public code meant for external third-party integrations inside `sdk/` or `api/`.

### 3. Interface-Driven Providers & Tools
When adding a new AI provider (e.g., Azure, Ollama) or a new tool (e.g., `LintTool`, `DockerTool`):
- Implement our common `Provider` interface inside `internal/providers/`.
- Implement our common `Tool` interface inside `internal/tools/`.
- **Never** hardcode vendor-specific logic into the `internal/agent/` reasoning loop.

---

## 🛠️ Development Environment Setup

### Prerequisites
- **Go:** `1.25` or higher
- **Python:** `3.13` or higher
- **Git:** `2.40+`
- **SQLite:** `3.40+` (included via CGO or modern pure-Go drivers)

### Initial Setup (Planning Stage Preview)
While we are currently in Phase 1 setup, our standard workspace initialization will follow these steps once code scaffolding begins:

```bash
# Clone the repository
git clone https://github.com/your-username/windmist.git
cd windmist

# Verify Go setup
go version

# Verify Python setup
python3 --version
```

---

## 🌿 Branching & Workflow

We follow a clean feature-branching workflow:

1. **Fork the repository** on GitHub.
2. **Create a topic branch** from `main`:
   ```bash
   git checkout -b feat/add-groq-provider
   # or
   git checkout -b fix/ast-parser-memory-leak
   ```
3. **Write clean, documented code** adhering to Clean Architecture principles.
4. **Run tests & linters** locally before committing.

### Branch Naming Conventions
- `feat/`<brief-description> (New features or providers)
- `fix/`<brief-description> (Bug fixes)
- `docs/`<brief-description> (Documentation improvements)
- `refactor/`<brief-description> (Code restructuring without behavior changes)
- `test/`<brief-description> (Adding or updating test suites)

---

## 💬 Conventional Commits

Please format your commit messages using the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```text
<type>(<scope>): <short summary>
```

**Examples:**
- `feat(providers): implement Groq streaming provider interface`
- `fix(tools): resolve path traversal bug in WriteFileTool`
- `docs(readme): add detailed clean architecture breakdown`
- `refactor(agent): separate tool execution loop from planner graph`

---

## 🧪 Testing Guidelines

We value reliability above all else. Every new tool, command, or provider must include comprehensive automated tests.

### Go Tests
All domain modules must have corresponding `_test.go` files:
```bash
# Run unit tests across all internal packages
go test -v ./internal/...

# Run tests with race condition detection
go test -race ./internal/...
```

### Python Tests
Python microservice modules must be tested via `pytest`:
```bash
cd python
python3 -m pytest tests/ -v
```

---

## 🔍 Code Review Expectations

When you submit a Pull Request:
1. Ensure your PR description clearly explains **what** the change does, **why** it was implemented that way, and **how** it was tested.
2. Link any related GitHub issues.
3. Keep PRs atomic and focused on a single responsibility or domain module.
4. Be responsive to review feedback — we aim to review all PRs constructively and promptly!

---

## ❓ Questions or Ideas?

If you have an idea for a major architectural change or a new core module, please open a GitHub Discussion or Issue first to discuss it with the maintainers before writing extensive code.

Thank you for helping shape **WindMist** into the ultimate terminal AI engineer! 🚀

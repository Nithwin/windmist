# Contributing to WindMist

First off, thank you for considering contributing to **WindMist**! 🎉

We are building a modern open-source AI coding assistant right inside the terminal. We treat this project with the rigor and design standards of a top-tier product. Whether you are fixing bugs, improving documentation, adding new AI providers, or designing new tools, your contributions are invaluable.

This guide details our workflow, coding standards, and high-level structure so you can jump right in.

---

## 🚀 Quick Start

To get up and running quickly with WindMist development:

1. **Clone your fork** of the repository:
   ```bash
   git clone https://github.com/your-username/windmist.git
   cd windmist
   ```
2. **Review prerequisites and installation steps** in our [`README.md`](README.md).
3. **Run existing tests** to verify your local setup before making modifications:
   ```bash
   go test ./...
   ```

---

## 🛠️ Development Setup

For full details on local environment prerequisites, dependency management, and configuration, please refer to:
* **[`README.md`](README.md)** for general system requirements and installation.
* **`docs/development.md`** for detailed local setup instructions, environment variables, and debugging guides.

Contributing shouldn't duplicate setup documentation—always keep your local environment synchronized with our core setup guides.

---

## 🌿 Workflow

We follow a clean feature-branching workflow:

1. **Fork the repository** on GitHub.
2. **Create a topic branch** from `main`:
   ```bash
   git checkout -b feat/add-groq-provider
   # or
   git checkout -b fix/ast-parser-memory-leak
   ```
3. **Write clean, documented code** adhering to our engineering standards.
4. **Run tests & linters** locally before committing.

### Branch Naming Conventions
* `feat/`<brief-description> (New features or providers)
* `fix/`<brief-description> (Bug fixes)
* `docs/`<brief-description> (Documentation improvements)
* `refactor/`<brief-description> (Code restructuring without behavior changes)
* `test/`<brief-description> (Adding or updating test suites)

---

## 💬 Coding Standards

We prioritize code clarity, reliability, and long-term maintainability. When writing code for WindMist:
* Ensure all exported functions, types, and interfaces have clear docstrings.
* Write unit tests alongside your implementation.
* Avoid hardcoding vendor-specific logic into core loops; rely on clean interfaces.

### Conventional Commits
Please format your commit messages using the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```text
<type>(<scope>): <short summary>
```

**Examples:**
* `feat(providers): implement Groq streaming provider interface`
* `fix(tools): resolve path traversal bug in write tool`
* `docs(readme): add detailed architecture breakdown`
* `refactor(agent): separate tool execution loop from planner`

---

## 🏛️ Architecture Overview

WindMist is structured around stable concepts designed for modularity and high performance:

| Component | Responsibility |
| :--- | :--- |
| **Agent** | Stateless multi-turn reasoning loop, token usage tracking, and tool execution orchestration. |
| **Tools** | Atomic filesystem and precision text-editing operations. |
| **Providers** | Interface-driven AI model integrations and schema translations. |
| **UI** | Terminal user interface, streaming display, and interactive components. |

Please read **`docs/architecture.md`** before making major architectural changes.

---

## 🧪 Testing

We value reliability above all else. Every new tool, command, or provider must include comprehensive automated tests.

Run:
```bash
go test ./...
```
before opening a PR, because that's what CI will execute.

You can also run tests across all internal domain modules with race detection enabled during local development:
```bash
go test -race ./...
```

For detailed testing guides and verification strategies, check **`docs/testing.md`**.

---

## 🔍 Submitting PRs

When you submit a Pull Request:
1. Ensure your PR description clearly explains **what** the change does, **why** it was implemented that way, and **how** it was tested.
2. Link any related GitHub issues or discussions.
3. Keep PRs atomic and focused on a single responsibility or stable component.
4. Verify all tests pass by executing `go test ./...` locally before submitting.
5. Be responsive to review feedback — we aim to review all PRs constructively and promptly!

If you have questions or ideas for a major architectural change or a new core module, please open a GitHub Discussion or Issue first to discuss it with the maintainers before writing extensive code.

Thank you for helping shape **WindMist** into a top-tier open-source project! 🚀

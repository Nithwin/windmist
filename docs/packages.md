# WindMist Package Overview

This document provides a high-level overview of the domain-driven Go architecture in WindMist (`v2.0.0`). It is designed to help open-source contributors navigate the codebase and understand the separation of concerns.

## `cmd/`
The entrypoint of the WindMist CLI. We use [Cobra](https://github.com/spf13/cobra) to define terminal commands.
- **`cmd/windmist`**: Contains the root command (`windmist`), alongside subcommands like `chat`, `set`, `show`, `uninstall`, and `version`. It parses flags and initializes the `internal/chat` TUI.

## `internal/`
The core application logic, protected from being imported by external Go modules.

### `internal/agent`
The brain of WindMist. This package contains the stateless, multi-turn reasoning loop. 
- Formulates the system prompt.
- Orchestrates tool execution (detecting when the AI requests a tool, running it, and feeding the result back).
- Manages sub-agents for heavy background tasks.
- Triggers summarization logic to compress long contexts.

### `internal/ai`
The provider-agnostic LLM interface.
- Defines standard interfaces (`Provider`, `Tool`, `Message`) that all models must implement.
- Houses the translation layers that convert WindMist's standard schemas into provider-specific API formats (e.g., Gemini's `v1beta` function calling API).

### `internal/chat`
The interactive Terminal User Interface (TUI). Built entirely with [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- Contains the central `Update` loop (event state machine) for handling keystrokes.
- Manages UI components: text areas, viewport scrolling, splash screens, and inline menus.
- Connects the UI events to the `internal/agent` execution logic.

### `internal/config`
Configuration management.
- Locates and parses the `~/.windmist/config.yaml` file.
- Defines standard struct definitions for Provider settings, UI themes, and Remote configurations.

### `internal/mcp`
The **Model Context Protocol (MCP)** implementation.
- Contains the JSON-RPC client capable of connecting to external standard MCP servers (like GitHub, Postgres, or SQLite).
- Dynamically translates remote MCP tools into WindMist `ai.Tool` schemas on the fly.

### `internal/rag`
The local memory and vector database engine.
- Implements a pure-Go **TF-IDF (Term Frequency-Inverse Document Frequency)** embedder and searcher.
- Slices codebase files into chunks and indexes them so the agent can perform semantic queries without relying on expensive remote embeddings.

### `internal/remote`
The headless remote-control architecture.
- **`remote.Hub`**: A thread-safe broadcaster that streams agent output to connected controllers.
- **`telegram`**: The first official remote controller, utilizing a Telegram Bot to receive commands (`/ask`, `/provider`) and forward them into the main `chat.Model`.

### `internal/store`
SQLite persistence layer.
- Manages the local `windmist.db` file.
- Stores historical chat sessions, messages, and timestamps so you can resume work after closing the terminal.

### `internal/tools`
The atomic capabilities granted to the AI.
- **`filesystem`**: Safe reading, writing, and directory traversal.
- **`editing`**: Precision code mutation tools (`replace_text`, `insert_text`).
- **`system`**: Bash execution and environment inspection (guarded by user-approval).
- **`web`**: HTTP requests and scraping.

### `internal/ui`
Styling and rendering utilities.
- Uses [Lip Gloss](https://github.com/charmbracelet/lipgloss) to define the color palette (`BrandCyan`, `MutedLight`, etc.).
- Configures [Glamour](https://github.com/charmbracelet/glamour) for rich markdown rendering (syntax highlighting, bolding, lists) tailored specifically for terminal outputs.

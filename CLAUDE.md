# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Lazygit - Terminal UI for Git Commands

Lazygit is a simple terminal UI for git commands, written in Go. It provides an interactive interface for common git operations like staging, committing, branching, rebasing, and more.

## Key Commands for Development

### Building and Running
- `make build` - Build the binary with debug flags
- `make run` - Build and run lazygit locally
- `make install` - Install lazygit to your Go bin directory
- `go run main.go` - Direct run during development

### Development and Debugging
- `make run-debug` - Run lazygit in debug mode (shows debug logs)
- `make print-log` - Show log output in separate terminal (run alongside run-debug)
- `lazygit --debug` - Alternative debug command
- `lazygit --logs` - Show logs from a completed lazygit session

### Testing
- `make unit-test` - Run unit tests (fast, uses `go test ./... -short`)
- `make test` - Run all tests (unit + integration)
- `go run cmd/integration_test/main.go tui` - Interactive integration test runner (TUI)
- `go run cmd/integration_test/main.go cli [testname]` - CLI integration test runner
- `go test pkg/integration/clients/*.go` - Integration tests for CI

### Code Quality
- `make format` - Format code using gofumpt (stricter than gofmt)
- `make lint` - Run golangci-lint via scripts/golangci-lint-shim.sh
- `make generate` - Generate auto-generated files (test lists, cheatsheets)

## Code Architecture

### High-Level Structure
Lazygit follows a layered architecture with clear separation between git operations, GUI logic, and presentation:

- **Commands Layer** (`pkg/commands/`): All git operations and OS interactions
- **GUI Layer** (`pkg/gui/`): User interface, event handling, and application logic  
- **Models** (`pkg/commands/models/`): Data structures for git objects (commits, branches, files)
- **Integration Tests** (`pkg/integration/`): End-to-end testing with simulated user input

### Key Architectural Concepts

**Contexts**: Each view (branches, commits, files, etc.) has a corresponding context that manages state and handles keypresses. Contexts are the primary way GUI state is organized.

**Controllers**: Define keybindings and their handlers. Controllers can be shared between multiple contexts (e.g., list navigation controller).

**Helpers**: Shared business logic used by multiple controllers. When a controller method needs to be used elsewhere, it gets extracted to a helper.

**Views**: The visual representation managed by the underlying gocui library. Views maintain their own content buffer.

**Common Pattern**: Most structs have a `c` field containing a "common" struct with shared dependencies (logger, config, i18n, etc.).

### Important Packages
- `pkg/commands/git_commands/` - All git binary communication
- `pkg/gui/controllers/` - Keybinding definitions and handlers
- `pkg/gui/context/` - State management for each view
- `pkg/gui/controllers/helpers/` - Shared business logic
- `pkg/gui/presentation/` - Data formatting for display
- `pkg/config/` - User configuration handling
- `pkg/i18n/` - Internationalization

### Configuration and Keybindings
- User config is defined in `pkg/config/user_config.go` with live reloading support
- Panel visibility can be controlled via config options (showStatusPanel, showFilesPanel, etc.)
- Custom keybindings and commands are supported via user configuration
- New keybindings should be added to controllers rather than the legacy `pkg/gui/keybindings.go`

### Development Workflow
1. New features typically start in controllers/helpers
2. Integration tests should be written for user-facing functionality
3. The codebase is transitioning from a "God Struct" pattern to controllers/contexts
4. Code formatting uses gofumpt (stricter than gofmt)
5. All new text should be internationalized via `pkg/i18n/english.go`

### Testing Strategy
- Unit tests for individual functions and components
- Integration tests that simulate full user sessions
- Integration tests can be run in "sandbox mode" for manual testing
- Test files end in `_test.go` for unit tests
- Integration tests live in `pkg/integration/tests/`

## Panel Toggle Feature Context
The codebase currently includes a panel toggle feature allowing users to hide/show different panels:
- **Shift+1**: Toggle status panel
- **Shift+2**: Toggle files panel  
- **Shift+3**: Toggle branches panel
- **Shift+4**: Toggle commits panel
- **Shift+5**: Toggle stash panel

Key files for panel management:
- `pkg/gui/controllers/helpers/window_arrangement_helper.go` - Layout calculations with safety checks
- `pkg/config/user_config.go` - Panel visibility configuration
- `pkg/gui/keybindings.go` - Panel toggle keybindings

Critical safety feature: Layout helper ensures at least one panel (status) remains visible to prevent crashes when too many panels are hidden.
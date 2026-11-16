# AGENTS.md — Lazygit AI Coding Agent

You are an expert in Go engineering and terminal-UI (TUI) development, and the
primary AI agent supporting the development of **lazygit** — a cross-platform
terminal UI for git, written in Go. Your goal is to write, debug, and refactor
lazygit code while rigorously following the project's portability, safety, and
coding standards.

**Core principle:** lazygit is not a single-platform tool. It is a Go program
that must run flawlessly on **Linux, macOS, and Windows** (including WSL, Git
Bash, and Cygwin environments). Every change must be validated against all
three platforms before it is considered done.

---

## 1. Portability & Cross-Platform Safety (Go)

- **Paths:** Never hardcode `/` or `\` as a path separator. Always use the
  `path/filepath` package (`filepath.Join`, `filepath.Separator`, etc.) so
  paths work on Windows. Be aware that Windows paths are case-insensitive and
  may contain spaces; macOS filesystems are typically case-insensitive too.
- **Line endings & text:** Do not assume `\n` line endings or a single shell.
  Git output parsing must tolerate `\r\n` on Windows.
- **Symlinks:** macOS and Windows handle symlinks and path resolution
  differently. Never assume a path is canonical; resolve it explicitly when
  correctness depends on it.
- **No OS-specific assumptions:** Avoid `syscall`, `golang.org/x/sys`, or
  build-tagged files unless the feature is genuinely platform-specific, and
  if so, provide implementations for every supported OS.
- **Go version:** The project targets **Go 1.25** (see `go.mod`). Do not use
  language features or stdlib APIs newer than that.

## 2. Coding Standards

- **Formatting:** Go uses tabs via `gofmt`, but this project enforces the
  stricter **`gofumpt`**. Never hand-format code — run `make format`
  (`gofumpt -l -w .`) or let your editor format on save.
- **Linting:** Your code must pass `make lint`
  (`./scripts/golangci-lint-shim.sh run`), which enables `revive`,
  `errorlint`, `staticcheck`, `nakedret`, `unparam`, `prealloc`, and more.
  If you need an exception, use a targeted `//nolint` with a justification —
  never silence a whole file or package.
- **Naming:**
  - Use `self` as the receiver name in struct methods — this is encouraged.
  - Prefix interfaces with `I` (e.g. `IHelper`) when the interface has
    several methods, rather than suffixing with `er`.
  - When a struct implements an interface, assert it explicitly:
    `var _ MyInterface = &MyStruct{}`.
  - Follow standard Go naming (no stutter, exported vs. unexported is
    meaningful). Note: `ST1003` (e.g. `remoteURL` vs `remoteUrl`) is
    intentionally disabled project-wide — do not fight it.
- **Variable declarations:** Prefer short declarations (`:=`) where the type
  is obvious; use `var` when the zero value is the intent. Avoid shadowing
  outer variables, and never use naked returns (`nakedret` is enforced with
  `max-func-lines: 0`).
- **Errors:** Return errors rather than panicking. Wrap with context
  (`fmt.Errorf("...: %w", err)`). Do not swallow errors silently — if an
  error is intentionally ignored, say so in a comment.

## 3. Command Execution & Security

This is the Go equivalent of "always quote your variables":

- **Never build shell command strings** by concatenating user input, branch
  names, file paths, or config values. This is a command-injection risk.
- Always invoke external commands with `os/exec.Command` (or the project's
  wrappers in `pkg/commands/oscommands`), passing each argument as a
  **separate slice element**. Let the OS handle argument boundaries.
- All communication with the `git` binary lives in
  **`pkg/commands/git_commands`** — do not shell out to git from GUI or
  controller code. All other OS interaction lives in
  **`pkg/commands/oscommands`**.
- Before invoking an external tool, verify it exists (e.g. `exec.LookPath`)
  and handle its absence gracefully with a user-facing error. Never assume
  `git` is on `PATH`.

## 4. Output, Logging & Error Reporting

- **Never use `fmt.Println` / `fmt.Printf` for diagnostics.** Use the project
  logger: `self.c.Log.Warn("...")`, `gui.Log.Warn("...")`, or — only from
  code inside `vendor/` where no logger is available —
  `logs.Global.Warn("...")` (requires `LAZYGIT_LOG_PATH` to be set).
- For user-facing errors, surface them through the GUI's error/toast
  mechanisms so the user sees them, rather than writing to stdout/stderr.
- Keep log output useful: log at `Warn`/`Error` for things a user or
  developer would want to see; avoid noisy `Info`/`Debug` in hot paths.

## 5. Project Structure (Single Source of Truth)

See `docs/dev/Codebase_Guide.md` for the full map. The essentials:

- `pkg/app` — startup, logging, and user-config initialisation.
- `pkg/commands/git_commands` — **all** git binary communication.
- `pkg/commands/oscommands` — OS interaction and general command invocation.
- `pkg/commands/models` — model structs (commits, branches, files, …).
- `pkg/config` — user config struct and defaults (`user_config.go`).
- `pkg/gui` — the GUI. Prefer the layered structure over the legacy God
  Struct: **controllers** (keybindings + handlers) → **helpers** (shared
  logic) → **contexts** (per-view state) → **views** (gocui). Lower layers
  must not depend on higher ones.
- `pkg/i18n` — internationalised strings (see §6).
- `pkg/common` — the `Common` struct (`c`) holding shared dependencies
  (logger, i18n, config). Most structs carry a `c` field; access shared
  services through it (e.g. `self.c.Helpers.MyHelper`).
- `pkg/utils` — low-level helpers. Put generic utilities here, not inline
  duplicates.
- `vendor/github.com/jesseduffield/gocui` — the underlying TUI library
  (event loop, rendering, input). See §9 before touching it.

## 6. Internationalisation (i18n)

- **Never hardcode user-facing strings** in GUI code. Add a new field to the
  `TranslationSet` struct in `pkg/i18n/english.go` and set its English value
  in `EnglishTranslationSet()`, then access it via `gui.Tr.YourNewText` (or
  `self.c.Tr.YourNewText`).
- Use **Sentence case** for all user-facing text (not Title Case).
- Translation files under `pkg/i18n/translations/` are managed via Crowdin —
  **do not edit them directly**.

## 7. Configuration

- New user-config options are defined in `pkg/config/user_config.go` with
  sensible defaults. Access config through `self.c.UserConfig` so values
  auto-reload when the user edits their config file.
- If a value cannot auto-reload, hook `Gui.onUserConfigLoaded`, or add it to
  `Gui.checkForChangedConfigsThatDontAutoReload` so the user is prompted to
  restart.
- After changing the config schema, regenerate the JSON schema and
  cheatsheets with `make generate` (`go generate ./...`).

## 8. Concurrency & the Event Loop

- The event loop runs in gocui's `MainLoop`. **Never block the UI thread**
  with slow I/O (git calls, file reads, network).
- For asynchronous work use `self.c.OnWorker(myFunc)`; to touch the UI again
  from a worker, use `self.c.OnUIThread(myOtherFunc)`.
- Be careful with shared mutable state across goroutines — prefer passing
  data through channels or confining mutation to the UI thread.

## 9. Vendored Dependencies (gocui & lazycore)

- `vendor/github.com/jesseduffield/gocui` and the lazycore packages are
  vendored forks. **Do not edit `vendor/` directly** as a permanent change.
- To change them: edit in `vendor/` to test locally, port the change to the
  upstream repo (gocui `awesome` branch / lazycore `master`), then bump the
  version with `./scripts/bump_gocui.sh` or `./scripts/bump_lazycore.sh`
  (`make bump-gocui` / `make bump-lazycore`).
- After any dependency change, run `make vendor` (`go mod vendor && go mod tidy`).

## 10. Testing

- **Unit tests:** Go tests in `*_test.go` files. Run with `make unit-test`
  (`go test ./... -short`).
- **Integration tests:** end-to-end tests in `pkg/integration` that record
  real lazygit sessions. See `pkg/integration/README.md` and
  `docs/dev/Integration_Tests.md`. Run with `make test` (unit + integration)
  or the targeted `make integration-test-tui` / `make integration-test-cli`.
- Every new feature or bug fix should come with a test. Integration tests
  must pass on **all platforms** (Linux, macOS, Windows) — avoid tests that
  depend on a specific shell, path separator, or case-sensitive filesystem.
- If a test needs external state, use the existing fakes/mocks in `pkg/fakes`
  rather than reaching for real git repos or the network.

## 11. Commits & Pull Requests

- Keep a **clean, minimal commit history**: one logical change per commit,
  with a good commit message explaining *why*, not just *what*.
- **Separate refactorings from behavior changes** into different commits.
- When iterating during review, use **fixup commits** (see
  `docs/Fixup_Commits.md`) rather than piling on top.
- All changes go through Pull Requests from a feature branch (never from
  `master`). If you add code that needs documentation, update the docs in
  the same PR.
- Before opening a PR, ensure `make format`, `make lint`, and `make test`
  all pass.

---

## Pre-Flight Checklist

Before writing any code, mentally run through:

1. **Windows / macOS impact:** path separators, case sensitivity, line
   endings, symlinks, shell differences (WSL / Git Bash / Cygwin).
2. **Injection safety:** is any user-controlled value reaching a shell
   command? If so, pass it as a separate `exec.Command` argument.
3. **Layering:** does this belong in `git_commands`, `oscommands`, a helper,
   a context, or a controller? (Lower layers must not depend on higher ones.)
4. **i18n:** is any new user-facing text going through `pkg/i18n`?
5. **Config:** does a new option need a default, JSON-schema regen, and
   auto-reload handling?
6. **Threading:** does this block the UI thread, or does it need
   `OnWorker` / `OnUIThread`?
7. **Tests:** is there a unit or integration test covering this, passing on
   all three platforms?
8. **Format & lint:** have you run `make format` and `make lint`?

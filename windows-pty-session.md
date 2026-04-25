# Session notes: Windows ConPTY support

Repo: `lazygit` (Go terminal UI for git). Branch: `windows-pty`. Base: `master`.

## Goal

Add ConPTY-based pseudo-terminal support on Windows so custom pagers (delta etc.)
and credential prompts work the same as on Unix. Pre-existing state had
creack/pty on Unix and a no-op stub on Windows that only set `LAZYGIT_COLUMNS`.

## Branch state (oldest → newest)

```
af8368718 Abstract task command over *exec.Cmd
c9ece72a9 Abstract pty startup behind a platform-specific primitive
d45388b7b Add pty support on Windows via ConPTY
62c4e8b36 Demonstrate that unknown escape sequences leak as literal text
b106d0ac4 Silently consume unrecognized escape sequences
1bae298a2 Stop leaking other malformed and unimplemented escape sequences
30eac7509 fixup! Add pty support on Windows via ConPTY    ← pending autosquash into d45388b7b
423cdadc8 Use ConPTY on Windows for pty-backed command execution
```

Pending: `git rebase -i --autosquash <base>` to fold `30eac7509` into `d45388b7b`.
The post-fold commit message of `d45388b7b` may want a `reword` — its scope grew
to include moving the pty primitive from `pkg/gui` to `pkg/commands/oscommands`,
not only adding ConPTY.

## Architecture

### Why the task API needed widening

Windows ConPTY can't attach a child via `os/exec`. Setting
`PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE` requires `CreateProcess` directly
(golang/go#62708). So the ConPTY path can't return an `*exec.Cmd` to the task
runner.

`pkg/tasks/tasks.go` defines:

```go
type Cmd interface {
    Wait() error
    String() string
    GetProcess() *os.Process
}

type ExecCmd struct{ *exec.Cmd }
func (c ExecCmd) GetProcess() *os.Process { return c.Process }
```

`NewCmdTask` accepts `func() (Cmd, io.Reader)`.
`oscommands.TerminateProcessGracefully` takes `*os.Process` (not `*exec.Cmd`).

### Shared pty primitive

`pkg/commands/oscommands/pty.go` (cross-platform):

```go
type Pty interface {
    io.ReadWriteCloser
    Resize(cols, rows uint16) error
}

type StartedPty struct {
    Pty     Pty
    Process *os.Process
    Wait    func() error
}

// func StartPty(cmd *exec.Cmd, cols, rows uint16) (StartedPty, error)
```

`pkg/commands/oscommands/pty_unix.go` — creack/pty impl. Wait = `cmd.Wait`.
`pkg/commands/oscommands/pty_windows.go` — ConPTY impl. Wait wraps `proc.Wait()`
and synthesizes a non-nil error on nonzero exit (mimics `*exec.Cmd.Wait`
semantics).

### Consumers

`pkg/gui/pty.go` — `newPtyTask`. Wraps `StartedPty` into `tasks.Cmd` via local
`ptyCmd` adapter (uses `cmd.String()` for description, the explicit
`sp.Process` for `GetProcess`).

`pkg/commands/oscommands/cmd_obj_runner.go` — `cmdHandler` struct grew a
`wait func() error` field because Windows ConPTY path never runs
`exec.Cmd.Start`. Non-pty path: `wait = cmd.Wait`. Pty path: `wait = sp.Wait`.
`runAndStreamAux` calls `handler.wait()` instead of `cmd.Wait()`.

The per-platform `cmd_obj_runner_{default,windows}.go` files are deleted; the
single cross-platform `getCmdHandlerPty` lives in `cmd_obj_runner.go` now.

### gocui escape interpreter changes

`pkg/gocui/escape.go` (in-tree, not vendored — pulled in from upstream gocui by
a sibling branch this branch was rebased on top of).

ConPTY emits a session-init stream of escape sequences (private modes, cursor
positioning, RIS, etc.). The pre-existing parser errored on anything outside
SGR / EL / OSC-8 hyperlinks; `view.go` then rendered the unparsed bytes as
literal cells. Visible as junk at the start of pager output.

Made the parser tolerant in `parseOne`:

- `stateEscape`: a single byte 0x30–0x7E (e.g. `ESC c`) is consumed as a
  complete Fs/Fp sequence.
- `stateCSI`: accept `;` (empty first param), accept private-prefix bytes
  0x3C–0x3F (`<`, `=`, `>`, `?`), accept any final byte 0x40–0x7E
  immediately after `[`.
- `stateParams`: accept any final byte 0x40–0x7E we don't implement; accept
  intermediate bytes 0x20–0x2F by transitioning to a new
  `stateCSIDiscard` state that swallows until the final byte.
- `stateOSCWaitForParams`: on non-`;`, transition to `stateOSCSkipUnknown`
  rather than erroring (would have leaked the rest of the OSC body).
- The CSI-too-long sanity checks at the top of `parseOne` now transition
  to `stateCSIDiscard` instead of returning an error.
- The `m` (SGR) case is transactional: snapshots `curFg/curBgColor` and
  restores them if `outputCSI` errors mid-loop, so a malformed SGR
  (`[1;;m`) doesn't leave a partial color state.

Removed: `errCSITooLong`, `errOSCParseError` (no callers).

`TestParseOneIgnoresUnknownSequences` covers all the new paths plus asserts
no side effects (instruction stays `noInstruction{}`, fg/bg stay
`ColorDefault`).

## Key behavioral changes for users

- `LAZYGIT_COLUMNS` is now set on every platform (was Windows-only). Kept for
  backwards-compat with the script in `docs/Custom_Pagers.md`. Harmless on
  Unix (no docs reference it there).
- Custom pagers now actually work on Windows via real terminal emulation.
- `cmd_obj_runner` pty path (used by custom commands with
  `output: logWithPty`) now uses ConPTY on Windows instead of falling
  through to non-pty.

## What's not done / known caveats

- **Windows runtime not validated** by me. Cross-compile passes. User
  confirmed the original ConPTY commit works against a real Windows lazygit
  build (escape-sequence junk was the visible bug; fixed). The post-refactor
  state (with primitive moved to oscommands and cmd_obj_runner rewired)
  hasn't been Windows-tested yet.
- **Risk hot spots if Windows misbehaves:**
  - Handle lifetime in `winPty.Close` — `ClosePseudoConsole` then close pipe
    fds. Order matters; closing pipes first can hang the child.
  - `os.FindProcess(int(pi.ProcessId))` after `CreateProcess` — narrow race
    where the pid is reused before we open it. Negligible in practice.
  - `createEnvBlock` errors on NUL in env values (`UTF16FromString` returns
    error). Not handled gracefully; would surface as a `StartPty` failure.
  - `cmd_obj_runner` pty handler hardcodes initial size 80×24. The gui side
    resizes via `Pty.Resize`; the runner side doesn't currently. May need a
    way to pipe view dimensions through if a custom command's output looks
    wrong.
- **`go vet` `unsafeptr` warning** on `unsafe.Pointer(hpc)` in
  `oscommands/pty_windows.go`. The cast is correct per Microsoft's ConPTY
  sample (HPCON is the attribute value, not a pointer to one). `go test`'s
  default vet set excludes `unsafeptr`, so CI doesn't flag it. Documented
  inline.
- **Spec-compliance edges** intentionally skipped in the escape parser:
  ESC + control bytes, double-ESC reset, CAN (0x18) / SUB (0x1A) abort.
  Discussed in conversation; deemed low-value.

## Build / test commands

```sh
go build ./...
GOOS=windows go build ./...
go test ./...
./scripts/golangci-lint-shim.sh run --timeout 5m
```

The lint pass is what CI gates on. `golangci-lint-shim.sh` wraps the project's
pinned version. `gofumpt -d <file>` to preview formatting; `gofumpt -w` to
apply (lint will catch violations otherwise).

## Repo conventions ([AGENTS.md](AGENTS.md))

- Every commit must compile + pass tests.
- Commit messages explain WHY (motivation, constraint, bug being fixed),
  not WHAT (the diff shows that).
- Separate prep refactors from behavior changes — land the refactor first.
- For bug fixes, demonstrate the bug in test(s) in commit N, fix in commit
  N+1. Use the EXPECTED/ACTUAL inline-comment pattern: live assertion is
  the broken behavior so the test passes against unfixed code; correct
  expectation lives in a `/* EXPECTED: ... ACTUAL: */` comment that gets
  flipped in the fix commit.
- Commit proactively when a logical unit is done.

## File index

| File | Role |
|------|------|
| `pkg/commands/oscommands/pty.go` | `Pty` interface, `StartedPty` struct, `StartPty` doc |
| `pkg/commands/oscommands/pty_unix.go` | creack/pty impl of `StartPty` |
| `pkg/commands/oscommands/pty_windows.go` | ConPTY impl + `createEnvBlock` helper |
| `pkg/commands/oscommands/cmd_obj_runner.go` | `cmdHandler{stdoutPipe,stdinPipe,close,wait}`, `getCmdHandlerPty`, `getCmdHandlerNonPty`, `runAndStreamAux` |
| `pkg/commands/oscommands/os_default_platform.go` | `TerminateProcessGracefully(*os.Process)` Unix impl |
| `pkg/commands/oscommands/os_windows.go` | `TerminateProcessGracefully(*os.Process)` no-op |
| `pkg/gui/pty.go` | `newPtyTask`, `onResize`, `ptyCmd` adapter |
| `pkg/gui/gui.go:88,735` | `viewPtmxMap map[string]oscommands.Pty` |
| `pkg/gui/tasks_adapter.go` | `newCmdTask`, returns `tasks.ExecCmd{Cmd: cmd}` |
| `pkg/gocui/escape.go` | `parseOne` state machine; `stateCSIDiscard` is the new state |
| `pkg/gocui/escape_test.go` | `TestParseOneIgnoresUnknownSequences` |
| `pkg/tasks/tasks.go` | `Cmd` interface, `ExecCmd`, `NewCmdTask` |

## Gotchas to remember

- The dependency direction is `tasks → oscommands`. So `oscommands` cannot
  import `tasks`. The pty primitive lives in `oscommands`; `pkg/gui`
  consumes it and adapts to `tasks.Cmd` locally.
- `pkg/gocui` is in-tree (not vendored). The escape interpreter changes are
  edits to our own code, not patches to a third-party dep.
- `pkg/gocui/view.go:946-957` is what renders the unparsed-bytes-as-cells
  fallback; that's the "leak" path the escape fixes are about.
- The `unsafeptr` cast pattern (`unsafe.Pointer(hpc)`) for ConPTY: matches
  Microsoft's sample and `microsoft/hcsshim`. `go vet` warns; ignore.
- Behavior contract for the escape parser: unrecognized but well-formed,
  AND outright malformed → silently consume. Leaking bytes to the view is
  never the right answer.

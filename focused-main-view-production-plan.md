# Focused main view — productionization plan

The plan for turning the `use-delta-hyperlinks-for-clicking-in-diff` prototype
into production PRs. The prototype branch is **throwaway** ([[prototype-branch-throwaway]]):
none of its history lands. Every PR below is **re-implemented from scratch on a
fresh branch off master**, using the prototype code as a reference to transcribe
from — not to cherry-pick. The knowledge lives in `focused-main-view-notes.md`
(referenced below as **N§x**) and `diff-line-metadata-notes.md` (**M§x**).

This document is the working plan for the (many) production sessions. Keep it
current: check off commits as they land, record deviations, and add findings.

---

## 1. How to use this document (read first, every session)

- **Ground rules:** AGENTS.md applies in full — small self-contained commits,
  every commit compiles + passes tests + `just format` + `just lint`, prep
  refactors split from behavior changes, "why not what" messages, `fixup!`/
  `amend!` for iteration, no conventional-commit prefixes, `just generate`
  after keybinding/test changes, docs in `docs-master/` only, translatable
  strings via Go templates, no PRs created by agents (the user opens them).
- **Terminology: say "diff renderer", never "pager".** "Pager" is a misnomer
  that leaked out of the implementation (the `GIT_PAGER` env var); the OSC 1717
  spec now says "diff renderer" throughout, and so do we — in code identifiers,
  commit messages, PR descriptions, docs, and user-facing strings. The
  prototype's identifiers and test names still say "pager" in places
  (`ProbePagerEmitsDiffMetadata`,
  `stage_from_main_view_with_unsupported_pager`, …); rename them during
  transcription. `GIT_PAGER` itself and the "stdin-pager invocation route"
  keep their technical names where the mechanism is literally meant. The
  user-facing config rename has **landed on master** (#5870 — see PR 3),
  which also restructured the config (`type` field); master and the rebased
  prototype already use the new names (`DiffRendererConfig`,
  `DiffRendererConfigManager`, `cycleDiffRenderers`, …).
- **Threading:** master has the landed main-thread rework — see §2.8 for the
  contract every PR must honor. The prototype is already rebased on top of it.
- **Prototype references** are given as *subject line* plus short SHA. The
  prototype branch gets rebased; **find commits by subject, not SHA**. The
  SHAs quoted in this plan predate the 2026-08-04 rebase onto master (past
  the gocui mouse-gesture PR #5854 and the diff-renderer config rework
  #5870) and resolve only on the pre-rebase copy kept at
  `fold-staging-functionality-into-main-view-plan`; the rebased branch is
  the live prototype and the source of truth for code shape. When a plan
  item says "reference: X", read that prototype commit (message + diff)
  before implementing — it usually contains the design rationale and the
  gotchas.
- **Transcribe the final state, not the journey.** The prototype iterated
  heavily; several mechanisms were built, reverted, and rebuilt differently.
  §3 below lists everything that must NOT be ported. When in doubt, the
  *current tree* of the prototype branch is the source of truth for code
  shape; the notes are the source of truth for *why*.
- **Branch naming/stacking:** each PR gets its own branch; PR N+1 branches off
  PR N's branch (linear stack). Rebase the stack when a PR merges.
- **Verification:** every PR runs `just build`, `just unit-test`, `just lint`,
  `just e2e`. PRs touching the async render or diff-renderer paths additionally
  need the interactive sign-off listed in §6 — the headless harness cannot
  exercise the pty/renderer path (N§13.1: cmd-path only, env allowlist blocks
  `LAZYGIT_SLOW_RENDER`).

## 2. Locked scope decisions (do not relitigate)

Decided with the user (2026-07-17/18 planning sessions, plus earlier locked
decisions from the notes):

1. **The staging and patch-building panels are removed.** `enter` on a file
   (files panel and commit-files panel) focuses the main view at that file's
   diff. The `Staging`/`StagingSecondary`/`CustomPatchBuilder` contexts, views,
   and `patch_exploring` machinery are deleted (PR 9). The prototype kept them
   as an A/B reference (N§21.24) — production does not.
2. **`enter` / double-click in the focused main view are dropped.** With the
   explorers gone the dive gesture has no target. Enter is unbound there;
   double-click behaves like a single click (select). Esc exits; space/d/e act
   on the selection.
3. **Sequencing:** stacked PRs, merged in quick succession; **no release ships
   with both staging UIs**. Brief coexistence on master between PRs 7/8 and
   PR 9 is fine.
4. **The nav/selection and position-preserve features land as their own PRs
   before the staging series** (PRs 5 and 6), independently releasable.
5. **Both extras are in scope** as final small PRs: alt/shift-click-to-edit
   (PR 10) and open-PR-at-selected-line (PR 11).
6. **The hyperlink backend is dropped** (N§14.5): no `lazygit-edit://`-based
   line identity resolution. (Master's existing click-a-delta-hyperlink feature
   in `pkg/gui/gui.go` is untouched — only the prototype's use of hyperlinks as
   an identity *backend* dies.)
7. **Selection is always shown in diff main views** (N§21.6), anchored at the
   first visible change line / hunk; no on-demand toggle, no config for it.
   Non-diff main views (branch log, …) keep no selection.
8. **Follow master's landed threading contract.** The main-thread-mutation
   rework has landed on master (Front A #5767; "Synchronize async view
   rendering" #5791; the popup, command-log, and status fronts;
   `RefreshFromWorker` + captured refresh inputs; `BatchUIUpdates` replacing
   BLOCK_UI), and the prototype is rebased on top of it (merge-base
   `4cf12a5b7`). The contract for all new code: **model/context/selection
   mutations are UI-thread-only** (asserted in `-debug` builds); **view
   geometry and origin (`ox`/`oy`) are UI-thread-only**; the view's **write
   buffer may be written off-thread under `writeMutex`, which is kept
   permanently** (the batch-to-UI-thread attempt deadlocked and was
   abandoned); task-side origin/dimension work bounces via
   `ViewBufferManager.onUIThread`; worker-issued refreshes use
   `RefreshFromWorker` with UI-owned inputs captured up front; incidental
   display work bounces as *background* tasks so it doesn't count toward
   `Busy()` and block repo switches. The prototype tree is already adapted
   (atomic `loading` flag, bounced `onEndOfInput`/`onNewKey`, `OnUIThread`
   hops in restore `Apply`, off-screen swap under `writeMutex`) — transcribe
   it as-is; where a threading question comes up, the answer is "match the
   current tree / master's contract". See memory
   [[main-thread-over-mutexes-direction]] and `docs/dev/Repo_Switch_Safety.md`.
9. **Side-by-side staging includes all records on a row** — no stage-one-side
   gesture (N§21.3, accepted restriction).
10. **Diff-renderer capability detection is the empty-input handshake probe**
    (N§21.30), not render observation. Non-conforming renderers get the
    raw-diff fallback at focus time; a `type: extDiff` entry with empty
    `command` (git's `diff.external`; formerly `useExternalDiffGitConfig`)
    is always-raw when focused (documented limitation).
11. **The escape/snapshot machinery is never built** (`FocusedMainViewSnapshot`,
    `EscapeFromPatchExplorer`, the N§12.2 escape routing): it existed only to
    return *from the explorers*, which no longer exist.

## 3. Prototype work that must NOT be ported

The branch contains superseded/reverted work. Do not transcribe any of these
(listed with the thing that replaced them):

| Not ported | Superseded by / reason |
|---|---|
| `gui.showSelectionInFocusedMainView` config; on-demand space-toggled selection | always-shown selection (N§21.6) |
| Middle-visible-line as the *selection* anchor ("Select the line in the middle…") | first-visible-change/hunk anchor. (`MiddleVisibleLineIdx` itself survives as the `-U`-preserve anchor when no selection shows — PR 6) |
| `enter`/double-click dive into staging/patch-building at a line; `CommitFilesHelper.EnterCommitFile` threading | gesture dropped (§2.2) |
| `FocusedMainViewSnapshot`, `EscapeFromPatchExplorer`, escape-restore-by-identity, N§12.2 routing | explorers removed (§2.11) |
| Hyperlink identity backend (`GetFileAndLineForClickedDiffLine` hyperlink parsing, `HyperLinkInLine` as a *backend*) | buffer-parse + OSC metadata only (§2.6) |
| `ScrollToOriginYForNextTask` / `thenForNextTask` / `KeepOriginForNextTask` / `LinesToRead.ApplyInitialScroll` | `RenderRestore` (PR 6) — build the final mechanism directly |
| "Hold the placeholder until first paint" + `freshViewLineCount` stale-tail guard (both reverted on the branch too) | off-screen render (PR 1) |
| Observe-at-focus diff-renderer detection (N§21.29) | handshake probe (N§21.30) |
| `HighlightInset` and `selectionBgColorEdgeWidth` experiments | `narrowSelectionHighlight` (N§21.34) |
| `matchByWorktreeChange` and `AdjacentChangeLine` reveal matchers | change-line-ordinal reveal (N§21.17) |
| The three separate handler channels (`onClick`/`onStage`/`onTogglePatch FocusedMainViewFn`) | build `FocusedMainViewActions` directly in its final one-interface shape (N§21.25) |
| Unconditional gutter reset-on-preview-render + the `strings.Join(cmd.Args)` content-equality hack | focused-pair-only gutter model (N§21.22(3), N§21.35) |
| `backUpOverHeader` file-nav landing | land on the first located row; f/h header records make headers resolvable ("Parse the f/h header records…", af98be48d) |
| "Pager" naming in identifiers, strings, docs | "diff renderer" (§1 terminology; the config rename landed on master as #5870) |
| OSC number `456`, env vars `EMIT_OSC1717_METADATA`/`OSC1717_METADATA` | OSC **1717**, env var **`OSC1717`** (final rename, 665149b11) |
| The in-repo spec file | the spec lives in its own repo, <https://github.com/stefanhaller/diff-line-metadata-spec> (`~/git-repos/diff-line-metadata-spec`); it was briefly on an `osc-1717-spec` branch in between, so notes naming that branch are stale too |
| Session-notes commits, `.claude/settings.json` commits, WIP commits | n/a |

## 4. The PR stack (overview)

PR titles become release-notes lines — they are written for users. Order is
dependency order; 1–6 are independently releasable; 7–9 merge in quick
succession (§2.3); 10–11 any time after their dependencies.

| # | Title (draft) | Depends on | Nature |
|---|---|---|---|
| 0 | Validate the context a custom command names — **DONE: landed on master as #5989** | — | fix (found reviewing PR 9) |
| 0b | Render the focused main view again while it is being searched | — | fixes, gocui/tasks (foot of the stack) |
| 0c | Render diffs without a pty on Windows | — | fix, oscommands/gui (foot of the stack; design-reviewed 2026-09-18) |
| 1 | Fix flicker, scroll glitches, and crashes in async diff rendering | — | fixes, gocui/tasks |
| 2 | Internal: resolve diff lines to (file, line, kind) identities | 1 | infra |
| 3 | Rename the "pagers" config to "diff renderers" — **DONE: landed on master as #5870** | — | rename + migration |
| 4 | Support diff renderers that emit OSC 1717 diff line metadata | 2, 3 | infra + protocol |
| 5 | Select, navigate and edit diff lines in the focused main view (copy moved to PR 6b) | 2 (4 for renderers) | feature |
| 6 | Keep your position in the diff when changing context size, ignoring whitespace, or switching diff renderers | 1, 2, 5 | feature |
| 6b | Copy the selected diff lines from the focused main view | 5, 6 | feature (split out of PR 7 on 2026-09-13) |
| 7 | Stage, unstage and discard changes directly from the focused main view | 4, 5, 6, 6b | feature |
| 8 | Build custom patches directly from a commit's diff view | 7 | feature |
| 9 | Replace the staging and patch-building panels with the focused main view | 7, 8 | removal + migration |
| 10 | Alt- or shift-click a diff line to open it in your editor | 2, 4, 5 | feature |
| 11 | Open the selected diff line in the branch's GitHub PR | 5 | feature |
| 12 | Jump to a file of the diff from a menu | 5 | feature |

---

## 5. Per-PR plans

### PR 1 — Fix flicker, scroll glitches, and crashes in async diff rendering

All standalone master-worthy fixes; users benefit regardless of the rest of
the series. Everything here lives in `pkg/gocui`, `pkg/tasks`, and the gui
layout/render plumbing.

**Re-validate each fix against current master before transcribing.** The
landed rework (#5791) restructured exactly this area — the read loop now
bounces `onEndOfInput`/`onNewKey` to the UI thread, `readLines` is an atomic
pointer, several buffer accessors gained locks. The prototype commits are
rebased on top of all that and are the authoritative shapes, but a fix may
have been absorbed or may need reshaping; check before porting.

**Status: DONE 2026-08-09** on branch `fix-async-diff-rendering`. 15 commits:
the 11 planned below plus four extra (two prep/demonstrate commits, the
loading-indicator gate, and the `FlushStaleCells` removal). No `fixup!`/
`amend!` left in the branch — see "History" below. `just build`, `just
unit-test`, `just lint`, `just e2e` all green, and every individual commit
builds and passes `pkg/tasks`, `pkg/gocui` and `pkg/gui`.
**§6 interactive sign-off: APPROVED (2026-08-09)** — tested with and without
`LAZYGIT_SLOW_RENDER`, no flicker. It turned up two things, both fixed:
deviation 8 (the scroll no longer reset in a VS Code terminal) and deviation 9
(the loading indicator blanking the view).

**History.** Both interactive findings were first written as a tip commit /
`amend!`, then folded so that no commit in the branch introduces a regression a
later one repairs — the user asked for the tidier history and suggested the
decomposition. `newContentPending` is now introduced **before** the off-screen
render, by its own commit that gates the loading indicator on it (a sensible
change in its own right at that point: don't clear the buffer and flash
"loading..." over content we're about to render identically). The off-screen
render then inherits the gate rather than introducing the blanking regression,
and the first-paint commit just adds the second reader — which is what makes
the `amend!` unnecessary. Verified by diffing the restructured tip against the
pre-restructure one: identical apart from two deliberately unified comment
wordings.

The branch base is **not master**: at the user's request the stack sits on
`fold-staging-into-main-view` (= the tip of `scroll-selection-into-view`),
which touches the same scroll code and merges to master soon. Since
2026-08-09 there is one branch below it: `fix-task-key-race`, an unrelated
data-race fix found during this work (see "Found and fixed alongside"), kept
separate so it can merge to master on its own.

Commits (in order):

1. **Route all view origin writes through `SetOriginX`/`SetOriginY`** — pure
   prep chokepoint refactor. Ref: b0a85eefb. ✅
2. **Add `LAZYGIT_SLOW_RENDER` debug knob** — sleeps N ms per written line so
   async render frames become visible; inert when unset. Needed by reviewers
   to see what the later commits fix. Ref: e8682b3fd. ✅
3. **Lock the event-thread viewLines readers that still skip `writeMutex`** —
   task goroutines rebuild `viewLines` under `writeMutex`; master still reads
   them unlocked on the event-handling thread in the click-path hyperlink
   lookup (`pkg/gocui/gui.go` ~1686, `viewLines[newY].line[newX].hyperlink`)
   and in `onMouseMove`/`findHyperlinkAt` (hover) — the rework locked
   `ClearViewLines`/`IsTainted`/`Buffer` but not these; audit for any others.
   The prototype fixed this class in its own `HyperLinkInLine` reader (part of
   the dropped hyperlink backend) and in `onMouseMove`; transcribe the lock
   pattern onto the readers master actually has. Use the demonstrate-then-fix
   pattern if a deterministic test is feasible. Refs: 2cc42fc81, a44bf5d05.
   ✅ (landed as **two** commits, see deviations)
4. **Fire queued ReadToEnd callbacks when the initial read reaches EOF** —
   the read loop abandoned queued `{Total:-1}` requests when the initial
   request hit EOF, silently dropping their `Then`. Ref: b6f99abc6.
   ✅ (landed as **two** commits — a deterministic test turned out feasible,
   so demonstrate-then-fix applied)
5. **Don't scroll a view up to fill blank space while its content is loading**
   — the layout clamp used the partially-loaded height; add the `loading`
   flag (an `atomic.Bool` in the rebased shape) and skip the clamp while
   `IsLoading()`. Ref: 695842291. ✅
6. **Reset other main views' scroll after copying content, not before** —
   `refreshMainViews` zeroed the source view's origin before `CopyContent`
   used it, so every cross-pair placeholder jumped to the top. *(Verified
   against master 2026-07-18: still needed — `CopyContent` still copies the
   source's `ox`/`oy` and `refreshMainViews` still resets first; master's
   `1efcfcc14` only stopped sharing the buffer slices.)* Ref: c35c9316c. ✅
7. **Bundle a view's cell buffer and write state into a `viewBuffer`** — prep.
   Ref: fd858cd98. ✅
8. **Make the buffer-writing methods operate on a `viewBuffer`** — prep.
   Ref: 2cfc0e24d. ✅
8b. *(added)* **Don't take the view over for "loading..." when the content
   isn't changing** — introduces `newContentPending`; must come **before** the
   off-screen render, see deviation 9. ✅
9. **Render async content into an off-screen buffer and swap it in** — the
   core mechanism: cmd/pty tasks write to `View.offscreen`; at first-paint
   (or EOF) the buffer swaps in atomically (under `writeMutex` — buffer
   writes are the writeMutex domain per §2.8); `refreshViewLinesIfNeeded`
   truncates view lines on swap (kills the stale-tail class);
   `clear()`/`Reset()` abandon an in-progress off-screen render. Includes the
   **scrollbar height freeze** (`FreezeScrollbarHeight` at `StartLoading`,
   release at EOF and in `clear()`) — in the prototype this was an `amend!`
   into the same commit precisely because the off-screen render introduces
   the scrollbar regression; keep them together here too. Tests:
   `TestOffscreenRender`, `TestBufferLineForViewLineStaleTail`,
   `TestScrollbarHeightHeldWhileLoading`,
   `TestScrollbarHeightReleasedWhenContentReplaced`. Refs: 27ce0a6bc + its
   scrollbar amend, N§13.5, N§13.6. ✅
10. **Don't run end-of-input handling for a render that was stopped** — the
    stopped-task EOF coin-flip (`select` between stop and closed `lineChan`)
    let a stopped task swap in a truncated buffer. No deterministic test (the
    bug *is* the nondeterministic select) — justified skip, N§13.6. Ref:
    8e3dc3eff. ✅
11. **Reset the scroll to the top at first paint, not when the task starts** —
    with the off-screen render the old content stays visible until the swap;
    resetting oy at task start made it jump first. Ref: 411681502. ✅
    (the reset keys off `newContentPending` so it outlives its task — see
    deviation 8)
12. *(added)* **Drop `FlushStaleCells`, which no longer has anything to
    flush** — see "Found and fixed alongside". ✅

Notes:
- If commit 9's stale-tail test needs `BufferLineForViewLine`, introduce that
  accessor here (PR 2 then reuses it) rather than contorting the test.
  *(Not needed — see deviations; PR 2 still introduces the accessor.)*
- Gotcha for the future: **fast renders unmask ordering transients that slow
  renders hide** (N§20.5). Re-test at normal speed *and* under slow render.

#### Deviations from the plan (2026-08-08, as implemented)

1. **Commit 3 split in two.** Master's click-path lookup indexes `viewLines`
   inline in `gui.go`'s event loop, so there was nowhere to take the view's
   lock. Prep commit "Move the click-path hyperlink lookup onto View" extracts
   a `View.hyperlinkAt(x, y)` (behaviour-preserving, sits next to
   `findHyperlinkAt`); the fix commit then locks it and `onMouseMove`. The
   audit for other unlocked `viewLines` readers found only dead code
   (`realPosition`/`Line`/`Word`/`ViewBuffer`/`LinesHeight` have no callers
   outside gocui), so nothing else was touched. No test: the bug is a data
   race, not reachable single-threaded.
2. **Commit 4 got a deterministic test after all.** A blocking reader
   (`BlockingLineReader`) holds a task in its still-loading state, so a
   `ReadToEnd` can be queued and *then* the task released to EOF. Landed as
   demonstrate-then-fix per AGENTS.md ("Add a test for a read request queued
   while a task reaches EOF", then the fix swapping EXPECTED/ACTUAL).
3. **Commit 9 fixed a `screenColMax` gap the prototype still has.**
   `escapeInterpreter.screenColMax` — the pty screen width soft-wrap tracking
   counts against — is set only on `v.buf.ei` (`NewView`, `SetContentWidth`).
   `BeginOffscreenRender` builds a *fresh* interpreter, and `SetContentWidth`
   runs inside `start()`, strictly before the first line triggers
   `beginRender()`, so the off-screen interpreter would have `screenColMax = 0`
   and count no soft wraps at all — ConPTY `CUP` escapes after a wrapped line
   would then land on the wrong row. `BeginOffscreenRender` now copies it over,
   guarded by `TestWriteCursorPositionEscapeInOffscreenRender`. **The prototype
   branch has the same gap** (`screenColMax` landed on master 2026-06-30 /
   07-08, long after the off-screen commit was written, and the rebase rewrote
   `v.ei` → `v.buf.ei` mechanically); the user decided to leave it there.
4. **Commit 9's truncation test doesn't use `BufferLineForViewLine`.** The
   truncation is observable through `ViewLinesHeight`/`ViewBufferLines` alone,
   so `TestBufferLineForViewLineStaleTail` became
   `TestViewLinesTruncatedByShorterRender` and the PR-2 accessor stayed in
   PR 2.
5. **Commit 10 calls `callThen()` before bailing.** The prototype's stopped-EOF
   bail drops the request's `Then`, unlike both explicit-stop branches right
   above it, which fire it. Matched those instead (it does *not* drain the
   queue — reaching EOF is what justifies the drain, and a stopped task hasn't).
6. **The origin reset is manager state, not a `LinesToRead` field.** The
   prototype has the cmd/pty wrappers compute `cmdStr != manager.GetTaskKey()`
   on the UI thread and pass it down as `LinesToRead.ResetOrigin`, which adds
   unsynchronized reads of `taskKey` (written by the `NewTask` goroutine under
   `taskIDMutex`). Instead `NewTask` keeps making the decision, under the lock
   exactly where it already did, and records it on the manager as
   `newContentPending` — the same flag the loading indicator keys off (see
   commit 8b); the gui wrappers need no change at all.
   **PR 6 note:** `RenderRestore` competes with that flag rather than reading a
   field — a restore places the scroll itself, so it wins:
   `if restore != nil { restore.Apply(swapIn) } else { swapIn(); if
   self.newContentPending.Swap(false) { self.resetOrigin() } }`. Note the
   restore must then *also* clear the pending reset, or the next same-key
   render would inherit it and jump to the top.
7. **Commit 11 runs the first paint on the UI thread**, resolving the
   prototype's `// TODO: should probably use OnUIThread?`. It has to: the paint
   now writes the origin, and swap + origin must land in one hop or a draw
   between them shows the new content at the previous render's scroll — the
   N§20.5 failure mode. `firstPaint` itself is therefore unwrapped and each
   call site wraps it (EOF already did).
8. **The pending origin reset outlives its task** (commit 11, 2026-08-09).
   Deferring the reset to the first paint would make it die with a task
   that was replaced before it ever painted. Found interactively, and only in a
   **VS Code terminal**: VS Code delivers the focus-in event before the click,
   the focus handler's `Refresh` does its git work on a worker, so its main-view
   re-render lands *after* the click's — same command key, so it decides on no
   reset of its own, and it stops the click's task before its first paint. (Zed
   delivers the click first, so that task paints and the bug never shows.) So
   the reset keys off the manager's `newContentPending` — set by `NewTask`,
   consumed by whichever first paint gets there — rather than per-task state
   (deviation 6 describes the shape, including what it means for PR 6). Guarded
   by `TestResetOriginSurvivesTaskReplacement`.
9. **The loading indicator only takes the view over for new content**
   (commit 8b, 2026-08-09). Also found interactively. The 200ms indicator clears
   the view to show its message; that was nearly invisible before, because the
   untruncated view lines kept the previous render below it, so only line 1
   changed — and scrolled down you never saw it at all. With the truncation it
   blanks the whole view, which is ugly precisely when it's least useful: a slow
   re-render of *unchanged* content (a background refresh over a repo with
   submodules that have uncommitted changes) blanks and then paints the same
   thing back. Restoring the old look is not an option — it *was* the stale-tail
   inconsistency. Gated the indicator on `newContentPending`, so it never
   clears content it is about to render identically. Guarded by
   `TestLoadingIndicatorOnlyTakesOverForNewContent`, which observes the
   indicator itself (a writer that watches for the "loading..." bytes) rather
   than the `beforeStart` callback, so it holds at every point in the branch —
   before the off-screen render, `beforeStart` is also called by the read loop.
   **Placement:** the blanking would arrive with the off-screen render, so the
   gate goes *before* it (commit 8b) — which is why `newContentPending` is
   introduced there rather than with the origin reset.

#### Found and fixed alongside

- **`GetTaskKey()` was an unsynchronized read of `taskKey`** (`tasks.go`), which
  the `NewTask` goroutine writes under `taskIDMutex`, while the string renders
  in `tasks_adapter.go` read it from the UI thread. Pre-existing on master and
  unrelated to async rendering, so it went on its own branch,
  `fix-task-key-race`, **below** this stack (2026-08-09) — a torn read of a
  two-word string value can index out of bounds, so it is worth landing on
  master without waiting for the stack. Rebased `fix-async-diff-rendering` onto
  it; two trivial conflicts on the one line both branches touch.
- **`FlushStaleCells` removed** (PR 1 commit 12). With
  `refreshViewLinesIfNeeded` truncating, the `onEndOfInput` call to it only
  forced a full re-wrap of the whole buffer. The prototype kept it; there is no
  reason to.

Interactive sign-off (user, `just debug` + `LAZYGIT_SLOW_RENDER` matrix):
flicking through commits/files scrolled down; the 10 s auto-refresh with
`refreshInterval: 3`; scrollbar stability. See §6.

### PR 2 — Internal: resolve diff lines to (file, line, kind) identities

The host-side primitive: rendered row → `(path, kind, new/old line)`. No
user-visible change; the PR description should say what it enables. Backends
in this PR: **buffer-parse only** (raw / `--color` / structure-preserving
renderer output). Precedence seam is built here; the metadata backend slots
in via PR 4.

Commits:

1. **patch package: line-number arithmetic + well-formedness** —
   `LineNumberOfLine`/`OldLineNumberOfLine` (quirk-free inverse maps),
   `PatchLineForLineNumber`/`PatchLineForOldLineNumber`, `Patch.IsWellFormed`
   (hunk-header lengths vs parsed body — the buffer-parse integrity gate,
   M§8), hunk-length capture in `parse.go`. **Must be rename-aware from day
   one**: master now has rename support in the patch builder (`f84ada494`),
   and the prototype's patch-package changes predate it — the prototype's
   failing `patch_building/renamed_file_whole` (N§21.36(2)) marks exactly
   where `Parse`/`Transform`/`FormatView` must reproduce rename headers.
   Write unit tests for renames here. Refs: 2e5151cdf, 9c0bb5357 (parser
   parts), N§21.36(2).
2. **gocui: displayed-buffer accessors** — `DiffLineContents` (text +
   metadata-slot + per-line data for unwrapped buffer lines),
   `BufferLineForViewLine` / `ViewLineForBufferLine` (wrapping-aware mapping,
   unless already landed in PR 1). Unit tests. Refs: ca095604c ("Add
   View.BufferLineForViewLine…"), 792c7a294.
3. **The resolver: `types.DiffLineInfo` + the batch buffer parser** —
   `pkg/gui/types/diff_line_info.go`, `pkg/gui/controllers/helpers/
   diff_line_parser.go` (`parseAllDiffLinesFromBuffer` → `parseFileSection`,
   one parse per file section — **O(n), never per-line**; the single-line
   resolver delegates to it), `StagingHelper.resolveDiffLines` /
   `GetDiffLineInfo` seam with backend precedence (metadata → buffer-parse;
   metadata arrives in PR 4). Port the prototype's unit tests
   (`diff_line_parser_test.go`, `diff_line_info_test.go`,
   `diff_line_navigation_test.go` comes with PR 5). Refs: 7cf9b5037,
   9c0bb5357, 556ba1213 (final O(n) shape — build it O(n) directly; N§20).
4. **Decode C-quoted paths in the buffer parser** — flagged as an unclosed
   prototype gap (M§8): git C-quotes unusual paths in `diff --git` headers;
   production must decode them.

Gotchas:
- The **two-call atomicity constraint** (M§8): never resolve via two separate
  locked gocui calls that can interleave with a re-render; snapshot content
  and index together (the `DiffLineContents` snapshot approach does this).
- Multi-file diffs: `fileSectionBounds` handles rows outside any file section.

#### Deviations from the plan (2026-08-09, as implemented)

Landed as 5 commits on branch `resolve-diff-lines-to-identities` (off PR 1):
old-file line numbers; the well-formedness gate; the gocui view/buffer line
mapping; the identity type + parser + helper; header-path decoding.

1. **The seam lives on a new `DiffLineHelper`**, not `StagingHelper` (decided
   with the user). `helpers/diff_line_helper.go` holds `GetDiffLineInfo` and
   needs nothing from the staging panel but the worktree path, so PR 9 doesn't
   have to relocate it out of a file it otherwise deletes. PR 5's navigation
   attaches here too.
2. **`PatchLineForLineNumber`/`PatchLineForOldLineNumber` are not ported**
   (decided with the user). Production resolves patch lines by *scanning* the
   parsed patch and matching `(line number, deletion?)` identities; the
   prototype's only callers were `patch_exploring`'s dive-at-a-line, which
   §2.2 drops. They'd be dead exported API.
3. **`DiffLineInfo` ships with `IsChange` only** (decided with the user).
   `PatchSelectLine`, `SamePatchLine` and `PullRequestAnchor` land with their
   first consumers (PRs 7, 6 and 11) rather than as an unused API surface.
4. **No `DiffLineContent`/`DiffLineContents` in gocui yet** (decided with the
   user). With no metadata to carry, the parser takes the `[]string` the
   existing `BufferLines()` already returns; PR 4 introduces the struct when
   there is a second field for it, and changes the two signatures that read it.
   So PR 2's gocui change is just `BufferLineForViewLine` /
   `ViewLineForBufferLine`.
5. **The helper-level batch resolver is deferred to PR 5.** The O(n) batch
   *parser* (`parseAllDiffLinesFromBuffer`) lands here as planned and is
   unit-tested against the single-line parser, but `resolveDiffLines` /
   `resolvedDiffLine` — which only add the backend precedence on top — wait for
   their first whole-buffer scan.
6. **The rename gap (N§21.36(2), §8's mandatory PR 2 row) does not
   reproduce.** `Parse`/`Transform`/`FormatView` already reproduce rename
   headers on a master-based branch: the header lines pass through verbatim and
   master's `StripRename` handles the partial case, so the whole
   `patch_building` e2e group — `renamed_file_whole` included — is green with
   the new patch-package code. The prototype's failure was an artifact of its
   pre-rename fork point, not of the parser changes. Rename coverage was added
   as unit tests anyway (`OldLineNumberOfLine` over a rename header; the parser
   over a pure rename, a rename with a content change, and a quoted rename).
7. **Header paths need more than C-quote decoding** (commit 5, found while
   implementing): git also terminates the path field of a `---`/`+++` line with
   a **tab** when the path contains a space, so a name like `with space.go`
   parsed to a path with a trailing tab. Both are decoded together, guarded by
   `TestPathFromDiffHeaderField` against real git output shapes.
8. **Commit 1 split in two** (line-number arithmetic; the well-formedness gate),
   and commit 2 needs no stale-tail guard: PR 1's truncating
   `refreshViewLinesIfNeeded` makes `viewLines` consistent with the buffer by
   construction, which is what M§8 predicted.
9. **Code comments carry no notes-file vocabulary.** The prototype's comments
   say "mechanism #1/#2", cite `diff-line-metadata-notes.md` §s, and say
   "pager"; none of that means anything in the repo (the notes and the spec live
   outside it). Transcribed comments describe the mechanism instead — worth
   doing on every later PR too.

#### Review round 1 (2026-09-05) — a diff the pane holds only part of

Focusing the main view over a long diff gave no selection at all with no diff
renderer configured, and none of the commands that act on one. A pane loads a
diff a screenful at a time and reads the rest as the user scrolls
(`linesToReadFromCmdTask` caps the first read at `height*(height-1)` lines), so
what it holds ends part way through a hunk. `Patch.IsWellFormed` counted that
hunk against the length its header declares, found it short, and refused the
whole file section — leaving every row of the diff on screen unresolvable.
Renderers that state each line's identity were unaffected, git's own diff being
the one the buffer parser is for.

Two `fixup!` commits, inserted mid-branch:

1. `Patch.IsWellFormedSoFar` on "Add a well-formedness check for a parsed
   patch": every hunk but the last has to match its header exactly, as before,
   and the last one only has to fit within what its header declares. The check
   exists to tell a faithful rendering from a restructured one, and it still
   does that. A rendering that moves the +/- marker off the start of the line
   makes a change read as context, and a context line counts towards both
   lengths, so such a hunk comes out **longer** than its header declares, never
   shorter.
2. `parseFileSection` takes `endsTheBuffer` on "Recover the identity of a diff
   line from the rendered diff", and holds a section in that position to what
   has arrived. Only that section: one that another section follows is all
   there, so a hunk short of its header there means the rendering restructured
   the diff.

The behavioural demonstration is an e2e test, and it has to live where the main
view gains a selection at all, so it went to PR 5's round 1 below.

`Rename parsedDiffLine.RelPath to Path` (PR 4) gained the one occurrence the
fixup's test added, since the branch below it now has one more.

#### Review round 2 (2026-09-18) — a submodule of the repo

The user found that a commit changing a submodule had no submodule in it as
far as the host was concerned: `n`/`N` passed over it, the menu of the diff's
files left it out, and clicking its name in the diffstat said the diff held no
such file. A submodule gets no `diff --git` header — lazygit asks for
`--submodule`, so git states which commits it moved between and lists them
("Submodule modules/xyz a32f27c..2d9f921:") — and the scan for the section a
row belongs to looked for that header alone.

One `fixup!` on "Recover the identity of a diff line from the rendered diff":
a `Submodule …` row opens a section too, and `pathFromDiffHeader` reads the
path off it. The section holds no hunks, so `patch.Parse` makes every row of
it a header row, which is what they are — what git states there is which
commits the submodule moved between, not lines of a file. Both shapes git
writes are recognized (the commit has moved, or the working tree is dirty);
`TestSubmodulePath` holds the six forms of the first and the two of the
second.

Two renderers needed the same thing said on their side, since nothing can
parse what they print (their branches, unpushed):

- **delta** tagged the row with no record at all, having been given one in the
  first place to stop it naming the file above it. It now reads the path off
  the row, the same way. A `fixup!` on "Emit `f` and `h` records on file and
  hunk header rows".
- **diff-so-fancy** passed the row through as an ordinary line, whose
  classifier answers for content lines only. A `fixup!` on its commit of the
  same name adds `submodule_header_path`, and the one row that names a
  submodule gets an `f`.

The user's delta binary predates the earlier fixup that stopped a submodule
row carrying the previous file's record, which is why the menu showed "." for
it (a record with an empty path, resolved against the worktree) rather than
the file above.

The first cut of the parser change let the section run to the next one, as a
file's does, and the user found what that costs under a renderer: the menu
listed the submodule again between every pair of files. A renderer prints no
`diff --git` line, so there was nothing below the submodule to stop at, and
every row carrying no record of its own — the blank rows delta puts between
one file and the next — came out as a row of the submodule. A second `fixup!`
on the same commit ends the section where what git writes for the submodule
ends: the line naming it, and the log of the commits it moved over. The rows
below are placed by the records the renderer states for them, and the log
lines now resolve under a renderer too, since it passes them through as git
wrote them.

Checked against all three backends over a commit whose submodule sorts first
and whose files have several hunks each, and over `--submodule=diff`, which
lazygit never asks for but which degrades sensibly (the submodule, then the
file inside it).

`Rename parsedDiffLine.RelPath to Path` (PR 4) again gained the occurrences
the fixups' tests added, as in round 1 — once per fixup, since each added one.

### PR 3 — Rename the "pagers" config to "diff renderers" — DONE (master #5870)

**Landed on master** as #5870 ("Rework the custom pager config (rename to
diff renderer)", merge `d8d09e1f9`), independently of this series, and it
went further than this plan asked: besides the rename
(`git.pagers` → `git.diffRenderers`, `cyclePagers[Reverse]` →
`cycleDiffRenderers[Reverse]`, docs, migration), it **restructured the
entries**:

- Each entry has an explicit **`type`** field: `stdinFilter` (default) |
  `extDiff` | `rawGit`.
- The old `pager:` / `externalDiffCommand:` fields are unified into one
  **`command`** field, interpreted per `type` (this dissolves the section's
  old open question about the `pager:` field name).
- `useExternalDiffGitConfig: true` became `type: extDiff` with empty
  `command` (= git's `diff.external`).
- New **`rawGit`** type with an `args` field (e.g. `--color-words`) that
  runs plain git with extra args — a renderer flavor the prototype never
  had to consider (see PR 7 commit 10).
- `pkg/config/diff_renderer_config_manager.go`:
  `DiffRendererConfigManager` with a `DiffRendererType` enum and typed
  accessors (`GetDiffRendererType`, `GetStdinFilterCommand`,
  `GetExternalDiffCommand`, `GetRawGitArgs`, …); no entries configured ⇒
  `RawGit`.

**Consequence for the remaining PRs:** wherever this plan says to key a
decision off "is a pager configured" / "is it an external diff command"
(probe route selection, env advertising, raw fallback, renderer cycling),
base it on `GetDiffRendererType()` instead of querying the pager and
ext-diff commands individually — mostly a simplification. The rebased
prototype is already adapted to the new config and manager names.

### PR 4 — Support diff renderers that emit OSC 1717 diff line metadata

Reads the protocol; sets the env var so conforming renderers emit. No
consumer behavior changes yet (consumers land in PRs 5–8), but the PR is the
public face of the protocol on the lazygit side — write the description for
diff-renderer authors, link the spec
(<https://github.com/stefanhaller/diff-line-metadata-spec>).

Commits:

1. **gocui: parse OSC 1717 records and attach them per cell** — escape.go
   parsing; payload attached to the following cell region; metadata cleared at
   line boundaries (no bleed); **keep a content-less sentinel cell when a
   blank line carries pending metadata** (delta renders some blank changed
   lines with no cells — N§21.15 bug 1); **swallow the version-only handshake
   record** (N§21.30, `TestDiffLineMetadataHandshakeSwallowed`); multi-record
   rows: a row's *all* distinct payloads are exposed per buffer line
   (side-by-side rows carry two — N§17.1, N§21.12) — as
   `DiffLineContent.Metadata []string`, not the prototype's separate
   `DiffLineMetadataPayloads()` accessor (deviation 1). Unit tests
   incl. wrapped rows (every renderer-wrapped output row carries the record —
   M§10.8). Include the **zero-width record regions** handling ("Keep OSC
   1717 records whose region is zero-width", fe8022827): back-to-back records
   with no cells between them — e.g. difftastic's combined file+hunk banner,
   or a modification pair collapsed to one column — get their orphaned
   payloads drained into content-less carrier cells, so every record on a row
   stays discoverable. Without it the `d` half of a collapsed modification
   row is invisible to every consumer (PR 8's secondary-pane identity bridge
   relies on it). Refs: 1cc7ecbbb, e8385b3cf, 13595f0a8, 3018289e8 (gocui
   half), fe8022827.
2. **The metadata backend, slotted ahead of buffer-parse** — resolve
   `a`/`d`/`c` records to `DiffLineInfo`; **accept the `f`/`h` header
   records** (file header: no line number; hunk header: first line of its
   hunk) mapping to the same header identities the buffer parser reports;
   `SamePatchLine` requires header kinds to match (a hunk header shares its
   number with the hunk's first content line — af98be48d's message has the
   full reasoning). Refs: 836f768cb, af98be48d.
3. **Advertise the protocol to diff renderers** — set `OSC1717=V1` in the
   environment of renderer/ext-diff invocations (pty task env + ext-diff cmd
   env); which route applies now follows from the entry's `type` (#5870).
   **Including `rawGit`:** git itself emits the records for its word-diff
   formats, so the advertisement must be set **before** `newPtyTask`'s
   early return for `rawGit` (that route needs no pty, since only a pager
   needs a terminal to be spawned — but git still needs to be asked). Put
   it next to `LAZYGIT_COLUMNS`, which sits there for the same reason.
   Getting this wrong is silent: the renderer's output simply carries no
   records, and every consumer falls back as if the renderer were
   non-conforming. Ref: 9975a8fac + 665149b11 (final name: `OSC1717`),
   "Advertise the metadata protocol to git as well, not only to a pager".

#### Deviations from the plan (2026-08-09, as implemented)

Landed as 7 commits on branch `support-osc-1717-diff-metadata` (off PR 2):
multi-digit OSC numbers; the record parsing; the handshake; the records that
cover no cell; the `RelPath`→`Path` rename; the metadata backend; the
advertisement.

1. **`DiffLineMetadataPayloads` is not built** (decided with the user).
   `gocui.DiffLineContent` is `{Text string, Metadata []string}` — *every*
   record on the row, in the one locked snapshot. The prototype split this in
   two (a first-payload field on `DiffLineContent`, plus an all-payloads
   accessor), which made PR 7's `ChangeLinesInViewRange` take two separately
   locked gocui calls that a re-render can interleave — the M§8 two-call
   hazard. **PR 7 commit 4 and PR 8 commit 7 read `contents[i].Metadata`**
   instead of calling a second accessor.
2. **Matching must compare all of a row's records — binds PR 6** (raised by the
   user during this PR). The prototype's restore resolves each row to *one*
   identity (`resolveDiffLines` → `findResolvedDiffLine` → `SamePatchLine`), so
   a target captured under a unified rendering (an `a` at line N) never matches
   the side-by-side row whose leftmost record is the `d` of the same
   modification — which is exactly the unified↔side-by-side renderer cycle PR 6
   exists to serve. PR 6's matcher has to ask "does *any* record on this row
   match", which deviation 1's field makes possible. Nothing to build in PR 4:
   its single consumer, `GetDiffLineInfo`, resolves a row to one identity by
   contract (the leftmost record), which is what `e` / alt-click / open-PR want.
3. **No `Hyperlink` field** on `DiffLineContent` (§2.6 drops the hyperlink
   backend), and **`DiffLineMetadataInLine` is not ported** — it had no
   production consumer even in the prototype.
4. **`parsedDiffLine.RelPath` renamed to `Path`**, in its own commit: a
   renderer states the path however it likes, absolute included. The
   prototype's two conversions (`diffLineInfoFromParsed` /
   `diffLineInfoFromMetadata`) collapse into one `diffLineInfo`, which joins
   the worktree path only when the path is relative.
5. **The OSC parser needed a prep refactor**: master dispatches on a single
   character, so only the single-digit OSC 8 could be recognized. Accumulating
   the number and dispatching on it is commit 1, behavior-preserving.
6. **A record carried at the line's end marks itself consumed**, so a line
   finished twice (a CR followed by an LF with nothing written between) doesn't
   get two carrier cells. Small addition over the prototype's shape.
7. **e2e coverage in CI**: `diff/diff_renderer_metadata` configures a fake
   conforming renderer that reports `$OSC1717` back and prefixes every line
   with a record. It proves the advertisement reaches the renderer through the
   pty path (verified: it fails without the `pty.go` change) and that records
   don't leak into the rendered text. The §6 interactive pass with real patched
   renderers is still owed.

Cross-repo note: the reference emitters live on `osc-1717-metadata` branches
in `/Users/stk/Stk/Dev/Builds/{delta,difftastic,diff-so-fancy}`; nothing is
upstreamed yet. difftastic's emitter commits were rewritten (2026-07-18) to
emit patch-space `d`+`a` records for collapsed modification rows — rebuild
before verifying. lazygit must remain fully functional without any conforming
renderer (buffer-parse + PR 7's raw fallback guarantee this). Interactive
verification of this PR needs locally built patched renderers.

**git is one of those emitters, and lazygit ships its side regardless of
whether git ever takes the patch.** The emitter lives on branch `osc-1717`
in `/Users/stk/Stk/Dev/Builds/git` (4 commits, word diffs only, unproposed
as of 2026-08-07); it may never be accepted, in which case the fallback is a
maintained fork for the users who want it. That doesn't change what lazygit
does, and nothing here is conditional on the outcome: the probe asks the
installed git what it can do, so a stock git answers "no records" and every
consumer degrades to exactly the behaviour of a non-conforming renderer —
the raw fallback when focused, buffer-parse otherwise. A git that speaks the
protocol turns `--color-words` from a browse-only renderer into a fully
usable one. Write the PR descriptions so that they don't promise a git
feature that doesn't exist upstream: describe the protocol and the probe,
not "works with git --color-words".

### PR 5 — Select, navigate, edit and copy diff lines in the focused main view

The focused main view (already reachable via `0`/click on master) gains a
real selection and line-level interactions. After this PR: diff main views
always show a selection when focused; ↑/↓/v/a move/extend it; `<left>`/
`<right>` jump hunks, `n`/`N` files, `f` opens a jump-to-file menu; `e` edits
the selected line; `ctrl+o` copies the selection; click/drag select.

Commits:

1. **Fold `ViewSelectionController` into `MainViewController`** — prep; the
   nav handlers need to live where the mode state is consulted. Ref:
   b92d71e29.
2. **Introduce the `DiffMainViewContext` classifier** —
   `GetDiffMainViewType() DiffMainViewType` (`None|Staging|PatchBuilding`) on
   the side-panel contexts (files=Staging; commitFiles/localCommits/
   subCommits/stash=PatchBuilding; reflog=PatchBuilding **from day one** —
   see PR 8 reflog item; others None). In this PR only "≠ None" is read (is
   this a diff main view). Refs: f470d870f (marker origin), a760f9ef5
   (classifier final shape), N§21.25.
3. **The selection model** — `DiffSelectState` on `MainContext`
   (`pkg/gui/context/main_context.go`): mode Line/Range/Hunk, sticky range,
   `userEnabledHunkMode`; selection rendering via the view's native cursor +
   `SetRangeSelectStart` (no new highlight machinery); always-shown selection
   anchored at first visible change line on focus (`0`, click, `<tab>`);
   hunk-default from `useHunkModeInStagingView` selects the first visible
   block; mode-aware ↑/↓ (hunk steps blocks; non-sticky range collapses;
   sticky extends), `v`, `a`, shift-↑/↓; pages/top/bottom drop hunk mode;
   clicks select (hunk-on-click when in hunk mode / config default, context
   lines stay single-line — N§21.32); `<tab>` seeds the landing pane's
   anchor and select state. **Mandatory: the hunk-default must implement the
   staging view's `IsSingleHunkForWholeFile` refinement** (a whole-file
   single-block diff — new/deleted file, no context — drops to line mode
   instead of select-everything); skipping it would regress vs master
   (decided 2026-07-19; overrides N§21.11's deferral). Whether the diff is a
   single file is known from the **side panel's selection** — don't derive
   that from the rendered content. For the single-block computation the lean
   is to compute it **in patch space from the file's raw diff**, fetched
   synchronously at focus (the same "the decision is synchronous" shape as
   N§21.3); the accepted fallback is a ReadToEnd of the rendered content on
   focus (user has OK'd the cost). Decide at implementation; surface other
   options if any appear. Refs: f470d870f, f4d5c79da (selection-model
   half), 5312357ce, 5688e8b87, 4e78aa4c4, 9b8249a60, N§21.10–21.11, N§21.32,
   [[diff-selection-state-home]].
4. **Selection visibility rules** — no selection over placeholders/no-diff
   content; hide when changes vanish (render-side hook in the panel's
   render-to-main + focus-side check via `ViewHasChangeLines`). e2e:
   `no_selection_when_no_changes`, `hide_selection_after_discarding_last_change`
   (adapt: discard via files panel until PR 7). Ref: 7901de3d4, N§21.27
   bug 4.
5. **Drag-to-range** — `dragAnchorViewLine` on MainContext; `MouseLeft` +
   `ModMotion` binding re-anchors at the mouse-down line. **Includes the gocui
   driver fix**: report the first drag movement as a drag, not a release
   (tcell_driver `MAYBE_DRAGGING→DRAGGING`). Refs: d6fd8c808, 0fa35ee42,
   N§21.32(5).
6. **Hunk and file navigation** — `<left>`/`<right>` change-block nav ("hunk"
   = lazygit change block, not `@@`), `n`/`N` file nav landing on the file's
   first located row (header under conforming sources; first content line
   otherwise — no `backUpOverHeader`), anchor's file found by scanning *down*
   (`anchorFilePath`, N§16.5); selection showing → move+scroll-into-view; the
   pure index arithmetic unit-tested (`diff_line_navigation_test.go`).
   **Mandatory: a target beyond the lazily-loaded portion must be found, not
   silently no-op'd** (decided 2026-07-19; overrides N§16.4's deferral): when
   the scan exhausts the loaded content without a match, ReadToEnd and
   re-scan (the `openSearch` shape). If commit 3's `IsSingleHunkForWholeFile`
   solution ends up reading to the end on focus anyway, this comes for free.
   Applies equally to the jump-to-file menu (commit 7 — its file list must
   cover the whole diff). Refs: 559955f7c, af98be48d (landing changes),
   N§16.2.
7. **Jump-to-file menu (`f`)** — menu of the diff's files in order,
   repo-relative; reuses the file-nav landing logic. **Production must add
   proper i18n strings** (prototype hardcoded English). Ref: 27b1012e1.
8. **Edit the selected line (`e`)** — resolve via `GetDiffLineInfo`,
   `AdjustLineNumber`, open editor; editing a file-header row opens the file
   without a line. Refs: 467806fba, af98be48d (header case).
9. **Copy the selection as raw diff lines (`ctrl+o`)** — copy the
   corresponding lines of the **original raw diff**, never the renderer's
   output (decided 2026-07-19; supersedes the prototype's verbatim copy with
   renderer-gated `dropDiffPrefix`, N§21.28 — you don't want a renderer's
   restructured text on your clipboard, and this dissolves the
   "can't tell whether the renderer preserves the +/− column" problem).
   Mechanism, all synchronous from existing pieces: selected view rows →
   buffer rows via `BufferLineForViewLine` (never `SelectedLines()`, it's
   wrapping-unaware) → identities (all payloads per row, so SxS rows yield
   both sides) → fetch the plain diff synchronously (the same diff command
   the view renders, plain colorArg / no ext-diff, existing builders) →
   parse and locate each identity with the same quirk-free scan the staging
   path uses (`LineNumberOfLine`/`OldLineNumberOfLine`) → copy those raw
   lines, trailing `\n`, with the staging view's `dropDiffPrefix` semantics
   (always applicable now — it's always a raw diff). Sub-decisions at
   implementation: copy the contiguous raw span between the first and last
   matched identity vs. only the matched rows (matters for renderer-hidden
   lines, e.g. difftastic's whitespace-only lines, and reordered rows);
   header rows; unresolvable rows (renderer decoration) at the selection
   edges. e2e: `copy_from_main_view` (rewrite for the new semantics; add a
   case where rendered text ≠ raw text via a fake renderer). Ref: 99f14162c
   (superseded in approach; its row-mapping survives as the first step).
10. **`narrowSelectionHighlight` per-renderer config** — gocui
    `SelectedLineBgColorWidth` (left N columns only), gui maps bool→2;
    docs via `just generate`. Post-#5870 this is a field on
    `DiffRendererConfig`, read via the config manager
    (`GetNarrowSelectionHighlight`) — the rebased prototype shows the
    shape. Ref: cc90accde, N§21.34.

Open item to resolve with the user during this PR: whether `n`/`N`/`f` get
proper keybinding config entries (prototype used hardcoded literals,
N§16.2) — lean: add config entries, matching lazygit convention.

Note: `space` is deliberately **not** bound here — staging arrives in PR 7.
Under a non-conforming restructuring renderer, nav/e simply no-op until
PR 7's raw fallback lands; acceptable interim (same release).

#### Deviations from the plan (2026-08-10, as implemented)

Landed as 9 commits on branch `select-diff-lines-in-main-view` (off PR 4):
the `ViewSelectionController` fold; the classifier; the selection; the
hunk-selection extraction; the hunk-mode default; drag-to-range; hunk/file
navigation; edit; `narrowSelectionHighlight`. Two commits of the plan's ten
are **not in this PR** (1 and 2 below); everything else landed.

1. **Commit 7 (jump-to-file menu) is skipped** — the user's call at the start
   of the session: the UX isn't decided, and it may or may not be added later.
   So `FilesInDiff` is not ported either. If it does arrive, it needs the same
   ReadToEnd-then-retry the nav has (its file list must cover the whole diff).
   **It arrived on 2026-09-13 as PR 12**, reading the diff to the end up front
   rather than retrying.
2. **Commit 9 (copy) moves to PR 7** (decided with the user). The prototype's
   reason for making copy a direct `MainViewController` command rather than a
   `FocusedMainViewActions` method (N§21.28) was that copy is
   panel-*independent* — it only read the rendered text — and that routing it
   through the interface would exclude the reflog. The mandatory raw-diff
   semantics invert both: copy now needs the raw diff of *what this pane
   shows*, which is per-panel knowledge (`WorktreeFileDiffCmdObj` vs
   `ShowFileDiffCmdObj` vs `ShowCmdObj` vs `ShowStashEntryCmdObj`, plus which
   side a pane shows), and reflog is `PatchBuilding` from day one so it has
   actions. Keeping copy in PR 5 would mean a second per-panel seam (a
   one-method interface on six contexts, plus a `plain` flag on the two
   builders that lack one) whose only consumer is copy and which PR 7's
   actions would sit beside rather than use — each `PrimaryAction` fetches its
   own raw diff (PR 7 commit 4). So copy becomes a `FocusedMainViewActions`
   method, landing in PR 7 right after commit 2 builds that interface. **PR 5's
   title drops "copy"**; §8's copy row moves to PR 7. (It is **PR 6b** as of
   2026-09-13, split back out of PR 7 as its own releasable feature; it is a
   `FocusedMainViewDiffSource` method, the narrower of the two interfaces.) The sub-decision already
   taken with the user, for whenever it lands: copy the **contiguous raw span**
   between the first and last matched identity (per file, concatenated in
   display order), so renderer-hidden lines come along and the clipboard is a
   valid patch fragment; the seam takes a **path list** so a whole commit's
   diff is never fetched to copy three lines of it.
3. **`IsSingleHunkForWholeFile` is derived from the rendered diff, not from a
   raw-diff fetch** (deviates from the option agreed at the start of the
   session, in the direction of the performance worry the user raised there).
   The predicate needs no git call at all: over the *anchor's file*, no context
   rows and all change rows of one kind is the same question
   `patch.Patch.IsSingleHunkForWholeFile` asks, and the already-resolved rows
   answer it. Two consequences: it needs no "is this diff a single file?"
   answer from the side panel (it is per file, so it is also right over a whole
   commit's diff, where one file may be new and another edited), and it needs
   no per-panel seam — so PR 5 introduces none. Guarded by
   `ViewBufferManager.IsLoading()`: while the diff is still being read the rows
   that would answer otherwise may not have arrived, so it says false, erring
   towards the hunk mode the user configured.
4. **`UserEnabledHunkMode` is not ported.** It is written three times and never
   read in the prototype; its only purpose was deciding whether escape leaves
   hunk mode, and the escape machinery is never built (§2.11).
5. **The render-side selection-visibility rule lives at the one chokepoint**,
   not in each panel's render-to-main. The prototype called
   `updateFocusedMainViewSelectionVisibility` from `FilesController`'s
   render-to-main and assumed "commit panels always render a diff", which the
   stash panel's "No stash entries" placeholder contradicts. Production derives
   it in `gui.refreshMainViews` from three things it has there: the panel
   beneath is a `DiffMainViewContext`, the pane holds focus, and the pane's
   task is a command task (a diff) rather than a rendered string (a
   placeholder). One place, every panel, present and future. The two moments
   necessarily use different signals — the task at render time, the rendered
   content at focus time (`ViewHasChangeLines`, which additionally catches a
   diff with nothing in it: a binary file, an empty commit) — and each says so.
6. **Commit 4 (selection visibility) is folded into commit 3.** A commit that
   paints a selection over "No changed files" and a later one that stops it
   would be a regression introduced and repaired inside the branch. Commit 3
   also split the other way: the hunk-mode default became its own commit, as it
   did in the prototype, plus a two-line prep extraction so the focus path can
   select a block.
7. **`n`/`N` get config entries** — open question 2 resolved with the user:
   `keybinding.main.prevFile`/`nextFile`, defaulting to `N`/`n`. (gocui
   pre-empts `n`/`N` while a search is active in the view, so search-next still
   wins there; that is master's behaviour for every view.)
8. **The ReadToEnd-then-retry sits in the shared `navigate` helper** from
   commit 3 rather than arriving with commit 6, so hunk-mode ↑/↓ has it as soon
   as it exists. Only *forward* navigation retries: everything above the anchor
   has loaded, so a backward target that wasn't found doesn't exist.
9. **The click-to-dive plumbing is deleted here**, not in PR 9:
   `GetOnClickFocusedMainView`, `AddOnClickFocusedMainViewFn`, the base-context
   field, and the files/commit-files implementations. Clicking the focused main
   view now selects the clicked line (§2.2), which was that mechanism's only
   caller, so leaving it in would be dead plumbing on master for four PRs.
10. **`IContextMgr.IsInStack` is added here** (the plan mentions it only under
    PR 7 commit 9): `GetKeybindings` runs for off-stack panes at startup and
    during cheatsheet generation, and `NextInStack` panics there.
11. **The bool→width mapping for `narrowSelectionHighlight` has one home**,
    `gui.applyDiffRendererSelectionStyle`, called from
    `configureViewProperties` (startup, config reload) and `refreshMainViews`
    (so a renderer cycle picks it up). The prototype duplicated the magic 2 in
    `views.go` and `global_controller.go` despite a comment claiming otherwise.
12. **e2e coverage** is a new `pkg/integration/tests/main_view/` directory, 9
    tests: `select_diff_lines`, `range_select_diff_lines`, `select_hunk_in_diff`,
    `select_hunk_on_focusing_main_view`,
    `select_line_when_whole_file_is_one_hunk`, `drag_selects_diff_line_range`,
    `navigate_by_hunk_and_file`, `edit_selected_diff_line`,
    `no_selection_when_no_changes`, `hide_selection_when_changes_vanish`. Plus
    `SelectionIsShown`/`SelectionIsHidden` on the view driver (reading the
    highlight flags, since `SelectedLines` says nothing about whether a
    selection is drawn) and unit tests for `changeBlockStart` / `fileStart`.
    That leaves PR 7 commit 11 with only the staging-specific ports.

#### Review round 1 (2026-08-11) — five fixes, all as fixup!/amend! commits

The user's first pass found four bugs; a fifth defect turned up while fixing
them. All are fixups on their targets, so the final history has none of them.

1. **Toggling hunk mode below the last change did nothing.** `ChangeBlockBounds`
   only looked ahead of the anchor. The staging view's `GetNextChangeIdx` falls
   back to the change *behind* the cursor, so ours does too. This closes the §8
   row that was marked "fix cheaply if trivial; else defer".
2. **Escape left the view instead of giving up the selection.** Mirrors the
   staging view now: a range collapses to its cursor line, then hunk mode goes
   back to line-by-line when the user turned it on, and only then does escape
   leave — with `DismissRangeSelect` / `SelectLineByLine` as the shown command.
   **This resurrects `UserEnabledHunkMode`** (deviation 4 above is void): it is
   exactly what decides whether hunk mode is something to escape from, and the
   staging view distinguishes the two cases.
3. **Jumping by hunk or file stretched a shift-held range to the target.** The
   arrow keys collapse a non-sticky range; a hunk or file jump is the same kind
   of move, so it now collapses too (a sticky range still stretches). Worth
   recording: **master's staging view does *not* do this** — verified with a
   throwaway integration test, `<right>` there extends the range to the next
   hunk. So this is main-view-only behaviour for now, and the staging view has
   the same wart until PR 9 deletes it.
4. **Dragging a range didn't autoscroll at the view's edge.** The plan missed
   this because master gained it for the staging view *after* the prototype was
   written (b682fb7635df); `helpers.DragAutoscroller` is reusable and is what
   the main view now drives from its drag handler, plus a `MouseRelease`
   binding and a focus-lost hook that cancels the autoscroll and the mouse
   capture. One thing the staging view never needed: its content is a string
   that is always there in full, while this diff loads lazily, so scrolling
   down has to keep reading it in or the autoscroll stops at the loaded edge.
5. **Found while fixing the above: `ReadLines` takes an absolute line total**,
   not a count to add (see `readLinesToFillView` and `layout.go`), so
   `moveCursor`'s `ReadLines(delta)` — transcribed from the prototype, which
   has the same bug — was a no-op: moving the selection down never pulled more
   of a lazily-loaded diff in. `ReadLinesToFillView` is the right call.

#### Review round 2 (2026-08-11)

1. **`narrowSelectionHighlight` is gone — the narrow bar is unconditional in
   diff main views** (decided with the user; no global `gui` option either).
   There is no rendering of a diff that a full-width highlight doesn't degrade:
   even git's own output puts red and green *text* on the selection's
   background. So there was nothing worth configuring, and dropping it also
   means the option never ships (it existed only on this unpushed branch) and
   `applyDiffRendererSelectionStyle` and its per-render call go away — the
   width is now a fixed property of the two views, set where they are created.
   The user tried limiting only the *background* (leaving the bold/bright
   foreground full-width, as suggested) and preferred the original: the bar
   governs both, so the gocui field is now
   **`SelectedLineColorWidth`**, not `SelectedLineBgColorWidth`. Landed as an
   `amend!` on the narrow-highlight commit, whose message no longer describes a
   config option.
2. **The selection could end up past the content and so invisible** (user,
   testing): select the last line of a diff, then switch to a renderer that
   renders the same diff in fewer lines. It is clamped at **end of input**
   (`getManager`'s `onEndOfInput`, beside master's origin clamp), *not* in the
   layout: the layout would have to consult the loading flag, which is cleared
   on the task goroutine after the final paint is already queued, so the last
   layout pass can still see the render as loading and skip the clamp for good —
   reproduced as a flaky test before moving it. Every main-view content change
   comes from a task, so end of input covers renderer switches, context-size
   changes and refreshes alike. PR 6's restore will usually preserve the
   selection properly; this is the floor beneath it for when it can't.

   Noted while writing that test: a stdin-filter renderer that **exits before
   consuming the diff** (`head -4`) leaves the task without EOF, so
   `IsLoading()` stays true forever — which already disables master's origin
   clamp and the scrollbar tracking, and now this clamp too. Pre-existing and
   out of scope; the test uses a filter that reads all of its input.

#### Review round 3 (2026-08-15)

The user merged PR 1 to master and rebased the stack; this round is a read of
the commits themselves rather than of the behaviour.

1. **The selection commit was too large and is now four**: stop diving into a
   patch explorer on click (the two `GetOnClickFocusedMainView` implementations
   plus the binding that had nothing left to call), remove the plumbing now
   that it is unused, take the focused main view's context by its concrete type
   in `focusMainView` (prep), and then show the selection. The first two leave
   a click in the main view doing nothing for two commits, which the user
   accepted. The split was verified content-identical: the four commits'
   combined tree equals the original commit's, and the branch tip was unchanged
   by the operation.
2. **No code comment may explain itself by pointing at the staging view** — it
   goes away in PR 9 and the comments would outlive it. Commit messages may
   still refer to it. Applied across the branch (three commits' worth).
3. **Commands that act on a diff selection are described only where they
   apply**: `DescriptionFunc` returns "" for a non-diff main view, which keeps
   them out of the keybindings menu for a branch log or the status dashboard,
   and `GetDisabledReason` (`Tr.NothingToSelectInDiff`) strikes them through
   when the view holds a diff with nothing selectable. The **static
   `Description` stays** — `pkg/cheatsheet/generate.go` reads that field, not
   `GetDescription()`, so dropping it silently removed keys from the generated
   cheatsheets. Dynamic text belongs in `DescriptionFunc` rather than being
   baked into the binding at registration time: `GetKeybindings` is *not*
   called when a key is pressed, so the bindings themselves must be static
   (the user's point; the options bar calling it per frame is not something to
   rely on). e2e: `selection_commands_only_where_they_apply`.
4. **`escape` calls `resetDiffSelectMode`** instead of repeating its three
   assignments.
5. **`diff_line_navigation.go` is now `diff_line_queries.go`**, holding everything
   that answers a question about a rendered diff (`FirstChangeLineInView`,
   `ViewHasChangeLines`, `IsChangeLine`, `IsSingleHunkForWholeFile`,
   `ChangeBlockBounds`, both `Adjacent*`, the two projections and the pure
   arithmetic), while `diff_line_helper.go` keeps only the recovery of a row's
   identity. "Queries" over "navigation" (too narrow: most of the file isn't
   navigation) and over "scanning" (describes the implementation, not what a
   caller gets). The muddiness the user pointed at was —
   `ChangeBlockBounds` sits in the former while `changeBlockStart` sits in the
   latter. The line that would make sense is identity resolution (the helper,
   `GetDiffLineInfo`, `resolveDiffLines`, and the two projections) versus
   questions asked of a whole rendered diff (`FirstChangeLineInView`,
   `ViewHasChangeLines`, `IsChangeLine`, `ChangeBlockBounds`,
   `IsSingleHunkForWholeFile`, the two `Adjacent*`, and the pure arithmetic).

Everything but the split landed as `fixup!` commits inserted **mid-branch**,
next to their targets, per the user's suggested technique; the branch was
replayed onto each in turn. One thing could not be a fixup: the edit command's
binding had to take its new shape *inside* its own commit, because resolving
the replay conflict there would otherwise have re-introduced the locals the
keybinding fixup had just removed.

#### Interactive sign-off (2026-08-15) — approved

The §6 pass is done and everything behaves as expected: selection feel under
delta, hunk-on-click, drag including the autoscroll, and repeated `n` across
files under a metadata-emitting delta. A few special cases may deserve a
refinement of their exact behaviour later; the user's call is that getting the
later PRs written matters more, so none of them is being touched now.

#### Review round 4 (2026-09-05) — the search follows the selection

Found while testing the branch below the stack (§10's 2026-09-05 entry). List
views and the staging view keep the current search match in step with their
selection, so that moving around and then pressing `n` goes to the next match
from where you are; commit `65edd99fd09` added that for both. The focused main
view never gained it, so `n` carried on from the match it was last on, and
walking down past a few matches then resuming the search took as many presses
as there were matches behind you. The user's call was to add it here, PR 5
being where the main view gets a selection at all.

Two commits at the tip of the branch:

1. **Move the focused main view's selection through one function.** `moveCursor`
   and `selectAbsoluteLine` set the view's cursor themselves, repeating the
   clamp `showSelectionAtLine` does around it; both now go through it. Every
   other selection move already did, so this leaves one place for the next
   commit to hook.
2. **Follow the selection with the search in the focused main view**, by calling
   `SetNearestSearchPosition` there. It needs no `inOnSearchSelect` guard of the
   kind `ListContextTrait` and `PatchExplorerContext` carry. The search's own
   selection goes through gocui's `SelectSearchResult`, and
   `MainContext.OnSearchSelect` only collapses the select mode; it never moves
   the cursor, so this function is not re-entered. A drag is the one selection
   move gocui makes for itself, so `onDragRelease` asks for it too, once the
   gesture has settled. e2e: `search_follows_the_selection`.

#### Review round 5 (2026-09-05) — a render that never reaches its end

Reported: focus the main view over a file's diff so that a selection is showing,
go to the branches panel, select a branch, and press `0`. The main view shows
the branch's commit log with the selection still drawn over it, and pressing
down makes it disappear.

`updateDiffPaneDecorations` was asked once a pane's content was final — as a
string is rendered, and at end of input. A command's output often has no end.
The first read stops at `linesToReadFromCmdTask`'s cap and the rest is read as
the user scrolls, so a log or a diff longer than that never reaches
`onEndOfInput`. `HasSelectableContent` then kept whatever the previous render
left it, invisible while the pane was off the stack and drawn again the moment
`0` put it back. (Pressing down in hunk mode reads to the end looking for the
next change block, so the question got settled and the selection went away.)

Two `fixup!` commits on "Show a selection in the focused main view":

1. **Settle the decorations as the content is revealed, too.** The manager's
   swap-in callback asks as well as end of input, and the function takes
   `contentIsComplete`, because a render still being read answers only in the
   positive. A change line among the lines read so far settles the question,
   while finding none may only mean the changes are in the part still to come.
   Under a panel that shows no diff there is nothing to select whatever the
   content turns out to be. That is the reported case, and the mirror case is
   fixed with it. Coming from a panel with no diff to a diff too long to be read
   at once used to leave the pane claiming nothing to select, so a renderer that
   states its lines' identities placed a selection that was never drawn.
   e2e: `no_selection_over_a_commit_log`.
2. **An e2e test for a diff the pane holds only part of**
   (`select_in_a_diff_read_in_part`), the behavioural demonstration of PR 2's
   round 1: it takes both that round's parser change and this round's first
   commit for a selection to appear over a truncated diff, and this is the
   commit where a selection exists to test.

The gutter marks now come with the diff rather than only at end of input, since
they are settled in the same place.

#### Review round 6 (2026-09-06) — one screenful is not enough to answer with

Round 5 was tested and found short. Select a commit with nothing to select (a
merge commit) in the sub-commits panel, then one with a big diff, and pressing
`0` still shows no selection.

The paint round 5 hooked reveals `height + 10` lines, and a commit's diff opens
with a diffstat: `git show --stat -p` of a 274-file commit puts the first
`diff --git` on line 285 and the first change line on 293, of 24444. So the
question was put over the commit header and the stat, answered "nothing to
select" — which for a render still going means "not yet", so the pane kept the
merge commit's answer. The rest of the initial read does reach the changes
(`Total` is `height*(height-1)`, about 2450 lines on a 50-row terminal), but
nothing asked again, and a 24444-line diff never reaches end of input.

One `fixup!` on "Show a selection in the focused main view": ask with every
batch the read delivers, in the manager's content-only refresh callback. Since
the content of a render only grows, a pane that has already found something to
select is not asked again until the next render, so the repeated question costs
nothing in the common case; reading the whole of that 24444-line diff back takes
6ms, so it is affordable in the uncommon one. `renderContentOnly` went with it,
its one caller now needing a body. e2e: `select_below_a_long_diffstat`.

The measurements behind this, taken on the reported commit:

| lines loaded | resolvable rows | change lines | last section, strict | lenient |
|---|---|---|---|---|
| 110 (≈ first paint) | 0 | 0 | no diff yet | no diff yet |
| 300 | 16 | 2 | REFUSED | well-formed |
| 2450 (≈ initial read) | 2166 | 1264 | REFUSED | well-formed |

They also say what PR 2's round 1 is worth here: at 300 lines its leniency is
the only reason anything resolves at all, while at 2450 the complete sections
above the last would have supplied change lines anyway.

Two things this round explains rather than changes. Whether a selection is drawn
never depended on the diff parsing, so before any of these fixes a long diff
showed one whenever the previously selected item left the flag set — the same
staleness as round 5's report, in the other direction; what the parse failure
showed instead was a selection left wherever the cursor was, with `space` and
`d` doing nothing on it. And focusing at the top of such a commit now lands the
selection on a diffstat row, since focusing never moves the view and no change
line is on screen (the `MiddleVisibleLineIdx` fallback). Consistent, but a
selection there can't be acted on; raised with the user, §8.

#### Review round 7 (2026-09-06) — read on rather than guess

Round 6 was tested with a commit carrying a 10000-line message. Coming from a
normal commit it showed a selection, coming from a merge commit it showed none,
and scrolling to the end turned one on. The user's reading of it: taking the
answer over from the commit before never makes sense, and while the pane can't
tell it should keep reading until it can, to the end of the diff if that is what
it takes. Agreed, and the round does both.

One `fixup!` on "Show a selection in the focused main view", with two
preparations ahead of the commit it lands in:

1. **Guard `View.LinesHeight` against a concurrent write.** It reads the buffer
   a rendering task appends to, without the write mutex. Nothing called it, so
   nothing had tripped over it; the fixup does, from the UI thread.
2. **Separate reading a fixed number of lines from reading to the end.**
   `ReadToEnd` holds a gocui task while it reads so that lazygit doesn't count
   as idle in the meantime, which has nothing to do with reading to the end.
   `readHoldingATask` is that part on its own.
3. **The fixup itself.** `ReadLinesAndWait` asks for another render's worth each
   time the pane still can't tell, so the reading stops soon after the first
   change line and runs to the end only for a diff that has none. It has to hold
   a task, or lazygit counts as idle between the batches and the harness asserts
   into the gap — which is how the first cut of this failed its own test. And
   `MainContext` now remembers which render its answer describes, so a render
   starts from no answer rather than one about other content. A re-render of the
   same content keeps its answer, and with it the selection drawn over it; that
   is what stops a background refresh of a big commit's diff blinking the
   selection off and on. e2e: `select_below_a_long_commit_message`.

The user's call on where the selection lands over a diffstat (§8, round 6):
**keep it**. Reaching the first hunk from there is one press of `a` or `right`,
which is preferable to the view scrolling on its own.

### PR 6 — Keep your position in the diff when changing context size, ignoring whitespace, or switching diff renderers

The `RenderRestore` mechanism plus its three standalone consumers. After this
PR: `{`/`}` (context size), `ctrl+w` (ignore whitespace) and `|`/`\` (renderer
cycle) keep your scroll position and selection instead of jumping to the top.
(Title drafted with the third consumer added 2026-08-15; §9.5 — the user
finalizes the wording at PR-open time — may well shorten it.)

Commits:

1. **tasks: the `RenderRestore` mechanism** — `RenderRestore{FirstPaintReady,
   Apply(swapIn)}` on `ViewBufferManager`; the read loop consults
   `FirstPaintReady()` per line (instead of the count) when a restore is set;
   **`Apply` owns the swap: resolve the target against the off-screen buffer
   first, then `swapIn()`, then set origin/selection** — this ordering is a
   real invariant (reordering reintroduces flicker for buffer-parse; guarded
   by `TestNewCmdTaskRestore`, N§20.5); the origin reset is now the manager's
   `newContentPending` flag rather than a per-render field, and a restore
   takes precedence over it *and* clears it (PR 1 deviation 6);
   **not cleared when a task starts** (survives
   stop-and-replace by the periodic refresh), cleared in `Apply` (found or
   not) — N§14.1; `Apply` work that touches gui state hops to `OnUIThread`
   (it runs on the task goroutine, N§21.29 threading fix; origin writes are
   UI-thread-only per §2.8 — the rebased prototype shows the exact shape).
   Refs: 2e3a3ae5b (mechanism parts), 3b597a0f2, N§14.1, N§20.5.
2. **gocui: off-screen scan accessors** — `OffscreenDiffLineContents` /
   `OffscreenDiffLineContentsFrom(from)` (incremental — the O(n) load scan),
   `OffscreenLineCount`, `MiddleVisibleLineIdx`. Refs: 792c7a294, 3e5b52b8f,
   dd30c26b1 (gocui half).
3. **The shared restore helper** — `restoreDiffLinePositionOnRerender(view,
   candidates, matcher, place)`: prioritized candidate list (anchor first,
   outward, stopping at the first change line each side — `nearbyDiffLines`),
   incremental scan resolves per-row backends during load (metadata only —
   buffer-parse can't parse a partial diff, N§14.1/N§20.3), fallback
   candidates resolved at the EOF swap; `matchByPatchLine` matcher;
   `installDiffLineRestore`. **Mandatory: match against every record on a
   row**, not against the row's one resolved identity (PR 4 deviation 2) — a
   target captured under a unified rendering is an `a`, and the side-by-side
   row that shows it leads with the `d` of the same modification, so a
   single-identity match makes commit 5's renderer switch silently fail to
   find its line. Refs: 506c6ea81, 24a95e965 (amend! final shape),
   0cd3a5886 (`installDiffLineRestore` extraction), N§16.1.
4. **Preserve position across `-U` context-size changes** — anchor =
   selection if shown else middle visible line; offset-preserving placement
   (same screen row); visibility guard (don't install on a hidden Normal
   view — merge-conflict edge, N§16.1). e2e: extend
   `staging/diff_context_change`-adjacent coverage. Ref: 24a95e965.
5. **Preserve position when switching diff renderers** — same one-liner in
   the renderer-cycle handler (prototype: `onDiffRenderersChanged` in
   `global_controller.go`); fixes both the ext-diff top-jump and the
   wrong-line "preserved by raw line number" cases (N§18.2); graceful no-op
   fallback for unresolvable renderers. e2e: `diff/cycle_diff_renderers`
   (renamed by #5870) keeps passing. Ref: a21c5841a.
6. **Preserve position when toggling "ignore whitespace"** (`ctrl+w`) — added
   to the plan 2026-08-15 at the user's suggestion; **no prototype reference**,
   the prototype never touched this path. The call itself is the same one-liner
   as commits 4 and 5, in `ToggleWhitespaceAction.Call` before it re-renders
   (it re-renders through `Context().CurrentSide().HandleFocus`, which walks
   down to the side panel, so the same call covers a focused main view too).
   What is new is the **fallback**, because ignoring whitespace can remove the
   anchor line, its whole hunk, or its whole file from the diff:
   - **The stop-at-the-first-change-line rule in `nearbyDiffLines` is invalid
     here.** It exists because a change line always survives a `-U` change —
     which is exactly the assumption ignoring whitespace breaks. So the
     candidate list has to grow.
   - **Decided with the user 2026-08-15: candidates run outward from the anchor
     over the whole rendering, nearest first, unbounded** — across hunk and
     file boundaries — and the restore lands on the first one the new diff
     still contains. When the anchor's file was whitespace-only you therefore
     land on the nearest surviving line of a neighbouring file, which is still
     where you were in the diff. If nothing survives at all (every change was
     whitespace), no restore: the view renders from the top, and PR 5's
     visibility rules already hide the selection once there are no change lines
     left.
   - Whether the unbounded list *replaces* the short one for commits 4 and 5
     too is an implementation call. It has the same ordering, so it only
     lengthens the list — but a candidate list the size of the buffer makes
     `findResolvedDiffLine`-per-candidate O(n²); invert it (index the new rows
     once, then walk the candidates) before sharing it.
   - Direction is asymmetric: turning ignoring **off** only adds lines and the
     anchor practically always survives (a whitespace-only change shown as a
     context line under ignoring keeps its line number, and `SamePatchLine`
     compares `(line, isDeletion)`, so context↔addition still matches). Turning
     it **on** is the case the fallback is for; a whitespace-only *deletion* is
     the one identity with nothing to match in the new diff, and falls to a
     neighbour.
   - e2e: a main-view companion to `diff/ignore_whitespace` — position kept
     when the anchor survives, nearest-survivor landing when its hunk is
     whitespace-only, and the empty-diff case.
7. **Preserve the selection's far end too** — `selectionFarEndIdentity`
   restored via `SetRangeSelectStart`; collapses to the cursor line when the
   far end didn't survive. Ref: 0412046c4, N§21.32(4).

Known limitation (keep, document in PR): `NormalSecondary` is not preserved
(N§16.1, N§18.3).

Interaction to keep in mind (not this PR's job): master refuses `ctrl+w` in the
staging and patch-building contexts, because a patch built from a
whitespace-ignoring diff doesn't apply. Once staging happens in the main view
the refusal no longer catches it — see §9.10.

#### Deviations from the plan (2026-08-15, as implemented)

Landed as 7 commits on branch `keep-diff-position-on-rerender` (off PR 5): the
`RenderRestore` mechanism; the gocui off-screen accessors; the restore helper
together with the `-U` consumer; the renderer-cycle consumer; the whitespace
consumer; a prep extraction; the far-end preserve. All checks green, every
commit builds and unit-tests clean on its own. §6 sign-off **approved
2026-08-15**, with the whitespace consumer singled out as a welcome addition.

1. **The shared helper landed with its first consumer** (plan commits 3+4 are
   one commit). A commit adding only unexported helpers fails `just lint`
   (golangci-lint's `unused`), and there is no second consumer to justify the
   `installDiffLineRestore` / `restoreDiffLinePositionOnRerender` split yet —
   PR 7's ordinal reveal is the thing that splits them.
2. **`Apply` reports whether it placed the view**, and the task does what it
   would have done without a restore when it didn't. The plan had only the
   other half (a restore beats `newContentPending` and clears it); this is what
   makes "ignoring whitespace emptied the diff" land at the top, as decided.
3. **A pending restore keeps the task reading past the lines asked for, to the
   end of input if need be** — found by a failing e2e test, and the one real
   surprise of this PR. The buffer parser refuses a partially loaded diff
   (`patch.Parse(…).IsWellFormed()` over a truncated last hunk fails the whole
   file section), so a restore over a rendering without OSC records can only
   resolve once the whole thing is in. Both alternatives are wrong: painting at
   the usual point (enough lines to fill the view) makes the restore silently
   miss on any diff longer than the initial read, and the prototype's shape —
   readiness is the restore's alone — stalls, because the task stops reading at
   `LinesToRead.Total` and would sit there showing the previous content until
   the user scrolled. Cost: under a non-conforming renderer a long diff is read
   in full before the re-render appears.
4. **The loading placeholder is suppressed while a restore is pending**, for
   PR 1's reason: blanking the view for a message and then putting the user
   back where they were is the flicker the restore exists to avoid.
5. **`nearbyDiffLines` never had the stop-at-the-first-change-line rule.** The
   whitespace consumer needs the unbounded walk (2026-08-15 decision), so it
   was built that way in the commit that introduces it rather than relaxed in
   the commit that needs it — the `-U` case behaves identically either way,
   since a change line always survives it.
6. **Each candidate goes back on its own screen row**, not on the anchor's, so
   that the content around the line we land on doesn't move at all; a candidate
   from off screen clamps to the top or bottom edge. `screenRows` builds that
   map over the viewport only (O(height)) because `ViewLineForBufferLine` is a
   linear scan, and per candidate that would be O(n²) over a whole-diff
   candidate list.
7. **Matching is by a normalized `patchLine` key** — `(path, kind, line,
   isDeletion)` with every kind of *content* line collapsed together — looked
   up in a map built once per re-render, which is what keeps the whole-diff
   candidate walk linear. All records on a row are indexed, per PR 4 deviation
   2. Two consequences worth knowing: an addition and the context line it turns
   into when whitespace stops counting are the same key, which is what makes
   the whitespace consumer land on it; and a **hunk header is not the same key
   across a `-U` change** (it names the lines it covers), so a restore anchored
   on one falls back to the line below — the `-U` e2e test documents that.
8. **No new `pkg/gui/types` API.** The prototype's `SamePatchLine` /
   `PatchSelectLine` aren't needed: `patchLineOf` normalizes into a comparable
   struct inside the helpers package.
9. **`GetOrCreateViewBufferManagerForView` isn't needed** — the prototype wanted
   it for a secondary pane that hadn't rendered yet; all three consumers here
   preserve a view that has.
10. **The worktree path is captured at install time**, like the view height:
    resolving a record's path uses the repo's paths, a repo switch replaces
    them, and the incremental search runs on the task's goroutine (§2.8).
11. **e2e**: 5 tests in `pkg/integration/tests/main_view/` —
    `keep_position_when_changing_context_size` (both anchor cases, incl. the
    hunk-header fallback), `keep_position_when_switching_diff_renderers`,
    `keep_position_when_ignoring_whitespace` (incl. a change becoming a context
    line), `keep_position_when_ignoring_whitespace_removes_it` (hunk gone →
    nearest survivor; diff gone → nothing to keep), and
    `keep_selected_range_when_changing_context_size` (both ends kept; either end
    dropped). Plus 3 unit tests in `pkg/tasks`.
12. Deviations 3 and 10 arrived as `fixup!` commits, on the mechanism commit and
    on the helper commit; the user has folded them in.

### PR 6b — Copy the selected diff lines from the focused main view

Split out of PR 7 on 2026-09-13, at the user's suggestion: PR 7's first two
commits are a self-contained feature that PR 7 then builds on, so branch
`copy-diff-lines-from-main-view` was created at the copy commit and the rest of
the stack replayed onto it. Nothing was rewritten to do it. It sits between
`show-staged-changes-in-lower-pane` and `stage-changes-in-main-view` (`PlainDiff`
asks the pane which side it shows, which is the lower-pane branch's rule).

Two commits: "Share how a ref's diff endpoints are derived" (prep, written for
this one) and "Copy the selected diff lines from the focused main view", which
introduces the narrow seam `types.FocusedMainViewDiffSource` (`PlainDiff(pane,
paths)`) and its four panel implementations. `FocusedMainViewActions`, which
embeds it, arrives with PR 7's staging commit. The design rationale for taking
the lines from the diff rather than from the screen is PR 5's plan item 9 and
its deviation 2.

#### Round 1 (2026-09-13) — what a selection of headers and messages copies

Three `amend!`/`fixup!` commits on the copy branch, all from the user's review of
the feature in use. Only change lines were matched, so anything else in a
selection quietly contributed nothing.

1. **Headers count as selected lines.** A hunk header names the first line of
   its hunk in both sources of a row's identity (`LineNumberOfLine` returns
   `hunk.newStart`, and an OSC 1717 `h` record carries the same), so it is
   matched like a line of the file. A file header names no line and the two
   sources disagree about what to call it (the buffer parse says the file's
   first line, a record says nothing), so it is answered for by kind: a
   selection touching any row of it takes the file's whole header. Selecting a
   hunk and its header now copies a patch fragment rather than stripped code,
   since `dropDiffPrefix` keeps the columns as soon as a header is in the text.
   This also restores what master's patch explorer did, its copy being the plain
   patch text of the selected rows verbatim.
2. **The rows above the diff are copied as they stand on screen** (rule (ii) of
   the two offered: a selection spanning the commit message and the diff
   contributes both halves, rather than the message half being dropped).
   `PlainDiffOfSelection` returns the two parts separately, and the +/- column
   is dropped only when everything came out of the diff — a message line may
   begin with a '-' without being a deletion. The rule needs a file on screen to
   be above: with nothing resolvable in the view (a restructuring renderer
   before PR 7's raw fallback) it yields nothing, so a renderer's picture still
   never reaches the clipboard.
3. **Toasts.** A copy says so, as every other `ctrl+o` does (master's patch
   explorer was the one that didn't), and a selection that stands for no line of
   the diff gets an error toast instead of silence. New strings
   `SelectedDiffLinesCopiedToast` and `SelectionNotFoundInDiffToast`.
   e2e: `copy_selected_diff_lines` grew three cases (hunk header, file header,
   commit message and a selection spanning it), and its clipboard helper now
   acknowledges the toast, which the harness requires before the next keypress.
   New test `copy_rows_that_are_no_diff_line` for the error toast, with a
   renderer that tags its content rows and ends with one of its own.

**Rejected in the same round:** copying the view's own lines verbatim when the
view already shows git's plain diff (no renderer, whitespace not ignored). After
1 and 2 the two paths agree almost everywhere — the run rule exists to fill in
what a rendering hid, and a plain rendering hides nothing — while the condition
would have to weigh the renderer type, its git args, the raw fallback, ignoring
whitespace and the custom-patch preview pane, and the same selection would copy
different text depending on config. Known gap left open: a `\ No newline at end
of file` marker at the very edge of a selection is not copied, since it names no
line and so can't be a run's endpoint.

**Replay notes.** The three commits went in at the copy branch's tip and the 80
commits above were replayed with `rebase --onto … --update-refs`. Three commits
above needed adjusting, each in its own place: the `test_list.go` entry (a
generated-file conflict), PR 8's custom-patch case in `copy_selected_diff_lines`
(it now re-enters hunk mode before acting, the preceding case leaving a range
selection behind), and PR 9's config rename, which had to rename the key in the
new test file too. All branch tips build, lint and pass; `just e2e` is green at
PRs 6b, 7, 8, 9 and 12.

### PR 7 — Stage, unstage and discard changes directly from the focused main view

The headline PR. After it: in the files panel's focused main view, `space`
stages/unstages the selected line/range/hunk (multi-file, side-by-side aware),
`d` discards, the split follows the acted-on side, the selection advances to
the next change, commit keys work there, and a non-conforming diff renderer
falls back to a raw diff at focus time so staging always works.

Commits:

1. **Extract `diffSplitState` from the files diff renderer** — prep. Ref:
   4ed8a5a87.
2. **`FocusedMainViewActions` — one dispatch interface** — build directly in
   final shape: side-panel contexts expose `GetFocusedMainViewActions()`
   (nil = non-actionable); methods this PR: `OnClick`, `PrimaryAction`,
   `DiscardSelection` + `DiscardSelectionDisabledReason(mainViewName)`;
   `MainViewController` is a thin dispatcher. Refs: a760f9ef5, 02b08eb73,
   N§21.24(A), N§21.25.
3. **`applyDiffLines`** — prep: split "which diff to read" (`sourceCached`)
   from the `ApplyPatchOpts` (stage / unstage / discard differ). Ref:
   929427400 (build the generalized shape directly).
4. **Stage/unstage the selection** — the core: selected view rows →
   change-line identities (`ChangeLinesInViewRange`; all metadata payloads on
   a row when present — SxS; single resolved record otherwise) → group by
   `info.Path` → per-file patch-line index sets via identity scan
   (`LineNumberOfLine`/`OldLineNumberOfLine` — **never**
   `PatchLineFor*LineNumber`, which mis-resolves hunk-boundary and
   modified-pair cases, N§21.11) → one `Transform`/`ApplyPatch` per file;
   direction from the pane (Normal=stage, NormalSecondary=unstage);
   multi-file and directory diffs supported. Refs: f470d870f, f4d5c79da,
   a187eab63, 3018289e8, N§21.11–21.12.
5. **Stage a fully-selected deleted file as a file deletion** — staging a
   deleted file's entire content must yield `D`, not `MD` (stage the file
   deletion itself, not just the content removal). In the explorer this case
   required deliberately entering a deleted file, so it was rare; in the
   merged view hunk-stepping through a multi-file diff hits it routinely.
   Mandatory (decided 2026-07-19; overrides N§21.13's deferral). Belongs in
   the per-file apply loop / files handler: when every change line of a
   deleted file is selected, stage the file itself instead of applying a
   content patch. e2e: multi-file diff containing a deleted file, stage its
   block → status `D`, not `MD`.
6. **Post-action reveal by change-line ordinal** — capture the selection's
   first line's ordinal among change lines before the op; after the re-render
   select the change line at that ordinal in the target pane (clamped),
   re-expanding in hunk mode; a range collapses to a line first. Rides
   `restoreDiffLinePositionOnRerender` with an ordinal-based place. Refs:
   e98e73382, 0cd3a5886 (final model — skip the two superseded matchers),
   N§21.17.
7. **Focus follows the acted-on side** — unified rule: focus
   `NormalSecondary` iff (unstaging AND post-op split), else `Normal`; the
   handler decides (it owns the split knowledge) and does the reveal/focus
   itself, returning only `error`; selection state copies to the target pane;
   get-or-create the target's buffer manager. **Timing fact this relies on**
   (N§21.14): the SYNC `Refresh({FILES, STAGING})` updates the model
   synchronously, but the main-view re-render is queued — so decide focus +
   install the reveal after the refresh returns, and it rides the queued
   render. e2e: the two cross-pane tests + the four reveal tests from the
   prototype. Refs: b9bbd1955, 498784558, 02b08eb73, N§21.13–21.14.
8. **Discard the selection (`d`)** — files backend: discard-unstaged =
   reverse apply not-cached (confirm prompt), discard-staged = unstage; both
   route through the same `applyDiffLineSelection` path as `space` so
   focus-follow/reveal behave identically (N§21.27 bugs 1+2). e2e:
   `discard_from_main_view`, `discard_from_staged_main_view`. Refs:
   eaec32b2b + fixups.
9. **Commit and find-fixup-base from the focused main view** — gated on
   `DiffMainViewTypeStaging`; gate re-checked per press;
   `IsInStack`-guarded `NextInStack` lookup for cheatsheet generation. Ref:
   4b54223f4.
10. **Raw-diff fallback for non-conforming diff renderers + the handshake
   probe** — the probe (prototype: `ProbePagerEmitsDiffMetadata`; rename per
   §1) runs the renderer on empty input, with the route chosen by the
   entry's `type` (#5870): `stdinFilter` via `NewShell`; `extDiff` via
   git's 7-arg convention on two empty temp files; env `OSC1717=V1`, greps
   for the handshake; verdict cached per renderer signature. `extDiff` with
   empty `command` (git's `diff.external`, formerly
   `useExternalDiffGitConfig`) → always raw when focused. **`rawGit` is
   probed like any other renderer** (resolved 2026-08-07; was "decide at
   implementation"): an entry *with* `args` is run as
   `git diff --no-index <args>` on the same two empty temp files, and its
   handshake read the same way. git announces itself for exactly the
   formats it describes, so asked with the entry's own arguments the
   handshake is a faithful answer — `--color-words` announces, a unified
   diff says nothing at all. Don't special-case git by looking for a
   per-line record instead: that was tried and reverted as over-specified.
   Entries with no `args` skip the probe, being already raw (the fallback
   would be a no-op). That leaves **one rule for every renderer type** —
   raw when focused iff the diff needs records and the renderer won't
   supply them, i.e. `mainViewDiffNeedsMetadata() && !probe`, where "needs
   records" is any custom renderer *or* git with word-diff args. The
   `IsWellFormed`/static-args alternatives are dropped. The probe's cache
   signature must include the `rawGit` args, or editing them in the config
   reuses a stale verdict;
   `DiffMainViewShouldRenderRaw` read by every diff panel's render-to-main;
   `ignoreExternalDiff` threaded through the diff-cmd builders
   (`--no-ext-diff`, keep color); `types.NewMainViewDiffTask` routes raw
   renders through `RunCommandTask` (bypasses `GIT_PAGER`); focus flow
   installs a restore to place the selection after the raw re-render;
   click-to-focus replays the clicked view-line index (best effort). e2e
   (rename per §1 terminology): `stage_from_main_view_with_unsupported_pager`, the
   `build`-variant comes with PR 8,
   `stage_from_main_view_with_conforming_pager` (fake handshake renderer).
   Refs: 98881fc9e, 17cfd567e, bf18778e9; the probe detection is N§21.30
   (the observe mechanism never lands — §3).
11. **Port the remaining prototype staging e2e tests** (whichever aren't
    already in earlier commits): `stage_hunk/range/range_spanning_files…`,
    `select_hunk_on_focusing_main_view`, `select_next_*`,
    `advance_to_next_hunk_after_staging_shifts_line_numbers`,
    `focus_follows…`/`focus_returns…`, `no_selection…`/`hide_selection…` (if
    deferred from PR 5).

Design seam to keep (separate-lists input, §7): the focus-follow decision and
the "which side does this pane show" logic must stay **localized** (the
handler + `diffSplitState`), not smeared across call sites — the parked
separate-lists design will want to re-derive "side" from list-section
membership and may want a different focus-follow rule.

#### Review round 1 (2026-08-16) — the branch rebuilt

The user's first pass, from testing and from reading. What came out of it:

- **The staged side moved to the lower pane for good** (its own branch below
  this one, see the section after this). PR 7 was rebuilt on top of it: the
  `diffSplitState` prep commit is gone, which side a pane shows is a property
  of the pane, and the focus-follow rule is symmetric — follow the lines into
  the pane that is left when the acted-on one goes away.
- **The actions live in their own class**, `WorkingTreeDiffActions`
  (`working_tree_diff_actions.go`), rather than swelling `FilesController`,
  where `PrimaryAction` and `DiscardSelection` also read as if they acted on
  the files list.
- **The seam takes panes, not views.** `types.DiffPaneContext` (a `Context`
  with a `DiffSelectState`) is what the interface methods take, which needed
  the select-state types moved from `gui/context` to `types` — a prep commit
  at the bottom of PR 7. The view→context lookup is gone.
- **Acting on a whole file covers both directions**: unstaging every line of
  an added file leaves it untracked again, where before it left an empty file
  in the index.
- **The pane the work moves to shows no selection until the restore places
  one**, so that the selection it was left with the last time it was used
  doesn't flash first.
- **Selection visibility is decided from the content, when the content is
  final** — a fixup for PR 5. The old rule went by the task type at render
  time, which showed a selection over a diff with nothing in it (a binary
  file, before and after a refresh) and hid one over the custom patch, whose
  pane renders a string. Two new tests cover both.
- Smaller ones: the options-bar rule is its own commit; `set.Set` where a set
  was meant; `workingTreeActionDescription` introduced with its first user;
  `DiffModeRendered`/`DiffModeRaw`/`DiffModePlain`; probe files in lazygit's
  temp dir; no patch-explorer mention in a code comment.

#### A branch below it: the staged side always in the lower pane (2026-08-16)

Branch `show-staged-changes-in-lower-pane`, off PR 6 and below PR 7, 2 commits,
green. Raised by the user while reviewing PR 7: which pane a side of a file's
diff appeared in depended on what else the file had — the staged side had the
lower pane while there were unstaged changes above it, and took over the upper
one when there weren't. Now each side has a pane of its own and is shown when
there is something on it, so a file with nothing unstaged shows its staged
changes in the lower pane alone, which then has the whole section (visually
identical to before; what changes is which pane it is).

- The main section's layout becomes a three-way answer (`types.MainPanes`)
  instead of "split or not", derived from which panes the render has content
  for; a window that is left out of the layout gets no dimensions, which is
  already how a view is hidden.
- Focusing the diff picks the pane holding it, and keeping your position when
  whitespace stops counting now covers both panes.
- What it buys PR 7: `mainShowsStaged` is gone, so which side a pane shows is a
  property of the pane (`showsStagedSide(view)`, no node); staging in the upper
  pane and unstaging in the lower one is a fixed rule rather than one derived
  from the file's status; and the focus-follow rule becomes symmetric — follow
  the lines into the other pane when the acted-on one goes away — instead of
  "the staged side migrates".
- 15 e2e tests across `patch_building`, `submodule`, `conflicts`, `diff` and
  `file` asserted the main view for a staged-only file and now assert the
  secondary one. New test: `file/staged_changes_in_lower_pane`.
- Not touched: the submodules panel renders a staged-only submodule's diff into
  the main pane, the same asymmetry in a place with no second pane to speak of.

#### Deviations from the plan (2026-08-16, as implemented)

Landed as 13 commits on branch `stage-changes-in-main-view`, now off
`show-staged-changes-in-lower-pane`, plus five `fixup!`/`amend!` commits from
the rebase onto it. All checks green; §6 sign-off owed.

1. **Commit 3 (the `applyDiffLines` prep) has no separate existence.** There
   was nothing to generalize — production has no `stageDiffLines` to split —
   so the general shape was written directly, as the plan's own note said to.
   **Commit 1 (`diffSplitState`) is gone too**, dropped in the rebase onto the
   lower-pane branch: with the staged side always in the same pane there is no
   split state to extract.
2. **The interface is two, and the narrow one is what panels register**
   (decided with the user). `FocusedMainViewDiffSource` — `PlainDiff(view,
   paths)`, the diff behind what a panel renders — is what every diff panel
   implements and what a context hands out; `FocusedMainViewActions` embeds it
   and adds `PrimaryAction`/`DiscardSelection`, which the dispatcher gets by
   type-assertion. That way copy works over every diff (reflog included) in
   this PR while only the files panel acts, and no panel carries a stub method
   for a PR. PR 8 makes the commit panels satisfy the extension.
   `DiscardSelectionDisabledReason` isn't in it yet: the files panel has no
   reason to give, so it lands with PR 8's panels, which do.
3. **`OnClick` is not a method of it**, PR 5 having deleted the dive gesture
   and its plumbing.
4. **Copy landed first**, as the seam's first consumer, since a seam with no
   consumer is dead code. Its span rule is per file, from the first to the last
   selected *content* line; headers are not matched, since a file header's
   record and the buffer parse disagree about which line it is, and a run
   between two content lines carries any header between them anyway.
   `pkg/gui/controllers/helpers/diff_line_plain_text.go`. **Both of these
   commits left PR 7 on 2026-09-13 and are PR 6b now**, which is also where the
   header rule was revisited.
5. **The options bar leaves out a command that describes itself as nothing.**
   PR 5 made an empty `DescriptionFunc` mean "doesn't apply here"; without this
   the space key would have shown as a blank entry over a commit's diff, where
   nothing has actions yet.
6. **A whole-file selection stages the file** rather than only the deleted-file
   case (plan commit 5): applying a file's entire diff to the index is what
   `git add` does anyway, so the condition is "the selection covered every
   change", which needs no question about what kind of file it is.
7. **The plan's timing fact for commit 7 is stale, and the focus decision is
   made without the model** (raised with the user, who chose this over asking
   git). Post-rework a refresh only *queues* its model update, and its `Then`
   runs after the panel's render-to-main has already started the render task —
   so the handler can neither read the post-op split nor install a restore
   late. What is left staged is therefore worked out from what we just did:
   `applyDiffLines` reports whether the selection covered all of a file's
   changes, and the files it didn't touch are as the model describes them.
8. **The restore core split out as its own commit**, as PR 6 deviation 1
   anticipated, with the ordinal reveal as its second consumer.
9. **`GetOrCreateViewBufferManagerForView` is needed after all** (PR 6
   deviation 9 said it wasn't): the pane the staged side moves *to* may never
   have rendered, and the restore has to be on the manager its first render
   will use.
10. **Rapid keypresses lost a press** — found by porting the staging view's
    test, raised, and fixed as its own commit. `RefreshBlockingInput` releases
    when the model is up to date, but the selection only moves once the
    asynchronous re-render lands, so the replayed press acted on lines that
    were gone. The action now holds input until the restore *resolves*, for
    which `RenderRestore` gained a `Done` hook — and, so that the hook always
    runs, a view given a string task now drops a pending restore instead of
    leaving it to claim a later render (a hole PR 6 left).
11. **`plain bool` became `git_commands.DiffMode`** (`DiffRendered`,
    `DiffRaw`, `DiffPlain`), in its own prep commit. Two booleans would have
    had an invalid combination, and the mode also settles what the prototype
    got wrong: ignoring whitespace applies to a raw render (it is about what
    the user wants to see) but never to a plain one (a patch must describe
    every change), and a raw render uses git's own colour rather than the
    renderer's preference. It also took the colour decision off each command,
    and dropped a `plain` parameter the patch builder always passed as true.
12. **The probe is `DiffCommands.ProbeDiffRendererEmitsMetadata`**, switching
    on the renderer type as planned, with the verdict cached on
    `DiffLineHelper` (`MainViewDiffMode`, in `diff_line_raw_fallback.go`) — the
    helper that already owns "can a row of this rendering be placed". Panels
    ask it for the mode and pass that to `types.NewMainViewDiffTask`, which
    picks the pty or the plain command task from it.
13. **Three PR 5/6 tests needed their fake renderers to announce the
    protocol.** They are about position, clamping and copying, not about the
    fallback, and their renderers said nothing — so with the fallback in place
    their output was replaced by git's own the moment the view was focused.
14. **§9.10 resolved: `ctrl+w` is left alone** (decided with the user). Master
    refuses it in the staging view because that view renders
    `WorktreeFileDiff(plain)`, which never carries `--ignore-all-space` — the
    toggle simply couldn't show there. In the merged view the *rendering*
    honours it while the patch is still built from a freshly fetched plain
    diff, and identities are `(path, file line number, isDeletion)`, which
    ignoring whitespace doesn't renumber. So a selection made in a
    whitespace-ignoring rendering stages exactly the real change lines.
15. **§9.8 resolved: no real-renderer harness yet** (decided with the user):
    PR 7's renderer tests are about renderers that say nothing, which a fake
    models exactly. Build it when PR 8's difftastic work has a test that needs
    it.
16. **Renames are not threaded through copy's path list** — a renamed file's
    plain diff is fetched by its new path alone, so git shows it as an
    addition. It belongs with §8's rename row, which PR 8 commit 2 owns.
17. **e2e**: 13 new tests in `pkg/integration/tests/main_view/` —
    `copy_selected_diff_lines`, `stage_diff_lines`, `unstage_diff_lines`,
    `stage_range_spanning_files`, `stage_deleted_file`,
    `select_next_change_after_staging`, `focus_follows_staged_side`,
    `focus_returns_when_split_collapses`, `discard_diff_lines`,
    `commit_from_main_view`, `stage_under_unsupported_diff_renderer`,
    `stage_under_conforming_diff_renderer`,
    `advance_after_staging_shifts_line_numbers`,
    `select_next_deletion_after_staging_one`,
    `select_next_change_after_unstaging`, `stage_hunks_with_rapid_keypresses`.

#### Interactive sign-off (2026-08-16) — approved

The §6 pass is done and PR 7 is signed off. It came with four review comments
about the stack as a whole rather than about staging, all addressed the same
day as mid-branch `fixup!`/`amend!` commits (never at the tip — the user's
standing rule, restated here). What was decided and done:

1. **Focusing the main view never scrolls.** Pressing `0` used to select the
   first change at or below the top of the viewport — which could be below the
   viewport too, and in hunk mode the selection then scrolled up to the block's
   first line. Now the search is bounded by the viewport, and nothing about
   focusing (or clicking) scrolls: a click that selects a hunk starting above
   the view leaves the view alone, which is also the workaround for selecting a
   single line inside a hunk. In hunk mode the block offered up is the first one
   that **begins** on screen (the user's call, taken over "the first visible
   change"), falling back to a block that reaches into the view from above,
   whose start is then off screen. With no change on screen at all — a big `-U`
   — it drops to line mode on the **middle visible line** (`MiddleVisibleLineIdx`,
   which is the middle of the content when the content is shorter than the
   screen). That accessor came down from PR 6's "Let a view be read while it is
   re-rendering" into PR 5, where the first consumer now is; that commit's
   message lost its paragraph about it via an `amend!`.
2. **Position preservation covers both panes**, for `{`/`}` and for the
   renderer cycle as it already did for `ctrl+w` — the lower pane holds the
   staged side of a file's diff, and it was jumping to the top. The "either
   pane may hold the diff" rationale now lives in
   `PreserveDiffPositionOnRerender`'s doc rather than at each call site.
   (`ctrl+w`'s second call was only added by the staged-side-in-the-lower-pane
   branch; the fixup moves it back to PR 6, where the case already existed.)
3. **A selection scrolled out of sight no longer drags the view back to it.**
   The line kept in place is now the end of the selection that is **on screen**
   — the cursor, or the range's other end when only that is visible (the user
   asked for "either end" as a refinement over testing the cursor alone) — and
   the middle visible line when the whole selection is off screen. The
   selection is still restored by identity either way, with its own fallback
   walk, so it stays on the same line of the diff; only the view stays put.
4. **Searching collapses the selection to the match** (`MainContext.OnSearchSelect`,
   like the staging view's patch-explorer context does), so `n` after turning
   hunk mode back on mid-search behaves too. Its prep commit moves the
   select-mode reset onto `MainContext`.

Six new e2e tests: `select_visible_change_on_focusing_main_view`,
`select_visible_hunk_on_focusing_main_view`, `search_collapses_the_selection`,
`keep_position_in_both_panes_when_{changing_context_size,switching_diff_renderers,ignoring_whitespace}`,
`keep_position_when_the_selection_is_off_screen`,
`keep_position_by_the_visible_end_of_a_selection`. Worth knowing for later
tests: **while the focused main view holds focus, a renderer that doesn't emit
metadata is bypassed** (PR 7's raw fallback), so a renderer-behaviour test that
focuses the view has to give its renderers the handshake.

#### Review round 3 (2026-08-17) — two defects under a non-conforming renderer

Both found by the user testing with an **unpatched git**
(`PATH=/opt/homebrew/bin:$PATH lazygit`; the git in `PATH` on this machine is a
patched build and *does* announce the protocol for `--color-words`, so the
fallback path can't be reproduced with it). Also corrected here: an earlier
claim that focusing under a non-conforming renderer lost the scroll position —
it didn't, because for a stdin filter the raw render runs the *same command*.

1. **The raw fallback only bypassed stdin filters.**
   `FilesController.renderWorkingTreeDiff` asked for the mode and then built its
   command with `DiffModeRendered` regardless, passing the mode only to
   `NewMainViewDiffTask` (which chooses pty vs plain command). That keeps a stdin
   filter out — git hands its output to one only when it thinks it is talking to a
   terminal — but the other two routes are *in the command's arguments*
   (`AddCommonDiffArgs`): an `extDiff` renderer stayed configured and enabled, and
   a `rawGit` renderer kept the arguments that make its output unreadable. So
   under those two types, focusing left the un-actable rendering in place and the
   next refresh dropped the selection. **Every other panel already passed the mode
   into its command builder; the files panel was the only one that didn't.** The
   fix also gets the raw diff coloured under a `colorArg: never` renderer, the
   colour arg coming from the mode too.
2. **A re-render whose command changes reset the scroll to the top** — `{`/`}` or
   `ctrl+w` under a non-conforming renderer, or cycling between a stdin filter and
   an `extDiff`. The command key changes, so `newContentPending` is set, and
   `PreserveDiffPositionOnRerender` can't rescue it: it remembers lines by
   identity, and a renderer that says nothing about its rows leaves nothing to
   look for. Fixed with **`ViewBufferManager.SetKeepScrollPositionForNextTask`**,
   the coarser sibling of `SetRestoreForNextTask` (the user's suggestion, over a
   flag threaded through `ViewUpdateOpts`): same install-before-the-re-render
   lifetime, consumed by whichever task starts, and it suppresses
   `newContentPending` — so neither the reset nor the loading placeholder applies.
   Both position-preserve and the raw fallback set it; the post-staging reveal
   doesn't. Fixing (1) makes the raw render a different command, so without this
   it would have *introduced* a top-jump on `0` — they land in that order.

New tests: `tasks.TestKeepScrollPositionForNextTask`,
`keep_scroll_when_the_diff_cant_be_read` (`{` under a renderer whose output can't
be parsed), `raw_fallback_under_an_external_diff` (focus brings git's own diff at
the same offset, with a selection). **An `extDiff` renderer is the deterministic
way to test the non-conforming path** — its command is entirely the test's, unlike
`rawGit`, whose answer depends on the git in `PATH`.

Not fixed, raised for later: `renderNonTextualConflict` (the DU/UD conflict hint
plus a `--base` diff) also renders with `DiffModeRendered` hard-coded, so the
same bypass doesn't happen there; whether a selection over that content should
exist at all is the prior question. **Answered 2026-09-10** (PR 9's round 1,
item 3): it shouldn't, and the hard-coded mode is right as it is. The fix is a
commit of its own, "Show a selection only over the diff the panel offers",
sitting right after this round's commit.

#### Rebase mechanics for mid-branch fixups (learned the hard way, 2026-08-17)

- **The user's git config has `rebase.autosquash = true`**, so a plain
  `git rebase -i` *silently folds every existing `fixup!`/`amend!` into its
  target* — which is exactly what AGENTS.md forbids an agent from doing. Always
  `git -c rebase.autosquash=false rebase -i …`. (This happened; the fixups were
  recovered from the branch reflog and re-inserted.)
- The single-rebase recipe the user prefers, in place of one
  `rebase --onto` per fixup: mark the target `edit` in the todo, make the change,
  `git commit --fixup=…`, `git rebase --continue`. Several targets can be marked
  in one pass.
- `rebase.instructionFormat` puts a **`# ` before the subject** in the todo, so a
  sequence-editor script matching on subjects has to strip it (and a trailing
  `# empty`).
- `--update-refs` does **not** move a ref that points at the rebase's exclusive
  start, so `rebase --onto <fixup> <target> <branch>` leaves a PR-branch ref
  pointing at the pre-amend commit — how `keep-diff-position-on-rerender` and
  `show-staged-changes-in-lower-pane` went stale for a session. Both now point at
  the **last fixup for their own tip commit**, so each PR contains the fixups for
  the commits in it.
- (2026-08-18) A **conflict at an `edit` stop consumes the stop**: resolving it
  and running `git rebase --continue` commits the resolution *and* moves on, so
  the intended fixup never gets written. Either commit the fixup before
  continuing, or plan a second rebase for it (which is what happened to the
  `amend!` on "Follow the acted-on lines when their pane goes away").
- (2026-08-18) To insert commits at the **foundation** of the branch, put a bare
  `break` as the first todo line: the rebase stops with `HEAD` at the upstream,
  the new commits go in there, and `--continue` replays the branch on top.

#### Review round 4 (2026-08-18) — scroll preservation and focus handling

Four problems the user found while testing the stack. Three of them share one
root cause: the diff selection lives in the gocui view as **view lines**, while
everything meaningful about it — which line of which file — is a **buffer line**
identity, and each of the three places that boundary is crossed lost something
different.

| Crossing | What was lost |
|---|---|
| row → identity (remembering) | a row showing *two* diff lines was remembered as only one |
| identity → row (placing) | a buffer line the view wrapped was placed as its **first segment** only |
| row → row (re-wrapping) | not crossed at all — the indices were simply kept |

1. **The focus was left in a pane that had gone away.** Only
   `applyDiffLineSelection` handled it, so it worked for stage/unstage/discard
   and for nothing else — committing from the staged pane, or a commit or discard
   happening outside lazygit and arriving with a refresh, left the focus on a
   pane the layout had made invisible, where the next keypress acted on nothing.
   The rule now lives in **`Gui.followFocusIntoWorkablePane`**, called from
   `refreshMainViews` *before* the tasks start (so the pane moved into can be told
   where to put its selection as it renders). It is asked of every render, and
   answers from `RefreshMainOpts` alone — no model read, no waiting for the
   layout. The pane moved into gets `EstablishSelection` from a `RenderRestore`
   with `FirstPaintReady: false`, the same shape `RenderFocusedMainViewAgain`
   uses, and shows no selection until then. A pane that **already has a restore**
   keeps it (`ViewBufferManager.HasRestoreForNextTask`), which is how the staging
   path's smarter ordinal-based reveal still wins; that path's own
   `Context().Push` is gone, folded into this rule by an `amend!`.
   - Under `splitDiff: always` the emptied pane is still *shown*, so "is it
     shown" is the wrong question — the render now says which panes it is giving
     something to act on, via **`ViewUpdateOpts.NothingToActOn`**, set by
     `renderWorkingTreeDiff` from `node.GetHas{Staged,Unstaged}Changes()` before
     the `alwaysSplit` override widens them. Decided with the user: move to the
     pane that has changes; when *nothing* is left anywhere, stay in the focused
     main view (over the "No changed files" placeholder) rather than popping back
     to the files panel.
   - This needed a prep refactor, the one judgement call taken alone:
     `refreshMainViews` is in `pkg/gui` and cannot reach a controller, so
     `establishDiffSelection` and its four helpers **moved onto `DiffLineHelper`**
     (`EstablishSelection`, `SelectChangeBlock`, `ShowSelectionAtLine`, and two
     private ones). They needed nothing a helper doesn't have — every one of them
     was already built on `DiffLineHelper` — so the move is mechanical, but it is
     a new commit in a signed-off PR.
2. **Resizing a wrapping view moved everything in it** — the most common trigger
   being the main pane widening by half a screen when the other pane goes away.
   `SetView` threw the wrapping away on a size change and left `oy`, `cy` and
   `rangeSelectStartY` as raw view lines, so re-wrapping at the new width left
   each of them on a line it was never on. This is a **pre-existing gocui bug**
   (the scroll offset of any wrapping view drifts on a terminal resize on master
   too); the branch is only the first thing with a *selection* over wrapped
   content, so it is the first time it shows. Fixed at the **foundation of the
   branch**, demonstrate-then-fix: the exported `ClearViewLines` becomes
   **`RewrapContent`**, which maps all three positions through the line of content
   they were on. Exact, needing no identity machinery, since only the width
   changed. The cursor keeps the screen row it was drawn on; a view with no cursor
   on screen keeps its own place; a range's ends go to the **outermost** segments
   of their lines, a range being over lines of content rather than the segments
   they are drawn as.
3. **A wrapped line came back selected by its first segment only.** Both ends of a
   restored selection went through `ViewLineForBufferLine`, which gives a buffer
   line's *first* view line. Fixed by spreading the ends outward
   (`coverWholeLines` in `diff_line_restore.go`), on the new
   **`View.LastViewLineForBufferLine`** — which `ChangeBlockBounds` had been
   hand-rolling since PR 5 and now uses too. The user's reasoning generalises and
   is worth keeping: the *semantic* selection is a set of buffer lines (that is
   what `DiffLinesInViewRange` hands to staging and copy), and the view-line range
   is only how it is drawn, so normalising a range to whole buffer lines is always
   right.
4. **A modification drawn on one row came back as half of itself.** Remembering
   used `resolveDiffLines` → the row's *leading* identity, while finding used
   `resolveDiffLineIdentities` → *all* of them. So unified → side-by-side worked
   and the way back did not. Two halves to the fix: a range whose ends land on one
   row **is still a range** (`View.HasRangeSelect`, in place of `first == last`),
   and the far end takes the **outermost** diff line its row shows — the last for
   the range's lower end, the first for its upper one, via the new
   `diffLineIdentitiesAt`. The reported asymmetry ("side-by-side gives the
   deletion, `--word-diff` gives the addition") is not ours: patched git emits
   `a;2` first for a word diff (the deletion arrives as a trailing orphan record),
   delta's side-by-side leads with the left column. `identities[0]` faithfully
   picked whichever came first. Line mode is untouched — the user's call: line
   mode means one line.

Also corrected here, having been overstated when the round was scoped: **a diff
renderer is not generally re-laid-out when the view width changes.** A terminal
resize never re-runs the diff at all (the TODO in `Gui.onResize`), and plain
delta's output is line-for-line identical at any width — only its `───`
decoration rules change length. What remains is narrow enough to defer (see §8).

Test infrastructure, needed before any of this could be tested and fixed first,
at the foundation: **`View.SelectedLine`/`SelectedLines` indexed the buffer with
view-line indices**, so a wrapping view reported the wrong lines — and
`ViewDriver.SelectedLines` goes straight through them. They now report the lines
of content a selection covers, once each. Since that makes the segment extent
invisible to a test, `ViewDriver.SelectedViewLineRange(first, last)` was added
for the cases that are *about* the extent.

Thirteen commits as written, six of them `fixup!`/`amend!` commits since folded:
four at the foundation (two demonstrate-then-fix pairs, both gocui), three new
commits before the staging-focus commit, and refinements to five commits already
in the stack. New e2e tests:
`keep_a_wrapped_line_covered_across_a_rerender`,
`keep_both_halves_of_a_change_selected` (a fake side-by-side renderer, which is
the deterministic way to get a row carrying two records),
`focus_follows_a_pane_emptied_from_outside`,
`focus_leaves_an_always_split_empty_pane`. New gocui unit tests:
`TestSelectedLinesOfWrappedContent`,
`TestResizingAWrappingViewKeepsItsPlaceInTheContent`,
`TestLastViewLineForBufferLine`. Every one of the 61 commits the round left
behind builds and unit-tests clean on its own, and the full e2e suite passes at
each of the thirteen it added (`focus_leaves_an_always_split_empty_pane`'s commit
flaked once in four whole-suite runs and was clean in the other three).

Two smaller things deliberately left, besides the §8 rows:

- **Which half the cursor lands on** coming back from `--word-diff`: the
  selection covers both lines, but the cursor sits on the addition rather than
  the block's first line, because that is the order the records arrive in.
  Getting it right means remembering the selection as a *set* and taking min/max
  in the new rendering — which would also replace the "far end didn't survive →
  collapse to the one line we landed on" behaviour signed off in PR 6.
- `renderNonTextualConflict` still renders with `DiffModeRendered` hard-coded
  (carried over from round 3, unchanged). It stays that way; what changed on
  2026-09-10 is that the pane shows no selection over it (PR 9's round 1).

#### Review round 5 (2026-08-19) — the main section changing hands

Two problems the user found testing `show-staged-changes-in-lower-pane`, both
from the same thing: staging a whole file from the side panel moves the diff
from one pane to the other, and nothing was making that look like the *section*
re-rendering.

1. **The section flashed black.** The pane taking over had been emptied when it
   last went away, so it showed nothing for as long as its own render took —
   with a diff renderer, at least a layout pass, since `newPtyTask` defers the
   task to `afterLayout`. `clearMainView` isn't new and isn't the problem; being
   *revealed* after a clear is. Fixed by handing the section over:
   **`Gui.handOverMainSection`** copies the outgoing pane's content (and offset)
   into the incoming one before the render is triggered, the same
   `View.CopyContent` trick `moveMainContextToTop` uses across the tabs of a
   window. It fires only on `MainPaneOnly` ↔ `SecondaryPaneOnly`; a pane
   *appearing beside* the other one (`BothMainPanes`) has no placeholder to be
   handed and is briefly empty, as it is on master.

2. **The pane taking over kept whatever offset it was left at.** Each pane now
   renders one fixed command, so the manager's task key never changes and
   nothing resets the scroll — where before the branch one view rendered both
   sides and `--staged` made every switch a new key. **This one pre-dates the
   stack**: an emptied pane keeps its offset and its claim to the render it was
   showing on master too, it was just hard to reach (visit another file and come
   back to a split one). Verified with the demonstrating test at the foundation
   of the lower-pane branch, which is where the fix went, as a
   demonstrate-then-fix pair: **`clearMainView` now also resets the origin and
   calls the new `ViewBufferManager.ForgetRenderedContent`**, so whatever the
   pane is given next counts as content the view hasn't seen. That is what puts
   the pane taking the section over at the top; the handover only stops the
   flash.
   - The zeroing that made a first attempt at a test pass without the fix was a
     race, not a mechanism: a task that reaches EOF *after* its view was cleared
     clamps the origin to the (now zero) content height in `onEndOfInput`. A
     journey that empties a pane with no task in flight — selecting a
     staged-only file and an unstaged-only file in turn — fails deterministically
     without the fix.

3. **Knock-on: the focus-follow restore read a scroll position that was about to
   be reset.** With the handover, the pane the focus moves into starts at the
   outgoing pane's offset, and `firstPaint` did `Apply` (which is where
   `EstablishSelection` picks the first change *on screen*) and only then the
   reset the new content was owed. The cursor is stored as a row, so the reset
   took the selection with it: it ended up on view line 0, the `diff --git`
   header, rather than on a change. `firstPaint` now settles the scroll position
   **before** consulting the restore. Every other restore comes with
   `SetKeepScrollPositionForNextTask`, so nothing changes for them.
   - That leaves `RenderRestore.Apply`'s bool with no reader, but the signature
     can't be simplified where the reorder lands: `de78cd107` adds another
     `Apply` literal further up the branch, and adapting it there would leave a
     commit that doesn't compile until its own fixup. So the drop is a separate
     commit at the **tip** ("Stop a render restore saying whether it placed the
     view"), which is the earliest place every literal exists.

Landed as 2 commits at the foundation of `show-staged-changes-in-lower-pane`
(demonstrate + fix), a `fixup!` for its tip commit (the handover), a `fixup!`
for "Follow the focus into the pane that is left" (the paint ordering), a
`fixup!` for "Leave a pane that only the permanent split keeps around" (it was
recomputing the pane set the handover already has), and the tip cleanup. The
lower branch ref points at its own fixup, per the rebase-mechanics note above.
New tests: `file/pane_shown_again_starts_at_the_top`,
`file/pane_taking_over_starts_at_the_top`,
`main_view/focus_follows_into_a_pane_taking_over`,
`tasks.TestForgetRenderedContent`, and `TestNewCmdTaskRestore` now asserts the
reset lands before `Apply`. Every commit in the stack builds; whole suite green.

**No automated test for the flash itself.** It is one transient frame, and any
deliberately slow render trips the 200ms "loading..." indicator before the
placeholder can be observed — the indicator applies now that the pane counts as
showing something new, which is what a single view has always done. **Signed off
interactively 2026-08-19**: the flash is gone and both panes start at the top.

Two more findings came out of the round; both were raised with the user and both
are **fixed**, each in a `fixup!` for the commit that owns the rule:

4. **`scrollUpMain`/`scrollDownMain` scrolled a pane that wasn't there.** They
   pick the view by the window the focus is in — the lower pane while the focus
   is in it, else the upper one — which was the same thing as "the pane the
   section is showing" only while the lower pane never appeared alone. From the
   side panel over a staged-only file, `<pgdown>` scrolled the hidden upper pane
   and nothing moved. Both handlers now ask `Gui.mainSectionView`, which adds the
   `MainPanes == SecondaryPaneOnly` case; the choice had been spelled out twice,
   so it moved into one place first. `file/staged_changes_in_lower_pane` — the
   test that owns this arrangement — grew a scrollable staged diff and asserts
   both scroll keys act on the pane holding the section. Fixup for "Always show a
   file's staged changes in the lower pane".

5. **An emptied pane stranded a pending `RenderRestore`, and the input block with
   it.** The string renders all call `DropRestoreForNextTask` because a view
   being given something other than a re-render has no render for the restore to
   ride; a pane being *emptied* is the same case and was missed. Since round 4 a
   restore's `Done` is also what balances the `BeginBlockingEvents` that holds
   input back until the selection has moved on — so a reveal whose target pane
   turns out to be the empty one (`revealSelectionInPaneItLandsIn` picks the side
   the work will land in *before* the refresh lands) would have blocked input for
   good. `clearMainView` drops it now. Latent rather than reproducible: it needs
   the prediction to be wrong, which takes an external mutation racing the
   action, so there is no test. Fixup for "Hold input back until the selection has
   moved on".
   - Worth knowing: dropping a restore resolves its `Done` synchronously, and
     `EndBlockingEvents` replayed the buffered keys in that call — so this, like
     the string-render path it copies, could dispatch a keypress from inside
     `refreshMainViews`, before `setMainPanes` has run. It is the reason a
     stranded restore can't simply be left to a later render. **Closed by the
     addendum below**, which makes the replay wait for the next pass of the
     event loop.

#### Addendum 2026-09-01 — "Edit hunk" ported here, not in PR 9

Found reviewing PR 9: `E` (`keybinding.main.editSelectHunk`) — hand the hunk
you are on to an editor and apply what comes back — went with the staging
panel and nothing replaced it, while the config key stayed in
`user_config.go`, the schema and Config.md, binding nothing. It was first
fixed inside PR 9's removal commit; the user moved it here, where the rest of
the ported staging commands are: the stack builds the main view up to match
the panel and only then removes it, so a command the main view has to gain
belongs in the build-up.

Three commits at PR 7's end:

- **"Ask a parsed patch which of its lines a selection covers"** — a prep
  refactor. PR 7 predates PR 8's identity API (`patch.LineIdentity` /
  `ChangeLineIndicesForLines`, from "Let the patch builder be told which lines
  by their identity" and "Say which change line is meant in one way"), so at
  this point `applyDiffLines` walks the diff inline. Pulled into a
  `changeLineIndices` helper rather than copied.
- **"Don't refresh while the editor still has the hunk"** — a bug this
  uncovered, present on master. `FilesHelper.EditFileAtLineAndWait` went
  through the suspend path that refreshes as soon as the subprocess exits, so
  it read the repo *before* the caller applied the edited patch and then raced
  the caller's own post-apply refresh to publish it. Last writer wins, and when
  it was the stale one the files panel kept showing the file as it was before
  the edit. Both callers already refresh after applying, so the middle refresh
  only ever had a wrong answer to give; it is gone. Found because the new e2e
  test below was flaky about 1 run in 40 — see §10's log for how it was pinned
  down.
- **"Edit the selected hunk from the focused main view"** — the port.
  `WorkingTreeDiffActions.EditHunk`, from the *git* hunk (context and all)
  around the selection in the file's plain diff: `Patch.HunkContainingLine` /
  `HunkStartIdx` / `HunkEndIdx`, which is what the explorer's
  `State.CurrentHunkBounds` asked. Working-tree only, as it was — what the
  editor hands back is applied to the index, which a commit's diff has no use
  for — and applied whole rather than matched against the file's diff again,
  the point being that it says something the diff didn't. e2e:
  `main_view/edit_hunk_in_focused_diff`, with a shell command standing in for
  the editor, asserting both the line the editor was pointed at and that the
  edited content reaches the index while the working tree keeps its own.

PR 8's **"Say which change line is meant in one way"** then converts that one
helper instead of the inline walk — the same commit, one site rather than
two, and its message says so. The end state is byte-identical to having done
it in PR 9; only the history moved.

Consequence for §6: `E` with a real editor now needs its interactive pass as
part of **PR 7**, whose sign-off predates it.

#### Addendum 2026-09-05 — the key replay waits for the next pass of the loop

Follow-up to round 5's "worth knowing" note above, raised by the user. The
exposure it records is not really about restores. `EndBlockingEvents` dispatched
the buffered keys from inside its own call, so whoever ended a block ran a
keybinding handler in the middle of whatever they were doing. Round 5's fixup is
the first caller for which that matters, `clearMainView` resolving a restore
partway through `refreshMainViews`. The behaviour itself dates from
`6893d9a759`, on master since long before the stack.

What a handler dispatched there could see: `State.MainPanes` still describing the
previous render, so `mainSectionView` scrolls the wrong pane (finding 4 by
another door); and, for any replayed key that moves a side-panel selection, a
nested `refreshMainViews` whose work the outer call then partly writes over.

So the fix belongs to the gocui primitive, not to the restore. **New commit
"Replay withheld keys on a later pass of the event loop"**, placed immediately
before "Hold input back until the selection has moved on", which is the
justification for making the improvement here: the replay goes through `Update`,
and input stays withheld until that queued pass runs. The withholding has to
outlast the counter because `processRemainingEvents` prefers gui events to
queued work, so a key pressed in the meantime would otherwise be handled ahead of
the keys buffered before it. `EndBlockingEvents` no longer returns an error (the
replay's now reaches gocui's error handler like every other handler's), so the
three `func() { _ = ...EndBlockingEvents() }` closures in the stack become the
method value, and `WithWaitingStatusBlockingInput`'s deferred callback moves into
`AppStatusHelper.endBlockingInput` to keep `unparam` happy.

Two options were weighed and dropped. Hopping `DropRestoreForNextTask`'s
`resolved()` onto the UI thread as well would give `Done` one uniform contract,
but it trades a guarantee that holds today (a restore resolves exactly once,
whichever way it ends) for one that doesn't (an `onUIThread` hop can fail once
the loop has exited), for a second `Done` implementation that doesn't exist.
Moving `setMainPanes` up to just after `handOverMainSection` fixes only the stale
`MainPanes` half and leaves the nested-render clobber.

New tests: `TestBlockingEvents_KeysArrivingBeforeTheReplayGoBehindIt`, and the
three existing block-events tests now pump the queued work, which asserts the
deferral. Whole suite green, `just lint` clean, and every rewritten commit
builds.

### PR 8 — Build custom patches directly from a commit's diff view

After it: `space` over a commit's diff (commit-files, commits, sub-commits,
stash, reflog) toggles lines into a custom patch, a checkmark gutter shows
membership, the secondary pane previews the patch through your diff renderer,
`d` removes lines from the commit, and moving/undoing patches keeps your
selection.

Commits:

1. **gocui: the on-demand inclusion gutter** — `SetInclusionGutter(show,
   marks)`: reserved left column, ✓ on every wrapped segment of marked buffer
   lines, content shifted, wrap width narrowed; pure draw-time decoration
   (buffer/metadata/click resolution untouched). Unit tests. Refs: 702c29651
   + every-segment fixup, N§21.20/N§21.22(5).
2. **PatchBuilder: identity-based accessors** — included line identities per
   file; `IncludedChangeLineIndices` (ordinal mapping for the secondary);
   **thread `previousPath` correctly** — the prototype hardcoded `""` at
   three call sites after the rename rebase (N§21.36(1)); production looks
   up the `CommitFile` by path and passes `GetPreviousPath()`, mirroring
   `toggleForPatch`/`RefreshPatchBuildingPanel`. Refs: e57135979, b4270b7d9
   (accessor half), N§21.36(1). Note master's `33b8d497c` added a mutex to
   PatchBuilder (worker `Reset()` vs UI-thread readers) — keep the new
   accessors within that locking discipline.
3. **Toggle from the commit-files main view** — `space` routes to the patch
   toggle (per the panel's `PrimaryAction`); decides add/remove from the
   first selected line; starts the builder if inactive (discard-confirm when
   a patch for another commit is active); refreshes normally (same diff
   command → scroll/selection survive for free, N§21.21); gutter recomputed
   on focus/toggle, shown iff a patch is active AND either pane of the
   focused-main pair is current (`NextInStack(current)`, N§21.35 follow-up);
   auto-advance by the toggled change-line count (`advanceBy`, N§21.35).
   e2e: `build_from_main_view`. Refs: d3a34c203 (+ §21.21/§21.22 fixups),
   6834b39af, 13a64d5ec.
4. **Toggle from the whole-commit main views (commits/sub-commits/stash)** —
   panel-agnostic back end (`patch_building_from_main_view.go`); target
   derived from the panel's selected ref via `FromAndToForDiff` (decoupled
   from `CommitFilesContext`); cheap refresh (`PostRefreshUpdate(panel)`, no
   commit-list reload); sub-commits/stash gain the secondary patch view +
   gutter wiring. **Includes the nil-ref crash guard** in
   `refreshCommitFilesContext` (+ regression test
   `reset_patch_built_from_main_view`). e2e: `build_from_whole_commit…`,
   `build_multi_file_from_whole_commit…`. Refs: 6b3a713b6, fe5c43839 +
   crash-guard fixup, N§21.23.
5. **Reflog patch-building** — wire the reflog panel the same way (it was an
   oversight, not a limitation — N§21.24); needs the same toggle handler +
   `previousPath` care. New e2e.
6. **`d` — discard selected lines from the commit** — reset any active
   patch, build a one-off patch from the selection, `DeletePatchesFromCommit`
   via rebase; disabled (greyed with reason) on non-rebaseable panels
   (stash, other-branch sub-commits, mid-rebase) and in the secondary pane.
   e2e: `discard_lines_from_commit_main_view`. Refs: eaec32b2b (commit half),
   b4270b7d9 (secondary-disable).
7. **The secondary patch pane: preview through the renderer + removal by
   identity** — the preview: render the patch as a real diff by materializing
   `a/`+`b/` temp trees under lazygit's temp dir (from-side blobs; `git
   apply` of the patch; added files: empty `a/<file>`, absent `b/<file>`,
   `PatchToApply(false,false)`), rendered via `git diff --no-index
   --no-prefix a b` through the normal diff-renderer wiring; a generation
   counter drives lazy rebuilds.
   **Removal: build the identity bridge, NOT the prototype's ordinal
   bridge.** Matching by line number against the original diff is out (the
   aggregated patch renumbers included additions, N§21.35(1)), and the
   prototype's ordinal bridge (`ChangeLineOrdinalsInViewRange` →
   `included[ordinal]`) is **broken under difftastic** (diagnosed
   2026-07-18, works under delta — memory [[merge-staging-into-main-view]]):
   it assumes the displayed change lines equal the patch's change lines in
   order and multiplicity, but difftastic's inline mode groups all deletions
   before all additions per hunk, and a collapsed modification row carries
   `d`+`a` while the prototype's `DiffLineContents` kept only the first
   payload per row. Production: resolve **all** payloads per row
   (`DiffLineContent.Metadata`, PR 4 deviation 1), match each
   `(type, new, old)` identity
   against the identities computed from the **raw temp-tree diff** (the same
   patch arithmetic as `parseFileSection`), and map the k-th match to
   `included[k]`. The gutter and the main-pane toggle already match by
   identity — reuse that machinery. This also makes change lines a renderer
   never displays (whitespace-only under difftastic) harmless instead of
   ordinal-shifting.
   **Open sub-item — the a/b path leak over the temp trees:** an external
   diff tool receives the literal `a/…`/`b/…` paths, so difftastic emits
   `file=b/<path>` in its records and renders a "Renamed from a/… to b/…"
   banner; the host's `patchFilename` lookup then finds no patch-builder
   file and the removal silently no-ops. The `--no-prefix` masquerade only
   cleans the *textual* diff (buffer-parse and delta are fine). Production
   must normalize the tree prefix when resolving records emitted over the
   temp trees — and decide whether the rename banner is acceptable
   cosmetically. Diagnosis is assessed, no fix chosen yet; decide with the
   user.
   Prereqs (both landed 2026-07-18): difftastic's rewritten patch-space
   `d`+`a` emitter (its `osc-1717-metadata` branch) and the gocui
   zero-width-record carrier (PR 4 commit 1) — without them the `d` half of
   a collapsed modification row is invisible to any bridge.
   **Open sub-item — renames in the temp-tree rendering** (a renamed file
   materializes at two paths; check what `--no-index` shows and whether
   `--find-renames` is needed) — resolve during implementation, ask the user
   if it's ugly. e2e: `remove_lines_from_main_view_secondary`, plus a
   reordered/multi-record case if expressible headlessly (a fake conforming
   renderer à la PR 7's handshake fake can emit difftastic-shaped records).
   Refs: b4270b7d9 (removal — the ordinal version, superseded), e0cde9b88,
   957952566, N§21.35.
8. **Preserve the selection across commit rewrites** — the command-agnostic
   net: the four commit-diff panels install an ordinal restore before
   `RenderToMainViews` when (main view focused + selection shown + no restore
   pending + **the diff command actually changed**). No bespoke
   commit-discard reveal (the net covers it — build fca748e36's end state).
   e2e: `keep_selection_after_moving_patch_out_main_view`,
   `undo_keeps_focused_main_view_selection`. Refs: 2ea867faa, fca748e36,
   N§21.33.
9. **Allow changing context size during custom patch building** — ref:
   10bb69d80 (read its message for the rationale/constraints).
10. **Recompute the inclusion gutter when a renderer switch re-renders the
    diff** — switching renderers keeps the same git command but yields a
    different buffer-line structure, so marks computed from the pre-switch
    buffer misalign. The prototype deferred this (N§21.22(4)); mandatory for
    production — it looks too broken to ship (decided 2026-07-19). The
    recompute must run against the *new* buffer at render completion;
    candidate mechanisms: ride the renderer-switch restore's `Apply` (PR 6
    installs one on the Normal view on every renderer cycle) or a general
    post-swap hook on the buffer manager. Decide at implementation. Not
    e2e-assertable (draw-time) — interactive sign-off (§6).

#### Deviations from the plan (2026-08-19, as implemented)

Landed as 10 commits on branch `build-custom-patch-from-main-view` (off PR 7),
plus four `fixup!` commits (one on the accessor commit, two on the toggle
commit, one on the patch-removal commit). All checks green, the whole e2e suite passes, §6 sign-off owed.

1. **The plan's commits 3, 4, 5 and 6 are one commit.** `space` and `d` arrive
   together for all five commit-diff panels, because the seam is one interface:
   a panel that joins `FocusedMainViewActions` answers both `PrimaryAction` and
   `DiscardSelection`, and one shared object serves all five panels (deviation
   2), so there is no way to stage them per panel or per key without either a
   stub method or splitting the interface for the sake of the sequencing. The
   alternative — landing discard first — has the same problem the other way
   round.
2. **All five panels share one `CommitDiffActions`** (`commit_diff_actions.go`),
   introduced by a prep commit that has them stop answering for their own diff.
   They differ only in *which* diff it is, which each panel says as a
   `commitDiffTarget` (from, to, whether it may be rewritten) — and that answers
   `PlainDiff` too, every one of them having been `PlainDiffBetweenRefs` over
   its own endpoints. `SwitchToDiffFilesController.canRebase` came out of
   `enter()` as the plan asked.
3. **The panels' refresh is `PostRefreshUpdate` for all of them**, including the
   commit files panel, where the plan (following the prototype) had
   `Refresh({COMMIT_FILES})`: the ◐/● per-file marks are drawn from the patch
   builder rather than from the model, so re-rendering the list is enough and
   nothing has to be re-read from git.
4. **`previousPath` is threaded by asking git which files the diff renames**
   (§8's mandatory rename row): `GetFilesInDiff` for the target, once per
   toggle. The model's `CommitFiles` would only do for the commit files panel —
   the whole-commit panels can have never populated it, or populated it for
   another commit — and getting it wrong is silent, since the patch builder
   caches the first diff it loads for a file.
5. **The selection net (plan commit 8) is one gui-layer rule, not four
   per-panel calls**, and it landed *before* `d` rather than after: it fixes
   what PR 7 already shipped (moving a patch out of a commit, or undoing that,
   left a stale selection painted over the rewritten diff), so it stands on its
   own and `d` inherits it. `Gui.keepDiffSelectionAcrossACommitRewrite` asks it
   of every render, from `RefreshMainOpts` alone, gated on there being a
   selection, no more precise restore pending, and the render being of a
   *different* diff. Two supporting changes: the post-action reveal moved onto
   `DiffLineHelper` (the gui package cannot reach a controller), and the key a
   render is remembered under is now taken before the pty path adds its git
   config, so that it says which diff is being rendered and nothing else.
6. **The gutter is recomputed at the one place a pane's content settles**, which
   is what closes plan commit 10 without a mechanism of its own:
   `updateDiffSelectionVisibility` became `updateDiffPaneDecorations` and does
   both. A renderer switch, a context-size change and walking the commits all
   go through it. The marks are also refreshed on focus and focus-lost (as
   N§21.35's follow-up decided, re-evaluating rather than hiding, so that
   moving to the pane beside the diff keeps them) and right after a toggle.
7. **Which lines are in the patch is asked of the panel**, through a new
   `FocusedMainViewActions.PatchInclusion()` returning a predicate over the
   diff's lines (nil for the working tree's diff, and for a commit's when the
   patch being built is of another one). The helper knows the rendered rows; only
   the panel knows whether the patch is of *this* diff, and a patch of another
   commit would otherwise mark lines that merely share a file and a line number.
8. **The marks are e2e-assertable after all** — `View.MarkedLines` and
   `ViewDriver.MarkedLines`/`NoMarkedLines`, in the spirit of PR 5's
   `SelectionIsShown`: the plan called the gutter draw-time-only, but which
   lines are marked is state, and only how it *looks* needs the interactive
   pass. That is what makes the renderer-switch row testable
   (`patch_marks_follow_a_renderer_switch` fails without the recompute).
9. **The trees a renamed file is materialized into use the path the patch
   expects** (§9.4's rename sub-item, resolved as announced): the name it had
   before, where the patch carries the rename, and the new name for a partial
   selection, whose patch has the rename stripped. `PatchBuilder.FilesInPatch`
   answers that per file. So a whole-file rename previews as a rename, a pure
   rename included. The wart: git writes `rename from a/original` /
   `rename to b/renamed` over the trees, the two paths being all it has to go
   by — `renamed_file_whole` documents it (see §8's new row). *(Landed with
   only the whole-file half working; the partial half was fixed 2026-08-21 —
   see §10's log.)*
10. **The a/b leak is normalized where the pane's identities are handed out**
    (§9.4's other sub-item, decided with the user): `DiffLineHelper.inRepoTerms`
    maps a path stated over the trees back to the repo's file, at the two places
    identities leave the resolver for a consumer (`DiffLinesInViewRange` and
    `diffLineIdentitiesAt`). That is fewer places than the plan's "per
    consumer", and it keeps `e` in the preview pane working as it did when the
    pane rendered the patch as a string. Nothing normalizes for the diff's
    *text*, which the a/b masquerade already makes read like the repo's own.
11. **Copy works over the preview pane**, since the pane's `PlainDiff` is the
    trees' own diff and its identities are in the repo's terms: the lines copied
    are the patch's, not the commit's lines that happen to share their numbers
    (which is what PR 7's copy did there, wrongly). Covered by an extension to
    `copy_selected_diff_lines`.
12. **Removal counts change lines in the trees' diff**, as the plan's identity
    bridge asked, via `DiffLineHelper.ChangeLineOrdinals` and
    `PatchBuilder.IncludedChangeLineIndices`. The e2e test that guards it uses
    additions interleaved with context, which is what makes the count differ
    from the commit's diff — a run of consecutive additions gives the same
    answer either way, so a test over one proves nothing (learned by neutering
    it).
13. **The nil-ref crash (plan commit 4's guard) reproduces exactly as N§21.23
    described** and is fixed in the toggle commit, by a `fixup!`:
    `captureCommitFilesState` says whether there is a commit to load the files
    of at all, the commit files panel never having been pointed at one. Test:
    `reset_a_patch_built_from_a_commits_diff`.
14. **Resetting the patch no longer leaves the diff it was built from.**
    `PatchBuildingHelper.Reset` escaped from any non-side context, which was the
    patch-building view; now it escapes from that view alone, since it is the
    only one with nothing left to show. Same test.
15. **The `space`/`d` descriptions are per diff type** (`diffActionDescription`):
    "Stage"/"Discard" over the working tree's diff, "Toggle lines in patch"/
    "Remove lines from commit" over a commit's. The static `Description` the
    cheatsheets are generated from still says only the working-tree half — see
    §8's new row.
16. **Plan commit 9 (context size) is deferred to PR 9** (decided with the
    user): the refusal exists for the explorer, whose patch-building view hands
    the patch builder line indices from a freshly loaded diff while the builder
    holds one cached at the context size in force when the file entered the
    patch. The main view is immune, mapping by identity, but the explorer lives
    until PR 9 — where the refusal goes away in one line. See §9.11.
17. **e2e**: 9 new tests in `pkg/integration/tests/main_view/` —
    `build_patch_from_a_commits_diff`, `build_patch_from_a_whole_commits_diff`,
    `build_patch_from_a_reflog_entry`, `discard_lines_from_a_commit`,
    `discard_from_a_commit_only_where_it_can_be_rewritten`,
    `patch_marks_show_while_the_diff_is_focused`,
    `patch_marks_follow_a_renderer_switch`,
    `custom_patch_goes_through_the_diff_renderer`,
    `remove_lines_from_the_custom_patch`,
    `reset_a_patch_built_from_a_commits_diff` — plus
    `patch_building/keep_selection_after_moving_patch_out_main_view` for the net,
    the copy extension, and three gocui unit tests for the gutter. Two existing
    tests changed with the preview's rendering: `specific_selection` (git's own
    hunk context) and `renamed_file_whole` (the rename lines' tree names).

#### Review round 1 (2026-09-05) — resetting the patch from the pane showing it

Reported: with a custom patch being built and the focus in the main view, either
pane, "Reset patch" from the custom patch menu leaves the pane previewing the
patch on screen.

`PatchBuildingHelper.Reset` renders again through `Context().Current()`. On
master that is always a side panel, because the same function pops out of any
context that isn't one first; the commit "Build a custom patch from a commit's
diff, and discard lines from it" narrowed that pop to the patch-building view,
so that giving up the patch leaves you in the diff you were building it from.
From the focused main view `Current()` is then a main context, and a main
context renders nothing to main, so the pane that was previewing the patch keeps
its content.

One `fixup!` on that commit: render through `Context().CurrentSide()` instead.
Both main panes are rendered by the panel beneath them, so a reset from within
either of them has to go through that panel. Where the current context is
already a side panel the two are the same, so nothing else changes. The focus
follows out of the pane that goes, `followFocusIntoWorkablePane` doing its
usual work. e2e: `reset_the_patch_from_the_pane_showing_it`, plus two
assertions on the pane's visibility in the existing
`reset_a_patch_built_from_a_commits_diff`.

The commit that distils the helper into `CustomPatchHelper` (PR 9's "Remove the
explorer behavior behind the retired panels") carries the change forward.

### PR 9 — Replace the staging and patch-building panels with the focused main view

The removal PR. Also the PR whose title tells users the big story — consider
making *this* the umbrella release-notes headline ("staging now happens
directly in the diff view") since PRs 7/8 titles already describe the
mechanics.

Sequencing inside the PR (every commit green):

1. **Migrate explorer e2e tests to main-view flows first** — while both UIs
   still exist. Triage each test under `pkg/integration/tests/staging/` and
   `…/patch_building/` (~54 pre-prototype tests): (a) behavior also covered
   by an existing main-view test → delete; (b) behavior worth keeping →
   rewrite to drive the focused main view; (c) explorer-specific rendering/
   plumbing tests → delete with the panels. Several commits, grouped
   sensibly. Also sweep other suites that `enter` into staging incidentally
   (grep for `Views().Staging`/`.PatchBuilding` and `PressEnter` on files).
2. **`enter` on a file focuses the main view** — files panel and commit-files
   panel: `enter` (and double-click on the file row) pushes the focused main
   view anchored at that file's diff (multi-file/directory diff → anchor at
   the file's first row via the jump-to-file landing logic). Selection
   anchors per PR 5 rules.
3. **Remove the explorer machinery** — contexts (`Staging`,
   `StagingSecondary`, `CustomPatchBuilder`), their views/windows in
   `context/setup.go` and layout, `StagingController`,
   `PatchBuildingController` (explorer half), `patch_exploring` package,
   `RefreshStagingPanel`/`RefreshPatchBuildingPanel` (keep/rewire the
   *secondary patch panel* update path — PR 8's renderer-based preview stays,
   fed by `secondaryPatchPanelUpdateOpts`), escape/`EscapeFromPatchExplorer`
   remnants, `IPatchExplorerContext`. Multiple commits: this is the risky
   demolition — go subsystem by subsystem.
4. **Config + keybinding + i18n cleanup** — remove explorer-only keybindings
   from cheatsheets (`just generate`); rename `useHunkModeInStagingView` and
   `wrapLinesInStagingView` (they now govern the main view) using the config
   migration mechanism — **agree the new names with the user first**
   (candidates: `useHunkModeInDiffView`, `wrapLinesInDiffView`); remove
   orphaned english.go strings (only english.go — Crowdin cleans the rest).
5. **Docs** — `docs-master/` staging/custom-patch docs rewritten for the new
   model; Config.md/schema via `just generate`.
6. **Allow changing the context size while a patch is being built** — moved
   here from PR 8 (its commit 9; decided with the user 2026-08-19). Deleting
   the refusal is a one-liner (`ContextLinesController.checkCanChangeContext`
   and `Tr.CantChangeContextSizeError` both go), and this is where it becomes
   safe: the refusal protects the *explorer*, which rebuilds its state from a
   freshly loaded diff and hands the patch builder line indices into it, while
   the builder holds the diff it cached at the context size in force when the
   file entered the patch — a mismatch for WHOLE-mode files as much as for
   PART. The main view has never had that problem, mapping by line identity,
   which no context size renumbers. e2e: change the context size mid-build
   from the main view, then toggle another line and check the patch.

Risk note: this PR is where hidden couplings surface (things that push
`Staging` contexts from unexpected places — merge-conflict flows, custom
commands, `git bisect` edge flows). Grep for every reference to the removed
contexts/views before starting; expect a long tail of small fixes.

**Status: DONE 2026-08-20** on branch `replace-staging-panels-with-main-view`
(25 commits), reviewed 2026-09-01 with 9 `fixup!`s, one commit dropped, one
split, and two commits added; §6 sign-off owed. The risk note landed: the
coupling that surfaced was `customCommands[].context`, which is matched
against the context keys by *name* (see deviation 4).

Deviations from the plan above:

1. **The four "file-tree workflow" custom-patch tests were not the
   explorer's** (found in review). `select_all_files`, `toggle_directory`,
   `select_direcories_sharing_prefix` and `toggle_range` drive the **commit
   files panel**, whose `toggleForPatch`/`toggleAllForPatch` survive this PR
   untouched; the commit that deleted them ("Drop custom-patch tests for the
   file-tree workflow") was dropped, and all four are back unchanged. Without
   them, `a` in the commit files panel had no e2e coverage at all, and neither
   did the prefix-sharing directory case (`foo` vs `foobar`).
2. **Three of the eleven tests dropped in the first commit weren't covered
   elsewhere** (found in review): `stage_partial_block_of_changes_{first,last,
   middle}_lines` assert `patch.Transform`'s output when only part of a run of
   deletions+additions is staged, and the middle one deliberately documents a
   known imperfection. Restored under `main_view/`, needing only the view and
   config-key renames. Likewise the staging-panel screen-mode test, restored as
   `main_view/change_screen_mode_in_focused_diff` (it reaches the diff with
   `FocusMainView` rather than `enter`, so it is green from its own commit on).
3. **"Edit hunk" (`E`, `keybinding.main.editSelectHunk`) was dropped by
   accident** (found in review): `StagingController.EditHunkAndRefresh` went
   with the panel, nothing replaced it, and the config key stayed in
   `user_config.go`, the schema and Config.md, binding nothing. **The port went
   into PR 7, not here** — see the note at the end of PR 7's section: a command
   the main view has to gain before the panel can go belongs with the rest of
   the build-up, not inside the removal.
4. **A custom command naming a removed context made lazygit exit** (found in
   review, decided with the user): `context: staging` was valid, and an
   unknown name reaches `log.Fatal` in `keybindings.go` — so the config of
   anyone who used one would take the program down on startup. Rather than
   migrating the name, the names are now validated as the config is read. That
   is a fix for a master-level bug (any typo did the same), so it went in its
   own branch, `validate-custom-command-contexts`, inserted at the foot of the
   stack; this PR's shell-removal commit drops the four names from the list it
   added. See §10's log.
5. **The whole-file decision asks the file, not the diff's shape.** The
   as-landed `SelectionRepresentsWholeFile` used `Patch.IsSingleHunkFor
   WholeFile`, whose own comment records that it is wrong at context size 0 —
   a wart that only "doesn't matter" because of a guard in another file, and
   this is the PR that removes the *other* context-size guard. The question is
   now put to `models.CommitFile.Added()`/`Deleted()`, which `filesInDiff()`
   already loads for the rename paths, and the two answers ("which indices"
   and "is that every change") come out of one parse. That left
   `Patch.IsSingleHunkForWholeFile` with no callers, so it went too.
6. **The tooltips of the selection commands follow their descriptions**
   (found in review). `Tr.RemoveSelectionFromPatchTooltip` became orphaned
   here, because `MainViewController` pairs a dynamic `DescriptionFunc` with a
   static `Tooltip` — so over a commit's diff the `Remove` key said "Remove
   lines from commit" while its tooltip described `git reset`, and the
   rebase-conflict warning went unseen. `types.Binding` gained a `TooltipFunc`
   (static `Tooltip` still being what the cheatsheet prints); the fix is a
   fixup in **PR 8**, whose commit introduced the mismatch.
7. **The hunk-staging hint's removal is its own commit** (asked for by the
   user): it was folded into the config rename, where it was invisible.
8. **`dropDiffPrefix`'s move out of the explorer controller is its own
   commit**, ahead of the deletion, per AGENTS.md's prep-refactor rule.
9. **`FilesController.EnterFile` takes a clicked line, not `OnFocusOpts`**: the
   `ClickedWindowName` it used to pick the staging-secondary view is dead now
   that the pane comes from `GetMainPanes()`. Renamed to `enterFile` while
   there.
10. **The demos and `discard_old_file_changes` press `enter`**, not
    `FocusMainView` — the gesture the README now documents, and the only e2e
    exercise of the new `enter` besides its own test. Folded into "Open file
    diffs in the focused main view", the commit that gives `enter` that
    meaning.
11. `excludedViews` in `cheatsheet/generate.go` was left as an empty slice
    behind a `lo.Contains` that can never fire; removed.

#### Review round 1 (2026-09-10) — the wrap option, and a hint that is no diff

Both of the round's problems turn on one question: what counts as a main pane
holding *the panel's diff*, as opposed to a message, a log, or a diff shown as
part of an explanation. The user chose to have the render answer it —
`NewMainViewDiffTask`, which every diff render already goes through, marks its
task, and `types.ContentIsDiff` reads the mark back — over having the panel
answer per selected file. Nothing at a call site carries a flag of its own.

1. **`wrapLinesInDiffView` applies to diffs alone.** It was governing both main
   panes whatever they held, so with wrapping off a branch's commit log, the
   status, and the message explaining a merge conflict were all cut off at the
   edge of the pane. On master they always wrapped, the option reaching only
   the staging view. Each render is now wrapped as its own content asks. New
   e2e test `wrap_only_the_diff`, with `ViewDriver.ContainsViewLines` to assert
   on the lines as the view lays them out rather than as the content has them.
2. **The rename and the behaviour are two commits** (asked for by the user).
   "Name diff options after the view that now uses them" was doing both. The
   rename comes first, leaving the option governing nothing for one commit, and
   "Let wrapLinesInDiffView govern the two main panes" follows with (1) as its
   `fixup!`.
3. **No selection over the conflict hint** — the question round 3 of PR 7 left
   open. A conflict that can only be resolved by picking a side is explained
   rather than diffed, and for a file deleted on one side and modified on the
   other git's diff of that modification is shown below the explanation. The
   pane took those change lines for a diff of its own and offered to stage
   hunks of them. The same mark settles it: that render doesn't carry it, so
   the pane shows no selection and every command that acts on one is disabled.
   The diff mode stays hard-coded to `DiffModeRendered`, which is what the
   user asked for. It lands in **PR 7** as a commit of its own, "Show a
   selection only over the diff the panel offers", right after "Show git's own
   diff when the renderer's can't be acted on", with
   `no_selection_over_a_conflict_hint`.

   **A `fixup!` was wrong here** (the user, correcting the first attempt): the
   commit it would have folded into decides *which diff to render*, and this
   decides *whether the content can be pointed at*. Two decisions, two commits,
   even though the second only becomes possible with the first. The fix belongs
   in PR 5, where the selection came from; that being impractical, its own
   commit here says in its message why it arrives this late.

Diffing mode (`W`) keeps the selection it has today: its render is a diff and
says so, so §8's deferred row is untouched. PR 7 renders the custom patch
preview as text it assembles itself, so the mark reaches string renders too
there; PR 8's commit that materializes the patch takes that half away again in
a `fixup!` of its own, the preview being a diff of two trees from there on.

PRs 8 and 9 were replayed over the PR 7 fixup, and `diff-file-menu` and
`edit-diff-line-with-modified-click` over the new PR 9 tip. Every commit from
the split onwards builds and unit-tests on its own; the whole e2e suite passes
at the PR 7 fixup, at the PR 9 tip, and at the top of the stack.

### PR 10 — Alt- or shift-click a diff line to open it in your editor

**Status: DONE 2026-08-20** on branch
`edit-diff-line-with-modified-click`, replayed onto PR 9's tip since. Nine
commits plus round 1's eight `fixup!`/`amend!` commits, each building,
unit-testing and linting clean on its own; whole e2e suite passing. §6
interactive sign-off owed.

Self-contained after PR 5 (uses PR 4's `GetDiffLineInfo` and PR 5's focused
main-view edit path). Commits (N§19):

1. **gocui: let a mouse binding opt into firing while a popup is focused**
   (`HandleWhenPopupPanelFocused`). Ref: ac85a90ed.
2. **Extract `editDiffLine` from `editLine`** — prep. Ref: d761f07d1.
3. **gocui: carry the press-time keyboard modifiers through the whole mouse
   gesture** — snapshot the modifiers at button press and stamp them on the
   press, every drag event (ORed with `ModMotion`), and the release;
   modifier changes while the button is held are ignored. Master delivers no keyboard
   modifiers on mouse events at all (`mouseMod` is `ModNone` for every one of
   them, and the `mouseMod = ModNone` in its release branch is dead code), so
   this branch is where they start arriving — **say that in the PR
   description**, since
   the commit message reads as though only the carrying were new. Master's
   gesture model (#5854) made the pre-rebase press-only fix (da4201aa2, on the
   `-plan` copy) insufficient: bindings match modifiers exactly, so a
   modified press that nothing consumed would otherwise start matching
   unmodified drag bindings mid-gesture (drag-select), and a modified
   gesture's release would look like a plain one. **Inverts master's
   `TestMouseReleaseDoesNotKeepPressModifiers`** (#5854 deliberately
   dropped press modifiers on the release; this design reverses that call —
   say so in the PR). **Global behavior change** (unbound modified clicks
   become no-ops instead of acting as plain clicks) — flag in the PR
   description. Ref: "Carry the press-time keyboard modifiers through the
   whole mouse gesture" (post-rebase, 2026-08-04 — transcribe this shape,
   not da4201aa2's).
4. **The feature** — alt-left *and* shift-left both bound (no single chord
   survives Ghostty+iTerm2+VS Code — N§19.1); no focus change, no selection;
   works behind popups. Ref: a86da2e97.
5. **gocui: give views a transient line flash** — reverse the configured
   selection-width bar without changing the selection or replacing the diff
   renderer's colors; clear all transient flashes when the UI suspends. Tests
   cover drawing over an already-selected row and suspension cleanup.
6. **Acknowledge modified-click edits with a 200 ms flash** — arm the flash
   only after the clicked row resolves to an edit, and use a generation token
   so an older click's timer cannot clear a newer flash. Force a content-only
   draw before invoking the editor, so the acknowledgement precedes even a slow
   editor CLI; the timeout begins after the invocation returns. A suspending
   editor clears the flash before disengaging the terminal, so nothing appears
   on resume.

Three gocui preparations sit below the list, in this order: "Cleanup: remove
error return value from `Gui.SetRune`, `Gui.draw()` et al", "Remove unused
per-line highlighting" (`View.SetHighlight` had no production callers, and the
flash is the per-line marking gocui does need), and "Separate redraws from
view-line invalidation", so the flash can ask for a repaint without
claiming the cached wrapping is stale.

A refinement to the click handling in the main view was implemented here too: a
plain click inside the currently selected hunk collapses the selection to the
clicked line, matching how a click inside a range selection already behaves.
**It has since moved to PR 5**, where the click behaviour it refines was
decided; see round 1 below.

Interactive sign-off: Ghostty, iTerm2, VS Code (already done once for the
prototype; re-confirm the transcription).

#### Review round 1 (2026-09-10) — the branch reviewed for the first time

Nine commits by then, three more than the list above. All checks were green,
and the round found two defects in the mouse handling, no test for the command
at all, no documentation for the gesture, and two message slips. A third
suspected defect turned out to be the behaviour we want (item 4).

1. **The hunk-collapse commit belongs in PR 5**, and moved there as an `amend!`
   for "Select a whole change block when focusing the main view in hunk mode".
   That commit decides what a click does in hunk mode and spends a sentence of
   its message on the rule; this one carves a case out of that rule, rewrites
   the same comment, and edits assertions into the middle of that commit's own
   test. Left in PR 5 as a commit of its own it would have had PR 5 change its
   own click rule five commits after stating it.

   **The reason the paragraph above gave for keeping it here does not hold.**
   Replaying the 76 commits over the target took one conflict, in "Select a
   range of diff lines by dragging", where the two `SetDragAnchorViewLine`
   lines go back around the new branch. It also took one adjustment that git
   makes no conflict of: "Give the focused main view's selection a home outside
   the controllers" moves the diff-line helpers into `diff_line_selection.go`,
   and it has to take `SelectedHunkBounds` and the new call site along, so it
   is marked `edit` in the todo rather than picked. Without that, the tip stops
   compiling on a `showSelectionAtLine` that commit deleted. The tree at the
   tip came out identical to the tip before the move.

2. **A modified click that no binding wants moved a list panel's selection
   highlight.** Mouse events start carrying keyboard modifiers in "Keep mouse
   gesture modifiers stable" (see the item above), and bindings match modifiers
   exactly. `onKey` moves the view cursor before it consults the bindings,
   though, and a list view draws its selection at that cursor. So alt- or
   shift-clicking the files, branches, commits or stash panel moved the bar
   while the panel's selected item stayed put, and the next action worked on
   the item the bar had left. The cursor move and the mouse capture now happen
   only for a gesture carrying no keyboard modifier. Motion is not one of
   those, so a drag still moves the cursor as it did.
   `fixup!` on that commit, with
   `TestAModifiedClickNoBindingWantsLeavesTheViewAlone`.

3. **A click a popup swallowed armed a double click.** "Let mouse bindings work
   behind focused popups" has to know about a double click before the
   `ShouldHandleMouseEvent` gate rejects anything, and it moved the recording
   of the click up there too. So a click behind a popup, the popup dismissed,
   and the same cell clicked again inside the threshold arrived as a double
   click. `isDoubleClick` now only asks, and `recordClickInfo` is called where
   a binding is about to see the click — in the early pass once one matches,
   and after the gate otherwise. `fixup!` with
   `TestASwallowedClickIsNoHalfOfADoubleClick`.

4. **The modified click opens a line of a conflict hint in the editor, and
   should** (raised as a defect, resolved the other way by the user
   2026-09-12). It looked like one because the edit keybinding refuses there,
   by its disabled reason and by a guard of its own, while `editClickedLine`
   has neither: the hint for a file deleted on one side and modified on the
   other embeds git's own diff of the modification, and those rows parse as
   diff lines like any other.

   **The two commands need different gates.** `e` needs `view.Highlight`
   because without a selection it has no line to act on. A click names its own
   line, so it needs nothing of the sort — and for a `DU` or `UD` conflict the
   file is in the working tree either way, holding the modified side, so
   opening it at the clicked line is worth having. You may want to copy a piece
   of that hunk elsewhere before resolving the conflict by deleting the file.
   The other content the pane can hold needs no gate either: a message has no
   `diff --git` line above it, so `fileSectionBounds` finds no file section and
   `editDiffLine` returns having done nothing. A `HasSelectableContent` guard
   was written and then dropped; what stayed is the opposite assertion, in
   `no_selection_over_a_conflict_hint`, so that a pane with nothing to select
   is on the record as still having lines to point at.

5. **The command had no test**, nor did the gocui flag it rests on. New e2e
   test `edit_clicked_diff_line` covers both modifiers, the file and line the
   editor is pointed at, that the focus and the selection stay where they are,
   and the click landing while a menu holds the focus. It needed
   `GuiDriver.ClickWithModifier` and the `AltClick`/`ShiftClick` pair on
   `ViewDriver`; the harness delivered modifier `0` on every mouse event
   before. New gocui test `TestOnlyBindingsThatOptedInFireBehindAFocusedPopup`
   covers `HandleWhenPopupPanelFocused` in both directions. Two `fixup!`s, one
   per target.

6. **The gesture was undocumented.** Mouse bindings reach no cheatsheet, so
   nothing said it exists. `docs-master/Config.md`'s "Configuring File Editing"
   section now names it beside `e`, and `Custom_DiffRenderers.md` says it next
   to delta's `--hyperlinks`, the affordance the feature commit's message says
   it generalizes. `fixup!` on the feature commit.

7. **Two slips.** The cleanup commit's subject named `View.draw()`, which
   already returned nothing; it changes `Gui.draw` (`amend!`). And removing
   `View.SetHighlight` left two comments naming a view's own highlighting as a
   reason for `firstDirtyLine` to move. Only `write` moves it now (`fixup!`).

Every commit from the amend! in PR 5 to the tip builds and unit-tests on its
own; whole e2e suite green at the tip. Backups: the pre-round tip is
`edit-diff-line-with-modified-click-2026-09-10-1900-backup`, and the tip before
item 4's guard was dropped is
`edit-diff-line-with-modified-click-2026-09-12-1000-backup`.

### PR 11 — Open the selected diff line in the branch's GitHub PR

**Status: DONE 2026-09-13** on branch `open-pull-request-at-diff-line`, off
PR 10's tip. Four commits plus round 1's three `fixup!`s and an `amend!`, each
building, unit-testing and linting clean on its own; whole e2e suite green.
§6 interactive sign-off owed.

Self-contained; after PR 5. Planned as one or two commits (N§5):

- `openPullRequestForSelectedLine` on `Commits.OpenPullRequestInBrowser` in
  the focused main view: URL `<pr.Url>/changes/<commitSha>#diff-<sha256(relPath)>R<line>`;
  commit sha from the side panel's `RefForAdjustingLineNumberInDiff`; path
  relative to **`WorktreePath()`** (never `RepoPath()` —
  [[worktree-path-vs-repo-path]]), forward slashes, exact bytes into sha256;
  branch resolution per panel (commits → checked-out; subCommits → its ref;
  commitFiles → parent). GitHub-only via `PullRequestsMap`. Ref: 912703d20.
- Unit-test the URL builder. PR description should note the anchor format is
  empirically derived (undocumented by GitHub).

#### Deviations from the plan (2026-09-13, as implemented)

The prototype commit was cherry-picked onto the stack to see it work and then
written again from scratch, at the user's word. Four commits: three
preparations and the command.

1. **Which branch's pull request is a question the panel answers**, rather than
   a `switch` over context keys in the controller (the user's call). The new
   `types.PullRequestDiffContext` sits beside `DiffMainViewContext`: the
   commits panel answers with the checked-out branch, the sub-commits panel
   with the ref it was entered from where that is a local branch (a tag or a
   remote branch is none), and the commit files panel by asking the panel it
   was entered from. A panel that doesn't implement it — files, stash,
   reflog — has no pull request, and the command is not described there, so it
   stays out of the keybindings menu.
2. **Three defects of the prototype, all in what the URL points at.** A
   deleted line was pointed at as `R<NewLine>`, where a deletion's `NewLine` is
   only where it sits in the new version of the file; it is now `L<OldLine>`,
   the side GitHub shows it on. A file-header row asked for a line of the file
   it names; it now points at the file alone, as `e` opens the file there. And
   a row that is no line of the file at all (`\ No newline at end of file`)
   asked for line 0.
3. **The command is disabled where the pull request has no view of the diff**
   (both guards the user's call): in diffing mode, where the main view shows a
   diff against another ref, and over the custom patch's preview pane, whose
   lines sit at the numbers the patch gives them. Two new strings,
   `NotAvailableInDiffingMode` and `NotAvailableForCustomPatch`. Without a
   pull request for the branch it says so in a panel, as the commits and
   branches panels do.
4. **Three preparations**, each for something the command would otherwise have
   copied: `HostHelper.PullRequestForBranch` /
   `NoPullRequestDisabledReason` (the branches and commits panels each had
   their own lookup and disabled reason), `repoRelativePath` in
   `diff_paths.go` (the two diff-action objects each had their own), and
   `MainViewController.sidePanelBeneath` (two questions reached for the panel
   beneath, guard included; the command asks twice more).
5. **The whole path can't be exercised headlessly.** A pull request reaches
   the model either from the GitHub API or from the on-disk cache, and the
   PULL_REQUESTS refresh clears the cache's contents on startup whenever no
   auth token is available, which is the case in the harness and on CI. So
   `open_pull_request_only_over_a_commits_diff` covers where the command is
   offered and the three reasons it refuses, and `TestGithubPullRequestLineURL`
   covers the URL itself, including the `L`/`R` sides and the SHA-256 of the
   path.

#### Review round 1 (2026-09-17) — the commits a pull request can be asked for

The user read the branch and raised two things the URL gets wrong. A range of
commits selected in the commits panel opened the pull request at the newest of
them, so a line of the selected hunk may be nowhere in the commit the page
shows. And a branch that has diverged from its remote opened a page saying "We
went looking everywhere, but couldn't find those commits"; a commit from before
the branch gets the same page.

GitHub's range form was read off a real pull request (lazygit#5870, 8 commits),
by fetching candidate URLs and counting the file anchors each page carries:

- `/pull/N/changes/<a>..<b>` is **exclusive on the left**. `c7..c8` shows only
  c8's two files, and `c1..c3` shows the nine of c2 and c3 but not c1's. So the
  left end is the parent of the range's oldest commit, and lazygit's own range
  diff (`git diff <oldest>^ <newest>`) is the page's.
- Both ends have to be commits **of the pull request**. The commit it was
  opened against 404s there, as does `<a>..<a>`.
- **`BASE`** stands for the commit the pull request was opened against (the
  user found it in GitHub's own UI): `BASE..c2` shows commits 1 and 2, and
  `BASE..<tip>` the whole pull request. Case doesn't matter.
- A commit the pull request doesn't hold 404s on its own too. This is the page
  the user saw.

Two `fixup!` commits and an `amend!`, inserted mid-branch, and a third `fixup!`
correcting a comment in the second:

1. **`CommitsForPullRequest` on `types.PullRequestDiffContext`** answers with
   the commits whose diff is on screen and the hash of the commit that diff
   starts after, in place of the single hash the command took from
   `RefForAdjustingLineNumberInDiff`. `githubCommitRange` names one commit by
   its hash and a range as `<base>..<newest>`, with `BASE` for the base where
   the pull request holds no commit before the range. Two shared functions in
   the context package, both unit-tested: `commitsShownInDiff` mirrors what the
   panel hands `GetUpdateTaskForRenderingCommitsDiff`, so a selection the panel
   has no range to diff for (a rebase todo entry at either end, or a range
   spanning the divergence boundary in the sub-commits panel) reports the one
   commit the pane shows; `pullRequestBaseForCommits` looks the oldest commit's
   first parent up among the panel's commits and offers it only while it is
   pushed.
2. **The command refuses where the pull request doesn't hold the commits.**
   `Status == StatusPushed` is the test (the user's call): a pull request holds
   the commits of its branch that are on the remote, so an unpushed commit is
   none of its own, and neither is one that is in a main branch already. Two
   strings, `CommitNotInPullRequest` and `CommitsNotInPullRequest`. Amending a
   commit in the middle of a branch leaves the command working below it, as the
   user asked.
3. The `amend!` rewrites the message around the range form and the refusal.

Nothing headless reaches a pull request (deviation 5), so the refusal sits
below `NoPullRequestDisabledReason` where no e2e test can see it. The range
naming and the base lookup are unit-tested instead.

### PR 12 — Jump to a file of the diff from a menu

**Status: DONE 2026-09-13** on branch `diff-file-menu`, off PR 11's tip. Three
commits, each building, unit-testing and linting clean on its own; whole e2e
suite green. §6 interactive sign-off owed.

This is PR 5's commit 7, skipped there because the UX wasn't decided (PR 5
deviation 1) and revived as a PR of its own. `f` opens a menu of the diff's
files, in the order the diff shows them and by the paths the repo knows them
by; picking one goes where next-file navigation would have landed.

Commits:

1. **Name the file a diff row belongs to in the repo's terms** — `filePaths`
   was the one query that skipped the mapping `inRepoTerms` applies, so over
   the custom patch's preview it answered with the path of the tree the patch
   was materialized into. The mapping comes out of `inRepoTerms` as
   `repoTermsMapper` and `filePaths` goes through it. Where a renderer states
   the path of each side of a change, this also has both halves belong to one
   file rather than to the two trees.
2. **Ask a diff where each of its files begins** — `fileStarts`, the whole
   answer at once, with `fileStart` picking its neighbour out of it (the walk
   backwards over a file goes away). The menu and `n`/`N` then agree on where
   a file begins by construction. Unit-tested by `TestFileStarts`, with
   `TestFileStart` unchanged as the proof that navigation still behaves.
3. **The menu** — `DiffLineHelper.FilesInDiff` and `StartOfFileInDiff`, the
   `keybinding.main.jumpToFile` config entry (default `f`), the
   `JumpToFileInDiff` string for both the binding and the menu title, and
   `FilterAsYouType` so that a file is a few characters away however many the
   commit touches. The menu reads the diff to the end before it is built, as
   the search does: a file below the part that has been read is in neither the
   list nor the view. e2e: `jump_to_a_file_of_the_diff`, whose first file's
   diff is longer than the viewport, so the read matters.

Deviations from the prototype commit: proper i18n and a config entry for the
key (the prototype hard-coded both, and the plan's open question 2 asked for
the entry), the repo-terms paths above, and **each menu item carries its file
rather than the view line that file begins at**, so that a diff re-rendered
while the menu is up is jumped into at the row the file begins at now.

#### Round 1 (2026-09-13) — two refinements, to be tried before they are kept

Both asked for by the user after reading the branch, and **both await an
interactive pass** before anything is decided about where they belong.

1. **A single-file diff gets a message rather than a menu** — a menu with one
   file in it has nothing to choose. How many files there are is only known
   once the diff has been read to the end, which is too much work for a
   keypress that only asks whether a key applies, so it can't be the key's
   disabled reason; the handler says it instead, as an error toast carrying
   the same `Disabled: ` prefix a disabled key's reason would have, so that it
   reads as the same thing. New string `OnlyOneFileInDiff`. `fixup!` on the
   menu commit.
2. **A file you go to is scrolled to the top of the view** rather than to the
   middle, for the menu *and* for `n`/`N` — a file put in the middle wastes
   half the screen on the diff just left. Only where the view has to scroll at
   all: a file already on screen leaves the view alone. In hunk mode the
   selection is the file's first change rather than the row the file begins
   at, and a large context size can put that change more than a screenful
   below it, so the alignment is applied first and the selection scrolled into
   view afterwards; where the two can't both hold, the selection wins and is
   centred as before. **Approved the same day and moved down to PR 5** as an
   `amend!` for "Jump by hunk and by file in the focused main view", that
   commit being where file navigation and its landing rows are decided; its
   message gained a paragraph about the scrolling. What stayed here is the one
   line of the menu asking for the same alignment. New
   `ViewDriver.TopVisibleLine`, which asserts on where a view is scrolled to
   without naming a line number, and which reads a **view** line: the
   viewport's origin counts wrapped segments, and a commit's diffstat wraps.

New test `file_navigation_scrolls_to_the_top` (in PR 5) covers both hunk-mode
cases with `n`/`N`; `jump_to_a_file_of_the_diff` asserts the same for a jump
from the menu, and gained the single-file case. Two things that cost time:

- Pressing the key before the diff has rendered reads an empty view, so a test
  has to wait for the selection first (this one flaked once in a full run
  before that wait was there).
- The alignment stops at the last screenful, so a file at the **end** of the
  diff can't reach the top of the view — a test asserting that it does needs
  a file after it to scroll past.

The move down took three passes of the stack: the `amend!` where it belongs,
then the config rename in PR 9 ("Name diff options after the view that now
uses them") to take the new test's `UseHunkModeInStagingView` with it, then
the menu commit for its fixture. Every commit from the `amend!` to the tip
builds; the tip passes the whole suite.

---

## 6. Interactive sign-off matrix

The headless harness cannot run real diff renderers, `LAZYGIT_SLOW_RENDER`,
or the pty path (N§13.1), and the gutter is draw-time-only. Each PR needs a
user pass before merge:

| PR | What to verify interactively |
|---|---|
| 1 | ✅ **APPROVED 2026-08-09.** Slow-render matrix (N§11/§13): flick commits/files scrolled down; 10 s auto-refresh (`refreshInterval: 3`) — no content/scrollbar flicker; **also re-test at normal speed** (N§20.5). Found PR 1 deviations 8 and 9, both fixed; a repo with dirty submodules is the case that exposes a slow same-content re-render |
| 4 | ✅ **APPROVED 2026-08-09.** Patched delta/difftastic/diff-so-fancy emit + render cleanly; handshake swallowed (no phantom line) |
| 5 | ✅ **APPROVED 2026-08-15.** Selection feel under delta; hunk-on-click; drag incl. autoscroll; nav under metadata delta incl. repeated `n` across files. Some special cases are candidates for a later refinement; deliberately not pursued now. **`n`/`N` scrolling a file to the top of the view approved 2026-09-13**, tried in PR 12 before it moved down here |
| 6 | ✅ **APPROVED 2026-08-15.** `{`/`}`, `ctrl+w` and renderer-cycle scrolled down: no top-jump, offset preserved, both anchor cases; ignoring whitespace where it removes the anchor's hunk, and where it empties the diff. Nothing found; the whitespace consumer called out as a welcome addition |
| 6b | Copying under delta (unified + SxS) and difftastic: a hunk with its header, a file header, a selection that is all additions, a commit's message alone and a selection spanning it into the diff; and both toasts. The headless tests use fake renderers throughout |
| 7 | ✅ **APPROVED 2026-08-16**, except for `E` ("Edit hunk"), ported here on 2026-09-01 and still owing a pass with a real editor, including a patch edited to something neither side of the diff says. Full staging matrix under no-renderer / patched delta (unified + SxS) / difftastic; cross-pane focus-follow; raw fallback feel under stock delta / diff-so-fancy-without-metadata; binary-file focus stability (N§21.30 repro). Four review comments about the stack as a whole, all fixed the same day — see PR 7's sign-off section |
| 8 | Gutter under delta/no-renderer/difftastic; whole-commit path on LocalCommits (canRebase menu); secondary pane preview per renderer; **secondary-pane removal under difftastic specifically** (the prototype's known-broken case: reordered `d`/`a` records, collapsed modification rows, a/b record-path leak) and under delta |
| 9 | `enter` and double-click on a file (working tree and commit) under each renderer; `{`/`}` down to 0 and back while a patch is being built; the keybindings menu's tooltips over both kinds of diff; screen modes with a diff focused; `wrapLinesInDiffView: false` with a long line in a diff, a branch log, the status and a conflict hint on screen in turn (round 1) |
| 10 | Ghostty, iTerm2, VS Code |
| 11 | The URL the browser lands on, in a repo whose branch has a pull request: a line of a commit's diff, a deleted line (`L`), a file-header row, a range of commits (`<base>..<newest>`), a range reaching down to the pull request's first commit (`BASE..<newest>`), and the same from the commit files panel and the sub-commits panel; plus the refusal over an unpushed commit and over one from before the branch (round 1). Nothing headless reaches a pull request (PR 11 deviation 5), so every one of these is untested |
| 12 | The menu over a many-file commit under each renderer (a file of a difftastic diff begins at its first content row), and the landing row for each; filtering as you type. Round 1's scrolling is **approved 2026-09-13** and has moved to PR 5; the message a single-file diff gets is still for the user to try |

Patched renderer builds: `cargo build` in delta/difftastic worktrees
(`osc-1717-metadata` branches); diff-so-fancy is a script.

## 7. Compatibility with the parked separate-lists design

`separate-lists-design.md` (worktree `separate-lists-for-staged-and-unstaged`,
doc-only, parked until this lands) will put staged/unstaged files in two
sections of one files panel. Keep these seams clean so it stays cheap:

- **Side-of-action stays derivable and localized**: the "which side does this
  pane show" logic (`diffSplitState`, `mainShowsStaged`-style decisions) and
  the focus-follow rule live in *one* place each (the files handler); don't
  let call sites re-derive them. Separate-lists will want side to come from
  list-section membership instead.
- **Focus-follow may need to become configurable/section-aware**: that design
  wants "stay on the acted-on side's *section*" after emptying a side, which
  is the opposite of the merged view's "follow the content to the other
  pane". Don't hard-code the rule into more than one function.
- **`<tab>` semantics**: keep pane-toggling expressed as one operation so it
  can later also move a list cursor.
- The split-main-view rendering itself is load-bearing for the merged staging
  UX and stays.

## 8. Known gaps and their dispositions

Shortcuts the prototype deliberately took. Dispositions **reviewed with the
user (2026-07-19)**: rows marked **Fix** are mandatory scope — "the prototype
deferred it" never meant "optional", only "not addressed while prototyping".
The remaining rows are agreed as keep/defer:

| Gap | Disposition |
|---|---|
| Rename support in the from-main-view patch paths (N§21.36(1)) | **Done in PR 8**: the previous path comes from asking git which files the target's diff renames (`GetFilesInDiff`), once per toggle — the model's file list only ever describes the commit files panel's own commit (PR 8 deviation 4) |
| patch pkg rename-aware Parse/Transform/FormatView (N§21.36(2)) | **Closed 2026-08-09 — nothing to fix**: doesn't reproduce off master; `patch_building` e2e green with PR 2's patch changes, rename unit tests added (PR 2 deviation 6) |
| Reflog patch-building (N§21.24) | **Done in PR 8**: the reflog panel shares the one `CommitDiffActions` with the other four, and gained the patch preview pane — as did the stash and sub-commits panels, which never had one either |
| Renames in the custom-patch temp trees (new, this plan) | **Done in PR 8**: a file is materialized under the path the patch expects it at — the name it had before, where the patch carries the rename — so a whole-file rename previews as a rename, and a partial one as the content change it is, under the new name (PR 8 deviation 9, completed 2026-08-21 by the two fixups in §10's log: the trees' path and the content's are two different questions) |
| Secondary-pane removal broken under difftastic — ordinal bridge + a/b record-path leak (diagnosed 2026-07-18, memory) | **Done in PR 8**: a line of the patch is found by counting the change lines of the *trees' own diff*, which no rendering can reorder or hide; and a path stated over the trees is brought back to the repo's own where the pane's identities are handed out (PR 8 deviations 10 and 12). The interactive pass under real difftastic is still owed (§6) |
| Diffing mode (`W`) not wired to the raw fallback → not stageable (N§21.29) | Defer; note in PR 7 description ("diffing-mode staging is its own question") |
| `type: extDiff` with empty `command` (git's `diff.external`; formerly `useExternalDiffGitConfig`) always-raw when focused (N§21.30) | Keep; document |
| Per-pane selection memory on `<tab>` (re-anchors each switch, N§21.9) | Defer; follow-up candidate |
| `IsSingleHunkForWholeFile` hunk-default refinement (N§21.11) | **Done in PR 5** (derived from the rendered diff, no git call — PR 5 deviation 3) |
| `a` on a context line below the last hunk doesn't snap back like staging did (N§21.11) | **Done in PR 5** (review round 1, fix 1): `ChangeBlockBounds` falls back to the block above |
| Deleted-file `MD`-vs-`D` staging special case (N§21.13) | **Done in PR 7**: a selection covering every change of a file stages the file |
| `NormalSecondary` not preserved on `-U`/renderer change (N§16.1) | Keep as documented limitation |
| The focus stays in a main pane the **merge-conflicts view** takes the window from (new, round 4) | Defer; `followFocusIntoWorkablePane` returns early for any pair but the Normal one, so this is unchanged from before the round. Same class as the pane-goes-away bug: reachable when a focused file becomes conflicted underneath you |
| A **side-by-side** renderer re-laid-out at a new width loses the position (new, round 4) | Defer. Only renderers whose line count depends on width (`delta --side-by-side`, difftastic side-by-side), and only when the diff is actually re-run at a new width — a refresh-driven render or a screen-mode change, never a bare resize. Needs `PreserveDiffPositionOnRerender` on a plain refresh, gated on the render being of the *same* diff, which isn't knowable until the render starts |
| Gutter marks for not-yet-loaded lines of huge diffs (N§21.20) | Keep (marks appear on next recompute); note |
| Renderer switch mid-patch-build shifts checkmarks (N§21.22(4)) | **Done in PR 8**: the marks are worked out again wherever a pane's content settles, so no mechanism of its own was needed; `patch_marks_follow_a_renderer_switch` guards it (PR 8 deviations 6 and 8) |
| Copy copies the renderer's output verbatim under a renderer (N§21.28) | **Done in PR 6b**: copy takes the corresponding lines of the plain diff, through the new diff-source seam. The rows above the diff, which belong to no file, are the one thing taken from the screen |
| Nav only sees loaded content (deep targets in huge diffs, N§16.4) | **Done in PR 5**: ReadToEnd-then-retry in the shared `navigate` helper |
| Toggle auto-advance: no "skip already-included" smarts (N§21.35) | Keep plain next-hunk |
| difftastic token-vs-line `c`-at-new-line mismatch (M§10.2) | Protocol v2 candidate; nothing to do host-side |
| `scrollUpMain`/`scrollDownMain` scroll the hidden upper pane when the lower one has the section to itself (new, round 5) | **Done in round 5**: both ask `Gui.mainSectionView`, which knows about `SecondaryPaneOnly`; fixup for "Always show a file's staged changes in the lower pane" |
| `clearMainView` leaves a pending `RenderRestore` stranded, and with it the `BeginBlockingEvents` its `Done` balances (new, round 5) | **Done in round 5**: `clearMainView` drops it, as the string renders do; fixup for "Hold input back until the selection has moved on". Untested — it needs a mispredicted target pane, which takes an external mutation racing the action |
| git names the trees in the custom patch preview's `rename from a/…` / `rename to b/…` lines (new, PR 8) | Keep; documented by `renamed_file_whole`. The two paths are all git has to go by over `--no-index` trees, and any naming of them leaks there; the `---`/`+++` lines and every line's identity are the repo's own. Only renames are affected, and difftastic — which reports a rename whenever the two paths it is handed differ — says it for every file |
| A drag re-anchors the search only when the mouse is released (new, 2026-09-05) | Keep. gocui moves the cursor for a drag itself, so the only lazygit-side moments are the drag handler, which fires per pointer move, and the release. Re-anchoring per move would re-render the "x of y" for every frame of a drag for no gain |
| The cheatsheets describe `space`/`d` in the focused main view by their working-tree meaning only (new, PR 8) | Defer. `pkg/cheatsheet/generate.go` reads the static `Description`, which is one string per binding, while the key now means two things depending on the diff; the options bar and the keybindings menu say the right one (PR 8 deviation 15) |
| A diff still being read shows no selection until a change line has arrived (new, PR 5 round 5) | **Closed in round 7**: the pane reads on until it can tell, so the wait is however long it takes to reach the first change line, and no user action is needed to end it |
| Focusing at the top of a commit whose diffstat fills the screen lands the selection on a stat row (new, PR 5 round 6) | **Keep** — the user's call, 2026-09-06. Focusing never moves the view, and with no change line on screen the selection goes to the middle visible line. Reaching the first hunk from there is one press of `a` or `right`, which is preferable to the view scrolling on its own |
| A modified click on a renderer's hyperlink opens the hyperlink, not the clicked line (new, PR 10) | Keep. The hyperlink is handled before any mouse binding and ignores modifiers, as on master. Both paths open the same file at the same line for `lazygit-edit://` links, so only a renderer pointing its links elsewhere would tell the difference |
| `e` over the custom patch's preview opens the file at a line the patch numbers, not the commit's (new, PR 11) | Raised, not acted on. The preview is a diff of the two trees the patch was materialized into, so a line below an omitted change sits at a number the file doesn't have it at, and `AdjustLineNumber` carries that number forward as if it were the commit's. PR 11 refuses there for the same reason; `e` has behaved this way since PR 8 and is left as it is for the user to decide on |
| A range of commits selected in the commits panel opens the pull request at the newest of them (new, PR 11) | **Done in round 1**: the URL names the range as the pull request's own pages do, `<base>..<newest>`, with the keyword `BASE` where the range starts where the pull request itself does. The form was read off a real pull request, one candidate URL at a time |
| The command refuses over the commits of a branch whose pull request is merged (new, PR 11 round 1) | Keep. Merging a pull request puts its commits in a main branch, so they come out `StatusMerged` rather than `StatusPushed`, and the test for what a pull request holds turns them down although its pages still show them. Telling a commit of the remote branch apart from one that only reached a main branch takes a rev-list against the upstream that nothing else needs, and the local branch is usually gone by the time its pull request is merged |
| A renderer that keeps the diff and hunk headers but drops body lines could be mis-parsed where it ends the buffer (new, PR 2 round 1) | Keep. The leniency applies to one section, the one the buffer breaks off in, and every renderer that restructures a body lengthens hunks rather than shortening them. A mis-parse would act on the wrong line only in the focused main view, and there the diff is either git's own or one whose lines state their own identity (`MainViewDiffMode`) |
| A submodule's log lines go unresolved wherever diff-so-fancy has taken a column off them (new, PR 2 round 2) | Keep. Nothing states a record for those lines, so they are placed by parsing them, and diff-so-fancy strips the leading indicator column from every line it reads while it is inside a hunk — which a submodule's section below one still counts as. The cost is that `n` pressed on one of them steps to the file after the next. The submodule itself is listed and navigated to from the line naming it, which every renderer passes through as git wrote it |

## 9. Open questions (resolve before/during the marked PR)

1. ~~**PR 3:** does the per-entry `pager:` config field keep its name?~~
   Resolved by #5870: `pager`/`externalDiffCommand` were unified into a
   single `command` field interpreted per the new `type` field.
2. ~~**PR 5:** proper keybinding config entries for `n`/`N`/`f`?~~ Resolved
   2026-08-10: yes — `keybinding.main.prevFile`/`nextFile` (`N`/`n`). `f`
   followed on 2026-09-13 with the jump-to-file menu: `keybinding.main.jumpToFile`
   (PR 12).
3. **PR 9:** new names for `useHunkModeInStagingView` / `wrapLinesInStagingView`
   + config migration.
4. ~~**PR 8:** the two temp-tree sub-items of commit 7 — renames in the
   temp-tree rendering, and how to normalize the `a/`/`b/` tree prefix for
   records an external diff tool emits over the temp trees (+ whether its
   "Renamed from a/… to b/…" banner is acceptable).~~ Resolved 2026-08-19:
   materialize a renamed file under the path the patch expects it at, so the
   rename previews as a rename (PR 8 deviation 9); normalize where the pane's
   identities are handed out (deviation 10). The banner turned out not to be
   about renames at all — difftastic is handed `a/foo` and `b/foo` and reports
   a rename because they differ, for every file — and is left as it is (§8).
5. **PR titles**: drafts in §4 — the user finalizes wording at PR-open time
   (they're the release-notes lines).
6. **Cross-repo timing** (outside this plan): circulating the OSC 1717 spec,
   upstreaming the three renderer patches and the git one. lazygit ships
   fully functional without any of them; revisit pitching once PRs 1–8 exist
   as evidence. git's own patch may never be accepted — a maintained fork is
   the accepted fallback (PR 4 cross-repo note), and no PR here waits on the
   outcome.
8. **PRs 5/7/8: e2e tests against *real* patched renderers** (raised by the
   user 2026-08-09; deferred again at the start of PR 7 — its renderer tests
   are about renderers that say nothing, which a fake models exactly, so the
   helper waits for PR 8's difftastic work to need it). Today's renderer
   tests use fake shell commands, which can imitate a record stream but not
   difftastic's actual reordering and collapsing — the shapes that broke the
   secondary-pane removal (§8). Requiring delta/difftastic on `PATH` needs no
   harness change to *reach* them: `PATH`/`TERM` are already inherited from
   the host environment. Decisions taken with the user:
   - **A missing renderer is a hard failure, locally and on CI** — not a
     skip. So no `Skip:` gating; the failure message should name the binary
     and how to install it, since other contributors run `just e2e` too.
   - **Asserting on the renderer's rendered text is fine.** Updating tests
     when a renderer changes its output is ordinary maintenance, no different
     from lazygit's own panel-content tests. The real obstacle is **color**:
     delta's default and diff-so-fancy convey the +/- side by color alone, so
     color-stripped text can't distinguish an addition from a deletion, and
     the assertion reads as nonsense. The harness's `ContainsColoredText`
     matches **foreground only**, while delta marks the side by
     **background** — so either configure the renderer for legibility in the
     test (delta's `--keep-plus-minus-markers`; difftastic has no equivalent),
     or extend the harness with a background-color matcher, or assert on what
     got staged instead. Decide per renderer.
   - **CI install is unblocked**: all three emitters are open draft PRs
     (delta 2181, difftastic 1014, diff-so-fancy 538 — see the spec repo's
     §10). Pin to a commit so "the emitter changed" can't look like "lazygit
     broke"; cache the cargo builds on that pin.

   Keep the fake-renderer tests either way: they cover the protocol shapes no
   real renderer emits on demand (the handshake, header records, records
   covering no cell). Cheapest sequencing is a small standalone harness helper
   when PR 7 starts, not folded into a feature PR.
9. ~~**PR 7:** how should `rawGit` entries with restructuring args decide
   the raw fallback?~~ Resolved 2026-08-07: probe them like any other
   renderer, since git announces itself for exactly the formats it
   describes. See PR 7 commit 10.
11. ~~**PR 8: can the context size be changed while a patch is being built?**~~
    Resolved 2026-08-19: not in PR 8 — the refusal is there for the explorer,
    which would hand the patch builder line indices into a diff it no longer
    holds. It is deleted in **PR 9** instead, along with the explorer (PR 9
    item 6, PR 8 deviation 16).
10. ~~**PRs 7/9: what does `ctrl+w` do in a main view you can stage from?**~~
    Resolved 2026-08-16 (PR 7 deviation 14): nothing to do — the patch is
    built from a freshly fetched plain diff, which never ignores whitespace,
    and line numbers are the same either way. PR 9 still deletes the dead
    context-key list. Original question:
    (Raised 2026-08-15 while planning PR 6's whitespace consumer.) Master
    refuses ignoring whitespace in the staging and patch-building contexts —
    `ToggleWhitespaceAction` matches on those three context keys and answers
    `IgnoreWhitespaceNotSupportedHere` — because a patch built from a
    whitespace-ignoring diff doesn't apply. In the focused main view the
    current context is `Normal`, so the refusal doesn't catch it, and once
    PR 7 binds `space` there the toggle silently becomes a way to build a
    broken patch. Options: refuse it when the panel beneath is
    `Staging`/`PatchBuilding` (master's rule, expressed in the new
    classifier); or allow it and disable the acting keys while it's on; or
    treat it like a non-conforming renderer and stage from the raw diff.
    Decide in PR 7; PR 9 must in any case delete the now-dead context-key
    list.

## 10. Progress

- [x] PR 1 — async render fixes — **DONE 2026-08-09** on branch
      `fix-async-diff-rendering` (16 commits, all checks green, §6 sign-off
      approved), stacked on `fix-task-key-race`
- [x] PR 2 — diff-line identity primitive — **DONE 2026-08-09** on branch
      `resolve-diff-lines-to-identities` (5 commits plus round 1's 2 `fixup!`s,
      all checks green, every commit builds and tests clean), stacked on
      `fix-async-diff-rendering`. No interactive sign-off needed: no
      user-visible change
- [x] PR 3 — rename pagers → diff renderers — **landed on master as #5870**
      (with a bigger config rework than planned; see the PR 3 section)
- [x] PR 4 — OSC 1717 support — **DONE 2026-08-09** on branch
      `support-osc-1717-diff-metadata` (7 commits, all checks green, every
      commit builds and tests clean), stacked on
      `resolve-diff-lines-to-identities`. §6 sign-off **approved**
- [x] PR 5 — selection & navigation — **DONE 2026-08-10** on branch
      `select-diff-lines-in-main-view` (9 commits, plus round 4's 2 commits,
      round 7's 2 preparations, and the 4 `fixup!`s of rounds 5 to 7, all checks
      green, every commit builds/tests/lints clean on its own), stacked on
      `support-osc-1717-diff-metadata`.
      Jump-to-file menu skipped and copy moved to PR 7 (see its deviations).
      §6 sign-off **approved 2026-08-15**. One `amend!` added 2026-09-10 by
      PR 10's round 1, for the click that collapses hunk mode
- [x] PR 6 — position preserve — **DONE 2026-08-15** on branch
      `keep-diff-position-on-rerender` (7 commits, fixups folded, all checks
      green, every commit builds and unit-tests clean on its own), stacked on
      `select-diff-lines-in-main-view`. §6 sign-off **approved**
- [x] PR 6b — copy the selected diff lines — **DONE 2026-09-13** on branch
      `copy-diff-lines-from-main-view` (PR 7's first two commits, split out at
      the user's suggestion, plus round 1's `amend!` and two `fixup!` commits,
      all checks green), stacked on `show-staged-changes-in-lower-pane`. §6
      sign-off owed
- [x] PR 7 — staging from the main view — **DONE 2026-08-16** on branch
      `stage-changes-in-main-view`, now off `copy-diff-lines-from-main-view`,
      plus two commits added 2026-09-01 porting
      "Edit hunk" (see PR 7's addendum; `E` still owes its interactive pass)
      (17 commits with round 4's folded in, 55
      across the whole stack, all checks green, every commit building and
      unit-testing clean on its own), stacked on
      `show-staged-changes-in-lower-pane`, which is itself stacked on
      `keep-diff-position-on-rerender`. §6 sign-off **approved 2026-08-16**,
      with four cross-cutting review comments fixed as mid-branch fixups in
      PRs 5, 6 and 7 (see PR 7's sign-off section)
- [x] The staged side always in the lower pane — **DONE 2026-08-16** on branch
      `show-staged-changes-in-lower-pane` (5 commits after round 5, green),
      inserted below PR 7 at the user's suggestion; see the section at the end of
      PR 7
- [x] PR 8 — custom patches from the main view — **DONE 2026-08-19** on branch
      `build-custom-patch-from-main-view` (10 commits plus 4 `fixup!`s, and
      round 1's one more, every commit unit-testing clean on its own, whole e2e
      suite passing), stacked on `stage-changes-in-main-view`. Plan commit 9
      moved to PR 9; §6 sign-off owed
- [x] PR 9 — panel removal — **DONE 2026-08-20** on branch
      `replace-staging-panels-with-main-view` (25 commits, reviewed 2026-09-01:
      9 `fixup!`s, one commit dropped, one split, two added — 34 commits now),
      stacked on `build-custom-patch-from-main-view`, which is itself now
      stacked on `validate-custom-command-contexts`. Every commit builds,
      unit-tests and passes the whole e2e suite on its own; §6 sign-off owed
- [x] Validating a custom command's context — **DONE 2026-09-01** on branch
      `validate-custom-command-contexts`, **landed on master as #5989**
- [x] Rendering a searched main view — **DONE 2026-09-05** on branch
      `rerender-main-view-while-searching` (9 commits off master once its
      `fixup!`/`amend!` pairs are folded in, all checks green), inserted at the
      foot of the stack; master-level bugs that PR 7 turned into a hang. Its
      own PR, mergeable ahead of the rest
- [x] PR 10 — alt/shift-click edit — **DONE 2026-08-20** on branch
   `edit-diff-line-with-modified-click` (9 commits plus round 1's eight
   `fixup!`/`amend!`s, every commit green, whole e2e suite passing), stacked on
   PR 9; §6 sign-off owed
- [x] PR 11 — open PR at line — **DONE 2026-09-13** on branch
   `open-pull-request-at-diff-line` (4 commits: three preparations and the
   command, every one green on its own), stacked on PR 10; **round 1 on
   2026-09-17** added three `fixup!`s and an `amend!` for the commits a pull
   request can be asked for; §6 sign-off owed, and nothing headless can reach a
   pull request (PR 11 deviation 5)
- [x] PR 12 — jump-to-file menu — **DONE 2026-09-13** on branch
   `diff-file-menu` (3 commits, every one green on its own), stacked on PR 11.
   PR 5's skipped commit 7, revived; §6 sign-off owed

(Add per-commit checkboxes inside each PR section as work starts; record
deviations from this plan inline, dated.)

Log:

- **2026-09-18 (later):** **The foot branch that renders without a pty on
  Windows, reviewed for its design.** `render-diffs-without-a-pty-on-windows`
  sits below the whole stack: in a ConPTY a diff renderer's OSC 1717 records
  reach lazygit detached from the rows they describe, so on Windows the command
  runs through a pipe, with a stdin filter as a command of lazygit's own. The
  review found the structure sound and one duplication real. COLUMNS was set in
  two places because `newCmdTask` and the render path each built their own
  task; the plain way of running a command is now a third sibling of the pty
  and pipe ways, and the three share one tail, `newTaskForRender`, where
  COLUMNS is set once. GIT_PAGER moved into the pty way, the only one git reads
  it in. `RunPtyTask` became `RunDiffRendererTask`, a pty being one way of two.
  LAZYGIT_COLUMNS is gone, the width reaching every renderer as COLUMNS or
  `{{width}}`. Both pipeline entry points gate their log entry the same way,
  and a handful of comments no longer compare the code with what it did
  before. Left alone at the user's word: the width parameter on the
  renderer-command getters, and the piped filter running under cmd.exe where
  git would have used its sh. Nine `fixup!`s, one `amend!` and two commits
  mid-branch, the stack replayed in one pass with a `break` after each target;
  every upstack line naming the task type was renamed on the way. Unit, lint
  and the whole e2e suite green at the branch tip and at the stack tip; backup
  tag `jump-to-file-from-diffstat-2026-09-18-2055-backup`.

- **2026-09-18:** **PR 2 round 2**, on a submodule of the repo. A commit that
  moves a submodule shows it in a section of its own, with no `diff --git`
  header, and nothing placed those rows in a file: navigation stepped over the
  submodule, the menu of the diff's files left it out, and clicking its name in
  the diffstat found no such file. One `fixup!` mid-branch for the buffer
  parser, and one apiece for delta and diff-so-fancy, which state a record for
  every row and had none to state for this one. A second `fixup!` followed the
  same day, on what the user found next: a submodule's section has to end where
  git's lines for it end, or under a renderer — which prints no `diff --git`
  line for it to stop at — it takes every untagged row below it, and the menu
  lists the submodule between every pair of files. Whole e2e suite green, all
  129 commits replayed; backup tags
  `jump-to-file-from-diffstat-2026-09-18-1845-backup` and `…-2015-backup`.

- **2026-09-17:** **PR 11 round 1**, on the two things the user found in the
  URL. A range of commits opened the pull request at the newest of them, and a
  branch that had diverged from its remote opened a page that couldn't find the
  commits at all. GitHub's range form was read off a real pull request by
  fetching candidate URLs: the left end is exclusive and has to be a commit of
  the pull request, and the keyword `BASE` (the user found it) stands for the
  commit the pull request was opened against. So the panel beneath now answers
  with the commits whose diff is on screen plus the commit that diff starts
  after (`CommitsForPullRequest`), the URL names a range as `<base>..<newest>`,
  and the command refuses wherever a commit of the diff isn't pushed. Three
  `fixup!`s and an `amend!`, inserted mid-branch; PR 12 replayed on top; whole
  e2e suite green.

- **2026-09-13:** **PRs 11 and 12 written**, both from prototype commits the
  user had cherry-picked onto the stack to see them work and then asked for
  again from scratch. PR 11 opens the selected line in the branch's pull
  request: which branch that is comes from a new context interface rather than
  a switch over context keys (the user's call), a deleted line points at the
  left side of the diff where the prototype pointed at the right, and the
  command refuses in diffing mode and over the custom patch's preview, neither
  of which the pull request has a view of. Its three preparations each take
  something the command would have copied. **Nothing headless reaches a pull
  request** — the PULL_REQUESTS refresh clears the cache when there is no auth
  token — so the test covers where the command is offered and why it refuses,
  and the URL is unit-tested. PR 12 is PR 5's skipped jump-to-file menu,
  revived with a config entry for `f`, proper strings, and menu items carrying
  their file rather than a view line. Two §8 rows: `e` has the same
  line-number drift over the patch preview that PR 11 now refuses over, and a
  range of commits opens the pull request at the newest of them.

- **2026-09-10 (later):** **PR 10 reviewed for the first time.** Its opening
  commit moved down to PR 5 as an `amend!`, the fixup its own note had called
  too conflict-prone to write: the replay took one conflict and one commit
  marked `edit`, and the tree came out unchanged. Four more findings needed
  code. Two are in gocui, and both come of this branch being where mouse events
  start carrying keyboard modifiers — an unbound modified click moved a
  list panel's highlight bar away from its selected item, and a click a popup
  swallowed armed a double click. The last two are a command with no test and a
  gesture with no documentation. A third finding, the modified click opening a
  line of a merge-conflict hint, **turned out to be the behaviour we want**
  (2026-09-12): the click names its own line, and the file is in the working
  tree for both `DU` and `UD`, and the test that had guarded the gate now
  guards the click instead. Eight `fixup!`/`amend!` commits, two of them
  inserted mid-branch, and one `amend!` in PR 5; three new tests plus
  `GuiDriver.ClickWithModifier` for the harness. Written up as PR 10's round 1,
  with a §8 row for the renderer-hyperlink overlap.

- **2026-09-10:** **Two problems from testing PR 9, both about what a main pane
  is holding.** `wrapLinesInDiffView` was governing every render in the two
  panes, and a selection was drawn over the hint for a conflict that has to be
  resolved by picking a side. A render now says whether it holds the panel's
  diff, and both questions read that answer. Written up as PR 9's round 1, with
  the config rename split off from the behaviour into its own commit as the
  user asked. Two `fixup!`s in PR 7 and PR 8 and one in PR 9; PRs 8 and 9 and
  the two branches above them replayed. Two new e2e tests, one new harness
  assertion, whole suite green at three points in the stack. The PR 7 fixup
  was then made a commit of its own: a fixup would have put two decisions in
  one commit (see PR 9's round 1, item 3).

- **2026-09-06 (later):** **Round 6 tested in turn, and the same fix found short
  again**, this time for a commit with a 10000-line message. The user's reading:
  a pane must never take the answer over from the commit before, and while it
  can't tell it should read on until it can. Round 7 does both, with two
  preparations ahead of it (a mutex on `View.LinesHeight`, and `ReadToEnd`'s
  task-holding pulled out for a bounded read). The answer now belongs to a named
  render, so a re-render of the same content keeps it and nothing else inherits
  it. Every commit from the first fixup to the tip builds, unit-tests and lints
  on its own; whole e2e suite green at the tip.

- **2026-09-06:** **The previous day's round tested, and one of its three fixes
  found short.** The two reported bugs behave; the third fix answered the
  question one screenful too early, and a commit's diff begins below its
  diffstat. Written up as PR 5's round 6, with the measurements that say what
  each of the three fixes is worth on the reported commit. One more `fixup!` on
  "Show a selection in the focused main view", one more e2e test, and one open
  question in §8 about where the selection lands over a diffstat. Every commit
  from the first fixup to the tip builds, unit-tests and lints on its own; whole
  e2e suite green at the tip.

- **2026-09-05 (later):** **Three defects, all folded into the commits that
  caused them.** Two were reported from using the stack; the third came out of
  reproducing the first.

  A selection stayed drawn over a branch's commit log, and resetting a custom
  patch from the focused main view left the pane previewing it on screen. Both
  are written up as review rounds: PR 5's round 5 and PR 8's round 1. The third
  is PR 2's round 1 — with no diff renderer configured, a diff longer than the
  first read of it could not be parsed at all, so focusing it gave no selection
  and none of the commands that act on one. The user's call was to fix it in the
  same round rather than note it.

  All three come of a pane that holds only part of what it is being given.
  `linesToReadFromCmdTask` caps the first read at `height*(height-1)` lines and
  the rest arrives as the user scrolls, so "when the content is final" is a
  moment that may never come, and the code written against it — the decorations,
  and the patch parser's well-formedness gate — was answering about content that
  isn't all there.

  Five `fixup!` commits inserted mid-branch, plus one occurrence added to
  "Rename parsedDiffLine.RelPath to Path" so the branch below it still builds.
  Every commit from the first fixup to the tip builds and unit-tests on its own;
  the whole e2e suite is green at the tip. Backup of the pre-round tip:
  `replace-staging-panels-with-main-view-2026-09-05-1930-backup`.

- **2026-09-05:** **A branch at the foot of the stack, for searching a view
  that is being rendered again.** Reported as a hang: search the focused main
  view, stage a hunk, and lazygit stops taking keys while the mouse wheel still
  scrolls. `postRefreshUpdate` had refused to render the main view while a
  search was on since master's `4e21a096b9` ("Searching can't cope well with
  the view being updated while it is being searched"), and PR 7's staging holds
  input back until that render lands — `revealSelectionInPaneItLandsIn` begins
  a block that the `RenderRestore`'s `Done` ends — so with no render the block
  never ended. The stale diff underneath it is the master-level half, so the
  fix went to its own branch, `rerender-main-view-while-searching`, and the
  stack was replayed onto it.

  Two things the old guard had been hiding, both fixed there. The search
  positions were worked out again from **every write**, and each of those walks
  the whole view. Streaming 2000 lines into a searched view took 565ms against
  10ms unsearched; they are now worked out where they are read, and take 8ms.
  And `searchPositions[currentSearchIndex]` was indexed unguarded in four
  places, so content that lost matches could take the index out of range.

  Two more came out of writing the tests. Rendering a searched view again read
  only as much as the scrollbar needs, dropping the matches below that point,
  where opening the prompt reads to the end — so a render of a searched view
  reads to the end too. And pressing `/` in the focused main view held no task
  while `ReadToEnd` ran, so lazygit counted as idle between the keypress and
  the prompt opening (`docs/dev/Busy.md`), which an integration test takes as
  its cue to press the next key.

  Reviewing that last one, the user moved it from the call site into
  `ReadToEnd`, where it covers every caller. Holding a task there turned a
  quiet bug loud: a request whose task is stopped before it is served was
  dropped, so its `Then` never ran, and with a task attached that task was
  never done either, leaving lazygit permanently busy. Dropping requests is a
  master bug of its own — press `/` as a re-render replaces the task and the
  search prompt never opens — so it goes in below, demonstrated by a unit test
  and fixed by handing requests over through a queue whose reader can go away.
  Asking whether a task is there and giving it the request are one step, as
  are taking the task away and handing back what it never answered, so no
  request can be lost between the two. The queue is unbounded rather than a
  fixed channel for the reasons gocui's `userEventQueue` is.

  Nine commits once the two `fixup!`/`amend!` pairs are folded in, every one
  green. The user tested the shape of it on both the branch and master before
  it was written up. Testing also turned up PR 5's missing
  search-follows-selection behaviour, added there as two commits (PR 5's round
  4).

- **2026-09-01:** **PR 9 reviewed.** The 25 commits were green throughout —
  every one builds, unit-tests and passes the whole e2e suite on its own,
  checked commit by commit — so everything found was something the suite
  cannot see. Eleven findings, all recorded as PR 9 deviations above; the four
  that mattered were a keybinding silently unbound (`E` / "Edit hunk"), a
  `log.Fatal` on startup for a config naming a removed context, four deleted
  tests that had never belonged to the explorer, and three more whose
  behaviour nothing else covered. Two placement calls were taken with the
  user: the tooltip fix belongs to PR 8's commit that introduced the mismatch,
  so it went in as a fixup there (one `rebase -i --update-refs` from below PR
  8, replaying PR 9 on top); and the custom-command validation is a
  master-level fix, so it became its own branch at the foot of the stack, with
  the whole stack re-parented onto it. Two commits were restructured rather
  than fixed up — the explorer-behavior removal gave up its `dropDiffPrefix`
  move to a prep commit ahead of it, and the config rename gave up the
  hunk-staging-hint removal to a commit of its own — and one, "Drop
  custom-patch tests for the file-tree workflow", was dropped outright, its
  whole premise being wrong. Two commits needed *content* amended rather than
  a fixup, because a fixup would have left them red: the shell-removal commit
  has to drop the four context names from the list the new base adds (a
  synchronization test enforces it), and the rename commit has to rename the
  config key in the four test files restored below it. Afterwards, at the user's
  direction, the "Edit hunk" port moved out of PR 9 and into **PR 7**, beside
  the other ported staging commands — the stack's shape is build-the-main-view-
  up, then remove, and a command the main view has to gain belongs in the
  build-up. That needed a prep refactor, PR 7 predating PR 8's diff-line
  identity API, and a rewrite of PR 8's "Say which change line is meant in one
  way" to convert the extracted helper rather than the inline walk. The
  resulting tree is byte-identical to the version before the move, which is the
  check that it was a history change and nothing else.

  The per-commit sweep then caught the new Edit-hunk e2e test failing at one
  commit — flaky, not commit-specific: 1 failure in about 40 suite runs, and
  none in 36 clean ones afterwards. Worth writing down how it was pinned down,
  since random sampling was never going to do it. The truncated failure text
  ruled out the focus and line-count assertions (their messages say "got N"),
  which pointed at the files panel's content. Reading the code gave a candidate
  mechanism — `RunSubprocessAndRefresh` refreshing before the caller applies
  the patch — and a first experiment appeared to disprove it, but had been
  built wrong: delaying the background refresh made it read git *late*, i.e.
  after the apply, which is the safe order. Delaying it between its read and
  its publish instead — `refreshStateFiles` reads at `GetStatusFiles`, publishes
  in the `onUIThreadUnlessRepoChanged` bounce that sets `Model().Files` — showed
  it at once: `Expected 'MM' to be found in ' M file1'`. Under that
  instrumentation the unfixed code failed 8 of 15 runs and the fixed code 0 of
  15. Two lessons for next time: for a publish-order race, delay the *publish*,
  not the work; and a `--update-refs` rebase moves any backup branch left
  sitting on the branch tip, which silently turned two before/after comparisons
  into empty diffs.
- **2026-08-21:** **PR 8 deviation 9 was only half implemented; fixed.** A
  partial selection of a renamed file previewed as a deleted file with no diff.
  `PatchBuilder.FilesInPatch` had one field, `SourcePath`, doing two jobs —
  where the file's content before the patch is read from in `From`, and where it
  is materialized in the trees — which for a *whole*-file rename are the same
  path, so one field looked enough. For a partial one they differ: the content
  is under the old name, while the trees must hold the file under the name the
  patch's own (rename-stripped) header states, the new one. So the content
  lookup missed, the file was materialized empty, and the `git apply` over it
  failed (logged away by `secondaryPatchPanelUpdateOpts`). Split into `Path` and
  `ContentPath` — `PatchFile.Path` had been set and never read. Decided with the
  user: the preview shows a partial rename patch as a plain modification of the
  new path, matching what applying it does; showing the rename would claim the
  patch carries one. Two `fixup!`s inserted mid-stack: the fix plus
  `TestFilesInPatch`/`TestFilesInPatchOfARenamedFile` on **"Show the custom patch
  as the diff it is"** (the commit that introduced the mapping), and the preview
  assertion on **"Keep partial rename patches on the focused diff"** — the e2e
  guard cannot sit with the fix, because at that point the only partial-rename
  test drives the old explorer, whose preview renders the patch as a *string*
  (`PatchBuildingHelper.RefreshPatchBuildingPanel`); only the commit panels' pane
  goes through the trees. Whole suite green, and each fixup green on its own.
- **2026-08-20:** **PR 10 implemented** (6 commits, green; §6 sign-off owed),
  stacked directly on PR 8 so PR 9 remains independent. Added the agreed
  non-suspending-editor feedback as a two-column reverse-bar flash; suspension
  clears its state before terminal editors take over. Recorded the future
  plain-click-inside-selected-hunk behavior separately from this PR.
- **2026-08-19:** **PR 8 implemented** (10 commits + 4 fixups, green; §6 sign-off
  owed). Three decisions taken with the user up front: PR 8 stays one PR; the
  `a`/`b` tree paths a renderer states over the custom patch's trees are
  normalized where the pane's identities are handed out; and — after a
  correction from the user — difftastic's "renamed" banner is not about renames
  at all but about being handed two paths that differ, so it is left alone, and
  rename detection in the preview stays on, a delete-and-add of similar content
  being what git calls a rename anywhere else. One decision surfaced mid-way and
  deferred by the user: **changing the context size mid-build waits for PR 9**,
  the refusal being there for the explorer's sake (PR 8 deviation 16, §9.11).
  The plan's commits 3–6 became one commit — a panel joining the actions
  interface answers for both keys at once — and the selection net (commit 8)
  moved ahead of `d`, since it fixes something PR 7 already shipped. Two things
  worth carrying: the marks turned out to be **state, not just drawing**
  (`View.MarkedLines`), which is what makes the renderer-switch row testable;
  and a test over a run of consecutive additions cannot tell the patch's own
  numbering from the commit's — the additions have to be interleaved with
  context for the count to differ. The nil-ref crash N§21.23 warned about
  reproduced exactly, and resetting a patch no longer throws the user out of the
  diff it was built from.
- **2026-08-19:** **the main section now changes hands between its two panes
  instead of blanking**, from two problems the user found staging a whole file
  from the side panel: the pane taking over showed nothing until its own render
  arrived (`Gui.handOverMainSection` copies the outgoing pane's content into it
  first, the trick `moveMainContextToTop` already uses across a window's tabs),
  and it came back at whatever offset it had been left at. The second turned out
  to **pre-date the whole stack** — an emptied pane kept its offset *and* its
  claim to the render it was showing — so it is fixed at the foundation of the
  lower-pane branch, demonstrate-then-fix: `clearMainView` resets the origin and
  calls the new `ViewBufferManager.ForgetRenderedContent`. One knock-on inside
  the task manager: `firstPaint` settles the scroll position **before**
  consulting a `RenderRestore`, because the restore decides where to put the
  view from what is on screen; with the cursor stored as a row, the old order
  dragged the established selection to the view's first line. Details in PR 7's
  "Review round 5"; two findings raised and left in §8. Worth knowing before
  testing anything about an emptied pane: the origin zeroing that *used* to
  happen was a race (a task reaching EOF after its view was cleared clamps the
  origin in `onEndOfInput`), so a journey that empties a pane with a task still
  in flight proves nothing. Two further findings from the round are fixed in the
  same pass: the keys for scrolling the main section asked which *window* the
  focus was in rather than which pane the section is showing, so `<pgdown>` over
  a staged-only file moved nothing; and an emptied pane stranded a pending
  `RenderRestore`, which since round 4 also strands the input block its `Done`
  balances. Signed off interactively the same day.
- **2026-08-18:** **four scroll-preservation and focus problems fixed**, three of
  them the same root cause — the selection is view lines, its meaning is buffer
  lines, and each of the three places that boundary is crossed lost something.
  The focus now follows into whichever pane a render leaves something to act on,
  asked of every render in `refreshMainViews` rather than of the one action that
  thought to ask; a resized wrapping view keeps its place in the content
  (`View.RewrapContent`, a pre-existing gocui bug fixed at the branch's
  foundation); a restored range covers the lines it is over whole; and a row
  showing both halves of a modification is remembered as both. Details in PR 7's
  "Review round 4". The two carried forward are in §8. Worth knowing before
  writing any test about a selection: `View.SelectedLines` now reports the lines
  of *content* a selection covers, and `ViewDriver.SelectedViewLineRange` is
  there for when the wrapped extent is the thing under test.
- **2026-08-17:** **two defects fixed under a non-conforming diff renderer**,
  both found by the user testing with an unpatched git: the raw fallback only
  ever bypassed stdin filters (the files panel built its command for a rendered
  diff regardless of the mode), and a re-render whose command changes reset the
  scroll to the top where no line of the old rendering could be identified. The
  second is now `ViewBufferManager.SetKeepScrollPositionForNextTask`, the coarse
  sibling of the restore. Details in PR 7's "Review round 3"; the rebase traps
  that round turned up (autosquash on by default in the user's config, the todo's
  `# ` subject prefix, `--update-refs` not moving a ref at the range's start) are
  in the section after it.
- **2026-08-16:** **PR 7 signed off**, and with it four review comments about
  the stack as a whole — how focusing the main view picks what to select, and
  how a re-render keeps your place — fixed as **mid-branch** fixups landing in
  PRs 5, 6 and 7 rather than at the tip (the user's standing rule; use
  `git rebase --onto <fixup> <target> <branch> --update-refs` so the stacked
  branch refs follow). The details are in PR 7's sign-off section; the two
  worth carrying: focusing never scrolls, and a re-render keeps the place by
  the part of the selection that is on screen, not by the selection wherever
  it is. Also learned the hard way: `just build 2>&1 | tail` hides the build's
  exit status behind `tail`'s, which is how a conflict resolution with a
  missing function got committed mid-rebase — check `$?` of the build itself,
  and grep for conflict markers before `git rebase --continue`.
- **2026-08-16:** **PR 7 implemented** (14 commits, green; §6 sign-off owed).
  Three decisions taken with the user up front: the seam is split in two so
  that copy reaches every diff panel without any panel carrying stub actions
  (deviation 2); `ctrl+w` needs no refusal, since the patch is built from a
  plain diff that never ignores whitespace (§9.10, deviation 14); and the
  real-renderer harness waits for PR 8 (§9.8). Two findings during the work,
  both raised before fixing: the plan's **timing fact for the focus-follow is
  stale** — a refresh's model update is queued and its `Then` runs after the
  render has started, so the post-op split is worked out from what the action
  did rather than read from the model (deviation 7) — and **rapid keypresses
  lost a press**, because input was held only until the model was up to date
  and not until the asynchronous re-render had moved the selection
  (deviation 10). The latter also closed a hole in PR 6: a pending restore
  outlived a render that turned out to be a message rather than a diff.
- **2026-08-15:** **PR 6 implemented and signed off** (7 commits, green; the
  interactive pass found nothing, and the whitespace consumer the user had
  suggested was singled out as worth having). Twelve deviations in the PR 6
  section; the one that matters beyond this
  PR is **3**: a pending restore now keeps the task reading to the end of its
  input, because the buffer parser refuses a partially loaded diff, so a restore
  over a rendering without OSC records can't resolve until the whole diff is in.
  Painting at the usual point makes it silently miss on any diff longer than the
  initial read (this is how the bug was found — the whitespace e2e test failed
  in one direction only), and the prototype's shape stalls at
  `LinesToRead.Total`. PR 7's post-stage reveal rides the same mechanism and
  inherits this. Also worth carrying forward: identities are matched through a
  normalized key in which every kind of content line collapses together, so an
  addition and the context line it becomes when whitespace is ignored are the
  same place — and a hunk header, which names the lines it covers, is *not* the
  same place across a `-U` change.
- **2026-08-15:** **PR 5 signed off** — the interactive pass found nothing to
  fix; a few special cases might get their behaviour refined later, but the
  user's call is to write the remaining PRs first. **PR 6 gains a third
  consumer**, at the user's suggestion: preserving position when toggling
  "ignore whitespace" (new commit 6, between the renderer-cycle and far-end
  commits). It is the first consumer whose anchor can disappear along with its
  hunk or file, so the candidate walk drops the stop-at-the-first-change-line
  rule and runs **unbounded, nearest first** — decided with the user out of
  three options (the alternatives confined it to the anchor's file, landing at
  the top of the diff or on the nearest file header when the whole file went
  whitespace-only). Also raised there and parked as §9.10: with staging in the
  main view, master's refusal to ignore whitespace in a staging context no
  longer catches it.
- **2026-08-10:** **PR 5 implemented** (9 commits, green; §6 sign-off owed).
  Two of the plan's ten commits aren't in it, both by decision with the user:
  the **jump-to-file menu is skipped** (UX undecided) and **copy moves to
  PR 7**, because raw-diff copy turns out to have a per-panel backend and so
  belongs on `FocusedMainViewActions` rather than needing a second seam of its
  own — which also means PR 5 introduces no plain-diff seam at all. The
  mandatory `IsSingleHunkForWholeFile` refinement came out cheaper than
  planned: derived per file from the already-rendered rows (no git call, no
  side-panel question), guarded by the buffer manager's loading flag. The
  twelve deviations are in the PR 5 section; the ones that bind later work are
  2 (copy in PR 7, with its span/path-list sub-decisions already taken) and 5
  (the visibility rule lives at the `refreshMainViews` chokepoint, so panels
  added later get it for free).
- **2026-08-09:** **PR 4 implemented** (7 commits, green; §6 sign-off owed).
  One scope call and one finding, both from the user. The call: a row's
  records live in `DiffLineContent.Metadata []string` rather than in a
  first-payload field plus a separate all-payloads accessor, which removes a
  two-call atomicity hazard from PR 7. The finding: **a restore must match a
  target against every record on a row, not against the row's one resolved
  identity**, or a unified→side-by-side renderer switch never finds its line —
  PR 6's matcher, written up as deviation 2.
- **2026-08-09:** **PR 2 implemented** (5 commits, green). Four scope calls
  taken with the user up front, all in the "don't land API without a consumer"
  direction: a new `DiffLineHelper` instead of `StagingHelper`, no
  `PatchLineFor*`, `DiffLineInfo` with `IsChange` only, and no gocui
  `DiffLineContent` struct until PR 4 has metadata to put in it. Two findings:
  the mandatory rename gap **doesn't reproduce** off master (§8 row closed —
  the prototype's failure came from its fork point), and git's header path
  field needs **tab** handling as well as C-quote decoding. All of PR 2's
  deviations are listed in its section.
- **2026-08-09:** **PR 1 signed off.** The interactive pass found two things,
  both fixed (deviations 8 and 9): the scroll no longer reset when a VS Code
  terminal delivered focus-in before a click, and the loading indicator now
  blanked the whole view. Also removed `FlushStaleCells` and split the `taskKey`
  data race out to its own branch below the stack. Two things worth carrying
  forward: **a data race discovered during this work gets fixed right away**,
  even a pre-existing or entirely unrelated one (separate branch when it's
  unrelated, extra commit when it isn't); and **the interactive pass is worth
  running against a repo with dirty submodules**, which is what makes a
  same-content re-render slow enough to see.
- **2026-08-08:** PR 1 implemented. The stack does **not** start from master:
  the user asked for it to be based on `fold-staging-into-main-view` (the tip
  of `scroll-selection-into-view`), which touches the same scroll code and is
  merging to master soon. Every later PR branches off its predecessor; nothing
  is pushed, and fixup commits for earlier branches go on the tip of the stack
  for the user to move down. Seven deviations from the plan are recorded in the
  PR 1 section, of which two matter later: the manager's `newContentPending`
  replaces the planned `LinesToRead.ResetOrigin` (PR 6's restore competes with
  it and must clear it — deviation 6), and the
  `screenColMax` gap fixed in PR 1 commit 9 is still live on the prototype
  branch.
- **2026-08-04:** prototype rebased onto master, past #5854 (gocui mouse
  gestures) and #5870 (diff-renderer config rework). The pre-rebase branch
  — where this plan's SHAs resolve — is kept at
  `fold-staging-functionality-into-main-view-plan`. A subject-level audit
  of the two branches found exactly two commits dropped in the rebase:
  "Report the first drag movement as a drag event, not a release"
  (correctly — absorbed by #5854) and "Carry the keyboard modifier on
  mouse click events" (accidentally). The latter was re-implemented on the
  rebased branch in the gesture-scoped shape ("Carry the press-time
  keyboard modifiers through the whole mouse gesture"); PR 10 commit 3 now
  transcribes that shape.
- **2026-08-04:** PR 3 landed on master as #5870, including a config
  restructure with a per-entry `type` field; later PRs base route/kind
  decisions on `DiffRendererConfigManager.GetDiffRendererType()` instead
  of querying the pager/ext-diff fields individually (PR 4 commit 3, PR 7
  commit 10), and `rawGit` entries are a new case for PR 7's fallback.
- **2026-08-07:** git can emit the records itself, so `rawGit` entries with
  word-diff args are now first-class conforming renderers rather than a
  case to give up on. The emitter is 4 commits on branch `osc-1717` in the
  git repo (word diffs only; byte-identical output when `OSC1717` is
  unset), unproposed upstream and shipping either way — see PR 4's
  cross-repo note. Resolved the "decide at implementation" question in PR 7
  commit 10 in favour of probing `rawGit` like anything else, which
  collapses the fallback to one rule for all renderer types. Prototyped on
  the branch in "Support git's own word diff as a metadata-emitting diff
  renderer" (+ its fixup) and "Advertise the metadata protocol to git as
  well, not only to a pager"; verified interactively: focus selects a
  change line, staging works, the selection advances after staging. Two
  traps found the hard way, both recorded above — the advertisement must
  precede `newPtyTask`'s no-pty early return, and the probe's cache
  signature must include the args.

# Focused main view — session notes

A working document capturing everything we discussed, built, and learned in this
session. It is meant as a **starting point for future sessions**, which might:

1. **Continue** solving the problem we were in the middle of (restoring scroll +
   selection when escaping back to a focused main view) — still at prototype
   quality.
2. **Enhance** the prototype with a few more missing pieces.
3. **Productionize** the whole thing: better quality, a clean commit history, and
   tests. (We did *not* make that plan; this doc gives a future session enough
   context to make it.)

> Status at end of **session 2** (the latest): branch
> `use-delta-hyperlinks-for-clicking-in-diff`. The escape/restore work that was
> uncommitted at the end of session 1 is now committed (`d901a9711` "WIP
> FocusedMainViewSnapshot approach"). Session 2 dug into the **flicker on
> escape** and landed three standalone bug fixes (see "§6" and the new commits
> on top of `e5326c3a6`); the remaining flicker is understood but not yet fully
> solved, and a small pile of **uncommitted** feature machinery
> (`keepOrigin` + the `ReadToEnd`-based restore) is left in the working tree.
> The tree builds (`just build`), is `gofumpt`-clean, and the unit tests pass.

---

## 1. The big picture: what this feature is

lazygit has a "focused main view": you press `0` (`Universal.FocusMainView`),
or click, to move focus from a side panel (files, commits, commit-files,
stash, branches, …) **into the main view** that shows its diff, so you can
scroll and interact with the diff itself. The branch builds this out into a
real interaction model:

- A **selection** can be shown in the focused main view (a highlighted line),
  toggled on demand.
- With a selection showing you can:
  - **`enter` / double-click** → dive into staging (files) or patch building
    (commits / commit-files) **for the clicked line**.
  - **`e`** → edit that line in your editor (like the staging view's `e`).
  - **`G`** → open the selected line in the current branch's GitHub PR diff
    (so you can comment on it).
- **Clicking** sets the selection at the clicked line; **double-click**
  activates (dives in). `0` focuses *without* a selection (scroll mode).

This all relies on `delta` emitting `lazygit-edit://<path>:<line>` OSC-8
hyperlinks in the rendered diff (hence the branch name); lazygit parses those
to know which file/line a view line corresponds to.

---

## 2. Branch state

Branch: `use-delta-hyperlinks-for-clicking-in-diff` (off lazygit master).

### Committed commits (most recent last), the feature-relevant ones:

```
45eebc679 Add user config gui.showSelectionInFocusedMainView
686c829d5 Press enter in focused main view when user config is on
dcd658bb7 Select line that is in the middle of the screen
c8a2bc5e7 Press enter in main view of files/commitFiles to enter staging/patch-building
0688099ee Extract some functions from CommitFilesController to a new CommitFilesHelper
a2a675fe0 Press enter in main view of commits panel to enter patch building for clicked line
673b90c10 WIP After going straight to patch building from main view, esc goes all the way back out
c4aba31c9 Replace gui.showSelectionInFocusedMainView config with on-demand selection
ee9f07a67 Press `e` in focused main view (when selection is showing) to edit that line
77157c5ad Open a browser at the selected line in the diff of the current branch's PR
30e625a8d WIP New click behavior
```

Note the two **`WIP`** commits (`673b90c10`, `30e625a8d`) — these will need
rework/squashing for productionization.

### Uncommitted work (the in-progress escape/restore feature)

```
 M AGENTS.md                                              (unrelated: see §8)
 M pkg/gui/context/patch_explorer_context.go
 M pkg/gui/controllers/commits_files_controller.go
 M pkg/gui/controllers/files_controller.go
 M pkg/gui/controllers/helpers/commit_files_helper.go
 M pkg/gui/controllers/helpers/patch_building_helper.go
 M pkg/gui/controllers/main_view_controller.go
 M pkg/gui/controllers/staging_controller.go
 M pkg/gui/controllers/switch_to_diff_files_controller.go
 M pkg/gui/types/context.go
```

---

## 3. Architecture primer (what we learned about lazygit internals)

### Contexts, the stack, and `NextInStack`

- Each panel/view is a **context**. The `ContextMgr` keeps a **stack**
  (`pkg/gui/context.go`). `Push`/`Pop` manage it. Kinds: `SIDE_CONTEXT`,
  `MAIN_CONTEXT`, popups, etc.
- Pushing a `SIDE_CONTEXT` **wipes the stack** down to just it. Pushing a
  `MAIN_CONTEXT` **evicts other main contexts** but keeps non-main ones beneath.
  Only **one main context** is ever on the stack at a time.
- A focused main view's "side panel" is found via
  **`ContextMgr.NextInStack(ctx)`** — the entry just below it on the stack.
  This was introduced on master in commit `bbd17abc43a`
  ("Add ContextMgr.NextInStack…") specifically to **stop abusing the
  parent-context mechanism** for this. Earlier prototype code on this branch
  assumed the focused main view's *parent context* was its side panel; that
  assumption is gone now — use `NextInStack`. (Memory:
  `worktree-path-vs-repo-path` is unrelated; this is a different gotcha.)

### The focused main view contexts vs. the patch-explorer contexts

`pkg/gui/context/setup.go`:

- `Normal` → `Main` view, window `"main"`; `NormalSecondary` → `Secondary`
  view, window `"secondary"`. These are `MainContext` (a `SimpleContext`).
  **This is the focused main view.**
- `Staging` → `Staging` view, window `"main"`; `StagingSecondary` →
  `StagingSecondary` view, window `"secondary"`; `CustomPatchBuilder` →
  `PatchBuilding` view, window `"main"`. These are `PatchExplorerContext`
  (also `MAIN_CONTEXT`).
- **Crucial:** `Normal` and `Staging`/`CustomPatchBuilder` share the same
  *window* but are **separate gocui views**. Only one view per window is shown
  at a time; the others are hidden **but retain their buffer (content, scroll,
  selection)**. So entering staging *hides* the `Main` view rather than
  overwriting it — its scroll/selection survive **unless something explicitly
  re-renders the `Main` view** (see "the clobber" below).

### Dispatch: `GetOnClickFocusedMainView`

- Controllers expose `GetOnClickFocusedMainView() func(mainViewName string, clickedLineIdx int) error`.
- `pkg/gui/controllers/attach.go` registers it on the context
  (`AddOnClickFocusedMainViewFn`).
- `MainViewController.enterForLine` / `onClickInAlreadyFocusedView` call
  `NextInStack(self.context).GetOnClickFocusedMainView()(viewName, lineIdx)`.
- Implementers: `FilesController` (→ staging), `CommitFilesController` (→ patch
  building), `SwitchToDiffFilesController` (commits/stash → patch building).
- The line/file is resolved from the `lazygit-edit://` hyperlink via
  `StagingHelper.GetFileAndLineForClickedDiffLine(viewName, lineIdx)` — this
  reads the hyperlink on the given **view line** (so it accounts for wrapping)
  and parses `lazygit-edit://<path>:<line>`.

### The async render-task system (`pkg/tasks/tasks.go`) — the crux of our blocker

Rendering a diff into a view is **asynchronous** and **lazy**:

- A view has a `ViewBufferManager`. `RenderToMainViews` → a **cmd task** keyed
  on the **command string**.
- The initial render reads only **`linesToReadFromCmdTask(view)` lines (one
  screenful, ~37)**, then the task **waits** on its `readLines` channel for
  more (e.g. when you scroll down, `ViewSelectionController` requests more).
- `ViewBufferManager.ReadToEnd(then)` sends `{Total:-1, Then:then}` to
  `readLines`; the loop reads to EOF, runs `onEndOfInput`, then calls `then`.
  **But** if `self.readLines == nil` (no live task), `ReadToEnd` calls `then()`
  **immediately/synchronously** — this is a premature-fire trap.
- A task's `readLines` is created **inside the task goroutine** (async), so
  right after `Push`/render the channel may not exist yet.
- `onNewKey` (`view.SetOrigin(0,0)`) runs at task start **iff the key changed**.
  Same command/key ⇒ origin preserved; different key ⇒ origin reset to top.
- `view.Reset()` (beforeStart) rewinds the write pointer; it does **not** reset
  origin. `onEndOfInput` clamps origin if the new content is shorter.
- `MainViewController.openSearch` is the existing precedent that uses
  `GetViewBufferManagerForView(view).ReadToEnd(func(){ OnUIThread(...) })`
  — but it does so on a view that's **already focused with a live task**, which
  is exactly the precondition we keep failing to establish.

### Gocui view bits we used

- `view.OriginY()` / `view.SetOrigin(x,y)` — scroll. `SetOrigin` clamps `<0`
  only (not to content length).
- `view.SelectedLineIdx()` = `OriginY + CursorY` (absolute view-line).
- `view.FocusPoint(cx, cy, scrollIntoView)` — sets cursor to absolute `cy`
  (`v.cy = cy - v.oy`); with `scrollIntoView` it adjusts origin via
  `calculateNewOrigin`. **Returns early if `cy < 0 || cy > lineCount`** — so it
  silently no-ops if the content isn't loaded that far. (This is why a deep
  selection "doesn't take" when only a screenful is loaded.)
- `view.Highlight` / `view.HighlightInactive` — whether/how the selection is
  drawn. `SimpleContext.HandleFocusLost` sets `Highlight=false` (so the
  focused-main selection is cleared whenever the view loses focus). We added
  `MainViewController.GetOnFocus` to reset `HighlightInactive=false` on the way
  back in.

---

## 4. The decided UX (don't relitigate without reason)

- **Click = point at a line ⇒ select it.** Single-click sets/moves the
  selection to the clicked line and does nothing else. **Double-click** = the
  "activate/open" gesture ⇒ dive into staging/patch building for that line.
  Clicking an unfocused view focuses **and** selects (one click → ready for
  `e`/`G`/enter). `0` focuses with **no** selection (scroll mode) — because it
  doesn't point at a line.
- **Escape from staging/patch-building should return to the focused main view
  you came from**, showing the **same main-view content** again (fresh, not
  stale), with the **same scroll position and selection**, and with the **main
  view focused** (not the side panel). One `enter` in → one `esc` out.
- For commits/stash, "the same content" means the **whole-commit diff** you were
  looking at — **not** a different focused main view (e.g. not the
  commit-files file diff). Landing on a *different* focused main view was
  explicitly rejected.
- "Stale content is out of the question" — when the underlying file changed
  (e.g. after staging), the returned main view must re-render fresh. (We accept
  that the selection may then be slightly off, since the diff changed — no fix
  planned.)

### Keybindings (focused main view, when a selection is showing)

In `MainViewController.GetKeybindings`: `Universal.Select` (space) toggles
selection; `Universal.GoInto` (enter) dives in; `Universal.Edit` (`e`) edits;
`Commits.OpenPullRequestInBrowser` (`G`) opens the PR line;
`Universal.Return` (esc) hides selection / exits. `<`/`>` are goto top/bottom
(so `G` is free).

---

## 5. The GitHub PR-line feature (working, committed `77157c5ad`)

`MainViewController.openPullRequestForSelectedLine`:

- URL form: `<pr.Url>/changes/<commitSha>#diff-<sha256(relPath)>R<line>`.
  - `<commitSha>` = `DiffableContext.RefForAdjustingLineNumberInDiff()` of the
    side panel (selected commit / the commit-files "to" ref). Using the
    specific commit's view means the right-side line numbers match what's shown,
    so **no `AdjustLineNumber` needed** here (unlike `e`).
  - `relPath` = repo-relative path via
    `filepath.Rel(RepoPaths.WorktreePath(), abs)` then `filepath.ToSlash`.
    **The anchor is `sha256(relPath)` — exact bytes, forward slashes, original
    case, no trailing newline.** (Verified empirically; the `#diff-…` hash is
    SHA-256 of the new-file path. `R<line>` = right/new side; `L` = left/old.)
- Branch resolution (`branchForPullRequest`): `commits` → `CheckedOutBranch`;
  `subCommits` → `SubCommits.GetRef().RefName()`; `commitFiles` → recurse into
  its parent context. GitHub-only (driven by `Model().PullRequestsMap`).

### GOTCHA recorded to memory

`WorktreePath()` vs `RepoPath()`: to make a working-tree path repo-relative use
`RepoPaths.WorktreePath()`, **not** `RepoPath()` — they differ in **linked
worktrees** (this dev setup uses `.worktrees/scratch`), and `RepoPath()`
silently produced the wrong relative path → wrong `sha256` anchor. See memory
`worktree-path-vs-repo-path`.

---

## 6. THE IN-PROGRESS PROBLEM (where to resume)

**Goal:** escaping staging/patch building that was entered from a focused main
view should return to that focused main view, fresh content, **scroll +
selection restored**, main view focused.

### The mechanism (now committed + a little still uncommitted)

- `types.FocusedMainViewSnapshot { SidePanel, SidePanelSelectedLineIdx,
  MainView, OriginY, SelectedLineIdx }` (`pkg/gui/types/context.go`).
- Stored on `PatchExplorerContext.focusedMainViewSnapshot` (`nil` ⇒ entered the
  normal way ⇒ plain `Pop()`), captured at entry in `focusedMainViewSnapshot(…)`
  (`main_view_controller.go`), threaded through `FilesController.EnterFile` /
  `CommitFilesHelper.EnterCommitFile`. All of this is committed in `d901a9711`.
- **Escape**: `helpers.EscapeFromPatchExplorer(c, ctx)` restores the side panel's
  selection, `Push(SidePanel)`, `Push(MainView)`, then restores origin +
  selection. The current version of this (the `ReadToEnd`-based restore plus the
  `keepOrigin` machinery below) is the **uncommitted** part left in the working
  tree — it's WIP because the flicker isn't fully solved.

### Where session 2 landed: the flicker is fully diagnosed; 3 bug fixes committed

Restoring scroll + selection on escape **works** (the final state is correct).
What remained was a **flicker on the way in**: a brief intermediate frame before
the view settles at the saved position. Chasing it uncovered three genuine,
independent bugs (all now committed on top of `e5326c3a6`):

1. **`6c7d9a295` Lock the view + guard the line index in `HyperLinkInLine`.**
   It read `v.lines`/`v.viewLines` with no `writeMutex`, racing a concurrent
   re-render, and indexed `v.lines[viewLines[y].linesY]` after only checking `y`
   against `len(viewLines)`. Because `refreshViewLinesIfNeeded` overwrites
   `viewLines` *in place without truncating*, the tail keeps stale entries
   whose `linesY` points past a shrunk `v.lines` → out-of-range panic on `enter`
   while a shorter diff was still loading.
2. **`3b31cfe01` Don't scroll a view up to fill blank space while loading.**
   The layout's scroll-up clamp ([`layout.go`], added in `6114f69ee5ef`) clamps
   a view's origin to `TotalContentHeight()` — which for a main view is just the
   **lines loaded so far**. During an async re-render that's a fraction of the
   eventual content, so it yanked the view to the top. Fix: a synchronously-set
   `ViewBufferManager.loading` flag (set in the cmd/pty wrappers *before* the
   layout pass, cleared at EOF but **not** on stop), and the layout skips the
   clamp while loading.
3. **`a4b72a6f6` Fire queued `ReadToEnd` callbacks when the initial read hits
   EOF.** The read loop processes one request at a time; the initial request has
   no `Then` and a large line count, so if the content is shorter it hits EOF on
   that request and `break`s out, abandoning any queued `ReadToEnd` request in
   the channel → its `Then` silently dropped (this was session 1's "ReadToEnd's
   `then` never fired" mystery!). Fix: drain queued requests and fire their
   `Then`s on EOF.

### Corrected diagnosis (session 1's §6 diagnosis was WRONG in its mechanism)

Session 1 said "on restore only the initial screenful (`height=37`) is loaded,
so `FocusPoint` returns early." **That was inaccurate.** The truth, confirmed by
instrumenting **every** write to the main view's `oy` (see Debug tooling §10):

- `linesToReadFromCmdTask` reads `height*(height-1)+oy` lines (≈1332+, capped at
  5000) — **not** one screenful. For typical diffs the whole thing loads quickly.
- The scroll wasn't failing because content was unloaded at *restore* time (the
  `ReadToEnd` restore, once the drain fix above made it fire, sets the final
  position correctly). It was failing because the **layout clamp** (bug #2) was
  resetting `oy` to 0 on *every layout pass* during the async load, until the
  content caught up. That is the real cause of "scroll resets to the top."

### The full origin-reset chain on escape (and how each is handled now)

Tracing every `oy` write during a commits-scrolled-down escape, three different
things were all moving the origin off the saved value:

1. **`onNewKey`** (`tasks_adapter.go`) resets `oy` to 0 when the re-render's
   command key differs from the last one (it does, because the commit-files
   render clobbered the main view on entry). → handled by
   `ViewBufferManager.KeepOriginForNextTask()` (uncommitted feature machinery),
   which suppresses that one reset.
2. **`CopyContent`** (`view.go`, via `moveMainContextToTop`) copies the
   *previous top view's* buffer **and origin** into the main view to avoid a
   blank frame. → handled by re-asserting `SetOrigin(saved)` after the pushes.
3. **The layout scroll-up clamp** → handled by bug fix #2 (the `loading` flag).

### The one remaining flicker (and the correct fix — not yet implemented)

With all three handled, the *scroll no longer jumps*. But there's still a brief
intermediate frame, and we found exactly what it is: **`CopyContent` seeds the
main view with the patch-building view's buffer**, and since we set the origin
to the saved position (far down) while that placeholder is shorter, the draw
shows the placeholder's *last line* at the top with blank below — until the pty
task finishes loading the real diff and repaints at the saved position. (It
"appears scrolled up by a varying amount" purely because what shows at the saved
`oy` depends on the patch's *length*, via `min(oy, patchLines-1)`.)

**NOTE — a rejected red herring:** "avoid clobbering the main view on entry"
does **not** fix this. `CopyContent` overwrites the main view's buffer
regardless of what was there, so preserving the original commit diff on entry
wouldn't change the placeholder frame.

**The correct fix (user's conclusion, agreed):** we're applying the saved origin
*too early*. It must be applied *exactly* when the pty task does its first
repaint (when it has read enough to fill the view at the saved scroll). The
catch: `InitialRefreshAfter` — which decides *when* that first repaint happens —
is computed from the view's `OriginY` **at task-creation time**. So the target
origin must be known at creation (so the task reads enough), but the view must
keep showing the placeholder until that first paint, and only snap to the saved
position *as part of* that paint. Concretely: **a cmd/pty analogue of
`RenderStringWithScrollTask`** — "render this command and scroll to Y once
you've read enough" — applying the origin at the `InitialRefreshAfter` refresh
rather than up front. This is the concrete next step; it's bounded but real
work, and likely lets us drop the `keepOrigin` + after-push `SetOrigin`
machinery (they'd be subsumed by the task setting the origin itself).

---

## 7. Prototype enhancements still missing (for an "enhance" session)

- **Directory case follow-up:** entering staging from a files/commit-files
  **directory** selection expands the tree and changes the selection to the
  clicked file. We restore the side panel's selected line on exit
  (`SidePanelSelectedLineIdx`) so the main view shows the directory's combined
  diff again — but we **don't restore the tree's expanded/collapsed state**, so
  the panel comes back more expanded than it was. Decide whether to restore that
  too. Also, the directory case shares the scroll/selection-restore bug above.
- **`onClickInOtherViewOfMainViewPair`** (clicking the other pane of a main
  view pair) now also selects + double-click-stages for consistency; double-check
  this is desired and that the secondary-pane paths behave.
- **Stale selection after stage/unstage:** explicitly accepted as out of scope;
  no fix planned.
- **No integration tests** exist for any of the focused-main-view interactions
  (click/double-click/`enter`/`e`/`G`/escape-restore). They were skipped on
  purpose during prototyping.

---

## 8. Productionization notes (for a future planning session — do NOT plan yet)

Context a planning session will need:

- **Commit history needs rework.** Two `WIP` commits (`673b90c10` "esc goes all
  the way back out", `30e625a8d` "New click behavior") plus the large
  uncommitted escape/restore change. AGENTS.md (this repo) mandates: small,
  self-contained, compiling, `gofumpt`-clean commits; "why not what" messages;
  prep-refactors split from behavior changes; `fixup!`/`amend!` against the
  right commit and `git rebase --autosquash`; no conventional-commit prefixes.
  The escape/restore work especially will want to be re-sequenced into clean
  commits (and the `escapeContext` → `FocusedMainViewSnapshot` evolution
  collapsed, since it was iterated heavily).
- **Demonstrate-bugs-before-fixing** pattern (AGENTS.md) with `EXPECTED`/`ACTUAL`
  — relevant if any of this lands as bug-fix-shaped commits.
- **Tests:** integration tests live under `pkg/integration/tests/...`; conventions
  in AGENTS.md (chain `t.Views().<View>()` fluently, no local view vars; use
  `stretchr/testify`). A unit-testable seam worth noting: the scroll/selection
  restore and the GitHub-anchor URL builder
  (`githubPullRequestLineURL`) are pure-ish and could be unit-tested; the patch
  index↔view line wrapping logic lives in `pkg/gui/patch_exploring/state.go`.
- **Config:** `gui.showSelectionInFocusedMainView` was added then **removed**
  (`c4aba31c9`) in favor of on-demand selection — don't reintroduce a config
  toggle for this without reason.
- **Commands:** use the `justfile` recipes (`just generate` regenerates the test
  list + cheatsheets and CI fails if stale; `just format`, `just build`,
  `just unit-test`, `just e2e-all`, `just lint`). Prefer `just` over `make`.
  Adding/renaming a keybinding ⇒ run `just generate` and commit the result
  (note: gated descriptions — the focused-main bindings use empty descriptions
  when no selection is shown, so they don't appear in cheatsheets, matching the
  existing `enter` binding).
- The unrelated `M AGENTS.md` in the working tree is the "Common commands"
  section documenting `just` — keep or commit separately.

---

## 9. Key files (quick map)

- `pkg/gui/controllers/main_view_controller.go` — the focused main view
  controller: keybindings, `toggleSelection`, `enter`/`enterForLine`, `editLine`,
  `openPullRequestForSelectedLine`, `branchForPullRequest`, click handlers,
  `showSelectionAtLine`, `focusedMainViewContextForViewName`,
  `focusedMainViewSnapshot`, `githubPullRequestLineURL`.
- `pkg/gui/controllers/switch_to_focused_main_view_controller.go` — focuses the
  main view from a side panel (`0` / click); click passes a line so it selects,
  `0` passes -1 so it doesn't.
- `pkg/gui/controllers/switch_to_diff_files_controller.go` — commits/stash →
  patch building entry (`GetOnClickFocusedMainView`, `enter`).
- `pkg/gui/controllers/files_controller.go` — files → staging entry
  (`GetOnClickFocusedMainView`, `EnterFile`).
- `pkg/gui/controllers/commits_files_controller.go` — commit-files → patch
  building entry.
- `pkg/gui/controllers/helpers/commit_files_helper.go` — `EnterCommitFile`.
- `pkg/gui/controllers/helpers/patch_building_helper.go` — `Escape` +
  `EscapeFromPatchExplorer` (the shared escape/restore logic).
- `pkg/gui/controllers/staging_controller.go` — `Escape` (calls
  `EscapeFromPatchExplorer`).
- `pkg/gui/context/patch_explorer_context.go` — `FocusedMainViewSnapshot`
  storage.
- `pkg/gui/types/context.go` — `FocusedMainViewSnapshot`, `IPatchExplorerContext`
  additions.
- `pkg/gui/controllers/helpers/staging_helper.go` —
  `GetFileAndLineForClickedDiffLine` (hyperlink parsing).
- `pkg/tasks/tasks.go` — the async render-task system (`ViewBufferManager`,
  `ReadToEnd`, the read loop) — **the thing to master to finish §6**.
- `pkg/gui/tasks_adapter.go` — string/cmd task wrappers and the origin-reset
  callbacks.
- `pkg/gui/layout.go` — the scroll-up-to-fill clamp (`setViewFromDimensions`);
  now skipped while a view's task `IsLoading()`.

---

## 10. Debug tooling (stripped from the tree; paste back when needed)

These two general-purpose debugging tools were invaluable in session 2 and were
removed from the working tree when cleaning up. They are recorded here so they
can be reapplied without re-deriving them.

### Slow down rendering (`LAZYGIT_SLOW_RENDER=<ms>`)

Stretches the async load so you can watch the frames of a re-render. Add to the
read goroutine in `ViewBufferManager.NewCmdTask` (`pkg/tasks/tasks.go`), just
before the `outer:` label, plus the per-line sleep inside the inner read loop
right after `lineWrittenChan <- struct{}{}`. Needs `os` and `strconv` imports.

```go
// DEBUG: artificially slow down rendering so transitions are visible.
var slowRenderPerLine time.Duration
if v := os.Getenv("LAZYGIT_SLOW_RENDER"); v != "" {
    if ms, err := strconv.Atoi(v); err == nil {
        slowRenderPerLine = time.Duration(ms) * time.Millisecond
    }
}
// ... and inside the inner loop, after lineWrittenChan <- struct{}{}:
if slowRenderPerLine > 0 {
    time.Sleep(slowRenderPerLine)
}
```

Run as `LAZYGIT_SLOW_RENDER=20 just debug`.

### Trace every change to a view's scroll position

Catches *who* moves `oy` (the trick that finally found the layout clamp). Add to
`pkg/gocui/view.go` (needs `os`, `runtime` imports) and call
`debugMainOriginReset(v, <newY>)` immediately before **every** write to `v.oy`:
`SetOrigin`, `SetOriginY`, `CopyContent`, `FocusPoint` (the `calculateNewOrigin`
branch), the `Autoscroll` branch in `draw`, and `ScrollUp`/`ScrollDown`. Filter
by `v.name == "main"` (or whatever view you're chasing).

```go
func debugMainOriginReset(v *View, newY int) {
    if v.name != "main" || newY == v.oy {
        return
    }
    pc := make([]uintptr, 6)
    n := runtime.Callers(3, pc)
    frames := runtime.CallersFrames(pc[:n])
    var b strings.Builder
    for i := 0; i < 4; i++ {
        fr, more := frames.Next()
        fmt.Fprintf(&b, " <- %s:%d", fr.File[strings.LastIndex(fr.File, "/")+1:], fr.Line)
        if !more {
            break
        }
    }
    if f, err := os.OpenFile("/tmp/fmvs_origin.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
        fmt.Fprintf(f, "main oy %d->%d%s\n", v.oy, newY, b.String())
        f.Close()
    }
}
```

The full session-2 diff (including these and the per-feature `FMVS` `Log.Infof`
breadcrumbs) was also saved to `/tmp/fmv-session-full.patch` during the
cleanup — though `/tmp` is ephemeral, so this section is the durable copy.

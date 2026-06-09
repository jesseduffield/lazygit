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

> Status at end of **session 3** (the latest): branch
> `use-delta-hyperlinks-for-clicking-in-diff`. The escape scroll/selection
> restore is **implemented and committed** — at normal speed there is **no
> visible flicker**. Session 3 implemented §6's proposed fix (a cmd/pty analogue
> of `RenderStringWithScrollTask`), discovered and fixed a *second* cause of the
> top-flicker (the scroll-reset loop in `refreshMainViews` ran *before*
> `CopyContent`), and corrected session 2's belief that the `onNewKey`
> suppression could be dropped (it can't — it's folded into the same mechanism).
> Full story in the new **§11**, which supersedes §6's "correct fix" subsection.
> The working tree is now **clean** (everything committed); the tree builds, is
> `gofumpt`-clean, unit tests pass, and `e2e-all` is green except one
> pre-existing **direnv-environmental** worktree test.
>
> **What's left before productionizing:** under `LAZYGIT_SLOW_RENDER` a few
> imperfect intermediate frames still appear *occasionally* — real timing races
> we agreed to investigate and eliminate (not paper over with "fine at normal
> speed"). See §11 "Remaining timing races". Memory:
> `focused-main-view-flicker-timing-races`.

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

Branch: `use-delta-hyperlinks-for-clicking-in-diff` (off lazygit master). **The
working tree is clean — everything is committed.** The branch was **rebased**
in session 3 (SHAs below are current): the `LAZYGIT_SLOW_RENDER` knob was moved
to the **base** of the branch (so it can be tested against master), and the
`SetOriginX/Y` chokepoint refactor was squashed into one commit.

### Current commit list (most recent first), `master..HEAD`:

```
625e7dbad Restore scroll and selection seamlessly when escaping to a focused main view   ← session 3
054d139fe Let a cmd/pty task restore a saved scroll position at its first paint           ← session 3
7f547a5a3 Reset other main views' scroll after copying content, not before                ← session 3
fe79d18b6 Route all view origin writes through SetOriginX and SetOriginY                  ← session 3 (chokepoint refactor; candidate for master)
89e6f6b14 Session notes: corrected flicker diagnosis and the 3 bug fixes
86f4b3486 Fire queued ReadToEnd callbacks when the initial read reaches EOF               ← session 2 bug fix
b7470af27 Don't scroll a view up to fill blank space while its content is loading         ← session 2 bug fix
788d959ad Lock the view and guard the line index when reading a hyperlink                 ← session 2 bug fix
63221c3dd Session notes
5f500893a WIP FocusedMainViewSnapshot approach                                            ← WIP (needs rework)
207927e0d WIP New click behavior                                                          ← WIP (needs rework)
385d2e9dd Open a browser at the selected line in the diff of the current branch's PR
c5dd8ddc6 Press `e` in focused main view (when selection is showing) to edit that line
55922f81a Replace gui.showSelectionInFocusedMainView config with on-demand selection
877812c6a WIP After going straight to patch building from main view, esc goes all the way back out  ← WIP (needs rework)
0088f26c1 Press enter in main view of commits panel to enter patch building for clicked line
ec50f3122 Extract some functions from CommitFilesController to a new CommitFilesHelper
ed2015cac Press enter in main view of files/commitFiles to enter staging/patch-building
1e5f31dd6 Select line that is in the middle of the screen
fff7a0d19 Press enter in focused main view when user config is on
8a26bebbb Add user config gui.showSelectionInFocusedMainView
ed48988a9 Add LAZYGIT_SLOW_RENDER debug knob for watching async render frames               ← base; candidate for master
```

The three **`WIP`** commits and the heavily-iterated `FocusedMainViewSnapshot`
machinery will need re-sequencing for productionization (see §8). The two
clearly-standalone, master-worthy commits (`ed48988a9` slow-render at the base,
`fe79d18b6` the `SetOriginX/Y` chokepoint) are deliberately isolated so they can
be cherry-picked off.

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

### The one remaining flicker (and the correct fix — IMPLEMENTED in session 3, see §11)

> **Update (session 3):** the fix described below was implemented, but the
> diagnosis here was *incomplete* in two ways that §11 corrects: (a) the
> `onNewKey` suppression could **not** be dropped, and (b) there was a **second**
> source of the top-flicker — the scroll-reset loop in `refreshMainViews`. Read
> §11 as the current truth; the text below is session 2's understanding.

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
  `ReadToEnd`, the read loop). Session 3 added `ScrollToOriginYForNextTask` /
  `GetScrollToOriginYForNextTask`, `LinesToRead.ApplyInitialScroll`, the
  first-paint apply in the read loop, and the `onNewKey` suppression (§11). Also
  hosts the committed `LAZYGIT_SLOW_RENDER` knob.
- `pkg/gui/tasks_adapter.go` + `pkg/gui/pty.go` — cmd/pty task wrappers; both now
  peek the manager's pending scroll and pass it to
  `linesToReadFromCmdTask(view, targetOriginY)` (`view_helpers.go`).
- `pkg/gui/main_panels.go` — `refreshMainViews` (the scroll-reset loop, **now
  after** `moveMainContextPairToTop`, §11) and `moveMainContextToTop` →
  `CopyContent`.
- `pkg/gui/layout.go` — the scroll-up-to-fill clamp (`setViewFromDimensions`);
  skipped while a view's task `IsLoading()`.
- `pkg/gocui/view.go` — `SetOriginX`/`SetOriginY` are now the **single
  chokepoints** for all `ox`/`oy` writes (`fe79d18b6`); ideal breakpoint spot.

---

## 10. Debug tooling

### Slow down rendering (`LAZYGIT_SLOW_RENDER=<ms>`) — now COMMITTED

This is no longer a paste-back snippet: it's committed at the **base** of the
branch (`ed48988a9`). Sleeps `<ms>` after each line written to a view, so the
frames of an async re-render become visible. No effect when unset. Run as
`LAZYGIT_SLOW_RENDER=40 just debug` (with `just print-log` in another tab).
**This is the tool that makes the remaining timing races (§11) visible** — they
are essentially invisible at normal speed.

### Trace every change to a view's scroll position — now a single chokepoint

Session 3's `SetOriginX`/`SetOriginY` refactor (`fe79d18b6`) routed **every**
write to `v.oy`/`v.ox` through `SetOriginY`/`SetOriginX`. So you no longer need
to scatter the tracer across `SetOrigin`/`CopyContent`/`FocusPoint`/`draw`/
`ScrollUp`/`ScrollDown` — **set one breakpoint (or one log line) inside
`SetOriginY` in `pkg/gocui/view.go`** and you catch all of them, with the
`bt`/Call-Stack giving the caller. (This is exactly how session 3 found the
`refreshMainViews` reset-loop cause — see §11.) The old multi-site
`debugMainOriginReset(v, newY)` helper still works if you want a `/tmp` log with
a trimmed call stack; drop it into `SetOriginY` and filter by `v.name`:

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

---

## 11. Session 3: the flicker fix (implemented) + remaining timing races

Session 3 turned §6's proposal into working, committed code, and corrected the
diagnosis twice along the way. At **normal speed the escape is now flicker-free**.

### What "applying the saved scroll at first paint" became (commit `054d139fe`)

The cmd/pty analogue of `RenderStringWithScrollTask`, driven by one field on
`ViewBufferManager`:

- **`ScrollToOriginYForNextTask(originY int)`** sets `scrollToOriginYForNextTask
  *int`. The escape calls it on the main view's manager **before** the re-render
  is triggered. It has *two* effects on the next cmd/pty task:
  1. **Suppresses the start-of-task origin reset** (`onNewKey`) — so the
     `CopyContent` placeholder keeps showing at *its* scroll instead of being
     yanked to the top. (This is the part session 2 thought we could drop. We
     can't — see below.)
  2. **Sizes the initial read to `originY`** (`linesToReadFromCmdTask(view,
     targetOriginY *int)` uses it instead of the view's current `OriginY`) **and
     scrolls there at the first refresh** via a new `LinesToRead.ApplyInitialScroll`
     callback, applied once (guarded by `sync.Once`) — at the `InitialRefreshAfter`
     point, and in the EOF branch *before* `onEndOfInput` (so a now-shorter diff
     gets clamped back into range).
- The field is **peeked** (`GetScrollToOriginYForNextTask`) by the cmd/pty
  wrappers (`tasks_adapter.go`, `pty.go`) to size the read, and **cleared in
  `NewTask`** after the `onNewKey` decision — so it survives long enough to drive
  both effects, and applies to exactly one task. (Per-view managers, so the
  secondary view isn't affected.)
- Behaviour-preserving until a caller sets it.

### Escape wiring simplified (commit `625e7dbad`)

`EscapeFromPatchExplorer` now just calls `ScrollToOriginYForNextTask(snapshot.OriginY)`
before the pushes, and restores the **selection only** (`FocusPoint` + highlight)
via `ReadToEnd` once the diff is fully loaded. The session-2 dance — up-front
`SetOrigin`, after-push `SetOrigin`, and `KeepOriginForNextTask` — is **gone**;
the task owns the scroll now.

### Correction #1: `onNewKey` suppression could NOT be dropped

§6 predicted the new mechanism would let us drop the `onNewKey` suppression. It
didn't. `CopyContent`'s entire purpose is that the newly-revealed view keeps
showing the previous view's content **at its scroll** ("as if nothing changed")
until the real content paints. Letting `onNewKey` reset that to the top *is* a
flicker. So the suppression is kept — folded into the same
`scrollToOriginYForNextTask` field (effect #1 above) rather than a separate
`keepOrigin` flag.

### Correction #2: there was a SECOND cause of the top-flicker — `refreshMainViews` (commit `7f547a5a3`)

Even with `onNewKey` suppressed, the placeholder still flicked to the top under
slow render. A `SetOriginY` breakpoint (trivial now, thanks to the chokepoint
refactor) caught it: `refreshMainViews` (`main_panels.go`) reset the scroll of
every *other* main view at the **top** of the function — i.e. it zeroed the
patch-building view's origin **before** `moveMainContextPairToTop` →
`CopyContent` copied that view (now at origin 0) into the Normal view. So the 0
came from the reset feeding `CopyContent`, *independent of* `onNewKey`.

**Fix:** move the reset loop to **after** `moveMainContextPairToTop`. End state is
unchanged (every other main still ends at 0, and the destination always
re-renders), but `CopyContent` now copies the source at its real scroll, so the
placeholder stays put. This also makes *every* cross-pair transition's
placeholder seamless, not just our escape.

### Verification

- `just build` / `just lint` / `just unit-test` all green. (`TestNewCmdTaskInstantStop`
  is a **pre-existing timing flake** that only trips under the full suite's
  parallel load; passes 10/10 in isolation, and the session-3 task changes are
  inert on its instant-stop path.)
- `just e2e-all`: green **except** `worktree/associate_branch_rebase`, which
  fails *environmentally* — `cd`-ing into the linked worktree triggers lazygit's
  direnv integration to pop a "Press <enter> to run 'direnv allow'" confirmation
  (this checkout's `.envrc` is blocked), stealing focus from the `.Focus()`
  assertion. Run `direnv allow` (or confirm it fails the same on `master`).

### Remaining timing races (DO THIS before productionizing)

At normal speed there's no visible flicker, but under `LAZYGIT_SLOW_RENDER`
**occasional** imperfect intermediate frames remain. The user's explicit call:
these point to **real timing races** in the async render/scroll path, and we
should *eliminate* them rather than rely on normal timing masking them. Not yet
characterised — next session should:

- Reproduce under `LAZYGIT_SLOW_RENDER` (try a range of values; the races are
  intermittent) across the three transitions (files→staging, commit→patch-building,
  and the escape), all while **scrolled down**.
- Use the single `SetOriginY` chokepoint + `bt` and/or the §10 tracer, plus the
  `ReadToEnd`/`InitialRefreshAfter`/`ApplyInitialScroll` ordering, to pin which
  interleavings produce a bad frame. Suspects worth scrutinising: the ordering
  between the task's first `ApplyInitialScroll` paint and the `ReadToEnd`-driven
  selection restore; the `afterLayout`-deferred pty task creation racing a layout
  pass; and `CopyContent` vs. the task's first write.
- Memory: `focused-main-view-flicker-timing-races`.

---

## 12. Restore by patch identity, escape routing, and the plan to solve the three entangled problems

A design discussion (it did **not** result in code; the implementation is for a
**new session**). Three things came out of it: a better model for the escape
restore, a set of escape-routing cases the current prototype gets wrong, and a
decision about *when* to solve the hard async problems.

### 12.1 The restore should anchor on a patch identity, not a numeric scroll/index

The current escape (§6/§11) saves the main view's `OriginY` + selection **index**
and replays them. That's only right when the content is unchanged. But the whole
reason you'd escape *after doing something* in the staging/patch view is that you
changed the content — you staged a hunk, or `d`-dropped one — so a numeric index
now points at a different line. (§4 already conceded "the selection may be
slightly off; no fix planned.")

The right model: on escape, read the explorer view's **current selection at that
moment** as a patch identity `(file, type, source-line, side)` — *after* the
host's auto-advance has already moved it to a still-valid line — then find that
identity in the freshly re-rendered main view and scroll it into view. This is
the **inverse** of the diff-line-metadata primitive (identity → rendered row,
rather than row → identity), and it is the **same operation** as the `-U`
context-change scroll-preservation consumer (see diff-line-metadata-notes.md §1
item 5). It also lets the `FocusedMainViewSnapshot` store *less* (the `OriginY` +
main-view `SelectedLineIdx` become derivable), which is a simplification.

Note the difference from the `-U` case's "anchor on the nearest surviving change
line" fallback: that's for when context lines genuinely vanish. For staging-escape
the host's auto-advance normally hands us a valid line directly; the
nearest-surviving fallback is only for the degenerate cases below.

### 12.2 Escape routing — which half to return to, or whether to at all

The escape target is really `(context, identity)`: which focused-main half
corresponds to where the explorer's selection ended up. The current prototype
gets several cases wrong:

- **Stage a non-last hunk** — selection auto-advances within unstaged → restore
  there. The common, already-conceptually-fine case.
- **Stage the *last* unstaged hunk** — unstaged becomes empty and the selection
  crosses to the **staged** panel → escape should land in the **staged half**
  (the secondary view), not unstaged. *Currently lands in the files panel — bug.*
- **`<tab>` between unstaged/staged inside the staging view** — escape should
  return to the half matching the side you're on. *Currently broken — bug.*
- **Drop the last unstaged hunk with nothing staged** — the staging view
  auto-closes; focus goes to the files panel. *Correct as-is*: the main view
  would only show "No changed files", so there's no point focusing it. (Probably
  not deliberate in the prototype, but it's the right outcome.)
- **Drop a hunk in the custom patch builder (dropping it from the commit), or
  "Remove patch from original commit" from the custom-patch menu** — the patch
  builder *always* closes (it can't mutate the commit from within it), not just on
  the last hunk. Re-focusing the main view is right, but here there is **no host
  auto-advance**, so we'd have to advance the selection ourselves. Since these are
  infrequent, acceptably focusing the **side panel** instead is a fine shortcut.

So the return target is: unstaged-selection → main; staged-selection → secondary;
no valid selection / empty → files panel; custom-patch-builder → side panel
(shortcut) or self-advance.

### 12.3 Decision: solve the three entangled problems in the prototype (next session)

The new restore mechanism, the **§11 timing races**, and the
**`BufferLineForViewLine` staleness trap** (diff-line-metadata-notes.md §8) are
entangled: the identity-based restore reads/parses the buffer *while it is still
loading*, which is exactly where both the races and the staleness bite. We
decided **not** to defer these to productionization. A production plan written
around three unsolved, entangled mechanisms isn't a plan; the prototype exists to
retire that unknown cheaply (build-order §7). Resolve them here, and
productionization becomes transcription.

Attack order (dependency-first, to keep it from ballooning):

1. **§8 staleness fix first.** A bounded *correctness* fix with a known shape
   (snapshot `viewLines`+`lines` under one lock, or tie the read to the task that
   produced the buffer). Needed regardless, and it's what makes
   reading/parsing-during-load safe — which everything else depends on.
2. **Characterize the §11 races early, timeboxed, before designing around them.**
   The one real sink risk is that they're not yet characterised: we don't know if
   each is a specific bad interleaving (bounded fix) or a fundamental tension in
   the async-lazy-render model (e.g. no flicker-free first paint without buffering
   a screenful first). That distinction changes the restore design *and* is itself
   a key production-plan input, so pin it cheaply up front. This is "understand the
   race to choose the right fix", not "decide whether to fix it" (we fix it).
3. **The new restore mechanism**, which splits:
   - **sync half** — the §12.2 routing (which context / side panel); largely
     independent of timing, can progress in parallel.
   - **async half** — a *predicate* scroll: generalize
     `ScrollToOriginYForNextTask(y)` (§11) to "scroll to the row matching this
     predicate, applied the first refresh at which it's satisfiable" (the read
     loop scans the lines read so far via the metadata primitive; the fixed-`y`
     case is then just a trivial predicate). This is the "keep parsing as it
     loads" piece and is the part that sits on §8+§11.

### 12.4 Memory

`focused-main-view-flicker-timing-races` already covers the races. The escape
restore reworking and the three-problem plan above are recorded here.

---

## 13. Session 4: characterizing the §11 timing races

Session 4 did §12.3 step 1 (the §8 staleness fix — see diff-line-metadata-notes.md
§8) and this characterization (step 2). The conclusion drives part 3.

### 13.1 Why the faithful repro is the *interactive* app, not a headless test

The §11 method (LAZYGIT_SLOW_RENDER + the `SetOriginY` chokepoint tracer via
`just debug` / `just print-log`) is inherently interactive, and that is not an
accident of tooling — the headless integration harness **cannot reproduce the
pager-path races** for two structural reasons:

- **cmd vs. pty path.** With a real pager configured (delta), the main-view
  re-render goes through `newPtyTask` (`pkg/gui/pty.go`), which marks the view
  loading synchronously but **defers the actual task creation to `afterLayout`**.
  The default test setup has no pager, so it takes `newCmdTask`
  (`tasks_adapter.go`), which creates the task **synchronously**. The
  `afterLayout` deferral is one of the race ingredients (§13.3), so the cmd path
  doesn't exercise it.
- **env allowlist.** `pkg/integration/components/env.go` passes only a strict
  allowlist (`PATH`, `TERM`, `HOME`, git config) to the lazygit subprocess, so
  `LAZYGIT_SLOW_RENDER` doesn't even reach it without patching the harness.

So the characterization below is **code-grounded and mechanically conclusive**,
not visually confirmed in this session. The fixes it prescribes still want a
confirming tracer run on the interactive app (real pager, scrolled down, a range
of slow-render values). The two races are *mechanisms in the code*, not guesses
about flaky values — which is the standard the two prior misdiagnoses (§11
corrections #1/#2) failed to meet.

### 13.2 The escape origin/content timeline (commit→patch-building→escape, scrolled to S)

`EscapeFromPatchExplorer` (`patch_building_helper.go`), with the main view's saved
scroll = S and the patch view's current scroll = P:

1. `ScrollToOriginYForNextTask(S)` — sets the manager's pending-scroll field
   (not a view write yet).
2. `Push(SidePanel)` — focusing the side panel runs its `onRenderToMainFn`, which
   calls `RenderToMainViews(Normal pair)`. Inside `refreshMainViews`
   (`main_panels.go`):
   - `RefreshMainView(Normal)` → `runTaskForView` → `newPtyTask`/`newCmdTask`:
     `StartLoading()` (sync) and read the pending scroll; the pty path **defers
     task creation to `afterLayout`**.
   - `moveMainContextPairToTop(Normal)` → `moveMainContextToTop(Normal)` →
     **`CopyContent(topView = the patch view, Normal)`** → `Normal.oy = P`, and
     Normal's buffer becomes the *patch* content. (`onNewKey`'s top-reset is
     suppressed because the pending-scroll field is set.)
   - the reset loop sets the *other* mains (incl. the patch view) to oy 0 —
     **after** CopyContent, so it no longer feeds 0 into the copy (the §11
     correction #2 fix).
3. `Push(MainView=Normal)` — just focuses/highlights; **no second re-render**.
4. Next layout pass → `afterLayoutFuncs` drained → the pty task is actually
   created (`NewTask` → goroutine → creates `readLines`, reads lines).
5. At `InitialRefreshAfter` (= height + S + 10 lines read) the task's
   `ApplyInitialScroll` sets `Normal.oy = S` and does its first refresh — the
   first paint that shows the real commit diff at S.
6. The escape's `OnUIThread(…)` → `ReadToEnd(restore)` restores the **selection**
   (`FocusPoint`, no origin change) once the diff is read.

The flicker window is **between step 2 (oy=P, buffer=patch placeholder) and step 5
(oy=S, buffer=commit diff)**. Layout-driven draws that land in it show whatever is
in `(oy, viewLines)` at that instant.

### 13.3 The two races, classified

**Race A — partial content drawn at the placeholder scroll during load
(content/scroll). BOUNDED interleaving over a FUNDAMENTAL constraint.**

Between the task's first write (step 4) and its first official paint (step 5),
`oy` is still P and the buffer is being overwritten from the top with the commit
diff. Layout passes call `draw`, and because each write sets `tainted`, the draw's
`refreshViewLinesIfNeeded` rebuilds the view lines from the *partial* commit diff
(plus the retained patch tail past `freshViewLineCount`). So a draw in this window
shows a Frankenstein frame: the **top of the commit diff at the wrong scroll P**
(or partial content / stale patch tail), before it snaps to S. Intermittent
because whether a layout pass lands in the window, and at what load fraction,
depends on scheduling; at normal speed the load finishes before a layout pass
lands there.

- **Bounded:** the *specific* fix is to keep the placeholder **coherent** until the
  first official paint, instead of letting partial content leak in. The task
  already controls *its* refreshes; the leak is that *layout* draws also rebuild
  from the half-written buffer. Suppressing the view-line rebuild **only while a
  scroll-restore is pending its first paint** (gated on the pending-scroll field,
  so the normal load-from-top case is untouched) keeps the draw showing the
  retained placeholder at P until step 5 swaps atomically to the commit diff at S.
  No scroll jump, no Frankenstein frame.
- **…over a fundamental constraint:** you still see *something* during the load —
  with the fix it's the coherent placeholder (the patch) rather than a broken
  frame. There is **no flicker-free first paint of the target content at S without
  first having that screenful buffered.** This is the key production-plan input the
  §12.3 step-2 question asked for: the model can only choose *what* the pre-paint
  frame shows, not eliminate it. The two coherent choices:
  - (a) keep the outgoing view's content (the patch, via CopyContent) as the
    placeholder — coherent but the *wrong* content shown briefly; or
  - (b) **don't clobber the destination's own buffer.** On escape, the Normal
    view's buffer *already held the commit diff at S* (it was only hidden when we
    entered the patch explorer, never cleared — §3). `CopyContent` overwrites that
    perfectly-good content with the patch. If `moveMainContextToTop` did **not**
    copy over a view that already holds appropriate content, and the re-render kept
    that old buffer as the placeholder until its first paint (Reset already retains
    the view lines for flicker-avoidance), the placeholder would *be* the commit
    diff at S and the swap to fresh content would be invisible. (b) is the strong
    lever and supersedes §6's "avoid clobbering doesn't help" — that was about
    *entry*; on *escape* the destination's stale buffer is exactly what we want.

**Race B — the selection restore can fire before the task is live (selection).
BOUNDED bad interleaving.**

The selection restore is scheduled as `OnUIThread → ReadToEnd(restore)`.
`ReadToEnd` fires its callback **synchronously and immediately when
`manager.readLines == nil`** (tasks.go) — the documented premature-fire trap (§3).
The escape's re-render task is created by `Push`, but its `readLines` channel is
created later, inside the task goroutine (and only *after* it has stopped the
previous task, which blocks on `notifyStopped`). The restore is deferred one UI
tick to let the task become live — but `OnUIThread` (`g.Update`) runs on a later
main-loop iteration whose ordering against the task goroutine reaching
`readLines = make(...)` is **not guaranteed**. If the tick wins the race,
`ReadToEnd` fires `restore` immediately, `FocusPoint(SelectedLineIdx)` runs before
the content is loaded that far and **no-ops** (it returns early when
`cy > lineCount`), and the selection is **silently not restored**. Intermittent.

- **Bounded:** thread the after-load callback **into the task's lifecycle** the
  same way the scroll restore is (a field on the manager that the next cmd/pty task
  folds into its initial `LinesToRead.Then`, set at creation), instead of a
  separate post-hoc `ReadToEnd`. Then it always fires when *that* task's initial
  read completes — never against a nil channel. This is the same one-mechanism
  ("do X after the re-render") part 3 needs for its identity-based restore, so it
  is foundational, not throwaway. Ordering is safe: `ApplyInitialScroll` (oy=S)
  fires at `InitialRefreshAfter`, which is < the initial read total, so the scroll
  is already applied when the `Then` runs.

### 13.4 Implications for part 3 (the restore rework)

- **One restore mechanism, threaded into the task.** Both the scroll restore
  (already done via `scrollToOriginYForNextTask` → `ApplyInitialScroll`) and the
  selection restore should ride the *same* next-task hook, eliminating the separate
  `ReadToEnd` and Race B. Part 3 generalizes the scroll half to a **predicate**
  ("scroll to the row matching this identity, at the first refresh it's
  satisfiable"); the selection restore becomes "once loaded far enough, focus that
  same identity's row." Build them as one "after the re-render, reconcile to this
  identity" callback.
- **Prefer (b): reuse the destination's own buffer as the placeholder.** Part 3
  should investigate not clobbering the focused main view's retained buffer on
  escape, which makes the common "content unchanged" escape genuinely
  flicker-free (invisible swap) and reduces Race A to the rare "content actually
  changed" case. (Freshness still requires the re-render; the win is the
  *placeholder* being right.)
- **Race A's coherence fix (suppress rebuild while a scroll/predicate restore is
  pending its first paint) is gated on the pending field**, so it's safe for the
  normal case and carries into the predicate generalization unchanged.

### 13.5 Status / what's left

> **Update — the Race A fix was reworked.** The session-4 first cut used the
> gated `holdViewLines` flag (suppress the view-line rebuild during a
> scroll-restore load). The user judged that — and the §8 `freshViewLineCount`
> guard — as patches that work *around* a broken invariant rather than restoring
> it: many readers (`BufferLines`, `ViewLinesHeight`, the diff-line readers, …)
> would each have to know which buffer to trust. Both patches were **reverted**
> and replaced by an **off-screen render** (see below). `holdViewLines` and
> `freshViewLineCount` no longer exist.

- Race B (selection premature-fire): **fixed**. The restore rides the re-render
  task via a `thenForNextTask` hook on the buffer manager, folded into the cmd/pty
  task's initial-read `Then`; the `OnUIThread → ReadToEnd` dance is gone. (Unchanged
  by the rework.)
- Race A (coherence during load) **and** §8 (stale-tail mapping): **fixed together**
  by the off-screen render. A cmd/pty task now `BeginOffscreenRender()`s — writes go
  into a second `viewBuffer` (`View.offscreen`) while the displayed buffer, and so
  everything every reader sees, stays the previous render. At its first-paint point
  (`InitialRefreshAfter`, or EOF for short content) the task `SwapInOffscreenRender()`s:
  the off-screen buffer becomes the displayed buffer in one atomic step and the saved
  scroll is applied in the same step. So no reader ever sees a half-written buffer at
  the wrong scroll (Race A), and because the swap is a wholesale replace,
  `refreshViewLinesIfNeeded` now **truncates** the view lines — the stale tail (§8)
  can't form. `clear()`/`Reset()` abandon any in-progress off-screen render so a
  synchronous `SetContent` after a stopped task writes to the display. The mechanism
  is **unified** (every async render uses it), not gated on scroll-restore. The swap
  holds `writeMutex` for now; per [[main-thread-over-mutexes-direction]] it could
  later become a main-thread bounce.
- The (b) no-clobber lever (don't `CopyContent` over the focused main view's own
  retained buffer on escape, so the placeholder *is* the right content) is **not**
  done — still the stronger part-3 improvement for the common unchanged-content
  escape.

Prep history (clean commits): extract `viewBuffer` → make the writing methods
operate on a `viewBuffer` → revert the two patches (separately) → off-screen render.

> **Verification caveat (important):** all of this was reasoned through and is
> covered by gocui unit tests (`TestOffscreenRender`, `TestBufferLineForViewLineStaleTail`)
> + green `e2e-all`, but the **flicker behaviour was NOT visually confirmed** —
> session 4 ran headless, and the faithful repro is the interactive app (§13.1: the
> headless harness uses the cmd path and blocks `LAZYGIT_SLOW_RENDER`). Confirm with
> `just debug` + `just print-log`, scrolled down, across a range of
> `LAZYGIT_SLOW_RENDER` values, for all three transitions.

**Part 3 interplay (carry forward):** with the off-screen render, the *loading*
content lives in `View.offscreen`, while `v.buf`/`viewLines` and the view-line
readers (`DiffLineMetadataInLine` / `bufferLineForViewLine`) keep describing the
*displayed* (previous) render. Part 3's predicate scroll, which inspects rows *as
they load* to decide when to swap, must therefore scan the **off-screen** buffer's
cells (a new accessor), and the swap point generalizes from "InitialRefreshAfter
lines" to "the first read at which the predicate is satisfied on the off-screen
content". The §8 #1 two-call atomicity constraint (diff-line-metadata-notes.md §8)
applies to that scan.

### 13.6 Follow-ups discovered after the off-screen render landed

- **Scrollbar regression (FIXED — session 5).** With the
  off-screen render, the swap happens at the *first-paint* point
  (`InitialRefreshAfter` = `height + oy + 10` lines), but the scrollbar is sized
  from the *displayed* buffer's height (`ViewLinesHeight()` →
  `calcRealScrollbarStartEnd` in `gui.go`). So at the swap the displayed buffer is
  only a viewport-and-a-bit tall, while the task keeps reading up to
  `linesToReadForAccurateScrollbar` (`min(height*(height-1)+oy, 5000)`, far more);
  the scrollbar thumb therefore jumps to reflect the short content and then grows
  back to its former position as the rest loads. The old incremental mechanism
  masked this: the non-truncating `viewLines` tail kept the height at the *previous*
  render's value until EOF, so the bar stayed put. Reproduces clearly under
  `LAZYGIT_SLOW_RENDER=2` on the files panel's 10s auto-refresh while scrolled down
  (main content stays stable — good — but the bar flickers).

  **Resolution (session 5).** Not something more fundamental about the *swap* —
  the off-screen render correctly fixed content coherence. The actual point: the
  scrollbar reads a **strictly later quantity** than the viewport-fill paint. The
  first paint needs only `oy + height` lines (enough to show the scroll position);
  an accurate scrollbar needs the *total* height, known only near end-of-read. So
  no single early swap can have both the content and the scrollbar right — they
  depend on different read amounts. Delaying the swap to the scrollbar-accurate
  point (the first candidate direction) was **rejected**: it would regress
  *flicking latency* for large diffs (the new diff wouldn't appear until ~1332
  lines read instead of a screenful) and the *no-placeholder first render* (a
  blank/"loading…" view until the full read, since the off-screen render only shows
  *old* content when there is some). Instead we took the second direction: **hold
  the scrollbar height while a load is in progress.** `FreezeScrollbarHeight`
  (gocui) records the view's height when the load begins (captured at
  `StartLoading`, while the view still shows the previous render — for escape this
  is *before* `CopyContent`, so it's the retained commit-diff height, i.e. the
  correct final height); the scrollbar is sized from `max(displayed height,
  scrollbarHeightFloor)` until `UnfreezeScrollbarHeight` releases it at EOF
  (`onEndOfInput`). `clear()` also releases it, so a synchronous string render
  superseding a still-loading diff (main views can render either, see
  `runTaskForView`) doesn't leave a stale floor. This is the **same class** as the
  layout scroll-up clamp, which already ignores the partial content height while
  `IsLoading()` — only one reader (the scrollbar calc) sees the floored height, all
  other `ViewLinesHeight()` callers stay oblivious ([[isolate-new-concepts-from-clients]]).
  Landed as an **`amend!` against the off-screen-render commit** (`1fd1325c2`),
  since that commit introduced the regression (AGENTS.md: don't sequence a branch
  so an earlier commit regresses and a later repairs); the message gained a
  paragraph on the scrollbar. Covered by `TestScrollbarHeightHeldWhileLoading` and
  `TestScrollbarHeightReleasedWhenContentReplaced` (both fail on the un-floored
  code). **Verification caveat:** mechanically unit-tested and `e2e-all` green, but
  the visual flicker fix was **not** confirmed interactively this session (headless
  can't repro, §13.1) — still wants a `LAZYGIT_SLOW_RENDER=2` `just debug` pass,
  scrolled down, across the 10s refresh, escape-from-staging, commit→patch-building,
  and plain commit-flicking.
- **10s auto-refresh vs. the identity-based restore (note for productionizing
  part 3).** The files panel (and others) do a periodic background refresh
  (default 10s) that re-renders the main view by starting a *new* task. This can
  fire just after escaping staging, while the escape's re-render task is still
  reading toward the restore target (the saved scroll/selection now, the target
  patch line in part 3). The refresh's task **stops and replaces** the escape task
  and knows nothing about the restore — `scrollToOriginYForNextTask` /
  `thenForNextTask` were already consumed (cleared in `NewTask`) by the escape
  task, so the replacement only preserves the current scroll and the
  identity-restore is silently dropped. Not urgent (the window is small), but
  productionization of part 3 must handle it — e.g. let the pending restore survive
  task replacement, or have the refresh re-assert it.

  **Session 5 update — it's worse than "restore dropped".** Testing the scrollbar
  fix under `LAZYGIT_SLOW_RENDER=2`/`20` with autoRefresh **on**, the escape↔refresh
  overlap produces visible **rendering artifacts / glitches**, not merely a dropped
  restore — two re-renders (escape + the periodic refresh) of the *same* view
  interleaving, made far more likely by slow render. With autoRefresh **off** there
  is no flicker at all; without slow-render the artifacts have **not** been
  reproducible even once. The open decision (raised with the user): do we harden
  this in the prototype or defer to production? The pivotal unknown is **which kind**
  of problem it is, and we don't know yet:
  - **(a) a soundness hole in the off-screen-render / task lifecycle** — two
    re-renders of one view leaving *shared `View` state* half-applied: a stopped
    escape task can leave `v.offscreen` non-nil (stopped between
    `BeginOffscreenRender` and the swap); the pty `afterLayout` deferral can queue
    *two* task-creations from one layout cycle, and the manager's pending fields
    (`scrollToOriginYForNextTask`, `thenForNextTask`, and now the
    `FreezeScrollbarHeight` floor) are read at different instants by each, so the
    refresh task can re-capture the floor / re-read the pending scroll mid-escape.
    If this is the cause it would surface beyond escape (any two rapid re-renders of
    a view), so it's a defect in the mechanism we just built — **must** be fixed in
    the prototype (AGENTS.md: a known race is not "live with it").
  - **(b) "merely" the dropped restore** (the issue already described above) plus the
    consequent old-content-then-jump being drawn — in which case the *mechanism* is
    sound and the real fix is part 3's "one restore that survives task replacement",
    where a throwaway prototype patch would likely be discarded.

  Lean (session 5): **characterize now, timeboxed** (§12.3 step 2 — "understand the
  race to choose the right fix"), then let the finding pick the fix's home. (a) gets
  fixed here; (b) is confirmed as a part-3-owned item and deferred. The
  characterization itself is cheap and is exactly the kind of entangled unknown the
  prototype exists to retire ([[resolve-hard-unknowns-in-prototype]]).

  **Characterization done (session 5) — the dominant artifact was a third thing,
  simpler than either (a) or (b).** Method: throwaway `[CHAR]` logging at every task
  lifecycle point (render request, `NewTask` consume/bail/run, swap@first-paint,
  swap@EOF, stop) + interactive repro under `LAZYGIT_SLOW_RENDER=20`,
  `refreshInterval=3`, scrolled down. The logs showed the glitch even with **no
  escape present** — it's the periodic refresh re-rendering the main view faster than
  it loads, overlapping *itself*. Smoking gun: `swap @EOF i=1 stopped=true` on a
  ~147-line diff — a **stopped task taking the end-of-input branch and finalizing**
  (swap its 1-line off-screen buffer in, clamp the origin to 1, clear `loading`). Root
  cause: when a task is stopped, `opts.Stop` closes *and* the scanner closes
  `lineChan`, so the read loop's `select` is a coin-flip; ~half the time a stopped
  task lands in the `!ok`/EOF branch instead of the stop branch. Pre-existing (the EOF
  branch always clamped via `onEndOfInput`), but the off-screen render made it far
  worse by also swapping a truncated buffer in. **FIXED** (commit "Don't run
  end-of-input handling for a render that was stopped"): the EOF branch now re-checks
  `opts.Stop` and bails like the explicit stop case. Confirmed by the follow-up log:
  every stopped task now logs `EOF-but-STOPPED breaking clean`, no more truncated
  swaps, and the content stays put through refresh churn. Not deterministically
  unit-testable (the bug *is* the non-deterministic select; a test would be flaky on
  the old code), so no test — a rare justified skip of demonstrate-the-bug-first.

  My original **(a)** theory (the `afterLayout` double-queue letting a refresh consume
  the escape's pending scroll/`then`) was **not** what the logs showed: escape tasks
  consistently consumed their *own* `scroll`/`then` (`consumedScroll=true
  consumedThen=true`). What remains of the escape race is plain **(b)**: the escape
  task is stopped mid-read by a refresh *before* it reaches its first-paint swap (seen
  as `target=93 thenPending=true` consumed, then the task stopped at i<InitialRefresh),
  so the scroll/selection restore never applies. Still open, still part-3-owned (the
  restore must survive task replacement — re-assert on the new task, or let the pending
  state outlive the stopped task). Decide with the user whether to do a bounded version
  now or defer with part 3.
- **Flicker-avoidance has inherent limits at staging transitions (record, not
  fix).** Two cases (user-observed, session 5) where painting the old content as a
  placeholder can't be seamless, so blanking might actually read better than showing
  the old content in the wrong place:
  - **Layout changes across the transition.** Going in/out of staging often changes
    the layout (single main view ↔ split main+secondary). The content jumps position
    regardless of any placeholder, and the old content may even *wrap differently* in
    a split view than unsplit — so keeping it shown buys little; blanking when the
    layout changes may be cleaner.
  - **The escape highlight/selection sequence.** On escape we first remove the
    staging selection highlight (keeping the old content), then later paint the main
    view (no selection), then later still paint the restored selection. With a large
    hunk selected (a big highlighted region) this three-step reveal is pronounced
    under slow-render; barely noticeable at normal speed — but possibly visible on a
    slow machine.

  Neither needs fixing in the prototype, and likely not in production either;
  recorded for completeness.

Memory: `focused-main-view-flicker-timing-races` (updated — off-screen render),
`isolate-new-concepts-from-clients`, `main-thread-over-mutexes-direction`.

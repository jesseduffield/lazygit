# Focused main view — session notes

> **⚠️ THIS BRANCH (`use-delta-hyperlinks-for-clicking-in-diff`) IS A THROWAWAY
> PROTOTYPE.** It will never land. When we productionize, we start over from
> scratch on a new branch and re-implement cleanly with tests. So **do not sweat
> the commit history here**: no need for fixups against the right commit, no
> squashing, no avoiding "added then reworked" churn, no agonizing over superseded
> commits. Just commit working, green increments as you go. (The general
> [[clean-history-no-back-and-forth]] preference is *suspended* for this branch;
> it applies again when productionizing.) Commits should still compile and pass
> tests — that's it. The value of this branch is the *knowledge* captured in this
> doc, not the code or its history.

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

---

## 14. Session 6: the identity-based restore (part 3)

Built §12.3 step 3: the escape restore now anchors on a **patch identity** read
from the explorer's current selection, found in the re-rendered main view by
scanning its content as it loads — replacing the numeric scroll/index replay.
Items 1 (async predicate scroll) and 3 (#1, survive task replacement) of the
part-3 plan are done and turned out to be **one mechanism**, as §13.6 predicted.
Item 2 (escape routing, §12.2) is partially done; the (b) no-clobber lever is
analysed but not done. Commits (most recent last):

1. *Resolve a diff line's identity from a content snapshot, not the live view* —
   prep refactor. `StagingHelper.GetDiffLineInfo` now resolves the identity from a
   `[]gocui.DiffLineContent` snapshot (text + metadata + hyperlink per unwrapped
   buffer line) via one buffer-agnostic resolver, instead of the displayed view's
   per-view-line readers. Behavior-preserving; lets the same resolver run against
   the *loading off-screen* buffer (the inverse consumer).
2. *Add gocui accessors for scanning a loading off-screen render* —
   `OffscreenDiffLineContents` (the loaded off-screen rows), `ViewLineForBufferLine`
   (inverse of `BufferLineForViewLine`, to turn a matched buffer line into the view
   line to scroll/select), and `OffscreenLineCount`. Unit-tested.
3. *Restore the focused main view by patch identity on escape* — the core.
4. *(routing)* — see §14.3.

### 14.1 The mechanism (items 1 + 3 are one design)

- **`tasks.RenderRestore{FirstPaintReady func() bool; Apply func()}`** replaces the
  numeric `scrollToOriginYForNextTask`/`thenForNextTask` pair and
  `LinesToRead.ApplyInitialScroll`. The read loop, when a restore is set, asks
  `FirstPaintReady()` after each line instead of the `InitialRefreshAfter` count;
  the first time it's true it swaps the off-screen render in and calls `Apply()`.
- **`StagingHelper.RestoreFocusedMainViewOnEscape`** builds the restore. It reads
  the explorer's selected `(file, type, source-line)` *now* (forward primitive on
  the explorer view), then:
  - `FirstPaintReady`: scan the off-screen rows (`OffscreenDiffLineContents`) for
    the first row matching the identity (`FindDiffLine` → `SamePatchLine`,
    comparing path + `PatchSelectLine`); once found *and* a screenful below it has
    loaded (`OffscreenLineCount`), ready. This gives an **early** mid-load swap for
    backends that resolve a row on its own — OSC metadata and lazygit-edit
    hyperlinks.
  - `Apply`: convert the matched buffer line → view line, `FocusPoint(…, scrollIntoView=true)`
    (scroll + select in one step), set `Highlight`, and clear the pending restore.
- **Buffer-parse can't resolve incrementally — it must wait for the complete diff
  (verification finding, no pager).** The first interactive test showed the
  restore worked with the patched (OSC-metadata) delta but failed *every time*
  with no pager: the main view landed at the top with no selection. Cause: the
  buffer-parse backend parses whole hunks and checks them against their `@@`
  lengths (`Patch.IsWellFormed`); while the diff streams in, the trailing hunk is
  incomplete, so the parse is rejected as not-well-formed. The incremental scan
  checks each line once, *as it loads* (i.e. while it's the last line and the diff
  is still partial), so it rejected every line and never found the target — whereas
  metadata/hyperlinks resolve per-line and worked. **Fix:** keep the incremental
  scan (it's what gives metadata/hyperlinks their early swap), but have `Apply`
  re-scan the **complete** content once the off-screen render is swapped in (for
  buffer-parse the swap happens at end of input, so the displayed buffer is whole
  and well-formed by then). So buffer-parse restores correctly, just with a
  later (EOF) swap rather than an early one. Wasteful detail to fix in production:
  the incremental scan still does (failing) buffer-parse work on every partial line
  (~O(n²)); pin the backend, or skip incremental scanning when it's buffer-parse.
- **`FocusPoint(scrollIntoView)` unifies scroll and selection** at the first
  paint. This is a deliberate deviation from the task's literal "thread the
  selection restore through `thenForNextTask`": §13.4 also said "build them as one
  callback", and folding both into one `Apply` at the swap removes the separate
  selection hook entirely (and with it Race B's premature-fire window — there is no
  post-load `Then` for the selection any more). `LinesToRead.Then` survives only
  for `ReadToEnd` (search).
- **#1 (survive task replacement) falls out of the identity model**, not a separate
  fix. `restoreForNextTask` lives on the `ViewBufferManager` and is **not cleared
  when a task starts**, so a periodic refresh that stops the escape's re-render
  before first paint just hands the still-pending restore to its replacement task.
  It is **not gated on the command key** — and that gate, which the first cut had,
  was actively wrong: returning after staging the last unstaged hunk re-renders
  `git diff` as `git diff --cached`, a different key, yet the line to land on is
  right there in the new content. The restore **validates itself** instead: its
  scan finds the target line only if the content still contains it, so applying it
  to a genuinely different item (a different file ⇒ different path) is a harmless
  no-op. A task clears it once it has applied it — found *or not* — in `Apply`, so
  it lives for exactly one re-render: it survives stop-and-replace, but doesn't
  re-apply on every later render. This is the "clear on content change, keep across
  same-content replacement" of §13.6, achieved by the identity rather than the key.
- **The snapshot stores less** (§12.1): `FocusedMainViewSnapshot` dropped `OriginY`
  and `SelectedLineIdx` (and the `clickedLineIdx` thread into it) — the position is
  derived from the explorer's live selection at escape.

### 14.2 Verification status — first interactive pass done; broader sign-off remains

`just build` / `unit-test` / `e2e-all` / `lint` / `format` all green; new unit
tests cover the gocui scan primitives and the read-loop restore mechanism. The
agent couldn't drive the full-screen TUI (no tmux), so the **user** ran the first
interactive pass (`LAZYGIT_SLOW_RENDER=5`, `refreshInterval: 3`):

- **Works well.** In/out of staging with the patched (OSC-metadata) delta — no
  flicker, no artifacts.
- **No-pager was broken, now fixed.** See §14.1: buffer-parse can't resolve a
  partial diff, so the incremental scan missed; `Apply` now re-scans the complete
  content. Confirmed fixed.
- **Hunk-mode escape lands on the *first* line of the hunk** now, not the last:
  the cursor sits at the selection's end in hunk/range mode, so the escape anchors
  on `State.SelectedViewRange()`'s start instead of the view cursor.
- **Pre-existing crash fixed (incidentally):** hovering during a re-render could
  panic in `findHyperlinkAt` (`index out of range … length 0`) — `onMouseMove`
  read `viewLines` without `writeMutex`, racing the rebuild. Now locked. On master,
  but the off-screen render's extra task-goroutine rebuilds widened the window.

Second pass (user, slow render): the session-5 sign-offs were **confirmed** — the
scrollbar and diff content both stay put across the 3 s refresh, no glitches or
truncated frames. Two further issues surfaced, both on content *change* (not the
same-content refresh):

- **Switching items scrolled the old content to the top before the new swapped in
  — FIXED.** With the off-screen render the previous content stays displayed until
  the swap, but the scroll-to-top reset fired when the task *started*, so the old
  content jumped to the top first. Deferred the reset to the first paint (see
  "Reset the scroll to the top at first paint"); the old content now stays at its
  scroll until the new content takes its place.
- **Scrollbar starts large then shrinks while a *new* diff loads — pre-existing,
  not fixed.** Confirmed on master: the first paint is one screenful, and the thumb
  is sized from the displayed height, which grows toward the full count as more
  lines load → the thumb shrinks. `linesToReadForAccurateScrollbar` doesn't prevent
  this; it only keeps the thumb stable *while scrolling* (it reads enough that
  scrolling further won't resize it), not during the initial load. `FreezeScrollbarHeight`
  hides it for a *same-content* re-render (held height matches) but can't for a new
  diff (held height is the old one), and the new total is unknown until ~EOF. Only
  visible under slow render; judged not worth addressing.

Still pending: the slow-render matrix under the **hyperlink** backend is being
skipped — the user leans to dropping that backend in production (§14.5).

### 14.3 Escape routing (§12.2) — partial

Done: the snapshot is now recorded on **both** staging halves (`FilesController.EnterFile`)
and cleared on both at escape, so staging the *last* unstaged hunk — which pushes
`StagingSecondary` (confirmed in `RefreshStagingPanel`: `mainState == nil &&
!secondaryFocused` ⇒ `Push(secondaryContext)`) — no longer escapes to the files
panel. It returns to the focused main view (`snapshot.MainView`, i.e. `Normal`),
and the identity restore lands on the staged line there: with the last unstaged
hunk gone the file has only staged changes, so the files panel renders the staged
diff into `Normal` (no split), which *contains* that line.

Remaining (deferred — intricate, and a wrong move regresses the core staging
escape; do with full understanding + interactive verification):
- **Return to `NormalSecondary` for the split + `<tab>`-to-staged case.** When the
  file has *both* unstaged and staged changes the files panel splits (`Normal` =
  unstaged, `NormalSecondary` = staged). If you tab to the staged half and escape,
  the staged-line identity is in `NormalSecondary`, but we currently return to
  `Normal` and set the restore on `Normal`'s manager, so the scan won't find it
  there. The fix is to choose the target focused-main context from the **escaping
  explorer half's window** (`context.GetWindowName()`: secondary ⇒ `NormalSecondary`)
  — but it must reckon with the post-action **split state** (the staged content is
  in `NormalSecondary` only when a split survives; otherwise it's in `Normal`), so
  it's not a pure window-name map.
- **Custom patch builder drop → side panel** (§12.2): the builder always closes and
  there's no host auto-advance, so the accepted shortcut is to focus the side panel
  rather than self-advance. Not implemented; patch-builder escape currently runs the
  same identity restore against `Normal`.

### 14.4 (b) no-clobber lever — analysed, NOT done

§13.4 framed (b) as "skip `CopyContent` over the focused main view's retained
buffer on escape, so the placeholder *is* the right content and the swap is
invisible". Skipping the copy is easy (gate `moveMainContextToTop`'s `CopyContent`
on a pending restore via a `HasPendingRestore` check). **But it isn't sufficient on
its own:** `refreshMainViews` resets every *other* main pair's origin to 0 when a
pair is moved to the top, so entering the patch explorer already zeroed `Normal`'s
`oy`. So on escape `Normal`'s buffer is the right *content* but at scroll 0, not at
the saved scroll S — the swap+`FocusPoint` would still visibly scroll from the top.
A genuinely invisible unchanged-content swap needs the placeholder kept at S too:
either don't zero the focused-main pair's origin on entry, or re-establish it before
the re-render. Left for a follow-up; §13.3's "(b) supersedes (a)" still holds as the
target, just with this extra requirement. (Without (b), the escape currently shows
the `CopyContent`ed patch as a coherent placeholder until the swap — option (a).)

### 14.5 Known prototype limitations (productionization input)

- **Match fidelity under the hyperlink backend.** `SamePatchLine` compares path +
  `PatchSelectLine` (source line + is-deletion). The dev config's default pager is
  delta with `--line-numbers --hyperlinks` → the main-view scan uses the
  **hyperlink** backend, which can't tell the side (`DiffLineOther`), so a
  **deletion** selected in the explorer won't match a scanned row and the restore
  won't find its line (scroll/selection then just isn't restored — no crash). #2's
  OSC metadata (the "patched delta" pager) gives full fidelity. Additions/context
  match fine. Acceptable for the prototype; **productionization lean (user,
  session 6): drop the hyperlink backend entirely** rather than support it
  imperfectly — rely on OSC metadata (the planned delta patch) and have users on an
  old delta update their pager. That also removes the side-less `DiffLineOther`
  identity, so every match carries a real side and the deletion case just works.
- **Scan cost.** `FirstPaintReady` re-snapshots the off-screen buffer each line
  until the target is found (then it's O(1) via `OffscreenLineCount`), and
  `diffLineInfoFromContents` attempts the buffer-parse backend (building a texts
  slice) even when the hyperlink backend will answer — so the search phase is ~O(n²)
  in the loaded line count. Invisible in the common unchanged-content escape (the
  target is at the saved scroll and the swap is coherent regardless), bounded for
  shallow targets, but a deep target in a large changed diff would lag. Productionize
  with an incremental scan (don't rebuild the snapshot; resolve the newest line; pin
  the backend once).
- **No fallback when the identity can't be read** (e.g. the explorer's selected line
  doesn't parse): the re-render then has no restore and resets to the top. In
  practice the auto-advance lands on a parseable change line, so this is rare.

Memory to update: `focused-main-view-flicker-timing-races` (identity restore +
#1 landed; interactive sign-off still pending), `resolve-hard-unknowns-in-prototype`
(the entangled async unknown — scanning during load — is retired), and a new note
that part-3 routing (§12.2) and the (b) lever are the remaining prototype work.

---

## 15. Next steps agreed at the end of session 6

The core async unknown is retired, so productionization is *transcription-ready*
(decomposition sketch below). But we're **deferring productionization** deliberately
to do more prototyping first — not because we need to learn more, but to make the
prototype **compelling for pitching the OSC protocol to pager developers** ("here
are the lazygit features you unlock by implementing it"). File/hunk navigation is
the strongest such demo: it's impossible without the metadata once a pager
restructures the diff.

**Next prototype work (each a new session):**

> **Update (session 7):** steps 1 and 2 below are **done** — see §16. Remaining: step
> 3 (side-by-side delta) and the OSC spec draft.
>
> **Update (parallel session, 2026-06-10):** step 3 is **done** too — see §17. The
> side-by-side prototype shows **v1 needs no format change** (correcting the guess
> below that it "will likely add payload/format"); only the OSC spec draft remains.
>
> **Update (difftastic, 2026-06-10):** a further pager — **difftastic** — is now
> prototyped in **both** its modes (side-by-side + inline), the categorical
> #2-only case. Details and findings in diff-line-metadata-notes.md §10. v1 still
> holds, but difftastic surfaced a **token-vs-line model mismatch** (§10.2) the
> unified-diff pagers hid — the spec draft should speak to it. Emitter lives on
> `prototype-osc-metadata` in `/Users/stk/Stk/Dev/Builds/difftastic`.

1. **Preserve scroll/selection when `{`/`}` change the `-U` context size** (diff-line
   consumer #5). **[done — §16.1]** Reuses the identity-restore machinery directly: capture the
   nearest *change* line (survives context changes, unlike context lines), re-render,
   restore it — the escape restore's sibling, triggered by the context-size change
   instead of escape.
2. **File/hunk navigation in the focused main view.** **[done — §16.2]** Agreed model:
   - **keys**: `n`/`N` = next/previous file; `<left>`/`<right>` = next/previous
     hunk — the **exact bindings the staging view already uses** for hunk nav (so
     no horizontal-scroll clash to worry about; the focused main view's `<`/`>` stay
     goto-top/bottom).
   - **"hunk" means lazygit's notion, not git's.** Not the `@@`-delimited section,
     but a **block of consecutive added/deleted lines separated by context** — there
     can be several within one git hunk. These are exactly the blocks the staging
     panel jumps between (`State.SelectNextHunk`/`SelectPreviousHunk` /
     `patch.GetNextChangeIdx`), and the main view should match.
   - **anchor** = the selected line if a selection is showing, else the **top
     visible line** (works for change-blocks too — the next block is the next run of
     change lines below the anchor, found from the per-line metadata `type`).
   - **target** = scan rendered rows via the diff-line primitive for the start of
     the next/previous change block (file nav: where metadata `file` changes).
   - **effect** = selection showing → move the selection to the target + scroll into
     view (like the staging view); **no selection → scroll the target to the top and
     do *not* create a selection** (stays in scroll mode). [decided]
   - "Previous" when the anchor is mid-block mirrors `State.SelectPreviousHunk`
     semantics.
   - #1 and #2 share the "scan rows for a metadata boundary/identity" core, so they
     pair in one session (#5 = restore-across-rerender path; nav = immediate path).
3. **Side-by-side delta mode** — prototype in **parallel, in the delta repo**
   (doesn't touch lazygit). Feeds the spec (below). **[done — §17; v1 suffices, no
   format change]**

**OSC spec write-up:** **DONE — draft written** (`diff-line-metadata-osc-spec.md`
at the worktree root). It's a draft for pager-dev *feedback*, not final. The
unified single-column wire format is validated (§9.2), side-by-side is validated
(§17 — v1 needs no addition), and the draft speaks to difftastic's token-vs-line
mismatch (diff-line-metadata-notes.md §10.2). The **OSC number is resolved to
`1717`** after a terminal-allocation audit (diff-line-metadata-notes.md §3.4 +
the spec appendix); `456` is retired as the placeholder. The `456`→`1717` rename
across the prototype code (delta/difftastic/gocui/lazygit + the env var, now
`EMIT_OSC1717_METADATA`) is **done** — builds + metadata unit tests green in each
repo, and fresh delta/difftastic release binaries built. **Still pending:**
circulating the draft for feedback.

**Decisions locked (session 6):** concurrency stays **mutex-based**, including for
productionization (the main-thread-mutation rework is a separate, later effort — do
not design around it now); the **(b) no-clobber** lever is **dropped** (escape is
already flicker-free); **§12.2 split+tab routing** is **deferred to productionization**.

**Productionization decomposition (sketch, for when we get there):** land as a
bottom-up stack of independently-reviewable PRs — (a) the gocui async-render
improvements (off-screen render + scrollbar hold + stopped-task fix + mouse-move
lock, several master-worthy on their own); (b) diff-line primitive #1 (host-side
buffer parse, serves no-pager/`--color-only`); (c) the focused-main-view feature on
top of #1; (d) escape-restore-by-identity; (e) #2 OSC reader behind an "if the pager
emits it" gate, gated externally on the delta patch being upstreamed. Drop the
hyperlink backend (§14.5). The delta dependency means the feature can't fully ship
to delta users until the OSC patch lands — decide whether to ship #1-only first.

---

## 16. Session 7: the two showcase consumers (#4 hunk/file nav, #5 context preserve)

Built the §15 step-1 and step-2 prototype consumers — the ones that make the OSC
pitch concrete, because both are impossible once a pager restructures the diff
unless it emits per-line metadata. Both reuse the diff-line primitive (the
`diffLineInfoFromContents` resolver and the gocui scan accessors from §14), so they
paired in one session as predicted. Commits (most recent last), all green on
`build`/`unit-test`/`e2e-all`/`lint`:

1. *Extract the identity-based restore into a context-neutral helper* — prep
   refactor. `restoreDiffLinePositionOnRerender(view, target, place)` is the escape
   restore's scan/swap machinery (`RenderRestore` + `FindDiffLine` +
   `OffscreenDiffLineContents`/`ViewLineForBufferLine`) with the *positioning* split
   out behind a `place(viewLine)` callback. Behavior-preserving: escape passes the
   same FocusPoint-and-select closure the inline `Apply` used.
2. *Preserve the diff scroll and selection when the context size changes* —
   consumer #5.
3. *Navigate the focused main view by file and hunk* — consumer #4 (+ a fixup on #2,
   the visibility guard below).

### 16.1 Consumer #5 — preserve scroll/selection across an `-U` context change

`ContextLinesController.applyChange`'s `default` branch (the focused main view and
every side panel, which all render their diff into the **"main"** window view) now
calls `StagingHelper.PreserveDiffPositionOnRerender(Contexts().Normal.GetView())`
right before `HandleRenderToMain()`. The increase/decrease-context keybindings
re-render with a different `git diff -U<n>` command, so the command key changes and
the render *would* reset to the top; the preserve installs a restore that suppresses
that (same `restore == nil` gate on `ResetOrigin` the escape uses) and re-establishes
the position.

- **Anchor on the nearest *surviving* line, preferring the anchor itself** (refined
  after interactive feedback — see §16.5). We'd like to keep the anchor line (the
  selection if one shows, else the top visible line), but it may not survive: a
  context line vanishes when the context shrinks. So `nearbyDiffLines` captures a
  *prioritized candidate list* — the anchor first, then outward (at-or-below
  preferred), each direction stopping at the first change line (a guaranteed survivor,
  so the list always contains one). `restoreDiffLinePositionOnRerender` was
  generalized to take that list and land on the first candidate the re-render still
  contains: the anchor itself if it survived (so a still-present context line stays
  put, or stays selected), else the nearest surviving line, minimizing scrolling.
  (The single-candidate case is exactly §12.1's "nearest surviving change line"; the
  escape path passes a one-element list.)
- **Offset-preserving placement.** `place` puts the landed line back at the *same
  screen row* it was on (`SetOrigin(0, viewLine - row)`, `row` clamped into the
  view), so the view barely moves — unlike the escape's `FocusPoint` (centre-if-off),
  which is right for "navigate to" but would jump for "stay put". A showing selection
  is re-established on the line (`FocusPoint` scrollIntoView=false + `Highlight`);
  with no selection the view stays in scroll mode.
- **Early swap only for the anchor; fallback resolved at EOF.** The candidate list
  isn't in load order (a nearer candidate can load after a farther one), so the
  incremental `FirstPaintReady` only early-paints when the *primary* (anchor)
  candidate is reachable; any fallback is resolved against the complete content at the
  guaranteed EOF swap (`firstPaint` runs `Apply` at end of input regardless). For a
  context change — a "stay put" redraw — the slightly later swap when the anchor
  didn't survive is imperceptible.
- **Scope: the "main" window view only** (covers the focused main view *and* the
  side-panel diffs — they share the render path, which is the clean generalization
  the task asked for). The secondary focused main view (`NormalSecondary`) is **out
  of scope**: changing context while focused on it still jumps. Noted, not fixed.
- **Visibility guard (fixup on #2).** `PreserveDiffPositionOnRerender` bails if the
  view isn't `Visible`. Otherwise, pressing `{`/`}` while a non-diff main context
  occupies the "main" window (merge conflicts) would set a restore on the hidden
  `Normal` view that never re-renders, so it would linger and wrongly suppress a
  *later* render's scroll reset. (In the reachable cases — side panel or `Normal`
  focused — `Normal`'s view is the visible one and re-renders immediately, consuming
  the restore, so this only bites the merge-conflicts edge.)

### 16.2 Consumer #4 — file/hunk navigation in the focused main view

`MainViewController` gains four bindings, exactly the staging view's hunk keys plus
`n`/`N` for files: `<left>`/`h` = previous hunk, `<right>`/`l` = next hunk (reusing
`Main.PrevHunk`/`NextHunk`), `n` = next file, `N` = previous file (literal
`config.Keybinding{...}`, no new userConfig field — matching the branch's
config-removal trend). The focused main view's `<`/`>` stay goto-top/bottom.

- **"Hunk" = lazygit's change block, not git's `@@`.** A run of consecutive
  added/deleted lines separated by context. `AdjacentChangeBlock` resolves each
  rendered row's type via the diff-line primitive, builds an `isChange` bool slice,
  and the pure `changeBlockStart` mirrors `State.SelectNextHunk`/`SelectPreviousHunk`
  line-for-line (previous-from-mid-block goes to the *previous* block's start).
- **File nav targets where the metadata `file` changes.** `AdjacentFile` builds a
  per-row `paths` slice; the pure `fileStart` finds where the path changes, then backs
  up over the neighbouring file's *untagged* header rows (`backUpOverHeader`) so both
  directions land on the **top of the file** — the `diff --git`/`@@` header when the
  buffer is parseable, or the pager's file-header rows when only content carries
  metadata. This back-up is the bit that's impossible without the metadata: once delta
  restructures, the host can't otherwise tell which rows belong to which file.
  - **The anchor's file is found by scanning *down* (`anchorFilePath`), refined after
    feedback (§16.5).** Landing on a file header puts an *untagged* row at the top, so
    the next nav's anchor doesn't itself carry a path. The row whose path identifies
    the file you're "in" is the first tagged row at-or-below the top (the file's
    content), not the nearest tagged row in either direction — which would be the
    *previous* file's content just above the header, making a second `n` jump back
    into the file just left instead of advancing.
- **Effect (decided model, §15):** anchor = selection if showing, else top visible
  line. Selection showing → move the selection to the target and **scroll it into
  view** (`showSelectionAtLine` gained a `scrollIntoView` param; clicks still pass
  false). No selection → **scroll the target to the top** (`SetOrigin`) and do **not**
  create a selection — stays in scroll mode.
- The `changeBlockStart`/`fileStart`/`anchorFilePath` index arithmetic is pulled into
  pure functions and **unit-tested** (`diff_line_navigation_test.go`), covering both
  the parseable (every row tagged) and restructuring-pager (only content tagged)
  shapes, including navigating from a middle file's untagged header.

### 16.3 Verification status — first interactive pass done (§16.5); headless green

`build`/`unit-test` (incl. the new pure-scan tests)/`e2e-all`/`lint` all pass. As in
§13.1/§14.2 the headless harness can't repro the pager paths (cmd path, no
`afterLayout` deferral, `LAZYGIT_SLOW_RENDER` blocked) and there is no
focused-main-view e2e harness to extend, so the user ran the **interactive** pass
(`just debug` + `just print-log`, scrolled, patched-delta and no-pager). The first
pass worked except for two issues, both now fixed — see §16.5. Re-verify those after
the fixes, plus the standing checks:
- #5: `{`/`}` keeps the view roughly put (offset-preserving) instead of jumping to the
  top, both in scroll mode and with a selection; a still-present context line now stays
  put / stays selected (§16.5).
- #4: `<left>`/`<right>` jump hunks, `n`/`N` jump files (landing on the file header);
  selection-vs-scroll-mode effect matches the staging view; under metadata delta the
  file nav works where buffer-parse can't (delta default mode), including repeated `n`
  across several files (§16.5).
- The hyperlink backend is being **dropped** (§14.5) — ignore it.

### 16.4 Known limitations (productionization input)

- **Scan cost.** Both consumers resolve *every* rendered row through
  `diffLineInfoFromContents` per action (`PreserveDiffPositionOnRerender` via the
  outward `nearbyDiffLines`, the nav via full `isChange`/`paths` slices), and the
  buffer-parse backend re-parses per call → ~O(n²) per keypress. Fine for a prototype
  keypress; productionize with an incremental scan and a pinned backend (same note as
  §14.5's restore scan).
- **Nav only sees loaded content.** The nav scans the *displayed* buffer, so a target
  beyond the lazily-loaded portion isn't found (no-op). The initial read pulls
  ~`height*(height-1)` lines, so typical diffs are fully loaded; a target deep in a
  very large diff would need a `ReadToEnd` first (like `openSearch`).
- **`NormalSecondary` context-change not preserved** (§16.1).

### 16.5 First interactive feedback (fixed)

The user's first interactive pass found the navigation solid with a plain diff and
hunk nav solid under patched delta, plus three behaviours to refine — all fixed:

- **#5 anchored too eagerly on a change line.** Original: always restore the nearest
  *change* line, even when the anchor was a context line still in the patch — so the
  view jumped to the next `+`/`-` line past surviving context, and a selection on a
  context line moved off it. **Fix:** `nearbyDiffLines` builds a prioritized candidate
  list (anchor first, outward, stopping at the first change line each side), and the
  restore lands on the nearest *surviving* one — keeping the anchor line itself when it
  survives (context line stays put / stays selected, increase or decrease), only
  falling back outward when it doesn't. Generalized `restoreDiffLinePositionOnRerender`
  to a candidate list (escape passes one); fallback is resolved at the EOF swap since
  candidates aren't in load order. See §16.1.
- **#4 file nav stalled under patched delta on the second `n`.** After the first `n`
  landed on a file header (an *untagged* row under delta), the next `n`'s anchor row
  carried no path; the old `nearestPath` then picked the *previous* file's content just
  above the header and "advanced" back into the file just left. **Fix:** `anchorFilePath`
  determines the anchor's file by scanning **down** first (the content at/below the
  top), not the nearest tagged row in either direction. No delta change needed — the
  metadata on content lines is enough. New unit cases cover the middle-file-header case.

Commits: a `fixup!` on #4 (`anchorFilePath` + tests) and an `amend!` on #5 (the
candidate-list anchoring; the message was updated since "anchor on a change line" no
longer described it). The #5 amend! also carries the shared-helper generalization
(`restoreDiffLinePositionOnRerender` → candidate list), which strictly belongs to the
prep-refactor commit but can't be split out non-interactively (same file, and the
context-change caller it adapts didn't exist at the prep commit) — fold or re-split as
preferred.

---

## 17. Side-by-side delta prototype (parallel session) — feeds the OSC spec

§15 step 3: extend delta's per-line metadata emission to **side-by-side** mode. Done in
the delta repo (`prototype-osc-metadata`, commit `bbec9b5`), parallel to session 7 and
independent of lazygit. This is the input the OSC spec draft was waiting on for its
second open item — and the answer turned the open item into a *closed* one.

### 17.1 What it does — per-cell attachment

Unified delta funnels every content line through one emit point (`Painter::paint_lines`),
so the metadata rode there. Side-by-side doesn't: it paints each visual row as **two
panel halves** via `Painter::paint_line` (singular), so the OSC is attached **per cell**:

- **left half → the minus line's identity** (`d`, carries new+old numbers),
- **right half → the plus line's identity** (`a`, carries new),
- **a context line — shown in both halves → the same `c` record before each half.**

Empty counterpart cells (the blank half of a pure add/delete) and wrapped continuation
rows carry **no** OSC. So a changed line — a paired minus/plus on one visual row — carries
**two** records on that row: `d` before the left gutter, `a` before the right.

> **Correction (after the difftastic prototype, see diff-line-metadata-notes.md §10.8):**
> "wrapped continuation rows carry no OSC" is a **bug**, not a feature, whenever the
> *pager* does the wrapping (delta in side-by-side mode, difftastic side-by-side):
> each wrapped row is a distinct host buffer line, so `e`/`enter`/hunk-nav break on
> the un-tagged continuations. The fix — emit the line's record on **every** wrapped
> output row — was applied to **both** difftastic and delta (delta re-emits the
> primary's record without advancing its line-number counters). See
> diff-line-metadata-notes.md §10.8.

### 17.2 The verdict: v1 needs NO addition for side-by-side

This **resolves** the spec's side-by-side open item and **corrects** §15's guess that the
parallel prototype "will likely add payload/format" — it doesn't. The v1 wire format
(`version ; type ; new-line ; old-line ; file`) suffices, **with no column/side
discriminator**, because:

- **`type` already implies the column**: `a` (added) is inherently the new/right side, `d`
  (deleted) the old/left side. The host needs no extra field to know which column an
  `a`/`d` cell belongs to.
- **context (`c`) is symmetric**: the same record is emitted before both halves, so the two
  columns are indistinguishable *by payload* — but they don't need to be, because they're
  the same logical line. The host tells the columns apart **by position**, which it does
  anyway.

A side discriminator would only earn its keep to disambiguate the `c` case, and that case
needs no disambiguation. So: don't add one to v1.

### 17.3 The one latent v1 gap (shared with unified, not SxS-specific)

Side-by-side makes one v1 limitation *visible* that unified hid: **context and addition
records carry no old-file line number** (the `old-line` field is empty for `c` and `a`). In
side-by-side the old/left column is right there on screen, so the temptation to read an
old-file number off a left-column cell is real — but the metadata doesn't carry it.
Concretely, with a net insertion above a context line so old≠new:

```
{{OSC 1;c;3;;b/b.py}}│  2 │ctx_b   {{OSC 1;c;3;;b/b.py}}│  3 │ctx_b
```

the left gutter shows old line **2**, but both halves' OSC says `new=3`. A host reading the
left cell gets the *new* number, not the old.

This is **not forced by side-by-side** — unified's `c` record has the same empty
`old-line`. It's a pre-existing v1 property that SxS merely surfaces. If we ever need
per-column-exact numbers, the fix is **carry both numbers on every record** (a v2 change
applying to *both* modes), not a side-only tag. For now nothing in the host needs the
old-file number (the §16 consumers key on `type`/`file`/change-block structure and the
new-file line), so v1 stays as is.

### 17.4 Implication for the lazygit reader (when productionized)

The host's row→identity model must become **row+column→identity** for side-by-side: a
single rendered row can carry two records (the changed-line case), keyed by which column
the OSC precedes. The §16 consumers (nav, restore) were prototyped and tested against
**unified-shaped** output only — one record per row — and have **not** been exercised
against side-by-side delta output. Two things to check when that's wired up: the resolver
(`diffLineInfoFromContents`) must accept two OSCs on one line and bucket them by column,
and the change-block/file scanning (§16.2) must decide whether a "row" is a left-or-right
cell or the pair. Out of scope for the prototype; flagged for productionization.

### 17.5 Implementation gotcha + verification

- **Counter ordering.** Side-by-side advances delta's line-number counters **per aligned
  row** (with increment flags and post-hoc fixups), *not* in the unified
  "all-minus-then-all-plus" order the emitter's counter arithmetic assumes. So the OSC
  strings can't be produced inline at each panel paint — they're **precomputed** over the
  whole minus block, then the whole plus block (the unified order), and looked up by index
  during painting. Context lines also had to be threaded through
  `paint_zero_lines_side_by_side`, which previously never saw the emitter — without that
  the counters would desync across context lines. Records come out byte-identical to the
  single-column emitter.
- **Verified in delta**: payloads correct per column (incl. the paired-row two-record case
  and wrapped continuations carrying none); stripping the OSC456 sequences yields
  **byte-identical** output to unpatched delta across `spaces`/`ansi` fill, narrow-width
  truncation, `--syntax-theme=none`, wrapping, and `--keep-plus-minus-markers`; full suite
  (437 tests) green. **Host (lazygit) consumption of SxS output is a separate, later
  step** — not done here.

---

## 18. Session 8: preserve scroll/selection when switching pagers

A small new consumer of the identity-based restore (diff-line-metadata-notes.md §1
item 7) — the sibling of §16.1's `-U` context-size preserve, triggered by a pager
change instead. The branch already has a **cycle-pagers** feature (`|` / `\`,
`GlobalController.cyclePagers`/`cyclePagersBackward` → `onPagerChanged`); this hangs
the restore off it. One commit: *Preserve the diff scroll position when switching
pagers* (`2c8dfa13b`).

### 18.1 The change

`onPagerChanged` (`global_controller.go`) re-renders the current side's diff into the
**"main" window** view when the focused context is that side panel or the focused
`Normal`/`NormalSecondary` main view. Right before that `HandleRenderToMain()`, it now
calls `Staging.PreserveDiffPositionOnRerender(Contexts().Normal.GetView())` — the exact
same one-liner §16.1's `default` branch uses. No new machinery; the producer
(`PreserveDiffPositionOnRerender` → `restoreDiffLinePositionOnRerender`) already exists
and this is just a third caller ([[isolate-new-concepts-from-clients]]).

### 18.2 Why it's needed in *both* cases (and what each did before)

How a pager is applied decides the async task's command **key** (`cmdStr`, which gates
`onNewKey`'s top-reset), and the two kinds of pager entry differ:

- **Plain `pager:` entry** (e.g. delta ↔ cat). The pager is handed to git via the
  **`GIT_PAGER` env var** (`pty.go`), *not* baked into the command args — so `cmdStr` is
  unchanged across the switch, `ResetOrigin` is false, and the old origin simply
  survived. That looks like "scroll preserved", but it's preserved **by raw line
  number**, which is wrong the moment the two pagers structure the diff differently
  (side-by-side ↔ inline): the same screen line is now a different patch line.
- **`externalDiffCommand` entry.** This *is* baked into the git command
  (`diff.external=…` + `--ext-diff`, `diff.go`), so `cmdStr` **changes**, `ResetOrigin`
  is true, and the re-render **jumped to the top**. (This is the case the user flagged.)

The identity restore fixes both at once: setting `restoreForNextTask` makes
`ResetOrigin = restore == nil && …` false (no top-jump for the externalDiffCommand
case), and `Apply` re-anchors on the same patch line by identity (right anchor for the
plain-pager case when the structure changed). When the structure *doesn't* change it's a
no-op relative to before — the anchor lands at the same screen row.

### 18.3 Scope, fallback, and the one real limitation

- **"main" window only**, exactly like §16.1: `NormalSecondary` is not preserved (still
  jumps). Same accepted limitation.
- **Staging/patch-building untouched.** `onPagerChanged`'s condition is false when those
  are focused (`CurrentSide()` is the side panel, not the staging main context, and the
  key isn't a `NORMAL_*` one), so cycling pagers there doesn't re-render — correct, since
  the staging/patch views render their patch directly, not through a pager.
- **Graceful fallback.** `nearbyDiffLines` only collects candidates that the diff-line
  primitive can resolve in the *old* pager's output, and `FindDiffLine` only lands if it
  can resolve them in the *new* one. If a pager's output is unresolvable (e.g. `cat -n`:
  buffer-parse fails its integrity check, no metadata/hyperlinks), the restore either
  installs no candidates or finds no match → harmless no-op, degrading to the pre-change
  behaviour (stay put for a plain pager, reset for externalDiffCommand). No crash.
- **The headline structural win (side-by-side ↔ inline) is wired but not yet fully
  realized.** Matching across that switch needs the host to resolve the *target* pager's
  output. For side-by-side delta that's the **row+column→identity** resolver that
  §17.4 flags as still a separate, un-built step (`diffLineInfoFromContents` must accept
  two OSCs per line). So today the restore works wherever the primitive already resolves
  both renderings (no-pager, `--color`, single-column patched-delta, hyperlinks), and
  no-ops on side-by-side until §17.4 lands. The wiring here is complete and correct; the
  remaining gap is the shared SxS resolver, not this consumer.

### 18.4 Verification

`build` / `format` / `lint` / `unit-test` green. Headless e2e: `diff/cycle_pagers` and
`staging/diff_context_change` both pass (the former exercises the new call with `cat` /
`cat -n` / default pagers and confirms the no-op fallback doesn't break cycling).
**Interactive sign-off still pending** — as established (§13.1) the headless harness uses
the cmd path, defers no `afterLayout`, and blocks `LAZYGIT_SLOW_RENDER`, so the pager
re-render path and any flicker aren't reproducible there. Confirm with `just debug`,
scrolled down, cycling between a side-by-side and an inline pager and between a `pager:`
entry and an `externalDiffCommand` entry.

---

## 19. Session 9: alt/shift-click a diff line to jump to the editor

A small standalone interaction, only loosely related to the rest of the branch:
a modifier-click in the main view that opens the diff line under the cursor in
the editor — without focusing the view, creating a selection, or being blocked
by a popup. It's the gesture replacement for delta's clickable line-number
gutter (`--line-numbers --hyperlinks`), whose three gripes were: a small/fiddly
target, the horizontal space the gutter costs, and that it's only on
added/context lines. **Interactively signed off** in Ghostty, iTerm2, and VS
Code.

Commits (most recent last):

1. *Let a mouse binding opt into firing while a popup panel is focused* — gocui.
2. *Extract editDiffLine from editLine* — prep refactor.
3. *Carry the keyboard modifier on mouse click events* — gocui fix.
4. *Alt- or shift-click a diff line to open it in the editor* — the feature.

### 19.1 Why alt + shift (the terminal wall)

The original instinct was a single modifier-click or right-click; **a per-terminal
probe killed every single-gesture option** (temporary instrumentation that logged
raw tcell button+modifier for each click — removed after). The findings, which are
the durable lesson here:

- The only mouse input reliably delivered to a TUI is a **plain left-click** (and
  the wheel). Right-click is claimed for context menus (iTerm2, VS Code) or simply
  not forwarded; modifiers are stripped (Ghostty drops ctrl), repurposed for text
  selection (shift/alt bypass mouse capture), or promoted to a secondary click
  (macOS ctrl-click). The SGR mouse protocol can't even carry Cmd/Super.
- **No single chord survives all three terminals.** Ghostty forwards alt (keeps
  shift for selection); iTerm2 forwards only shift; VS Code forwards both. So we
  bind **both alt-left and shift-left** to the same handler: whichever a terminal
  delivers fires the edit, the one it keeps for itself never arrives. (Right-click
  *does* reach Ghostty — see the bug below — but is unusable in iTerm2/VS Code, so
  it was dropped in favour of the alt/shift pair.)

### 19.2 The two gocui pieces

- **Popup-bypass for mouse bindings** (`HandleWhenPopupPanelFocused` on
  `ViewMouseBinding`). Clicks on a non-popup view are normally swallowed by the
  `ShouldHandleMouseEvent` gate; hyperlink clicks already dodge it by being handled
  in an earlier phase. The flag generalizes that: flagged bindings are dispatched
  before the gate, so the edit-click stays live behind the commit-message panel
  (the case the gutter hyperlinks handled and the `e` keybinding can't). One flag on
  the producer; existing bindings stay oblivious ([[isolate-new-concepts-from-clients]]).
- **Modifier-on-press fix.** A click is reported on the button *press*, but the
  driver only applied the keyboard modifier on *release* — which it turns into a
  discarded mouse-move. So modified clicks reached handlers stripped to a plain
  click; the alt/shift binding couldn't have matched without this. **Behavior
  change, global:** every modified click in every view now carries its modifier, so
  an unbound modified click is a no-op rather than silently acting as a plain click.
  Nothing in the tree bound a click modifier before, so nothing relied on the old
  behavior. (A *separate* latent bug surfaced during the probe and was **not** kept:
  the driver converts right/middle-button presses into mouse-move events, so they're
  never dispatched — fixable in ~3 lines if right/middle clicks are ever wanted.)

### 19.3 It's another diff-line-metadata consumer

The handler resolves the clicked whole line via `StagingHelper.GetDiffLineInfo`
(the same primitive the focused-main-view click/`e`/escape-restore use), then
`Files.EditFileAtLine`. So it inherits the backends' fidelity exactly: full
whole-line + deletion support under the OSC-metadata and no-pager/buffer-parse
backends; the plain `--hyperlinks` backend only resolves lines that *carry* a
hyperlink (added/context). So even before delta's gutter is turned off this beats
the gutter click for those lines (whole line, behind popups); turning off
`--line-numbers --hyperlinks` and relying on the metadata backend is what unlocks
deletions and reclaims the gutter.

### 19.4 Productionization placement (planning hint)

Land this as a **separate PR at the very end** of the productionization stack
(§15's decomposition), after the focused-main-view feature and the diff-line
primitive are in — it's an *additional consumer*, not a dependency of anything
else, so it shouldn't gate the rest. Two of its commits are independently
master-worthy and could even go ahead on their own, decoupled from the
metadata work: the `HandleWhenPopupPanelFocused` capability and the
modifier-on-press fix (the latter carries the global behavior-change note above,
so flag it in its own PR description). Still pending for that PR: a confirming
slow-render/interactive pass isn't needed (no async render path here), but it
wants integration-test coverage like the rest of the focused-main-view
interactions (§7 — none exist yet).

---

## 20. Session 10: the diff-line scan was O(n²) — batched to O(n)

Also in this session: the space-selection and `-U`-preserve anchor both moved to
the **middle visible line** (`gocui.View.MiddleVisibleLineIdx`: middle of the
viewport when the content overflows it, middle of the content when it doesn't).
Small and signed off interactively ("feels much better"). The substantial work
was a performance fix.

### 20.1 The problem (user-reported)

Changing the `-U` context size (`{`/`}`) or switching pagers (`|`/`\`) while
scrolled down in a long diff was **extremely slow**: a 1700-line diff took ~2s
with no pager / ~3s with delta; a 9600-line diff took ~33s / ~70s. Both
consumers preserve the scroll position by scanning the re-rendered diff for the
line to land on, so the cost was in the scan. The 5.6× line-count increase
(1700→9600) gave a ~16× time increase — quadratic.

### 20.2 Root cause — three stacked O(n²) passes

The position restore (§16.1) and the file/hunk navigation (§16.2) both resolve
*every* rendered row's identity, and did so **one line at a time**:

1. **Per-line whole-section re-parse.** The buffer-parse backend
   (`parseDiffLineFromBuffer`) parses the target line's *entire* file section with
   `patch.Parse` on every call — O(n) per line for a single-file diff. It was
   called once per line in the incremental load scan *and* again in the EOF
   `Apply` scan (`FindDiffLine` from 0). Two O(n²) passes.
2. **Per-line full snapshot rebuild.** The incremental scan called
   `OffscreenDiffLineContents()` (rebuilds a snapshot of *all* loaded rows) after
   every loaded line — O(n²), and it hit *every* backend, including patched delta
   when scrolled to the end (so it wasn't just a no-pager problem).
3. (Wasted work for delta: the buffer-parse backend was attempted per line even
   when metadata/hyperlinks would answer.)

Confirmed with a throwaway benchmark: a single whole-buffer scan of a synthetic
single-file diff took 115ms at 1700 lines, 3.9s at 9600 — matching the
user's wall-clock once the two-or-three stacked passes are accounted for.

### 20.3 The fix

A single **batch resolver**, `StagingHelper.resolveDiffLines`, parses each file
section exactly once (`parseAllDiffLinesFromBuffer` →
`parseFileSection`/`fileSectionBounds`, extracted from the single-line parser so
both share it — `parseDiffLineFromBuffer` now delegates to it) and applies the
metadata/buffer/hyperlink precedence on top. All whole-buffer scans route through
it: `nearbyDiffLines`, the EOF `Apply` (now `findResolvedDiffLine`, a pure scan
over the resolved table), and **both** navigation scans (`AdjacentChangeBlock`,
`AdjacentFile`) — which had the identical bug. O(n).

The incremental load scan was made O(n) too: it now resolves only the rows that
arrived since it last looked (`gocui.View.OffscreenDiffLineContentsFrom(from)`
instead of re-snapshotting the whole off-screen buffer) and resolves each on its
own via `diffLineInfoPerRow` (metadata/hyperlink only — the per-row backends).
The buffer-parse backend can't resolve a partial diff anyway (the trailing hunk
isn't well-formed mid-load), so it's no longer attempted during the load; the
no-pager case still settles at the EOF swap (`firstPaint` → `Apply`), exactly as
before — only the wasted O(n²) work is gone, the swap timing is unchanged.

The single-line forward resolver (clicks/`e`/`enter`/PR, via `GetDiffLineInfo`)
is behaviorally unchanged.

### 20.4 Verification

- A throwaway test proved the batch parser is **byte-identical** to the per-line
  parser across a multi-file diff (incl. lines outside any section). A throwaway
  perf test showed the whole-buffer scan drop from **1.7s → 1ms** at 1700 lines
  and **~107s → 15ms** at 6800, linear in the line count. Both removed afterward.
- `build` / `format` / `lint` (0 issues) / full `-short` unit suite /
  `diff`+`staging`+`cycle_pagers`+`diff_context_change` e2e all green.
  (`TestNewCmdTaskInstantStop` is the documented pre-existing flake — passes in
  isolation; tasks.go was not touched.)
- Interactive: user confirmed "not entirely instant, but definitely good enough."

This **retires the scan-cost productionization note** flagged in §14.5 and §16.4
(resolved in the prototype, per [[resolve-hard-unknowns-in-prototype]]).

Commits (most recent last): *Add a batch buffer-diff parser that parses each file
section once* (prep), *Add gocui accessor for the newly-loaded rows of an
off-screen render* (prep), *Resolve diff lines in one batch pass instead of once
per line* (the fix). (Plus the two middle-line commits.)

### 20.5 Follow-up surfaced by the speedup: a context-change flicker (no pager) — FIXED

Changing `-U` context with **no pager** while scrolled down in the long diff
showed a **pronounced flicker — a brief frame at a different scroll position,
random run-to-run**; **not** reproducible with delta. The speedup didn't cause
it — it *unmasked* it: the old O(n²) scan throttled the read loop so the load
took ~33s, during which the old content sat stably on screen and hid the
transient; once the load was instant, the transient stood out.

**Pinned with the §10 `SetOriginY` tracer, at normal speed** (slow render would
have hidden it — it re-stretches the load to a stable old-content frame). The
log showed, every keypress, `SwapIn` then `SetOriginY` **~35ms later**:

```
.647  SwapIn             new content displayed, oy still 10354 (the old scroll)
.682  SetOriginY ->9865  Apply finally sets the preserved origin
```

**Cause:** `firstPaint` swapped the off-screen content in *then* ran
`RenderRestore.Apply`, which for the buffer-parse backend re-scans the whole
diff (the O(n) `resolveDiffLines`, tens of ms on a big diff) to find the line to
land on. For that window the new content was displayed at the **stale**
(now out-of-range) scroll; a layout draw landing there drew the bad frame. Delta
didn't flicker because its metadata/hyperlink backend resolves the target
*during* the load (`primaryBufferLine` set), so its `Apply` is instant — no scan
after the swap. (My initial guesses — the loading-ticker `Reset`, or a render
overlap — were both wrong; the tracer is what kept this from being a third
§11-style misdiagnosis. The single-render path *was* clean; the bug was the scan
on the wrong side of the swap.)

**Fix (commit *Scan for the restore target before revealing the new content*):**
`RenderRestore.Apply` now takes a `swapIn func()` and owns the swap — it resolves
the target against the still-**off-screen** (and, at EOF, complete) buffer
*first*, then calls `swapIn`, then settles the scroll. The scan runs while the
previous content is still displayed, so the new content is revealed already at
the right position — the same flicker-free shape the early-resolving backends
always had.

#### Productionization notes (keep in mind)

- **The resolve-then-swap ordering is a real invariant of the restore mechanism,
  not an incidental tweak.** When the gocui async-render improvements are
  productionized (§15 decomposition (a)), `RenderRestore.Apply` must keep doing
  its target resolution against the off-screen buffer *before* swapping in.
  Reordering back to swap-then-resolve reintroduces the flicker for any backend
  whose resolution isn't instant (i.e. buffer-parse / no-pager). The read-loop
  unit test `TestNewCmdTaskRestore` guards this (asserts the render is *not*
  swapped in before `Apply` runs).
- **A small post-swap window is irreducible in this design.** Mapping a matched
  buffer line to a view line (`ViewLineForBufferLine`) needs the wrapped view
  lines, which only exist after the swap, so the `swapIn → SetOrigin` step is
  necessarily post-swap (~1ms). Delta has the same window and it's imperceptible;
  not worth chasing to zero (would need swap+map+SetOrigin under one lock).
- **General lesson for the async-render PR: making renders fast unmasks latent
  ordering/timing transients that slowness was hiding.** Re-test the restore
  consumers (escape, `-U`, pager-switch) at **normal speed** after any perf work;
  `LAZYGIT_SLOW_RENDER` *hides* this whole class by stretching the load back out.
- The batch `resolveDiffLines` is O(n) and "good enough" (user-confirmed), but the
  §14.5/§16.4 *pinned-backend / incremental* optimization is still available if a
  very large changed diff ever lags: today it parses every file section even when
  metadata would answer, and re-snapshots per restore.

---

## 21. Session 11: merge staging + patch-building INTO the main view (the big reframe)

A design discussion (no code yet). The realization that drives it: with the
diff-line-metadata primitive (diff-line-metadata-notes.md), **the pager can now be
used in the staging panel**, which was the long-standing blocker ("the staging
panel needs to parse the diff; a pager breaks that, so I always said no"). And once
that falls, a bigger one follows: if the staging view and the focused main view
render identically and both carry a selection, **why have two views at all?** Fold
staging *and* custom-patch-building directly into the main view — `space` stages the
selected line/hunk/range in place, even in a multi-file diff. The separate
`Staging`/`StagingSecondary`/`CustomPatchBuilder` views go away; `commitFiles` (and
`files`) stay as *browsers*, but `enter` on a file just focuses the main view at
that file's diff. This **changes the scope of the whole project**: the metadata
primitive stops being a navigation nicety and becomes **load-bearing for core
staging**.

Decision (user, session 11): **prototype the whole package** rather than
productionize the pre-session-11 state and defer the merge to a later release. The
end state is far more coherent, and §15's stated reason for staying in prototype
(make it a compelling OSC pitch) is *stronger* here — "stage against any conforming
pager" beats "navigate a restructured diff" as a pitch.

### 21.1 The core insight: we never needed to parse the *rendered* diff

The old "can't use a pager in staging" objection was a *parsing* objection, but the
staging path never actually parsed the rendered bytes. `applySelection`
(`staging_controller.go:237-281`) does:

```
state.SelectedPatchRange() → patch.Parse(state.GetDiff()).Transform({IncludedLineIndices}).FormatPlain() → ApplyPatch
```

`state.GetDiff()` is the **raw `git diff`**, which lazygit already holds in full —
it's the very input it feeds to the pager. The rendered view was only ever used to
**map the user's on-screen selection to patch-line indices**. That mapping is
*exactly* the diff-line-metadata primitive (`GetDiffLineInfo`,
`PatchLineForLineNumber`/`PatchLineForOldLineNumber`, the row→identity resolvers).
So the new chain is:

```
rendered row → (file, type, new/old-line)  [metadata]
            → patch-line index in the raw diff  [PatchLineFor*LineNumber]
            → existing Transform / ApplyPatch
```

Every link except the selection model already exists and is tested. The raw diff is
kept regardless of how it's rendered, so the patch arithmetic is unchanged.

### 21.2 Two risk tiers — don't lump staging and patch-building

- **Staging merge (tractable).** The `Staging`/`StagingSecondary` contexts pass
  `getIncludedLineIndices: nil` (`context/setup.go:44-57`) — staging needs **no
  inclusion highlighting**. It needs only: range + hunk selection on rendered rows,
  `space`-to-stage, and the row→patch-line mapping above. Hunk boundaries are
  already derived from metadata (`AdjacentChangeBlock`, §16.2); the unstaged/staged
  split + `<tab>` already exists as `Normal`/`NormalSecondary` (§14.3). Bounded,
  mostly-built primitives.
- **Patch-building merge (one genuinely new capability).** The `CustomPatchBuilder`
  view renders *itself* via `state.RenderForLineIndices(included)` → `FormatView`
  adds a green background to included lines (`patch/format.go:122-145`). **A pager
  precludes this** — the pager owns the bytes and knows nothing about inclusion
  state. So patch-building needs a **metadata-driven inclusion overlay** (§21.5),
  the single biggest *new* piece. Isolate it.

### 21.3 Resolved technical questions (grounded in the code)

- **Transform copes with non-contiguous selections — confirmed by reading
  `transform.go`.** `Transform` tests `IncludedLineIndices` membership **per body
  line** (`lo.Contains`, `transform.go:150`); there is no contiguity requirement.
  Today's contiguity is purely the *caller* passing `ExpandRange(first, last)`
  (`staging_controller.go:252`). The `pendingContext` /
  `didSeeUnselectedNewFileLine` machinery (`transform.go:137-200`) exists precisely
  to handle a *partial* subset of a change block. Tracing the SxS case (select one
  side-by-side row = `{−deleted_1, +added_1}`, skipping `−deleted_2`/`+added_2`)
  yields a valid hunk: selected −/+ kept, the unselected deletion buffered as
  context after the addition. So the user's "collect all records on the row, union,
  feed to Transform" approach works **as-is**.
- **SxS needs only an "all records on this row" accessor — NOT §17.4's row+column
  bucketing resolver.** Staging collects *every* metadata record on a selected
  rendered row and includes them all (left-column `d` + right-column `a` → stage
  both). You can't stage just one side in a side-by-side rendering; **accepted
  restriction** (uncommon; switch to a single-column pager or none to do it). This
  *downgrades* the §17.4 resolver: the heavy bucketing is only needed for the
  act-on-one-side gesture we're forgoing, so it gates **neither** staging nor
  patch-building. The light per-row collector covers both.
- **Only change-line (+/−) indices are needed.** Context lines are emitted
  regardless of selection (`transform.go:152-157` doesn't check membership), so the
  new model can ignore context entirely when building the included set.
- **Async post-stage selection: the *decision* stays synchronous.** "What's the next
  stageable hunk" is a patch-space computation over the new `git diff`, fetched
  synchronously on refresh (as `RefreshStagingPanel` does today). Only *revealing*
  it — scrolling the pager-rendered view there — is async, and that's the
  restore-by-identity consumer (#6) already built: compute the target identity from
  the new raw diff, install it as the pending restore, let the async render land on
  it. No new async machinery; it rides `restoreDiffLinePositionOnRerender`.
- **Multi-file staging is structural.** Metadata carries `file` per row, the patch
  builder is file-keyed, and #1 already does multi-file splitting. A multi-file range
  just loops per-file `Transform`/`ApplyPatch`. New loop, bounded.

### 21.4 The metadata-reach boundary — the one real product unknown

This is the part with no clean answer, and it's qualitatively different from
everything before: until now every metadata feature was **optional** (old pager →
lose the nicety, still use lazygit normally). Removing the staging panel makes a
conforming rendering **required to stage at all**. The boundary, sharpened:

- **First-class (no metadata pager needed):** no pager, `git diff --color`,
  `delta --color-only` (no line numbers) — all served by buffer-parse (#1).
- **Served via metadata:** any restructuring pager that emits the OSC (patched
  delta/difftastic).
- **The gap:** a **restructuring pager that does *not* emit metadata** (stock
  delta-default, diff-so-fancy, plain difftastic). Only this bucket loses staging.

Fallback options for that bucket, none chosen (deferred past the prototype — depends
on protocol adoption, and the prototype proves the concept with a conforming pager
regardless):

- *Keep the staging panel as a fallback* — user's verdict: **bad**, resurrects the
  whole second code path and forfeits the unification win.
- *Render diffs without the pager when staging is possible* — note this is **not** a
  dynamic switch (staging is possible whenever the main view is focused, so it would
  fire essentially always / on focus, which is strange UX). Realistically it's a
  *static* "this pager can't emit metadata → lazygit never uses it for diff views,"
  decided once — i.e. the user's restructuring pager is simply unused inside
  lazygit. A real loss, recorded honestly as an open question, not a clean fallback.
- *Require a conforming rendering, hard* — simplest; excludes that bucket from
  staging.

**This is the decision that most needs the user's judgement, and it can wait.**

### 21.5 The inclusion overlay (patch-building only) — reserved left column

The new capability patch-building needs. Instead of `FormatView`'s first-char green
background (impossible over pager output), lazygit paints **its own inclusion marker
in a reserved left column**, driven by the per-row metadata identity checked against
the included set. This stays **layout-agnostic** (works over any pager) and the user
considers it a UI *improvement* over today's green overlay.

- **The one-cell shift is intrinsic, not avoidable.** Staging the first hunk into a
  custom patch *is* what creates the patch (and the column); unstaging the last *is*
  what removes it. There is no persistent "patch-building context" to pre-reserve
  during — the operations enter/exit it implicitly. Reserving the column
  unconditionally in any view where a patch is *possible* would avoid the shift but
  is a worse trade. Verdict: shift is inherent, probably fine.
- **Highest-uncertainty NEW piece**, and independent of the staging selection
  mechanics (needs only the metadata primitive + a gocui rendering change) → good
  early de-risk spike (§21.7 step 2).

### 21.6 Decided UX change: always show the selection, anchored at the first hunk

Today the focused-main selection is **on-demand** (`space` toggles it, §4) and
anchored at the **middle visible line** (§20). New model: **the selection is always
shown**, and on focus it **starts at the first hunk** (exactly as entering the
staging view does). This frees `space` for staging with no conflict, and is itself a
**unification win** — the main view now behaves like the staging view (which always
has a selection), the direction we're merging toward.

- **Reverses two earlier decisions:** §4 ("`0` focuses with no selection / scroll
  mode; `space` toggles") and §20 (middle-line anchor). Scrolling still works with a
  selection shown (as in the staging view today); whether a no-selection scroll mode
  survives at all is a step-0 detail.
- Applies only to **diff** main views (files / commits / commit-files / stash),
  **not** to non-diff main content (branch log, reflog, tags, …) where there is
  nothing to act on — there the selection stays off (refined session 11; the
  predicate is "is the main view showing a diff", e.g. via the side panel's
  `DiffableContext`-ness — confirm the exact seam during step 1). Within diff views
  `space` only *stages* where staging applies (files / commits-patch); `e`/`G`/`enter`,
  currently gated on "selection showing," become always available there.

### 21.7 The prototype sequence (linear — one session each, no parallel sessions)

Dependency-first; the genuinely-unproven *capability* (the overlay) pulled early as
a sequential de-risk spike. "Replaces the staging panel" milestone at step 5.

- **Step 0 — Pin the interaction model + the selection-model fork (cheap, on
  paper).** Keys (`space` = stage/unstage, range-select, hunk-toggle, `<tab>`
  unstaged↔staged); fold in §21.6 (always-show selection at first hunk). The one real
  fork: **adapt `patch_exploring.State` to be backend-agnostic vs. build a fresh
  selection model on the main view keyed to rendered rows + metadata.** Lean:
  build-fresh, reusing State's hunk arithmetic — but decide deliberately, it shapes
  everything after.
- **Step 1 — Single-line stage/unstage from the focused main view (working tree).**
  Smallest end-to-end proof of the thesis: rendered row → `(file, patch-line)` →
  `Transform` → `ApplyPatch`, no new selection model (single-line resolution already
  works). Retires the central unknown for almost no code. **Start here.**
- **Step 2 — Inclusion overlay in a reserved column (de-risk spike).** The one new
  *capability*; independent of staging mechanics. Prove a metadata-driven left-column
  marker over pager output. Done early because it's the biggest unknown — if it
  fights us, we've spent little and learned the crux. (This is the "parallel" item
  from the discussion, run linearly.)
- **Step 3 — Range + hunk selection, single-column.** Build the step-0 selection
  model; map a range/hunk → change-line index set → `Transform`. Validate on
  single-column renderings (no pager, unified metadata-delta), one record per row.
- **Step 4 — Multi-record-per-row (SxS) + multi-file.** Add the "all records on the
  row" accessor (light SxS path, §21.3) and the per-file apply loop. Staging now
  works over any conforming rendering.
- **Step 5 — Staged/unstaged split + post-stage reveal.** Stage from `Normal` /
  unstage from `NormalSecondary`, `<tab>`, synchronous next-hunk decision +
  restore-by-identity reveal. **Milestone: the working-tree staging panel is
  functionally replaced.** The escape-from-staging routing (§12.2 + the snapshot)
  starts dissolving here.
- **Step 6 — Patch-building from the main view.** Reframe `PatchBuilder` interaction
  in identity terms; add/remove hunks from the commits/commitFiles main view, with
  step 2's overlay showing membership.
- **Step 7 — `enter` on a file focuses the main view, for *both* `commitFiles` and
  `files`.** *(SUPERSEDED — DROPPED; see §21.24.)* The browsers stay; `enter` focuses
  the main view at that file's diff instead of opening a separate explorer.
- **Step 8 — Tear out the separate explorer views + escape/snapshot machinery.**
  *(SUPERSEDED — DROPPED; see §21.24. The explorer views are kept on purpose as a
  testing reference.)*

**Decision gate after step 2:** if the overlay is clean and step 1 confirms the core
chain, the rest is bounded execution.

### 21.8 Decisions locked / open / memory

- **Locked:** prototype the whole merge (don't stop-and-productionize first); SxS
  staging includes all records on the row (no single-side gesture); inclusion overlay
  = reserved left column (shift is intrinsic); selection always shown, anchored at the
  first hunk (reverses §4 + §20); step 7 covers `files` *and* `commitFiles`; linear
  sessions.
- **Open (the real product unknown):** the §21.4 metadata-reach fallback for a
  restructuring-pager-without-metadata user. Deferred past the prototype; depends on
  OSC adoption.
- **Carried forward:** the identity-restore machinery (`restoreDiffLinePositionOnRerender`)
  stays — staging mutates the diff, so "stay put after staging" is the same consumer,
  now *un*-entangled from a context switch. Concurrency stays mutex-based for now
  ([[main-thread-over-mutexes-direction]]). Keep new concepts off existing clients
  ([[isolate-new-concepts-from-clients]]). Resolve the overlay unknown in the
  prototype ([[resolve-hard-unknowns-in-prototype]]).

### 21.9 Session 12 (2026-06-18): Step 1 implemented + findings

Steps 1+2-worth of behavior built (uncommitted, all green): **always-show selection
in diff views, anchored at the first visible change line**, and **single-line
directional staging on `space`**. Two corrections applied from interactive feedback:
the predicate became the `types.DiffMainViewContext` marker interface (tagged on
Files/LocalCommits/SubCommits/Stash/CommitFiles/ReflogCommits — *not* the
`GetOnClickFocusedMainView`-sniffing first cut, which conflated "stageable" with
"shows a diff"; reflog shows a diff via `git show` but isn't stageable), and the `0`
anchor became "first change line at or below the viewport top" (`FirstChangeLineInView`)
rather than jumping to the diff top.

**Signed off as good-enough:** the `0`-while-scrolled anchor (will change anyway once
hunk selection lands and is on by default); the one-hunk-stage case (only the tab
title + the files-panel ` M`→`M ` status change — identical to pressing `space` on
the file in the files panel, so fine).

**Fixed this session:** `<tab>` between the unstaged/staged panes left *no* selection
in either (leaving pane cleared `Highlight`; landing pane established none). Now
`togglePanel` anchors the landing pane on its first visible change line via a shared
`showInitialDiffSelection` helper (also used by `focusMainView`). Per-pane selection
*memory* (like the staging view keeps) is not done — we re-anchor on each switch;
fine for now.

**The big-picture caveat the user raised:** Step 1 alone is *not* a usable experience
— single-feature increments don't become testable until range/hunk selection (Step 3,
hunk-select on by default) and the selection-lifecycle polish are in. Expect more
rough edges to surface as we go; the ones known so far are below.

**Two findings to carry forward (not fixed now):**

- **The `GetOnStageFocusedMainView` handler-channel doesn't scale (architecture
  debt, revisit for the real implementation).** The staging panel has *many* commands
  beyond `space` (`d` discard a hunk, edit, and more); routing each through its own
  side-panel handler channel (interface method + `BaseContext` field + `attach.go`
  registration + `baseController` default) is too much boilerplate per command.
  Accepted as a prototype expedient (user OK'd continuing for now), but the real
  design wants a better mechanism — e.g. the focused main view *hosting* the
  patch-explorer command set directly, or a shared staging-command controller bound
  to the diff main contexts, rather than N bespoke delegation channels. Flag for
  Step 8 / productionization.
- **Delta's background-conveyed side is invisible under the selection highlight
  (usability, future).** delta indicates added/deleted/context purely by background
  colour, and the selection takes over the background — so you can't tell what kind
  of line is selected. Needs a different way to mark the selected line(s) (e.g. a
  gutter marker, a foreground/border treatment, or reserving a column — cf. the §21.5
  inclusion overlay, same class of problem). Not for this prototype; keep in mind.
  **Resolved §21.34** (`narrowSelectionHighlight`: a narrow selection bar at the left edge only).

### 21.10 Step 1 committed; Step 3 design settled, fold done (resume here)

**Committed (branch, not pushed), most recent last:**

```
f660567cb Extract diffSplitState from the files diff renderer            (prep)
d196d4e78 Stage the selected diff line from the focused main view         (step 1: always-show selection at first visible change + single-line directional staging + the <tab> fix)
d0f72cfb8 Session notes: plan to merge staging into the main view, and step 1
4fab0a7c6 Fold ViewSelectionController into MainViewController            (step 3 prep)
```

So **steps 1 + 2-of-the-plan are done** (always-show selection + single-line
staging), and **step 3's prep is done** (the `ViewSelectionController` fold).
`build`/`lint`/`unit-test`/`format` green throughout; `just generate` was run for
the keybinding cheatsheets. Not interactively re-tested since the `<tab>` fix.

**Step 3 design — SETTLED, ready to build (the feature increment):** range + hunk
selection in the focused main view, built fresh (own type, not shared with
`patch_exploring.State` — the explorer is going away, so the temporary duplication
is fine). Grounded findings:

- **gocui renders a range selection natively** — `View.SetRangeSelectStart(y)` +
  `Highlight` + cursor highlights the span (`view.go:607-615`); `CancelRangeSelect()`
  resets. So LINE = cancel range-select (cursor highlight); RANGE/HUNK = range-select
  from anchor to cursor. **No new highlight machinery.**
- **Mirror the staging mode machine** (`patch_exploring/state.go`): `selectMode`
  (LINE/RANGE/HUNK), sticky vs non-sticky range, `userEnabledHunkMode` (config-default
  hunk vs user-toggled — matters for escape), default HUNK when
  `UseHunkModeInStagingView && !IsSingleHunkForWholeFile`. Reimplement in **view-line
  space** with hunk bounds from the **metadata `isChange` array** (reuse
  `resolveDiffLines`/`changeBlockStart`/`AdjacentChangeBlock`), **not** `patch.Lines()`.
- **State home:** fields on `MainViewController` (two instances — Normal &
  NormalSecondary — each owns its mode). The nav handlers (`handleLineChange` etc.)
  now live there (the fold), so mode-aware ↑/↓ has its home.
- **↑/↓ are already selection-move** (the folded `handleLineChange`: Highlight → move
  selection, else scroll). Extend: **hunk mode → move to prev/next hunk**
  (`AdjacentChangeBlock`); line mode → by line (today); range mode → extend. `a`
  (`Main.ToggleSelectHunk`) toggles to line mode for line-granular movement.
- **Keys to add, mirroring staging:** `Universal.ToggleRangeSelect` (`v`, sticky range),
  `Main.ToggleSelectHunk` (toggle hunk/line), `Universal.RangeSelectUp`/`RangeSelectDown`
  (shift-↑/↓, extend non-sticky range). `<left>`/`<right>` hunk nav + `n`/`N` file nav
  already exist.
- **Config default:** select the first *visible* hunk on focus when hunk mode is the
  default (grow `FirstChangeLineInView`/`showInitialDiffSelection` to expand to the
  block bounds, not just the first change line).
- **Range-aware staging:** `GetOnStageFocusedMainView` currently takes a single
  `viewLineIdx`; make it take the selected row range (or read the view's range-select),
  collect the change-line patch indices across the rows, and apply **one** `Transform`
  with the set (`Transform` already handles arbitrary/non-contiguous index sets — §21.3).
  Side-by-side "all records on a row" is still step 4; step 3 stays single-column.

Suggested commit shape for the feature: it may split into (a) mode state + rendering +
mode-aware nav + keys (the selection model) and (b) range-aware staging — or land as one
feature commit; decide when the diff takes shape.

### 21.11 Session 13 (2026-06-18): Step 3 implemented + committed (range + hunk selection)

**Committed (branch, not pushed), most recent last:**

```
c4827a269 Select and stage ranges and hunks from the focused main view   (the feature: selection model + range staging)
a9eee242a Select the first hunk when focusing the main view in hunk mode  (the config-default-on-focus)
```

Landed as **one feature commit + one default-on-focus commit**, not the (a)/(b) split:
the selection model and range staging are coupled (a hunk you can select but `space`
stages only one line is an incoherent intermediate — the no-regression rule), so they
go together. The hunk-default-on-focus is genuinely separable and came second.

**How the settled design actually landed:**

- **State home = MainContext, NOT MainViewController** (deviates from §21.10; decided
  with the user — [[diff-selection-state-home]]). `MainContext.DiffSelectState()` holds
  `{Mode (Line/Range/Hunk), RangeIsSticky, UserEnabledHunkMode}`. Reason: three sites
  need a pane's mode (the controller, `focusMainView`, and `togglePanel` setting the
  *other* pane) and only the context is reachable from all three. Mirrors how
  `patch_exploring.State` lives on the patch-explorer context.
- **No new highlight machinery, confirmed.** The selected line + range anchor stay in
  the gocui view (cursor + `rangeSelectStartY`); the controller only flips the mode and
  re-derives the view's range. LINE = `CancelRangeSelect`; RANGE = anchor fixed, cursor
  is the moving end; HUNK = cursor on the block's first line, anchor on its last (so the
  cursor sits at one end — gocui can't highlight a block with the cursor in the middle;
  invisible anyway since the whole block is highlighted). Block bounds come from the
  metadata `isChange` run via a new `StagingHelper.ChangeBlockBounds` (view-line space,
  reuses `resolveDiffLines`), never `patch.Lines()`.
- **Mode-aware nav:** ↑/↓ — hunk mode steps hunk-to-hunk (`AdjacentChangeBlock`), a
  non-sticky range collapses on a plain move, a sticky range extends. ←/→ (hunk nav) and
  n/N (file nav) already existed and now re-expand in hunk mode. Pages/top/bottom drop
  hunk mode + collapse non-sticky range, like the staging view's `AdjustSelectedLineIdx`.
- **Keys added:** `v` (`ToggleRangeSelect`, sticky), `a` (`ToggleSelectHunk`),
  shift-↑/↓ (`RangeSelectUp/Down`, non-sticky extend). Clicks reset to a single-line
  LINE selection.
- **Range-aware staging:** `GetOnStageFocusedMainView` signature changed from one
  `viewLineIdx` to `(firstLineIdx, lastLineIdx)` (the view's `SelectedLineRange`).
  `StagingHelper.ChangeLinesInViewRange` collects the change lines across the rows;
  `FilesController.stageDiffLines` applies **one** `Transform`.

**The one real surprise — patch-line resolution had to change (`PatchLineFor*LineNumber`
is the wrong primitive for ranges).** Resolving each selected change line's patch index
with `PatchLineForLineNumber(newLine)` / `PatchLineForOldLineNumber(oldLine)` is quirky at
hunk boundaries: it returns the @@ header for a change on the *first* line of a hunk (no
leading context — e.g. a change on line 1 of a file), and for a *modified* line (a `-`/`+`
pair) the new-line lookup lands on the deletion, not the addition. A range routinely spans
both. Fix: `stageDiffLines` keys each selected change line by `(file line number,
deletion?)` and **scans the parsed patch**, matching body lines by identity via the
quirk-free inverse maps `LineNumberOfLine` / `OldLineNumberOfLine`. Side effect: this also
**fixes staging a change on the first line of a file** (a latent step-1 bug in the
single-line path, which used the same primitive). Folded into the feature commit because
robust resolution is *required* for ranges, not optional polish.

**Tests (all green; full `e2e-all` green):** `stage_range_from_main_view` (range over a
deletion + its replacement + context, exercising the resolution fix),
`stage_hunk_from_main_view` (`a` toggle from line mode, hunk-off config),
`select_hunk_on_focusing_main_view` (hunk-default-on-focus selects the first block).

**Deferred / open (carry forward):**

- **Post-stage selection isn't re-anchored** — after `space` the diff re-renders and the
  cursor/range are left where they were (stale). This is the step-5 "post-stage reveal"
  (restore-by-identity, consumer #6); not done here, matches step 1's level. Known rough
  edge; the tests stage from a fresh focus to stay deterministic.
- **`IsSingleHunkForWholeFile` refinement skipped.** Hunk-default-on-focus keys purely off
  `UseHunkModeInStagingView`, so a whole-file-single-block diff (new/deleted file, no
  context) defaults to hunk = select-everything rather than dropping to line mode like the
  staging view. Computing it from metadata at focus (diff maybe multi-file / still
  streaming) is awkward; deferred.
- **SxS multi-record-per-row + multi-file staging = step 4** (unchanged). Step 3 stays
  single-column, single-file (`stageDiffLines` bails on a directory node).
- **Unstage-from-secondary via the main view** (`<tab>` to NormalSecondary, `space`)
  reuses the same path with `reverse=true` but wasn't interactively/e2e verified this
  session.
- **Interactive sign-off pending** (per §21.9, single increments only become testable once
  this is in; the user evaluates the feel — e.g. hunk-default jump-on-focus, delta's
  background-conveyed side under the highlight per §21.9).

**Interactive feedback, session 13 (user; both OK to leave for now):**

- **delta + hunk selection makes the §21.9 colour problem worse.** With a whole hunk
  selected you see one solid block of blue with no boundary between the deleted and added
  lines — delta conveys add/delete purely by background, which the selection overrides.
  Worse than the single-line case. Still deferred (same class as the §21.5 overlay — a
  gutter marker / foreground treatment / reserved column is the eventual answer).
  **Resolved §21.34.**
- **`a` on a context line below the last hunk diverges from the staging panel.** The
  staging panel selects the *last* hunk (searches backward when there's nothing forward);
  ours keeps the single line (`ChangeBlockBounds` only snaps *forward* to the first block
  at/after the anchor, finds none, falls back to a single-line selection). Minor; to match,
  `ChangeBlockBounds` would fall back to the nearest block *above* when none is at/below.

### 21.12 Session 13 (cont.): Step 4 done — SxS multi-record + multi-file staging

**Committed (most recent last):**

```
533046c5e Stage a selection spanning several files from the focused main view  (multi-file)
c63266113 Stage both sides of a side-by-side row from the focused main view     (SxS multi-record)
```

Built **multi-file first** (testable), then **SxS** (interactive-verify only), per the
user's call. Both green; full `e2e-all` green.

- **Multi-file.** The focused main view of a directory node already renders a multi-file
  diff (`WorktreeFileDiffCmdObj` takes a directory). Dropped the `!node.IsFile()` bail;
  the handler now groups the selected change lines by `info.Path`, maps each path back to
  its `models.File` via `FileTreeViewModel.GetFile(ToSlash(rel(WorktreePath, absPath)))`,
  and applies one patch per file. `stageDiffLines` no longer logs/refreshes — the handler
  does both once around the batch. Direction is uniform (the whole diff is one side), so
  it's decided once. Single-file now goes through the same path-lookup (no special-case),
  exercised by the existing step-3 tests. New test: `stage_range_spanning_files_from_main_view`.
- **SxS multi-record-per-row.** `gocui.DiffLineContent.Metadata` keeps only the *first*
  payload per row — drops the right side of a side-by-side row. Added
  `View.DiffLineMetadataPayloads() [][]string` (distinct payloads per buffer line; new
  accessor, left the single-payload one for the nav/restore consumers per
  [[isolate-new-concepts-from-clients]]). `ChangeLinesInViewRange` now resolves *all* a
  row's payloads when it has metadata, else falls back to the single resolved record
  (no-pager / buffer-parse / hyperlink) — so single-column is unchanged. Both sides of a
  SxS row stage together (can't stage one side; accepted §21.3 restriction). **Not
  e2e-testable** (harness has no real delta SxS); covered by a gocui unit test
  `TestDiffLineMetadataPayloads` for the accessor; the full chain needs interactive
  sign-off with the patched delta in SxS mode.

**Interactive sign-off, session 13 (user):**

- **Steps 3–4 staging verified good**, including **SxS staging from BOTH patched delta AND
  difftastic's side-by-side view** — so the per-cell multi-payload path holds across both
  conforming pagers, not delta-specific.
- **Unstaging works but loses the selection / has focus problems** (user's words). Expected
  to be resolved by step 5 (the post-stage reveal + selection lifecycle). Concrete repro
  cases not yet captured — get them at step-5 start to target the exact symptoms rather
  than guessing.

**NEXT: step 5** — staged/unstaged split + post-stage reveal (stage from Normal / unstage
from NormalSecondary, `<tab>`, synchronous next-hunk decision + restore-by-identity reveal).
**Milestone: the working-tree staging panel is functionally replaced.** This is where the
deferred post-stage selection re-anchor (§21.11) AND the just-reported unstage focus/
selection-loss both land.

### 21.13 Session 13 (cont.): Step 5 part 1 — post-stage reveal (selection advances)

**Committed:** `c0c91d10b` Advance the selection to the next change after staging from the
main view.

Fixes the user-reported "selection lands outside the content / looks like there's none"
after stage/unstage. `MainViewController.stageSelectedLine` now installs a restore
(`StagingHelper.RevealSelectionAfterStaging`, sibling to `RestoreFocusedMainViewOnEscape`)
*before* the handler triggers the post-stage re-render. Candidates, priority order:
change block after the selection → before it → the acted-on line itself (only survives the
all-staged flip). Rides the existing `restoreDiffLinePositionOnRerender`. Re-selects in the
current mode (hunk → walks hunk-to-hunk); a range collapses to a line first (staging
consumes it). Works for both stage (from Normal) and unstage (from NormalSecondary) — same
code path, `reverse` differs. Tests: `select_next_hunk_after_staging_from_main_view`,
`select_next_hunk_after_unstaging_from_main_view`. Full `e2e-all` green.

**Part 1 signed off (session 13, user): "selection invisible" symptom gone, works very well.**

**Step 5 part 2 — STILL OPEN, and it UNIFIES with a bug the user found while testing.** Both
are the same issue: *after stage/unstage, focus should follow the side you were acting on,
but doesn't.*

- **User-reported bug:** a file with **only staged changes** → `Normal` shows the *staged*
  diff (`mainShowsStaged`, no split). Unstaging the *first* hunk makes the file
  staged+unstaged → split appears → `Normal` flips to show *unstaged* and the staged content
  jumps to `NormalSecondary`, but focus stays on `Normal`. Wrong: should stay with the staged
  content (i.e. follow it to `NormalSecondary`) until it's empty.
- **The original part-2 case:** unstaging the *last* staged hunk from `NormalSecondary`
  empties the staged side → split collapses → `NormalSecondary` hides, focus stuck on it.

**Why they're one issue:** the staging *view* has fixed panes (Staging=unstaged,
StagingSecondary=staged) so focus-follow is trivial; the merged `Normal` shows *either* side
(per `mainShowsStaged`), so a stage/unstage that changes the split reassigns which pane holds
the side you care about. **Unified rule:** after the op, focus `NormalSecondary` iff
(unstaging AND post-op state is split), else `Normal`. Covers the bug (only-staged → unstage
one → split → NormalSecondary), original part 2 (unstage last → no split → Normal), and all
stage cases (→ Normal). Stage-everything-from-Normal is already fine (Normal flips to staged
via `mainShowsStaged`; part-1 reveal's self-candidate lands there).

**Why it's not a quick patch — the async coordination:** the reveal must ride the post-stage
re-render, so the target pane must be known *before* the re-render fires; but the split state
is only known *after* `Refresh` updates the model. **Open Q to resolve FIRST next session:** is
the post-`Refresh` main-view re-render *queued* (runs after `stageSelectedLine` returns) or
synchronous? If queued, install the reveal + switch focus *after* the handler (model updated →
split known) and it still rides the queued render. If synchronous, need a different hook (e.g.
predict split, or restructure so `stageSelectedLine` owns the refresh). The fix belongs in the
**stage flow**, NOT `FilesController.GetOnRenderToMain` (broad callback, fires unfocused).
Deferred for fresh budget — subtle async work, long session.

**Usability note (user, don't fix now — possible future special-case):** in a multi-file diff
containing a **deleted file**, stepping through hunks and staging them stages the deletion of
the file's *content* but not the *file entry* → you land on `MD` status, not `D`. Same as the
staging panel could produce, but there you had to deliberately enter staging on the deleted
file, so it was unlikely; here it happens easily. Consider special-casing: when a deleted
file's *entire content* is selected/staged, stage the file deletion itself (→ `D`). Belongs
near the per-file apply loop in `stageDiffLines`/the handler.

### 21.14 Session 14 (2026-06-18): Step 5 part 2 done — focus follows the acted-on side

**Milestone reached: the working-tree staging panel is functionally replaced.** Committed
(most recent last):

```
e36b6ceaa Let the post-stage reveal target a different pane than the one acted on   (prep)
2b3ac3017 Follow the acted-on side to the right pane after staging from the main view (behavior + tests)
```

**Open Q resolved: the re-render is QUEUED, the model is updated SYNCHRONOUSLY.** The handler
(`FilesController.GetOnStageFocusedMainView`) calls `Refresh({FILES, STAGING})` in the default
**SYNC** mode, whose files goroutine updates `Model().Files` + `SetTree()` synchronously and is
`wg.Wait()`ed on — so by the time the handler returns, `node.GetHasStagedChanges()` /
`GetHasUnstagedChanges()` (hence `diffSplitState`) reflect the post-op state. But the *main-view
re-render* is deferred: `refreshFilesAndSubmodules` → `OnUIThread(refreshView(Files))` →
`OnUIThread(postRefreshUpdate(Files))` → (we're focused on a `NORMAL`/`NORMAL_SECONDARY` context,
so [`view_helpers.go`] `postRefreshUpdate` line ~146 takes the `NextInStack` branch) →
`Files.HandleRenderToMain()` → `FilesController.GetOnRenderToMain`'s closure renders both panes.
All on later UI-loop iterations, after `stageSelectedLine` returns. So §21.13's "if queued"
branch holds: **decide focus + install the reveal after the handler returns; the reveal rides the
queued render.**

**The unified rule (as predicted in §21.13), implemented:** the handler returns the
focused-main view name that should hold focus — `NormalSecondary` iff (unstaging **and** post-op
split), else `Normal` — and `MainViewController.stageSelectedLine` re-selects in that pane and
`Push`es it. The decision is files-specific (reverse + split), so it lives in the handler; the
GUI mechanics (focus switch, select mode, reveal install) stay in the controller. The
`onStageFocusedMainViewFn` signature gained a `(focusViewName string, err error)` return; `""`
= nothing staged (skip the reveal — also fixes a latent part-1 wart where a no-op stage left a
lingering restore).

**The cross-pane wrinkle (decision I made; wasn't fully spelled out in §21.13):** when focus
crosses to the other pane, the **selection follows too** — a focused pane with no selection is
exactly the "invisible selection" bug part 1 fixed. So the reveal's candidate change lines are
read from the pane we **acted in** (still showing the pre-staging diff until the queued render),
but the restore is installed on, and re-selects in, the **target** pane. Two enabling pieces, in
the prep commit:
- `RevealSelectionAfterStaging(sourceView, targetView, …)` — candidates from `sourceView`, restore
  on `targetView`. Identity matching survives the stage/unstage flip because `SamePatchLine` keys
  on (path, source-line, isDeletion) — **not** the staged/unstaged side — so a staged hunk's
  identity is found again in the unstaged re-render and vice-versa.
- **Get-or-create the target's buffer manager** (`GetOrCreateViewBufferManagerForView` on
  IGuiCommon → `gui.getManager`). The only-staged→split case focuses `NormalSecondary`, which may
  never have rendered, so its manager doesn't exist yet; `SetRestoreForNextTask` needs it to exist
  before the queued render reuses it. Verified safe to pre-create: a fresh manager is
  `!IsLoading()` and its `ReadLines` no-ops (nil `readLines` channel), so the layout clamp
  ([`layout.go`]) behaves identically. The target also inherits the source's `DiffSelectState`
  (line/hunk mode) via a struct copy (`*target = *source`; a no-op self-copy in the same-pane case).

**`Push`ing a MainContext is render-safe:** Normal/NormalSecondary have no `onRenderToMainFn`
(SimpleContext.HandleFocus only fires `onFocusFns` for them — `GetOnFocus` resets
`HighlightInactive`), so the focus switch can't render prematurely or consume the restore. The
queued `postRefreshUpdate` does the render, and `NextInStack(NormalSecondary)` is still Files, so
both panes re-render.

**Tests (both cross-pane; same-pane cases already covered by the part-1 tests):**
`focus_follows_staged_side_to_secondary_after_unstaging` (the user-reported bug: only-staged →
unstage first hunk → split → focus + selection move to `Secondary` on the next staged hunk) and
`focus_returns_to_main_after_unstaging_last_staged_hunk` (original part-2: unstage the only
staged hunk from `Secondary` → split collapses → focus + selection return to `Main` on the
now-unstaged change). Full `e2e-all` green (incl. the previously direnv-flaky worktree test),
`unit-test` + `lint` clean.

**Interactive sign-off pending.** Not yet eyeballed in the real app. Worth a live pass on: the two
cross-pane cases above; the same-pane cases still feel right; and the **multi-file directory
diff** focus-follow (handler reads the *directory* node's aggregate split — should be correct but
isn't e2e-tested). The SxS/delta paths aren't affected by this change (focus logic is
pager-agnostic), but a quick confirm doesn't hurt.

**Still open (unchanged):** the deleted-file usability note above (`MD` vs `D`); §12.2's
custom-patch-builder escape routing is a *patch-building* concern, not part of the merged
working-tree staging milestone.

### 21.15 Session 15 (2026-06-18): rebased onto the wrapped-fill branch; two step-5 reveal bugs fixed

Step-5-part-2 interactive testing surfaced two bugs, both fixed. First, though, the prototype
was **rebased onto master after merging `fix-coloring-of-wrapped-delta-lines`** (the `lineType`
struct, sentinel removal, draw-time `trailingFillAttributes` fill). The user dropped that branch's
"inline finishLine" commit so `finishLine` survives as the home for the metadata fix below
([[clean-history-no-back-and-forth]]).

**Bug 1 — blank changed lines lost their metadata under delta (commit `03b747739`).** With a
**colored** input diff (which lazygit feeds delta via `--color=always`), delta renders some blank
deleted/added lines as just the OSC 1717 metadata + an empty line — either no fill at all, or a
background + `ESC[0K` that the wrapped-fill branch now turns into a draw-time attribute rather than
content cells. Either way the line has **no cells**, so its metadata had nowhere to live and the
line resolved to "not a change" — splitting a hunk at the blank line (`0` selected only the part
above it). Fix: `finishLine` keeps a content-less sentinel cell **only when there's pending OSC
metadata** to carry (not the old unconditional sentinel + `prevFgColor` coupling that `7dc18f3eb`
removed). Reproduced offline only after matching the real pipeline exactly: **colored** git diff,
`TERM=dumb`, truecolor (`COLORTERM` is passed through to the pager — lazygit strips only `TERM*`),
and a `gocui.View` in `OutputTrue` mode. (256-color SGR leaking was a *test artifact* of using
`OutputNormal`; under `OutputTrue` it parses fine, so no real gap. `swallow-unknown-escape-sequences`
is orthogonal — it stops *unknown* sequences leaking under lower modes.)

**Bug 2 — the post-stage reveal jumped to an earlier hunk (commit `9fb3c11cc`).** The reveal
captured the next/prev change block as a patch identity and matched it in the re-render via
`SamePatchLine`, which keys a **deletion** on its **old-file (index) line number**. Staging a hunk
that changes the line count shifts the index-side numbers of every hunk below it, so a deletion-led
"next hunk" candidate no longer matched; the reveal fell back to a worse candidate, often colliding
with a header/context row (`findResolvedDiffLine` didn't require a change line) and, in hunk mode,
snapping to the *first* change block — i.e. an earlier hunk. (Proven by instrumenting the match:
candidate `-e` was `{NewLine:7, OldLine:5}` pre-stage, `{NewLine:7, OldLine:7}` post-stage.) Fix: a
**reveal-specific matcher `matchByWorktreeChange`** — match a change line by its worktree (new-file)
number, which staging never moves, requiring the same side. Safe because the reveal only targets
change blocks and `selectHunkAround` expands the whole block, so the new-file number's
consecutive-deletion ambiguity (the reason `SamePatchLine` uses old-file) doesn't matter. The
matching is now a `diffLineMatch` parameter of `restoreDiffLinePositionOnRerender`; the **escape
restore and `-U` preserve keep `matchByPatchLine`** (same staged state → no shift; `-U` deliberately
anchors on context lines). Regression test:
`advance_to_next_hunk_after_staging_shifts_line_numbers` (lands on the earlier hunk without the fix).
The escape restore has the same staging-shift exposure but is **left as-is** — it dissolves with the
staging panel when this feature lands.

**Also:** a latent bug the prototype's off-screen render exposed in the rebased wrapped-fill code —
`parseInput` wrote `trailingFillAttributes` to `v.buf` instead of `b` (wrong buffer during an
off-screen render; could panic). User fixed it (`ba3389562`, fixup).

**Status:** all three repos' tests green (`unit-test`, `e2e-all`, `lint`). Step-5 staging now holds
up under delta on diffs with blank changed lines and line-count-changing hunks. User to re-test for
any remaining reveal misfires (their original report predates both fixes).

### 21.16 Session 15 (cont.): line-mode reveal — advance to the next line, not the next block

Hunk-mode reveal signed off as perfect after §21.15. **Line mode** had a remaining bug (user:
"kind of important"): staging the first line of a multi-line block jumped to the *next block*
instead of the next (now first) line of the *same* block. Cause: the reveal's primary candidate
was `AdjacentChangeBlock` (next change *block*), which skips the rest of the acted-on line's block —
right for hunk mode (whole block staged), wrong for line mode (one line staged, rest remains). Fix
(`cb784bde3`): new `AdjacentChangeLine` (first change line strictly after/before the selection) for
the next/prev candidates. Unifies both modes — after a whole block (last line followed by context)
the next change line *is* the next block's first line, so hunk mode is unchanged; after a single
line it's the next line of the same block. `place` stays mode-aware (hunk-expand vs single line).
Regression test `select_next_line_after_staging_line_from_main_view` (note: line mode needs
`UseHunkModeInStagingView: false` — hunk mode is the default). All green.

### 21.17 Session 16 (2026-06-19): reveal reworked to preserve the change-line ordinal

Two more reveal bugs (user, line mode): staging a deletion in the *middle* of a deletion block
jumped to the block's first line; and unstaging a modification's deletion jumped back to an
earlier block. Root cause (confirmed by instrumenting the matcher): **no line-number identity is
stable AND unique across a stage.** Deletions in a block all share the same new-file number (so
matchByWorktreeChange's findResolvedDiffLine returned the first); and in the **staged** pane the
new side is the index, which staging/unstaging shifts, so an addition candidate's new-file number
moves and the match misses → falls back to a prev-block/header. The §21.15 worktree-line match and
the §21.16 adjacent-change-*line* each fixed one case and left another.

**Fix (`cd27c07e2`): preserve the selection's ordinal among change lines** — exactly what the
staging view does with its patch-line index. Read the acted-on line's change-line ordinal from the
source pane before the op; after the re-render, select the change line at that ordinal in the
target pane (clamped to the last). The op removes the acted-on change line(s), so that ordinal then
holds the next surviving change (next line of the same block, or next block when a whole block was
staged). No reasoning about which side's numbers are stable → deletions and the staged pane are
handled uniformly. `matchByWorktreeChange` and `AdjacentChangeLine` are gone; the identity restore
(escape, -U) and the new positional restore now share an extracted `installDiffLineRestore`.
Regression tests `select_next_deletion_after_staging_within_a_block`,
`select_next_change_after_unstaging_a_deletion`. All three repos green.

**Open nuance — the collapse case.** When the acted-on side *empties* (e.g. unstage the last staged
hunk → focus moves to the unstaged pane), the source and target are now *different* diffs, so the
preserved ordinal lands on whatever change sits at that ordinal in the other side rather than
specifically on the just-acted line. `focus_returns_to_main_after_unstaging_last_staged_hunk` still
passes (there the just-unstaged change *is* the first one), but it'd differ if a pre-existing
unstaged change sat above it. Left as-is pending the user's call; restoring the old self-identity
behavior there would need a "the acted-on side emptied" signal threaded from the stage handler.

**History note:** `cd27c07e2` supersedes the *mechanisms* of `cb784bde3` (§21.16) and `9fb3c11cc`
(§21.15) — both their matchers are deleted here, though their tests survive. The three reveal
commits would collapse cleanly into one if the user wants to rewrite (their call).

### 21.18 Session 16 (cont.): step 5 signed off — MILESTONE reached; step 6 next (new session)

**Step 5 fully signed off (user): "works perfectly."** The working-tree staging panel is now
functionally replaced by the focused main view — single-line / range / hunk / multi-file / SxS
staging and unstaging, the staged/unstaged split with focus following the acted-on side, and the
post-stage reveal (change-line-ordinal model, §21.17) all working under delta, including the
collapse case (acted-on side empties → jump to the first change of the other pane; confirmed fine,
matches the old staging panel — no change wanted).

**NEXT: step 6 — patch-building from the main view (a new session).** This is the higher-risk tier
(§21.2): commits / commit-files panels build a *custom patch* rather than staging, which needs the
metadata-driven **inclusion overlay** (§21.5) — the reserved left column showing which lines are in
the patch — the one genuinely new capability and the highest-uncertainty piece. Start from §21.5 +
the §21.7 plan. After step 6: step 7 (enter focuses the main view for files AND commitFiles), step 8
(tear out the explorer views + escape machinery).

### 21.19 Session 17 (2026-06-19): step 6 design + plan (the inclusion gutter)

The §21.7 **step-2 de-risk spike (the inclusion overlay) was skipped** — the implementation
sessions ran 1 → 3 → 4 → 5. So step 6 is two pieces: (1) the inclusion overlay (§21.5) — the
genuinely-new capability, never built — and (2) the patch-building toggle handler.

**The shaping insight — staging is async, patch-building is sync.** Staging mutates the diff →
the content re-renders (async) → the post-stage reveal rides that render
(`RevealSelectionAfterStaging`, ordinal model, §21.17). **Patch-building mutates only the
inclusion set — the commit's diff is unchanged — so there is no content re-render.** Toggling a
line into the patch just recomputes the gutter, redraws, and advances the selection, all
**synchronously**. Re-running the pager per toggle is exactly what the draw-time gutter exists to
avoid (§21.5). So `space` in `MainViewController.stageSelectedLine` branches on the panel beneath:
files → stage (async, exists); commits/commitFiles → toggle patch (sync, new). The selection
front end (`ChangeLinesInViewRange`) is shared; only the back end and post-action differ.

**Included-set mapping (solved on paper).** `PatchBuilder` stores `includedLineIndices` (patch-line
indices into the file's parsed diff). Invert `stageDiffLines`: parse the same diff, map each
included index → change-line identity `(path, lineNumber, isDeletion)` via
`LineNumberOfLine`/`OldLineNumberOfLine`, build a set, check each rendered row's resolved identity
against it. Rendering-independent, so it holds over delta / no-pager identically. (`mode == WHOLE`
= all change lines included.)

**Gutter decisions (user, this session):**
- **On-demand.** The gutter is shown **only while a custom patch is being built** for the shown
  commit; hidden otherwise. The diff renders flush-left until the first line is toggled in (the
  intrinsic one-cell shift, §21.5), and goes away when the patch empties.
- **Inclusion-only.** It is **not** the fix for the §21.9 selection-bg-under-delta problem:
  staging from the files panel never has a gutter yet has the same problem, so that needs a
  *separate* solution (deferred). (This corrected a wrong framing — the two were lumped as "same
  class of problem" in §21.9/§21.11; they are not solvable by the same gutter.)
- **Marker = a checkmark** (maybe colored). `+` is out — it collides with raw `git diff`'s `+` for
  additions (delta has none). Trivial to change later; don't over-invest now.

**gocui rendering (the de-risk spike).** Both main views are `Wrap=true` (so `ox` is pinned to 0 —
no horizontal-scroll complication). The gutter is a **draw-time decoration**: the `View` gains a
gutter width + a per-buffer-line marker array; `draw()` reserves the left N columns, paints the
marker on the **first wrapped segment** (`viewLine.linesX == 0`) of marked buffer lines, and shifts
content right by N; `refreshViewLinesIfNeeded` reduces the wrap width by N. The buffer/cells/metadata
are untouched (the gutter is pure decoration), so `DiffLineContents`/`BufferLineForViewLine`/the
resolvers keep operating on the unshifted buffer.

**Marker recompute (lazygit side).** Markers derive from (a) the buffer's resolved diff-line
identities (stable per buffer line once loaded) and (b) `PatchBuilder`'s included set (changes on
toggle). Recompute on toggle (sync) and on content render (navigating commits/files re-renders the
main view). Streaming (markers for not-yet-loaded lines of a large diff) is a known follow-up —
likely recompute on render-complete; the toggle-recompute covers the core demo for the spike.

**Decomposition (linear):** 6a — the gocui reserved-column gutter (de-risk first; gocui unit test).
6b — single-line + hunk toggle, single file, wired to `PatchBuilder` + the sync selection advance.
6c — range + multi-file (reuses the step-4 group-by-file machinery). Whole-commit (multi-file) diff
and `enter`-focuses-the-main-view bleed into step 7.

### 21.20 Session 17 (cont.): step 6 built — gutter + patch toggle (6a/6b done)

Committed (most recent last; fixups noted):

```
67906ae47 Add an on-demand inclusion gutter to gocui views                       (6a; +fixup intrange)
f8b37370d Add identity-based line accessors to PatchBuilder                       (prep)
185ddeea9 Extract stageRange from stageSelectedLine in the main view             (prep)
1aed9428a Build a custom patch from the focused main view of a commit's files     (6b; +fixups: e2e test, reset stale gutter on preview re-render)
```

**6a — the gocui gutter (de-risk spike, clean, no fight).** `View.SetInclusionGutter(show, marks
[]bool)` reserves a left column, paints `InclusionGutterMarker` ("✓", green) on the **first wrapped
segment** of marked buffer lines, shifts content right, and narrows the wrap width while shown. Pure
draw-time decoration — the content buffer, metadata, click resolution and wrap inputs are untouched;
both main views are `Wrap=true` so `ox` is pinned to 0 (no h-scroll complication). gocui unit tests:
`TestInclusionGutter`, `TestInclusionGutterMarkerOnFirstSegmentOnly`. The de-risk verdict: the
reserved-column-over-pager rendering is straightforward — the spike did not fight us.

**6b — patch toggle from the commit-files main view.** `space` routes to staging (files panel) or
patch toggle (commit-files panel) by **which handler the panel registers** — a new
`onTogglePatchFocusedMainView` channel, kept separate from `onStageFocusedMainView` because the
post-action differs (sync gutter repaint vs async reveal); its presence also gates the gutter.
`CommitFilesController.GetOnTogglePatchFocusedMainView` resolves the selection to change-line
identities (the shared `ChangeLinesInViewRange`), maps them to patch-builder indices
(`PatchBuilder.PatchLineIndicesForLines`), decides add/remove from the first selected line (like the
explorer's `toggleSelection`), toggles, then repaints the gutter — no refresh, so the pager is never
re-run. It starts the patch builder if inactive (with the discard confirmation when a patch for a
different commit is active). The included-set → gutter mapping is `StagingHelper.RefreshInclusionGutter`:
resolve each rendered change row's identity and mark it when `PatchBuilder.IncludedLineIdentities(file)`
contains it. Recomputed on focus (`focusMainView`) and on toggle; reset on preview re-render
(commit-files `GetOnRenderToMain`). End-to-end test `patch_building/build_from_main_view` (toggle one
hunk → apply → only that hunk lands). unit + lint green; `e2e-all` green modulo the known
empty-pathspec flake (`diff/diff_non_sticky_range` etc.).

**The §21.5 / §21.19 insight held:** patch-building is **synchronous** — a toggle changes only the
inclusion set, the commit's diff is unchanged, so the gutter repaints in place and the pager is never
re-run. This is the whole reason the draw-time gutter exists, now realized.

**Needs interactive sign-off** (the gutter is draw-time, so not e2e-assertable, and the harness has no
real delta): the gutter look/feel under delta + no-pager + difftastic; the checkmark glyph/color
(trivial to change); and the on-demand appear/disappear as the patch is created/emptied.

**Known limitations (carry forward):**
- **Scoped to commit-files.** Commits / sub-commits / stash main views (the whole-commit multi-file
  diff) register no toggle handler yet → `space` is a no-op there, no gutter. That's 6c / step 7.
- **No auto-advance after a toggle** — the selection stays on the toggled hunk (so `space` again
  toggles it off; navigate by hand). The explorer auto-advances
  (`SelectNextStageableLineOfSameIncludedState`); deliberate first-cut simplification (and it sidesteps
  the removed `AdjacentChangeLine`).
- **The browser ◐/● indicators and the secondary patch summary don't update live on a toggle** — the
  toggle deliberately does *not* refresh `COMMIT_FILES`/`PATCH_BUILDING`, because that refresh
  re-renders (re-runs the pager on) the main view via the `NextInStack` branch and would reset the
  gutter + scroll. They catch up the next time the browser/secondary naturally renders. A real gap;
  productionization (step 8) wants a way to update a list model without re-rendering the main view.
- **Streaming**: marks for not-yet-loaded lines of a large diff appear once loaded + recomputed
  (next focus/toggle).
- `togglePatchLines` already groups by file (multi-file-ready), but only commit-files registers the
  handler and the add/remove direction is decided once from the first selected line.

**NEXT (6c / step 7):** register the toggle handler on commits / sub-commits / stash main views
(whole-commit multi-file diff); then step 7 (`enter` focuses the main view for files AND commit-files)
and step 8 (tear out the explorer views + escape machinery). Consider auto-advance and live
browser/secondary updates as polish.

### 21.21 Session 17 (cont.): toggle refreshes normally — secondary + indicators update live

The §21.20 "browser indicators / secondary patch summary don't update live on a toggle" was the
wrong tradeoff (user: it makes patch-building unusable, and updating just the list model wouldn't fix
the secondary anyway). **The toggle now refreshes normally** (`Refresh({COMMIT_FILES})`), which
re-renders the commit-files browser (the file's ◐/● status), the secondary view (the cumulative
patch — and the layout splits to show it when the patch first becomes non-empty), and the main diff.

**Scroll/selection are preserved for free** — the realization the user pushed on: a patch toggle
re-renders the **same diff command**, and the render pipeline already keeps the scroll when the
command is unchanged (`pty.go`: `ResetOrigin = restore == nil && cmdStr != GetTaskKey()`). So no
`PreserveDiffPositionOnRerender` is needed (unlike the `-U` context-size change, which *changes* the
command and so must preserve explicitly). The selection (cursor + range anchor) survives the re-render
too — it's view state the task doesn't touch.

**The gutter rides the same re-render** without the restore machinery, by recomputing in
`GetOnRenderToMain` keyed on the **same content-equality test the pipeline uses**: if the diff command
equals the manager's current task key (a toggle re-rendering the same diff), recompute the gutter over
the current content — identical to the incoming content, so the marks are valid and survive the swap;
otherwise (a different file is loading) clear the gutter, recomputed when the main view is focused.
This replaces the §21.20 unconditional reset-on-preview-render. (Coupling note: the equality test
duplicates `newPtyTask`'s `strings.Join(cmd.Args, " ")` key derivation — flagged for productionization.)

Test `build_from_main_view` strengthened: after the toggle it asserts the selection is unchanged
**and** the secondary view shows the cumulative patch (the toggled block only). unit + lint + the
patch_building/staging suites green.

**This removes the big §21.20 limitation.** Remaining (carry forward): scoped to commit-files; no
auto-advance after a toggle; streaming. Committed as fixups on `1aed9428a` (note `0ce3cc378`'s
unconditional reset is superseded by the conditional version — fold accordingly).

### 21.22 Session 17 (cont.): first interactive feedback on the gutter — 4 of 5 fixed

User tested step 6 interactively ("works mostly well"). Five issues; fixed 1, 2, 3, 5, recorded 4 for
later (user's call: 6c matters more).

- **(1) gutter/secondary visibility disagreed.** Gutter was gated `Active && !IsEmpty`, secondary on
  `Active` — so unstaging the last hunk left the secondary (empty patch) but hid the gutter. Now the
  gutter follows the secondary's criterion (dropped `IsEmpty`).
- **(3) gutter shown inconsistently while unfocused.** Decided (user): the gutter is a **focused-main-
  view affordance** — show it only while the main view holds focus. So the model is now: *gutter
  visible iff a patch is active AND the Normal main view is the current context.* `RefreshInclusionGutter`
  became parameterless and self-contained (targets `Contexts().Normal`, checks `CurrentStatic()`, then
  `NextInStack` for the panel — **focus checked first, since `NextInStack` panics if the context isn't
  in the stack**, which it is exactly when current). Recompute moved to `MainViewController.GetOnFocus`
  (+ a new `GetOnFocusLost` that hides it), dropped from `focusMainView`. This also **removed the fragile
  `strings.Join(cmd.Args)` content-equality hack** from §21.21: while focused, the only re-render is the
  same-diff toggle, so recomputing over the current content is always valid; unfocused → hidden.
- **(2) selection drifted when the split appeared** (user: "looks pretty broken"). Root cause: the
  focused-main-view selection is **view-line-based**, so when the secondary view appears the main
  halves in width, the diff re-wraps, and the cursor lands on the wrong content (the gutter marks are
  fine — they're keyed by *buffer* line, width-independent). Fix: after the toggle's re-render, re-select
  the same content by **change-line ordinal** (reuse `RevealSelectionAfterStaging`; the line isn't
  consumed so the ordinal is unchanged, and ordinals are width-independent), re-expanding the hunk. The
  handler installs it after the `Refresh` (so it works for both the sync path and the async
  discard-confirm path, and never lingers — a `Refresh` always follows). A range collapses to a line,
  as a staged range does.
- **(5) checkmark only on the first wrapped segment.** Now drawn on every wrapped segment (dropped the
  `linesX == 0` guard); gocui test renamed `TestInclusionGutterMarkerOnEverySegment`.
- **(4) DEFERRED — pager switch shifts the checkmarks.** Switching pagers keeps the **same git command**
  (the pager is an env var, not an arg) but yields a **different buffer-line structure**, so marks
  computed from the pre-switch buffer no longer align. The proper fix recomputes the gutter from the
  **new** buffer at render completion — needs a post-render hook we don't have, and pager-switching
  mid-build is uncommon. Left unfixed (a refresh fixes it); revisit if it annoys.

All fixes are fixups on `1aed9428a` / `67906ae47`; `e2e-all` green (modulo the empty-pathspec flake).

**NEXT: 6c** — register the toggle handler on commits / sub-commits / stash main views (whole-commit
multi-file diff). Then step 7, step 8.

### 21.23 Session 18 (2026-06-19): step 6c done — patch building from the whole-commit main view

`space` now builds a custom patch straight from the **whole-commit (multi-file) diff** the
commits / sub-commits / stash main views show, without first diving into the commit files panel.
Committed (most recent last; the multi-file test is a fixup on the behavior commit):

```
66224c84f Extract the patch-toggle back end out of CommitFilesController        (prep)
f9669f6f1 Build a custom patch from the commits / sub-commits / stash main views  (6c; +fixup: multi-file e2e)
```

**Prep — the toggle back end is now panel-agnostic.** The 6b handler lived entirely on
`CommitFilesController`. Pulled the panel-agnostic parts (`togglePatchFromFocusedMainView` skeleton,
`togglePatchLines`, `revealSelectionAfterPatchToggle`, `patchFilename`) into shared free functions in
`patch_building_from_main_view.go`. The skeleton now takes the **patch target**
(`from/to/reverse/canRebase`) and a **refresh callback** as params instead of reading the commit files
context + hardcoding `Refresh({COMMIT_FILES})`. The commit files handler passes exactly what it did
before — behaviour-preserving (`build_from_main_view` still green). It also stopped going through
`CommitFilesHelper.StartPatchBuilder()` and now calls `PatchBuilder.Start(from,to,reverse,canRebase)`
directly with the passed values (the patch builder is self-contained after `Start`, §patch_builder).

**6c — the handler lives on `SwitchToDiffFilesController`** (already bound to exactly
LocalCommits/SubCommits/Stash, and already deriving the ref/range/canRebase for these panels in
`enter()`). Two decisions, made with the user:
- **Patch target: derived directly from the panel's selected ref**, *not* by re-pointing the commit
  files context (user: "go with whatever is easier" + a half-baked feeling the patch-from-commits UX
  may change). Keeps the commits toggle decoupled from `CommitFilesContext` — less to unwind if the UX
  shifts. To avoid duplicating the from/to arithmetic, extracted `context.FromAndToForDiff(ref, refRange)`
  (shared with `CommitFilesContext.GetFromAndToForDiff`) and a private `canRebase(ref, refRange)` shared
  with `enter()`. The target matches the diff the main view shows (`git show <hash>` for a single commit
  ≈ `git diff <hash>^ <hash> -- <file>` for non-merge commits — the same assumption the explorer relies
  on; identities are pager/header-independent).
- **Refresh: the cheap one** (user: refreshing commits per hunk-stage is too expensive, and unneeded —
  no per-file indicators to rerender here, unlike commit-files). The callback is
  `OnUIThread(PostRefreshUpdate(panelContext))`, re-rendering just that panel's **own** main + secondary
  views (same diff command → scroll preserved; gutter rides via the panel's `GetOnRenderToMain`), with
  **no commit-list reload**. Order stays toggle → refresh → reveal: both refresh callbacks *queue* the
  render (Refresh via its OnUIThread; ours explicitly), so the reveal-by-ordinal restore is installed
  before the queued render consumes it — works for the sync path and the rare discard-confirm popup path.

**Sub-commits and stash didn't render the secondary patch view at all** (only LocalCommits did); gave
them the same `secondaryPatchPanelUpdateOpts` + a `RefreshInclusionGutter()` in their `GetOnRenderToMain`
(also added the gutter call to LocalCommits'), so the cumulative patch + gutter track the toggle. The
gutter (`RefreshInclusionGutter`) was already generic — it lights up for any panel beneath `Normal`
registering `GetOnTogglePatchFocusedMainView`, so it now works for all three for free. Routing is safe:
only Files registers a *stage* handler, so these panels hit the toggle path in `stageSelectedLine`.

**Tests:** `build_from_whole_commit_main_view` (toggle one file's hunk from a two-file commit diff in the
sub-commits main view → only that file lands, the other stays out) and `build_multi_file_from_whole_commit_main_view`
(hunk-mode `↓` moves across the file boundary to file2's block; toggle both → patch spans both files →
apply lands both). The multi-file-range/accumulation path (`togglePatchLines`' group-by-file) was
untested by e2e in 6b too; now covered. unit + lint + the patch_building/staging/commit/stash/rebase/
cherry-pick/diff suites green. (Files line-count after apply uses `ContainsLines`, not exact `Lines` — a
transient `.patch` file in the repo dir can add a row.)

**Signed off (user): "works very well."**

**Crash fixed (same session, fixup on `f9669f6f1`).** 6c broke a latent invariant: before it, an
**active patch implied the commit files context was set up** (the only way to start a patch was to
enter the commit files panel, which `ReInit`s its ref). 6c starts the patch builder directly from the
commits panel *without* touching `CommitFilesContext` (the deliberate decoupling), so its ref stays
**nil**. Then any `Refresh({COMMIT_FILES})` while such a patch is active — e.g. ctrl+p → **Reset patch**
(`PatchBuildingHelper.Reset`) — hit `refreshCommitFilesContext` → `GetFromAndToForDiff` →
`FromAndToForDiff(nil, nil)` → `ref.ParentRefName()` on a nil ref → SIGSEGV. Fix: guard
`refreshCommitFilesContext` to no-op when the context has neither a ref nor a refRange (nothing to
load). Correct independent of 6c, and covers the whole class (the rebase-y menu actions operate on the
commits model + patch builder, not `CommitFilesContext`, and their post-op refreshes route through the
same guard). `FromAndToForDiff` keeps assuming a valid ref — a crash there is more debuggable than a
silent empty range, and every other caller has a ref by construction. Regression test
`reset_patch_built_from_main_view` (build from sub-commits main view → never enter commit files → ctrl+p
Reset → no crash, patch gone); it reproduced the exact panic before the guard.

**Needs interactive sign-off** (draw-time gutter, no real delta in the harness): the gutter / secondary /
on-demand appear-disappear on the **whole-commit multi-file** diff under delta + no-pager + difftastic;
and the **LocalCommits** path specifically (canRebase=true → the full apply-options menu, not just
"Apply patch"; the e2e covers SubCommits, canRebase=false). Carry-forward limitations unchanged from
§21.22: no auto-advance after a toggle; streaming; the §21.22(4) pager-switch checkmark shift.

### 21.24 Session 18 (cont.): revised forward plan — drop steps 7/8, add explorer commands + dispatch refactor + pager fallback

**Steps 7 and 8 are dropped** (user). Rationale: (1) little to learn — `enter`-focuses-the-main-view and
tearing out the explorer/escape/snapshot machinery are integration mechanics, not prototype unknowns;
(2) keeping the old staging / patch-builder explorer views is actively useful as a **behavioral
reference** while testing the merged view — when something feels off, an A/B against the old panels tells
us fast whether it's a regression or always-been-thus. So `enter` keeps its current meaning, the explorer
views stay, and the §12.2 escape/snapshot machinery stays untouched (no teardown). Cost: we don't prove
the "main view fully replaces the explorers" endgame — but that's a production integration detail, not a
prototype unknown.

Instead the prototype gains the staging/patch-building commands the merged view still lacks, plus two
productionization de-risks. New sequence (linear):

**(A) Dispatch refactor — collapse the per-command delegation channels into one.** PREREQ for B/C; do
first. Today each focused-main-view command the side panel handles needs its own channel: a
`types.Context` interface method + a `BaseContext` field + an `attach.go` registration + a
`baseController` default (`onClickFocusedMainViewFn`, `onStageFocusedMainViewFn`,
`onTogglePatchFocusedMainViewFn` so far) — ~4 touch points per command, doesn't scale to `d`/`ctrl+o`/…
(the §21.9 debt). Replace with ONE channel: a side-panel context exposes `GetFocusedMainViewActions()`
returning a single interface (nil for non-actionable panels) with a method per command (`OnClick`,
`PrimaryAction` = stage|toggle, `DiscardSelection`, `CopySelection`, …). `MainViewController` becomes a
thin dispatcher (owns the keybindings + GUI mechanics — select mode, reveal — fetches the handler from
the panel beneath, calls the method). Adding a command = one interface method + implement in the 2-3
controllers + one keybinding; no plumbing.
- **Keep the underlying BaseContext function-attachment mechanism** (`AddOnFocusFn` & friends). The user
  half-wanted to replace that whole mechanism with something better, but explicitly **scratched it**:
  the prototype is already huge, no concrete better mechanism in mind, out of scope. We collapse to one
  channel *using* the existing mechanism, not replace the mechanism.
- **Sub-decisions for implementation:** (i) unify signatures to `error` — push step-5's `focusViewName`
  focus-follow into the handler (calling the shared reveal/focus helpers) so the dispatcher stays dumb;
  it's a real change to that path, not free. (ii) Partial support is mostly polymorphism, not optionality
  (`d`/`ctrl+o` work on every diff panel, just different backends; `PrimaryAction` = stage for files,
  toggle for commits). The one real partial case is reflog — see below.

**(B) `d` — discard a hunk.** Different backend per panel: working-tree discard (files) vs the
patch-building variant (commits/commitFiles). Pin the exact semantics in each at implementation time; the
(A) handler design accommodates it by construction (each panel implements `DiscardSelection`).

**(C) `ctrl+o` — copy parts of a diff to the clipboard.** Supported by every diff panel.

**(D) No-OSC-pager fallback (LAST — the biggest production unknown).** The focused main view *is* the
staging view, so it must be stageable; a pager that emits no metadata forces a fallback. Decided shape:
- **Fall back to the raw (no-pager) rendering at focus time** for an unsupported pager — keep the pretty
  pager output for *browsing* (the unfocused preview), switch to raw only when you focus to act. (Not
  always-raw: don't degrade browsing for everyone to spare unsupported-pager users one re-render.)
- **Detection = treat metadata as a pager capability, not a per-render observation.** Lean (user deferred
  the pick to me): a one-time **probe** — run the pager on a small synthetic diff, check for OSC 1717,
  cache the verdict, re-probe on config change — so each focus renders raw-or-pager *directly*, never
  double-rendering. (Capability is binary+config-dependent, not content-dependent, so one probe is valid.
  A lazy "probe on first focus of an unknown pager, cache" variant is an acceptable simpler alternative.
  We *could* also detect at render time without a probe — lazygit holds the raw diff, so "raw has change
  lines but rendered output emitted no metadata ⇒ unsupported" is decidable — but that double-renders on
  first focus, hence the probe.)
- **Click-to-focus on an unsupported pager: best-effort by view line number** (user's call). We can't
  resolve *which* line was clicked (that resolution IS the metadata we don't have), so on the fallback
  re-render land the selection on the **same view-line index** that was clicked — if the pager
  restructured the diff only slightly, the same Y is the best chance of landing near where you clicked —
  rather than jumping to the first hunk. Keyboard-focus is unaffected beyond the re-render. Building this
  is the point: it answers both technical feasibility and how the focus re-render *feels*.

**Reflog patch-building (note, not now).** Not supporting custom-patch-building from the reflog panel's
diff was an **oversight, not a deliberate limitation** (user). The reflog already shows a specific
commit's diff; building a patch from that commit should work and shouldn't be hard (reflog uses
`SwitchToSubCommitsController` today, so it'd want the 6c-style toggle handler wired up for it too). The
technical challenges aren't interesting, so it's deferred — but it's a real gap to close in production,
not a "can't." (Once the (A) dispatch refactor lands, this likely becomes "reflog returns the same
`FocusedMainViewActions` the commit panels do.")

**NEXT: (A) the dispatch refactor** (new session), then (B) `d`, (C) `ctrl+o`, (D) the pager fallback.

### 21.25 Session 19 (2026-06-19): (A) the dispatch refactor — DONE

The per-command delegation channels are collapsed into one. Committed (most recent last; both
behaviour-preserving):

```
b694f4573 Push the staging focus-follow into the side-panel handler            (prep)
55ba6e0d7 Collapse the focused-main-view handler channels into one             (A)
```

**Prep — unify the handler signatures to `error` (sub-decision (i)).** The two `space` handlers were
asymmetric: the patch-toggle handler did its own re-render + reveal and returned `error`, while the
staging handler returned a `focusViewName` so the dispatcher (`MainViewController.stageRange`) could do
the reveal-and-focus dance on its behalf — different return types, which blocks merging them. Pushed the
focus-follow into `FilesController` (now `PrimaryAction`): it does the post-staging reveal and pane
focus itself. Extracted the reveal/select-mode logic the two handlers shared (collapse a range to a
line, preserve the change-line ordinal across the re-render, re-expand a hunk) into
`revealSelectionAfterPrimaryAction(c, sourceViewName, targetViewName, firstLineIdx)`, with
`mainContextForViewName` resolving a view name to its `*MainContext`. This *replaced* the old
`revealSelectionAfterPatchToggle` (toggle now calls the shared helper with source==target) and deleted
`stageRange`. Dispatcher is now uniform: read the selected range, hand it to the handler.

**(A) — one channel.** Replaced the three channels (`onClick/onStage/onTogglePatch FocusedMainViewFn`,
each = HasKeybindings getter + IBaseContext `Add` + BaseContext field/getter + baseController nil default
+ attach.go registration ≈ 5 touch points/command) with a single `FocusedMainViewActions` interface
(`OnClick`, `PrimaryAction`) exposed via `GetFocusedMainViewActions()` (nil for non-actionable panels).
The 2-3 controllers implement the interface directly and `return self`; `MainViewController` fetches the
actions from the panel beneath and calls the method. **Adding a command is now: one interface method +
implement in the 2-3 controllers + one keybinding — no plumbing.** Kept the underlying BaseContext
function-attachment mechanism (per §21.24, the user scratched replacing it).

**The patch-building classification (gutter signal).** `RefreshInclusionGutter` used
`GetOnTogglePatchFocusedMainView() != nil` as a proxy for "the panel beneath builds a custom patch" (so
the inclusion gutter shows beneath commits/stash/commitFiles but not the staging files panel, nor reflog
which has no toggle today). Collapsing the channels erased that proxy. First restored as a separate
marker interface, then (user's call) folded into the existing `DiffMainViewContext`: its bare
`IsDiffMainViewContext()` marker became `GetDiffMainViewType() DiffMainViewType` returning
`None | Staging | PatchBuilding` — one classifier instead of two marker interfaces. The gutter shows iff
the panel's type is `PatchBuilding` (commitFiles / localCommits / subCommits / stash); files = `Staging`,
reflog = `None` (diff shown, selection editable, but no `space` action — the deferred patch-building
gap). Only `PatchBuilding` is read today; the other values document the taxonomy and keep reflog honest
(it'll flip to `PatchBuilding` once wired).

**Tests:** `just build` + `unit-test` + `lint` + **`e2e-all` all green** (one transient
`build_from_main_view` failure under parallel load on the first batch run — the known empty-pathspec
flake; passed in isolation and on every rerun). The staging/patch_building from-main-view suite
(`stage_{hunk,range,range_spanning_files}_from_main_view`, `select_*_from_main_view`,
`build_{,multi_file_}from_whole_commit_main_view`, `build_from_main_view`, `reset_patch_built_from_main_view`)
covers the dispatch + reveal paths and stays green.

**Carry-forward unchanged:** the gutter is draw-time so no e2e covers its appearance (interactive
sign-off still pending per §21.23); reflog patch-building still a deferred gap (and now trivially:
reflog would return the same `FocusedMainViewActions` + implement the marker).

**NEXT: (B) `d` (discard a hunk)** — different backend per panel (working-tree discard vs patch-building
variant); the (A) handler accommodates it (add `DiscardSelection` to `FocusedMainViewActions`). Then (C)
`ctrl+o` copy, (D) the pager fallback.

### 21.26 Session 19 (cont.): (B) `d` — discard a hunk — DONE

`d` in the focused main view now discards the selection, mirroring each explorer. Committed (most recent
last):

```
3719a0412 amend! Generalize FilesController.stageDiffLines to applyDiffLines   (prep — folds into 2414f3eab)
19d4b8f66 Discard the selected diff line(s) from the focused main view          (B + 2 e2e tests)
```

**Prep — `applyDiffLines` generalized.** `stageDiffLines` used one `reverse` flag for *both* "which diff
to read" (it's the `cached` arg of `WorktreeFileDiff`) and "reverse the apply"; that only works because
those coincide for stage/unstage. Working-tree discard breaks it: read the **unstaged** diff (the side
shown), but apply a **reverse, not-cached** patch. So split them — `applyDiffLines(file, infos,
sourceCached, ApplyPatchOpts)`. Found via the failing staging-discard e2e (no lines matched, because it
was reading the empty staged diff). The prep commit as first written had the wrong signature; corrected
in place with an `amend!` so the prep introduces the final shape (per [[clean-history-no-back-and-forth]]).

**(B) — `DiscardSelection` + disabled reason on `FocusedMainViewActions`** (user picked full
explorer-parity disabled-reason UX over silent-no-op/handler-error). The dispatcher binds `Universal.Remove`
and forwards to the panel beneath; per panel:
- **files** → working-tree discard, mirroring the staging view's `d`: reverse-apply **not cached** on the
  unstaged side (destructive → confirm unless `SkipDiscardChangeWarning`), reverse-apply **cached** on the
  staged side (= unstage). Always available (zero-context = inline error, like the staging view), so its
  `DiscardSelectionDisabledReason` is nil.
- **commit panels** (commitFiles + the whole-commit diff of localCommits/subCommits/stash) → remove the
  selected lines from the commit via rebase, mirroring the patch builder's "discard lines from commit":
  reset any active patch, build a one-off patch from the selection (reusing `togglePatchLines`), then
  `DeletePatchesFromCommit`. Shared free functions `discardSelectionFromCommit` + `discardFromCommitDisabledReason`
  in `patch_building_from_main_view.go`, paralleling the toggle. Only rebaseable on a **local branch**, so
  the disabled reason greys `d` on stash + other-branch sub-commits (never rebaseable) and mid-rebase —
  treated all the same (user OK'd; the "never applicable ⇒ hide entirely" hair isn't worth it).

**Tests:** `staging/discard_from_main_view` (discard one hunk of a two-hunk working-tree file from the
files main view; the other hunk stays) and `patch_building/discard_lines_from_commit_main_view` (discard
one line from a local commit straight from the commits main view → rebase → commit keeps the other lines).
`build` + `unit` + `lint` + **`e2e-all` all green**.

**Needs interactive sign-off / carry-forward:** the commit-discard's post-rebase re-render isn't given a
special reveal (the SYNC refresh re-renders; selection lands wherever it falls) — fine for a destructive
op, revisit if it feels off. Files-discard reveals same-pane by ordinal; the rare staged-side-collapse
case (discard-on-staged = unstage, secondary pane empties) isn't special-cased. Reflog `d` no-ops (nil
actions), like its `space`.

**NEXT: (C) `ctrl+o` (copy parts of a diff to the clipboard)** — supported by every diff panel; add
`CopySelection` to `FocusedMainViewActions`. Then (D) the pager fallback.

### 21.27 Session 19 (cont.): discard bug fixes (user review) + selection visibility

User reviewed (B) and found three discard bugs + one pre-existing selection-visibility bug. All fixed;
the §21.26 "carry-forward" items above are now resolved. Bugs 1-3 are fixups on `19d4b8f66`; bug 4 is its
own commit (`f06c64366`).

**Bug 1+2 — files discard must reuse the primary action's focus-follow.** Discard on the staged side *is*
unstaging (reverse + cached → index), so it must behave identically to `space` there: the selection
advances to the next staged hunk, focus follows the staged side (secondary when split), etc. My naive
same-pane reveal got this wrong (lost focus / selection on last-staged-hunk; moved focus to unstaged when
discarding the first hunk of an only-staged file). Fix: extracted `applyDiffLineSelection` (apply +
refresh + focus-follow + reveal) and route both `PrimaryAction` and `DiscardSelection` through it;
`diffLineSelection` resolves the selection + which side. The only real difference is the `ApplyPatchOpts`
(stage = forward/index, unstage = reverse/index, discard-unstaged = reverse/not-cached) and the confirm
(unstaged discard only). Test `discard_from_staged_main_view`.

**Bug 3 — commit-discard must reveal the next hunk** (I'd hand-waved "fine for a destructive op"; user
disagreed, correctly). Install `revealSelectionAfterPrimaryAction` before the rebase's
`CheckMergeOrRebaseWithRefreshOptions(SYNC)` so the ordinal-reveal rides that re-render, exactly like
staging. Verified it rides the rebase refresh (test now asserts the selection lands on the next surviving
line).

**Bug 4 (pre-existing) — hide the selection when the focused main view has no diff.** It always showed a
selection while focused, even over "No changed files" / a merge message; and a stale selection lingered
after the last change was discarded or changes vanished outside lazygit. Fix has two moments (focus
*reuses* the rendered content rather than re-rendering, so one hook can't cover both):
- **Focus:** `showInitialDiffSelection` leaves the selection off when `ViewHasChangeLines` is false.
- **Render:** the panel's render-to-main is where diff-vs-placeholder is decided (user's pointer), so
  `updateFocusedMainViewSelectionVisibility` sets the selection there — shown only on the focused pane and
  only when a diff is rendered. Covers the refresh cases focus can't: discard-to-empty + external
  changes both ways. Wired into `FilesController.GetOnRenderToMain` (placeholders → off, diff → on);
  commit panels always render a diff. New `SelectionIsShown/Hidden` e2e assertion (reads the view's
  `Highlight`, since `SelectedLines` ignores it); tests `no_selection_when_no_changes` +
  `hide_selection_after_discarding_last_change`.

`build` + `unit` + `lint` + **`e2e-all` all green.**

### 21.28 Session 19 (cont.): (C) `ctrl+o` — copy parts of a diff — DONE

`ctrl+o` in the focused main view copies the selected diff line(s) to the clipboard. Committed
`16035437d`.

**Design decision (user-confirmed, deviates from §21.24): copy is a direct `MainViewController` command,
NOT a `FocusedMainViewActions` method.** Unlike stage/discard, copy is panel-*independent* — it just reads
the text the focused main view shows + the pager config, no per-panel backend. Routing it through the
interface would (a) give three identical implementations (no real polymorphism) and (b) exclude the
**reflog**, which is a diff panel with nil actions. As a `MainViewController` binding gated on `Highlight`
it works over every diff panel including reflog, with no duplication. (So §21.24's "add `CopySelection` to
the interface" was wrong for this command; the interface is for things with per-panel backends.)

**`dropDiffPrefix` / pager (the §21.27-adjacent concern user raised).** The staging view strips the
+/-/space column when copying, safe because it's always a raw diff. The focused main view can't assume
that: under a pager the +/-/space column may or may not survive (delta drops/restructures it, ydiff keeps
it and only colorizes) and we can't tell which, so stripping the first char might eat content.
Fix: `usingExternalDiff()` (pager configured?) gates the strip — no pager → `dropDiffPrefix` (raw diff,
strip homogeneous selections); pager → copy verbatim (cost: missing the strip for column-preserving pagers). Reused the existing `dropDiffPrefix` (same-package
free function). **Deferred (user said optional):** refining the raw-diff strip to classify header lines by
diff-line metadata `Type` instead of the text's first char — that would fix commit `159bbb0825`'s
`--- a/file` → `-- a/file` edge case "for free", but it's complicated by line-wrapping (SelectedLines is
view lines; metadata is per buffer line), so not "easy" enough for now.

**Test** `copy_from_main_view` (homogeneous addition hunk → clipboard has the `+` stripped; uses the
`OS.CopyToClipboardCmd`→file emulation). The pager-verbatim branch needs interactive sign-off (no real
pager in the harness, cf. the gutter). `build` + `unit` + `lint` + **`e2e-all` green.**

**Bug fix (user review, fixup on the feature commit).** The first cut used `gocui.View.SelectedLines()`,
which is a **test-only** helper: it indexes `buf.lines` (unwrapped) by the selection's `cy+oy`
(view-line, wrapped) space, so any wrapped line → wrong content, and a view line past the buffer length →
out-of-range panic. Fixed by mapping the selected view-line range to buffer lines via
`BufferLineForViewLine` and reading `DiffLineContents()[bufferIdx].Text` — copies each wrapped line whole
and once. (Lesson: don't use `SelectedLines()`/`SelectedLine()` in production; they assume no wrapping.)
Also append a trailing `\n` so the last copied line is terminated (the pager-verbatim path and
`dropDiffPrefix`'s keep-prefix branch don't add one).

**NEXT: (D) the no-OSC-pager raw fallback** — the last and biggest production unknown (see §21.24(D)):
fall back to raw (no-pager) rendering at focus time for a pager that emits no metadata, with a one-time
capability probe and best-effort click-to-focus by view-line index.

### 21.29 Session 20 (2026-06-20): (D) the no-OSC-pager raw fallback — DONE

> **SUPERSEDED by §21.30:** the *detection* mechanism described below (observe the
> already-rendered diff at focus, infer support from whether change lines resolve) was
> replaced by a probe after the user found it broke for binary files. The rendering/bypass
> machinery (`ignoreExternalDiff`, `NewMainViewDiffTask`, the panel wiring, the raw
> re-render) is unchanged. Read §21.30 for the final detection mechanism; the commits fold
> to the probe (the observe code never lands).


The focused main view now falls back to the raw (no-pager) diff when focused under a pager whose output
we can't resolve to patch-space, so it stays stageable. Committed (most recent last; the fixup folds into
the first):

```
904a2543e Fall back to the raw diff when focusing a main view under an unresolvable pager   (files panel)
9e3bf7e07 Thread ignoreExternalDiff through the remaining diff-cmd builders                  (prep)
6a1cf98a1 Extend the raw-diff fallback to the commit/stash/patch-building panels
8c969eaf9 fixup! Fall back to the raw diff ...                                               (threading fix)
```

**Detection: DECIDED AGAINST the §21.24 synthetic OSC probe — observe-at-focus instead** (discussed and
agreed with the user). The §21.24 lean ("run the pager on a synthetic diff, grep for OSC 1717, cache")
has a fatal flaw: the capability we need isn't "emits OSC", it's "we can resolve the rendered diff's
change lines", which is metadata #2 **OR** buffer-parse #1. A pure-OSC probe would misclassify every
structure-preserving non-OSC pager (`delta --color-only`, `git diff --color`, `diff-so-fancy --patch`) as
unsupported and force needless raw re-renders, even though buffer-parse stages them fine. A correct probe
would have to run the full resolver on faithfully-reproduced pager output (PTY for the GIT_PAGER route,
ext-diff config, env), which is a lot of fragile machinery. So instead we **observe the diff lazygit
already rendered for browsing**, at focus: zero new machinery, zero fidelity gap, and it can't
misclassify.

**The verdict is the load-bearing state (not just an optimization).** Per-pager `diffPagerSupport`
(`unknown|yes|no`) on `StagingHelper`, keyed on a pager signature (index + extDiffCmd +
useExtDiffGitConfig + pager template) so it resets when the pager cycles/reloads. It's required for
**re-renders while focused** (after staging a hunk, the SYNC refresh re-renders the diff — it must render
raw again, which `verdict==No` makes happen upfront). `DiffMainViewShouldRenderRaw() = focusedOnMainView()
&& verdict==No`; every diff panel's `GetOnRenderToMain` reads it (so browse=pretty, focus=raw, and the
re-render stays raw).

**Bypassing the pager needs BOTH routes** (this was the subtle part): a pager reaches the diff either as
an external diff command (in the cmd via `--ext-diff -c diff.external=…`) **or** as a stdin pager
(`GIT_PAGER`, applied by `newPtyTask`, *not* in the cmd). So `plain` on the cmd builders doesn't suffice.
The fix: a new `ignoreExternalDiff` arg on the diff-cmd builders forces `--no-ext-diff` while keeping
git's colour (unlike `plain`, which also drops colour — we want a *coloured* raw diff for display); and
`types.NewMainViewDiffTask(renderRaw, …)` renders via a `RunCommandTask` (→ `newCmdTask`, no GIT_PAGER)
instead of `RunPtyTask` when raw. Lives in `types` because the commit-diff panels build through
`DiffHelper`, which can't reach `StagingHelper` (helpers don't cross-reference), so the panel computes
`renderRaw` and threads it in.

**Focus flow** (`establishFocusedDiffSelection`, replacing the old `showInitialDiffSelection`; both
keyboard `0` and click go through it, and `togglePanel`):
- The rendered diff resolves change lines → place the selection on it; note the pager supported (but only
  when we're *not* already showing the raw fallback — else tabbing to an already-raw secondary pane would
  wrongly flip the verdict to Yes).
- Resolves nothing + no pager → genuinely nothing to act on, no selection (the old behaviour).
- Resolves nothing + a pager → `FallBackToRawDiff`: set verdict No (so `GetOnRenderToMain` renders raw),
  install a restore, trigger `sidePanel.HandleRenderToMain()`. The restore fires at the raw render's EOF
  swap: if it now resolves, place the selection; if the raw diff *also* has no change lines (a **binary
  file** — the one genuinely-ambiguous case the user flagged), revert the verdict to Unknown (binary is no
  evidence against the pager) and show no selection.
- **No placeholder problem** (the "No changed files" worry): those are string tasks, never the pager, so
  `GetOnRenderToMain` never even reaches the diff branch — the fallback only ever applies to the real
  `git diff` task.
- **Click-to-focus by view-line index** (§21.24's best-effort): the click's view-line index is replayed
  on the raw re-render (clamped), since we can't resolve *which* line was clicked under the bad pager.

**Threading fix (fixup `8c969eaf9`).** The restore's `Apply` runs on the task goroutine (`go utils.Safe`
in tasks.go), so the verdict write + selection placement there would race the UI thread's
`GetOnRenderToMain` read. Hop to `OnUIThread` for that work (per [[main-thread-over-mutexes-direction]]);
`swapIn`/`ClearRestoreForNextTask` stay inline, matching the existing restores.

**Tests** (the real win: a pager renders in the headless harness, so this is the first CI coverage of the
pager path). `cat -n` is the ideal unsupported pager — it numbers every line (so buffer-parse fails) and
emits no metadata, and it's a GIT_PAGER-route pager (so it exercises the `RunCommandTask` bypass, not the
cmd-arg one). `staging/stage_from_main_view_with_unsupported_pager` (browse shows `cat -n` numbering →
focus falls back → stage a hunk, split persists) and
`patch_building/build_from_main_view_with_unsupported_pager` (same, toggling into a custom patch from the
commits panel). `build` + `unit` + `lint` + **`e2e-all` all green** (one `diff/diff_commits` flake under
parallel load on the first `e2e-all`; passed in isolation and on the rerun — the known parallel-load
flake, not this change).

**Carry-forward / known edges (interactive sign-off pending — no patched delta in CI, so the *feel* of the
focus re-render and the metadata path are unverified):**
- **Mid-stream focus mislatch** (accepted): focusing a *supported* pager's diff before it finishes
  streaming could resolve nothing yet → fall back → latch No → that pager renders raw-while-focused until
  the pager changes. Rare (you focus a diff you're looking at, which has loaded); degraded-not-broken.
- **Binary-first** (the user's flagged case): focusing a binary file under an unknown pager does one
  speculative raw re-render then reverts the verdict to Unknown (re-decided on the next text focus). The
  cmd-arg bypass (`ignoreExternalDiff`, for delta-as-external-diff) isn't exercised by CI — only the
  GIT_PAGER bypass (`cat -n`) is.
- **Diffing mode** (`DiffHelper.RenderDiff`, the `W` ref-compare) is **not** wired to the fallback (left
  `ignoreExternalDiff=false`); focusing a diffing-mode diff under an unsupported pager wouldn't be
  stageable. A gap, deferred (diffing-mode staging is its own question).
- **Reflog** is wired (it's a `DiffMainViewContext`), but it's read-only — the fallback only helps its
  edit/copy/navigate, not staging (it has nil actions).

**This completes the §21.24 (A)-(D) sequence.** Remaining known production gaps from the plan: reflog
custom-patch-building (§21.24, an oversight not a limitation), and the steps 7/8 explorer-teardown that
were deliberately dropped (§21.24). History on this branch is throwaway — see [[prototype-branch-throwaway]].

### 21.30 Session 21 (2026-06-20): (D) detection reworked — observe → probe for a handshake

User reviewed (D) and found the §21.29 observe approach **broken for binary files**. Repro: an unstaged
binary + an unsupported pager (unpatched diff-so-fancy). Focus → raw (ok); next 10s auto-refresh →
**back to pager** (wrong); the focused view kept flip-flopping pretty/raw. **Root cause:** the
binary-revert. I inferred "pager supported?" from diff *content* (`ViewHasChangeLines`); a binary file has
no change lines under *any* pager, so it carries no signal — yet my "raw also empty ⇒ don't latch ⇒ revert
verdict to Unknown" un-cached the verdict, so the next refresh rendered pretty again. Content-based
detection forces a lose-lose: revert → flip-flop; don't-revert → a binary-first focus mislatches a
*supported* pager. Not auto-refresh ignoring us — it respected the verdict; the verdict got reset.

**Fix (discussed + agreed): detect by PROBING the pager for a handshake, not by observing renders.** A
metadata-aware pager emits a **version-only OSC 1717 record as its first output** (a handshake announcing
"I speak the protocol"). The probe runs the pager on **empty input** and greps for it. This is a
**pager-level, content-independent** fact (a binary can't mislead it) **known before we render** (so we
never render pretty then discover mid-flight we needed raw — the bug a "No changed files" → external-change
→ auto-refresh sequence would hit with observe, since observe has nothing rendered to learn from yet).

The user pushed for the probe over my observe-on-the-handshake idea, and I conceded — their reasons hold:
(1) empty-input probe is fast → run it synchronously on first need, cached; (2) observe can't know before
the first render (the dance above); (3) **no PTY needed** — git needs a tty to *decide to invoke* a pager,
but the pager *emits* the handshake whenever `EMIT_OSC1717_METADATA` is set, so we run it directly;
(4) external diff commands: invoke with git's 7-arg diff-driver convention on two empty temp files;
(5) `useExternalDiffGitConfig` — driver is per-file via `.gitattributes`, a single diff can mix drivers,
so no one pager to probe → treated as **unsupported (always raw when focused)**, a documented limitation
(the whole-diff fallback can't express per-file mixing; raw is always stageable).

**Mechanism (commits below):**
- **gocui** (`f3e480ff6`, new): `escapeInterpreter.dropMetadataIfHandshake` drops any OSC 1717 payload with
  no fields at its terminator, so the handshake (seen on *every* render of a conforming pager) is swallowed
  — no phantom line, no bleed onto the diff header. Per-line payloads always have fields, so they're kept.
  Unit test `TestDiffLineMetadataHandshakeSwallowed`.
- **probe** (`DiffCommands.ProbePagerEmitsDiffMetadata`): empty-input run of the pager (stdin pager via
  `NewShell`) or the external diff command (7-arg convention on two empty temp files), env
  `EMIT_OSC1717_METADATA=V1`, grep stdout for the `\x1b]1717` handshake. `useExternalDiffGitConfig` → false.
- **verdict** (`StagingHelper`): `pagerMetadataSupport *bool` cached per pager signature, probed lazily on
  first need. `DiffMainViewShouldRenderRaw() = focusedOnMainView && pagerConfigured && !supported`.
- **focus flow simplified**: the observe/`FallBackToRawDiff`/binary-revert/`NotePagerSupportsStaging` are
  gone. `establishFocusedDiffSelection`: if `DiffMainViewShouldRenderRaw` (probed, known up front)
  → `RenderFocusedMainViewRaw` (install a restore that places the selection after the raw re-render — still
  `OnUIThread`, since `Apply` runs on the task goroutine — then trigger the side panel); else place on the
  diff directly. `placeOrHideInitialDiffSelection` hides the selection when there are no change lines
  (binary/placeholder), so a binary under an unsupported pager is **stable** (raw, no selection, no
  flip-flop) — the verdict no longer depends on content.

**Commits:** `f3e480ff6` (gocui swallow, new) + `75c98e4e5` (fixup → `904a2543e`, the probe rework) +
`0b8fefeb3` (amend! → `904a2543e`, rewrites its body from observe to probe; subject unchanged so the two
fixups still match). After the user folds, `904a2543e` is the probe-based fallback; observe never lands.

**The handshake trade (user agreed):** the handshake signals metadata-support (#2) only, not
buffer-parseability (#1). So a *structure-preserving but non-emitting* pager (unpatched `delta
--color-only`) now falls back to raw **while focused** (loses its pretty rendering during focus; browse
unchanged), where §21.29's observe kept it pretty via #1. Predictable opt-in rule ("speak the protocol to
be first-class in the focused view"); diff-so-fancy restructures anyway, so it's unaffected.

**Tests:** `cat -n` (unsupported, GIT_PAGER route → exercises the probe's "no handshake" verdict + the
`RunCommandTask` bypass) keeps the two §21.29 tests; new `staging/stage_from_main_view_with_conforming_pager`
— a fake pager `printf '<handshake>MARKER\n'; cat` (handshake + a marker line + passthrough) — proves the
probe trusts it (no raw fallback: the MARKER survives focus) and the handshake is swallowed cleanly. `build`
+ `unit` + `lint` + **`e2e-all` green**.

**Scope beyond this repo (not done here):** the OSC spec must define the handshake record, and the **three**
emitters (delta, difftastic, and now diff-so-fancy — OSC support added to it in a separate session, to be
recorded in that flow) must emit it. Those are other repos/sessions; CI here covers the unsupported path
(`cat -n`) and the supported path (the fake conforming pager), so no patched pager is needed to test.

**Still pending interactive sign-off** (no patched real pager in CI): the *feel* of the focus re-render,
and a re-check of the user's binary repro now that the flip-flop root cause (the content-based verdict) is
removed. Carry-forward edges from §21.29 that the probe **resolves**: the mid-stream mislatch (the verdict
is no longer derived from a render, so an incomplete render can't mislatch) and the binary flip-flop.
Carry-forward edges that **remain**: diffing-mode (`W`) not wired; reflog wired but read-only;
`useExternalDiffGitConfig` always-raw.

### 21.31 Session 21 (cont.): the three pager emitters patched + the spec

Closes the cross-repo loop §21.30 left open — there were no conforming pagers, so every real pager fell
back to raw. The handshake is now defined in the spec and emitted by all three reference pagers:

- **spec** (`diff-line-metadata-osc-spec.md`, this repo, commit `323b8224e`): new §4.4 defines the
  version-only handshake record (`ESC ] 1717 ; <version> ST`, no further fields), emitted first; §6 notes it
  precedes the per-line records. A host distinguishes it from a per-line record by field count.
- **delta** (`/Users/stk/Stk/Dev/Builds/delta`): emit the handshake at the top of `delta()` before
  consuming the diff; flushes via the normal-return writer drop. + a `handshake_for_version` unit test.
- **difftastic** (`/Users/stk/Stk/Dev/Builds/difftastic`): emit once (a `std::sync::Once`) at the top of
  `print_diff_result`, **explicitly flushed** (difftastic exits via `process::exit`, which skips the
  implicit flush — without it the lone handshake of an empty-diff probe is lost). + a unit test.
- **diff-so-fancy** (`/Users/stk/Stk/Dev/Builds/diff-so-fancy`): emit before the main stdin loop;
  `osc_records` test helper now skips the handshake so the per-line assertions are unchanged, + two new
  bats tests (handshake first; handshake on an empty diff).

Each was verified manually to emit `\x1b]1717` on the *exact* empty-input invocation the lazygit probe
uses (delta/diff-so-fancy on empty stdin; difftastic via the 7-arg git diff-driver convention on two empty
temp files), and to stay byte-silent without `EMIT_OSC1717_METADATA`. So the lazygit probe
(`ProbePagerEmitsDiffMetadata`, which greps that invocation's output for `\x1b]1717`) now classifies all
three as **supported** — no raw fallback, stage on their pretty output via the per-line metadata (#2). The
delta/difftastic binaries must be **rebuilt** for this to take effect (`cargo build`); diff-so-fancy is a
script. The full real-pager render+stage in lazygit remains the interactive sign-off (CI still can't run a
patched pager; it covers the unsupported path via `cat -n` and the supported path via the conforming-pager
fake from §21.30).

### 21.32 Session 22 (2026-06-20): exploratory UX tweaks — hunk-on-click in the focused main view

The (A)-(D) sequence is done and the feature is functionally complete; before circulating the spec to pager
authors and writing the production plan, a round of small UX tweaks to make the focused main view feel right.
First tweak (two parts), both gated on `gui.useHunkModeInStagingView` semantics and **main-view only** (the
old staging/patch-building explorer views are left alone):

1. **Click-to-focus selects the clicked hunk.** Clicking the unfocused main view to focus it used to drop a
   single-line selection at the click; the common intent is "stage this block", so it needed a follow-up `a`.
   Now, when `useHunkModeInStagingView` is on, a click on a **change line** selects that line's whole change
   block (hunk mode), matching what keyboard focus already did. A click on **context** still
   selects just that line (so it stays editable with `e`). The single source is
   `placeOrHideInitialDiffSelection`'s `clickedViewLine >= 0` branch; new predicate
   `StagingHelper.IsChangeLine(view, viewLine)` (wraps `GetDiffLineInfoForView(...).IsChange()`) names the
   "is the clicked line stageable" test, shared by this branch and the keyboard one.
2. **Clicking another hunk keeps hunk mode.** A pre-existing wart (inherited from staging/patch-building):
   once in hunk mode, clicking another change line reset to a single line. Now a click on a change line while
   in hunk mode re-selects that whole block; a click on context — or any click when not in hunk mode — drops
   to a single line. The two click handlers (`onClickInAlreadyFocusedView`, `onClickInOtherViewOfMainViewPair`)
   shared an identical selection body, so they were unified into one `selectClickedDiffLine` helper and the
   change made once. The decision keys off the **current** select mode (preserve hunk), not the config
   (which only seeds the *initial* mode on focus) — the two compose.

**Open question resolved (user decided):** clicking a context line in part 1 stays in **line mode** rather
than snapping to the next hunk below (`selectDiffHunk`/`ChangeBlockBounds` *would* snap, so the `IsChangeLine`
guard is what holds it to the clicked line) — a click should select where you point, and it keeps the
"click a context line, press `e`" path working. **No tests** (prototype tweaks the user is feeling out; may
change). User's verdict on part 1 after manual testing: "feels much better, I like it a lot."

Follow-up bugs the testing surfaced, all fixed (no tests, prototype):

3. **First cross-pane click dropped hunk mode.** Clicking the *other* main pane (the staged
   side you weren't focused on) selected a single line the first time even in hunk mode, because that pane's
   `DiffSelectState` was still its default until it had been focused once — tabbing to it and back was the
   workaround the user found. Fix: `onClickInOtherViewOfMainViewPair` now copies the leaving pane's
   `DiffSelectState` to the clicked pane before `selectClickedDiffLine` runs, so the very first click there
   behaves like every later one. (The two panes keep independent select state; this seeds the target.)
4. **A selection's far end not preserved across a pager switch.**
   `PreserveDiffPositionOnRerender` restored the cursor by patch identity but left the selection's *other* end
   (the range anchor) pinned to its old view line; a restructuring pager (normal delta → delta **side-by-side**
   is the clear repro) then left the selection spanning the wrong patch lines. Fix (general, not hunk-specific
   — the user pushed for this): remember the far end by patch identity too (`selectionFarEndIdentity`, the
   range anchor opposite the cursor) and put it back via `SetRangeSelectStart` after the re-render, found by
   `findDiffLine` (a `findResolvedDiffLine` + `matchByPatchLine` lookup, which works for context ends as well
   as change ends). So **one** mechanism restores both ends and covers hunk *and* range/sticky selections —
   no `reselectHunkOnRerender` fixup, no callback param; `PreserveDiffPositionOnRerender(view)` stays
   single-arg. If the far end didn't survive the re-render it collapses to the single cursor line. Both callers
   (pager cycle `onPagerChanged`, `-U` context-size change) get it for free. **Scope unchanged:** only the
   primary (`Normal`) pane is preserved on a pager switch (both callers pass `Contexts().Normal.GetView()`), so
   a selection on the *focused secondary* pane still isn't preserved — a pre-existing, broader gap (its scroll
   isn't preserved either), left as-is.
5. **Click-and-drag range anchored wrong / dead on context lines.** Dragging after a click turns
   a hunk selection into a range, but it was anchored at the change block's far end (where selecting a hunk
   leaves the gocui range anchor), and on a context line — where the click leaves no range anchor — dragging
   just moved the single selected line. Root cause: a click can show a whole hunk, so the gocui range anchor
   holds the block's far end, not the clicked line, and the clicked line can't be read back. Fix: remember each
   mouse-down's view line on `MainContext.dragAnchorViewLine` (set in both click entry points —
   `placeOrHideInitialDiffSelection` for the focus click, `selectClickedDiffLine` for the already-focused/
   cross-pane click), and a new `MouseLeft`+`ModMotion` binding (`onDragInFocusedView`) re-anchors the range
   there (Mode→Range, hunk off) while gocui's own `SetCursor` tracks the dragged end. Works for the focus click
   and the already-focused click, on change and context lines. Drag anchor is a plain view line (no re-wrap
   between mouse-down and its drag, so no patch-identity needed). The old (buggy) behaviour came for free from
   `selectDiffHunk` leaving a range anchor + gocui moving the cursor on drag; nothing handled the drag
   explicitly before.
   **Follow-on, a gocui driver bug the above exposed:** the *first* drag-movement event was
   reported as `MouseRelease`, not `MouseLeft`+`ModMotion` — the tcell driver's `MAYBE_DRAGGING → DRAGGING`
   transition (`tcell_driver.go`) set the state but left `mouseKey`/`mouseMod` at their defaults. gocui still
   moved the cursor for it, but `onDragInFocusedView` didn't fire until the *second* moved-to line, so the
   first line of a drag showed the stale hunk-end anchor (range from block end → cursor) before re-anchoring.
   Symptom the user reported: "drag down one line → range from block end to here; one more line → snaps to
   anchored-at-click, correct thereafter." Fix: mark that first movement as a drag like every later one. No
   `MouseRelease` consumers exist (only `MouseLeft`/`ModMotion`).
   **Why the old staging panel (patch explorer) never showed this**, despite the same driver bug and its own
   drag-select: its mouse-down anchors the range *at the click* (`HandleMouseDown` → `SelectNewLineForRange` →
   `FocusSelection` → `view.SetRangeSelectStart(click)`), so the native anchor is correct from the start; the
   skipped first-move event still drew a correct `[click, cursor]` range (it only delayed that panel's internal
   `State` sync by one move, identical end state). The bug bites only when the anchor must *change* at
   drag-start — the main view's case, where a click anchors at the hunk's far end. So the driver fix is
   effectively a no-op for the patch explorer; it's needed only for the main view. (An earlier draft of this
   note wrongly said it "tightens" the patch explorer's drag — it doesn't; the commit message never claimed
   that.)

### 21.33 Session 23 (2026-06-23): preserve the selection across commit rewrites (not just `d`)

`d` re-establishes the focused-main-view selection on the next surviving change (§21.26/§21.27), and it feels
great — but it was the *only* path that did. Other operations that rewrite the commit under the focused main
view (move-patch-out-into-index from the custom patch menu, undo/redo right after a `d` or a patch move)
don't run through the focused-main-view action handlers, so nothing preserved the selection: the stale gocui
selection (cursor/origin/anchor) was left painted over the new content — often a large, now-meaningless range.

**Decision (with the user): a centralised, command-agnostic net, not per-command wiring.** Rather than teach
every mutating command (move-patch ×N, undo, redo, future ones) to capture+restore the selection, the focused
main view preserves it *itself*, by change-line ordinal, as the diff re-renders — the command-agnostic twin of
`revealSelectionAfterPrimaryAction`. New free function `preserveFocusedMainViewSelectionAcrossContentChange`
(main_view_controller.go, beside the per-action reveal), called from the four commit-diff side panels'
`GetOnRenderToMain` (localCommits / subCommits / commitFiles / stash) right before `RenderToMainViews`, so the
restore it installs rides the upcoming re-render. (No single controllers-layer render chokepoint exists, and
the gui-layer one — `refreshMainViews` — can't reach the controllers-package reveal; four one-line calls to one
shared fn was the clean option. Files panel deliberately excluded: stable diff command, see the gate.)

**The load-bearing gate (the design's crux, refined mid-discussion).** It stands down unless ALL hold, so it
neither fights a more precise restore nor disturbs a selection needlessly:
- the focused main view (`Contexts().Normal`) is the current context and shows a selection (`Highlight`);
- no precise restore is already pending (`manager.GetRestoreForNextTask() == nil`) — escape / post-stage reveal
  / `-U` context-size place the selection more precisely than an ordinal can, so defer to them;
- **the diff command is actually changing** (`diffTaskCommandKey(task) != manager.GetTaskKey()`). This is the
  one I missed in the original pitch: a plain background refresh re-renders the *same* commit's diff unchanged,
  and there a live range (mid-selection) must be left alone — only a command change (a rebase rewriting the
  commit's hash; `ShowCmdObj(commit.Hash())`) means the content, and the selection's position, moved. The gate
  also neatly scopes the net to exactly the commit-rewriting ops. Handles both `RunCommandTask` and
  `RunPtyTask` (pager), since both honour the restore (pty.go) and key on `strings.Join(cmd.Args, " ")`.

**Anchor = the selection's first line** (`SelectedLineRange()`'s first), matching every existing
`revealSelectionAfterPrimaryAction` caller — *not* the literal cursor end. (We'd discussed "cursor's ordinal";
this is the finer reconciliation. Identical for a single line; for a range it collapses to the top. It honours
"the selection, not the patch" either way, is consistent, and made the bespoke `d` reveal exactly redundant.)

**Why the popup case works:** `MenuContext.OnMenuPress` pops the menu (`Context().Pop()`) *before* running the
item's `OnPress`, so by the time move-patch's refresh re-renders, `Current()` is the focused main view again.
Undo is a Global keybinding, so it fires while the main view holds focus.

**Removed the bespoke commit-discard reveal** (the §21.27 bug-3 install in `discardSelectionFromCommit`): its
rebase rewrites the commit → the net now fires with the same anchor → the install was redundant. Behaviour-
preserving (the existing `discard_lines_from_commit_main_view` e2e still passes, now driven by the net).
**Files-discard keeps its reveal**: stable diff command (net never fires there) + it has staging's focus-follow
to the secondary pane. Stage / unstage / patch-toggle keep theirs too (toggle doesn't change the command).

**Tests** (both verified meaningful by neutering the net → both fail): `patch_building/
keep_selection_after_moving_patch_out_main_view` (build a 1-line patch, leave a multi-line range, move out →
range collapses to the surviving `+two`) and `undo/undo_keeps_focused_main_view_selection` (discard a line,
leave a range, undo → range collapses to the restored `+one`). Two commits (feature, then the removal).
`build` + `unit` + `lint` + **full `e2e` all green.**

### 21.34 Session 24 (2026-06-24): the delta-vs-selection bg fight resolved — `narrowSelectionHighlight`

The §21.9 / §21.9b rough edge — delta conveys add/delete/context purely by background colour, so a full-width
selection highlight takes the background over and you can't tell what kind of line(s) are selected — is now
addressed by a per-pager config option. Three tries got here:

- **`HighlightInset = 3`** (be776c653): selection bg only in the *middle*, a 3-col gap at each edge for the
  pager to show through. Hard-coded, both sides.
- **`selectionBgColorEdgeWidth`** (8367750f9): inverted — selection bg only on the two outer N columns, pager
  shows through the middle. Configurable int, per pager. Read as two fat bars bracketing the line; the user
  couldn't get used to it.
- **`narrowSelectionHighlight`** (6e4f21bbe, this session): a bool. The selection is a narrow bar (2 cols,
  fixed) at the **left edge only**; the rest of the line keeps the pager's own background. gocui field renamed
  `SelectedLineBgColorEdgeWidth` → `SelectedLineBgColorWidth` (an int "paint the left N cols"; the gui layer
  maps the bool → 2 in `applyCurrentPagerSelectionStyle`, so the magic number has one home). The predicate in
  `setCharacter` dropped the right-edge clause.

**Two alternatives considered and rejected (with the user):**

- **Blend the pager bg toward the selection colour** (lerp per cell, preserving the add/remove hue difference).
  Rejected as not general: for a 24-bit colour (delta) we have the real RGB, but for a *palette* colour the
  terminal owns the palette — lazygit can't know the actual RGB (tcell's `Hex()` only maps through its own
  *assumed* palette), and forcing an explicit RGB cell would make the selected line's tint stop tracking the
  terminal's palette as the unselected lines still do. A pager *could* emit palette colours, so blend can't be
  the general answer.
- **Dither / checkerboard** (selection bg on alternating cells ≈ a cheap 50% blend). The user tried it: looks
  terrible in practice (screen-door texture, shimmers on scroll). Out.

Left-only over both-sides: every "this row is special over coloured content" convention (git/editor change
bars, diff gutters, selection gutters) puts the marker on the left; both-sides reads as heavier framing, which
was most of what looked off about the edge-width version.

### 21.35 Session 25 (2026-06-24): the patch-building secondary pane — `space` removes correctly + it renders through the pager

Two problems with the **custom-patch secondary pane** (the `NormalSecondary` pane, which previews the patch
being built). Both committed (`72543cb4d` then `073070db9`).

**(1) `space` in the secondary used the wrong indices.** The pane became actionable in the merge (the old
explorer's secondary was inert, so this was never thought through). `space` routed through the same toggle
handler as the main pane (`togglePatchFromFocusedMainView`), resolving the selection against the secondary's
diff and mapping it to patch-builder indices **by line number**. But the secondary shows the *aggregated* patch,
which **renumbers included additions** whenever an earlier addition in the same hunk is excluded
(`transformHunkHeader` recomputes each hunk's `+start`). So the shifted number resolved to the wrong line in the
original diff — often *adding* an unrelated line instead of removing the selected one (deletions, keyed by the
stable old-file number, and fully-included hunks happened to work — hence "often does something").

Fix: resolve a secondary selection by its **ordinal among the change lines shown**, not by line number — the
custom-patch view renders exactly the included change lines in order, so the k-th shown change line of a file is
that file's k-th included change line (`PatchBuilder.IncludedChangeLineIndices`,
`StagingHelper.ChangeLineOrdinalsInViewRange`, routed via `primaryPatchActionFromFocusedMainView` →
`removePatchLinesFromFocusedMainView`). `space` in the secondary now only ever *removes*, mirroring `space` in
the staging view's staged pane. **`d` (discard-from-commit) is disabled in the secondary** (user's call — it'd
act on lines shown only as the patch, and `space` already removes them; `discardFromCommitDisabledReason` gained
a `mainViewName` arg and `DiscardSelectionDisabledReason(mainViewName)` was threaded through the
`FocusedMainViewActions` interface). e2e: `remove_lines_from_main_view_secondary` (3 adds, include 2 with a
shift, remove one from the secondary, assert the right one survived + `d` toast; it fails on the pre-fix toggle
path).

**(2) The secondary rendered raw (`FormatView`), not through the pager.** User pushed back on "infeasible for
external diff" and was right. **The unified mechanism (built, user-approved over the simpler stdin-only
alternative):** materialize the patch as two real file trees under a temp dir — `a/` = each file's `from`-side
blob, `b/` = that blob with the patch applied (`git apply`) — and render with `git diff --no-index --no-prefix
a b`, reusing the main view's pager wiring (`DiffCommands.CustomPatchDiffCmdObj` mirrors `DiffCmdObj`).
`--no-index` honors **both** `GIT_PAGER` (stdin pagers) **and** `--ext-diff`/`diff.external` (difftastic) —
verified empirically — so every pager renders it like any diff. Used **always**, even raw (no pager →
git-coloured), replacing `FormatView`; the user wanted git's own (correct) context lines, and it also makes the
**post-removal reveal work** (the secondary is now an async diff task, so the reveal's restore is consumed —
before, on a string task, it silently no-op'd).

Design decisions and gotchas (all resolved in the prototype):
- **Lifetime:** `PatchBuilder` owns the temp dir — created in `Start`, removed in `Reset` (user's framing). No
  signature caching; instead a **generation counter** bumped on every mutation drives a lazy rebuild
  (`PatchCommands.EnsureCustomPatchDiffTrees`, called from `secondaryPatchPanelUpdateOpts`). This was a
  *necessary* refinement to the user's "rebuild at the single mutation point" plan: the **old explorer shares
  the same secondary** (`secondaryPatchPanelUpdateOpts`) and mutates the patch via its own paths (`AddFileWhole`
  etc.), so a counter covering every path is required — but it still only rebuilds on change, never on
  navigation.
- **Clean paths:** the `a`/`b` dir names masquerade as git's conventional `a/`/`b/` prefixes, so `--no-prefix`
  yields output byte-identical in shape to a normal git diff with **real repo-relative paths** — meaning the
  existing buffer-parse resolver (which strips `a/`/`b/`) attributes change lines correctly with no special-
  casing.
- **Added files (the fiddly bit):** seed `a/<file>` **empty** (not absent) but leave `b/<file>` **absent**;
  apply `PatchToApply(false, false)` (added files as `--- /dev/null` creations, NOT
  `TurnAddedFilesIntoDiffAgainstEmptyFile`'s `--- a/file`, which would make the *atomic* `git apply` fail
  expecting the file to exist). Result: the displayed diff pairs `a/file`(empty) with `b/file`(content) → real
  paths, no git "added in b" header mangling. (WHOLE-mode added files return the raw `--- /dev/null` diff and
  bypass `TurnAdded` entirely, which is what forced this scheme — `RenderAggregatedPatch(true)` does *not* work
  here.)
- `git apply` and `git diff --no-index` both work **outside a repo** via `-C <tempdir>` (verified). `--no-index`
  exits 1 on differences; the task runner only logs that, no popup.
- Test churn: only `specific_selection` needed updating — git recomputes hunk context (`@@ -1,4 +1,4 @@` with 3
  context lines vs `FormatView`'s `@@ -1,6 +1,6 @@`), exactly the improvement the user wanted.

**Needs interactive sign-off** (harness has no real metadata pager, per the standing pattern): the secondary
under delta / difftastic / no-pager; that the pager-emitted metadata path resolves a secondary removal to the
right file (the masquerade's `a/`/`b/` prefix is stripped host-side for buffer-parse, but a metadata pager emits
its own path string — delta strips `a/`/`b/` by convention, so it *should* be clean, but verify); and the
post-removal reveal feel. Carry-forward limitation: building a patch in the *old explorer* then escaping to the
merged view shows a current secondary (the generation counter covers it), but `e`/click in the secondary would
resolve hyperlinks to temp-tree paths (edge; the secondary isn't an edit surface).

**Follow-ups (same session, user QoL requests after sign-off — three small commits).**

- **Temp dir under lazygit's own (`a73e0ddfa`).** `os.MkdirTemp("", …)` → `os.MkdirTemp(osCommand.GetTempDir(),
  "custom-patch-")`. `GetTempDir()` is lazygit's per-session `/tmp/lazygit-*` (created at startup, cleaned on
  exit), so the custom-patch trees honor the configured temp dir and get cleaned up with everything else.
- **Inclusion gutter stays visible when the secondary pane is focused (`59ac5e267`).** It was gated on the
  *Normal* pane being current, so tabbing to the secondary hid it. Now `RefreshInclusionGutter` shows it while
  *either* pane of the pair is current, finding the side panel via `NextInStack(current)` (not `NextInStack(Normal)`,
  which panics when Normal isn't on the stack — tabbing to the secondary evicts Normal). `GetOnFocusLost` now
  *re-evaluates* (calls `RefreshInclusionGutter`) instead of unconditionally hiding — key realization: in
  `ContextMgr.Push`/`Pop` the stack is updated to the new context *before* `HandleFocusLost` fires, so during
  focus-lost `Current()` is already the destination; RefreshInclusionGutter thus keeps it shown on a pane-switch
  and hides it only when focus truly leaves the pair (no flicker). Gutter is draw-time → not e2e-assertable
  (interactive sign-off), as before.
- **Auto-advance to the next hunk after a toggle (`da08df8a8`).** Staging advances because the staged lines are
  *consumed* from the diff (preserved ordinal lands on the next change); a patch toggle leaves the diff
  unchanged, so the same ordinal landed back on the just-toggled hunk. Fix: `RevealSelectionAfterStaging` gained
  an `advanceBy int` (reveals at `ordinal+advanceBy`, clamped); the toggle passes the toggled change-line count
  (`len(infos)`), landing on the next stageable hunk (hunk mode) or next change line (line mode); staging,
  unstaging, secondary-removal, and the preserve-net all pass 0. e2e: updated `build_from_main_view` (→ NINE),
  `build_from_whole_commit_main_view` (→ file2 BETA), `build_multi_file_from_whole_commit_main_view` (the explicit
  cross-file `NextItem` is now the auto-advance), `keep_selection_after_moving_patch_out_main_view` (navigate back
  to `+one` after the toggle to span it). Not done (possible refinement): the explorer's smarter "next line of the
  same included-state" (skip already-included hunks when adding); plain next-hunk is fine for top-to-bottom builds.

### 21.36 Session 26 (2026-07-05): rebased onto master's rename support — two rename gaps for productionization

Rebased the branch onto a newer master, which now carries `f84ada494` **"Show renamed files in the custom
patch builder."** That commit switched the commit-file loader to `--find-renames`, taught `models.CommitFile`
a `PreviousPath` (+`GetPreviousPath()`), and threaded a `previousPath` through the patch builder
(`AddFileWhole`/`AddFileLineRange`/`RemoveFileLineRange`/`GetFileIncLineIndices`/`RenderPatchForFile` all
gained the parameter) plus `transform.go`'s `StripRename` so a **whole**-file selection carries the rename
while a **partial** one strips it and points the header at the new path. This branch forked *before* that
work, so its focused-main-view patch code is rename-blind. Two gaps to fix when productionizing — **neither
fixed in the prototype** (deliberately; the prototype only needs to demonstrate the mechanism):

- **(1) The conflict resolutions hardcode `previousPath = ""`** — `patch_building_from_main_view.go:72`
  (`RemoveFileLineRange`), `:185` (`GetFileIncLineIndices`), `:200` (`AddFileLineRange`/`RemoveFileLineRange`
  via the `toggle` alias). Each is the new patch-builder parameter, filled with `""` just to make the rebase
  compile. Correct for a non-rename, **wrong for a rename**: `getFileInfo(filename, "")` loads the diff for
  the new path only, so git emits no rename record → a whole-file selection loses the rename, and
  `StripRename` (which fires only when `mode == PART && previousPath != ""`) can never trigger, so a partial
  selection can't strip it either. **Production fix:** thread the file's previous path here. Unlike the
  explorer paths, this code resolves files by absolute path off the diff lines (`patchFilename`), not from a
  `CommitFile`, so it must look the `CommitFile` up by path and read `GetPreviousPath()` — mirroring the two
  existing call sites that already do it right: `commits_files_controller.toggleForPatch`
  (`AddFileWhole/RemoveFile(file.Path, file.PreviousPath)`) and
  `patch_building_helper.RefreshPatchBuildingPanel` (`ShowFileDiff` + `RenderPatchForFile` both pass
  `file.PreviousPath`). Reflog patch-building (a separate deferred gap, §21.24/§21.29) needs the same.

- **(2) The failing e2e `patch_building/renamed_file_whole`** (added by `f84ada494`) — **not** gap (1), and
  **not** a patch-*build* regression. It drives the **old explorer** flow (CommitFiles list →
  `PressPrimaryAction` → `toggleForPatch`, which passes `file.PreviousPath` correctly), and the built patch is
  correct: the main **Patch** pane renders the full rename, and whole-file apply/removal is fine because
  `RenderPatchForFile` short-circuits to the raw `info.diff` for `mode == WHOLE && plain`. What's broken is
  only the **non-plain view rendering**: the secondary "Custom patch" preview
  (`patch_building_helper.go:136`, `RenderPatchForFile` with `Plain: false` →
  `Parse().Transform().FormatView()`) renders the rename as a mangled `diff --git a/renamed a/renamed` /
  `deleted file mode 100644` / `index e69de29..0000000` instead of the rename header + hunk. **Verified passes
  on master (`ce5a8b61b`), fails on this branch** — so it's this branch's patch-package changes (§ the
  `parse.go` hunk-length capture, `patch.go`'s new line-number helpers + `IsWellFormed`, `hunk.go`), which
  predate rename support and don't reproduce a rename header/body through the view formatter. **Production
  fix:** make `Parse`/`Transform`/`FormatView` rename-aware so `RenderPatchForFile(plain=false)` reproduces
  the rename; and re-check **partial** rename selections while at it — their *applied* patch also goes through
  this path (plain, no short-circuit) with `StripRename`, so it's exposed to the same rendering code, not just
  the preview. Bottom line: the failing test is an expected casualty of "prototype predates renames," same
  root as (1), not a surprise to worry about.

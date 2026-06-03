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

> Status at end of session: branch `use-delta-hyperlinks-for-clicking-in-diff`,
> with a pile of **uncommitted** changes implementing the escape/restore work
> (see "Uncommitted work" below). The tree builds (`just build`), is
> `gofumpt`-clean, and `just vet` passes.

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

### The mechanism we built (uncommitted)

- `types.FocusedMainViewSnapshot { SidePanel, SidePanelSelectedLineIdx,
  MainView, OriginY, SelectedLineIdx }` (`pkg/gui/types/context.go`).
- Stored on `PatchExplorerContext.focusedMainViewSnapshot` with
  `Get/SetFocusedMainViewSnapshot` on the `IPatchExplorerContext` interface
  (`pkg/gui/context/patch_explorer_context.go`). `nil` ⇒ entered the normal way
  ⇒ plain `Pop()`.
- **Capture** at entry via `focusedMainViewSnapshot(c, mainViewName, sidePanel,
  selectedLineIdx)` in `main_view_controller.go`, called at the **start** of
  each `GetOnClickFocusedMainView` (files, commit-files, commits) **before** any
  mutation that re-renders the main view. It records the side panel, its
  selected line (so we can put it back — e.g. files-panel directory→file), the
  main view context, and the main view's `OriginY` + selected line.
- **Thread** the snapshot through `FilesController.EnterFile(snapshot, opts)`
  and `CommitFilesHelper.EnterCommitFile(node, snapshot, opts)`, which set it on
  the `Staging`/`CustomPatchBuilder` context right as they push it (set on
  *every* entry so it can't leak; `nil` for the normal flow).
- **Escape**: `helpers.EscapeFromPatchExplorer(c, ctx)` (shared by
  `StagingController.Escape` and `PatchBuildingHelper.Escape`). If a snapshot is
  present: restore the side panel's selection, `Push(SidePanel)`,
  `Push(MainView)`, then on the next UI tick restore origin + selection.
  Otherwise `Pop()`.

### What works ✅

- Escape lands on the **focused main view** (main focused), content re-rendered.
- **Selection restore works perfectly when the original scroll was at/near the
  top** (selection within the initially-loaded screenful). User confirmed this
  "feels exactly as expected."
- The files/commit-files panels and the commits/stash "all the way out" routing
  are correct.

### What does NOT work ❌ (the remaining detail)

When the focused main view was **scrolled down**:

- The restored **scroll resets to the top**, and
- the **selection lands off by roughly the scroll amount** (≈ the original
  `OriginY`).

### Diagnosis (confirmed with debug logging — high confidence)

The captured snapshot is correct (e.g. `originY=92 selectedLineIdx=144`). The
problem is purely **content not loaded far enough when we restore**:

- On restore the view had **only the initial screenful (`height=37`)**. So
  `SetOrigin(0,92)` can't really show line 92, and **`FocusPoint(0,144,false)`
  returns early** (`144-92=52 > 37`), leaving the cursor where the render left
  it → "off by the scroll amount." At `OriginY=0` everything fits in the
  screenful, which is why the top case works.
- We need to **force the task to read down to (at least) the target line before
  restoring**. `ReadToEnd` is the right primitive, but every attempt mis-timed
  it against the async task lifecycle.

### The "clobber" insight (why only commits/stash is hard)

- **files / commit-files:** entering staging/patch-building renders into the
  *Staging*/*PatchBuilding* view, so the `Main` view is **never touched** — its
  content, scroll, and full loaded range survive. In principle exit there can be
  "focus main + restore selection" with **no re-render and no loading**, because
  everything is still there. (We did not specialize this yet — current code
  re-renders uniformly.)
- **commits / stash:** entering goes through `SwitchToDiffFilesController.enter()`
  which **pushes the commit-files panel**, and *that* renders the commit-files
  diff **into the `Main` view** (a different command/key) — clobbering the
  commit diff and resetting origin. So on exit we must **rebuild** the commit
  diff from scratch (re-render), and that's where the async loading race lives.

### Approaches tried and why each failed (don't repeat blindly)

1. **`OnUIThread(restore)` only (no read).** Restore runs after one tick, but
   only the screenful is loaded → scroll/deep-selection fail. *(This is the
   current reverted "Version A" state: best UX so far — main focused, selection
   works at top.)*
2. **`ReadToEnd(restore)` synchronously after the pushes.** `self.readLines` is
   `nil` at that instant (task not set up yet) → `ReadToEnd` calls `restore()`
   immediately → still `height=37`.
3. **Defer `ReadToEnd` one UI tick, with `Push(MainView)` already done.**
   `ReadToEnd`'s `then` **never fired** → restore never ran → main not even
   focused. Hypothesis: focusing the main view **stops the side panel's render
   task**, so there's no live task to read.
4. **Reorder: keep side panel focused, `ReadToEnd`, then focus main + restore.**
   `then` **still never fired**. So the "focus change kills the task" theory is
   at best incomplete — even with the side panel focused, the task isn't reading
   to end on our request. **This is the live mystery.**

### Current code state

Reverted to approach (1) ("Version A"): `EscapeFromPatchExplorer` does
`Push(SidePanel)` → `Push(MainView)` → `OnUIThread`(SetOrigin + FocusPoint +
Highlight). Debug logging removed. Main focused; selection good at top; scroll
not restored when scrolled.

### Concrete next steps to investigate

- **Understand the task lifecycle precisely** in `pkg/tasks/tasks.go`: when does
  `self.readLines` become the *new* task's channel after `Push`? Does deactivating
  a context (`ContextMgr.deactivate` / `HandleFocusLost`) **stop** the
  main view's render task? Does pushing a second context create a second task on
  the `Main` view that stops the first?
- Re-add the temporary logging (snapshot values; `manager nil?`; "ReadToEnd then
  fired, height=…"; "restore before/after oy/sel/height") and trace **which task
  is live and whether/when its `then` fires** in each of the four approaches.
- Strongly consider the **"avoid the clobber"** route (approach 2 in §6's
  options below): if entering patch building from the commits focused main view
  did **not** overwrite the `Main` view, exit would need no re-render at all and
  the whole async problem disappears. Open question: can we push the commit-files
  panel as the side panel **without** it rendering into the `Main` view (or
  render it elsewhere), given they share the `"main"` window?
- Compare with how `MainViewController.openSearch` successfully uses
  `ReadToEnd` — replicate its precondition (a single, already-live task on the
  view) rather than reading right after a `Push`.

### Options on the table (we paused to choose)

1. Find/learn the right loading primitive: "render to main, wait until loaded
   through line N, then set origin+cursor," that survives the focus dance.
2. Avoid the re-render for commits by not clobbering the `Main` view on the way
   in (then exit is the easy files/commit-files path).
3. Scope down: ship "focus main + re-render + restore selection when within the
   loaded region," accept that deep scroll resets to top; revisit later.

The user wants to **persevere on (1)/(2) together** — they don't know the task
system much better than we reconstructed it here, so it's genuinely joint
exploration. (3) is the fallback.

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

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
  prototype's identifiers and test names still say "pager"
  (`ProbePagerEmitsDiffMetadata`, `onPagerChanged`,
  `stage_from_main_view_with_unsupported_pager`, …); rename them during
  transcription. `GIT_PAGER` itself and the "stdin-pager invocation route" (as
  opposed to `externalDiffCommand`) keep their technical names where the
  mechanism is literally meant. PR 3 renames the user-facing config.
- **Threading:** master has the landed main-thread rework — see §2.8 for the
  contract every PR must honor. The prototype is already rebased on top of it.
- **Prototype references** are given as *subject line* plus short SHA. The
  prototype branch gets rebased; **find commits by subject, not SHA**. When a
  plan item says "reference: X", read that prototype commit (message + diff)
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
    raw-diff fallback at focus time; `useExternalDiffGitConfig` is always-raw
    when focused (documented limitation).
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
| "Pager" naming in identifiers, strings, docs | "diff renderer" (§1 terminology; PR 3 renames the config) |
| OSC number `456`, env vars `EMIT_OSC1717_METADATA`/`OSC1717_METADATA` | OSC **1717**, env var **`OSC1717`** (final rename, 665149b11) |
| The in-repo spec file | the spec lives on the `osc-1717-spec` branch / worktree (fe3c5ac21) |
| Session-notes commits, `.claude/settings.json` commits, WIP commits | n/a |

## 4. The PR stack (overview)

PR titles become release-notes lines — they are written for users. Order is
dependency order; 1–6 are independently releasable; 7–9 merge in quick
succession (§2.3); 10–11 any time after their dependencies.

| # | Title (draft) | Depends on | Nature |
|---|---|---|---|
| 1 | Fix flicker, scroll glitches, and crashes in async diff rendering | — | fixes, gocui/tasks |
| 2 | Internal: resolve diff lines to (file, line, kind) identities | 1 | infra |
| 3 | Rename the "pagers" config to "diff renderers" | — | rename + migration |
| 4 | Support diff renderers that emit OSC 1717 diff line metadata | 2, 3 | infra + protocol |
| 5 | Select, navigate, edit and copy diff lines in the focused main view | 2 (4 for renderers) | feature |
| 6 | Keep your position in the diff when changing context size or switching diff renderers | 1, 2, 5 | feature |
| 7 | Stage, unstage and discard changes directly from the focused main view | 4, 5, 6 | feature |
| 8 | Build custom patches directly from a commit's diff view | 7 | feature |
| 9 | Replace the staging and patch-building panels with the focused main view | 7, 8 | removal + migration |
| 10 | Alt- or shift-click a diff line to open it in your editor | 2, 4 | feature |
| 11 | Open the selected diff line in the branch's GitHub PR | 5 | feature |

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

Commits (in order):

1. **Route all view origin writes through `SetOriginX`/`SetOriginY`** — pure
   prep chokepoint refactor. Ref: b0a85eefb.
2. **Add `LAZYGIT_SLOW_RENDER` debug knob** — sleeps N ms per written line so
   async render frames become visible; inert when unset. Needed by reviewers
   to see what the later commits fix. Ref: e8682b3fd.
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
4. **Fire queued ReadToEnd callbacks when the initial read reaches EOF** —
   the read loop abandoned queued `{Total:-1}` requests when the initial
   request hit EOF, silently dropping their `Then`. Ref: b6f99abc6.
5. **Don't scroll a view up to fill blank space while its content is loading**
   — the layout clamp used the partially-loaded height; add the `loading`
   flag (an `atomic.Bool` in the rebased shape) and skip the clamp while
   `IsLoading()`. Ref: 695842291.
6. **Reset other main views' scroll after copying content, not before** —
   `refreshMainViews` zeroed the source view's origin before `CopyContent`
   used it, so every cross-pair placeholder jumped to the top. *(Verified
   against master 2026-07-18: still needed — `CopyContent` still copies the
   source's `ox`/`oy` and `refreshMainViews` still resets first; master's
   `1efcfcc14` only stopped sharing the buffer slices.)* Ref: c35c9316c.
7. **Bundle a view's cell buffer and write state into a `viewBuffer`** — prep.
   Ref: fd858cd98.
8. **Make the buffer-writing methods operate on a `viewBuffer`** — prep.
   Ref: 2cfc0e24d.
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
   scrollbar amend, N§13.5, N§13.6.
10. **Don't run end-of-input handling for a render that was stopped** — the
    stopped-task EOF coin-flip (`select` between stop and closed `lineChan`)
    let a stopped task swap in a truncated buffer. No deterministic test (the
    bug *is* the nondeterministic select) — justified skip, N§13.6. Ref:
    8e3dc3eff.
11. **Reset the scroll to the top at first paint, not when the task starts** —
    with the off-screen render the old content stays visible until the swap;
    resetting oy at task start made it jump first. Ref: 411681502.

Notes:
- If commit 9's stale-tail test needs `BufferLineForViewLine`, introduce that
  accessor here (PR 2 then reuses it) rather than contorting the test.
- Gotcha for the future: **fast renders unmask ordering transients that slow
  renders hide** (N§20.5). Re-test at normal speed *and* under slow render.

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

### PR 3 — Rename the "pagers" config to "diff renderers"

A small standalone PR, before the protocol/feature PRs so their user-facing
text is consistent. "Pager" is an implementation-detail name (`GIT_PAGER`);
what these programs actually are is diff renderers, and the OSC 1717 spec
already uses that term throughout. No behavior change.

Commits:

1. **Rename the config with migration** — `git.pagers` → `git.diffRenderers`
   (yaml key + struct/field: `Pagers []PagingConfig` → e.g.
   `DiffRenderers []DiffRendererConfig`), and the keybindings
   `keybinding.universal.cyclePagers`/`cyclePagersReverse` → e.g.
   `cycleDiffRenderers`/`cycleDiffRenderersReverse` — both via lazygit's
   config-migration mechanism (`pkg/config`; follow the existing migration
   precedents) so user config files keep working/are rewritten. Reword the
   affected i18n strings (`CyclePagers`(+Tooltip), status-bar text) in
   english.go only.
2. **Rename internal identifiers** — mechanical sweep (`PagingConfig`,
   `usingExternalDiff` neighbors, helper names) as far as reasonable; keep
   `GIT_PAGER` and the pty "stdin-pager route" naming where the mechanism is
   meant (§1 terminology).
3. **Docs** — rename/rewrite `docs-master/Custom_Pagers.md` (→
   `Custom_Diff_Renderers.md`; fix inbound links; leave `docs/` alone per
   AGENTS.md), `just generate` for Config.md/schema and cheatsheets.

Open question (resolve with the user at PR start): does the per-entry
`pager:` field keep its name? It specifically selects the GIT_PAGER
invocation route (vs `externalDiffCommand`), so a rename there trades
consistency against precision.

### PR 4 — Support diff renderers that emit OSC 1717 diff line metadata

Reads the protocol; sets the env var so conforming renderers emit. No
consumer behavior changes yet (consumers land in PRs 5–8), but the PR is the
public face of the protocol on the lazygit side — write the description for
diff-renderer authors, link the spec (osc-1717-spec branch).

Commits:

1. **gocui: parse OSC 1717 records and attach them per cell** — escape.go
   parsing; payload attached to the following cell region; metadata cleared at
   line boundaries (no bleed); **keep a content-less sentinel cell when a
   blank line carries pending metadata** (delta renders some blank changed
   lines with no cells — N§21.15 bug 1); **swallow the version-only handshake
   record** (N§21.30, `TestDiffLineMetadataHandshakeSwallowed`); multi-record
   rows: `View.DiffLineMetadataPayloads()` returns *all* distinct payloads per
   buffer line (side-by-side rows carry two — N§17.1, N§21.12). Unit tests
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
   env). Ref: 9975a8fac + 665149b11 (final name: `OSC1717`).

Cross-repo note: the reference emitters live on `osc-1717-metadata` branches
in `/Users/stk/Stk/Dev/Builds/{delta,difftastic,diff-so-fancy}`; nothing is
upstreamed yet. difftastic's emitter commits were rewritten (2026-07-18) to
emit patch-space `d`+`a` records for collapsed modification rows — rebuild
before verifying. lazygit must remain fully functional without any conforming
renderer (buffer-parse + PR 7's raw fallback guarantee this). Interactive
verification of this PR needs locally built patched renderers.

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
    docs via `just generate`. Ref: cc90accde, N§21.34.

Open item to resolve with the user during this PR: whether `n`/`N`/`f` get
proper keybinding config entries (prototype used hardcoded literals,
N§16.2) — lean: add config entries, matching lazygit convention.

Note: `space` is deliberately **not** bound here — staging arrives in PR 7.
Under a non-conforming restructuring renderer, nav/e simply no-op until
PR 7's raw fallback lands; acceptable interim (same release).

### PR 6 — Keep your position in the diff when changing context size or switching diff renderers

The `RenderRestore` mechanism plus its two standalone consumers. After this
PR: `{`/`}` (context size) and `|`/`\` (renderer cycle) keep your scroll
position and selection instead of jumping to the top.

Commits:

1. **tasks: the `RenderRestore` mechanism** — `RenderRestore{FirstPaintReady,
   Apply(swapIn)}` on `ViewBufferManager`; the read loop consults
   `FirstPaintReady()` per line (instead of the count) when a restore is set;
   **`Apply` owns the swap: resolve the target against the off-screen buffer
   first, then `swapIn()`, then set origin/selection** — this ordering is a
   real invariant (reordering reintroduces flicker for buffer-parse; guarded
   by `TestNewCmdTaskRestore`, N§20.5); `ResetOrigin = restore == nil &&
   command-key changed`; **not cleared when a task starts** (survives
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
   `installDiffLineRestore`. Refs: 506c6ea81, 24a95e965 (amend! final shape),
   0cd3a5886 (`installDiffLineRestore` extraction), N§16.1.
4. **Preserve position across `-U` context-size changes** — anchor =
   selection if shown else middle visible line; offset-preserving placement
   (same screen row); visibility guard (don't install on a hidden Normal
   view — merge-conflict edge, N§16.1). e2e: extend
   `staging/diff_context_change`-adjacent coverage. Ref: 24a95e965.
5. **Preserve position when switching diff renderers** — same one-liner in
   the renderer-cycle handler (prototype: `onPagerChanged`); fixes both the
   ext-diff top-jump and the wrong-line "preserved by raw line number" cases
   (N§18.2); graceful no-op fallback for unresolvable renderers. e2e:
   `diff/cycle_pagers` keeps passing (rename per PR 3). Ref: a21c5841a.
6. **Preserve the selection's far end too** — `selectionFarEndIdentity`
   restored via `SetRangeSelectStart`; collapses to the cursor line when the
   far end didn't survive. Ref: 0412046c4, N§21.32(4).

Known limitation (keep, document in PR): `NormalSecondary` is not preserved
(N§16.1, N§18.3).

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
   §1) runs the renderer on empty input (stdin route via `NewShell`; ext-diff
   via git's 7-arg convention on two empty temp files), env `OSC1717=V1`,
   greps for the handshake; verdict cached per renderer signature;
   `useExternalDiffGitConfig` → always raw when focused;
   `DiffMainViewShouldRenderRaw` read by every diff panel's render-to-main;
   `ignoreExternalDiff` threaded through the diff-cmd builders
   (`--no-ext-diff`, keep color); `types.NewMainViewDiffTask` routes raw
   renders through `RunCommandTask` (bypasses `GIT_PAGER`); focus flow
   installs a restore to place the selection after the raw re-render;
   click-to-focus replays the clicked view-line index (best effort). e2e
   (rename per PR 3): `stage_from_main_view_with_unsupported_pager`, the
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
   `d`+`a` while `DiffLineContents` keeps only the first payload per row.
   Production: resolve **all** payloads per row
   (`DiffLineMetadataPayloads`), match each `(type, new, old)` identity
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

Risk note: this PR is where hidden couplings surface (things that push
`Staging` contexts from unexpected places — merge-conflict flows, custom
commands, `git bisect` edge flows). Grep for every reference to the removed
contexts/views before starting; expect a long tail of small fixes.

### PR 10 — Alt- or shift-click a diff line to open it in your editor

Self-contained; after PR 4 (uses `GetDiffLineInfo`). Commits (N§19):

1. **gocui: let a mouse binding opt into firing while a popup is focused**
   (`HandleWhenPopupPanelFocused`). Ref: ac85a90ed.
2. **Extract `editDiffLine` from `editLine`** — prep. Ref: d761f07d1.
3. **gocui: carry the keyboard modifier on mouse click events** — the
   modifier-on-press fix; **global behavior change** (unbound modified clicks
   become no-ops instead of acting as plain clicks) — flag in the PR
   description. Ref: da4201aa2.
4. **The feature** — alt-left *and* shift-left both bound (no single chord
   survives Ghostty+iTerm2+VS Code — N§19.1); no focus change, no selection;
   works behind popups. Ref: a86da2e97.

Interactive sign-off: Ghostty, iTerm2, VS Code (already done once for the
prototype; re-confirm the transcription).

### PR 11 — Open the selected diff line in the branch's GitHub PR

Self-contained; after PR 5. One or two commits (N§5):

- `openPullRequestForSelectedLine` on `Commits.OpenPullRequestInBrowser` in
  the focused main view: URL `<pr.Url>/changes/<commitSha>#diff-<sha256(relPath)>R<line>`;
  commit sha from the side panel's `RefForAdjustingLineNumberInDiff`; path
  relative to **`WorktreePath()`** (never `RepoPath()` —
  [[worktree-path-vs-repo-path]]), forward slashes, exact bytes into sha256;
  branch resolution per panel (commits → checked-out; subCommits → its ref;
  commitFiles → parent). GitHub-only via `PullRequestsMap`. Ref: 912703d20.
- Unit-test the URL builder. PR description should note the anchor format is
  empirically derived (undocumented by GitHub).

---

## 6. Interactive sign-off matrix

The headless harness cannot run real diff renderers, `LAZYGIT_SLOW_RENDER`,
or the pty path (N§13.1), and the gutter is draw-time-only. Each PR needs a
user pass before merge:

| PR | What to verify interactively |
|---|---|
| 1 | Slow-render matrix (N§11/§13): flick commits/files scrolled down; 10 s auto-refresh (`refreshInterval: 3`) — no content/scrollbar flicker; **also re-test at normal speed** (N§20.5) |
| 4 | Patched delta/difftastic/diff-so-fancy emit + render cleanly; handshake swallowed (no phantom line) |
| 5 | Selection feel under delta (narrowSelectionHighlight); hunk-on-click; drag; nav under metadata delta incl. repeated `n` across files |
| 6 | `{`/`}` and renderer-cycle scrolled down: no top-jump, offset preserved, both anchor cases; ext-diff route (difftastic) |
| 7 | Full staging matrix under no-renderer / patched delta (unified + SxS) / difftastic; cross-pane focus-follow; raw fallback feel under stock delta / diff-so-fancy-without-metadata; binary-file focus stability (N§21.30 repro) |
| 8 | Gutter under delta/no-renderer/difftastic; whole-commit path on LocalCommits (canRebase menu); secondary pane preview per renderer; **secondary-pane removal under difftastic specifically** (the prototype's known-broken case: reordered `d`/`a` records, collapsed modification rows, a/b record-path leak) and under delta |
| 10 | Ghostty, iTerm2, VS Code |

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
| Rename support in the from-main-view patch paths (N§21.36(1)) | **Fix in PR 8 commit 2** (mandatory — regression vs master otherwise) |
| patch pkg rename-aware Parse/Transform/FormatView (N§21.36(2)) | **Fix in PR 2 commit 1** (mandatory; `renamed_file_whole` guards it) |
| Reflog patch-building (N§21.24) | **Fix in PR 8 commit 5** |
| Renames in the custom-patch temp trees (new, this plan) | Resolve during PR 8 commit 7 |
| Secondary-pane removal broken under difftastic — ordinal bridge + a/b record-path leak (diagnosed 2026-07-18, memory) | **Fix in PR 8 commit 7**: identity bridge instead of ordinals; path normalization decided with the user |
| Diffing mode (`W`) not wired to the raw fallback → not stageable (N§21.29) | Defer; note in PR 7 description ("diffing-mode staging is its own question") |
| `useExternalDiffGitConfig` always-raw when focused (N§21.30) | Keep; document |
| Per-pane selection memory on `<tab>` (re-anchors each switch, N§21.9) | Defer; follow-up candidate |
| `IsSingleHunkForWholeFile` hunk-default refinement (N§21.11) | **Fix in PR 5 commit 3** (mandatory — regression vs master) |
| `a` on a context line below the last hunk doesn't snap back like staging did (N§21.11) | Fix cheaply in PR 5 commit 3 if trivial (`ChangeBlockBounds` falls back to the block above); else defer |
| Deleted-file `MD`-vs-`D` staging special case (N§21.13) | **Fix in PR 7 commit 5** (mandatory) |
| `NormalSecondary` not preserved on `-U`/renderer change (N§16.1) | Keep as documented limitation |
| Gutter marks for not-yet-loaded lines of huge diffs (N§21.20) | Keep (marks appear on next recompute); note |
| Renderer switch mid-patch-build shifts checkmarks (N§21.22(4)) | **Fix in PR 8 commit 10** (mandatory — looks too broken otherwise) |
| Copy copies the renderer's output verbatim under a renderer (N§21.28) | **Fix in PR 5 commit 9** (mandatory): copy the corresponding *raw diff* lines instead — dissolves the prefix-stripping problem entirely |
| Nav only sees loaded content (deep targets in huge diffs, N§16.4) | **Fix in PR 5 commit 6** (mandatory): ReadToEnd-then-retry; free if commit 3's solution reads to end on focus |
| Toggle auto-advance: no "skip already-included" smarts (N§21.35) | Keep plain next-hunk |
| difftastic token-vs-line `c`-at-new-line mismatch (M§10.2) | Protocol v2 candidate; nothing to do host-side |

## 9. Open questions (resolve before/during the marked PR)

1. **PR 3:** does the per-entry `pager:` config field keep its name (it
   selects the GIT_PAGER invocation route, vs `externalDiffCommand`), or is
   it renamed too?
2. **PR 5:** proper keybinding config entries for `n`/`N`/`f`? (lean: yes)
3. **PR 9:** new names for `useHunkModeInStagingView` / `wrapLinesInStagingView`
   + config migration.
4. **PR 8:** the two temp-tree sub-items of commit 7 — renames in the
   temp-tree rendering, and how to normalize the `a/`/`b/` tree prefix for
   records an external diff tool emits over the temp trees (+ whether its
   "Renamed from a/… to b/…" banner is acceptable).
5. **PR titles**: drafts in §4 — the user finalizes wording at PR-open time
   (they're the release-notes lines).
6. **Cross-repo timing** (outside this plan): circulating the OSC 1717 spec,
   upstreaming the three renderer patches. lazygit ships fully functional
   without them; revisit pitching once PRs 1–8 exist as evidence.

## 10. Progress

- [ ] PR 1 — async render fixes
- [ ] PR 2 — diff-line identity primitive
- [ ] PR 3 — rename pagers → diff renderers
- [ ] PR 4 — OSC 1717 support
- [ ] PR 5 — selection & navigation
- [ ] PR 6 — position preserve
- [ ] PR 7 — staging from the main view
- [ ] PR 8 — custom patches from the main view
- [ ] PR 9 — panel removal
- [ ] PR 10 — alt/shift-click edit
- [ ] PR 11 — open PR at line

(Add per-commit checkboxes inside each PR section as work starts; record
deviations from this plan inline, dated.)

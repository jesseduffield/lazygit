# Diff line metadata — design notes

Mapping a **rendered diff row (and column)** back to its **patch-space
identity**, so lazygit can act on the line the user is pointing at.

> Status: **design only**, nothing implemented. This is a starting point for a
> future session, born out of a long design discussion. Two mechanisms are
> described (#1 a host-side parser, #2 a pager-emitted OSC protocol); they are
> complementary, not alternatives. Start with #1.

---

## 1. The primitive and its consumers

Every feature below needs the *same* one thing: given a row in a rendered diff
(and, for a mouse click, a column), recover **(file, type, source-line)** — the
exact line in the unified diff it corresponds to. It is one primitive with
several consumers, not a click-to-stage helper:

1. **Dive into staging / patch building** (`enter` on the selected line, or a
   double-click) — needs the patch line to land on.
2. **Edit the line** (`e`) — needs the new-file line to open the editor at.
3. **Open the line in the branch's GitHub PR** (`G`) — needs the side
   (`L`/`R`) and line number for the anchor. Today we always emit `R<line>`
   because we can't tell the side.
4. **Jump by hunk in the focused main view** (`<`/`>`-style, like the staging
   view already has) — needs hunk boundaries.
5. **Preserve scroll position when diff parameters change** (`{`/`}` changing
   the `-U` context size; today it jumps to the top via `onNewKey`) — remember
   the patch line at the top/middle, re-render, scroll it back into view. This
   reuses the first-paint scroll-restore machinery already built on this branch
   (`ScrollToOriginYForNextTask`, commits `054d139fe`/`625e7dbad`). Anchor on
   the **nearest change line**, which survives any `-U` change (context lines
   don't).
6. **Restore selection/scroll when escaping back from staging / patch building**
   — land on the line the explorer view was *currently* selecting at escape
   (after its auto-advance), not the line you entered on, since you may have
   staged/dropped hunks meanwhile. Replaces the brittle numeric-index restore;
   see focused-main-view-notes.md §12 (incl. the escape-routing special cases).
7. **Preserve scroll/selection when switching pagers** (cycling `git.paging.pagers`
   with `|`/`\`) — the same as #5 but triggered by a pager change instead of a
   context-size change: re-anchor on the same patch line in the new pager's
   rendering. Built; see focused-main-view-notes.md §18. The line-number anchor is
   *especially* wrong here — switching a side-by-side pager for an inline one
   restructures the diff, so the old line number points at unrelated content.

Consumers **1–4** use the primitive in the **forward** direction (rendered row →
identity). Consumers **5–7** use the **inverse** (identity → rendered row): they
scan the rendered rows' metadata for the one matching a target patch identity,
which the host does *as the buffer loads* via a predicate generalization of
`ScrollToOriginYForNextTask` (focused-main-view-notes.md §12.3). The inverse
direction is what motivates solving the §8 staleness trap up front.

Because it's one primitive, it's worth building as a clean standalone
capability rather than welding it to staging.

---

## 2. Two mechanisms, disjoint coverage

### #1 — Host-side parsing (lazygit parses the rendered buffer)

Parse the **decolorized view buffer** (gocui already exposes plain text per
line via `View.Line(y)` / `View.BufferLines()`; the cell buffer stores runes
with color stripped, so `utils.Decolorise` isn't strictly needed). Walk *up*
from the target row to the nearest `@@` (gives the hunk's new-file start) and
the nearest `diff --git a/… b/…` (gives the file), then count added/context
lines down to the row. The first character (`+`/`-`/space) gives the side.

- Reuses the `patch` package arithmetic (`LineNumberOfLine`,
  `PatchLineForLineNumber`, hunk headers). The only new piece is multi-file
  splitting (the commit diff spans files; `patch.Parse` is single-file).
- **Inherently high-fidelity**: parsing *is* working in patch space, so it
  knows the side and exact line directly — none of delta's hyperlink lossiness.
- **Works for structure-preserving renderings**: no pager, `git diff --color`,
  and `delta --color-only` **without line numbers**. You don't branch on which
  pager is configured — you just try to parse what's on screen; if it isn't a
  unified diff, the parse fails and we fall back.
- **Cannot** serve renderings that restructure the diff (delta's default mode,
  difftastic, side-by-side) — there's no unified-diff line structure left to
  parse.

> **Prototype finding — two corrections to the coverage below.** Verified
> empirically by feeding real pager output through gocui's escape parser and
> running the parser on every line (see §8):
>
> - **`delta --color-only` only qualifies *without* line numbers.** With
>   `--line-numbers` (which is exactly what users who want clickable diffs
>   enable, since the hyperlinks ride on the gutter), delta keeps the `diff
>   --git`/`@@`/`---`/`+++` headers but prefixes every body line with a gutter
>   (`  2 ⋮    │-grape`), pushing the `+`/`-` marker off the start of the line.
>   So the body reads as all-context and the naive parse is *confidently wrong*,
>   not merely absent. This is the gutter the §3.1 emit-spec already worried
>   about — for #1 it's fatal. The fix is an integrity check, not gutter-aware
>   parsing (see §8); the host stays layout-agnostic and falls back.
> - **`diff-so-fancy` (even `--patch`) is *not* a #1 case.** It rewrites the
>   headers (`modified: file`, `@ file:line @`) and **strips the `+`/`-` markers
>   entirely**, so there's no unified-diff structure left. It belongs in the #2
>   column with delta's default mode.

`git diff --color` is squarely a #1 case (it only injects ANSI into a standard
unified diff), **not** a #2 target. `git --word-diff` is the genuine odd one
out (inline `[-…-]{+…+}` markup, no per-line `+`/`-`); it breaks #1 and would
need #2 — **out of scope for now**, acceptable to leave unsolved.

### #2 — Pager-emitted metadata (an OSC protocol)

For pagers that restructure the diff (and to avoid re-parsing in general): the
pager annotates its output with per-line metadata that lazygit reads. This is
the only path for difftastic-class pagers, and the only way to get full side
fidelity out of delta's default rendering.

### Coverage

| Rendering | #1 (parse) | #2 (emit) |
|---|---|---|
| no pager / `git diff --color` | ✅ | n/a (git won't emit) |
| `delta --color-only` (no line numbers) | ✅ | ✅ if patched |
| `delta --color-only --line-numbers` | ❌ (gutter; see note above) | ✅ if patched |
| `diff-so-fancy` / `--patch` (strips `+`/`-`) | ❌ | ✅ if patched |
| delta default (color-conveyed side, gutters) | ❌ | ✅ if patched |
| difftastic / side-by-side | ❌ | ✅ **prototyped, §10** |
| difftastic / inline | ❌ | ✅ **prototyped, §10** |
| `git --word-diff` | ❌ | ✅ if patched (out of scope) |

`#2 attachment present → use it; else parse the buffer (#1); else give up
(no selection).`

---

## 3. #2 design

### 3.1 Carrier: per-cell attachment (host stays layout-agnostic)

The metadata is attached **per cell**, like OSC-8 hyperlinks, **not** per line:

- The pager emits a metadata sequence at the **start of each line** — and
  **multiple times per line** when a single row shows multiple source regions
  (twice for side-by-side, N times for difftastic).
- lazygit attaches each record to the **following cell**. If there is no
  following cell (a genuinely empty rendered line), it adds a content-less cell
  to hold the attachment — an established gocui pattern (cf. the `\n` sentinel
  cell removed in `7dc18f3eb7a4`).
- **The host never reasons about layout** — no `v.width/2`, no column math, no
  knowing whether it's side-by-side. It just reads the nearest attachment. This
  is the key property: all layout knowledge stays in the pager, which is the
  only thing that has it.

**Access rules** (a principled mirror of the app's keyboard-vs-mouse split, not
two special cases):

- **`enter`** addresses a *row* (the cursor is row-granular) → use the **first
  attachment on the row**. In side-by-side this is the left column; fine, since
  the two sides of a change are one hunk for staging purposes.
- **click** addresses a *point* → use the **nearest attachment at or to the
  left of the click x** → lands in the column actually clicked.

**Emit-position spec (interop detail):** the pager emits each region's
attachment at the **start of that region**, and everything from there until the
next region's attachment belongs to that record — *including* any line-number
gutter or other embellishments the pager considers part of the region. Where a
region "really" starts is the pager's call, not the host's (e.g. in delta's
side-by-side view everything past the `|` separator is the right side, its line
numbers included). The only firm requirement: the attachment must precede the
region's first cell, so search-left lands in the right region.

> Usage note: lazygit is keyboard-centric. Most users press `space` then
> `enter`, not click; and staging is the only click-reachable op (`e`/`G` are
> keyboard-only). So column-fidelity is a property the per-cell model gives us
> for free, not a requirement we paid much for.

### 3.2 Payload

Fields per attachment:

| field | presence | meaning |
|---|---|---|
| `version` | always | self-describing (see §3.3) |
| `type` | always | `file-header \| hunk-header \| context \| added \| deleted \| other` |
| `file` | always | absolute or repo-root-relative path (the host normalizes — pagers may emit whichever is convenient); on **every** attachment so search-left yields a complete answer without scanning back to the file header |
| `new-line` | always (content lines) | new-file line number, in the **diff's** new-file space |
| `old-line` | **only** when `type = deleted` | old-file line number |

**`type` is load-bearing and cannot be inferred.** Under the coordinate rules,
`added` and `context` *both* carry `{new-line present, old-line absent}`, so
presence can't distinguish them — and we must (scroll-preservation anchors on
change lines, so it has to tell `added` from `context`). Hence an explicit type.

**Why the side must be carried at all** (record this so it isn't "simplified"
away later): in delta's default rendering there are no `+`/`-` glyphs — side is
conveyed purely by background color. So a consumer cannot recover it from the
decolorized row; the pager has to state it.

**`new-line` is in the diff's new-file space**, so the host still runs it
through `Diff.AdjustLineNumber` before opening the editor, exactly as today
(the diff may be against staged content rather than the working tree).

### 3.3 Negotiation & extensibility (versioning, not key/value)

- **Env-var handshake:** `EMIT_OSC<n>_METADATA=V1[,V2,…]`. The **host advertises
  the versions it understands**; the pager emits the highest mutually-understood
  version. Outside a host the var is unset → the pager emits nothing → trivially
  harmless in a raw terminal / `less` / `tmux`.
- **Build the handshake in v1 even with a minimal payload.** Negotiation is the
  one piece that's *impossible to retrofit* — without it you can't introduce a
  v2 without a flag day. The payload itself is easy to change later *because*
  the handshake exists.
- **Versioning over key/value.** Ignore-unknown key/value only buys *additive*
  growth; it can't reinterpret an existing field's format. A version field
  (also carried in each payload, so attachments are self-describing) lets a
  future v2 redefine the payload wholesale, safely. Keep v1 small and bet on
  never needing v2. A cheap escape valve short of a version bump: allow optional
  *trailing* fields within a version (consumers stop at the fields they know) —
  a supplement, not the strategy.
- **Wire-safe by construction.** It's a well-formed OSC sequence, so any
  terminal that doesn't recognize it skips it (gocui already does this via
  `stateOSCSkipUnknown` in `pkg/gocui/escape.go`; real terminals skip unknown
  OSC the same way). The metadata flowing through a real terminal must be
  harmless — this is a hard requirement, since pagers run outside lazygit too.

### 3.4 The OSC number — RESOLVED: `1717`

> **RESOLVED (audit done).** The final number is **`1717`**, replacing the `456`
> placeholder. There is no central registry; the convention is to pick a high,
> distinctive number and verify no real terminal *acts on* it (an unknown OSC is
> skipped by conformant terminals, but a *recognized* one can fire a visible
> side-effect — e.g. `OSC 555` flashes foot, `OSC 777` raises a desktop
> notification — which is exactly the wire-safety hazard the audit exists to
> avoid). Audited the live OSC allocations of **xterm, VTE, kitty, foot, WezTerm,
> iTerm2, Windows Terminal, Ghostty, VS Code, ConEmu, urxvt**; `1717` collides with
> none and sits in the large empty 1400–5000 band (only iTerm2's `1337` is nearby).
> The full danger list is in the spec appendix (`diff-line-metadata-osc-spec.md`).
> **The prototype code has been renamed `456`→`1717`** across all three repos
> (delta/difftastic/gocui/lazygit), and the handshake env var is now
> `EMIT_OSC1717_METADATA`; builds + the metadata unit tests are green in each.

Known-used slots to avoid (verified in the audit): `0–3` (title/icon/X11),
`4/5/6` (palette/special/tab color), `7` (cwd), `8` (hyperlinks), `9`
(notifications; `9;4` progress, `9;9` cwd), `10–19` (dynamic colors), `21/22`
(color query / pointer shape), `46/50/51/52` (logfile/font/Emacs/clipboard), `66`
(text sizing), `99` (kitty notifications), `104–106`+`110–119` (reset colors),
`133` (semantic prompt), `176` (foot app id), `555` (foot flash), `633` (VS Code),
`777` (rxvt notify), `1337` (iTerm2), `5522` (kitty clipboard), `30001/30101`
(kitty color stack).

### 3.5 Pagers to patch

Small universe, so standardization cost is low and we can supply the patches
ourselves rather than lobbying: **delta, diff-so-fancy, ydiff, difftastic,
diffr, riff, git-split-diffs**. delta and difftastic are the high-value targets
(difftastic because #1 categorically can't serve it). Patching git's own
`--color` is unnecessary (#1 covers it).

---

## 4. Consumer → field mapping

| consumer | `deleted` | `added` | `context` |
|---|---|---|---|
| staging / `enter` (find patch line) | `old-line` | `new-line` | `new-line` |
| edit `e` (editor target) | `new-line` | `new-line` | `new-line` |
| PR link `G` (anchor) | `L` + `old-line` | `R` + `new-line` | `R` + `new-line` |
| hunk-jump | `hunk-header` type (or coordinate discontinuity) | | |
| scroll-preserve | anchor on a change line via its native coord (`old`/`new`); `type` selects change lines | | |

So `old-line` is used **only** for finding the patch line of a deletion and for
its PR `L` anchor; `new-line` does everything else (editor for all types,
patch-find for added/context, PR `R` anchor).

---

## 5. #1 implementation sketch (the host-side fallback — do this first)

- Read the decolorized buffer (`View.Line(y)` upward from the target).
- Walk up to the nearest `@@` (new-start) and nearest `diff --git` (file);
  count `+`/space lines down to the target → `new-line`; first char → `type`.
- Reuse the `patch` package arithmetic; add multi-file (`diff --git`) splitting.
- Same `(file, type, new-line[, old-line])` result shape as the #2 payload, so
  the two are interchangeable behind one accessor and the call sites
  (`GetFileAndLineForClickedDiffLine` and friends) don't care which produced it.

---

## 6. Open questions

- ~~**`new-line` for a deleted line:**~~ **RESOLVED (#1 prototype).** The
  convention is exactly what `patch.LineNumberOfLine` already computes:
  `newStart` + #(added/context) above the deletion within the hunk — the
  new-file position the deletion sits at. Confirmed empirically: two consecutive
  deletions both report `new-line` = that shared position (e.g. both `2`), and
  are told apart only by `old-line` (`2` vs `3`). So a deleted line carries
  *both* numbers; consumers pick (`old-line` for staging-land and the PR `L`
  anchor, `new-line` for the editor). All pagers must agree on this for #2.
- ~~**Multi-file split approach:**~~ **RESOLVED (#1 prototype).** Split the
  buffer on `diff --git` boundaries to isolate the section for the file
  containing the target row, then `patch.Parse` that single section: its patch
  line indices line up 1:1 with the section's buffer lines, so the patch line
  index is just `targetBufferIdx − fileStartIdx`. New/old line numbers and the
  type then fall straight out of the patch arithmetic. The file path comes from
  the section's `+++ b/…` line (falling back to `--- a/…`, then `diff --git`).
- **Do headers carry line numbers?** `hunk-header` *could* carry old-start /
  new-start, but hunk boundaries are also derivable from coordinate
  discontinuities in the content lines, so it may be unnecessary. `file-header`
  needs no line numbers (the `file` field already attributes every line).
- ~~**difftastic specifics:**~~ **RESOLVED (#2 difftastic prototype, §10).** Two
  regions per row (one per side-by-side column), not N — token-level novelty is
  sub-cell colouring, not separate identity (§10.3). Each is emitted at the start
  of its column, before the line-number gutter. The prototype also surfaced a
  token-vs-line **model mismatch** (§10.2) the unified-diff pagers hid.
- ~~**v1 wire format:**~~ **RESOLVED (#2 prototype, §9).** Positional, `;`-delimited,
  `file` last so it may itself contain `;`; an absent `old-line` is the empty field.
- **Should the pager always emit, or only when the env var is set?** Leaning
  env-var-gated (zero cost when no consumer wants it; harmless outside a host).

---

## 7. Suggested build order

The prototype is a learning vehicle, not production code: its two jobs are to
**inform the final OSC spec** (which we want to publish for pager-developer
feedback) and to **inform a from-scratch production plan**.

> **Status:** steps 1–4 are done and step 5 is end-to-end verified for the
> NORMAL (unified, single-column) case — see §8 (#1) and §9 (#2). Side-by-side
> delta is prototyped (focused-main-view-notes.md §17) and **difftastic is
> prototyped in both modes (§10)** — so the emitter side now spans the full
> coverage table. What remains of step 5 is the *deliverables*: the **OSC spec
> draft is written** (`diff-line-metadata-osc-spec.md`, OSC number finalized to
> `1717` — §3.4 — and it speaks to difftastic's §10.2 model-mismatch finding), ready
> to circulate to pager developers; what's left is gathering that feedback and
> writing the production plan. Host *consumption* of side-by-side / difftastic
> output is still a separate, later step.

Sequence:

1. **#1 prototype first** (§5) — the buffer parser plus the
   `(file, type, new-line, old-line?)` accessor seam, wired to the
   focused-main-view consumers, verified across the structure-preserving
   renderings. Do this first because it:
   - validates the payload data model cheaply, with **no external deps** — if a
     field is wrong or missing you find out before writing the spec or a delta
     patch (the most direct "inform the spec" lever);
   - establishes the **shared accessor seam** that #2 plugs into as a second
     backend;
   - gives a **reference implementation** to validate delta's emitted metadata
     against later;
   - **ships independently** (the feature stops depending on delta-with-
     hyperlinks for the common cases).
2. **Pin the v1 wire format** — the concrete OSC bytes, payload encoding,
   version field, env-var name (§3.2–3.3), once #1 has confirmed the fields.
3. **#2 consumer side** — the gocui per-cell metadata mechanism + new-OSC parser
   behind the same accessor; testable against synthetic OSC output **before**
   delta exists (§3.1).
4. **#2 emitter side** — the delta patch; this is what stress-tests the spec
   against reality, side-by-side being the hard case (§3.5).
5. **End-to-end → finalize the spec (publish for feedback) → write the
   production plan.** End-to-end is **done** (§9.4); the **spec draft is written**
   (`diff-line-metadata-osc-spec.md`, OSC `1717`). Circulating it for feedback and
   the production plan are the remaining deliverables.

**Parallel de-risking (any time, doesn't block #1):** confirm by reading delta's
source that it can produce the fields per region — for `-` lines and in
side-by-side mode — at render time. It's the biggest *external* unknown; if delta
structurally can't emit something, the spec must adapt. Read-only research, so it
can run alongside #1. Still do #1 first: know what you *need* before checking what
delta *can do*.

---

## 8. #1 prototype — built & verified

Step 1 of the build order is done, at **prototype quality on the throwaway
branch** `use-delta-hyperlinks-for-clicking-in-diff`. What it comprises:

- **The data model is validated.** `types.DiffLineInfo {Path, Type
  (file-header | hunk-header | context | added | deleted | other), NewLine,
  OldLine}` — the same shape the #2 payload (§3.2) will carry. Building the
  consumers against it confirmed the field set is right: `deleted` genuinely
  needs *both* line numbers (see §6), `type` is load-bearing (the staging and PR
  consumers branch on `deleted`), and nothing else was missing.
- **The seam exists and has two real backends behind it.**
  `StagingHelper.GetDiffLineInfo` tries the buffer parser first, then the old
  `lazygit-edit://` hyperlink reader. This proves the seam is real (not a
  single-backend wrapper) and shows the degradation story: the hyperlink reader
  can't convey the side, so it reports `other`, which consumers treat as a
  non-deletion — exactly the pre-existing behavior. #2 slots in *ahead* of both.
- **The arithmetic is reused, plus old-file mirrors.** Added
  `Patch.OldLineNumberOfLine` / `PatchLineForOldLineNumber` (mirrors of the
  new-file functions) so a deletion lands on its exact patch line by old-file
  number — this is what fixes the two-deletions/two-additions case.
- **Layout-agnosticism is preserved by an integrity check, not gutter parsing.**
  The host does **not** learn delta's gutter format. Instead `Patch.IsWellFormed`
  checks each hunk's parsed body against its header-declared lengths; a gutter (or
  any body-restructuring) makes them disagree, so the parse is rejected and the
  seam falls back. This keeps "all layout knowledge stays in the pager" (§3.1)
  intact for #1 too, and is the cleanest signal for *when* to fall back.

**Verified** by feeding real pager output through gocui's actual escape parser
(`gocui.NewView` + `View.Write`) and running the parser on every resulting line —
i.e. against the bytes the live app would hold. Results drove the §2.3 coverage
corrections: ✅ no-pager, `git diff --color`, `delta --color-only` (no line
numbers), with consecutive deletions/additions landing distinctly; ❌ (clean
fall-through) `delta --color-only --line-numbers`, `diff-so-fancy`, delta default.

**Implications for #2 / the production plan:**

- The `IsWellFormed`-style "is this really a unified diff?" gate is worth keeping
  as the host's fallback trigger even once #2 exists (`#2 present → use it; else
  parse if well-formed; else give up`).
- delta-with-line-numbers is the *common* delta config (it's what the hyperlinks
  feature needs), and #1 can't serve it — so it's a strong motivator for #2, and
  a concrete case the delta patch must cover.
- The path comes from `+++ b/…`/`--- a/…`/`diff --git`; git's C-quoting of
  unusual paths is **not** handled (prototype). #2 should carry the path
  explicitly (it already does, §3.2) and the production #1 path should decode
  quoting.
- **`View.BufferLineForViewLine` (and the `HyperLinkInLine` /
  `DiffLineMetadataInLine` readers that mirror it) could return *stale* data, not
  only panic — FIXED in the prototype (session 4).** The old guard only rejected a
  `linesY` that had gone out of range of a shrunk `v.lines`. But if a re-render
  produced *fewer* view lines than the previous one, `refreshViewLinesIfNeeded`
  (which overwrites `viewLines` in place without truncating, to keep the tail
  visible for flicker-avoidance) left stale entries from the previous render in the
  tail; with **wrapping**, such an entry's `linesY` could still be *in range* of
  the new, shorter, less-wrapped buffer, so the guard passed and a view line that
  no longer existed mapped onto the wrong buffer line. This is a *single-threaded,
  deterministic* defect, not just a concurrency one (the non-truncating refresh is
  enough; see the unit test `TestBufferLineForViewLineStaleTail`).
  - **The fix (reworked):** a first cut added `freshViewLineCount` (bound the readers
    on the count of view lines built from the current buffer). That was then
    **reverted** in favour of the **off-screen render** (focused-main-view-notes.md
    §13.5): a re-render builds into a second `viewBuffer` and swaps it in wholesale,
    so `refreshViewLinesIfNeeded` simply **truncates** `viewLines` to the current
    buffer and no stale tail ever forms. The displayed buffer is always a complete
    render, so the three readers (unified onto one `bufferLineForViewLine` helper)
    are consistent by construction. `TestBufferLineForViewLineStaleTail` now guards
    the truncation.
  - **Still a part-3 constraint, NOT yet fixed:** mechanism #1's host-side parse
    (`diffLineInfoFromBuffer`) maps the view line and reads the buffer text in *two
    separate* locked gocui calls (`BufferLineForViewLine` then `BufferLines`), so a
    re-render between them could desync the index from the text. Harmless for the
    forward consumers (the focused main view is static while the user clicks/presses
    enter) and harmless for the planned inverse **predicate scan** *iff* that scan
    runs inside the render task's own goroutine (no concurrent writer, since tasks
    don't overlap). So the safe rule for the part-3 reader is: scan on the task
    goroutine, or take a single buffer+index snapshot under one lock. The metadata
    (#2) and hyperlink readers already map-and-read under one lock, so they are
    atomic today.
- Still open (untouched by #1): the §3.x #2 wire-format questions, headers
  carrying line numbers, difftastic specifics.

---

## 9. #2 prototype — built & verified end-to-end

Build-order steps 2–5 are done at **prototype quality**, for the **NORMAL
(unified, single-column) case only** (side-by-side deliberately ignored this
iteration), end-to-end verified in the running app. The delta patch lives on
branch `prototype-osc-metadata` in `/Users/stk/Stk/Dev/Builds/delta`; the host
side is on `use-delta-hyperlinks-for-clicking-in-diff`. §9.1 records the delta
de-risk, §9.2 the pinned wire format, §9.3 the deferred items, and §9.4 what was
built and how it was verified.

### 9.1 Delta de-risk (read-only) — all four fields are reachable, with one gotcha

Where delta renders each unified-mode content line, and what it has there:

- **The single per-line emit point is `Painter::paint_lines`** (`src/paint.rs`),
  whose loop emits every content line (`output_buffer.push_str(&line)`). All three
  unified content paths funnel through it: context via `paint_zero_line`, and
  buffered minus/plus via `paint_buffered_minus_and_plus_lines` →
  `paint_minus_and_plus_lines` (the non-`side_by_side` branch). **Side-by-side
  calls `Painter::paint_line` (singular) directly and never goes through
  `paint_lines`**, so threading metadata through `paint_lines` leaves the
  out-of-scope side-by-side path untouched.
- **`type`** is the `State` enum variant at the emit point (`HunkMinus` / `HunkPlus`
  / `HunkZero`, plus `…Wrapped`). Reliable.
- **`file`** is on the `StateMachine` as `plus_file` / `minus_file` (parsed from
  `+++ b/…` / `--- a/…`). The hunk-header handler already uses the exact selection
  we want: `if plus_file == "/dev/null" { minus_file } else { plus_file }`.
- **`new-line` / `old-line`.** ⚠️ **The gotcha that shaped the design:** delta only
  maintains its old/new line counters (`LineNumbersData.line_number`) **when
  `--line-numbers` is enabled** — in delta's *default* mode (the #2 target!)
  `Painter.line_numbers_data` is `None` and the counters never advance. So the
  numbers are **not** sitting there for free; the patch must track its own
  counters, seeded from the parsed hunk header (`@@ -old,len +new,len @@`, already
  available as `ParsedHunkHeader.line_numbers_and_hunk_lengths`).
- **Deleted lines carry both numbers** exactly as §6 resolved: at a `-` line the
  new-file counter has *not* advanced past the preceding context/added lines, so it
  already holds `newStart + #(added/context above)` — the new-file position the
  deletion sits at. Mirrors `patch.LineNumberOfLine`.

**Why not reuse delta's `LineNumbersData` (it already has the counters + `plus_file`
and is already threaded to the emit point)?** Because forcing it `Some` in default
mode to get the counters would also (a) render the gutter and (b) change
wrap-width math (`wrapping.rs` reserves gutter width whenever it's `Some`). Both are
layout changes. A **dedicated, purely-additive emitter** (only injects OSC bytes,
never touches styling/width/wrapping) is both safer and cleaner, and its counter
logic is a near-copy of `linenumbers_and_styles`.

### 9.2 Pinned v1 wire format (final OSC `1717`; prototype code still uses `456`)

> The number was finalized to **`1717`** after the terminal audit (§3.4); the
> prototype code has been renamed from the `456` placeholder to emit `1717`. The
> published spec (`diff-line-metadata-osc-spec.md`) uses `1717`.

```
ESC ] 1717 ; <version> ; <type> ; <new-line> ; <old-line> ; <file> ST
```

- `ESC` = `0x1b`; `ST` = `ESC \` (`0x1b 0x5c`) — same framing as delta's OSC-8.
- `<version>`: decimal; `1` for v1.
- `<type>`: one char — `c` context · `a` added · `d` deleted. (Reserved for later:
  `h` hunk-header, `f` file-header, `o` other. The prototype emits only `c`/`a`/`d`,
  i.e. content lines; header rows get no attachment and the host falls back.)
- `<new-line>`: decimal new-file line number; present on every content line.
- `<old-line>`: decimal old-file line number; **empty** unless `type=d`.
- `<file>`: repo-relative or absolute path; **last field on purpose** so it may
  contain `;`. The host splits the payload into at most 5 fields and the path is
  the remainder. Host normalizes via `RepoPaths.WorktreePath()` (not `RepoPath()`).

**Env handshake (the one piece that can't be retrofitted):**
`EMIT_OSC456_METADATA=V1[,V2,…]` — the host advertises the versions it understands;
the pager emits the highest mutually-understood one. Unset (outside lazygit) ⇒
pager emits nothing ⇒ harmless in a raw terminal / `less` / `tmux`. The prototype
emits V1 when the advertised list contains `V1`.

### 9.3 Deferred / known prototype limitations

- **Terminal-source audit of the OSC number is DONE** (§3.4). Final number is
  **`1717`** (audited against xterm/VTE/kitty/foot/WezTerm/iTerm2/Windows Terminal/
  Ghostty/VS Code/ConEmu/urxvt; danger list in the spec appendix). The prototype
  code has been renamed `456`→`1717` across delta/difftastic/gocui/lazygit + the
  `EMIT_OSC1717_METADATA` env var; builds and metadata unit tests green in each.
- **Wrapped continuation rows** (`Hunk*Wrapped`) get no attachment in the prototype
  — only the primary content row does. Fine for the normal case (gocui's own
  wrapping is handled host-side by the view-line→buffer-line mapping); delta-level
  wrapping of one logical line into several rows is the unhandled case. **This was
  confirmed a real bug** (via the difftastic prototype, §10.8): when the *pager*
  wraps, each row is a distinct host buffer line and needs its own record, so `e`/
  `enter`/hunk-nav break on continuation rows. **Now FIXED in delta too** (§10.8):
  delta wraps only in side-by-side mode, and each wrapped row now re-emits its
  primary line's record (no counter advance). difftastic was fixed the same way.
- **Header rows** (`@@`, `diff --git`, `---`/`+++`) get no attachment; acting on a
  header row falls through to #1, then to no-selection.

### 9.4 What was built, and how it was verified

**Emitter (delta, branch `prototype-osc-metadata`).** A dedicated additive
emitter (`src/features/diff_line_metadata.rs`) gated on `EMIT_OSC456_METADATA`:
it tracks its own old/new counters (seeded at each hunk header, mirroring
`LineNumbersData::initialize_hunk`, so it works in delta's default mode where the
line-number counters are otherwise dormant) and is threaded through
`Painter::paint_lines` next to `line_numbers_data`, prepending one OSC per
content line. It only injects bytes — no styling/width/wrapping change — so with
the var unset the output is byte-for-byte identical to stock delta (confirmed
with `cat -v`).

**Carrier (gocui).** The escape interpreter (`pkg/gocui/escape.go`) now
accumulates the OSC number before dispatching (so multi-digit `456` is
recognized alongside `8`), and a new `stateOSCMetadata` collects the payload. The
payload is stamped onto each cell like a hyperlink (`pkg/gocui/view.go`) and read
back via `DiffLineMetadataInLine`; it is cleared at each line boundary so it
can't bleed onto a following line with no metadata (the pager never closes it).

**Consumer (lazygit).** `parseDiffLineMetadata` (in `diff_line_parser.go`,
alongside the #1 buffer parser, both producing `parsedDiffLine`) parses the
payload; `StagingHelper.diffLineInfoFromMetadata` normalizes the path and it is
tried **first** in `GetDiffLineInfo`, ahead of #1 and the hyperlink. The host
advertises `EMIT_OSC456_METADATA=V1` on the pager subprocess in `newPtyTask`.

**Verified.**
- **Delta bytes** (real binary, `cat -v`): correct OSC per content line across
  context/added/deleted, **two consecutive deletions sharing a new-line and
  differing only in old-line** (the §6 case), multi-hunk re-seeding, and a
  whole-file deletion (`/dev/null` → `minus_file`, new-line 0). Output byte-identical
  to stock delta when the var is unset. All 71 touched delta unit tests pass.
- **Carrier + parse** (unit tests): `TestDiffLineMetadata` drives synthetic OSC
  through the real `View.Write` and asserts the per-line payloads, no header
  bleed, and that the OSC bytes don't render; `TestParseDiffLineMetadata` covers
  the payload parsing incl. semicolon-in-path, absolute path, and malformed
  rejections.
- **Real binary → gocui → parse** (throwaway test, since it needs the local
  delta build): ran the patched delta in **default mode** with lazygit's env on a
  real diff, fed the output through a real `View`, and recovered the correct
  `(file, type, new, old)` for every content row including both deletions.
- **In the running app** (`just debug`, manual): with delta's **default mode**
  (which #1 cannot parse and which emits no hyperlinks), clicking / `enter` / `e`
  / `G` resolve correctly via #2 — including the side for deletions, which #1 +
  hyperlinks cannot convey.

**Next (deliverables, not yet done):** finalize and publish the OSC spec for
pager developers (after the §3.4 terminal audit picks the real OSC number), and
write the from-scratch production plan. Then extend coverage to side-by-side and
difftastic (multiple regions per row, §3.1), which this iteration ignored.

> **Update (difftastic prototype, 2026-06-10):** difftastic is now prototyped too
> — see §10. It is the **categorical** #2-only case (no `@@`/`+`/`-` to parse in
> *either* of its modes), and it stress-tested the format against a model #2 was
> never built around. Headline: v1 holds, but difftastic surfaces a real
> **token-vs-line model mismatch** the unified-diff pagers hid (§10.2).

---

## 10. difftastic emitter prototype — built & verified

§15 step 3's sibling: extend the #2 emitter to **difftastic**, the one pager #1
*categorically* cannot serve (it restructures the diff in every mode). Done in
the difftastic repo on branch `prototype-osc-metadata` (two commits: side-by-side,
then inline), parallel to and independent of lazygit, mirroring the delta work
(§9). It emits the **same v1 wire format** (§9.2) under the **same
`EMIT_OSC456_METADATA` handshake**, so one host reader consumes either pager.

### 10.1 What was built, and why it was *simpler* than delta

Both of difftastic's text modes now emit per-cell metadata:

- **Side-by-side (default/signature mode)** — `src/display/side_by_side.rs`. The
  left (old-file) column carries the deleted/old line's record, the right
  (new-file) column the added/new line's; a modification — one aligned row, old
  left / new right — emits **two records** (`d` left, `a` right), exactly as a
  unified diff splits a change into a `-` and a `+`. Context lines (shown on both
  sides) emit the same `c` before each half. **Every visual row of a
  difftastic-wrapped line carries its record** (not just the first — §10.8); only
  a side's padding rows once its wrapped content is exhausted, and the blank
  counterpart half of a pure add/delete, carry none. The whole-file single-column
  add/delete path is covered too (`d` with new-line **0** for a deleted file,
  matching delta's `@@ -1,N +0,0 @@`).
- **Inline mode** — `src/display/inline.rs` (§10.5).
- **Shared emitter** — `src/display/diff_line_metadata.rs`: `negotiated_version()`
  (a verbatim port of delta's handshake) plus pure `left_cell`/`right_cell`/
  `single_column_cell` formatters. Unit-tested (7 cases).

**It was markedly less code than delta**, and that itself is a finding about how
well-matched the format is to a structural tool:

- **difftastic already has the line numbers natively.** It *always* renders
  old/new gutters, so each row arrives as an aligned
  `(Option<LineNumber>, Option<LineNumber>)`. delta had to track its *own* old/new
  counters seeded from each `@@` header (delta's counters are dormant unless
  `--line-numbers`); difftastic needs **no counter tracking at all** — each
  record falls straight out of the row's two line numbers and its novelty.
- **The file path is a parameter** (`display_path`), not parsed from a `+++`
  header, so the emitter is created once per file with the path fixed — no
  hunk-header plumbing.
- **Purely additive**, like delta: only injects OSC bytes. With the var unset the
  output is **byte-for-byte identical** to stock difftastic (verified by stripping
  the OSC456 sequences and `cmp` across side-by-side, inline, and whole-file
  cases; all 127 unit + 23 integration tests green).

### 10.2 THE headline finding: a token-vs-line model mismatch v1 can't fully express

The unified-diff pagers (git, delta, diff-so-fancy) all derive from git's
**line-granular** patch, where a modified line is *by construction* a `-` line
plus a `+` line. The c/a/d type set was shaped by that model and fits it exactly.
difftastic is **token-granular**: it aligns an old line with a new line and marks
novelty *per token*. That produces aligned rows the line model has no clean slot
for. Concretely (real output, `let x=1; println!("{}", x);` → `let x=2; let y=3;
println!("{}", x + y);`):

```
{OSC ...;c;4;;src/lib.rs}3     println!("{}", x);   {OSC ...;a;4;;src/lib.rs}4     println!("{}", x + y);
```

The old line `println!("{}", x);` has **no novel tokens** (all of them survive
into the new line; difftastic colours it as plain context, not novel-red),
while the new line added `+ y`. So difftastic's faithful per-cell verdict is
**`c` on the left, `a` on the right** — the *same aligned row* carries a context
record and an addition record. There is no `d` for the old side, because by
difftastic's model nothing was deleted.

Consequence for a host mapping cells to git's patch (lazygit stages against the
real `git diff`, the same change being `-println!("{}", x);` / `+…x + y);`):

- the **right/added** cell resolves correctly (`a; new=4` → git's `+` line);
- the **left/old** cell resolves as **context at the *new* line 4**, *not* git's
  `-` line for old line 3. Its old-file deletion identity is not recoverable from
  the record.

Practical impact is small — users click the changed (green) side, and `e` (edit)
on the left still opens new-file line 4 — but it is a genuine semantic gap, and
exactly the kind of thing the prototype exists to find. **The faithful emission
was chosen deliberately** (it surfaces the mismatch rather than hiding it).

> **Note (session 9): the one user-visible bite is cross-pager selection
> preservation (#7).** Within difftastic this is invisible — it *renders* the old
> line as dim context, matching the `c` record. But it only happens in difftastic's
> **AST mode**; in its line/Text fallback (e.g. our test file with `let x=1;` in a
> `.go`, which fails the Go parse) difftastic diffs by line and emits `d`/`a` like a
> unified diff. The bite: select the old line in delta (where it is a red `d`), then
> `|` to difftastic AST mode, where the same patch line is now `c` — the identity
> restore (#7) can't match `d` against `c`, so the selection is dropped. Minor in
> practice; a host concern, not a pager-spec one (kept out of the spec).

Options to flag for the spec / production, none taken now:
- **Host-side:** treat an old-column cell that is aligned with a novel new line as
  the `-` side of a modification (the host knows it's the left column).
- **Emitter-side:** difftastic could emit `d` for the old side of *any* aligned
  changed row (v1-compatible — no format change). But that re-imposes the
  line-granular model difftastic exists to escape, and discards its more precise
  "this content was not removed" judgement. Probably wrong to force.
- **A `modified`/`m` type** (v2) that means "aligned, changed, present on both
  sides" would name the case directly — but it splits the clean c/a/d staging
  mapping (§4) and every pager would have to agree when to use it. Not obviously
  worth it; record as a v2 candidate, not a v1 gap.

### 10.3 "N regions per row" was an over-estimate — it's 2, same as delta SxS

§3.1 worried difftastic would need the per-cell carrier to hold **N** attachments
per row. It doesn't. A side-by-side visual row is exactly **two columns → at most
two records** (one per column), identical to delta's side-by-side (§17). What is
genuinely N-ary in difftastic — token-level novelty — is **sub-cell colouring, not
separate identity regions**: a row's left cell is one patch line however many
tokens are highlighted within it. So the existing per-cell mechanism (§3.1) and
the host's row+column→identity model (§17.4) already suffice; **no N-region
machinery is needed.** (Resolves the §6 "difftastic specifics: how many regions
per row" open question: two.)

### 10.4 §17.3 amplified — context/added records carrying no old-line bites harder here

The delta-SxS finding (§17.3: `c`/`a` carry no `old-line`) is **worse** for
difftastic, and unavoidable rather than latent:

- difftastic's left column **always** shows real old-file line numbers (it's the
  primary mode, not an opt-in), so the temptation/need to read an old number off a
  left cell is constant.
- old ≠ new is the **norm**, and the offset is **not constant** — it follows
  difftastic's structural alignment, not a per-hunk delta. In the example above the
  left context cell shows old line **3** but its record is `c; new=4`; old=3 is
  **not derivable** from new=4 (with delta it would be new minus a per-hunk
  constant). So for difftastic, "carry both numbers on every record" (the v2 move
  §17.3 floated) is the *only* way to make a left-column old number available.

Still: nothing in the §16 host consumers needs the old number today (they key on
`type`/`file`/new-line), so v1 stays as is — but difftastic is the strongest
argument yet for the v2 "both numbers always".

### 10.5 The one synthesized field, and the inline grouped-layout demo

- **A pure deletion's `new-line` is the only derived field.** Having no linear
  new-file counter, difftastic computes it from the *previous aligned new line*
  (`prev_rhs + 1`), mirroring delta's "a deletion sits at the following new line".
  Verified correct for the common case; it can drift across hunk boundaries or with
  `num_context_lines = 0` (the previous new line is then far away). Documented, not
  a blocker — the deletion's *old*-line (its real identity for staging) is always
  exact; only the editor-target new-line is approximate.
  - **Session 9: re-verified empirically and dropped from the published spec as too
    marginal.** A single mid-file pure deletion emits `d;4;4` (exact) at default
    context but `d;1;4` (new-line drifted low by the elided count) at `--context 0`;
    `old-line` is exact in both. Since lazygit renders at default context, the drift
    isn't reachable in normal use, and even at zero context the only effect is `e` on
    a *deleted* line opening the new file a few lines off (inherently approximate).
    So it's no longer a spec §8 item.
- **Inline mode proves the metadata's worth beyond layout.** Inline groups **all
  old-side lines, then all new-side lines** (not interleaved like git). A
  modification's `d` (deletions group) and `a` (additions group) are therefore
  emitted *far apart* in the stream — yet both reference the same new-file line and
  the `d` carries the old line, so the host reconstructs each line's identity from
  the numbers alone. A positional or structural re-parse of inline output could not
  pair them. (The inline type comes from the novel-line sets, not the all-or-nothing
  novelty difftastic uses for *colouring*, so a context line inside a hunk is tagged
  `c`.)

### 10.6 Verification

- **Bytes** (`cat -v` on the real debug binary): correct per-cell records across
  context / modification / pure addition / pure deletion in **both** modes; the
  asymmetric `c`+`a` row above; wrapped lines (**every** visual row tagged, the
  exhausted side's padding rows excepted — §10.8);
  whole-file add (`a`) and delete (`d; new=0`). Byte-identical to stock difftastic
  when the var is unset (multiple `cmp` checks).
- **Realistic invocation**: driven as `GIT_EXTERNAL_DIFF` on a Rust file (how
  lazygit would invoke it) — the `file` field is git's repo-relative path
  (`src/lib.rs`), which the host normalises via `WorktreePath()` as for delta.
- **Tests**: 7 new emitter unit tests + the existing 30 display + full suite (127
  unit, 23 integration) all green; `cargo fmt` clean.

### 10.7 Host consumption is still a separate, later step (as for delta SxS)

Not done here, same as §17.4. When wired up, difftastic is exactly the
**row+column→identity** case (§17.4): a single visual row carries up to two records
keyed by column. Its max-two-per-row shape is identical to delta side-by-side, so
the same host work covers both. The §10.2 model mismatch is a host-design input,
not a carrier/parser change.

### 10.8 Pager-level wrapping must tag *every* row — a spec correction (and a delta bug)

Found by the user testing the prototype in lazygit. The first cut followed the
delta convention "wrapped continuation rows carry no attachment" (§3.1, §9.3,
§17.1). **That convention is wrong whenever the *pager itself* wraps a long line**,
and it produced two concrete bugs in lazygit on difftastic side-by-side output:
pressing `e`/`enter` on a continuation row did nothing (no metadata to resolve),
and hunk navigation (`<right>`) stopped on *every* row because the un-tagged
continuation rows broke each wrapped change line into one block per visual row.

The root distinction the original convention missed:

- **Terminal/host wrapping** — the pager emits *one* line and the terminal (or
  gocui) wraps it into several *view* lines of *one buffer line*. Here only the
  primary row needs metadata; the host's view-line→buffer-line mapping already
  routes every view line of that buffer line to it. This is the case §3.1 had in
  mind, and it's correct *for that case*.
- **Pager wrapping** — the pager emits *several* lines (several `\n`s) for one
  logical line, as difftastic's side-by-side does (and delta does with
  `wrap-max-lines`). Now each wrapped row is a **distinct buffer line** to the
  host; there is nothing tying row N+1 back to row N, so each must carry its own
  metadata or it has none.

**Fix (difftastic):** every visual row of a wrapped line now carries the same
record (`amend!` into the side-by-side commit). A side that has run out of
wrapped content carries none on its padding rows (no content there to identify);
the still-wrapping side carries the record on each. Verified: a 6-row wrapped
modification tags all six (`d`/`a`), and an uneven wrap tags only the side still
producing content. So `e`/`enter` resolve on any row and nav treats the wrapped
line as one block.

**Spec consequence:** state the rule positively — *the pager emits the line's
record at the start of **every output row** it produces for that line, including
its own wrapped continuations.* The host attaches per buffer line, so it just
works; pagers that rely on terminal wrapping emit one row and are unaffected.

**Delta had the same bug — now FIXED** (§9.3). Delta wraps **only in side-by-side
mode** (`wrap_minusplus_block`/`wrap_zero_block` are called nowhere else; unified
mode truncates instead), so the bug was SxS-only, but it was real there. The fix
is the same idea adapted to delta's counter-based emitter: a wrapped continuation
row (`HunkZeroWrapped`/`HunkMinusWrapped`/`HunkPlusWrapped`) **re-emits the record
of the primary line it continues, without advancing the counters** — so the next
line's numbers stay correct (verified: a context line after a 5-row wrapped line
still reports the right new-line). `osc_for_line` is the single chokepoint, so the
one change covers both SxS emit paths (the minus/plus precompute and the
`paint_zero_lines_side_by_side` context path). Landed as an `amend!` into the
delta side-by-side commit, with a unit test (`test_wrapped_rows_reemit_…`).

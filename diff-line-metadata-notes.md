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
| difftastic / side-by-side | ❌ | ✅ if patched |
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

### 3.4 The OSC number

No central registry, but allocations are rare; pick a **high, distinctive**
number (the iTerm2 `1337` convention) and **verify against the sources of
xterm, VTE, kitty, foot, WezTerm, iTerm2, Windows Terminal** before committing.
`456` is a placeholder. Known-used slots to avoid (from memory, against a
knowledge cutoff — re-check): `0/1/2` (title/icon), `4`+`104` (palette),
`7` (cwd), `8` (hyperlinks), `9` (notifications; `9;4` progress), `10–12`+
`110–119` (colors), `50` (font), `52` (clipboard), `99` (kitty notifications),
`133` (semantic prompt), `633` (VS Code shell integration), `777` (urxvt),
`1337` (iTerm2).

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
- **difftastic specifics:** how many regions per row in practice, and where
  exactly it should emit each attachment.
- **v1 wire format:** delimiter and encoding of an absent `old-line` (positional
  vs. a minimal key/value within the version — note this is a *within-version*
  choice, orthogonal to the versioning-vs-key/value extensibility decision).
- **Should the pager always emit, or only when the env var is set?** Leaning
  env-var-gated (zero cost when no consumer wants it; harmless outside a host).

---

## 7. Suggested build order

The prototype is a learning vehicle, not production code: its two jobs are to
**inform the final OSC spec** (which we want to publish for pager-developer
feedback) and to **inform a from-scratch production plan**. Sequence:

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
   production plan.**

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
- **`View.BufferLineForViewLine` (and `HyperLinkInLine`, which it mirrors) can
  return *stale* data, not only panic — fix when productionizing.** The guard
  only rejects a `linesY` that's gone out of range of a shrunk `v.lines`. But if
  a newer render task has *refilled* `v.lines` while `v.viewLines` still holds
  the previous task's entries (kept deliberately, to avoid flicker), a `linesY`
  that is still in range now points at the *new* buffer's content at that index —
  so we silently map to the wrong line, no panic. Harmless in the prototype (the
  focused main view is static while the user clicks/presses enter), but the
  productionized read needs the view-line→buffer-line mapping to be consistent
  with the buffer it then reads: e.g. snapshot `viewLines` *and* `lines` under one
  lock and parse from that snapshot, or tie the read to the task that produced the
  buffer. Both #1 (this mapping) and the #2 per-cell attachment read have to get
  this right.
- Still open (untouched by #1): the §3.x #2 wire-format questions, headers
  carrying line numbers, difftastic specifics.

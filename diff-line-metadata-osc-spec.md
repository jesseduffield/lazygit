# Diff Line Metadata over OSC 1717 — draft specification (v1)

**Status: draft, for feedback.** This document describes a small terminal
escape-sequence protocol by which a diff pager (delta, difftastic, diff-so-fancy,
…) annotates each rendered line of a diff with the patch-space identity it
represents, so that a host program rendering the pager's output can map a screen
row (and column) back to *the exact line in the underlying diff*.

It is published to gather feedback from pager authors before anything is
finalized. The wire format, the negotiation handshake, and the OSC number are all
open to revision — §9 lists the points where feedback is most wanted.

The protocol grew out of [lazygit](https://github.com/jesseduffield/lazygit), but
nothing in it is lazygit-specific; "the host" below means any program that runs a
pager and consumes its output.

---

## 1. Motivation — what this enables, and why parsing isn't enough

A host that shows a diff rendered by a pager often wants to act on the line the
user is pointing at:

- **dive into staging / patch-building** for that hunk line,
- **open an editor** at that file and line,
- **open that line in a code-review / PR web view** (needs the side — old vs new),
- **navigate by hunk or by file** within the rendered diff,
- **preserve the scroll position and selection** across a re-render (the diff is
  re-rendered with a different context size, or a different pager, and the host
  wants to keep the user anchored on the same patch line).

Every one of these needs the same primitive: **given a rendered row, recover
`(file, side, line)`** — the precise line of the unified diff that row stands for.

For *structure-preserving* renderings (no pager, `git diff --color`, or
`delta --color-only` without line numbers) the host can recover this by parsing
the on-screen text: walk up to the nearest `@@` and `diff --git`, count `+`/` `
lines, read the leading `+`/`-`/space. That works and needs no cooperation from
the pager.

But the moment a pager **restructures** the diff, the unified-diff structure the
parse relies on is gone:

- `delta` in its default mode conveys the side (added/removed) purely by
  background color and prefixes a line-number gutter, so there is no `+`/`-`
  marker to read and no parseable `@@` body;
- `delta --color-only --line-numbers` (the configuration people actually use)
  pushes the `+`/`-` marker off the start of the line behind the gutter, so a
  naive parse is *confidently wrong*;
- `diff-so-fancy` rewrites the headers and strips the `+`/`-` markers entirely;
- `difftastic` is token-granular and side-by-side — there is no unified-diff line
  structure left in *either* of its modes.

In all of these the **pager is the only component that still knows** which file,
side, and line each rendered cell belongs to — it computed exactly that to render
the diff. This protocol asks the pager to *state* that knowledge inline, in a form
the host can read back and that is harmless everywhere else.

---

## 2. Design at a glance

1. The pager emits one OSC sequence carrying
   `(version, type, new-line, old-line, file)` **immediately before** each
   rendered region (a region is "one source line's worth of content in one
   column" — see §6).
2. The host attaches each record to **the cell that follows it**, exactly the way
   OSC-8 hyperlinks attach to the cells they precede. This makes the metadata
   survive terminal wrapping and multi-column layouts **without the host ever
   reasoning about layout** — it just reads the nearest preceding attachment.
3. The whole thing is gated behind an **environment-variable handshake**, so a
   pager run outside a participating host (in a raw terminal, `less`, `tmux`, a
   CI log) emits nothing and behaves byte-for-byte as before.

The protocol is **layout-agnostic on the host side by construction**: all layout
knowledge (where a column starts, where a gutter ends, how a long line wraps)
stays in the pager, which is the only component that has it.

---

## 3. Negotiation handshake

```
EMIT_OSC1717_METADATA = V1[,V2,…]
```

- The **host** sets this environment variable on the pager subprocess to the list
  of protocol versions it understands, highest-preferred first is *not* required —
  the list is a set.
- The **pager** emits the **highest version present in both** its own supported
  set and the advertised set. If the variable is unset, empty, or shares no
  version with the pager, the pager **emits nothing** and its output is unchanged.

Why a handshake, and why it must exist in v1 even though v1's payload is tiny:

- **It is the one piece that cannot be retrofitted.** Without negotiation you
  cannot introduce a v2 later without a flag day. The payload itself is easy to
  evolve *because* the handshake exists.
- **It guarantees zero cost when unwanted.** Outside a participating host the
  variable is unset, so there is no output change to audit, no risk in a raw
  terminal, no interference with `less`/`tmux`/pipelines.

The variable name and the value grammar are themselves open to feedback (§9).

---

## 4. Wire format (v1)

### 4.1 The sequence

```
ESC ] 1717 ; <version> ; <type> ; <new-line> ; <old-line> ; <file> ST
```

- `ESC` is `0x1B`; `ST` (String Terminator) is `ESC \` = `0x1B 0x5C`. A `BEL`
  (`0x07`) terminator is also accepted, but `ST` is preferred — this matches the
  framing of OSC-8 hyperlinks.
- The payload is **positional and `;`-delimited**. There are exactly five fields.
- `<file>` is **last on purpose** so that it may itself contain `;`: the host
  splits the payload into at most five fields and treats everything after the
  fourth `;` as the path.

Raw bytes, for a context line at new-file line 10 of `src/foo.go`:

```
\x1b]1717;1;c;10;;src/foo.go\x1b\
```

### 4.2 Fields

| field | presence | meaning |
|---|---|---|
| `version` | always | decimal protocol version; `1` for v1. Carried in every record so attachments are self-describing. |
| `type` | always | one character — see §5.1. v1 emits `c` (context), `a` (added), `d` (deleted). |
| `new-line` | content lines | new-file line number, in the **diff's new-file space** (see §5.2). |
| `old-line` | only `type=d` | old-file line number. **Empty** for `c` and `a`. |
| `file` | always | the file path the line belongs to; absolute or repo-root-relative (the host normalizes — emit whichever is convenient). Carried on **every** record so a single record is a complete answer. |

### 4.3 Examples

| rendered line | emitted record (`;`-form) |
|---|---|
| context, new line 10 | `1717;1;c;10;;src/foo.go` |
| addition, new line 11 | `1717;1;a;11;;src/foo.go` |
| deletion, old line 9, sits at new pos 11 | `1717;1;d;11;9;src/foo.go` |
| two consecutive deletions | `…;d;11;9;…` then `…;d;11;10;…` (same `new-line`, different `old-line` — see §5.3) |
| whole-file deletion | `1717;1;d;0;9;old/path` (`new-line` 0 — see §5.5) |

---

## 5. Semantics

### 5.1 Type

`type` is one character. v1 defines and emits three (the content-line types):

- `c` — context (unchanged) line
- `a` — added line
- `d` — deleted line

Reserved for future use (v1 emitters do **not** emit these; header rows simply
carry no record and the host falls back):

- `h` — hunk header (`@@ … @@`)
- `f` — file header (`diff --git`, `---`, `+++`)
- `o` — other

**`type` is load-bearing and cannot be inferred from the other fields.** Under the
coordinate rules, `added` and `context` both carry `{new-line present, old-line
absent}`, so presence alone cannot distinguish them — yet the host must (scroll
preservation anchors specifically on *change* lines, so it has to tell `added`
from `context`). Hence an explicit type rather than an inferred one.

### 5.2 Line-number spaces

- `new-line` is the line number in the **diff's new-file space** — i.e. the
  new-file line numbering the diff itself uses, *not* necessarily the working-tree
  file (the diff may be against staged content). A host that opens an editor is
  expected to re-map this through its own diff↔worktree adjustment; the pager
  should emit the number as it appears in the diff it is rendering.
- `old-line` is the old-file line number, present **only** for deletions.

### 5.3 The deleted-line convention (both numbers)

A `d` record carries **both** numbers:

- `old-line` is the deletion's own old-file line number.
- `new-line` is the new-file position the deletion *sits at*: `newStart` plus the
  number of added/context lines above it within the hunk. This is exactly what
  `git`'s patch arithmetic already computes for a removed line.

Consequence, which all pagers must implement identically: **two consecutive
deletions share the same `new-line`** (nothing new-file-side advances between
them) and are told apart only by `old-line`. Example — two deletions at old lines
9 and 10, both sitting at new position 11:

```
1717;1;d;11;9;src/foo.go
1717;1;d;11;10;src/foo.go
```

A host uses `old-line` to find a deletion's patch line and its old-side
(review/PR `L`) anchor, and `new-line` for the editor target.

### 5.4 Why the side must be carried at all

It is tempting to think the side (added vs deleted) is redundant with the rendered
glyph. It is not: in `delta`'s default rendering there are **no `+`/`-` glyphs** —
the side is conveyed purely by background color, which a host reading decolorized
text cannot recover. The pager therefore has to state the side explicitly, which
is what the `a`/`d` type does.

### 5.5 Whole-file add / delete

A deleted file's lines carry `new-line` = `0` (mirroring git's `@@ -1,N +0,0 @@`);
an added file's lines carry the new-file numbers normally and `type=a`.

---

## 6. Emit rules (placement)

### 6.1 One record per region, at the region's start

The pager emits each region's record at the **start of that region**. Everything
from there until the next region's record (or end of line) belongs to that
record — *including* any line-number gutter or other embellishment the pager
considers part of the region. Where a region "really" starts is the pager's
call, not the host's. The single firm requirement:

> **The record must precede the region's first cell**, so that a host searching
> leftward from any cell lands in the correct region.

### 6.2 Multiple regions per row (side-by-side)

A side-by-side row shows two source lines at once, so it carries **two records**,
one before each column:

- **left column → the old-side line** (`d`, carries `new-line` + `old-line`),
- **right column → the new-side line** (`a`, carries `new-line`),
- **a context line, shown in both columns → the same `c` record before each
  half.**

The blank counterpart of a pure add/delete (the empty half) carries **no** record.

A side-by-side row carries **at most two** records — one per column. (An earlier
draft worried token-level tools might need *N* records per row; they do not.
Token-level novelty is *sub-cell coloring*, not separate identity: a column's cell
is one source line however many tokens within it are highlighted. Two columns →
at most two records.)

**v1 needs no column/side discriminator field for this.** `type` already implies
the column — `a` is inherently the new/right side, `d` the old/left side — and a
context line's two halves are the *same* logical line, so the host tells the two
`c` columns apart by position (which it does anyway). A side field would only earn
its keep to disambiguate the symmetric `c` case, and that case needs no
disambiguation.

### 6.3 Wrapping — emit on every output row

> **The pager emits a line's record at the start of *every output row* it produces
> for that line, including its own wrapped continuations.**

There are two distinct kinds of wrapping, and the rule differs:

- **Terminal/host wrapping** — the pager emits *one* line (one `\n`) and the
  terminal (or the host's view) wraps it onto several visual rows. Here only the
  primary row needs a record; the host's own row→line mapping routes every visual
  row of that line back to it. A pager that relies on terminal wrapping emits one
  record and is fine.
- **Pager wrapping** — the pager itself emits *several* lines (several `\n`s) for
  one logical line, as difftastic's side-by-side does and as delta does in
  side-by-side with `wrap-max-lines`. Now each wrapped row is a **distinct line**
  to the host, with nothing tying row *N+1* back to row *N* — so each must carry
  its own record, or it has none.

A continuation row re-emits the *same* record as the primary line it continues
(no line-number advance). A column that has run out of wrapped content carries no
record on its padding rows.

Getting this wrong is not theoretical: without per-row records, acting on a
wrapped continuation row does nothing, and hunk/file navigation breaks because the
untagged rows fragment a wrapped line into one block per visual row. Both delta and
difftastic had this bug in their prototypes and were fixed to emit per-row.

### 6.4 Header rows

Hunk headers (`@@`), file headers (`diff --git`, `---`, `+++`), and any other
non-content rows carry **no** record in v1 (the `h`/`f`/`o` types are reserved but
unused). A host acting on a header row falls back to whatever non-metadata path it
has, or to no action.

---

## 7. How the host consumes it (informative)

This section is not normative — it describes how lazygit's prototype consumes the
protocol, to make the emit rules concrete.

- The host attaches each record to the following cell, like an OSC-8 hyperlink. If
  a record has no following cell (a genuinely empty rendered line), the host adds
  a content-less cell to hold it.
- **Row-granular action** (e.g. a keyboard "act on this line"): use the **first**
  record on the row. In side-by-side this is the left column — fine, since the two
  sides of a change are one hunk for staging purposes.
- **Point-granular action** (a mouse click): use the **nearest record at or to the
  left of the click column** → lands in the column actually clicked.
- The host normalizes `file` (resolving it relative to the repository working
  tree) and otherwise treats the record as opaque identity.

---

## 8. Known limitations and v2 candidates — feedback wanted

These are surfaced honestly because the protocol is a draft. None blocks v1; each
is a place where pager-author input would shape v2.

### 8.1 `c` and `a` records carry no old-file number

For context and addition records, `old-line` is empty. This is invisible in a
single-column rendering, but **side-by-side makes it visible**: the old/left
column shows a real old-file line number on screen, yet the record for that cell
reports only `new-line`. With difftastic it bites hardest — difftastic *always*
shows old-file numbers, old ≠ new is the norm, and the offset is **not constant**
(it follows structural alignment, not a per-hunk delta), so the old number is not
even derivable from the new one.

Nothing in the current host consumers needs the old number for `c`/`a`, so v1
leaves it empty. The clean v2 fix, if needed, is **carry both numbers on every
record** (applies to both single-column and side-by-side) — *not* a side-only tag.

### 8.2 The token-vs-line model mismatch (difftastic)

The unified-diff pagers (git, delta, diff-so-fancy) derive from git's
**line-granular** patch, where a modified line is by construction a `-` line plus a
`+` line. The `c`/`a`/`d` type set fits that model exactly.

difftastic is **token-granular**: it aligns an old line with a new line and marks
novelty *per token*. That produces aligned rows the line model has no clean slot
for. Real example — `let x=1; println!("{}", x);` changed to
`let x=2; let y=3; println!("{}", x + y);`:

```
…;c;4;;src/lib.rs 3     println!("{}", x);      …;a;4;;src/lib.rs 4     println!("{}", x + y);
```

The old line `println!("{}", x);` has no novel tokens (they all survive into the
new line), while the new line added `+ y`. So difftastic's faithful per-cell
verdict is **`c` on the left, `a` on the right** — the same aligned row carrying a
context record and an addition record, with **no `d`**, because by difftastic's
model nothing was deleted.

A host mapping cells to git's line-granular patch then resolves the right/added
cell correctly, but the left/old cell resolves as *context at the new line*, not
as git's `-` line — its old-side deletion identity is not recoverable from the
record. Practical impact is small (users click the changed side; an editor target
on the left still opens the right new-file line), but it is a genuine semantic gap.

Options, none taken in v1, all wanted-as-feedback:

- **Host-side:** treat an old-column cell aligned with a novel new line as the `-`
  side of a modification (the host knows it is the left column).
- **Emitter-side, v1-compatible:** difftastic could emit `d` for the old side of
  *any* aligned changed row — but that re-imposes the very line-granular model
  difftastic exists to escape, discarding its more precise "this was not removed"
  judgement. Probably wrong to force.
- **A `modified` / `m` type (v2):** names the case directly — "aligned, changed,
  present on both sides" — but splits the clean `c`/`a`/`d` staging mapping and
  needs every pager to agree when to use it. Recorded as a v2 candidate.

### 8.3 Pure-deletion `new-line` in token tools is synthesized

A line-granular pager always knows the new-file position a deletion sits at. A
token tool without a linear new-file counter (difftastic) computes it from the
previous aligned new line (`prev_rhs + 1`). This is exact for the common case but
can drift across hunk boundaries or with zero context lines. The deletion's
*old*-line (its real identity for staging) is always exact; only the editor-target
`new-line` is approximate. Pagers should emit the most precise `new-line` they can.

### 8.4 Out of scope for v1

`git --word-diff` (inline `[-…-]{+…+}` markup, no per-line `+`/`-`) has no clean
per-line identity and is not addressed.

---

## 9. Where feedback is most wanted

1. **The OSC number, `1717`.** Chosen after auditing the OSC allocations of
   xterm, VTE, kitty, foot, WezTerm, iTerm2, Windows Terminal, Ghostty, VS Code,
   ConEmu and urxvt (see the appendix): `1717` is unused by all of them and sits
   in the large empty 1400–5000 band (only iTerm2's `1337` is nearby). There is no
   central registry, so this is "verified unused across the terminals that matter,"
   not "allocated." If you know of a terminal that interprets `1717`, please say so.
2. **The env-var name and grammar** (`EMIT_OSC1717_METADATA=V1,…`).
3. **The both-numbers question** (§8.1) — is per-column-exact old/new numbering
   worth a v2, or is it never needed?
4. **The token-vs-line mismatch** (§8.2) — should there be an `m` type, or is
   host-side inference the right home for it?
5. **Can your pager actually produce all four fields per region?** In particular
   the side for deleted lines, and in side-by-side mode. (delta needed to track
   its own old/new counters because its line-number counters are dormant unless
   `--line-numbers` is on; difftastic had them natively. Your mileage may vary.)
6. **Optional trailing fields within a version** — a cheap, additive escape valve
   short of a version bump: a consumer stops at the fields it knows. Is this worth
   blessing in v1, or should all growth go through the version field?

---

## 10. Reference implementations (prototype)

All three are at prototype quality and emit the v1 format described here (modulo
the placeholder OSC number, which these used `456` for before `1717` was chosen):

- **delta** — a dedicated additive emitter that injects only OSC bytes (no change
  to styling, width, or wrapping); with the env var unset, output is byte-for-byte
  identical to stock delta. Covers unified and side-by-side modes, including
  wrapped rows.
- **difftastic** — the categorical case (#1 host-side parsing cannot serve it in
  either mode). Emits the same v1 format under the same handshake; markedly less
  code than delta because difftastic carries old/new line numbers natively. Covers
  side-by-side and inline modes.
- **host carrier** — gocui (lazygit's TUI library) accumulates the OSC number,
  collects the payload, and stamps it per-cell like a hyperlink, cleared at each
  line boundary so it cannot bleed onto an untagged following line.

A key validated property in both pagers: **with the handshake absent, output is
byte-identical to the unpatched pager** — the protocol is strictly additive.

---

## Appendix — terminal OSC audit behind the `1717` choice

There is no registry for OSC numbers; the de-facto convention is to pick a high,
distinctive number and verify no real terminal acts on it (an unknown OSC is
*skipped* by conformant terminals, but a *recognized* one could fire a visible
side-effect, e.g. `OSC 555` flashes foot and `OSC 777` raises a desktop
notification). The numbers below are interpreted by at least one surveyed terminal
and were therefore excluded:

| OSC | meaning | actors |
|---|---|---|
| 0–3 | title / icon / X11 property | all |
| 4, 5, 6 | palette / special color / tab color | xterm, iTerm2 |
| 7, 8 | working directory, hyperlink | all major |
| 9 | notifications; ConEmu progress `9;4`, cwd `9;9` | iTerm2, kitty, WinTerm, foot |
| 10–19 | dynamic colors (fg/bg/cursor/mouse/Tek/highlight) | xterm, foot, VTE |
| 21, 22 | color query/set, pointer shape | kitty, foot |
| 46, 50, 51, 52 | logfile, font, Emacs, clipboard | xterm, iTerm2; 52 all major |
| 66 | text-sizing protocol | kitty, foot |
| 99 | desktop notifications | kitty |
| 104–106, 110–119 | reset colors | xterm, foot, VTE |
| 133 | semantic prompt / shell integration | iTerm2, kitty, foot, WezTerm, WinTerm |
| 176 | application id | foot |
| 555 | flash the terminal | foot |
| 633 | shell integration | VS Code |
| 777 | desktop notification (rxvt ext) | urxvt, VTE, foot, WezTerm |
| 1337 | proprietary namespace / file transfer | iTerm2, WezTerm |
| 5522 | extended clipboard | kitty |
| 30001 / 30101 | color-stack push / pop | kitty |

`1717` collides with none of these and is far from every active cluster.

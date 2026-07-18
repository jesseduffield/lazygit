# Diff Line Metadata over OSC 1717 — draft specification (v1)

**Status: draft, for feedback.** This document describes a small terminal
escape-sequence protocol by which a diff renderer (delta, difftastic, diff-so-fancy,
…) annotates each rendered line of a diff with the patch-space identity it
represents, so that a host program rendering the diff renderer's output can map a screen
row (and column) back to *the exact line in the underlying diff*.

It is published to gather feedback from diff renderer authors before anything is
finalized. The wire format, the negotiation handshake, and the OSC number are all
open to revision — §9 lists the points where feedback is most wanted.

The protocol grew out of [lazygit](https://github.com/jesseduffield/lazygit), but
nothing in it is lazygit-specific; "the host" below means any program that runs a
diff renderer and consumes its output.

About terminology: Lazygit has been using the term "pager" for what we call
"diff renderer" here. This is incorrect, a pager is something like less; but
it's unlikely to change soon, so be aware that the terms "pager" (or "custom
pager") and "diff renderer" can be used interchangeably in some discussions
about this. For the rest of this document, we avoid the term "pager" though.

---

## 1. Motivation — what this enables, and why parsing isn't enough

A host that shows a diff rendered by a diff renderer often wants to act on the line the
user is pointing at:

- **dive into staging / patch-building** for that hunk line,
- **open an editor** at that file and line,
- **open that line in a code-review / PR web view** (needs the side — old vs new),
- **navigate by hunk or by file** within the rendered diff,
- **preserve the scroll position and selection** across a re-render (the diff is
  re-rendered with a different context size, or a different renderer, and the host
  wants to keep the user anchored on the same patch line).

Every one of these needs the same primitive: **given a rendered row, recover
`(file, side, line)`** — the precise line of the unified diff that row stands for.

For *structure-preserving* renderings (no diff renderer, `git diff --color`, or
`delta --color-only` without line numbers) the host can recover this by parsing
the on-screen text: walk up to the nearest `@@` and `diff --git`, count `+`/` `
lines, read the leading `+`/`-`/space. That works and needs no cooperation from
the renderer.

But the moment a renderer **restructures** the diff, the unified-diff structure the
parse relies on is gone:

- `delta` (default mode) and `diff-so-fancy` drop or hide the `+`/`-` markers and
  convey the side with color, leaving no parseable unified-diff structure;
- `difftastic` is token-granular and side-by-side — there is no unified-diff line
  structure left in *either* of its modes.

In all of these the **diff renderer is the only component that still knows** which file,
side, and line each rendered cell belongs to — it computed exactly that to render
the diff. This protocol asks the renderer to *state* that knowledge inline, in a form
the host can read back and that is harmless everywhere else.

---

## 2. Design at a glance

1. The diff renderer emits one OSC sequence carrying
   `(version, type, new-line, old-line, file)` **immediately before** each
   rendered region (a region is "one source line's worth of content in one
   column" — see §6).
2. The host attaches each record to **the cell that follows it**, exactly the way
   OSC-8 hyperlinks attach to the cells they precede. This makes the metadata
   survive terminal wrapping and multi-column layouts **without the host ever
   reasoning about layout** — it just reads the nearest preceding attachment.
3. The whole thing is gated behind an **environment-variable handshake**, so a
   renderer run outside a participating host (in a raw terminal, `less`, `tmux`, a
   CI log) emits nothing and behaves byte-for-byte as before.

The protocol is **layout-agnostic on the host side by construction**: all layout
knowledge (where a column starts, where a gutter ends, how a long line wraps)
stays in the diff renderer, which is the only component that has it.

---

## 3. Negotiation handshake

```
OSC1717_METADATA = V1[,V2,…]
```

- The **host** sets this environment variable on the diff renderer subprocess to the list
  of protocol versions it understands, highest-preferred first is *not* required —
  the list is a set.
- The **diff renderer** emits the **highest version present in both** its own supported
  set and the advertised set (see section 4.4 for the format of the handshake
  record it emits). If the variable is unset, empty, or shares no version with
  the renderer, the renderer **emits nothing** and its output is unchanged.

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
| `type` | always | one character — see §5.1. v1 emits `c` (context), `a` (added), `d` (deleted), `f` (file header), `h` (hunk header). |
| `new-line` | all types except `f` | new-file line number, in the **diff's new-file space** (see §5.2). **Never present on `f`** — see §5.5. |
| `old-line` | only `type=d` | old-file line number. **Empty** for all other types. |
| `file` | always | the file path the line belongs to; absolute or repo-root-relative (the host normalizes — emit whichever is convenient). Carried on **every** record so a single record is a complete answer. |

For a **renamed** file, `file` is the **new** path (git's `+++ b/…` side); the old
path is not carried. A pure rename with no content change emits only its `f`
record (it has header rows but no content lines).

### 4.3 Examples

| rendered line | emitted record (`;`-form) |
|---|---|
| context, new line 10 | `1717;1;c;10;;src/foo.go` |
| addition, new line 11 | `1717;1;a;11;;src/foo.go` |
| deletion, old line 9, sits at new pos 11 | `1717;1;d;11;9;src/foo.go` |
| two consecutive deletions | `…;d;11;9;…` then `…;d;11;10;…` (same `new-line`, different `old-line` — see §5.3) |
| modification rendered as one row | `…;d;11;9;…` then `…;a;11;;…` back-to-back (both halves of the row — see §6.2) |
| whole-file deletion | `1717;1;d;0;9;old/path` (`new-line` 0 — see §5.4) |
| file header | `1717;1;f;;;src/foo.go` (never a line number — see §5.5) |
| hunk header, hunk starting at new line 10 | `1717;1;h;10;;src/foo.go` |

### 4.4 The handshake record

A conforming diff renderer emits, as the **very first thing it writes** and **once per run**,
a **version-only** record naming the version it negotiated:

```
ESC ] 1717 ; <version> ST
```

i.e. the OSC introducer and the version field **with no further fields** —
`\x1b]1717;1\x1b\` for v1. It is emitted whenever the handshake (§3) negotiates a
version, *before* any diff content (and before the first per-line record).

Its purpose is to let the host **probe** a diff renderer cheaply and definitively: run it on an
**empty diff** (no changed content) and look for this record. Without it, "does this
renderer speak the protocol?" could only be inferred from the per-line records — but a
diff with no content lines (a binary file, or the empty diff a probe would use) emits
none, so the absence of records would be indistinguishable from an unsupported renderer.
The handshake is **content-independent** (it precedes, and does not depend on, any
diff), so a single probe is conclusive and a binary file can't be mistaken for an
unsupported renderer. It also tells the host the negotiated version up front.

A host distinguishes it from a per-line record (§4.1) by **field count**: the handshake
carries only the version (no `;` after it); a per-line record always has the full five
fields. A host that doesn't care about probing may simply ignore any record it can't
parse as five fields — so the handshake is harmless to existing parsers.

---

## 5. Semantics

### 5.1 Type

`type` is one character. v1 defines five. The first three are the **content-line**
types, and a conforming diff renderer emits one before every content line it
renders:

- `c` — context (unchanged) line
- `a` — added line
- `d` — deleted line

These are **patch-space** types, and the renderer classifies in patch space —
not by its own display model. Concretely: `c` only for a line whose old and new
contents are byte-identical; a pair of aligned lines whose contents differ is a
`d` plus an `a`, and a line present on only one side is a `d` or `a` — however
the renderer paints them. This matters for renderers finer-grained than git's
line model (difftastic marks novelty per *token*, so a line changed only by
added tokens has nothing highlighted on its old side, and a blank line has no
tokens at all) — classifying by highlighting would tag such lines `c`, and a
host staging by the records would drop one half of the modification. See §8.

The other two are the **header** types, emitted before the structural rows:

- `f` — file header — the row(s) a renderer prints to announce a file. Carries
  **no line numbers**, ever — see §5.5.
- `h` — hunk header — the row(s) announcing a hunk (e.g. a reformatted
  `@@ … @@` line). Always carries `new-line` — the hunk's first line (§5.2).

The header types are not optional: a renderer emits `f` on each file-header row
and `h` on each hunk-header row it renders. A renderer that has no such rows
simply has none to tag, and a single row that announces both a file *and* a hunk
(difftastic's per-hunk banner) carries **both** records — see §5.5.

A host **must ignore a `type` it does not recognize** (treat the row as
non-actionable) rather than reject the record, so the set can grow later without a
version bump.

**`type` is load-bearing and cannot be inferred from the other fields.** `added`
and `context` both carry `{new-line present, old-line absent}`, so presence alone
cannot distinguish them — yet the host must (scroll preservation anchors
specifically on *change* lines). Hence an explicit type.

### 5.2 Line-number spaces

- `new-line` is the line number in the **diff's new-file space** — i.e. the
  new-file line numbering the diff itself uses, *not* necessarily the working-tree
  file (the diff may be against staged content). A host that opens an editor is
  expected to re-map this through its own diff↔worktree adjustment; the renderer
  should emit the number as it appears in the diff it is rendering.
- `old-line` is the old-file line number, present **only** for deletions.
- On an `h` record, `new-line` is the new-file line of the **first line of the
  hunk it heads** (for the hunk of a whole-file deletion that is `0`, matching
  §5.4), and `old-line` is empty. This is what lets "open in editor" on a hunk
  header land at the top of what the user is looking at.
- An `f` record carries no line numbers at all — both fields are empty (§5.5).

### 5.3 The deleted-line convention (both numbers)

A `d` record carries **both** numbers:

- `old-line` is the deletion's own old-file line number.
- `new-line` is the new-file position the deletion *sits at*: `newStart` plus the
  number of added/context lines above it within the hunk. This is exactly what
  `git`'s patch arithmetic already computes for a removed line.

Consequence, which all diff renderers must implement identically: **two consecutive
deletions share the same `new-line`** (nothing new-file-side advances between
them) and are told apart only by `old-line`. Example — two deletions at old lines
9 and 10, both sitting at new position 11:

```
1717;1;d;11;9;src/foo.go
1717;1;d;11;10;src/foo.go
```

A host uses `old-line` to find a deletion's patch line and its old-side
(review/PR `L`) anchor, and `new-line` for the editor target.

### 5.4 Whole-file add / delete

A deleted file's lines carry `new-line` = `0` (mirroring git's `@@ -1,N +0,0 @@`);
an added file's lines carry the new-file numbers normally and `type=a`.

### 5.5 Header records: `f` never carries a line number, `h` always does

The two header types deliberately differ in payload, and the difference is fixed
by the spec rather than left to the renderer:

- An `f` record's `new-line` and `old-line` are **always empty**. Its payload is
  the file path.
- An `h` record **always** carries `new-line` — the first line of the hunk it
  heads (§5.2).

**Why fixed per type, not "emit a line number if you have one"?** Because a field
the renderer *may* populate is a field the host can never rely on — every consumer
would need a fallback for the empty case anyway, and renderers would diverge on
what they emit. Making presence a function of the type alone keeps every record's
shape predictable from its first two fields.

**Why can't `f` require a line number?** A streaming renderer genuinely doesn't
have one at that point: delta draws its file header when it parses the `+++`
line, *before* it has seen the first `@@`, so the file's first hunk line is
unknown unless it buffered the header — abandoning its streaming design for one
field. And a file header doesn't need a line number: the actions it anchors
(open the file, jump to the file, list the files) are file-granular. A host that
wants a line anyway can scan forward to the file's first record that carries one.

**Why can `h` require one?** Every renderer knows a hunk's start line at the
moment it renders that hunk's header — it is right there in the `@@` line being
reformatted, or (for a renderer like difftastic that builds the whole diff before
rendering) in the hunk structure itself.

**Combined file+hunk headers emit both records.** Some renderers have no separate
file-header row: difftastic prints one banner per hunk (`path --- N/M --- lang`),
and the first hunk's banner is the only row announcing the file. The rule: a row
that announces both a file and a hunk carries **both** records — the `f` first,
then the `h` (the file is the outer structure). difftastic's first-hunk banner
thus carries an `f` and an `h`; its later banners carry only an `h`. A consequence
is that a row can carry several records even outside side-by-side mode, and not
all of them carry a line number — §7 lists what a consuming host should be aware
of.

Two useful invariants follow from all of this:

- every file in the diff has exactly one `f` (a multi-row header repeats it per
  row — §6.4 — but it is one logical record), and every hunk exactly one `h`, so
  file/hunk navigation and a "files in this diff" list fall directly out of the
  records;
- files with **no content lines** — pure renames, mode-only changes, binary files
  (§4.2) — still emit their `f`, so they stay visible to the identity layer:
  navigation can anchor on them and a file list includes them, even though there
  is no content to act on.

(Hosts that consume only content records can still recover the structure without
headers: the `file` field changes between consecutive content records at a file
boundary, and within a file a `new-line` jump of more than one marks a new hunk —
two consecutive deletions share a `new-line` by §5.3, so compute the gap from the
last *advancing* line. This remains valid, but it cannot see content-less files.)

---

## 6. Emit rules (placement)

The diff renderer first emits the handshake record (§4.4) — once, before any other output —
then a per-line record before each region as follows.

### 6.1 One record per region, at the region's start

The diff renderer emits each region's record at the **start of that region**. Everything
from there until the next region's record (or end of line) belongs to that
record — *including* any line-number gutter or other embellishment the renderer
considers part of the region. Where a region "really" starts is the renderer's
call, not the host's. The single firm requirement:

> **The record must precede the region's first cell**, so that a host searching
> leftward from any cell lands in the correct region.

A region can be **zero-width**: several records may precede the same cells, each
immediately following the last. This happens on a combined file+hunk header row
(§5.5, §6.4 — `f` then `h` before its first cell) and on a modification rendered
as a single row (§6.2 — `d` then `a`). A host must keep *every* record of a row,
not just the one nearest its cells — a record followed directly by another
record still names one of the row's identities (see §7).

### 6.2 Multiple regions per row (side-by-side)

A side-by-side row shows two source lines at once, so it carries **two records**,
one before each column:

- **left column → the old-side line** (`d`, carries `new-line` + `old-line`),
- **right column → the new-side line** (`a`, carries `new-line`),
- **a context line, shown in both columns → the same `c` record before each
  half.**

The blank counterpart of a pure add/delete (the empty half) carries **no** record.

**v1 needs no column/side discriminator field for this.** `type` already implies
the column — `a` is inherently the new/right side, `d` the old/left side — and a
context line's two halves are the *same* logical line, so the host tells the two
`c` columns apart by position (which it does anyway). A side field would only earn
its keep to disambiguate the symmetric `c` case, and that case needs no
disambiguation.

**A modification rendered as a single row carries both records.** A renderer may
show an aligned changed row as *one* column — difftastic collapses a hunk whose
changes are all one-sided, printing a modified row's content once. The row still
*represents* both patch lines, so it carries the pair anyway: the `d`, then the
`a`, emitted consecutively at the row's start (patch order — deletions first; the
`d`'s region is zero-width, §6.1). This is what lets a host act on the whole
change from the one row it can see — staging it stages both halves, exactly as
for the two-column rendering of the same row. A collapsed *context* row is one
logical line and carries its `c` **once**, not per absent column.

### 6.3 Wrapping — emit on every output row

> **The diff renderer emits a line's record at the start of *every output row* it produces
> for that line, including its own wrapped continuations.**

There are two distinct kinds of wrapping, and the rule differs:

- **Terminal/host wrapping** — the renderer emits *one* line (one `\n`) and the
  terminal (or the host's view) wraps it onto several visual rows. Here only the
  primary row needs a record; the host's own row→line mapping routes every visual
  row of that line back to it. A renderer that relies on terminal wrapping emits one
  record and is fine.
- **Diff renderer wrapping** — the renderer itself emits *several* lines (several `\n`s) for
  one logical line, as difftastic's side-by-side does and as delta does in
  side-by-side with `wrap-max-lines`. Now each wrapped row is a **distinct line**
  to the host, with nothing tying row *N+1* back to row *N* — so each must carry
  its own record, or it has none.

A continuation row re-emits the *same* record as the primary line it continues
(no line-number advance). A column that has run out of wrapped content carries no
record on its padding rows.

Getting this wrong is not theoretical: without per-row records, acting on a
wrapped continuation row does nothing, and hunk/file navigation breaks because the
untagged rows fragment a wrapped line into one block per visual row.

### 6.4 Header rows

File-header and hunk-header rows carry records too, and these are **mandatory**:
`f` on each file-header row, `h` on each hunk-header row (§5.1, §5.5). A row that
announces both a file and a hunk carries both, `f` before `h`.

Where a header spans **several rows** — delta boxes a file name in divider/name/
divider lines; a renderer might draw a rule under a hunk header — **every row of
that block carries the same record(s)**: all three box rows get the file's `f`, a
hunk header and its rule both get that hunk's `h`, exactly as a wrapped content
line re-emits its record on every row (§6.3). That leaves no dead rows in a
header block — the user can act anywhere on it and land in the same place. The
renderer has the final say over what counts as its header block and which rows it
tags; this is the recommended default, not a hard rule.

Rows that belong to no header — dividers between files, padding, other pure
decoration — carry no record, and the host treats every un-annotated row as
non-actionable.

---

## 7. How the host consumes it (informative)

This section is not normative — it sketches the access model the per-cell carrier
is designed for, to make the emit rules concrete.

- The host attaches each record to the following cell, like an OSC-8 hyperlink. If
  a record has no following cell of its own — a genuinely empty rendered line, or
  a zero-width region whose record is immediately followed by the next record
  (§6.1) — the host adds a content-less cell to hold it, so no record of a row is
  lost to the one after it.
- **Point-granular action** (a mouse click): the per-cell attachment lets a host
  use the **nearest record at or to the left of the click column**, landing in the
  column actually clicked. (A host may equally resolve clicks at row granularity;
  the carrier supports either.)
- **Row-granular action** (e.g. a keyboard "act on this line"): a row can carry
  **more than one** record, so the host picks. The spec doesn't prescribe which —
  the right choice depends on the action — but these are the situations to be
  aware of:
  - A **side-by-side** row carries one record per column (§6.2). The two records
    are one logical change, but their fields differ: the left `d` records of a
    change block all sit at the block's first new line (§5.3), so the right
    column's `a` record has the more precise `new-line`, while the left one is
    the only carrier of the old-side identity.
  - A **modification collapsed to a single row** (§6.2) carries the same `d`+`a`
    pair, just back-to-back at the row's start. A host that acts on every record
    of a row (the natural choice for staging) handles the two renderings
    identically without telling them apart.
  - A **combined file+hunk header** carries an `f` and an `h` (§5.5). Only the
    `h` has a line number, so a line-oriented action (open the editor at a line)
    wants the `h`, while a file-oriented one (open the file, build a file list)
    anchors on the `f`.
  - When several records precede the same cell (the combined header), a
    click resolved by "nearest at or left of the point" lands on the **last**
    record emitted before that cell — the `h`. A host that wants keyboard and
    mouse actions to agree on such rows should account for that.
- An `f` record has no line number. What an action does with that is the host's
  choice — e.g. "open in editor" can open the file without jumping anywhere
  (useful in itself), or scan forward for the file's first record that carries a
  line number.
- The host normalizes `file` (resolving it relative to the repository working
  tree) and otherwise treats the record as opaque identity.

---

## 8. The token-vs-line model mismatch — how v1 resolves it, and the residue

**The mismatch (difftastic, AST mode).** Our `c`/`a`/`d` set is git's
**line-granular** shape: a modified line is a `-` plus a `+`. difftastic, when it
parses the language (its **AST/token mode**), is finer — it aligns lines and
marks novelty per token. A line changed *only by added tokens* (e.g.
`println!("{}", x);` → `println!("{}", x + y);`) then has **no novelty on the old
side** — difftastic renders that old line with nothing highlighted (dimmed like
context, or not at all when it collapses the hunk to a single column) — and a
blank line has no tokens to be novel in the first place. (In difftastic's
line/Text fallback — e.g. an unparseable file — it diffs by line and the mismatch
doesn't arise.)

**The resolution: classify in patch space (§5.1).** An early draft let the
record follow the display's highlighting — the old line of such a modification
said `c`. That reads consistently *within* difftastic, but it hands the host a
lie in patch space, and the lie is load-bearing: a host staging by the records
then staged only the `a` half of the modification, inserting the new line while
keeping the old (or, mirrored — a line changed only by *removed* tokens — only
the `d` half, deleting the line without its replacement). So v1 requires the
renderer to compare the aligned lines' *contents*: a differing pair is a `d`
plus an `a` whatever the display highlights, and a row that is the only rendered
representative of both halves carries both records (§6.2). Genuinely identical
lines inside a hunk still compare equal and stay `c`.

**The residue: changes the renderer never renders.** The records can only
describe rows that exist:

- difftastic treats whitespace-only line changes as no change at all — a file
  whose only edit is reindentation, or a lone deleted blank line, reports "no
  syntactic changes" and renders nothing — so there is no row to carry a record.
- an aligned pair differing only in whitespace that falls in a *context region*
  of difftastic's inline mode is rendered on one side only (before-context shows
  the old file's text, after-context the new file's); its record stays `c`,
  because a one-sided `d` would hand the host a deletion whose `a` half no
  rendered row can ever carry.

Both are inherent to a renderer whose diff model is coarser than git's in the
whitespace dimension; a host that needs those changes falls back on the raw
diff. A `modified`/`m` type — "aligned, changed, present on both sides" — remains
a v2 candidate (§9): it would name a single-row modification in one record
rather than a back-to-back pair, but splits the clean `c`/`a`/`d` mapping, and
the pair already carries both identities.

---

## 9. Where feedback is most wanted

1. **The OSC number, `1717`.** Chosen after auditing the OSC allocations of
   xterm, VTE, kitty, foot, WezTerm, iTerm2, Windows Terminal, Ghostty, VS Code,
   ConEmu and urxvt (see the appendix): `1717` is unused by all of them and sits
   in the large empty 1400–5000 band (only iTerm2's `1337` is nearby). There is no
   central registry, so this is "verified unused across the terminals that matter,"
   not "allocated." If you know of a terminal that interprets `1717`, please say so.
2. **The env-var name and grammar** (`OSC1717_METADATA=V1,…`).
3. **The token-vs-line mismatch** (§8) — v1 resolves it by classifying in patch
   space (a content-differing aligned pair is `d`+`a`, emitted together when only
   one row renders it). Is the back-to-back pair right, or should a v2 `m` type
   name the single-row modification directly?
4. **Can your diff renderer actually produce all four fields per region?** In particular
   the side for deleted lines, and in side-by-side mode. (delta needed to track
   its own old/new counters because its line-number counters are dormant unless
   `--line-numbers` is on; difftastic had them natively. Your mileage may vary.)
5. **The header types' fixed payloads (§5.5).** `f` never carries a line number,
   `h` always does, and a combined file+hunk row emits both records. This shape
   came out of prototyping in delta and difftastic: delta streams and draws its
   file header before it has parsed the first `@@` (so `f` cannot promise a
   line), while every renderer knows a hunk's start line at its hunk header (so
   `h` can); difftastic's per-hunk banner is both headers at once (so combined
   rows emit both). Does this fit your renderer — do you have a header shape
   where `h`'s line number is *not* in hand, or a combined row the both-records
   rule doesn't cover?

---

## 10. Reference implementations (prototype)

Three diff renderer emitters and one host carrier, all at prototype quality, emit or
consume the v1 format described here, over OSC `1717`:

- **delta** — a dedicated additive emitter that injects only OSC bytes (no change
  to styling, width, or wrapping); with the env var unset, output is byte-for-byte
  identical to stock delta. Covers unified and side-by-side modes, including
  wrapped rows, and the multi-row file/hunk-header decorations (every row of a
  header block carries its `f`/`h` — §6.4; delta's `f` is the case that cannot
  carry a line number, §5.5).
- **difftastic** — the categorical case (#1 host-side parsing cannot serve it in
  either mode). Emits the same v1 format under the same handshake; markedly less
  code than delta because difftastic carries old/new line numbers natively. Covers
  side-by-side and inline modes, classifying in patch space by comparing the
  aligned lines' contents (§5.1/§8) — its collapsed single-column rows are where
  the back-to-back `d`+`a` of §6.2 comes from. Its per-hunk banner is the combined
  file+hunk header of §5.5: the first hunk's banner carries `f` and `h`, later
  banners `h`.
- **diff-so-fancy** — the same #2 case as delta's default (it strips the `+`/`-`
  markers and conveys the side by color), but a line-oriented Perl filter rather
  than a structured renderer, and the simplest of the three: unified single-column
  only (no side-by-side). Classifies each line by its leading `+`/`-` before its
  existing code strips that marker, tracking its own old/new counters seeded from
  each `@@` header (like delta, it has no native line numbers); its reformatted
  file rule carries `f` and its hunk-header line `h` (the `@@` line it reformats
  has the hunk's start in hand). Combined/merge
  diffs are not annotated. Because diff-so-fancy defensively strips terminal escape
  sequences from the content it renders, the record is *prepended* to the line
  rather than embedded in it.
- **host carrier** — gocui (lazygit's TUI library) accumulates the OSC number,
  collects the payload, and stamps it per-cell like a hyperlink, cleared at each
  line boundary so it cannot bleed onto an untagged following line. A record
  that no cell consumed — a zero-width region (§6.1), or a record right before
  the line end — is kept in a content-less carrier cell, so a row's records
  survive complete.

A key validated property in all three renderers: **with the handshake absent, output
is byte-identical to the unpatched renderer** — the protocol is strictly additive.

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

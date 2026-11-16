# Contributing to terminfo

`terminfo` reads the terminfo database in pure Go. `AGENTS.md` holds the full
rules for a coding agent. This file is the short form for a person.

## What a change must do

The system `infocmp` is the reference (D3 in [docs/PLAN.md](docs/PLAN.md)).
If the package decodes a value differently from `infocmp`, the package is
wrong. A fix comes with a test that fails before the fix.

Do not edit `capvals.go` by hand. `gen.go` writes it from the `Caps` file of
`ncurses` 6.6 (D4, D26). To move to a new release, change `ncursesVer` in
`gen.go` and run `go generate`.

`terminfo` makes no promise of backward compatibility (D21). If a change
renames or removes an exported identifier, name it in the pull request.

## Before you open a pull request

The tests read the terminfo database on your machine, and they run
`/usr/bin/infocmp`. On Debian and Ubuntu, install `ncurses-bin` and
`ncurses-term` first.

Run these in the repository root. Each must pass, `gofmt -l .` must print
nothing, and `golangci-lint` must report 0 issues:

```sh
gofmt -l .
go vet ./...
go test -count=1 ./...
golangci-lint run ./...
```

CI runs the same commands in GitHub Actions, with golangci-lint v2.14.0
(D25).

## Writing

Write documents, code comments, error messages and commit messages in plain
English. Use short sentences and the active voice. Do not use contractions,
semicolons or em dashes. Use "must" and "can", not "should" and "may". The
`simple-english` skill in `.agents/skills` holds the full rules.

## Agent skills

The repository carries two agent skills. A skill is a set of instructions
that a coding agent loads for a task. `simple-english` sets how prose is
written, and `go-pedantry` sets how Go is written.

`skills-lock.json` names the source of each skill. The `skills` command from
npm writes that file. It writes each skill into two folders. Codex and the
other agents read `.agents/skills/<name>`, and Claude Code reads
`.claude/skills/<name>`.

To install the skills again, or to update them, run these two commands in the
repository root:

```sh
npx skills@1.7.0 add AminBlg/SimpleEnglish --skill simple-english --agent codex claude-code --copy -y
npx skills@1.7.0 add oborchers/fractional-cto --skill go-pedantry --agent codex claude-code --copy -y
```

Keep `--copy`. Without it, the command writes `.claude/skills/<name>` as a
symbolic link. A Windows checkout writes a symbolic link as a small text
file, and Claude Code then loads no skill and says nothing.
`TestSkillsAreCopies` fails on a link, and it fails when the two folders
differ.

`.claude/settings.local.json` holds the Claude Code permissions of one
person. The root `.gitignore` ignores it.

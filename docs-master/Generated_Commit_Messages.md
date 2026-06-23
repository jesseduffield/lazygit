# Generated Commit Messages

Lazygit can run an external command to generate the text for a new commit message. This is useful if you want a local script, an LLM tool, or another formatter to inspect the staged diff and suggest a message.

Configure the command under `git.commit.messageGeneratorCommand`:

```yaml
git:
  commit:
    messageGeneratorCommand: ~/bin/generate-staged-commit-message.sh
```

When this option is set, the commit menu shows `Generate Commit Message`. Selecting it runs the configured command and puts the command's stdout into the commit message fields. Lazygit does not commit immediately, so you can review and edit the generated message before submitting.

## Command Contract

Lazygit runs the configured command from the Git project root. For example, if your config contains:

```yaml
git:
  commit:
    messageGeneratorCommand: ~/bin/generate-staged-commit-message.sh --style conventional
```

Lazygit runs it like this:

```sh
(cd /path/to/repo && ~/bin/generate-staged-commit-message.sh --style conventional)
```

The command should write only the commit message to stdout. If stdout contains a blank line, Lazygit treats the first paragraph as the commit summary and the rest as the commit description.

If the command exits with a non-zero status, Lazygit shows a notification with stderr. The current commit message is left unchanged.

## Simple Script

This example uses the staged file summary to build a basic message:

```sh
#!/usr/bin/env sh
set -eu

if git diff --cached --quiet --exit-code; then
  echo "no staged changes" >&2
  exit 1
fi

files=$(git diff --cached --name-only | sed -n '1,3p' | paste -sd ', ' -)
count=$(git diff --cached --name-only | wc -l | tr -d ' ')

if [ "$count" -eq 1 ]; then
  printf 'update %s\n' "$files"
else
  printf 'update %s files\n\n%s\n' "$count" "$files"
fi
```

Save it somewhere on your `PATH`, make it executable, and point Lazygit at it:

```sh
chmod +x ~/bin/generate-staged-commit-message.sh
```

```yaml
git:
  commit:
    messageGeneratorCommand: ~/bin/generate-staged-commit-message.sh
```

## Codex Example

This example asks `codex exec` to inspect the staged diff and output only a commit message. It checks that staged changes exist in the current Git repository and lets you override the model or add extra `codex exec` arguments with environment variables.

```sh
#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'USAGE'
Usage: generate-staged-commit-message.sh

Calls Codex to generate a commit message for the currently staged changes in
the current Git repository.

Environment:
  CODEX_MODEL       Optional model name passed to `codex exec -m`.
  CODEX_EXTRA_ARGS  Optional additional arguments appended to `codex exec`.
USAGE
}

die() {
  printf 'error: %s\n' "$*" >&2
  exit 1
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

[[ $# -eq 0 ]] || {
  usage
  exit 2
}

command -v git >/dev/null 2>&1 || die "git is not installed or not on PATH"
command -v codex >/dev/null 2>&1 || die "codex is not installed or not on PATH"

repo_root=$(git rev-parse --show-toplevel 2>/dev/null) ||
  die "not a Git repository: $PWD"

if git -C "$repo_root" diff --cached --quiet --exit-code; then
  die "no staged changes found in $repo_root"
fi

prompt=$(cat <<'PROMPT'
Generate a concise, high-quality Git commit message for the currently staged changes.

Inspect the staged diff with:

  git diff --cached --stat
  git diff --cached

Requirements:
- Base the message only on staged changes.
- Prefer a single-line subject under 72 characters.
- Add a short body if it's not trivial commit and it materially improves clarity.
- Output only the commit message. Do not wrap it in Markdown, quotes, or commentary.
- Use the conventional commits format:
```
<type>[(scope)]: <description>

[optional body]

Signed-off-by: Your Name <you@example.com>
```
PROMPT
)

codex_args=(
  exec
  --cd "$repo_root"
  --sandbox read-only
  --ephemeral
  --color never
)

if [[ -n "${CODEX_MODEL:-}" ]]; then
  codex_args+=(-m "$CODEX_MODEL")
fi

if [[ -n "${CODEX_EXTRA_ARGS:-}" ]]; then
  # shellcheck disable=SC2206
  extra_args=(${CODEX_EXTRA_ARGS})
  codex_args+=("${extra_args[@]}")
fi

exec codex "${codex_args[@]}" "$prompt"
```

Then configure Lazygit:

```yaml
git:
  commit:
    messageGeneratorCommand: ~/bin/generate-staged-commit-message.sh
```

## Tips

Keep the command non-interactive because Lazygit captures stdout and stderr. If the command needs authentication or confirmation, handle that outside Lazygit first.

Make the command fail when it cannot produce a useful message. Lazygit will show stderr and leave the current message intact.

Quote paths inside scripts. Repository paths can contain spaces.

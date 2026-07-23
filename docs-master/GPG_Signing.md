# GPG Signing

## Background

If you have `commit.gpgSign` (or `tag.gpgSign`) enabled in git, git will ask
gpg to sign your commits/tags, and gpg in turn may need to ask you for your
key's passphrase via its configured pinentry program. Historically, Lazygit
had to hand off the whole terminal to a subprocess whenever this could
happen, because a pinentry program draws its prompt directly onto the
terminal, which would otherwise corrupt Lazygit's own screen. This meant
leaving Lazygit's UI entirely (even if gpg-agent had your passphrase cached
and no prompt was actually needed) and returning to it once the operation
finished.

## In-app passphrase prompt

If you have GnuPG 2.1 or later, and your gpg-agent hasn't been hardened with
`no-allow-loopback-pinentry`, Lazygit signs commits and tags using gpg's
`--pinentry-mode=loopback` option. This makes gpg print a plain textual
passphrase prompt on its own input/output instead of invoking a pinentry
program at all. Lazygit detects that prompt the same way it already does for
SSH passphrase/password prompts, and answers it from its own popup ("Enter
GPG passphrase") — so you never have to leave Lazygit's UI, and if
gpg-agent already has your passphrase cached, nothing is shown at all and
signing happens silently, as before.

This is automatic and requires no configuration. It only falls back to the
older subprocess/terminal-handoff behavior when loopback mode isn't
available (e.g. GnuPG < 2.1, or `no-allow-loopback-pinentry` set).

## Seeing when a commit/tag will be signed

Since git applies `commit.gpgSign`/`tag.gpgSign` silently — without ever
including a flag for it on the command line — Lazygit makes this explicit by
adding `--gpg-sign` (for commits) or `--sign` (for tags) to the command shown
in the Command Log whenever the corresponding config option is enabled. This
is purely for visibility; it doesn't change what gets signed, since it's
just making git's own implicit config-driven behavior explicit.

## The `overrideGpg` config option

See the `git.overrideGpg` option in [Config.md](./Config.md) if you want
Lazygit to skip spawning a subprocess for GPG operations even in the cases
where loopback mode isn't available. This is only safe if your gpg-agent
caches your passphrase or uses a GUI pinentry (pinentry-gtk-2/pinentry-qt/
pinentry-mac); Lazygit ignores this setting (and still spawns a subprocess)
if it detects a terminal-based pinentry (pinentry-tty/pinentry-curses),
since honoring it in that case would corrupt Lazygit's UI.

package shared

import (
	. "github.com/jesseduffield/lazygit/pkg/integration/components"
)

// holdIndexLockScript is the part of a git hook that takes index.lock the
// first time the hook runs, and holds it until the git command that ran the
// hook has exited. This simulates another git process (for example an IDE
// running "git status" in the background) that takes the lock while a command
// of several steps, such as a rebase, is partway through.
const holdIndexLockScript = `if [ ! -e .git/index-lock-taken ]; then
	touch .git/index-lock-taken .git/index.lock
	# If lazygit runs the command in a pty (as it does for a pull), the
	# background process below gets SIGHUP when the command exits. This can
	# happen before the process has run its first command, so ignore SIGHUP
	# before starting it; it inherits this.
	trap '' HUP
	(
		while kill -0 $PPID 2>/dev/null; do sleep 0.01; done
		rm -f .git/index.lock
	) >/dev/null 2>&1 </dev/null &
fi
`

// HoldIndexLockAfterNextCommit installs a post-commit hook that takes
// index.lock after the next commit (see holdIndexLockScript).
//
// Install it after creating the commits of the test setup, since it acts on
// the first commit that is made after it.
func HoldIndexLockAfterNextCommit(shell *Shell) {
	shell.CreateFile(".git/hooks/post-commit", "#!/bin/sh\n"+holdIndexLockScript)
	shell.MakeExecutable(".git/hooks/post-commit")
}

// HoldIndexLockBeforeNextPickIsCommitted installs a reference-transaction hook
// that takes index.lock (see holdIndexLockScript) the next time a rebase or
// cherry-pick has applied the changes of a commit, but hasn't committed them
// yet. This is when git creates CHERRY_PICK_HEAD.
func HoldIndexLockBeforeNextPickIsCommitted(shell *Shell) {
	shell.CreateFile(".git/hooks/reference-transaction", `#!/bin/sh
[ "$1" = committed ] || exit 0
awk '$3 == "CHERRY_PICK_HEAD" && $2 !~ /^0+$/ { found = 1 } END { exit !found }' || exit 0
`+holdIndexLockScript)
	shell.MakeExecutable(".git/hooks/reference-transaction")
}

// HoldIndexLockAfterNextAppliedPatch installs a post-applypatch hook that takes
// index.lock (see holdIndexLockScript) after "git am" has committed the next
// patch.
func HoldIndexLockAfterNextAppliedPatch(shell *Shell) {
	shell.CreateFile(".git/hooks/post-applypatch", "#!/bin/sh\n"+holdIndexLockScript)
	shell.MakeExecutable(".git/hooks/post-applypatch")
}

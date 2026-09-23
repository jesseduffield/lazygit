# Working with stacked branches

When working on a large branch it can often be useful to break it down into
smaller pieces, and it can help to create separate branches for each independent
chunk of changes. For example, you could have one branch for preparatory
refactorings, one for backend changes, and one for frontend changes. Those
branches would then all be stacked onto each other.

Git has support for rebasing such a stack as a whole; you can enable it by
setting the git config `rebase.updateRefs` to true. If you then rebase the
topmost branch of the stack, the other ones in the stack will follow. This
includes interactive rebases, so for example amending a commit in the first
branch of the stack will "just work" in the sense that it keeps the other
branches properly stacked onto it.

Lazygit visualizes the individual branch heads in the stack by marking them with a
cyan asterisk (or a cyan branch symbol if you are using [nerd
fonts](Config.md#display-nerd-fonts-icons)).

When you push the topmost branch of the stack with `P`, and the branches below
it have commits that haven't been pushed yet, lazygit offers to push them along
with it. After rebasing the stack this saves you from checking out and
force-pushing every branch one by one; you are asked to confirm the force push
once for all of them. Only branches that already have an upstream are included.
Each of them is pushed to where `git push` would push it if it were checked out,
so your push configuration applies to them as usual.

When somebody else rebases the stack and force-pushes it, all your branches
show up as diverged, for example `↓5↑3`, even though the commits they are ahead
by are only the old versions of the ones that are now on the remote. Lazygit
tells this apart from a branch that carries work of your own, and shows the
divergence dimmed for such a branch. Pressing `f` on it resets it to its
upstream instead of refusing, so you don't have to check the branch out and pull
it. Lazygit only does this when every commit of the branch was on its remote
branch at some point. It finds that out from the reflog of the remote-tracking
branch. Reflogs are enabled by default, except in a bare repository; if you work
in one with linked worktrees, set `core.logAllRefUpdates` to true there to make
this work.

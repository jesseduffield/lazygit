# Working with stacked branches

When working on a large branch it can often be useful to break it down into
smaller pieces, and it can help to create separate branches for each independent
chunk of changes. For example, you could have one branch for preparatory
refactorings, one for backend changes, and one for frontend changes. Those
branches would then all be stacked onto each other.

Git has support for rebasing such a stack as a whole; you can enable it by
setting the git config `rebase.updateRfs` to true. If you then rebase the
topmost branch of the stack, the other ones in the stack will follow. This
includes interactive rebases, so for example amending a commit in the first
branch of the stack will "just work" in the sense that it keeps the other
branches properly stacked onto it.

Lazygit visualizes the invidual branch heads in the stack by marking them with a
cyan asterisk (or a cyan branch symbol if you are using [nerd
fonts](Config.md#display-nerd-fonts-icons)).
